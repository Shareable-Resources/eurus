package conf

import (
	"encoding/json"
	"errors"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"math"
	"math/big"

	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"

	"github.com/shopspring/decimal"
)

func GetServerSetting(db *database.ReadOnlyDatabase, req *conf_api.QueryServerSettingRequest) *response.ResponseBase {
	var err error

	data, err := GetServerSettingFromDB(db)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get server setting"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, data)
	return res
}

func GetConfig(db *database.ReadOnlyDatabase, req *conf_api.QueryConfigRequest) *response.ResponseBase {
	var err error
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	data, err := GetConfigFromDB(db, req)
	configInfo := conf_api.ConfigInfo{ConfigData: *data}
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get config info"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, configInfo)
	return res
}

func GetConfigAuth(db *database.ReadOnlyDatabase, req *conf_api.QueryConfigAuthInfoRequest) *response.ResponseBase {
	var err error
	data, err := GetConfigAuthFromDB(db, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get config and auth info"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, data)
	return res
}

func UpdateConfig(server *ConfigServer, req *conf_api.AddConfigRequest) *response.ResponseBase {
	err := req.ValidateSetField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	err = checkServerIdFromAuth(&req.RequestBase)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("checkServerIdFromAuth(&req.RequestBase) Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.ServerIdUnmatch, err.Error())
		return res
	}
	isExist, err := ConfigIsExist(&server.DefaultDatabase.ReadOnlyDatabase, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to check if the config is exist: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	if isExist {
		err = UpdateConfigToDB(server.DefaultDatabase, req)
	} else {
		err = InsertConfigToDB(server.DefaultDatabase, req)
	}
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to update new config info: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	if req.Id == 0 && req.Key == "internalSCConfigAddress" && !*req.IsService {
		server.ServerConfig.InternalSCConfigAddress = ethereum.ToLowerAddressString(req.Value)
	} else if req.Id == 0 && req.Key == "eurusInternalConfigAddress" && !*req.IsService {
		server.ServerConfig.EurusInternalConfigAddress = ethereum.ToLowerAddressString(req.Value)
	} else if req.Id == 0 && req.Key == "externalSCConfigAddress" && !*req.IsService {
		server.ServerConfig.ExternalSCConfigAddress = ethereum.ToLowerAddressString(req.Value)
	}

	var mqEvent *conf_api.MQConfigServiceEvent
	if isExist {
		mqEvent = conf_api.NewMQConfigServiceEvent(conf_api.PublishActionUpdate, *req.ConfigMap)
	} else {
		mqEvent = conf_api.NewMQConfigServiceEvent(conf_api.PublishActionInsert, *req.ConfigMap)
	}
	server.mqPublisher.PublishJson("config.config", mqEvent, nil)

	res := response.CreateSuccessResponse(req, nil)
	return res
}

func DelConfig(server *ConfigServer, req *conf_api.DeleteConfigRequest) *response.ResponseBase {

	err := req.ValidateDeleteField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	err = checkServerIdFromAuth(&req.RequestBase)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("checkServerIdFromAuth(&req.RequestBase) Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.ServerIdUnmatch, err.Error())
		return res
	}

	err = DelConfigFromDB(server.DefaultDatabase, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to delete config info: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	publishEvent := conf_api.NewMQConfigServiceEvent(conf_api.PublishActionDelete, req)

	server.mqPublisher.PublishJson("config.config", publishEvent, nil)
	res := response.CreateSuccessResponse(req, nil)
	return res
}

func GetSystemConfig(db *database.ReadOnlyDatabase, req *conf_api.GetSystemConfigRequest) *response.ResponseBase {
	found, systemConfig, err := TryGetSystemConfigFromDB(db, req.Key)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get system config: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	if !found {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Requested system config not found\nRequest Params: " + string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RecordNotFound, "Requested system config not found")
		return res
	}

	res := response.CreateSuccessResponse(req, systemConfig)
	return res
}

func AddOrUpdateSystemConfig(server *ConfigServer, req *conf_api.AddOrUpdateSystemConfigRequest) *response.ResponseBase {
	serviceID, err := extractServiceIDFromLoginToken(req.LoginToken)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to extract ServiceID from LoginToken: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.BadRequest, err.Error())
		return res
	}

	isInGroup, groupID, err := TryGetGroupIDOfService(&server.DefaultDatabase.ReadOnlyDatabase, serviceID)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to check GroupID: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	// Determine need to insert or update system config
	found, systemConfig, err := TryGetSystemConfigFromDB(&server.DefaultDatabase.ReadOnlyDatabase, req.Key)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get system config: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	if !found {
		// The OwnerID to be inserted depends on whether this ServiceID is in group
		var ownerID int64
		if isInGroup {
			ownerID = groupID
		} else {
			ownerID = serviceID
		}

		systemConfig, err = InsertSystemConfigToDB(server.DefaultDatabase, ownerID, !isInGroup, req.Key, req.Value)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("Unable to insert new system config: "+err.Error(), "\nRequest Params: "+string(reqStr))
			res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
			return res
		}

		res := response.CreateSuccessResponse(req, conf_api.SystemConfigWithUpdateStatus{SystemConfig: *systemConfig, IsNewConfig: true})
		publishEvent := conf_api.NewMQConfigServiceEvent(conf_api.PublishActionInsert, req.Key)
		_ = server.mqPublisher.PublishJson("config.system", publishEvent, nil)
		return res
	}

	// System config is found, need to check if this ServiceID or GroupID is allowed to update config
	if systemConfig.IsService {
		if systemConfig.OwnerID != serviceID {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("ServiceID not allowed to update this system config\nRequest Params: " + string(reqStr))
			res := response.CreateErrorResponse(req, foundation.BadRequest, "ServiceID not allowed to update this system config")
			return res
		}
	} else {
		if !isInGroup || systemConfig.OwnerID != groupID {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("GroupID not allowed to update this system config\nRequest Params: " + string(reqStr))
			res := response.CreateErrorResponse(req, foundation.BadRequest, "GroupID not allowed to update this system config")
			return res
		}
	}

	systemConfig, err = UpdateSystemConfigInDB(server.DefaultDatabase, req.Key, req.Value)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to update system config: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	publishEvent := conf_api.NewMQConfigServiceEvent(conf_api.PublishActionUpdate, req.Key)
	_ = server.mqPublisher.PublishJson("config.system", publishEvent, nil)
	res := response.CreateSuccessResponse(req, conf_api.SystemConfigWithUpdateStatus{SystemConfig: *systemConfig, IsNewConfig: false})

	return res
}

