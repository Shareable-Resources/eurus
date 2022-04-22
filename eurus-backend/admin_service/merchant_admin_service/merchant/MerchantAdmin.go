package merchant_admin

import "eurus-backend/foundation/database"

type MerchantAdmin struct {
	database.DbModel
	OperatorId   uint64
	MerchantId   uint64
	Username     string
	PasswordHash string
	Email        string
	Status       MerchantAccountState
}
