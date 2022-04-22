package background

import "eurus-backend/foundation/server"

type BackgoundServerConfig struct {
	server.ServerConfigBase
}

func NewBackgroundServerConfig() *BackgoundServerConfig {
	config := new(BackgoundServerConfig)
	config.ActualConfig = config
	return config
}
