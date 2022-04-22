package conf_api

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"fmt"
	"net/http"
	"strconv"
)

type QueryConfigRequest struct {
	request.RequestBase
	Id  int64  `json:"id"`
	Key string `json:"key"`
}
type QueryConfigResponse struct {
	response.ResponseBase
	Data *ConfigInfo `json:"data"`
}
type ConfigInfo struct {
	ConfigData []ConfigMap `json:"config"`
}

func NewQueryConfigRequest() *QueryConfigRequest {
	req := new(QueryConfigRequest)
	req.RequestPath = "/config"
	req.Method = http.MethodPost
	return req
}

type AddConfigRequest struct {
	request.RequestBase
	*ConfigMap
}

func NewAddConfigRequest() *AddConfigRequest {
	req := new(AddConfigRequest)
	req.RequestPath = "/config/update"
	req.ConfigMap = new(ConfigMap)
	req.Id = -1
	req.Method = http.MethodPost
	return req
}

func (me *AddConfigRequest) ValidateSetField() error {
	fmt.Println("check error ValidateSetField")
	fmt.Println(me.IsService)
	if me.Key == "" && me.Value == "" {
		return errors.New("key and value should not be null")
	} else if me.Id == -1 {
		return errors.New("serverId should not be null")
	} else if me.IsService == nil {
		return errors.New("is_service must be include")
	}
	return nil
}

type DeleteConfigRequest struct {
	QueryConfigRequest
	IsGroup bool
}

func NewDeleteConfigRequest(serverId int, key string, isGroup bool) *DeleteConfigRequest {
	req := new(DeleteConfigRequest)
	req.RequestPath = "/config/" + strconv.Itoa(serverId) + "/" + key
	req.Method = http.MethodDelete
	req.IsGroup = isGroup
	return req
}

func (me *DeleteConfigRequest) ValidateDeleteField() error {
	if me.Key == "" {
		return errors.New("Key should not be null")
	}
	return nil
}

type QueryFaucetConfigRequest struct {
	request.RequestBase
}

func NewQueryFaucetConfig() *QueryFaucetConfigRequest {
	req := new(QueryFaucetConfigRequest)
	req.Method = http.MethodGet
	req.RequestPath = "/config/asset/faucet-config"

	return req
}

type QueryFaucetConfigFullResponse struct {
	response.ResponseBase
	Data []*FaucetConfig `json:"data"`
}

type SetServerWalletAddressRequest struct {
	request.RequestBase
	WalletAddress string `json:"walletAddress"`
}

func NewSetServerWalletAddressRequest() *SetServerWalletAddressRequest {
	req := new(SetServerWalletAddressRequest)
	req.Method = http.MethodPost
	req.RequestPath = "/config/auth/setWalletAddress"
	return req
}

type SetServerWalletAddressFullResponse struct {
	response.ResponseBase
}
