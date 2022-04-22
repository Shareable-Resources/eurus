package deposit

import (
	"errors"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type DepositObserverConfig struct { //implements IRetrySetting
	server.ServerConfigBase
	MainnetEthClientChainID           int    `json:"mainnetEthClientChainId"`
	MainnetEthClientProtocol          string `json:"mainnetEthClientProtocol"`
	MainnetEthClientIP                string `json:"mainnetEthClientIP"`
	MainnetEthClientPort              int    `json:"mainnetEthClientPort"`
	MainnetEthClientWebSocketProtocol string `json:"mainnetEthClientWebSocketProtocol"`
	MainnetEthClientWebSocketPort     int    `json:"mainnetEthClientWebSocketPort"`
	MainnetEthClientWebSocketIP       string `json:"mainnetEthClientWebSocketIP"`
	MainnetBlockConfirmCount          int    `json:"mainnetBlockConfirmCount"`
	LocalDbFileName                   string `json:"localDbFileName"`
	SideChainGasLimit                 uint64 `json:"sideChainGasLimit"`
	RegistrationRewardSetting         string `json:"registrationRewardSetting"`

	//From system config
	SweepServiceInvokerAddress common.Address `json:"-"`

	AssetSettings map[string]conf_api.AssetSetting `json:"-"`
}

func NewDepositObserverConfig() *DepositObserverConfig {
	config := new(DepositObserverConfig)
	config.AssetSettings = make(map[string]conf_api.AssetSetting)
	config.ActualConfig = config
	return config
}

func (me *DepositObserverConfig) ValidateEthClientField() error {
	var err error = nil

	err = me.ServerConfigBase.ValidateEthClientField()
	if err != nil {
		return err
	}

	if me.MainnetEthClientWebSocketProtocol == "" {
		err = errors.New("MainnetEthClientWebSocketProtocol should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.MainnetEthClientWebSocketIP == "" {
		err = errors.New("MainnetEthClientWebSocketIP should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.MainnetEthClientWebSocketPort == 0 {
		err = errors.New("MainnetEthClientWebSocketPort should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.MainnetEthClientPort == 0 {
		err = errors.New("MainnetEthClientPort should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.MainnetEthClientIP == "" {
		err = errors.New("MainnetEthClientIP should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	}
	return err
}

func (me *DepositObserverConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *DepositObserverConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *DepositObserverConfig) GetRetryCount() int {
	return me.RetryCount
}

func (me *DepositObserverConfig) GetRetryInterval() time.Duration {
	return time.Duration(me.RetryInterval)
}
