package auth

import "eurus-backend/foundation/server"

type AuthServerConfig struct {
	server.ServerConfigBase
	ResponseDelayInterval int64 `json:"responseDelayInterval"` //In milliseconds
}

func NewAuthServerConfig() *AuthServerConfig {
	config := new(AuthServerConfig)
	config.ActualConfig = config
	return config
}

func (me *AuthServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *AuthServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
