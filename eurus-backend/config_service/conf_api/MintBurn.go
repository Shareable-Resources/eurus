package conf_api

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
)

type MintBurnServerSetting struct {
	EthClientIP       string      `json:"ethClientIP"`
	EthClientPort     int         `json:"ethClientPort"`
	Name              string      `json:"name"`
	MiscellaneousData []ConfigMap `json:"miscellaneousData"`
}

type AddMintInfoRequest struct {
	request.RequestBase
	MintBurnServerSetting
}

type AddMintBurnInfoResponse struct {
	HdWalletAddress string `json:"hdWalletAddr"`
	RsaPublicKey    string `json:"rsaPubKey"`
	RsaPrivateKey   string `json:"rsaPrivKey"`
	ServerId        int    `json:"serverId"`
}

type AddMintInfoResponse struct {
	AddMintBurnInfoResponse
	response.ResponseBase
}

func NewAddMintInfoRequest() *AddMintInfoRequest {
	req := new(AddMintInfoRequest)
	req.RequestPath = "/config/mint"
	req.Method = http.MethodPost
	return req
}

func (me *MintBurnServerSetting) ValidateField() error {
	var err error
	if me.EthClientIP == "" {
		err = errors.New("ethClientIP should not be null")
	} else if me.EthClientPort == 0 {
		err = errors.New("ethClientPort should not be null")
	} else if me.Name == "" {
		err = errors.New("name should not be null")
	}
	return err
}

type AddBurnInfoRequest struct {
	request.RequestBase
	MintBurnServerSetting
}

type QueryBurnInfoResponse struct {
	response.ResponseBase
	AddMintBurnInfoResponse
}

func NewAddBurnInfoRequest() *AddBurnInfoRequest {
	req := new(AddBurnInfoRequest)
	req.RequestPath = "/config/burn"
	req.Method = http.MethodPost
	return req
}
