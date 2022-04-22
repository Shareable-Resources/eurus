package admin

import (
	eurus_ethereum "eurus-backend/foundation/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AdminSCProcessor struct {
	AdminContractAbi             *abi.ABI
	Config                       *AdminServerConfig
	sidechainEthClient           *eurus_ethereum.EthClient
	mainnetEthClient             *eurus_ethereum.EthClient
	MainnetPlatformWalletAddress *common.Address
}

func NewAdminSCProcessor(config *AdminServerConfig) *AdminSCProcessor {
	processor := new(AdminSCProcessor)
	processor.AdminContractAbi = new(abi.ABI)
	// *processor.AdminContractAbi, _ = abi.JSON(strings.NewReader(string(contract.AdminSmartContractABI)))
	processor.Config = config
	return processor
}
