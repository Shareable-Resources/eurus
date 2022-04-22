package conf

import (
	"encoding/json"
	"errors"
	"eurus-backend/auth_service/auth"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service_server"
	"strconv"
	"strings"

	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type ConfigServer struct {
	service_server.ServiceServer
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	mqPublisher     *network.MQPublisher
	Config          *ConfigServerConfig
}

func NewConfigServer() *ConfigServer {
	configServer := new(ConfigServer)
	configServer.Config = NewConfigServerConfig()
	configServer.ServerConfig = &configServer.Config.ServerConfigBase
	return configServer
}

func (me *ConfigServer) InitHDWallet() {
	wallet, account, hdWalletAddr, _ := secret.GenerateMintBurnKey(me.ServerConfig.MnemonicPhase, me.ServerConfig.PrivateKey, strconv.FormatInt(me.ServerConfig.GetServiceId(), 10))
	me.ServerConfig.HdWalletAddress = hdWalletAddr
	me.ServerConfig.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*account)

	err := DBUpdateAuthServiceWalletAddress(me.DefaultDatabase, uint64(me.ServerConfig.ServiceId), ethereum.ToLowerAddressString(me.ServerConfig.HdWalletAddress))
	if err != nil {
		log.GetLogger(log.Name.Root).Warn("Unable to UpdateAuthServiceWalletAddress: ", err.Error())
	}
}

func (me *ConfigServer) LoadConfigFromDB() error {

	req := &conf_api.QueryConfigRequest{}
	req.Id = -1
	result, err := GetConfigFromDB(me.SlaveDatabase, req)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load InternalSmartContractConfigAddress: ", err.Error())
		panic(err)
	}

	var parseConfig interface{} = me.ServerConfig
	for parseConfig != nil {
		err = conf_api.ConfigMapListToServerConfig(*result, parseConfig.(server.IServerConfig))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
			return err
		}
		parseConfig = parseConfig.(server.IServerConfig).GetParent()
	}

	if me.ServerConfig.InternalSCConfigAddress == "" {
		log.GetLogger(log.Name.Root).Error("Unable to load InternalSmartContractConfigAddress")
		return errors.New("Unable to load InternalSmartContractConfigAddress")
	}

	if me.ServerConfig.ExternalSCConfigAddress == "" {
		log.GetLogger(log.Name.Root).Error("Unable to load ExternalSCConfigAddress")
		return errors.New("Unable to load ExternalSCConfigAddress")
	}

	return nil
}

func (me *ConfigServer) LoadConfigServerConfigFromDB() error {
	var err error
	req := &conf_api.QueryConfigRequest{}
	req.Id = me.ServerConfig.ServiceId
	result, err := GetConfigFromDB(me.SlaveDatabase, req)
	var parseConfig interface{} = me.Config
	for parseConfig != nil {
		err = conf_api.ConfigMapListToServerConfig(*result, parseConfig.(server.IServerConfig))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
			return err
		}
		parseConfig = parseConfig.(server.IServerConfig).GetParent()
	}
	return err
}

func (me *ConfigServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	err := me.ServerBase.InitHttpServer(httpConfig)
	if err != nil {
		return err
	}
	me.loginMiddleware.AuthClient = me.AuthClient

	me.SetupRouter()
	return err
}

func (me *ConfigServer) SetupRouter() {
	me.HttpServer.Router.Post(RootPath+GetConfigAuthPath, me.getConfigAuth)
	me.HttpServer.Router.Get(RootPath+GetServerSettingPath, me.getServerSetting)
	me.HttpServer.Router.Mount(RootPath, me.ConfigRouter())
}

