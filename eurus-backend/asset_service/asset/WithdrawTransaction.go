package asset

import (
	"eurus-backend/foundation/database"
	"time"

	"github.com/shopspring/decimal"
)

type WithdrawTransactionData struct {
	withdrawTransactionList []*WithdrawTransaction
}

type WithdrawTransaction struct {
	database.DbModel
	Id                    uint64
	CustomerId            uint64
	CustomerType          CustomerType
	AssetName             string
	Amount                decimal.Decimal `gorm:"type:numeric"`
	ApprovalWalletAddress string
	RequestTransId        *uint64
	RequestTransHash      string
	RequestDate           time.Time
	ReviewDate            *time.Time
	ReviewedBy            string
	ReviewTransHash       string
	InnetFromAddress      string
	MainnetFromAddress    string
	MainnetToAddress      string
	MainnetTransHash      string
	MainnetTransDate      *time.Time
	BurnTransId           uint64
	BurnTransHash         string
	BurnDate              *time.Time
	Status                WithdrawStatus
	Remarks               string
	SidechainGasUsed      decimal.Decimal `gorm:"type:numeric"`
	SidechainGasFee       decimal.Decimal `gorm:"type:numeric"`
	AdminFee              decimal.Decimal `gorm:"type:numeric"`
	UserGasUsed           decimal.Decimal `gorm:"type:numeric"`
	GasPrice              decimal.Decimal `gorm:"type:numberic"`
}

type WithdrawStatus int16
type CustomerType int16

const (
	StatusError           WithdrawStatus = iota - 1
	StatusPendingApproval WithdrawStatus = iota * 10
	StatusApproved
	StatusRejected
	StatusBurnConfirming
	StatusBurned
	StatusConfirmingTransfer
	StatusTransferProcessing
	StatusCompleted
)

const (
	CustomerUser CustomerType = iota
	CustomerMerchant
)
