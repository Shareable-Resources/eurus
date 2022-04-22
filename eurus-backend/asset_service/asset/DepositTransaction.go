package asset

import (
	"eurus-backend/foundation/database"
	"time"

	"github.com/shopspring/decimal"
)

type DepositStatus int16

const (
	DepositError            DepositStatus = -1
	DepositReceiptCollected DepositStatus = 10
	DepositAssetCollected   DepositStatus = 20
	DepositMintRequesting   DepositStatus = 30
	// DepositMintConfirming   DepositStatus = 30
	DepositCompleted DepositStatus = 40
	SweepTrans       DepositStatus = 0x7FFF
)

type DepositTransaction struct {
	database.DbModel
	Id                      uint64
	MainnetTransHash        string
	Amount                  decimal.Decimal `gorm:"type:numeric"`
	AssetName               string
	MainnetFromAddress      string
	MainnetToAddress        string
	MainnetTransDate        *time.Time
	MainnetGasFee           decimal.Decimal `gorm:"type:numeric"`
	MainnetGasUsed          decimal.Decimal `gorm:"type:numeric"`
	MainnetCollectTransHash string
	MainnetCollectTransDate *time.Time
	MintTransId             *uint64
	MintTransHash           string
	MintDate                *time.Time
	InnetToAddress          string
	InnetFromAddress        string
	CustomerId              uint64
	CustomerType            CustomerType
	Status                  DepositStatus
	Remarks                 string
}