func (me *ConfigServer) ConfigRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(me.loginMiddleware.VerifyServiceLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)

	r.Post(GetConfigPath, me.getConfig)
	r.Post(UpdateConfigPath, me.updateConfig)
	r.Delete(DelConfigPath, me.delConfig)
	r.Delete(DelGroupConfigPath, me.delConfig)
	r.Get(QueryServiceGroupAddressPath, me.getServiceGroupAddress)
	r.Get(GetServiceGroupIdPath, me.getServiceGroupId)

	r.Get(GetSystemConfigPath, me.getSystemConfig)
	r.Post(AddOrUpdateSystemConfigPath, me.addOrUpdateSystemConfig)

	// r.Post(AddMintConfigPath, me.addMintInfo)
	// r.Post(AddBurnConfigPath, me.addBurnInfo)
	// r.Post(AddApprovalConfigPath, me.addApprovalInfo)
	// r.Post(AddUserConfigPath, me.addUserInfo)

	r.Post(GetAuthPath, me.getAuth)
	// r.Post(UpdateAuthPath, me.updateAuth)
	// r.Put(AddAuthPath, me.addAuth)
	// r.Delete(DelAuthPath, me.delAuth)
	r.Post(SetServerWalletAddressPath, me.setServerWalletAddress)

	r.Post(GetKeyPairPath, me.getKeyPair)

	// r.Post(SetMainnetWalletPath, me.setMainnetWalletAddress)
	// r.Post(SetInnetWalletPath, me.setInnetWalletAddress)
	r.Get(GetMainnetWalletPath, me.getMainnetWalletAddress)
	r.Get(GetInnetWalletPath, me.getInnetWalletAddress)

	r.Post(SetWithdrawFeeETH, me.setWithdrawFeeETH)

	r.Post(SetExchangeRate, me.setExchangeRate)

	// r.Post(SetAsset, me.setAsset)
	r.Get(GetAsset, me.getAsset)

	//coldwallet has no getter
	// r.Post(SetColdWalletPath, me.setColdWalletAddress)

	// r.Post(AddAssetPath, me.addCurrencyInfo)
	// r.Delete(DelAssetPath, me.delCurrencyInfo)
	r.Get(GetERC20SCPath, me.getCurrencySCAddr)
	r.Get(GetERC20NamePath, me.getCurrencyNameByAddr)

	// r.Post(AddPublicAssetPath, me.addCurrencyInfo)
	// r.Delete(DelPublicAssetPath, me.delCurrencyInfo)
	r.Get(GetPublicERC20SCPath, me.getCurrencySCAddr)
	r.Get(GetPublicERC20NamePath, me.getPublicCurrencyNameByAddr)

	// r.Get(QueryAssetConfigPath, me.getAssetConfig)
	r.Get(GetFaucetConfigPath, me.getFaucetConfig)

	r.Get(GetAssetSettingsPath, me.getAssetSettings)

	return r
}

// func (me *ConfigServer) setColdWalletAddress(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewSetColdWalletRequest()
// 	reqObj.ValidateField()

// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = SetColdWalletAddress(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)

// }

// func (me *ConfigServer) setMainnetWalletAddress(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewSetMainnetWalletRequest()
// 	reqObj.ValidateField()

// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = SetMainnetWalletAddress(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)

// }

