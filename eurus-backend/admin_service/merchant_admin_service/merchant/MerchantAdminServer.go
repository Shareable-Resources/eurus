package merchant_admin

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type MerchantAdminServer struct {
	service_server.ServiceServer
	loginMiddleware        request.LoginMiddleware
	errorMiddleware        response.ErrorMiddleware
	corsMiddleware         response.CORSMiddleware
	processor              *MerchantAdminProcessor
	smartContractProcessor *MerchantAdminSCProcessor
	dbProcessor            *MerchantAdminDBProcessor
	config                 *MerchantAdminServerConfig
}

func NewMerchantAdminServer() *MerchantAdminServer {
	server := new(MerchantAdminServer)
	server.config = NewMerchantAdminServerConfig()
	server.ServerConfig = &server.config.ServerConfigBase

	return server
}

func (me *MerchantAdminServer) LoadConfig(args *service.ServiceCommandLineArgs) error {
	var err error = me.ServiceServer.LoadConfig(args, me.config)
	return err
}

func (me *MerchantAdminServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	if httpConfig == nil {
		httpConfig = me.ServerConfig
	}
	var err error
	me.HttpServer, err = network.NewServer(httpConfig)
	me.loginMiddleware.AuthClient = me.AuthClient
	me.setupRouter()

	return err
}

func (me *MerchantAdminServer) InitAll() {
	me.ServiceServer.InitAuth(me.authLoginHandler, nil)

}

func (me *MerchantAdminServer) authLoginHandler(authClient auth_base.IAuth) {

	_, err := me.QueryConfigServer(me.config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	// wallet, acc, walletAddr, err := secret.GenerateMintBurnKey(me.Config.PrivateKey, strconv.Itoa(int(me.Config.ServiceId)))
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Load Hdwallet error: ", err.Error())
	// 	panic("Load Hdwallet error: " + err.Error())
	// }

	// me.Config.HdWalletPrivateKey, err = wallet.PrivateKeyHex(*acc)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Load wallet key error: ", err.Error())
	// 	panic("Load Hdwallet key error: " + err.Error())
	// }
	// me.Config.HdWalletAddress = walletAddr
	// log.GetLogger(log.Name.Root).Infoln("Sign server Wallet address: ", walletAddr)

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to DB: ", err)
		panic("Unable to connect to DB: " + err.Error())
	}

	me.dbProcessor = NewMerchantAdminDBProcessor(me.DefaultDatabase)
	me.smartContractProcessor = NewMerchantAdminSCProcessor(me.config)
	me.processor = NewMerchantAdminProcessor(me.dbProcessor, me.smartContractProcessor, me.config)

	// _, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
	// 	panic("Unable to Init EthClient: " + err.Error())
	// }

	err = me.InitHttpServer(nil)
	if err != nil {
		panic(err)
	}

	go func() {
		err := me.HttpServerListen()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("MerchantAdminServer start listening at %s:%d\r\n", me.ServerConfig.HttpServerIP, me.ServerConfig.HttpServerPort)
}

func (me *MerchantAdminServer) setupRouter() {
	me.HttpServer.Router.Use(me.corsMiddleware.Handler)
	me.HttpServer.Router.Post(RootPath+DummyLoginPath, me.dummyLogin)
	me.HttpServer.Router.Post(RootPath+MerchantLoginPath, me.login)

	me.HttpServer.Router.Mount(RootPath, me.setupAuthRouter())

}

func (me *MerchantAdminServer) setupAuthRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Post(RefreshLoginToken, me.refreshLoginToken)
	r.Get(RefundRequestListPath, me.getRefundRequestList)
	r.Post(RefundRequestPath, me.processRefundRequest)
	return r
}

func (me *MerchantAdminServer) getRefundRequestList(writer http.ResponseWriter, req *http.Request) {

	reqObj := NewGetMerchantRefundRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	var err error
	_, err = reqObj.ParseQueryString(req.URL)
	if err != nil {
		res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, err.Error())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = me.processor.GetRefundRequestList(reqObj)
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *MerchantAdminServer) dummyLogin(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		if secret.Tag != "dev" {
			res = response.CreateErrorResponse(reqObj, foundation.MethodNotFound, "Method not found")
		} else {
			res = me.processor.DummyLogin(me.AuthClient, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *MerchantAdminServer) login(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(MerchantLoginRequest)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		res = me.processor.processLogin(me.AuthClient, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *MerchantAdminServer) processRefundRequest(writer http.ResponseWriter, req *http.Request) {
	requestIdStr := chi.URLParam(req, "RequestId")

	requestId, err := strconv.ParseUint(requestIdStr, 10, 64)
	var res response.IResponse
	if err != nil {
		reqObj := new(request.RequestBase)
		res = response.CreateErrorResponse(reqObj, foundation.BadRequest, "Invalid request Id")
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	reqObj := NewRefundRequest(requestId)

	res = api.RequestToModel(req, reqObj)
	if res == nil {
		res = me.processor.processRefundRequest(reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *MerchantAdminServer) refreshLoginToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(MerchantRefreshLoginTokenRequest)
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
		res = me.processor.ProcessRefreshLoginToken(me.AuthClient, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}
