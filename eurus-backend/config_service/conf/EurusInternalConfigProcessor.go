package conf

import (
	"encoding/json"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

// func SetColdWalletAddress(server *ConfigServer, req *conf_api.SetColdWalletAddressRequest) *response.ResponseBase {
// 	var err error
// 	err = req.ValidateField()
// 	if err != nil {
// 		reqStr, _ := json.Marshal(req)
// 		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
// 		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
// 		return res
// 	}
// 	txHash, err := SetColdWalletAddressFromSC(server, req.Address)
// 	if err != nil {
// 		reqStr, _ := json.Marshal(req)
// 		log.GetLogger(log.Name.Root).Error("Unable to set cold wallet address to smart contract", err.Error(), "\nRequest Params: "+string(reqStr))
// 		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
// 		return res
// 	}
// 	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
// 	res := response.CreateSuccessResponse(req, txHashRes)
// 	return res
// }

func AddPublicCurrencyInfo(server *ConfigServer, req *conf_api.AddPublicCurrencyRequest) *response.ResponseBase {
	err := req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	txHash, err := AddPublicCurrencyInfoFromSC(server, req.Address, req.Asset)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to add new public asset to smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res
}

// func DelPublicCurrencyInfo(server *ConfigServer, req *conf_api.DelPublicCurrencyRequest) *response.ResponseBase {
// 	err := req.ValidateField()
// 	if err != nil {
// 		reqStr, _ := json.Marshal(req)
// 		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
// 		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
// 		return res
// 	}
// 	txHash, err := DelPublicCurrencyInfoFromSC(server, req.Asset)
// 	if err != nil {
// 		reqStr, _ := json.Marshal(req)
// 		log.GetLogger(log.Name.Root).Error("Unable to delete public asset from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
// 		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
// 		return res
// 	}
// 	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
// 	res := response.CreateSuccessResponse(req, txHashRes)
// 	return res
// }

func GetPublicCurrencySCAddr(server *ConfigServer, req *conf_api.QueryPublicCurrencySCAddrRequest) *response.ResponseBase {
	err := req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	address, err := GetCurrencySCAddrFromSC(server, req.Asset)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get erc20 smart contract address from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	if err != nil && err.Error() == "No such asset" {
		//if such an asset does not exist in the mapping list
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	addrRes := conf_api.QueryAddressResponse{Address: address}
	res := response.CreateSuccessResponse(req, addrRes)
	return res
}

func GetPublicCurrencyNameByAddr(server *ConfigServer, req *conf_api.QueryPublicCurrencyNameRequest) *response.ResponseBase {
	assetName, err := GetPublicCurrencyNameByAddrFromSC(server, req.CurrencyAddress)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get asset name from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	ares := conf_api.QueryPublicCurrencyNameResponse{AssetName: assetName}
	res := response.CreateSuccessResponse(req, ares)
	return res
}
