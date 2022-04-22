package kyc_model

//UserKYCStatusDetail
type UserKYCStatusDetail struct {
	UserKYCStatus
	WalletAddress string `gorm:"column:wallet_address" json:"walletAddress"`
	Email         string `gorm:"column:email" json:"email"`
}
