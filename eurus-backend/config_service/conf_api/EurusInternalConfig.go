package conf_api

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type SetColdWalletAddressRequest struct {
	request.RequestBase
	Address string `json:"address"`
}

func NewSetColdWalletRequest() *SetColdWalletAddressRequest {
	req := new(SetColdWalletAddressRequest)
	req.RequestPath = "/config/asset/cold-wallet"
	req.Method = http.MethodPost
	return req
}

func (me *SetColdWalletAddressRequest) ValidateField() error {
	if me.Address == "" {
		return errors.New("address should not be null")
	} else if !common.IsHexAddress(me.Address) {
		return errors.New("address does not in correct format")
	}
	return nil
}

type AddPublicCurrencyRequest struct {
	request.RequestBase
	Currency
}

func NewAddPublicCurrencyRequest() *AddPublicCurrencyRequest {
	req := new(AddPublicCurrencyRequest)
	req.RequestPath = "/config/asset/currency"
	req.Method = http.MethodPost
	return req
}

func (me *AddPublicCurrencyRequest) ValidateField() error {
	if me.Address == "" {
		return errors.New("address should not be null")
	} else if !common.IsHexAddress(me.Address) {
		return errors.New("address does not in correct format")
	} else if me.Asset == "" {
		return errors.New("asset should not be null")
	}
	return nil
}

type QueryPublicCurrencySCAddrRequest struct {
	request.RequestBase
	Asset string `json:"asset"`
}

func NewQueryPublicCurrencySCAddrRequest() *QueryPublicCurrencySCAddrRequest {
	req := new(QueryPublicCurrencySCAddrRequest)
	req.RequestPath = "/config/asset/currency"
	req.Method = http.MethodGet
	return req
}

func (me *QueryPublicCurrencySCAddrRequest) ValidateField() error {
	if me.Asset == "" {
		return errors.New("asset should not be null")
	}
	return nil
}

type DelPublicCurrencyRequest struct {
	request.RequestBase
	Asset string `json:"asset"`
}

func NewDelPublicCurrencyRequest(asset string) *DelPublicCurrencyRequest {
	req := new(DelPublicCurrencyRequest)
	req.RequestPath = "/config/asset/public-currency/" + asset
	req.Method = http.MethodDelete
	req.Asset = asset
	return req
}

func (me *DelPublicCurrencyRequest) ValidateField() error {
	if me.Asset == "" {
		return errors.New("asset should not be null")
	}
	return nil
}

type QueryPublicCurrencyNameRequest struct{
	request.RequestBase
	CurrencyAddress string `json:"currencyAddress"`
}

func NewQueryPublicCurrencyNameRequest(currencyAddress string) *QueryPublicCurrencyNameRequest {
	req := new(QueryPublicCurrencyNameRequest)
	req.RequestPath = "/config/asset/public-currency/" + currencyAddress
	req.Method = http.MethodGet
	req.CurrencyAddress = currencyAddress
	return req
}

type QueryPublicCurrencyNameResponse struct{
	AssetName string `json:"assetName"`
}