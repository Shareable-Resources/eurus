package auth

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/env"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"eurus-backend/foundation/ws"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"time"

	"github.com/gorilla/websocket"
)

type AuthServer struct {
	ws.WebSocketServer
	dispatcher *ws.WebSocketMessageDispatcher
	DataSource *AuthDataSource
	Config     *AuthServerConfig
}

type pendingConnection struct {
	Conn        *websocket.Conn
	SessionId   int64
	ConnectTime time.Time
}

const (
	ApiAuth                     = "authenticate"
	ApiRequestLoginToken        = "requestLoginToken"
	ApiRequestPaymentLoginToken = "requestPaymentLoginToken"
	ApiVerifyLoginToken         = "verifyLoginToken"
	ApiRefreshLoginToken        = "refreshLoginToken"
	ApiRevokeLoginToken         = "revokeLoginToken"
	ApiVerifySign               = "verifySign"
)

func NewAuthServer() *AuthServer {
	authServer := new(AuthServer)
	authServer.dispatcher = ws.NewWebSocketMessageDispatcher(ApiAuth)
	authServer.WebSocketServer = *ws.NewWebSocketServer(authServer.dispatcher)
	authServer.Config = NewAuthServerConfig()
	authServer.ServerConfig = authServer.Config.GetServerConfigBase()
	authServer.registerRequestHandler()
	authServer.ActualServer = authServer
	authServer.DataSource = &AuthDataSource{DB: nil, ServiceInfo: make(map[int64]conf_api.AuthService), Config: authServer.Config}

	//debug
	// authServer.DataSource.ServiceInfo[0] = new(ServiceInfo)
	// authServer.DataSource.ServiceInfo[0].PublicKey = "MIGJAoGBAOR0tmvxWuG8PntMAUjZFEhvmbh0w8+kjOoHKrNk5Tz53YLr2Hvda/lreozY0vOZKfOnoCVWi+EEPzj2NPagpSbJxSrx8XeH71kVzG9bZ9VGKGvRBBi7SWIMeCNH8n/UHRj71sUuY5+yk5V4jzWWIhdxUkqH8wQrPbB3ZVzyvWK9AgMBAAE="

	// authServer.DataSource.ServiceInfo[1] = &ServiceInfo{PublicKey: "MIGJAoGBALJJy2p4lwXtRk7zS4qUeWoxUtgRzdNHMwOUQ78PRbMAO/O3phvF2ptGgeDmMCtvVpYejAfso9EgBmbNZdT1Q2CfMxi0UG2uDWLCfq8BxWLgwvd7liIMJT7eR8L/OWmimmY8kr79yG9zFf/wGfZNYun9YXHCvZB0Fp9Z3TSO1c6/AgMBAAE="}
	return authServer
}

func (me *AuthServer) registerRequestHandler() {
	me.dispatcher.AddMessageHandler(ApiAuth, &ws.RequestHandlerInfo{Handler: AuthenticateHandler})
	me.dispatcher.AddMessageHandler(ApiRequestLoginToken, &ws.RequestHandlerInfo{Handler: RequestLoginTokenHandler})
	me.dispatcher.AddMessageHandler(ApiVerifyLoginToken, &ws.RequestHandlerInfo{Handler: VerifyLoginTokenHandler})
	me.dispatcher.AddMessageHandler(ApiRefreshLoginToken, &ws.RequestHandlerInfo{Handler: RefreshLoginTokenHandler})
	me.dispatcher.AddMessageHandler(ApiRevokeLoginToken, &ws.RequestHandlerInfo{Handler: RevokeLoginTokenHandler})
	me.dispatcher.AddMessageHandler(ApiRequestPaymentLoginToken, &ws.RequestHandlerInfo{Handler: RequestNonRefreshableLoginTokenHandler})
	me.dispatcher.AddMessageHandler(ApiVerifySign, &ws.RequestHandlerInfo{Handler: VerifySignHandler})
}

func (me *AuthServer) InitHttpServer(httpConfig network.IHttpConfig) {
	if httpConfig == nil {
		httpConfig = me.ServerConfig
	}
	//var err error
	//me.HttpServer, err = network.NewServer(httpConfig)

	//me.setupRouter()
	//err = me.HttpServer.Listen()

}

type QueryConfigAuthInfoResponse struct {
	response.ResponseBase
	Data *conf_api.ConfigAuthInfo `json:"data"`
}

func (me *AuthServer) QueryConfigServer() {
	for {
		authServiceInfoReq := conf_api.NewQueryConfigAuthInfoRequest()
		authServiceInfoReq.ServiceId = me.ServerConfig.ServiceId

		resp := new(QueryConfigAuthInfoResponse)
		reqRes := api.NewRequestResponse(authServiceInfoReq, resp)

		_, err := me.SendConfigApiRequest(reqRes)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.GetReturnCode() < int64(foundation.Success) {
			log.GetLogger(log.Name.Root).Errorln(resp.GetMessage())
			time.Sleep(5 * time.Second)
			continue
		}

		var parseConfig interface{} = me.Config

		for parseConfig != nil {
			err = conf_api.ConfigMapListToServerConfig(resp.Data.ConfigData, parseConfig.(server.IServerConfig))
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
				panic(err)
			}
			parseConfig = parseConfig.(server.IServerConfig).GetParent()
		}

		for _, authInfo := range resp.Data.AuthData {
			me.DataSource.ServiceInfo[int64(authInfo.Id)] = authInfo
		}

		log.GetLogger(log.Name.Root).Infoln("Service info count: ", len(resp.Data.AuthData), " loaded")
		break
	}
}

func (me *AuthServer) InitDBFromConfig(config *server.ServerConfigBase) error {
	err := me.WebSocketServer.InitDBFromConfig(config)
	if err != nil {
		return err
	}

	me.DataSource.DB = me.DefaultDatabase
	me.DataSource.SlaveDB = me.SlaveDatabase
	return nil
}

func (me *AuthServer) LoadConfig(commandLineArgs *service.ServiceCommandLineArgs, config server.IServerConfig) error {

	err := me.ServerBase.LoadConfig(&commandLineArgs.CommandLineArguments, config, func(configPath string) ([]byte, error) {
		return service.LoadConfigFile(configPath, env.IsConfigEncrypted, env.IsDeleteConfigAfterUsed, commandLineArgs.PasswordServerPath, commandLineArgs.UDSPath)
	})
	if err != nil {
		return err
	}

	return secret.DecryptSensitiveConfig(config.GetServerConfigBase())
}
