package sign

import (
	"eurus-backend/foundation/server"
)

type SignServerConfig struct {
	server.ServerConfigBase
	UserWalletOwnerPrivateKey          string `json:"userWalletOwnerPrivateKey" eurus_conf:"noPrint"`
	InvokerPrivateKey                  string `json:"invokerPrivateKey" eurus_conf:"noPrint"`
	invokerAddress                     string
	SideChainGasLimit                  uint64 `json:"sideChainGasLimit"`
	CentralizedUserWalletMnemonicPhase string `json:"centralizedUserWalletMnemonicPhase" eurus_conf:"noPrint"`
}

func NewSignServerConfig() *SignServerConfig {
	config := new(SignServerConfig)
	config.ActualConfig = config
	return config
}

func (me *SignServerConfig) GetSignServerConfig() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *SignServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *SignServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *SignServerConfig) GetConfigFileOnlyFieldList() []string {
	return []string{"centralizedUserWalletMnemonicPhase", "userWalletOwnerPrivateKey", "invokerPrivateKey"}
}

func (me *SignServerConfig) ValidateField() {
	if me.CentralizedUserWalletMnemonicPhase == "" {
		panic("CentralizedUserWalletMnemonicPhase is mandatory field in config file")
	}
	if me.InvokerPrivateKey == "" {
		panic("InvokerPrivateKey is mandatory field in config file")
	}
	if me.UserWalletOwnerPrivateKey == "" {
		panic("UserWalletOwnerPrivateKey is mandatory field in config file")
	}

}
