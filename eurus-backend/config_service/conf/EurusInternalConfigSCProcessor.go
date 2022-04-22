package conf

import (
	"errors"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// func SetColdWalletAddressFromSC(server *ConfigServer, addr string) (string, error) {
// 	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(server.ServerConfig.EurusInternalConfigAddress), server.EthClient.Client)
// 	txHash := ""
// 	if err != nil {
// 		return txHash, err
// 	}
// 	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
// 	if err != nil {
// 		return txHash, err
// 	}

// 	tx, err := instance.SetPlatformWalletAddress(auth, common.HexToAddress(addr))
// 	if err != nil {
// 		return "", err
// 	}
// 	txHash = tx.Hash().Hex()
// 	return txHash, nil
// }

func AddPublicCurrencyInfoFromSC(server *ConfigServer, addr string, asset string) (string, error) {
	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(server.ServerConfig.EurusInternalConfigAddress), server.EthClient.Client)
	txHash := ""
	if err != nil {
		return txHash, err
	}
	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return txHash, err
	}
	tx, err := instance.AddCurrencyInfo(auth, common.HexToAddress(addr), asset)
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

// func DelPublicCurrencyInfoFromSC(server *ConfigServer, asset string) (string, error) {
// 	instance, err := contract.NewEurusInternalConfig(common.HexToAddress(server.ServerConfig.EurusInternalConfigAddress), server.EthClient.Client)
// 	txHash := ""
// 	if err != nil {
// 		return txHash, err
// 	}
// 	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
// 	if err != nil {
// 		return txHash, err
// 	}
// 	tx, err := instance.RemoveCurrencyInfo(auth, asset)
// 	if err != nil {
// 		return txHash, err
// 	}
// 	txHash = tx.Hash().Hex()
// 	return txHash, nil
// }

func GetPublicCurrencySCAddrFromSC(server *ConfigServer, address string) (string, error) {
	instance, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(server.ServerConfig.EurusInternalConfigAddress), server.EthClient.Client)
	asset, err := instance.GetErc20SmartContractByAddr(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return "", err
	}
	if asset == "" {
		return "", errors.New("No such asset")
	}
	return asset, nil
}

func GetPublicCurrencyNameByAddrFromSC(server *ConfigServer, address string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	asset, err := instance.GetErc20SmartContractByAddr(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return "", err
	}
	if asset == "" {
		return "", errors.New("No such asset")
	}
	return asset, nil
}
