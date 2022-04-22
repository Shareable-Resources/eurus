package sweep

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/server"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type SweepServiceConfig struct { //implements IRetrySetting
	server.ServerConfigBase
	DBPollingInterval                  int    `json:"dbPollingInterval"`
	MainnetEthClientChainID            int    `json:"mainnetEthClientChainId"`
	MainnetEthClientIP                 string `json:"mainnetEthClientIP"`
	MainnetEthClientPort               int    `json:"mainnetEthClientPort"`
	MainnetEthClientProtocol           string `json:"mainnetEthClientProtocol"`
	MainnetEthClientWebSocketIP        string `json:"mainnetEthClientWebSocketIP"`
	MainnetEthClientWebSocketPort      int    `json:"mainnetEthClientWebSocketPort"`
	MainnetEthClientWebSocketProtocol  string `json:"mainnetEthClientWebSocketProtocol"`
	SweepERC20Workers                  int    `json:"sweepErc20Workers"`
	SweepETHWorkers                    int    `json:"sweepEthWorkers"`
	SweepExtraGasFee                   int64  `json:"sweepExtraGasFee"`
	SweepExtraGasLimit                 uint64 `json:"sweepExtraGasLimit"`
	CentralizedUserWalletMnemonicPhase string `json:"centralizedUserWalletMnemonicPhase"`
	InvokerPrivateKey                  string `json:"invokerPrivateKey" eurus_conf:"noPrint"`
	QueryReceiptRetryCount             int    `json:"queryReceiptRetryCount"`
	CurrencyToSymbol                   map[string][]string
	SymbolToCurrency                   map[string]string
	AssetSettings                      map[string]conf_api.AssetSetting
	InvokerAddress                     common.Address
}

func NewSweepServiceConfig() *SweepServiceConfig {
	config := new(SweepServiceConfig)
	config.CurrencyToSymbol = make(map[string][]string)
	config.SymbolToCurrency = make(map[string]string)
	config.AssetSettings = make(map[string]conf_api.AssetSetting)
	config.ActualConfig = config
	return config
}

func (conf *SweepServiceConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &conf.ServerConfigBase
}

func (conf *SweepServiceConfig) GetParent() interface{} {
	return &conf.ServerConfigBase
}

func (conf *SweepServiceConfig) GetRetryCount() int {
	return conf.RetryCount
}

func (conf *SweepServiceConfig) GetRetryInterval() time.Duration {
	return time.Duration(conf.RetryInterval)
}

func (conf *SweepServiceConfig) GetConfigFileOnlyFieldList() []string {
	return append(conf.ServerConfigBase.GetConfigFileOnlyFieldList(), "centralizedUserWalletMnemonicPhase", "invokerPrivateKey")
}

func (conf *SweepServiceConfig) ValidateField() {
	if conf.CentralizedUserWalletMnemonicPhase == "" {
		panic("CentralizedUserWalletMnemonicPhase is mandatory field in config file")
	}

	if conf.InvokerPrivateKey == "" {
		panic("InvokerPrivateKey is mandatory field in config file")
	}

	conf.ServerConfigBase.ValidateField()
}
