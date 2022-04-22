package test

import (
	"context"
	"encoding/json"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestMarketingWallet_RegAddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")
	wallet, err := contract.NewMarketingWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	for _, depositObserver := range depositObserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " for deposit observer: ", depositObserver)
			continue
		}
		tx, err := wallet.AddWalletOperator(transOpt, common.HexToAddress(depositObserver))
		if err != nil {
			fmt.Println("AddWriter error: ", err, " for deposit observer: ", depositObserver)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " for deposit observer: ", depositObserver)
			continue
		}

		receiptData, _ := json.Marshal(receipt)
		fmt.Println("Receipt: ", string(receiptData))
	}

	for _, userServerAddr := range userServerHDWalletAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " for user server: ", userServerAddr)
			continue
		}
		tx, err := wallet.AddWalletOperator(transOpt, common.HexToAddress(userServerAddr))
		if err != nil {
			fmt.Println("AddWriter error: ", err, " for user server: ", userServerAddr)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " for user server: ", userServerAddr)
			continue
		}

		receiptData, _ := json.Marshal(receipt)
		fmt.Println("Receipt: ", string(receiptData))
	}

	for _, blockChainIndexerAddr := range blockChainIndexerHDWalletAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " for block chain indexer: ", blockChainIndexerAddr)
			continue
		}
		tx, err := wallet.AddWalletOperator(transOpt, common.HexToAddress(blockChainIndexerAddr))
		if err != nil {
			fmt.Println("AddWriter error: ", err, " for block chain indexer: ", blockChainIndexerAddr)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " for block chain indexer: ", blockChainIndexerAddr)
			continue
		}

		receiptData, _ := json.Marshal(receipt)
		fmt.Println("Receipt: ", string(receiptData))
	}

}

func TestMarketingWallet_RegFundingWallet(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")

	_, tx, err := ethClient.TransferETH(ownerPrivKey, addr.Hex(), big.NewInt(6900000000000000000))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}
func TestMarketingWallet_GetWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")
	wallet, err := contract.NewMarketingWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	walletOperatorList, err := wallet.GetWalletOperatorList(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, addr := range walletOperatorList {
		fmt.Println(addr.Hex())
	}
}

func TestMarketingWallet_GetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	impl, err := proxy.Implementation(&bind.CallOpts{})
	fmt.Println("Implementation: ", impl.Hex())
}

func TestMarketingWallet_RegWalletBalance(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")

	amount, err := ethClient.Client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("EUN Amount: ", amount.String())
}

func TestMarketingWallet_MigrateEUN(t *testing.T) {
	ethClient := initEthClient(t)
	oldAddr := getAddressBySmartContractName("Old_OwnedUpgradeabilityProxy<MarketingRegWallet>")
	newAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")

	fmt.Println("Old wallet addr: ", oldAddr.Hex())
	fmt.Println("New wallet addr: ", newAddr.Hex())

	marketingWallet, err := contract.NewMarketingWallet(oldAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	balance, err := ethClient.GetBalance(oldAddr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Old Balance: ", balance)
	estimateBalance := balance
	transOpt.NoSend = true
	tx, err := marketingWallet.TransferETH(transOpt, newAddr, estimateBalance)
	if err != nil {
		t.Fatal(err)
	}
	gasLimit := big.NewInt(0)
	gasLimit.SetUint64(tx.Gas())

	eunUsed := tx.GasPrice().Mul(tx.GasPrice(), gasLimit)
	final := balance.Sub(balance, eunUsed)
	fmt.Println("Final balance: ", final.String())
	transOpt.NoSend = false
	tx1, err := marketingWallet.TransferETH(transOpt, newAddr, final)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx1)

	newBalance, err := ethClient.GetBalance(newAddr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("New wallet balance: ", newBalance.String())

}

func TestMarketingWallet_RegFundingTransferAway(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")
	marketingWallet, err := contract.NewMarketingWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	transOpt.GasLimit = 1000000
	tx, err := marketingWallet.TransferETH(transOpt, common.HexToAddress(testOwnerAddr), big.NewInt(500000000000000000))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInit_MarketRegWallet(t *testing.T) {

	TestMarketingWallet_RegAddWalletOperator(t)
	TestMarketingWallet_RegFundingWallet(t)
}
