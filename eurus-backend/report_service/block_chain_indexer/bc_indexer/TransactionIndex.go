package bc_indexer

import "time"

type TransactionIndex struct {
	TxHash        string    `json:"txHash"`
	WalletAddress string    `json:"walletAddress"`
	UserId        uint64    `json:"userId"`
	CreatedDate   time.Time `json:"createdDate"`
	AssetName     string    `json:"assetName"`
	Status        bool      `json:"status"`
}
