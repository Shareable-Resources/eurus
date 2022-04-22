package approval

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"

	"github.com/shopspring/decimal"
)

type PendingPrewithdraw struct {
	database.DbModel
	Id                    uint64 `"gorm:default:0"`
	CustomerId            uint64
	CustomerType          asset.CustomerType
	InnetFromAddress      string
	MainnetToAddress      string
	ApprovalWalletAddress string
	RequestTransId        *uint64
	RequestTransHash      string
	AssetName             string
	Amount                decimal.Decimal `gorm:"type:numeric"`
	Status                WithdrawStatus
	SidechainGasUsed      decimal.Decimal `gorm:"type:numeric"`
	SidechainGasFee       decimal.Decimal `gorm:"type:numeric"`
	AdminFee              decimal.Decimal `gorm:"type:numeric"`
	UserGasUsed           decimal.Decimal `gorm:"type:numeric"`
	GasPrice              decimal.Decimal `gorm:"type:numeric"`
}

type WithdrawStatus = asset.WithdrawStatus
