package merchant_common

import (
	"eurus-backend/foundation/database"

	"github.com/shopspring/decimal"
)

type RefundRequestStatus int

const (
	RefundPending RefundRequestStatus = iota
	RefundAccepted
	RefundRejected
)

type MerchantRefundRequest struct {
	database.DbModel
	Id                 uint64
	DestAddress        string
	AssetName          string
	UserId             *uint64
	Amount             decimal.Decimal `gorm:"type:numeric"`
	PurchaseTransHash  string
	RefundReason       string
	OperatorComment    string
	Status             RefundRequestStatus
	MerchantId         uint64
	MerchantOperatorId *uint64
	RefundTransHash    string
}
