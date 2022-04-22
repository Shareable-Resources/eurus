package withdrawal

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/server"

	"github.com/ethereum/go-ethereum/common"
)

type WithdrawObserverConfig struct { // implements IRetrySetting
	server.ServerConfigBase
	WithdrawEventTopic        common.Hash
	BurnCompletedEventTopic   common.Hash
	LocalStateFilePath        string //full path, relative
	WithdrawSmartContractAddr common.Address

	//Ethereum client
	MainnetEthClientProtocol string `json:"mainnetEthClientProtocol"`
	MainnetEthClientIP       string `json:"mainnetEthClientIP"`
	MainnetEthClientPort     int    `json:"mainnetEthClientPort"`

	MainnetEthClientChainID  int    `json:"mainnetEthClientChainID"`
	MainnetTransferGasTipCap int64  `json:"mainnetTransferGasTipCap"` //Unit is wei
	MainnetTransferGasFeeCap int64  `json:"mainnetTransferGasFeeCap"`
	MainnetTransferGasLimit  uint64 `json:"mainnetTransferGasLimit"`

	SideChainGasLimit int64 `json:"sideChainGasLimit"`
}

func NewWithdrawObserverConfig() *WithdrawObserverConfig {
	config := new(WithdrawObserverConfig)

	abi := ethereum.DefaultABIDecoder.GetABI("WithdrawSmartContract")
	event := abi.Events["BurnCompletedEvent"]
	config.BurnCompletedEventTopic = event.ID

	event = abi.Events["WithdrawEvent"]
	config.WithdrawEventTopic = event.ID
	config.ActualConfig = config
	return config
}

func (me *WithdrawObserverConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
