package conf

import (
	"encoding/json"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"strconv"
)

func GetCurrencySCAddr(server *ConfigServer, req *conf_api.QueryCurrencySCAddrRequest) *response.ResponseBase {
	err := req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	address, err := GetCurrencySCAddrFromSC(server, req.Asset)
	if err != nil && err.Error() == "No such asset" {
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	} else if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get erc20 address from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	addrRes := conf_api.QueryAddressResponse{Address: address}
	res := response.CreateSuccessResponse(req, addrRes)
	return res
}

func AddCurrencyInfo(server *ConfigServer, req *conf_api.AddCurrencyRequest) *response.ResponseBase {
	err := req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	decimalDigit, err := strconv.ParseInt(req.Decimal, 10, 64)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}

	txHash, err := AddCurrencyInfoFromSC(server, req.Address, req.Asset, decimalDigit, req.ListID)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to add new asset to smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}

	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res
}

func DelCurrencyInfo(server *ConfigServer, req *conf_api.DelCurrencyRequest) *response.ResponseBase {
	err := req.ValidateField()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res
	}
	txHash, err := DelCurrencyInfoFromSC(server, req.Asset)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to delete asset from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	txHashRes := conf_api.TxHashResponse{TxHash: txHash}
	res := response.CreateSuccessResponse(req, txHashRes)
	return res
}

func GetCurrencyNameByAddr(server *ConfigServer, req *conf_api.QueryCurrencyNameRequest) *response.ResponseBase {
	assetName, err := GetCurrencyNameByAddrFromSC(server, req.CurrencyAddress)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get asset name from smart contract"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		return res
	}
	ares := conf_api.QueryCurrencyNameResponse{AssetName: assetName}
	res := response.CreateSuccessResponse(req, ares)
	return res
}
