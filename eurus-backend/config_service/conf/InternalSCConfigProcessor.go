package conf

import (
	"encoding/json"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

func SetMainnetWalletAddress(server *ConfigServer, req *conf_api.SetWalletAddressRequest) *response.ResponseBase {
	var err error
	err = req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	txHash, err := SetMainnetWalletAddressFromSC(server, req.Address)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to set mainnet wallet address to smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res
}

func GetMainnetWalletAddress(server *ConfigServer, req *request.RequestBase) *response.ResponseBase {
	var err error
	address, err := GetMainnetWalletAddressFromSC(server)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get mainnet wallet address from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	addrRes := conf_api.QueryAddressResponse{Address: address}
	res := response.CreateSuccessResponse(req, addrRes)
	return res
}

func SetInnetWalletAddress(server *ConfigServer, req *conf_api.SetWalletAddressRequest) *response.ResponseBase {
	var err error
	err = req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	txHash, err := SetInnetWalletAddressFromSC(server, req.Address)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to set innet wallet address to smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res
}

func GetInnetWalletAddress(server *ConfigServer, req *request.RequestBase) *response.ResponseBase {
	var err error
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	address, err := GetInnetWalletAddressFromSC(server)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get innet wallet address from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	addrRes := conf_api.QueryAddressResponse{Address: address}
	res := response.CreateSuccessResponse(req, addrRes)
	return res
}
