package wallet_bg

import "eurus-backend/foundation/server"

type WalletBackgroundIndexerConfig struct {
	server.ServerConfigBase
	MainnetEthClientChainID           int      `json:"mainnetEthClientChainId"`
	MainnetEthClientProtocol          string   `json:"mainnetEthClientProtocol"`
	MainnetEthClientIP                string   `json:"mainnetEthClientIP"`
	MainnetEthClientPort              int      `json:"mainnetEthClientPort"`
	MainnetEthClientWebSocketProtocol string   `json:"mainnetEthClientWebSocketProtocol"`
	MainnetEthClientWebSocketPort     int      `json:"mainnetEthClientWebSocketPort"`
	MainnetEthClientWebSocketIP       string   `json:"mainnetEthClientWebSocketIP"`
	UserWalletOwnerWalletAddr         string   `json:"userWalletOwnerWalletAddr"`
	InvokerWalletAddrJson             string   `json:"invokerWalletAddrJson"`
	InvokerWalletAddrList             []string `json:"-"`
}

func NewWalletBackgroundIndexerConfig() *WalletBackgroundIndexerConfig {
	config := new(WalletBackgroundIndexerConfig)
	config.ActualConfig = config
	config.InvokerWalletAddrList = make([]string, 0)
	return config
}

func (me *WalletBackgroundIndexerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *WalletBackgroundIndexerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
