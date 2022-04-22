package wallet_bg_model

import (
	"eurus-backend/foundation/database"
)

type WalletBalanceConfig struct {
	database.DbModel
	Id             uint64
	ConfigType     WalletBalanceConfigType
	ServiceGroupId int
	ServiceId      uint64
	WalletAddress  string
	WalletType     int
	ChainId        int
	AssetName      string
	Description    string
}

type WalletBalanceConfigType int

const (
	WalletServiceGroupId WalletBalanceConfigType = iota
	WalletServiceId
	WalletAddress
	WalletConfigResolved //Internal used only
)
