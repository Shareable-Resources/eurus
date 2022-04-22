package asset

import (
	"eurus-backend/foundation/database"

	"github.com/shopspring/decimal"
)

type PendingSweepWallet struct {
	database.DbModel
	ID                   uint64
	UserID               *uint64
	MainnetWalletAddress string
	AssetName            string
	PreviousGasFeeCap    *decimal.Decimal `gorm:"type:numeric"`
	PreviousGasTipCap    *decimal.Decimal `gorm:"type:numeric"`
	PreviousGasLimit     *uint64
}