func extractServiceIDFromLoginToken(loginToken auth_base.ILoginToken) (int64, error) {
	userIdMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(loginToken.GetUserId()), &userIdMap)
	if err != nil {
		return 0, err
	}

	serviceID, ok := userIdMap["serviceId"].(float64)
	if !ok {
		return 0, errors.New("Unable to read ServiceID from LoginToken")
	}

	return int64(serviceID), nil
}

//check serverID from AuthServer == 0  , if yes, it is the operator
//reqObj.GetLoginToken().GetUserId() ["serviceId"] : ILoginToken from Auth server
func checkServerIdFromAuth(reqObj *request.RequestBase) error {

	if reqObj.GetLoginToken() == nil {
		return errors.New("reqObj.GetLoginToken() == nil")
	}

	userIdMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(reqObj.GetLoginToken().GetUserId()), &userIdMap)
	if err != nil {
		return err
	}

	iLoginTokenFromAuth, ok := userIdMap["serviceId"].(float64)
	if !ok {
		return errors.New(`userIdMap["serviceId"].(float64) error`)
	}

	if !(iLoginTokenFromAuth == 0) {
		return errors.New(`!(iLoginTokenFromAuth == 0)`)
	}
	return nil
}

func GetFaucetConfig(server *ConfigServer, req *request.RequestBase) *response.ResponseBase {
	faucetConfigs, err := GetFaucetConfigFromDB(server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to GetFaucetConfig: " + err.Error())
		return response.CreateErrorResponse(req, foundation.MethodNotFound, err.Error())
	}

	return response.CreateSuccessResponse(req, faucetConfigs)

}

