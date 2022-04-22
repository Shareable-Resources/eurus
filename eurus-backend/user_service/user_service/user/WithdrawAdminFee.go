package user

import "math/big"

type WithdrawAdminFee struct {
	Currency string  `json:"currency"`
	Fee      big.Int `json:"fee"`
	Decimal  int64   `json:"decimal"`
}
