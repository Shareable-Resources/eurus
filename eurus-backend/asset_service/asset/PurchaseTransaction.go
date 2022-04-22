package asset

import (
	"eurus-backend/foundation/database"

	"github.com/shopspring/decimal"
)

type PurchaseStatus int16

const (
	PurchaseStatusError     PurchaseStatus = -1
	PurchaseStatusConfirmed PurchaseStatus = 50
)

type PurchaseTransaction struct {
	database.DbModel
	UserId       uint64
	TxHash       string
	Chain        int
	FromAddress  string
	ToAddress    string
	AssetName    string
	Amount       decimal.Decimal `gorm:"type:numeric"`
	ProductId    *uint64
	Quantity     *decimal.Decimal `gorm:"type:numeric"`
	GasFee       decimal.Decimal  `gorm:"type:numeric"`
	TransGasUsed uint64
	UserGasUsed  uint64
	GasPrice     decimal.Decimal `gorm:"type:numeric"`
	Status       PurchaseStatus
	PurchaseType TransType
	Remarks      string
}
