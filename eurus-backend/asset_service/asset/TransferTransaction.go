package asset

import (
	"eurus-backend/foundation/database"
	"time"

	"github.com/shopspring/decimal"
)

type TransType int8

const (
	Transfer         TransType = 1
	Deposit          TransType = 2
	Withdraw         TransType = 3
	Purchase         TransType = 4
	MerchantDeposit  TransType = 5
	DistributedToken TransType = 6
	TopUp            TransType = 7
)

type TransferTransaction struct {
	database.DbModel
	UserId           uint64 `json:"userId"`
	AssetName        string `json:"assetName"`
	FromAddress      string `json:"fromAddress"`
	ToAddress        string `json:"toAddress"`
	RequestTransId   *uint64
	IsSend           bool
	TxHash           string          `json:"txHash"`
	Chain            int             `json:"chain"`
	Amount           decimal.Decimal `gorm:"type:numeric"`
	GasFee           decimal.Decimal `gorm:"type:numeric"`
	GasPrice         decimal.Decimal `gorm:"type:numeric"`
	TransGasUsed     *uint64
	UserGasUsed      *uint64
	Status           TransferStatus
	TransactionDate  time.Time `json:"transactionDate"`
	ConfirmTransHash string
	Remarks          string
}
type TransferStatus int16

const (
	TransferStatusError     TransferStatus = -1
	TransferStatusPending   TransferStatus = 10
	TransferStatusConfirmed TransferStatus = 50
)
