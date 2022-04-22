package merchant_admin

type MerchantLoginToken struct {
	MerchantId   uint64               `json:"merchantId"`
	OperatorId   uint64               `json:"operatorId"`
	AccountState MerchantAccountState `json:"accountState"`
}
