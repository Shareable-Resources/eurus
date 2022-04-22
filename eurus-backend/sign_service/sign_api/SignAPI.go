package sign_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
)

type WalletKeyType int

const (
	WalletKeyUserWalletOwner WalletKeyType = iota
	WalletKeyInvoker
)

type SignTransactionRequest struct {
	request.RequestBase
	Token string `json:"-"`
	Data  []byte `json:"data"`
	Sign  string `json:"verifySign"`
}

type SignTransactionResponse struct {
	SignedTrans []byte `json:"signedTrans"`
}

type SignTransactionFullResponse struct {
	response.ResponseBase
	Data *SignTransactionResponse `json:"data"`
}

func NewSignTransactionRequest() *SignTransactionRequest {
	req := new(SignTransactionRequest)
	req.RequestPath = RootPath + SignTransactionPath
	req.Method = http.MethodPost
	return req
}

type QueryAddressRequest struct {
	request.RequestBase
	Token   string        `json:"-"`
	KeyType WalletKeyType `json:"walletKeyType,string"`
}

func NewQueryAddressRequest(token string) (me *QueryAddressRequest) {
	req := new(QueryAddressRequest)
	req.RequestPath = RootPath + GetAddressPath
	req.Method = http.MethodGet
	req.Token = token
	return req
}

type QueryAddressResponse struct {
	Address string `json:"address"`
}

type GetPendingNonceRequest struct {
	request.RequestBase
	Token   string        `json:"-"`
	KeyType WalletKeyType `json:"walletKeyType,string"`
}

func NewGetPendingNonceRequest(token string) (me *GetPendingNonceRequest) {
	req := new(GetPendingNonceRequest)
	req.RequestPath = RootPath + GetPendingNoncePath
	req.Method = http.MethodGet
	req.Token = token
	return req
}

type GetNonceResponse struct {
	PendingNonce uint64 `json:"pendingNonce"`
}

type GetPendingNonceFullResponse struct {
	response.ResponseBase
	Data *GetNonceResponse `json:"data"`
}

type CalibrateNonceRequest struct {
	request.RequestBase
	KeyType WalletKeyType `json:"walletKeyType"`
	Token   string        `json:"-"`
}

func NewCalibrateNonceRequest() *CalibrateNonceRequest {
	return new(CalibrateNonceRequest)
}

type CalibrateNonceResponse struct {
	Nonce   uint64        `json:"nonce"`
	KeyType WalletKeyType `json:"walletKeyType"`
}

type SignUserWalletTransactionRequest struct {
	request.RequestBase
	Value         string `json:"value"`
	To            string `json:"to"`
	GasPrice      uint64 `json:"gasPrice"`
	InputFunction string `json:"inputFunction"`
	LoginToken    string `json:"-"`
}

func NewSignUserWalletTransactionRequest() *SignUserWalletTransactionRequest {
	req := new(SignUserWalletTransactionRequest)
	req.Method = http.MethodPost
	req.RequestPath = RootPath + SignUserWalletTransPath
	return req
}

type SignedUserWalletTransactionResponse struct {
	SignedTx string `json:"signedTx"`
}

type SignedUserWalletTransactionFullResponse struct {
	response.ResponseBase
	Data *SignedUserWalletTransactionResponse `json:"data"`
}

type GetCentralizedUserMainnetAddressRequest struct {
	request.RequestBase
	UserId uint64 `json:"userId,string"`
}

func NewGetCentralizedUserMainnetAddressRequest() *GetCentralizedUserMainnetAddressRequest {
	req := new(GetCentralizedUserMainnetAddressRequest)
	req.Method = http.MethodGet
	req.RequestPath = RootPath + GetCentralizedUserMainnetAddressPath
	return req
}

type GetCentralizedUserMainnetAddressFullResponse struct {
	response.ResponseBase
	Data string `json:"data"`
}
