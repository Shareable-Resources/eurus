package sign

import (
	"encoding/base64"
	"eurus-backend/auth_service/auth"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"eurus-backend/sign_service/sign_api"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type SignServer struct {
	service_server.ServiceServer
	Config          *SignServerConfig
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	corsMiddleware  response.CORSMiddleware
	isSuccess       chan bool
	// UserWalletOwnerNonce uint64
	InvokerNonce uint64
	Mutex        *sync.Mutex
}

func NewSignServer() *SignServer {
	signServer := new(SignServer)
	signServer.Config = NewSignServerConfig()
	signServer.ServerConfig = &signServer.Config.ServerConfigBase
	signServer.isSuccess = make(chan bool)
	signServer.Mutex = new(sync.Mutex)
	return signServer
}

func (me *SignServer) LoadConfig(args *service.ServiceCommandLineArgs) error {

	var err error = me.ServiceServer.LoadConfig(args, me.Config)
	if err == nil {
		val, err := secret.DecryptConfigValue(me.Config.CentralizedUserWalletMnemonicPhase)
		if err != nil {
			return errors.Wrap(err, "Failed to decrypt centralized user wallet mnemonic phase")
		}
		data, _ := base64.StdEncoding.DecodeString(val)
		me.Config.CentralizedUserWalletMnemonicPhase = string(data)

		val, err = secret.DecryptConfigValue(me.Config.InvokerPrivateKey)
		if err != nil {
			return errors.Wrap(err, "Failed ot decrypt Invoker private key")
		}
		data, _ = base64.StdEncoding.DecodeString(val)
		me.Config.InvokerPrivateKey = string(data)

		val, err = secret.DecryptConfigValue(me.Config.UserWalletOwnerPrivateKey)
		if err != nil {
			return errors.Wrap(err, "Failed ot decrypt User wallet owner private key")
		}

		data, _ = base64.StdEncoding.DecodeString(val)
		me.Config.UserWalletOwnerPrivateKey = string(data)

		priKeyObj, err := crypto.HexToECDSA(me.Config.InvokerPrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Invalid invoker private key: ", err.Error())
			return errors.Wrap(err, "Invalid invoker private key")
		}
		me.Config.invokerAddress = crypto.PubkeyToAddress(priKeyObj.PublicKey).Hex()

	}
	return err
}

func (me *SignServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	err := me.ServerBase.InitHttpServer(httpConfig)
	if err != nil {
		return err
	}
	me.loginMiddleware.AuthClient = me.AuthClient

	me.setupRouter()
	return err
}

func (me *SignServer) InitAll() {
	me.ServiceServer.InitAuth(me.processInit, nil)
	for {
		select {
		case <-me.isSuccess:
			fmt.Println("Success")
			break
		}
		break
	}

}

func (me *SignServer) processInit(authClient auth_base.IAuth) {

	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	if me.Config.SideChainGasLimit == 0 {
		me.Config.SideChainGasLimit = 10000000
	}

	// wallet, acc, walletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.Itoa(int(me.Config.ServiceId)))

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
	// fmt.Println("Sign server wallet address: " + walletAddr)
	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic("Unable to Init EthClient: " + err.Error())
	}
	fmt.Println("Invoker address: ", me.Config.invokerAddress)
	log.GetLogger(log.Name.Root).Infoln("Invoker address: ", me.Config.invokerAddress)
	err = me.InitHttpServer(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Chain ID: ", me.Config.EthClientChainID)

	go func() {
		err := me.HttpServer.Listen()
		if err != nil {
			panic(err)
		}
		fmt.Printf("SignServer start listening at %s:%d\r\n", me.ServerConfig.HttpServerIP, me.ServerConfig.HttpServerPort)
	}()

	me.isSuccess <- true
}

func (me *SignServer) setupRouter() {
	me.HttpServer.Router.Use(me.corsMiddleware.Handler)
	me.HttpServer.Router.Mount(sign_api.RootPath, me.signRouter())
}

func (me *SignServer) signRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyServiceLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Post(sign_api.SignTransactionPath, me.signTransaction)
	r.Get(sign_api.GetAddressPath, me.getWalletAddress)
	r.Get(sign_api.GetPendingNoncePath, me.getPendingNonce)
	r.Post(sign_api.CalibrateNoncePath, me.calibrateNonce)
	r.Post(sign_api.SignUserWalletTransPath, me.handleSignUserWalletTransaction)
	r.Get(sign_api.GetCentralizedUserMainnetAddressPath, me.handleGetCentralizedUserMainnetAddress)
	return r
}

