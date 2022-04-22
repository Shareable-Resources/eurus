package conf_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
)

type ConfigAuthInfo struct {
	ConfigData []ConfigMap   `json:"config"`
	AuthData   []AuthService `json:"auth"`
}

type QueryConfigAuthInfoRequest struct {
	request.RequestBase
	ServiceId int64 `json:"serviceId"`
}

type QueryConfigAuthResponse struct {
	response.ResponseBase
	Data *ConfigAuthInfo `json:"data"`
}

func NewQueryConfigAuthInfoRequest() *QueryConfigAuthInfoRequest {
	req := new(QueryConfigAuthInfoRequest)
	req.RequestPath = "/config/auth/all"
	req.Method = http.MethodPost
	return req
}
