package conf

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/database"
)

func UpdateMintBurnInfoResponseHandler(configServer *ConfigServer, config *conf_api.ConfigMap,
	keyPair *conf_api.KeyPair, walletAddr string) *conf_api.AddMintBurnInfoResponse {
	var updateMintBurnInfoResponse = &conf_api.AddMintBurnInfoResponse{}
	updateMintBurnInfoResponse.HdWalletAddress = walletAddr
	updateMintBurnInfoResponse.RsaPublicKey = keyPair.PublicKey
	updateMintBurnInfoResponse.RsaPrivateKey = keyPair.PrivateKey

	updateMintBurnInfoResponse.ServerId = config.Id
	return updateMintBurnInfoResponse
}

func AddMintInfo(configServer *ConfigServer, db *database.Database, req *conf_api.AddMintInfoRequest) *response.ResponseBase {
	return response.CreateErrorResponse(req, foundation.NotImplementError, foundation.NotImplementError.String())
	// var err error
	// err = req.ValidateField()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
	// 	return res
	// }
	// keyPair, err := GenerateKeyPair()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("RSA Key Pair Generation Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// config, err := AddMintBurnInfoToDB(configServer, db, keyPair.PublicKey, &req.MintBurnServerSetting, conf_api.ServiceGroupDeposit)
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to add burn info"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	// 	return res
	// }
	// _, _, walletAddr, err := secret.GenerateMintBurnKey(keyPair.PrivateKey, strconv.Itoa(config.Id))
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to generate HDWallet key"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }

	// updateMintBurnInfoResponse := UpdateMintBurnInfoResponseHandler(configServer, config, keyPair, walletAddr)
	// res := response.CreateSuccessResponse(req, updateMintBurnInfoResponse)
	// return res
}

func AddBurnInfo(configServer *ConfigServer, db *database.Database, req *conf_api.AddBurnInfoRequest) *response.ResponseBase {
	return response.CreateErrorResponse(req, foundation.NotImplementError, foundation.NotImplementError.String())
	// var err error
	// err = req.ValidateField()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
	// 	return res
	// }

	// keyPair, err := GenerateKeyPair()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("RSA Key Pair Generation Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// config, err := AddMintBurnInfoToDB(configServer, db, keyPair.PublicKey, &req.MintBurnServerSetting, conf_api.ServiceGroupWithdraw)
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to add burn info"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }

	// _, _, walletAddr, err := secret.GenerateMintBurnKey(keyPair.PrivateKey, strconv.Itoa(config.Id))
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to generate HDWallet key"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// updateMintBurnInfoResponse := UpdateMintBurnInfoResponseHandler(configServer, config, keyPair, walletAddr)
	// res := response.CreateSuccessResponse(req, updateMintBurnInfoResponse)
	// return res
}

func AddApprovalInfo(configServer *ConfigServer, db *database.Database, req *conf_api.AddBurnInfoRequest) *response.ResponseBase {
	return response.CreateErrorResponse(req, foundation.NotImplementError, foundation.NotImplementError.String())
	// var err error
	// err = req.ValidateField()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
	// 	return res
	// }

	// keyPair, err := GenerateKeyPair()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("RSA Key Pair Generation Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// config, err := AddMintBurnInfoToDB(configServer, db, keyPair.PublicKey, &req.MintBurnServerSetting, conf_api.ServiceGroupApproval)
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to add burn info"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }

	// _, _, walletAddr, err := secret.GenerateMintBurnKey(keyPair.PrivateKey, strconv.Itoa(config.Id))
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to generate HDWallet key"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// updateMintBurnInfoResponse := UpdateMintBurnInfoResponseHandler(configServer, config, keyPair, walletAddr)
	// res := response.CreateSuccessResponse(req, updateMintBurnInfoResponse)
	// return res
}

func AddUserInfo(configServer *ConfigServer, db *database.Database, req *conf_api.AddBurnInfoRequest) *response.ResponseBase {
	return response.CreateErrorResponse(req, foundation.NotImplementError, foundation.NotImplementError.String())
	// var err error
	// err = req.ValidateField()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
	// 	return res
	// }

	// keyPair, err := GenerateKeyPair()
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("RSA Key Pair Generation Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// config, err := AddMintBurnInfoToDB(configServer, db, keyPair.PublicKey, &req.MintBurnServerSetting, conf_api.ServiceGroupUser)
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to add burn info"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }

	// _, _, walletAddr, err := secret.GenerateMintBurnKey(keyPair.PrivateKey, strconv.Itoa(config.Id))
	// if err != nil {
	// 	reqStr, _ := json.Marshal(req)
	// 	log.GetLogger(log.Name.Root).Error("Unable to generate HDWallet key"+err.Error(), "\nRequest Params: "+string(reqStr))
	// 	res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	// 	return res
	// }
	// updateMintBurnInfoResponse := UpdateMintBurnInfoResponseHandler(configServer, config, keyPair, walletAddr)
	// res := response.CreateSuccessResponse(req, updateMintBurnInfoResponse)
	// return res
}
