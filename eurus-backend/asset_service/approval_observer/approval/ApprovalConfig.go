package approval

import (
	"eurus-backend/foundation/server"

	"github.com/ethereum/go-ethereum/common"
)

type ApprovalObserverConfig struct {
	server.ServerConfigBase
	ApprovalWalletAddress common.Address
	SideChainGasLimit     uint64 `json:"sideChainGasLimit"`
}

func NewApprovalObserverConfig() *ApprovalObserverConfig {
	config := new(ApprovalObserverConfig)
	config.ActualConfig = config
	return config
}
