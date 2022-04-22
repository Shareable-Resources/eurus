package conf_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
	"strconv"
)

type ServiceGroup struct {
	GroupId   int `json:"group_id"`
	ServiceId int `json:"service_id"`
}

type QueryServiceGroupDetailRequest struct {
	request.RequestBase
	GroupId uint `json:"groupId,string"` //Since this is a query string parameter at URL, it is a string, need to tell Unmarshaller to convert from string to integer
}

type QueryServiceGroupDetailResponse struct {
	GroupId    uint            `json:"groupId"`
	ServerList []*ServerDetail `json:"serverList"`
}

type QueryServiceGroupDetailFullResponse struct {
	response.ResponseBase
	Data *QueryServiceGroupDetailResponse `json:"data"`
}

type ServerDetail struct {
	Name          string `gorm:"column:service_name" json:"name"`
	ServiceId     uint64 `gorm:"column:id" json:"serviceId"`
	WalletAddress string `json:"walletAddress"`
}

func NewQueryServiceGroupDetailRequest(groupId uint64) *QueryServiceGroupDetailRequest {
	req := new(QueryServiceGroupDetailRequest)
	req.RequestPath = "/config/group/address/" + strconv.FormatUint(uint64(groupId), 10)
	req.Method = http.MethodGet
	req.GroupId = uint(groupId)
	return req
}

type GetServiceGroupIdRequest struct {
	request.RequestBase
	ServiceId int64 `json:"serviceId,string"` //Optional , 0 means get current service id from login token //Since this is a query string parameter at URL, it is a string, need to tell Unmarshaller to convert from string to integer
}

func NewGetServiceGroupIdRequest() *GetServiceGroupIdRequest {
	req := new(GetServiceGroupIdRequest)
	req.RequestPath = "/config/group/getGroupId"
	req.Method = http.MethodGet
	return req
}

type GetServiceGroupIdFullResponse struct {
	response.ResponseBase
	Data int64 `json:"data"`
}
