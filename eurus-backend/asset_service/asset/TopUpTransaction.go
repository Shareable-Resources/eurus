package asset

import (
	"eurus-backend/foundation/database"
	"time"

	"github.com/shopspring/decimal"
)

type TopUpStatus int16

const (
	TopUpError   TopUpStatus = -1
	TopUpSuccess TopUpStatus = 1
)

type TopUpTransaction struct {
	database.DbModel
	TxHash          string
	CustomerId      uint64
	CustomerType    CustomerType
	FromAddress     string
	ToAddress       string
	TransferGas     decimal.Decimal `gorm:"type:numeric"`
	TargetGas       decimal.Decimal `gorm:"type:numeric"`
	Status          TopUpStatus
	IsDirectTopUp   bool
	Remarks         string
	TransGasUsed    uint64
	UserGasUsed     uint64
	GasPrice        decimal.Decimal `gorm:"type:numeric"`
	TransactionDate time.Time
}