func (me *ConfigServer) getMainnetWalletAddress(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryMainnetWalletRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetMainnetWalletAddress(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *ConfigServer) setWithdrawFeeETH(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewSetWithdrawFeeETH()
	res := RequestToModelHandler(req, reqObj)
	res = SetWithdrawFee(me, reqObj)

	api.HttpWriteResponse(writer, reqObj, res)

}
func (me *ConfigServer) setExchangeRate(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewSetExchangeRate()
	res := RequestToModelHandler(req, reqObj)
	res = SetExchangeAmount(me, reqObj)

	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *ConfigServer) setAsset(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewSetAssetRequest()
	res := RequestToModelHandler(req, reqObj)
	err := AddAssetToDB(me, reqObj)
	if err == nil {
		api.HttpWriteResponse(writer, reqObj, res)
	}

}

func (me *ConfigServer) getAsset(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewGetAssetRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetAssetListFromDB(me.SlaveDatabase, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

// func (me *ConfigServer) setInnetWalletAddress(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewSetInnetWalletRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = SetInnetWalletAddress(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)

// }

func (me *ConfigServer) getInnetWalletAddress(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryInnetWalletRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetInnetWalletAddress(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

// func (me *ConfigServer) addPublicCurrencyInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddPublicCurrencyRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = AddPublicCurrencyInfo(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)

// }

// func (me *ConfigServer) addCurrencyInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddCurrencyRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = AddCurrencyInfo(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) delPublicCurrencyInfo(writer http.ResponseWriter, req *http.Request) {
// 	asset := chi.URLParam(req, "asset")

// 	reqObj := conf_api.NewDelPublicCurrencyRequest(asset)
// 	res := RequestToModelHandler(req, reqObj)
// 	reqObj.Asset = asset
// 	if res == nil {
// 		res = DelPublicCurrencyInfo(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) delCurrencyInfo(writer http.ResponseWriter, req *http.Request) {
// 	asset := chi.URLParam(req, "asset")

// 	reqObj := conf_api.NewDelCurrencyRequest(asset)
// 	res := RequestToModelHandler(req, reqObj)
// 	reqObj.Asset = asset
// 	if res == nil {
// 		res = DelCurrencyInfo(me, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

func (me *ConfigServer) getCurrencySCAddr(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryCurrencySCAddrRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetCurrencySCAddr(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getPublicCurrencySCAddr(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryPublicCurrencySCAddrRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetPublicCurrencySCAddr(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getServerSetting(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryServerSettingRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetServerSetting(me.SlaveDatabase, reqObj)
	} else {
		log.GetLogger(log.Name.Root).Errorln(res.GetMessage(), " nonce: ", res.GetNonce())
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *ConfigServer) getConfig(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryConfigRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetConfig(me.SlaveDatabase, reqObj)
	}

	//check serverID from UserServer == serverID from AuthServer
	// 	reqObj.ServerId : loginToken from user server
	//	reqObj.GetLoginToken().GetUserId() ["serviceId"] : ILoginToken from Auth server

	if reqObj.GetLoginToken() == nil {
		errStr := fmt.Sprintf("%v", "reqObj.GetLoginToken() == nil")
		reqObj := request.RequestBase{}
		res = response.CreateErrorResponse(&reqObj, foundation.InternalServerError, errStr)
		api.HttpWriteResponse(writer, &reqObj, res)
		return
	}

	userIdMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(reqObj.GetLoginToken().GetUserId()), &userIdMap)
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		reqObj := request.RequestBase{}
		res = response.CreateErrorResponse(&reqObj, foundation.InternalServerError, errStr)
		api.HttpWriteResponse(writer, &reqObj, res)
		return
	}

	iLoginTokenFromAuth, ok := userIdMap["serviceId"].(float64)
	if !ok {
		errStr := `Login token is not a server login token`
		reqObj := request.RequestBase{}
		res = response.CreateErrorResponse(&reqObj, foundation.LoginTokenInvalid, errStr)
		api.HttpWriteResponse(writer, &reqObj, res)
		return
	}
	if float64(reqObj.Id) != iLoginTokenFromAuth {
		errStr := `Request service id not match with login token`
		reqObj := request.RequestBase{}
		res = response.CreateErrorResponse(&reqObj, foundation.ServerIdUnmatch, errStr)
		api.HttpWriteResponse(writer, &reqObj, res)
		return
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getConfigAuth(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryConfigAuthInfoRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetConfigAuth(me.SlaveDatabase, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) updateConfig(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewAddConfigRequest()

	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = UpdateConfig(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) delConfig(writer http.ResponseWriter, req *http.Request) {
	serverId := chi.URLParam(req, "serverId")
	realServerId, err := strconv.Atoi(serverId)
	var resBase response.ResponseBase
	isGroup := strings.Contains(req.URL.Path, "config/group/")
	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		reqObj := request.RequestBase{}
		resBase = *response.CreateErrorResponse(&reqObj, foundation.InternalServerError, errStr)
		api.HttpWriteResponse(writer, &reqObj, &resBase)
		return
	}

	key := chi.URLParam(req, "key")
	reqObj := &conf_api.DeleteConfigRequest{IsGroup: isGroup}
	res := RequestToModelHandler(req, reqObj)
	reqObj.Id = int64(realServerId)
	reqObj.Key = key
	if res == nil {
		res = DelConfig(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *ConfigServer) getKeyPair(writer http.ResponseWriter, req *http.Request) {
	reqObj := &request.RequestBase{}

	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetKeyPair(reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getSystemConfig(writer http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	reqObj := conf_api.NewGetSystemConfigRequest(key)
	res := RequestToModelHandler(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	reqObj.Key = key
	res = GetSystemConfig(me.SlaveDatabase, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) addOrUpdateSystemConfig(writer http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	reqObj := conf_api.NewAddOrUpdateSystemConfigRequest(key)
	res := RequestToModelHandler(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	reqObj.Key = key
	res = AddOrUpdateSystemConfig(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) addAuth(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewAddAuthRequest()

	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = AddAuth(me.DefaultDatabase, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

// func (me *ConfigServer) updateAuth(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewUpdateAuthRequest()

// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = UpdateAuth(me.DefaultDatabase, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) delAuth(writer http.ResponseWriter, req *http.Request) {
// 	id := chi.URLParam(req, "id")
// 	realID, err := strconv.Atoi(id)
// 	var res *response.ResponseBase

// 	if err != nil {
// 		errStr := fmt.Sprintf("%v", err)
// 		reqObj := request.RequestBase{}
// 		res = response.CreateErrorResponse(&reqObj, foundation.InternalServerError, errStr)
// 		api.HttpWriteResponse(writer, &reqObj, res)
// 	} else {
// 		reqObj := conf_api.NewDeleteAuthRequest(realID)
// 		reqObj.Id = realID
// 		res := RequestToModelHandler(req, reqObj)
// 		if res == nil {
// 			res = DelAuth(me.DefaultDatabase, reqObj)
// 		}
// 		api.HttpWriteResponse(writer, reqObj, res)
// 	}

// }

func (me *ConfigServer) getAuth(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewQueryAuthInfoRequest()
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetAuth(me.SlaveDatabase, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

// func (me *ConfigServer) addMintInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddMintInfoRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	if res == nil {
// 		res = AddMintInfo(me, me.DefaultDatabase, reqObj)
// 	}
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) addBurnInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddBurnInfoRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	res = AddBurnInfo(me, me.DefaultDatabase, reqObj)

// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) addApprovalInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddBurnInfoRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	res = AddApprovalInfo(me, me.DefaultDatabase, reqObj)
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

// func (me *ConfigServer) addUserInfo(writer http.ResponseWriter, req *http.Request) {
// 	reqObj := conf_api.NewAddBurnInfoRequest()
// 	res := RequestToModelHandler(req, reqObj)
// 	res = AddUserInfo(me, me.DefaultDatabase, reqObj)
// 	api.HttpWriteResponse(writer, reqObj, res)
// }

func (me *ConfigServer) getPublicCurrencyNameByAddr(writer http.ResponseWriter, req *http.Request) {
	assetAddress := chi.URLParam(req, "assetAddress")
	reqObj := conf_api.NewQueryPublicCurrencyNameRequest(assetAddress)
	reqObj.CurrencyAddress = assetAddress
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetPublicCurrencyNameByAddr(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getCurrencyNameByAddr(writer http.ResponseWriter, req *http.Request) {
	assetAddress := chi.URLParam(req, "assetAddress")
	reqObj := conf_api.NewQueryCurrencyNameRequest(assetAddress)
	reqObj.CurrencyAddress = assetAddress
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetCurrencyNameByAddr(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getFaucetConfig(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := RequestToModelHandler(req, reqObj)
	res = GetFaucetConfig(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getServiceGroupAddress(writer http.ResponseWriter, req *http.Request) {

	reqObj := conf_api.NewQueryServiceGroupDetailRequest(0)
	res := RequestToModelHandler(req, reqObj)
	if res != nil {
		log.GetLogger(log.Name.Root).Errorln(res.GetMessage(), " nonce: ", res.GetNonce())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	groupIdStr := chi.URLParam(req, "groupId")
	groupId, err := strconv.ParseUint(groupIdStr, 10, 32)

	if err != nil {
		res := response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Invalid group Id")
		api.HttpWriteResponse(writer, reqObj, res)
		return
	} else {
		reqObj.GroupId = uint(groupId)
	}

	res = GetServiceGroupDetail(me.SlaveDatabase, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getAssetSettings(writer http.ResponseWriter, req *http.Request) {
	reqObj := conf_api.NewGetAssetSettingsRequest()
	res := RequestToModelHandler(req, reqObj)

	if res == nil {
		res = GetAssetSettings(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) getServiceGroupId(writer http.ResponseWriter, req *http.Request) {

	reqObj := new(conf_api.GetServiceGroupIdRequest)
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = GetServiceGroupId(me, reqObj)

	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) setServerWalletAddress(writer http.ResponseWriter, req *http.Request) {

	reqObj := new(conf_api.SetServerWalletAddressRequest)
	res := RequestToModelHandler(req, reqObj)
	if res == nil {
		res = processSetServerWalletAddress(me.DefaultDatabase, me.AuthClient.(*auth.AuthClient), reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *ConfigServer) InitConfigMQ() error {
	me.mqPublisher = network.NewMQPublisher(me.ServerConfig.MQUrl, network.MQModeTopic, "config_server_exchange", log.GetLogger(log.Name.Root))
	return me.mqPublisher.InitPublisher(true)
}
