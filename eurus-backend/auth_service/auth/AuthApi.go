package auth

import (
	"eurus-backend/foundation/api/request"
	"net/http"
	"strings"
)

type RequestLoginTokenByLoginRequestTokenRequest struct {
	Token string
	request.RequestBase
}

func NewRequestLoginTokenByLoginRequestTokenRequest(token string) *RequestLoginTokenByLoginRequestTokenRequest {
	req := new(RequestLoginTokenByLoginRequestTokenRequest)
	req.RequestPath = RootPath + RequestLoginRequestTokenPath
	req.Method = http.MethodPost
	req.Token = strings.Split(token, " ")[1]
	return req
}

type RequestLoginRequestTokenResponse struct {
	LoginRequestToken string `json:"loginRequestToken"`
}
