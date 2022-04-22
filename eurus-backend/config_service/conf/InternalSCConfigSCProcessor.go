package conf

import (
	"eurus-backend/smartcontract/build/golang/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func SetMainnetWalletAddressFromSC(server *ConfigServer, addr string) (string, error) {
	instance, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.ServerConfig.InternalSCConfigAddress), server.EthClient.Client)
	txHash := ""
	if err != nil {
		return txHash, err
	}
	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return txHash, err
	}

	tx, err := instance.SetMainnetWalletAddress(auth, common.HexToAddress(addr))
	if err != nil {
		return "", err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

func GetMainnetWalletAddressFromSC(server *ConfigServer) (string, error) {
	instance, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.ServerConfig.InternalSCConfigAddress), server.EthClient.Client)
	address, err := instance.GetMainnetPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

func SetInnetWalletAddressFromSC(server *ConfigServer, addr string) (string, error) {
	instance, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.ServerConfig.InternalSCConfigAddress), server.EthClient.Client)
	txHash := ""
	if err != nil {
		return txHash, err
	}
	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return txHash, err
	}
	tx, err := instance.SetInnetWalletAddress(auth, common.HexToAddress(addr))
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

func GetInnetWalletAddressFromSC(server *ConfigServer) (string, error) {
	instance, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.ServerConfig.InternalSCConfigAddress), server.EthClient.Client)
	address, err := instance.GetInnetPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}
