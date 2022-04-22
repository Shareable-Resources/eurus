package test

import (
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestGeneralMultiSigWallet_SubmitTranaction(t *testing.T) {
	var ownerPrivateKey1 = "64f8ba795cf8f78e9c3c7a1b154326ba6e0e6f994e4853f0a551c15519fb438e"
	addr := getAddressBySmartContractName("UnitTest")
	ethClient := initEthClient(t)
	unitTest, err := contract.NewUnitTest(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivateKey1, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	transOpt.NoSend = true
	val := big.NewInt(10)
	tx, err := unitTest.SetValue(transOpt, val)
	if err != nil {
		t.Fatal(err)
	}

	ownerWalletAddr := getAddressBySmartContractName("GeneralMultiSigWallet")
	ownerWallet, err := mainnet_contract.NewGeneralMultiSigWallet(ownerWalletAddr, ethClient.Client)
	transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivateKey1, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	tx1, err := ownerWallet.SubmitTransaction(transOpt1, addr, big.NewInt(0), tx.Data())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Trans hash: ", tx1.Hash().Hex())
	receipt, err := ethClient.QueryEthReceipt(tx1)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Status == 0 {
		receiptData, _ := json.Marshal(receipt)
		t.Fatal("Receipt status is 0: ", string(receiptData))
	}

	ethereum.DefaultABIDecoder.ImportABIJson("GeneralMultiSigWallet", mainnet_contract.GeneralMultiSigWalletABI)

	abi := ethereum.DefaultABIDecoder.GetABI("GeneralMultiSigWallet")

	for _, receiptLog := range receipt.Logs {
		if receiptLog.Topics[0] == abi.Events["Submission"].ID {
			fmt.Println("Submission ID: ", receiptLog.Topics[1].Hex())
			break
		}
	}
}

func TestGeneralMultiSigWallet_Confirmation(t *testing.T) {
	ethClient := initEthClient(t)

	ownerWalletAddr := getAddressBySmartContractName("GeneralMultiSigWallet")
	ownerWallet, err := mainnet_contract.NewGeneralMultiSigWallet(ownerWalletAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	ownerPrivateKey2 := "8497f70e6e69704a68f9a598f738c13198acf3ea4d8059b94c3fc678c2049211"
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivateKey2, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	submissionId := big.NewInt(1)
	tx, err := ownerWallet.ConfirmTransaction(transOpt, submissionId)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestGeneralMultiSigWallet_GetOwnerList(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	ownerWalletAddr := getTestNetAddressBySmartContractName("GeneralMultiSigWallet")
	fmt.Println("Chain id: ", ethClient.ChainID.String())
	fmt.Println("Smart contract address: ", ownerWalletAddr.Hex())
	ownerWallet, err := mainnet_contract.NewGeneralMultiSigWallet(ownerWalletAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	callOpts := &bind.CallOpts{
		// BlockNumber: big.NewInt(0x9238c7),
	}

	addrList, err := ownerWallet.GetOwners(callOpts)
	for _, addr := range addrList {
		fmt.Println("Address: ", addr.Hex())
	}

}
