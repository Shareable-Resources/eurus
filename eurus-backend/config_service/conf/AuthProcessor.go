package conf

import (
	"encoding/json"
	"eurus-backend/auth_service/auth"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
)

func GetAuth(db *database.ReadOnlyDatabase, req *conf_api.QueryAuthInfoRequest) *response.ResponseBase {
	var err error
	data, err := GetAuthFromDB(db)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get auth info"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, data)
	return res
}

func AddAuth(db *database.Database, req *conf_api.AddAuthRequest) *response.ResponseBase {
	var err error
	err = req.ValidateSetField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	auth, err := InsertAuthToDB(db, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to insert auth info to db"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		return res
	}
	addedRes := conf_api.AddAuthResponse{Id: auth.Id}
	res := response.CreateSuccessResponse(req, addedRes)
	return res
}

func UpdateAuth(db *database.Database, req *conf_api.UpdateAuthRequest) *response.ResponseBase {
	var err error
	err = req.ValidateSetField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	err = UpdateAuthToDB(db, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to insert auth info to db"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, nil)
	return res
}

func DelAuth(db *database.Database, req *conf_api.DeleteAuthRequest) *response.ResponseBase {
	var err error
	err = req.ValidateDeleteField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}

	err = DelAuthFromDB(db, req)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to delete auth info to db"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, nil)
	return res
}

func processSetServerWalletAddress(db *database.Database, authClient *auth.AuthClient, req *conf_api.SetServerWalletAddressRequest) *response.ResponseBase {

	serviceId, err := authClient.GetServiceIdFromServerLoginToken(req.LoginToken)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Invalid login token: "+err.Error())
	}
	if serviceId > 0 {
		err = DBUpdateAuthServiceWalletAddress(db, uint64(serviceId), ethereum.ToLowerAddressString(req.WalletAddress))
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to update DBUpdateAuthServiceWalletAddress error: ", err)
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}
		return response.CreateSuccessResponse(req, nil)
	}
	return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "No service Id found")
}
