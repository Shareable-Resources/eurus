package asset

import (
	"time"
)

type TransactionIndex struct {
	TxHash        string    `json:"tx_hash"`
	WalletAddress string    `json:"wallet_address"`
	UserId        uint64    `json:"userId"`
	CreatedDate   time.Time `json:"created_date"`
	AssetName     string    `json:"asset_name"`
	Status        bool      `json:"status"`
}
