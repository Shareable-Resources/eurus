package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestMultiOwnerable_SideChainAddOwner(t *testing.T) {

	newOwner := common.HexToAddress("")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	fmt.Println("Smart contract address: ", addr.Hex())
	multiOwnerable, _ := contract.NewMultiOwnable(addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := multiOwnerable.AddOwner(transOpt, newOwner)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestMultiOwnerable_TestnetAddOwner(t *testing.T) {

	newOwner := common.HexToAddress("")
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	fmt.Println("Smart contract address: ", addr.Hex())
	multiOwnerable, _ := contract.NewMultiOwnable(addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := multiOwnerable.AddOwner(transOpt, newOwner)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestMultiOwnerable_SideChainRemoveOwner(t *testing.T) {

	newOwner := common.HexToAddress("")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	fmt.Println("Smart contract address: ", addr.Hex())
	multiOwnerable, _ := contract.NewMultiOwnable(addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := multiOwnerable.RemoveOwner(transOpt, newOwner)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestMultiOwnerable_TestnetRemoveOwner(t *testing.T) {

	newOwner := common.HexToAddress("")
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	fmt.Println("Smart contract address: ", addr.Hex())
	multiOwnerable, _ := contract.NewMultiOwnable(addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := multiOwnerable.RemoveOwner(transOpt, newOwner)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}
