package test

import "eurus-backend/foundation/api/response"

type TestData struct {
	ResponseData string
}

type TestResponse struct {
	response.ResponseBase
	Data TestData `json:"data"`
}
type TestRequest struct {
	Nonce      string `json:"nonce"`
	MethodName string `json:"methodName"`
}
