package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var ownerPrivKey string = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
var userPrivKey string = "6c4298fd60836228304d19ba21268ea58a1e5a708c8c3986560d17d864ae939f"
var ownerAddr common.Address = common.HexToAddress("0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93")
var userAddr common.Address = common.HexToAddress("0x39EB6463871040f75C89C67ec1dFCB141C3da1cf")
var toAddr common.Address = common.HexToAddress("0x39eb6463871040f75c89c67ec1dfcb141c3da1cf")

func TestWrappedEUN_Implementation(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	proxy, err := contract.NewReceiveFallbackOwnedUpgradeabilityProxy(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewReceiveFallbackOwnedUpgradeabilityProxy:", err)
	}

	addr, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Error("Implementation:", err)
	}

	fmt.Println(addr)
}

func TestWrappedEUN_Upgrade(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	proxy, err := contract.NewReceiveFallbackOwnedUpgradeabilityProxy(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewReceiveFallbackOwnedUpgradeabilityProxy:", err)
	}

	tx, err := proxy.UpgradeTo(trans, common.HexToAddress("0x79D1d66231026e09957aADeB41772790EBbe44ee"))
	if err != nil {
		t.Error("UpgradeTo:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_Info(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	weun, err := contract.NewERC20(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	name, err := weun.Name(&bind.CallOpts{})
	if err != nil {
		t.Error("Name:", err)
	}

	fmt.Println(name)

	symbol, err := weun.Symbol(&bind.CallOpts{})
	if err != nil {
		t.Error("Symbol:", err)
	}

	fmt.Println(symbol)

	decimals, err := weun.Decimals(&bind.CallOpts{})
	if err != nil {
		t.Error("Decimals:", err)
	}

	fmt.Println(decimals)
}

func TestWrappedEUN_TotalSupply(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	totalSupply, err := weun.TotalSupply(&bind.CallOpts{})
	if err != nil {
		t.Error("TotalSupply:", err)
	}

	fmt.Println(totalSupply)
}

func TestWrappedEUN_GetBalance(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	eun, err := ethClient.GetBalance(userAddr)
	if err != nil {
		t.Error("GetBalance:", err)
	}

	fmt.Println("EUN:", eun)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	balance, err := weun.BalanceOf(&bind.CallOpts{}, userAddr)
	if err != nil {
		t.Error("BalanceOf:", err)
	}

	fmt.Println(balance)
}

func TestWrappedEUN_Allowance(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	allowed, err := weun.Allowance(&bind.CallOpts{}, userAddr, ownerAddr)
	if err != nil {
		t.Error("Allowance:", err)
	}

	fmt.Println(allowed)
}

func TestWrappedEUN_Approve(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	allowance := big.NewInt(-1)
	// allowance.SetString("1000000000000000000000", 10)
	tx, err := weun.Approve(trans, ownerAddr, allowance)
	if err != nil {
		t.Error("Approve:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_Deposit(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	trans.Value = big.NewInt(0)
	trans.Value.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	tx, err := weun.Deposit(trans)
	if err != nil {
		t.Error("Deposit:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_DepositTo(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	trans.Value = big.NewInt(0)
	trans.Value.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Deposit to 0 = deposit to myself
	// tx, err := weun.DepositTo(trans, common.HexToAddress("0x0000000000000000000000000000000000000000"))
	tx, err := weun.DepositTo(trans, toAddr)
	if err != nil {
		t.Error("DepositTo:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_DepositBySend(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	amount := big.NewInt(0)
	amount.SetString("1000000000000000000", 10)

	_, tx, err := ethClient.TransferETHToSmartContract(userPrivKey, contractAddr.Hex(), amount, false, nil, nil)
	if err != nil {
		t.Error("TransferETHToSmartContract:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_Transfer(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	amount := big.NewInt(0)
	amount.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Transfer all
	amount, err = weun.BalanceOf(&bind.CallOpts{}, userAddr)
	if err != nil {
		t.Error("BalanceOf:", err)
	}

	tx, err := weun.Transfer(trans, common.HexToAddress("0x0"), amount)
	if err != nil {
		t.Error("Transfer:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_TransferFrom(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	amount := big.NewInt(0)
	amount.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Transfer all
	// amount, err = weun.BalanceOf(&bind.CallOpts{}, userAddr)
	// if err != nil {
	// 	t.Error("BalanceOf:", err)
	// }

	tx, err := weun.TransferFrom(trans, userAddr, common.HexToAddress("0x0"), amount)
	if err != nil {
		t.Error("TransferFrom:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_Withdraw(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	amount := big.NewInt(0)
	amount.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Withdraw all
	// amount, err = weun.BalanceOf(&bind.CallOpts{}, userAddr)
	// if err != nil {
	// 	t.Error("BalanceOf:", err)
	// }

	tx, err := weun.Withdraw(trans, amount)
	if err != nil {
		t.Error("Withdraw:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_WithdrawTo(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(userPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	amount := big.NewInt(0)
	amount.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Withdraw all
	// amount, err = weun.BalanceOf(&bind.CallOpts{}, userAddr)
	// if err != nil {
	// 	t.Error("BalanceOf:", err)
	// }

	// Withdraw to 0 = withdraw to myself
	// tx, err := weun.WithdrawTo(trans, common.HexToAddress("0x0000000000000000000000000000000000000000"), amount)
	tx, err := weun.WithdrawTo(trans, toAddr, amount)
	if err != nil {
		t.Error("WithdrawTo:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWrappedEUN_WithdrawFrom(t *testing.T) {
	ethClient := initEthClient(t)
	contractAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WEUN>")

	trans, err := ethClient.GetNewTransactorFromPrivateKey(ownerPrivKey, ethClient.ChainID)
	if err != nil {
		t.Error("GetNewTransactorFromPrivateKey:", err)
	}

	amount := big.NewInt(0)
	amount.SetString("1000000000000000", 10)

	weun, err := contract.NewWrappedEUN(contractAddr, ethClient.Client)
	if err != nil {
		t.Error("NewWrappedEUN:", err)
	}

	// Withdraw all
	amount, err = weun.BalanceOf(&bind.CallOpts{}, userAddr)
	if err != nil {
		t.Error("BalanceOf:", err)
	}

	tx, err := weun.WithdrawFrom(trans, userAddr, userAddr, amount)
	if err != nil {
		t.Error("WithdrawFrom:", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}
