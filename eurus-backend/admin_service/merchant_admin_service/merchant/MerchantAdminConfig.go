package merchant_admin

import (
	"eurus-backend/foundation/server"
)

type MerchantAdminServerConfig struct {
	server.ServerConfigBase
	PlatformWalletAddress string `json:"platformWalletAddress"`
}

func NewMerchantAdminServerConfig() *MerchantAdminServerConfig {
	config := new(MerchantAdminServerConfig)
	config.ActualConfig = config
	return config
}

func (me *MerchantAdminServerConfig) GetMerchantAdminServerConfig() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *MerchantAdminServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *MerchantAdminServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
