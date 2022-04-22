package admin

import (
	"eurus-backend/foundation/server"
)

type AdminServerConfig struct {
	server.ServerConfigBase
	PlatformWalletAddress string `json:"platformWalletAddress"`
}

func NewAdminServerConfig() *AdminServerConfig {
	config := new(AdminServerConfig)
	config.ActualConfig = config
	return config
}

func (me *AdminServerConfig) GetAdminServerConfig() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *AdminServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *AdminServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
