package asset

import (
	"eurus-backend/foundation/database"
	"github.com/shopspring/decimal"
)

type AssetAllocationCost struct {
	database.DbModel
	Id             uint64
	TransHash      string
	AllocationType string
	GasUsed        uint64  `gorm:"type:numeric"`
	GasPrice       decimal.Decimal `gorm:"type:numeric"`
}
