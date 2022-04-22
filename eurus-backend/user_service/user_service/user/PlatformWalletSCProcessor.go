package user

import (
	"eurus-backend/smartcontract/build/golang/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
)

func GetPlatformWalletAddress(server *UserServer) (*common.Address, error) {
	internalSC, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.Config.InternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		return nil, err
	}

	addr, err := internalSC.GetInnetPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &addr, err
}

func GetMarketRegWalletAddress(server *UserServer) (*common.Address, error) {
	internalSC, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.Config.InternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		return nil, err
	}
	addr, err := internalSC.GetMarketingRegWalletAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &addr, err
}
