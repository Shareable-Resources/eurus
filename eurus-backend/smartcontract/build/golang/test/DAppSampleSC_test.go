package test

import (
	"encoding/json"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestDAppSampleSC_DepositERC20ToDApp(t *testing.T) {

	ethClient := initEthClient(t)

	usdtAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	usdt, err := contract.NewEurusERC20(usdtAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	balance, _ := usdt.BalanceOf(&bind.CallOpts{}, common.HexToAddress(testOwnerPrivateKey))
	fmt.Println(balance)
	var extraData [32]byte
	transOpt.GasLimit = 1000000
	tx, err := usdt.DepositToDApp(transOpt, big.NewInt(1000000), dappAddr, extraData)
	if err != nil {
		t.Fatal(err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestDAppSample_GetWriter(t *testing.T) {
	ethClient := initEthClient(t)

	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	dappToken, err := contract.NewERC20(dappAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	addrList, err := dappToken.GetWriterList(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Writer list: ")
	for _, addr := range addrList {
		fmt.Println(addr)
	}
}

func TestDAppSample_SetWriter(t *testing.T) {
	ethClient := initEthClient(t)

	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	dappToken, err := contract.NewERC20(dappAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	externalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")

	extSC, err := contract.NewExternalSmartContractConfig(externalSCAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	assetNameList, assetAddrList, err := extSC.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for i, assetAddr := range assetAddrList {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey on asset: ", assetNameList[i], " error: ", err)
			continue
		}
		tx, err := dappToken.AddWriter(transOpt, assetAddr)
		if err != nil {
			fmt.Println("Add writer failed on asset: ", assetNameList[i], " error: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Cannot get receipt on asset: ", assetNameList[i], " error: ", err)
			continue
		}

		receiptData, _ := json.Marshal(receipt)
		fmt.Println(string(receiptData))
	}
}

func TestDAppSample_AddOwner(t *testing.T) {
	ethClient := initEthClient(t)

	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	dappToken, err := contract.NewERC20(dappAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	ownerAddr := "0xc316a94074fe5c387d6099e286e7f05a4dfc599f"

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := dappToken.AddOwner(transOpt, common.HexToAddress(ownerAddr))

	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestDAppSample_Refund(t *testing.T) {
	ethClient := initEthClient(t)

	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	dappToken, err := contract.NewDAppSample(dappAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	// ownerAddr := "0xc316a94074fe5c387d6099e286e7f05a4dfc599f"

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	transOpt.GasLimit = 1000000
	var extraData [32]byte
	tx, err := dappToken.Refund(transOpt, "USDT", big.NewInt(1000000), common.HexToAddress(testOwnerAddr), extraData)

	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestDAppSample_BalanceOf(t *testing.T) {
	ethClient := initEthClient(t)
	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	fmt.Println(dappAddr.Hex())
	dapp, _ := contract.NewDAppSample(dappAddr, ethClient.Client)
	balance, err := dapp.BalanceOf(&bind.CallOpts{}, common.HexToAddress(testOwnerAddr))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance)
	supply, err := dapp.TotalSupply(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(supply.String())
}

func TestDAppSample_SetExternalSmartContractConfig(t *testing.T) {
	ethClient := initEthClient(t)
	dappAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppSampleToken>")
	extAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")

	dapp, _ := contract.NewDAppSample(dappAddr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := dapp.SetExternalSmartContractConfig(transOpt, extAddr)
	if err != nil {
		t.Fatal(err)
	}

	queryEthReceipt(t, &ethClient, tx)
}
