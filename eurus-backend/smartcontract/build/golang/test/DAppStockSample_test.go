package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestDAppStockSample_AddStock(t *testing.T) {
	ethClient := initEthClient(t)
	stockAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppStockSample>")
	stockInstance, err := contract.NewDAppStockSample(stockAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	usdtAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := stockInstance.AddWriter(transOpt, usdtAddr)
	if err != nil {
		t.Fatal(err)
	}

	queryEthReceipt(t, &ethClient, tx)

	transOpt, err = ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx1, err := stockInstance.UpdateProduct(transOpt, big.NewInt(1), big.NewInt(1000), big.NewInt(10000000), "Bear")
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx1)
}

func TestDAppStockSample_Purchase(t *testing.T) {

	ethClient := initEthClient(t)

	usdtAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	stockAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppStockSample>")

	usdtInstance, err := contract.NewEurusERC20(usdtAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt.GasLimit = 500000
	tx, err := usdtInstance.Purchase(transOpt, big.NewInt(1), big.NewInt(2), big.NewInt(20000000), stockAddr, [32]byte{})
	if err != nil {
		t.Fatal(err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestDAppStockSample_GetProductList(t *testing.T) {
	ethClient := initEthClient(t)
	stockAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<DAppStockSample>")
	stockInstance, err := contract.NewDAppStockSample(stockAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	productList, err := stockInstance.GetProductList(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, product := range productList {
		fmt.Println("Product Id: ", product.ProductId)
		fmt.Println("Name: ", product.Name)
		fmt.Println("On shelf: ", product.OnShelf)
		fmt.Println("Price: ", product.Price)

		stock, err := stockInstance.GetProductStock(&bind.CallOpts{}, product.ProductId)
		if err != nil {

			fmt.Println("err: ", err)
			continue
		}
		fmt.Println("Stock: ", stock.String())
	}

	list, _ := stockInstance.GetPurchaseList(&bind.CallOpts{}, big.NewInt(1))
	for _, item := range list {
		fmt.Println("Buyer: ", item.Buyer)
		fmt.Println("Quantity: ", item.Quantity)
		fmt.Println("Extra data: ", item.ExtraData)
	}
}
