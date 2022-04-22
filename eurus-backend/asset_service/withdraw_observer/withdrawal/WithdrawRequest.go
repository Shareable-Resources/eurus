package withdrawal

import (
	"eurus-backend/foundation/database"
)

type WithdrawRequestData struct {
	withdrawRequestList []*WithdrawRequest
}

type WithdrawRequest struct {
	database.DbModel
	WithdrawId            uint64
	RequestTransHash      string
	ServiceId             uint64
	Data				  []byte
	Signature			  []byte
}