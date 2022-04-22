package conf

import (
	"eurus-backend/foundation/server"
)

// Extends ServerConfigBase struct, should add new attributes in here
// Database Config > .json config > base class config(Default config --- ServerConfig)
type ConfigServerConfig struct {
	server.ServerConfigBase
	SideChainGasLimit uint64 `json:"sideChainGasLimit"`
}

func NewConfigServerConfig() *ConfigServerConfig {
	config := new(ConfigServerConfig)
	config.ActualConfig = config
	return config
}

func (me *ConfigServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *ConfigServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *ConfigServerConfig) GetConfigFileOnlyFieldList() []string {
	return []string{}
}
