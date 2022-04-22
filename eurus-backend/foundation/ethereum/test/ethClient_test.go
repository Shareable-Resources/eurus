package test

import (
	"fmt"
	"math/big"
	"testing"

	"eurus-backend/foundation/ethereum"

	"github.com/ethereum/go-ethereum/common"
)

func TestEthClient(t *testing.T) {
	ethClient := ethereum.EthClient{Protocol: "http",
		IP:      "13.228.80.104",
		Port:    8545,
		ChainID: big.NewInt(2018),
	}
	ethClient.Connect()
	/*
		if err != nil {
			t.Error("cannot set ethClient")
		}
		token, err := NewToken(common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40"), client)
		if err != nil {
			t.Errorf("Failed to instantiate a Token contract: %v", err)
		}
		name, err := token.Name(&bind.CallOpts{})
		if err != nil {
			t.Errorf("Failed to retrieve token name: %v", err)
		}
		t.Logf("Succeed to retrieve token name: %s", name)
	*/

	transaction, _ := ethClient.GetTransactionByHash("0x9561a2870707fbd91baa65b73b1f0fdccece2dea1fed41f24d19b510bb67d6b1")
	txByte, err := transaction.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to retrieve token name: %v", err)
	} else {
		fmt.Println(string(txByte))
	}
}

func TestReceipt(t *testing.T) {
	var ethClient ethereum.EthClient = ethereum.EthClient{
		Protocol: "https",
		IP:       "rinkeby.infura.io/v3/7d2a9aab75194b6c8597521994ade28d",
		Port:     443,
		ChainID:  big.NewInt(4),
	}
	ethClient.Connect()
	receipt, err := ethClient.GetConfirmedTransactionReceiptByHashStr("0x5b606b31fe7de14bdbf3017649ecd05ad348d32def8ae5fa2f1d2aad412f9de1")
	if err != nil {
		t.Errorf("Failed to retrieve token name: %v", err)
	}
	fmt.Println(len(receipt.Logs))
}

const targetAddress string = "0x4DfB6d6790054F3EB68324BC230E3104137CA8Db"

func TestTransferETH(t *testing.T) {
	var ethClient ethereum.EthClient = ethereum.EthClient{
		Protocol: "http",
		IP:       "13.228.80.104",
		Port:     8545,
		ChainID:  big.NewInt(2018),
	}
	ethClient.Connect()
	_, _, err := ethClient.TransferETH("5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99", targetAddress, big.NewInt(1000000))
	if err != nil {
		t.Errorf("Failed to Transfer ETH: %v", err)
	}
}

func TestTransferETHInEIP1559(t *testing.T) {
	var ethClient ethereum.EthClient = ethereum.EthClient{
		// Rinkeby testnet
		Protocol: "http",
		IP:       "54.254.253.134",
		Port:     8545,
		ChainID:  big.NewInt(4),
	}

	privKey := "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
	to := targetAddress
	toSmartContract := false
	amount := big.NewInt(1000000000)

	_, err := ethClient.Connect()
	if err != nil {
		t.Error(err)
	}

	_, tx, err := ethClient.TransferETHInEIP1559(privKey, to, amount, toSmartContract, false, false, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to Transfer ETH: %v", err)
	}

	t.Log(tx.Hash())
}

func TestSuggestGasFeeCap(t *testing.T) {
	var ethClient ethereum.EthClient = ethereum.EthClient{
		// Rinkeby testnet
		Protocol: "https",
		IP:       "ethnodes.eurus.network",
		Port:     443,
		ChainID:  big.NewInt(1),
	}

	_, err := ethClient.Connect()
	if err != nil {
		t.Error(err)
	}

	gasFeeCap, err := ethClient.SuggestGasFeeCap()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Gas fee cap: ", gasFeeCap.String())
}

func TestGetBalanceByBlockNumber(t *testing.T) {
	// var ethClient ethereum.EthClient = ethereum.EthClient{
	// 	// Rinkeby testnet
	// 	Protocol: "https",
	// 	IP:       "ethnodes.eurus.network",
	// 	Port:     443,
	// 	ChainID:  big.NewInt(1),
	// }
	ethClient := ethereum.EthClient{Protocol: "http",
		IP:      "13.228.80.104",
		Port:    8545,
		ChainID: big.NewInt(2018),
	}

	_, err := ethClient.Connect()
	if err != nil {
		t.Error(err)
	}

	var address common.Address
	b := [20]byte("0x4ccb4ef3027ac20bd1c20faed18084d940704c21")
	address = b
	// Address address = "0x4ccb4ef3027ac20bd1c20faed18084d940704c21"

	gasFeeCap, err := ethClient.GetBalanceByBlockNumber(address, "1010")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("balance: ", gasFeeCap.String())
}
