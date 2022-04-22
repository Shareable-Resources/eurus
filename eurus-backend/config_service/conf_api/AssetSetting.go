package conf_api

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/database"
	"net/http"

	"github.com/shopspring/decimal"
)

type AssetSetting struct {
	database.DbModel           `json:"-"`
	ID                         uint64           `json:"id"`
	AssetName                  string           `json:"assetName"`
	KYC0MaxDailyWithdrawAmount *decimal.Decimal `json:"kyc0MaxDailyWithdrawAmount" gorm:"type:numeric"`
	KYC1MaxDailyWithdrawAmount *decimal.Decimal `json:"kyc1MaxDailyWithdrawAmount" gorm:"type:numeric"`
	KYC2MaxDailyWithdrawAmount *decimal.Decimal `json:"kyc2MaxDailyWithdrawAmount" gorm:"type:numeric"`
	KYC3MaxDailyWithdrawAmount *decimal.Decimal `json:"kyc3MaxDailyWithdrawAmount" gorm:"type:numeric"`
	SweepTriggerAmount         decimal.Decimal  `json:"sweepTriggerAmount" gorm:"type:numeric"`
	IsEnabled                  bool             `json:"isEnabled"`
}

type GetAssetSettingsRequest struct {
	request.RequestBase
}

func NewGetAssetSettingsRequest() *GetAssetSettingsRequest {
	req := new(GetAssetSettingsRequest)
	req.RequestPath = "/config/asset/settings"
	req.Method = http.MethodGet
	return req
}

type GetAssetSettingsFullResponse struct {
	response.ResponseBase
	Data []AssetSetting `json:"data"`
}
