package conf_api

import (
	"net/http"

	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
)

type KeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type QueryKeyPairRequest struct {
	request.RequestBase
	KeyPair
}

type QueryKeyPairResponse struct {
	response.ResponseBase
	Data *KeyPair `json:"data"`
}

func NewKeyPair() *QueryKeyPairRequest {
	req := new(QueryKeyPairRequest)
	req.RequestPath = "/config/key"
	req.Method = http.MethodPost
	return req
}
