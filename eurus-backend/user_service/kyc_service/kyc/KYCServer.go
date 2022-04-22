package kyc

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/secret"
	"eurus-backend/service_base/service_server"
	kyc_const "eurus-backend/user_service/kyc_service/kyc/const"
	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
)

type KYCServer struct {
	service_server.ServiceServer
	Config          *KYCConfig
	scProcessor     *KYCSCProcessor
	LogFilePath     string
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	corsMiddleware  response.CORSMiddleware
	taskQueue       *network.MQConsumer
	initChannel     chan bool
}

func NewKYCServer() *KYCServer {
	obj := new(KYCServer)
	obj.Config = NewKYCConfig()
	obj.ServerConfig = &obj.Config.ServerConfigBase
	obj.initChannel = make(chan bool, 2)
	obj.taskQueue = new(network.MQConsumer)

	return obj
}

func (me *KYCServer) processInit(authClient auth_base.IAuth) {
	fmt.Println("1. func - processInit - start")
	var err error
	//1. starts getting config from config server (This steps require config server running)
	//Send Service Id to config server, gets the configs_maps related to this service id, then assign to me.Config
	_, err = me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}

	hdWallet, acc, addr, err := secret.GenerateServerHDWallet(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.ServiceId, 10))
	if err != nil {
		panic(err)
	}
	me.Config.HdWalletPrivateKey, err = hdWallet.PrivateKeyHex(*acc)
	if err != nil {
		panic(err)
	}

	val, err := secret.DecryptConfigValue(me.Config.AdminAESKey)
	if err != nil {
		panic("Unable to decrypt AdminAESKey: " + err.Error())
	}

	data, _ := base64.StdEncoding.DecodeString(val)
	me.Config.AdminAESKey = string(data)

	me.Config.HdWalletAddress = addr
	fmt.Println("KYC Wallet address: ", me.Config.HdWalletAddress)

	//2. Init Db From Config
	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic(err)
	}

	_, err = me.InitEthereumClient(me.Config.EthClientProtocol, me.Config.EthClientIP, me.Config.EthClientPort, int64(me.Config.EthClientChainID))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init ETH client: ", err)
		panic(err)
	}
	me.scProcessor = new(KYCSCProcessor)
	err = me.scProcessor.Init(me.EthClient, me.Config, log.GetLogger(log.Name.Root))
	if err != nil {
		panic(err)
	}

	//3. emit signal to initChannel for operation after me.ServiceServer.InitAuth(me.processInit)
	me.initChannel <- true
}

