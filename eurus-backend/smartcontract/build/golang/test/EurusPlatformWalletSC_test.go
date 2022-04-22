package test

import (
	"encoding/json"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// //////THIS CONTRACT IS FOR TESTNET Rinkeby!!!!
func TestEurusPlatformWallet_TransferUSDT(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	platformWallet, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet ", err)
	}
	fmt.Println("OwnedUpgradeabilityProxy<EurusPlatformWallet> address: ", addr.String())
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error", err)
	}
	amount := big.NewInt(1)
	transOpt.GasLimit = 100000
	tx, err := platformWallet.Transfer(transOpt, common.HexToHash("0xe939d334b3a3c2e36055d066c7c6b60e3497ebbec3e4a71c788a4cc62454e20f"), common.HexToAddress("0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"), "USDT", amount, 6, nil)
	if err != nil {
		t.Fatal("TestRemoveCurrencyInfo: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusPlatformWallet_TransferETH(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	platformWallet, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet ", err)
	}
	fmt.Println("OwnedUpgradeabilityProxy<EurusPlatformWallet> address: ", addr.String())
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerAddr, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error", err)
	}
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("10", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}
	tx, err := platformWallet.Transfer(transOpt, common.HexToHash("0xe939d334b3a3c2e36055d066c7c6b60e3497ebbec3e4a71c788a4cc62454e20f"), common.HexToAddress("0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"), "ETH", amount, 18, nil)
	if err != nil {
		t.Fatal("Transfer: ", err)
	}
	fmt.Println("Trans Hash: ", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusePlatformWallet_GetETHBalance(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	amount, err := ethClient.GetBalance(addr)
	if err != nil {
		t.Fatal("GetBalance: ", err)
	}
	fmt.Println("ETH Balance: ", amount.String())
}

func TestEurusPlatformWalletSetWriter(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	platformWallet, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet ", err)
	}
	for _, withdrawObsAddr := range withdrawobserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error", err)
		}
		tx, err := platformWallet.AddWriter(transOpt, common.HexToAddress(withdrawObsAddr))
		if err != nil {
			t.Fatal("platformWallet.AddWriter error", err)
		}
		fmt.Println("tx Hash: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
		data, _ := json.Marshal(receipt)
		fmt.Println(string(data))
		fmt.Println()
	}
}

func TestEurusPlatformWallet_GetUSDTBalance(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	erc20Addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	erc20, err := contract.NewERC20(erc20Addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewERC20: ", err)
	}
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	balance, err := erc20.BalanceOf(&bind.CallOpts{From: addr}, addr)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error", err)
	}
	fmt.Println("Balance: ", balance.String())
}

func TestEurusPlatformWallet_SetEurusInternalConfigAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	walletSC, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error", err)
	}
	configAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	tx, err := walletSC.SetEurusInternalConfig(transOpt, configAddr)
	if err != nil {
		t.Fatal("SetEurusInternalConfig error", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusPlatformWallet_GetWriterList(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	walletSC, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet error: ", err)
	}

	writerList, err := walletSC.GetWriterList(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Writer List:")
	for _, writer := range writerList {
		fmt.Println(writer.Hex())
	}
}

func TestEurusPlatformWallet_GetInternalSCAddr(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	walletSC, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet error: ", err)
	}

	intcConfigAddr, err := walletSC.GetEurusInternalConfig(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(intcConfigAddr)
}

func TestEurusPlatformWallet_Upgrade(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	implAddr := getTestNetAddressBySmartContractName("EurusPlatformWallet")
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		t.Fatal("UpgradeTo error : ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusPlatformWallet_SetWriter(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	walletSC, err := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewEurusPlatformWallet error: ", err)
	}

	for _, writer := range withdrawobserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error", err)
		}
		tx, err := walletSC.AddWriter(transOpt, common.HexToAddress(writer))
		if err != nil {
			fmt.Printf("Add writer %s error: %s\r\n", writer, err.Error())
			continue
		} else {
			fmt.Printf("Add writer %s, trans hash: %s\r\n", writer, tx.Hash().Hex())
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error for withrawobserver: ", writer, " Error: ", err.Error())
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Receipt  for withdraw observer: ", writer, " receipt: ", string(receiptJson))
		}
	}

}

func TestEurusPlatformWallet_GetAllBalance(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	for _, assetName := range erc20Currency {
		if assetName == "ETH" {
			continue
		}
		tokenAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
		erc20, err := mainnet_contract.NewERC20(tokenAddr, ethClient.Client)
		if err != nil {
			t.Fatal(assetName + ": " + err.Error())
		}

		balance, err := erc20.BalanceOf(&bind.CallOpts{}, addr)
		if err != nil {
			t.Fatal(assetName + ": " + err.Error())
		}

		fmt.Println("Asset: ", assetName, ": ", balance.String())
	}
}

func TestEurusPlatformWallet_IsWriter(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	mainnetInstance, _ := mainnet_contract.NewEurusPlatformWallet(addr, ethClient.Client)
	isValid, err := mainnetInstance.IsWriter(&bind.CallOpts{}, common.HexToAddress("0xa345973fa94e786dc483e96fa73c5fa232d11858"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("isValid :", isValid)
}
