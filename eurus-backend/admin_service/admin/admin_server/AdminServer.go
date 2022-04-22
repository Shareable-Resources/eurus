package admin

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"strconv"
)

type AdminServer struct {
	service_server.ServiceServer
	Config          *AdminServerConfig
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	corsMiddleware  response.CORSMiddleware
}

func NewAdminServer() *AdminServer {
	adminServer := new(AdminServer)
	adminServer.Config = NewAdminServerConfig()
	adminServer.ServerConfig = &adminServer.Config.ServerConfigBase
	return adminServer
}

func (me *AdminServer) LoadConfig(args *service.ServiceCommandLineArgs) error {
	var err error = me.ServiceServer.LoadConfig(args, me.Config)
	return err
}

func (me *AdminServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	if httpConfig == nil {
		httpConfig = me.ServerConfig
	}
	var err error
	me.HttpServer, err = network.NewServer(httpConfig)
	me.loginMiddleware.AuthClient = me.AuthClient

	err = me.HttpServer.Listen()
	return err
}

func (me *AdminServer) InitAll() {
	me.ServiceServer.InitAuth(me.authLoginHandler, nil)
}

func (me *AdminServer) authLoginHandler(authClient auth_base.IAuth) {

	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	wallet, acc, walletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.Itoa(int(me.Config.ServiceId)))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load Hdwallet error: ", err.Error())
		panic("Load Hdwallet error: " + err.Error())
	}

	me.Config.HdWalletPrivateKey, err = wallet.PrivateKeyHex(*acc)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load wallet key error: ", err.Error())
		panic("Load Hdwallet key error: " + err.Error())
	}
	me.Config.HdWalletAddress = walletAddr
	log.GetLogger(log.Name.Root).Infoln("Sign server Wallet address: ", walletAddr)

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to DB: ", err)
		panic("Unable to connect to DB: " + err.Error())
	}
	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic("Unable to Init EthClient: " + err.Error())
	}

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

	fmt.Printf("AdminServer start listening at %s:%d\r\n", me.ServerConfig.HttpServerIP, me.ServerConfig.HttpServerPort)
}
