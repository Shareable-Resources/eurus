package bc_indexer

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func GetMainNetAssetList(eurusInternalSCAddress string, ethClient *ethereum.EthClient) ([]string, []common.Address, error) {
	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(eurusInternalSCAddress), ethClient.Client)
	if err != nil {
		return nil, nil, err
	}
	assetName, assetAddr, err := instance.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		return nil, nil, err
	}
	return assetName, assetAddr, nil
}

func GetEurusUserDepositAddress(eurusInternalSCAddress string, ethClient *ethereum.EthClient) (*common.Address, error) {
	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(eurusInternalSCAddress), ethClient.Client)
	if err != nil {
		return nil, err
	}

	addr, err := instance.EurusUserDepositAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func GetMainnetPlatformWalletAddress(eurusInternalSCAddress string, ethClient *ethereum.EthClient) (*common.Address, error) {

	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(eurusInternalSCAddress), ethClient.Client)
	if err != nil {
		return nil, err
	}

	addr, err := instance.PlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &addr, nil
}
