package conf_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
)

type ServerSetting struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Path string `json:"path"`
}

type QueryServerSettingResponse struct {
	response.ResponseBase
	Data *ServerSetting `json:"data"`
}
type QueryServerSettingRequest struct {
	request.RequestBase
}

func NewQueryServerSettingRequest() *QueryServerSettingRequest {
	req := new(QueryServerSettingRequest)
	req.RequestPath = "/config/setting"
	req.Method = http.MethodGet
	return req
}