func (me *SignServer) signTransaction(writer http.ResponseWriter, req *http.Request) {
	reqObj := sign_api.NewSignTransactionRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		_ = api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	actualAuthClient := me.AuthClient.(*auth.AuthClient)
	serviceId, err := actualAuthClient.GetServiceIdFromServerLoginToken(reqObj.LoginToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get service ID from login token. Error: ", err)
		res = response.CreateErrorResponse(reqObj, foundation.UnauthorizedAccess, "Invalid login token")
	}

	serviceReq := conf_api.NewGetServiceGroupIdRequest()
	serviceReq.ServiceId = serviceId
	serviceRes := new(conf_api.GetServiceGroupIdFullResponse)

	reqRes := api.NewRequestResponse(serviceReq, serviceRes)
	_, err = me.SendConfigApiRequest(reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get group id from config server: ", err)
		res = response.CreateErrorResponse(reqObj, foundation.InternalServerError, "Internal network error")
	}

	if serviceRes.ReturnCode == int64(foundation.Success) {
		if serviceRes.Data == int64(conf_api.ServiceGroupUser) {
			res = signTransaction(me, reqObj)
		} else {
			res = response.CreateErrorResponse(reqObj, foundation.UnauthorizedAccess, "Request is not from user server")
		}
	}
	_ = api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) getWalletAddress(writer http.ResponseWriter, req *http.Request) {
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	reqObj := sign_api.NewQueryAddressRequest(loginToken.GetToken())
	res := getWalletAddress(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) getPendingNonce(writer http.ResponseWriter, req *http.Request) {
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	reqObj := sign_api.NewGetPendingNonceRequest(loginToken.GetToken())
	res := getPendingNonce(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) calibrateNonce(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(sign_api.CalibrateNonceRequest)
	res := calibrateNonce(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) handleSignUserWalletTransaction(writer http.ResponseWriter, req *http.Request) {
	reqObj := sign_api.NewSignUserWalletTransactionRequest()
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = signUserWalletTransaction(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) handleGetCentralizedUserMainnetAddress(writer http.ResponseWriter, req *http.Request) {
	reqObj := sign_api.NewGetCentralizedUserMainnetAddressRequest()
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = getCentralizedUserMainnetAddress(me.Config, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *SignServer) Calibrate(req *server.ControlRequestMessage) (bool, string, error) {
	var output string
	if req.MethodName == "calibrate" {
		args := req.Data
		if len(args) == 0 {
			output = me.PrintUsage()
			return true, output, nil
		}

		walletKeyType, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			output = me.PrintUsage()
			return true, output, nil
		}
		nonce, err := calibrate(me, sign_api.WalletKeyType(walletKeyType))
		if err != nil {

			return true, output, err
		}
		output = fmt.Sprintln("nonce is : ", nonce)
		return true, output, nil
	} else if req.MethodName == "help" {
		output = "calibrate [0/1] - Calibrate nonce of user wallet owner key / invoker key\r\n"
	}
	return false, output, nil

}

func (me *SignServer) PrintUsage() string {
	var output string
	output = fmt.Sprintln("Usage: calibrate [wallet key type]")
	output += fmt.Sprintln("Wallet key type: ")
	output += fmt.Sprintln("0 - user wallet owner key")
	output += fmt.Sprintln("1 - invoker key")
	return output
}
