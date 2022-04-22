package database

import (
	"time"

	"github.com/shopspring/decimal"
)

type DBModel struct {
	// Not using gorm.Model becaus no DeletedAt column
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DBNetwork struct {
	DBModel
	NetworkCode string
	NetworkName string
	ChainID     *int
	RpcURL      *string
}

type DBToken struct {
	DBModel
	NetworkID int64
	Address   string
	Symbol    string
	Name      string
	Decimals  int
}

type DBMerchant struct {
	DBModel
	MerchantCode    string
	MerchantName    string
	TagDisplayName  *string
	TagDescription  *string
	MerchantLastSeq int64
}

type DBMerchantWallet struct {
	DBModel
	MerchantID int64
	TokenID    int64
	Address    string
}

type DBSubmission struct {
	DBModel
	SubmitTime    time.Time
	NetworkID     int64
	TokenID       int64
	FromAddress   string
	Amount        decimal.Decimal `gorm:"type:numeric"`
	MerchantID    int64
	Tag           string
	TxHash        string
	TxStatus      int
	PaymentStatus int
	Signature     string
	MessageBody   string
}

type DBTransaction struct {
	DBModel
	SubmitTime    time.Time
	ConfirmedTime time.Time
	NetworkID     int64
	TokenID       int64
	FromAddress   string
	Amount        decimal.Decimal `gorm:"type:numeric"`
	MerchantID    int64
	Tag           string
	MerchantSeqNo int64
	SubmissionID  int64
	OnchainStatus int
	ConfirmStatus int
	Signature     string
	TxHash        string
	BlockHash     string
	BlockNumber   int64
}

type DBMerchantAPIKey struct {
	DBModel
	MerchantID int64
	APIKey     string
	Salt       string
}
