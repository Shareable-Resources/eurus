package wallet_bg_model

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type WalletBalance struct {
	WalletType    string          `gorm:"column:wallet_type" json:"walletType"`
	WalletAddress string          `gorm:"column:wallet_address" json:"walletAddress"`
	AssetName     string          `gorm:"column:asset_name" json:"assetName"`
	Balance       decimal.Decimal `gorm:"column:balance;type:numeric" json:"balance"`
	UserId        *uint64         `gorm:"column:user_id" json:"userId"`
	CreatedDate   time.Time       `gorm:"column:created_date" json:"createdDate"`
	MarkDate      time.Time
	ChainId       int `gorm:"column:chain_id" json:"chainId"`
}

func (t WalletBalance) TableName() string {
	return "wallet_balances"
}

const oneDaySeconds int64 = 3600 * 24

func NewWalletBalance(walletType string, walletAddress string, assetName string, balance *big.Int, userId *uint64, chainId int) *WalletBalance {
	obj := new(WalletBalance)
	obj.WalletType = walletType
	obj.WalletAddress = walletAddress
	obj.AssetName = assetName
	obj.Balance = decimal.NewFromBigInt(balance, 0)
	obj.UserId = userId
	obj.ChainId = chainId
	today := time.Now()
	currentDayUnix := (today.Unix() / oneDaySeconds) * oneDaySeconds
	currentDate := time.Unix(currentDayUnix, 0)
	yesterday := currentDate.AddDate(0, 0, -1)
	obj.CreatedDate = today
	obj.MarkDate = yesterday
	return obj
}
