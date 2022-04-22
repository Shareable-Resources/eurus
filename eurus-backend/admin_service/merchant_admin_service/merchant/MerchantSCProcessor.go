package merchant_admin

type MerchantAdminSCProcessor struct {
	config           *MerchantAdminServerConfig
}

func NewMerchantAdminSCProcessor(config *MerchantAdminServerConfig) *MerchantAdminSCProcessor{
	merchantAdminSCProcessor := new(MerchantAdminSCProcessor)
	merchantAdminSCProcessor.config = config
	return merchantAdminSCProcessor
}