package conf_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
)

type SystemConfig struct {
	OwnerID   int64  `json:"ownerId"`
	IsService bool   `json:"isService"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type SystemConfigWithUpdateStatus struct {
	SystemConfig
	IsNewConfig bool `json:"isNewConfig"`
}

type GetSystemConfigRequest struct {
	request.RequestBase
	Key string `json:"key"`
}

func NewGetSystemConfigRequest(key string) *GetSystemConfigRequest {
	req := new(GetSystemConfigRequest)
	req.RequestPath = "/config/system/" + key
	req.Method = http.MethodGet
	return req
}

type GetSystemConfigFullResponse struct {
	response.ResponseBase
	Data *SystemConfig `json:"data"`
}

type AddOrUpdateSystemConfigRequest struct {
	request.RequestBase
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddOrUpdateSystemConfigFullResponse struct {
	response.ResponseBase
	Data *SystemConfigWithUpdateStatus `json:"data"`
}

func NewAddOrUpdateSystemConfigRequest(key string) *AddOrUpdateSystemConfigRequest {
	req := new(AddOrUpdateSystemConfigRequest)
	req.RequestPath = "/config/system/" + key
	req.Method = http.MethodPost
	return req
}
