package main

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func NewERC20(address string, client *ethereum.EthClient) (*contract.EurusERC20, error) {
	var err error
	instance, err := contract.NewEurusERC20(common.HexToAddress(address), client.Client)
	if err != nil {
		return nil, err
	}

	return instance, nil

}

func MintToken(client *ethereum.EthClient, instance *contract.EurusERC20, privateKeyHex string, to string, amount *big.Int) (string, error) {
	var err error
	txHash := ""
	auth, err := client.GetNewTransactorFromPrivateKey(privateKeyHex, client.ChainID)
	if err != nil {
		return "", err
	}
	tx, err := instance.Mint(auth, common.HexToAddress(to), amount)
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

func TransferToken(client *ethereum.EthClient, instance *contract.EurusERC20, privateKeyHex string, to string, amount *big.Int) (string, error) {
	var err error
	txHash := ""
	auth, err := client.GetNewTransactorFromPrivateKey(privateKeyHex, client.ChainID)
	if err != nil {
		return "", err
	}
	tx, err := instance.Transfer(auth, common.HexToAddress(to), amount)
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}