func (me *KYCServer) InitAll() error {
	// 1. Get Auth Server config
	var err = me.QueryAuthServerInfo()
	if err != nil {
		fmt.Println("Unable to get auth server IP and port")
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}
	// 2. Get Config Server config
	me.ServiceServer.InitAuth(me.processInit, nil)
	<-me.initChannel
	fmt.Println("2. func - processInit - receive finish signal")

	me.taskQueue.Logger = log.GetLogger(log.Name.Root)
	err = me.taskQueue.SubscribeTaskQueue(me.ServerConfig.GetMqUrl(), &kyc_const.TaskQueueMetaData, me.mqTaskReceived)
	if err != nil {
		panic("Unable to subscribe MQ: " + err.Error())
	}
	// 3. Init Http Server
	go func() {
		err = me.InitHttpServer(nil)
		if err != nil {
			panic(err)
		}
		err = me.HttpServer.Listen()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (me *KYCServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	err := me.ServerBase.InitHttpServer(httpConfig)
	if err != nil {
		return err
	}
	me.loginMiddleware.AuthClient = me.AuthClient
	fmt.Println(me.loginMiddleware.VerifyServiceLoginToken)
	me.setupRouter()
	return err
}

func (me *KYCServer) setupRouter() {
	fmt.Println("3. func - setupRouter - start")
	//General API
	me.HttpServer.Router.Use(me.corsMiddleware.Handler)
	me.HttpServer.Router.Post(kyc_const.RootPath+kyc_const.AdminServerPath+kyc_const.EndPoint.LoginAdminUser, me.loginAdminUser) //done
	//Admin API with login token
	routerAdminToken := chi.NewRouter()
	routerAdminToken.Use(me.loginMiddleware.VerifyUserLoginToken)
	routerAdminToken.Use(me.errorMiddleware.ErrorHandler)
	routerAdminToken.Get(kyc_const.EndPoint.GetKYCStatusList, me.getKYCStatusList)        //done
	routerAdminToken.Post(kyc_const.EndPoint.UpdateKYCStatus, me.updateKYCStatus)         //done
	routerAdminToken.Post(kyc_const.EndPoint.ResetKYCStatus, me.resetKYCStatus)           //done
	routerAdminToken.Post(kyc_const.EndPoint.CreateAdminUser, me.createAdminUser)         //done
	routerAdminToken.Post(kyc_const.EndPoint.ChangeAdminPassword, me.changeAdminPassword) //done
	routerAdminToken.Get(kyc_const.EndPoint.GetKYCStatusOfUser, me.getKYCStatusOfUser)    //done
	routerAdminToken.Post(kyc_const.EndPoint.RefreshToken, me.refreshToken)               //done
	me.HttpServer.Router.Mount(kyc_const.RootPath+kyc_const.AdminServerPath, routerAdminToken)
	//User API with service token
	routerUserToken := chi.NewRouter()
	routerUserToken.Use(me.loginMiddleware.VerifyServiceLoginToken)
	routerUserToken.Use(me.errorMiddleware.ErrorHandler)                              //Error Handler must be after other middleware
	routerUserToken.Get(kyc_const.EndPoint.GetKYCCountryList, me.getKYCCountryList)   //done
	routerUserToken.Post(kyc_const.EndPoint.CreateKYCStatus, me.createKYCStatus)      //done
	routerUserToken.Get(kyc_const.EndPoint.GetKYCStatusOfUser, me.getKYCStatusOfUser) //done
	routerUserToken.Post(kyc_const.EndPoint.SubmitKYCApproval, me.submitKYCApproval)  //done
	me.HttpServer.Router.Mount(kyc_const.RootPath+kyc_const.UserServerPath, routerUserToken)

	fmt.Println()
	fmt.Println("Routes registered------")
	fmt.Println()
	//Print registered route
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(me.HttpServer.Router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
	fmt.Println()
}

func (me *KYCServer) mqTaskReceived(message *amqp.Delivery, topic string, contentType string, content []byte) {

	if contentType == "multipart/form-data" {
		boundary, ok := message.Headers["Boundary"]
		if !ok {
			log.GetLogger(log.Name.Root).Errorln("Missing boundary header")
			message.Reject(false)
			return
		}
		contentReader := bytes.NewReader(content)
		reader := multipart.NewReader(contentReader, boundary.(string))
		err, canRetry := ProcessKYCImage(me, reader)
		if err != nil {
			message.Reject(canRetry)
		} else {
			message.Ack(false)
		}

	}
}

// Controller
func (me *KYCServer) getKYCCountryList(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestGetKYCCountryList()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = RequestGetKYCCountryList(me, &reqObj.RequestBase)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) getKYCStatusOfUser(writer http.ResponseWriter, req *http.Request) {
	var err error
	reqObj := kyc_model.NewRequestGetKYCStatusOfUser()
	userId := chi.URLParam(req, "userId")
	reqObj.UserId, err = strconv.ParseUint(userId, 10, 64)
	if err != nil {
		api.HttpWriteResponse(writer, reqObj, response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "User Id is not a number"))
		return
	}
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = RequestGetKYCStatusOfUser(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) getKYCStatusList(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestGetKYCStatusList()

	pageQuery := req.URL.Query().Get("page")
	pageSizeQuery := req.URL.Query().Get("page_size")
	statusIdQuery := req.URL.Query().Get("status")
	reqObj.Email = req.URL.Query().Get("email")
	reqObj.WalletAddress = req.URL.Query().Get("walletAddress")
	page, err := strconv.Atoi(pageQuery)
	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil {
		api.HttpWriteResponse(writer, reqObj, response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Page and page size is not a number"))
		return
	}

	if pageQuery == "" || pageSizeQuery == "" {
		api.HttpWriteResponse(writer, reqObj, response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Page and page size is required"))
		return
	}
	if statusIdQuery != "" {
		status, err := strconv.Atoi(statusIdQuery)
		if err != nil {
			api.HttpWriteResponse(writer, reqObj, response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Status is not a number"))
			return
		}
		statusId := kyc_model.KYCStatusType(status)
		reqObj.KYCStatus = &statusId
	}

	res := api.RequestToModel(req, reqObj)

	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = RequestGetKYCStatusList(me, reqObj, page, pageSize)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) createKYCStatus(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestCreateKYCStatus()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = RequestCreateKYCStatus(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) loginAdminUser(writer http.ResponseWriter, req *http.Request) {

	reqObj := kyc_model.NewRequestLoginAdminUser()

	res := api.RequestToModelNoLoginToken(req, reqObj)

	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
	}
	res = LoginAdminUser(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) submitKYCApproval(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestSubmitKYCApproval()

	res := api.RequestToModel(req, reqObj)

	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = SubmitKYCApproval(me, reqObj)

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) updateKYCStatus(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestUpdateKYCStatus()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	userIdStrJson := loginToken.GetUserId()
	adminUserObj := new(kyc_model.AdminUser)
	json.Unmarshal([]byte(userIdStrJson), adminUserObj)
	operatorId := adminUserObj.Username
	reqObj.OperatorId = database.NullString{sql.NullString{Valid: true, String: operatorId}}
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = UpdateKYCStatus(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}
func (me *KYCServer) resetKYCStatus(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestResetKYCStatus()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	userIdStrJson := loginToken.GetUserId()
	adminUserObj := new(kyc_model.AdminUser)
	json.Unmarshal([]byte(userIdStrJson), adminUserObj)
	reqObj.AdminUser.Username = adminUserObj.Username
	res = ResetKYCStatus(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) createAdminUser(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestCreateAdminUser()
	res := api.RequestToModel(req, reqObj)

	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = CreateAdminUser(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) changeAdminPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRequestChangeAdminPassword()
	res := api.RequestToModel(req, reqObj)

	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = ChangeAdminPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *KYCServer) refreshToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := kyc_model.NewRefreshTokenRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	if res == nil {
		res = RefreshToken(me.AuthClient, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}
