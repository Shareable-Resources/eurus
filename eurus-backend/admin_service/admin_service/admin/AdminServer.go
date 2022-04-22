package admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/elastic"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type AdminServer struct {
	service_server.ServiceServer
	Config             *AdminServerConfig
	loginMiddleware    request.LoginMiddleware
	errorMiddleware    response.ErrorMiddleware
	corsMiddleware     response.CORSMiddleware
	dbProcessor        *AdminDBProcessor
	roleDbProcessor    *AdminRoleDBProcessor
	accountDbProcessor *AdminAccountDBProcessor
	elasticSearch      *elastic.ElasticSearch
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

	return err
}

func (me *AdminServer) InitAll() {
	me.ServiceServer.InitAuth(me.processPostLogin, nil)
}

func (me *AdminServer) processPostLogin(authClient auth_base.IAuth) {

	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
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
	// log.GetLogger(log.Name.Root).Infoln("Sign server Wallet address: ", walletAddr)

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to DB: ", err)
		panic("Unable to connect to DB: " + err.Error())
	}

	me.dbProcessor = NewAdminDBProcessor(me.Config, me.DefaultDatabase)
	me.roleDbProcessor = NewAdminRoleDBProcessor(me.Config, me.DefaultDatabase)
	me.accountDbProcessor = NewAdminAccountDBProcessor(me.Config, me.DefaultDatabase)

	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic("Unable to Init EthClient: " + err.Error())
	}

	me.elasticSearch = elastic.NewElasticSearch(me.ServerConfig.ElasticSearchPath)

	err = me.InitHttpServer(nil)
	if err != nil {
		panic(err)
	}

	go func() {
		me.initRouter()
		err := me.HttpServerListen()
		if err != nil {
			panic(err)
		}
	}()

	// fmt.Printf("AdminServer start listening at %s:%d\r\n", me.ServerConfig.HttpServerIP, me.ServerConfig.HttpServerPort)
}

func (me *AdminServer) initRouter() {
	me.HttpServer.Router.Post(RootPath+LoginPath, me.adminLogin)
	me.HttpServer.Router.Mount(RootPath, me.setupAdminRouter())
}

func (me *AdminServer) setupAdminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Post(VerifyGAPath, me.verifyGA)
	r.Get(FeatureListPath, me.getFeatureList)

	//Role
	r.Get(RoleListPath, me.getRoleList)
	r.Get(RoleDetailPath, me.getRoleDetail)
	r.Post(RoleCreatePath, me.createRole)
	r.Post(RoleUpdatePath, me.updateRole)
	r.Delete(RoleDeletePath, me.deleteRole)

	//Account
	r.Get(AccountListPath, me.getAccountList)
	r.Post(AccountCreatePath, me.createAccount)
	r.Post(AccountEditPath, me.updateAccount)
	r.Delete(AccountDeletePath, me.deleteAccount)
	r.Post(AccountChangePassord, me.changePassword)
	return r
}

func (me *AdminServer) adminLogin(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(AdminLoginRequest)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		res = processAdminLogin(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) verifyGA(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(VerifyGARequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() != int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			forward, ok := req.Header["X-Forwarded-For"]
			if !ok || len(forward) == 0 {
				reqObj.LoginIp = req.RemoteAddr
			} else {
				reqObj.LoginIp = forward[0]
			}
			res = processVerifyGA(me, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) getFeatureList(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			res = processQueryFeatureList(me.dbProcessor, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) getRoleList(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.QueryRoleListRequest)
	reqObj.State = admin_common.RoleAll

	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			res = processQueryRoleList(me.roleDbProcessor, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) getRoleDetail(writer http.ResponseWriter, req *http.Request) {

	reqObj := new(admin_common.QueryRoleDetailRequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			roleIdStr := chi.URLParam(req, "roleId")
			roleId, err := strconv.ParseUint(roleIdStr, 10, 64)
			if err != nil {
				res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid role id")
			} else {
				reqObj.RoleId = roleId
				res = processQueryRoleDetail(me.roleDbProcessor, reqObj)
			}
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) createRole(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.CreateRoleRequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			res = processCreateRole(me.roleDbProcessor, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) updateRole(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.UpdateRoleRequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			roleIdStr := chi.URLParam(req, "roleId")
			roleId, err := strconv.ParseUint(roleIdStr, 10, 64)
			if err != nil {
				res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid role id")
			} else {
				reqObj.RoleId = roleId
				res = processUpdateRole(me.roleDbProcessor, reqObj)
			}
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) deleteRole(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.DeleteRoleRequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			roleIdStr := chi.URLParam(req, "roleId")
			roleId, err := strconv.ParseUint(roleIdStr, 10, 64)
			if err != nil {
				res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid role id")
			} else {
				reqObj.RoleId = roleId
				res = processDeleteRole(me.roleDbProcessor, reqObj)
			}
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) getAccountList(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.QueryAccountListRequest)
	reqObj.State = admin_common.AccountAll
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {

			res = processQueryAccountList(me.accountDbProcessor, reqObj)

		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) createAccount(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.CreateAccountRequest)

	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {

			res = processCreateAccount(me.accountDbProcessor, reqObj)

		}
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *AdminServer) updateAccount(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.UpdateAccountRequest)

	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			id := chi.URLParam(req, "adminId")
			var err error
			reqObj.AdminId, err = strconv.ParseUint(id, 10, 64)
			if err != nil {
				res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid admin id")
			} else {
				res = processUpdateAccount(me.accountDbProcessor, reqObj)
			}

		}
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *AdminServer) deleteAccount(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.DeleteAccountRequest)

	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			id := chi.URLParam(req, "adminId")
			var err error
			reqObj.AdminId, err = strconv.ParseUint(id, 10, 64)
			if err != nil {
				res = response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid admin id")
			} else {
				res = processDeleteAccount(me.accountDbProcessor, reqObj)
			}

		}
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *AdminServer) changePassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(admin_common.ChangePasswordRequest)

	res := api.RequestToModel(req, reqObj)
	if res == nil {
		if reqObj.LoginToken.GetTokenType() == int16(auth.NonRefreshableToken) {
			res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, "Invalid login token")
		} else {
			res = processChangePassword(me.accountDbProcessor, reqObj)
		}
	}
	api.HttpWriteResponse(writer, reqObj, res)

}
