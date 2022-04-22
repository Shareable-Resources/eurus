package bc_indexer

import (
	"errors"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func GetCurrencyNameByAddrFromSC(externalSmartContractAddress string, ethClient *ethereum.EthClient, addr string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(externalSmartContractAddress), ethClient.Client)
	if err != nil {
		return "", err
	}
	asset, err := instance.GetErc20SmartContractByAddr(&bind.CallOpts{}, common.HexToAddress(addr))
	if err != nil {
		return "", err
	}
	if asset == "" {
		return "", errors.New("No such asset")
	}
	return asset, nil
}

func GetAssetList(externalSmartContractAddress string, ethClient *ethereum.EthClient)([]string, []common.Address,error){
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(externalSmartContractAddress), ethClient.Client)
	if err != nil {
		return nil, nil,err
	}
	assetName, assetAddr,err:=instance.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		return nil, nil, err
	}
	return assetName,assetAddr,nil
}

