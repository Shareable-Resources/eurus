package conf_api

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"
)

type QueryAddressResponse struct {
	Address string `json:"address"`
}

type TxHashResponse struct {
	TxHash string `json:"txHash"`
}

type SetWalletAddressRequest struct {
	request.RequestBase
	Address string `json:"address"`
}

type SetWithdrawFeeETH struct {
	request.RequestBase
	Value string `json:"value"`
}

type SetExchangeRate struct {
	request.RequestBase
	AssetName string          `json:"assetName"` //This parameter is currency ID
	Rate      decimal.Decimal `json:"rate" gorm:"type:numeric"`
}

type Asset struct {
	Decimal    int64  `json:"decimal"`
	AssetName  string `json:"assetName"`
	CurrencyId string `json:"currencyId"`
	AutoUpdate bool   `json:"auto_update"`
}

type GetAssetRequest struct {
	request.RequestBase
}

type GetAssetFullResponse struct {
	response.ResponseBase
	Data []*Asset `json:"data"`
}

type SetAssetRequest struct {
	request.RequestBase
	Asset
}

type SetWall struct {
	request.RequestBase
	Address string `json:"address"`
}

type ExchangeRates struct {
	AssetName        string
	Rate             decimal.Decimal `gorm:"type:numeric"`
	CreatedDate      time.Time
	LastModifiedDate time.Time
}

func NewSetInnetWalletRequest() *SetWalletAddressRequest {
	req := new(SetWalletAddressRequest)
	req.RequestPath = "/config/asset/innet-wallet"
	req.Method = http.MethodPost
	return req
}

func (me *SetWalletAddressRequest) ValidateField() error {
	if me.Address == "" {
		return errors.New("address should not be null")
	} else if !common.IsHexAddress(me.Address) {
		return errors.New("address does not in correct format")
	}
	return nil
}

type SetWalletAddressResponse struct {
	response.ResponseBase
	TxHash string `json:"txHash"`
}

func NewQueryMainnetWalletRequest() *request.RequestBase {
	req := new(request.RequestBase)
	req.RequestPath = "/config/asset/mainnet-wallet"
	req.Method = http.MethodGet
	return req
}

func NewQueryInnetWalletRequest() *request.RequestBase {
	req := new(request.RequestBase)
	req.RequestPath = "/config/asset/innet-wallet"
	req.Method = http.MethodGet
	return req
}

func NewSetMainnetWalletRequest() *SetWalletAddressRequest {
	req := new(SetWalletAddressRequest)
	req.RequestPath = "/config/asset/mainnet-wallet"
	req.Method = http.MethodPost
	return req
}
func NewSetWithdrawFeeETH() *SetWithdrawFeeETH {
	req := new(SetWithdrawFeeETH)
	req.RequestPath = "/config/asset/withdraw-fee-eth"
	req.Method = http.MethodPost
	return req
}
func NewSetExchangeRate() *SetExchangeRate {
	req := new(SetExchangeRate)
	req.RequestPath = "/config/asset/exchange-rate"
	req.Method = http.MethodPost
	return req
}
func NewSetAssetRequest() *SetAssetRequest {
	req := new(SetAssetRequest)
	req.RequestPath = "/config/asset/public-currency-db"
	req.Method = http.MethodPost
	return req
}
func NewGetAssetRequest() *GetAssetRequest {
	req := new(GetAssetRequest)
	req.RequestPath = "/config/asset/public-currency-db"
	req.Method = http.MethodGet
	return req
}

func NewGetAssetResponse() *GetAssetFullResponse {
	res := new(GetAssetFullResponse)
	return res
}

type Currency struct {
	Address string `json:"address"`
	Asset   string `json:"asset"`
	Decimal string `json:"decimal"`
	ListID  string `json:"id"`
}

type AddCurrencyRequest struct {
	request.RequestBase
	Currency
}

func (me *AddCurrencyRequest) ValidateField() error {
	if me.Address == "" {
		return errors.New("address should not be null")
	} else if !common.IsHexAddress(me.Address) {
		return errors.New("address does not in correct format")
	} else if me.Asset == "" {
		return errors.New("asset should not be null")
	} else if me.Decimal == "" {
		return errors.New("decimal should not be null")
	} else if me.ListID == "" {
		return errors.New("id should not be null")
	}
	return nil
}

func NewAddCurrencyRequest() *AddCurrencyRequest {
	req := new(AddCurrencyRequest)
	req.RequestPath = "/config/asset/currency"
	req.Method = http.MethodPost
	return req
}

type QueryCurrencySCAddrRequest struct {
	request.RequestBase
	Asset string `json:"asset"`
}

func NewQueryCurrencySCAddrRequest() *QueryCurrencySCAddrRequest {
	req := new(QueryCurrencySCAddrRequest)
	req.RequestPath = "/config/asset/currency"
	req.Method = http.MethodGet
	return req
}

func (me *QueryCurrencySCAddrRequest) ValidateField() error {
	if me.Asset == "" {
		return errors.New("asset should not be null")
	}
	return nil
}

type DelCurrencyRequest struct {
	request.RequestBase
	Asset string `json:"asset"`
}

func NewDelCurrencyRequest(asset string) *DelCurrencyRequest {
	req := new(DelCurrencyRequest)
	req.RequestPath = "/config/asset/currency/" + asset
	req.Method = http.MethodDelete
	req.Asset = asset
	return req
}

func (me *DelCurrencyRequest) ValidateField() error {
	if me.Asset == "" {
		return errors.New("asset should not be null")
	}
	return nil
}

type QueryCurrencyNameRequest struct {
	request.RequestBase
	CurrencyAddress string `json:"currencyAddress"`
}

func NewQueryCurrencyNameRequest(currencyAddress string) *QueryCurrencyNameRequest {
	req := new(QueryCurrencyNameRequest)
	req.RequestPath = "/config/asset/currency/" + currencyAddress
	req.Method = http.MethodGet
	req.CurrencyAddress = currencyAddress
	return req
}

type QueryCurrencyNameResponse struct {
	AssetName string `json:"assetName"`
}

func (me *SetWithdrawFeeETH) ValidateSetWithdrawFeeField() error {
	var err error
	if me.Value == "" {
		err = errors.New("Value should not be null")
		return err
	}
	_, err = strconv.ParseInt(me.Value, 10, 64)

	return err
}

func (me *SetExchangeRate) ValidateSetExchangeRate() error {
	var err error
	if me.AssetName == "" {
		err = errors.New("Asset name should not be null")
	} else if me.Rate == decimal.NewFromInt(0) {
		err = errors.New("rate should not be null")
	}

	return err
}
func (me *SetAssetRequest) ValidateSetAsset() error {
	var err error
	if me.AssetName == "" {
		err = errors.New("Asset name should not be null")
	} else if me.CurrencyId == "" {
		err = errors.New("Currency should not be null")
	}

	return err
}
