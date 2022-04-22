package wallet_bg_model

import (
	"eurus-backend/foundation/database"
	"time"

	"github.com/shopspring/decimal"
)

type AssetTotalSupply struct {
	database.DbModel
	Id           uint64
	AssetName    string
	TotalSupply  *decimal.Decimal `gorm:"type:numeric"`
	ChainId      int
	BlockNumber  *decimal.Decimal
	AssetAddress string
	MarkDate     time.Time
}
