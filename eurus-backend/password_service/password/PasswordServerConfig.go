package password

import "eurus-backend/foundation/server"

type PasswordServerConfig struct {
	server.ServerConfigBase
	UDSPath  string `json:"udsPath"`
	Password string `json:"-" eurus_conf:"noPrint"`
}

func NewPasswordServerConfig() *PasswordServerConfig {
	config := new(PasswordServerConfig)
	config.ActualConfig = config
	return config
}

func (me *PasswordServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *PasswordServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *PasswordServerConfig) ValidateField() {
	if me.UDSPath == "" {
		panic("Empty UDSPath")
	}
}
