package conf_api

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
	"strconv"
)

type QueryAuthInfoRequest struct {
	request.RequestBase
	AuthInfo
}

type QueryAuthInfoResponse struct {
	response.ResponseBase
	Data *AuthInfo `json:"data"`
}

type AuthInfo struct {
	AuthData []AuthService `json:"auth"`
}

func NewQueryAuthInfoRequest() *QueryAuthInfoRequest {
	req := new(QueryAuthInfoRequest)
	req.RequestPath = "/config/auth"
	req.Method = http.MethodPost
	return req
}

type UpdateAuthRequest struct {
	request.RequestBase
	AuthService
}

func (me *UpdateAuthRequest) ValidateSetField() error {
	if me.ServiceName == "" || me.PubKey == "" {
		return errors.New("service_name and pubKey should not be null")
	} else if me.Id == -1 {
		return errors.New("id should not be null")
	}
	return nil
}

func NewUpdateAuthRequest() *UpdateAuthRequest {
	req := new(UpdateAuthRequest)
	req.RequestPath = "/config/auth/update"
	req.Id = -1
	req.Method = http.MethodPost
	return req
}

type AddAuthRequest struct {
	request.RequestBase
	AuthService
}

type AddAuthResponse struct {
	Id int `json:"id"`
}

func (me *AddAuthRequest) ValidateSetField() error {
	if me.ServiceName == "" || me.PubKey == "" {
		return errors.New("serviceName and pubKey should not be null")
	}
	return nil
}

func NewAddAuthRequest() *AddAuthRequest {
	req := new(AddAuthRequest)
	req.RequestPath = "/config/auth/insert"
	req.Method = http.MethodPut
	return req
}

type DeleteAuthRequest struct {
	request.RequestBase
	Id int `json:"id"`
}

func (me *DeleteAuthRequest) ValidateDeleteField() error {
	if me.Id == -1 {
		return errors.New("Id should not be null")
	}
	return nil
}

func NewDeleteAuthRequest(id int) *DeleteAuthRequest {
	req := new(DeleteAuthRequest)
	req.Id = -1
	req.RequestPath = "/config/auth/" + strconv.Itoa(id)
	req.Method = http.MethodDelete
	return req
}
