package reward

import (
	"eurus-backend/foundation/database"

	"github.com/shopspring/decimal"
)

type TokenDistributedType int

const (
	DistributedUnknown TokenDistributedType = iota
	DistributedRegistration
)

type TokenDistributedTriggerType int

const (
	TriggerNotApplicable     = 0
	TriggerDeposit           = 10
	TriggerWithdraw          = 20
	TriggerSideChainTransfer = 30
	TriggerSideChainAirDrop  = 40
)

type TokenDistributedStatus int

const (
	DistributedStatusError   TokenDistributedStatus = -1
	DistributedStatusPending TokenDistributedStatus = 10
	DistributedStatusSuccess TokenDistributedStatus = 20
)

type DistributedTokenError struct {
	database.DbModel
	Id              *uint64
	AssetName       string
	Amount          decimal.Decimal `gorm:"type:numeric"`
	Chain           *uint64
	DistributedType TokenDistributedType
	TriggerType     TokenDistributedTriggerType
	UserId          uint64
	TxHash          string
	FromAddress     string
	ToAddress       string
	GasPrice        *decimal.Decimal `gorm:"type:numeric"`
	GasUsed         uint64
	GasFee          *decimal.Decimal `gorm:"type:numeric"`
}

type DistributedToken struct {
	DistributedTokenError
	Status TokenDistributedStatus
}
