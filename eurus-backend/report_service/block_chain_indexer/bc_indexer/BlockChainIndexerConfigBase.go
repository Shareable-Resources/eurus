package bc_indexer

import (
	"eurus-backend/foundation/server"

	"github.com/sirupsen/logrus"
)

type BlockChainIndexerConfigBase struct { //implements IServerConfig
	server.ServerConfigBase
	MainnetEthClientChainID           int    `json:"mainnetEthClientChainId"`
	MainnetEthClientProtocol          string `json:"mainnetEthClientProtocol"`
	MainnetEthClientIP                string `json:"mainnetEthClientIP"`
	MainnetEthClientPort              int    `json:"mainnetEthClientPort"`
	MainnetEthClientWebSocketProtocol string `json:"mainnetEthClientWebSocketProtocol"`
	MainnetEthClientWebSocketPort     int    `json:"mainnetEthClientWebSocketPort"`
	MainnetEthClientWebSocketIP       string `json:"mainnetEthClientWebSocketIP"`
	RegistrationCriteriaListJsonStr   string `json:"registrationRewardSetting"`
}

func NewBlockChainIndexerConfigBase() *BlockChainIndexerConfigBase {
	config := new(BlockChainIndexerConfigBase)
	config.ActualConfig = config
	return config
}

func (me *BlockChainIndexerConfigBase) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me BlockChainIndexerConfigBase) ValidateField() {
	//me.ServerConfigBase.ValidateField()
}

func (me BlockChainIndexerConfigBase) SetHttpErrorLogger(logger *logrus.Logger) {
	me.ServerConfigBase.SetHttpErrorLogger(logger)
}

func (me *BlockChainIndexerConfigBase) GetParent() interface{} {
	return &me.ServerConfigBase
}