func GetServiceGroupDetail(db *database.ReadOnlyDatabase, req *conf_api.QueryServiceGroupDetailRequest) *response.ResponseBase {
	serverList, err := GetServerDetailByGroupFromDB(db, req.GroupId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetServerDetailByGroupFromDB error: ", err)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	res := new(conf_api.QueryServiceGroupDetailResponse)
	res.GroupId = uint(req.GroupId)
	res.ServerList = serverList

	return response.CreateSuccessResponse(req, res)
}

func GetServiceGroupId(server *ConfigServer, req *conf_api.GetServiceGroupIdRequest) *response.ResponseBase {

	var serviceId int64
	var err error
	if req.ServiceId == 0 {
		serviceId, err = server.GetServiceIdFromServerLoginToken(req.GetLoginToken())
	} else {
		serviceId = req.ServiceId
	}
	if err != nil {
		res := response.CreateErrorResponse(req, foundation.ServerTokenInvalid, err.Error())
		return res
	}

	groupId, err := DBGetGroupIdByServiceId(server.SlaveDatabase, serviceId)
	if err != nil {
		res := response.CreateErrorResponse(req, foundation.ServerTokenInvalid, err.Error())
		return res
	}
	return response.CreateSuccessResponse(req, groupId)
}

func SetExchangeAmount(server *ConfigServer, req *conf_api.SetExchangeRate) *response.ResponseBase {

	err := req.ValidateSetExchangeRate()
	reqStr, _ := json.Marshal(req)

	if err != nil {

		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	err = SetExchangeRateToDB(server.DefaultDatabase, req.Rate, req.AssetName)
	if err != nil {

		log.GetLogger(log.Name.Root).Error("Unable to update DB : "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	withdrawFee, err := GetWithdrawFeeFromSC(server)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to GetWithdrawFeeFromSC : "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res

	}
	decimalCount, err := GetAssetDecimalFromDB(server.SlaveDatabase, req.AssetName)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to GetAssetDecimalFromDB : "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	adminFee, err := CalculateAdminFeeAmount(server, withdrawFee, req.Rate, decimalCount)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("CalculateAdminFeeAmount error : "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	tx2, err := AddAdminFeeToSC(server, req.AssetName, adminFee)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("AddAdminFeeToSC error : "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	log.GetLogger(log.Name.Root).Debugln("AddAdminFee trans hash: ", tx2.Hash().Hex(), " asset name: ", req.AssetName, " value: ", adminFee)
	receipt, err := server.EthClient.QueryEthReceipt(tx2)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get receipt : "+err.Error(), " tran hash: ", tx2.Hash().Hex(), " Request Params: "+string(reqStr))

		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	if receipt.Status == 0 {
		log.GetLogger(log.Name.Root).Debugln("AddAdminFee status is 0. trans hash: ", tx2.Hash().Hex(), " asset name: ", req.AssetName)
		res := response.CreateErrorResponse(req, foundation.EthereumError, "Receipt status error: "+receipt.RevertReason)
		return res
	}

	res := response.CreateSuccessResponse(req, tx2.Hash().Hex())
	return res
}

func CalculateAdminFeeAmount(me *ConfigServer, ethFee *big.Int, inputRate decimal.Decimal, decimal int64) (*big.Int, error) {

	FormatWithdrawFee := new(big.Float)
	FormatWithdrawFee.Quo(new(big.Float).SetInt(ethFee), new(big.Float).SetFloat64(math.Pow(10, 18)))

	targetFeeFloat := new(big.Float)
	targetFeeFloat.Quo(FormatWithdrawFee, inputRate.BigFloat()) //admin fee / rate = new Currency amount (unit in ETH(gel))

	newDecimal := new(big.Float).SetFloat64(math.Pow(10, float64(decimal)))

	targetFeeNoDecimal := new(big.Float).Mul(newDecimal, targetFeeFloat)
	targetFeeBigInt, _ := targetFeeNoDecimal.Int(nil)

	return targetFeeBigInt, nil

}

func SetWithdrawFee(me *ConfigServer, req *conf_api.SetWithdrawFeeETH) *response.ResponseBase {

	err := req.ValidateSetWithdrawFeeField()
	reqStr, _ := json.Marshal(req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	txHash, err := SetWithdrawFeeFromSC(me, req)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	receipt, err := me.EthClient.QueryEthReceipt(txHash)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query receipt: ", err, "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}

	if receipt.Status == 0 {
		log.GetLogger(log.Name.Root).Errorln("Receipt status is 0. Revert reason: ", receipt.RevertReason, " \nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, "Receipt status error: "+receipt.RevertReason)
		return res
	}

	err = SetETHWithdrawFeeToDB(me.DefaultDatabase, req.Value)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	txHashRes := conf_api.TxHashResponse{TxHash: txHash.Hash().Hex()}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res

}

func GetAssetSettings(server *ConfigServer, req *conf_api.GetAssetSettingsRequest) *response.ResponseBase {
	settings, err := DBGetAssetSettings(server.DefaultDatabase)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get user asset settings: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	res := response.CreateSuccessResponse(req, settings)
	return res
}
