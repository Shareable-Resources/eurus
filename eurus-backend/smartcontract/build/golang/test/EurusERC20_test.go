package test

import (
	"encoding/json"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestEurusERC20_Init(t *testing.T) {

	amount := big.NewInt(0)

	initEurusERC20(t, "USDT", amount, 6)

	initEurusERC20(t, "USDC", amount, 6)

	initEurusERC20(t, "LINK", amount, 18)

	initEurusERC20(t, "UNI", amount, 18)

	initEurusERC20(t, "BNB", amount, 18)

	initEurusERC20(t, "BUSD", amount, 18)

	initEurusERC20(t, "YFI", amount, 18)

	initEurusERC20(t, "DAI", amount, 18)

	initEurusERC20(t, "OMG", amount, 18)

	initEurusERC20(t, "VEN", amount, 18)

	initEurusERC20(t, "AAVE", amount, 18)

	initEurusERC20(t, "HT", amount, 18)

	initEurusERC20(t, "SUSHI", amount, 18)

	initEurusERC20(t, "TUSD", amount, 18)

	initEurusERC20(t, "cDAI", amount, 8)

	initEurusERC20(t, "SXP", amount, 18)

	initEurusERC20(t, "BAT", amount, 18)

	initEurusERC20(t, "USDK", amount, 18)

	initEurusERC20(t, "WBTC", amount, 8)

	initEurusERC20(t, "ZIL", amount, 12)

	initEurusERC20(t, "SNX", amount, 18)

	initEurusERC20(t, "OKB", amount, 18)

	initEurusERC20(t, "BAND", amount, 18)

	initEurusERC20(t, "HUSD", amount, 8)

	initEurusERC20(t, "MKR", amount, 18)

	initEurusERC20(t, "ZRX", amount, 18)

	initEurusERC20(t, "PAX", amount, 18)

	initEurusERC20(t, "COMP", amount, 18)

	initEurusERC20(t, "RSR", amount, 18)

	initEurusERC20(t, "BAL", amount, 18)

	initEurusERC20(t, "ETH", amount, 18)

	initEurusERC20(t, "PLA", amount, 6)

}

func initEurusERC20(t *testing.T, assetName string, supply *big.Int, decimals uint8) {
	ethClient := initEthClient(t)
	erc20Instance, err := contract.NewEurusERC20(getAddressBySmartContractName("OwnedUpgradeabilityProxy<"+assetName+">"), ethClient.Client)
	if err != nil {
		t.Fatal("New ERRC20 error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	intcADDR := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	extcADDR := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	tx, err := erc20Instance.Init0(transOpt, intcADDR, assetName, assetName, supply, decimals, extcADDR)
	if err != nil {
		fatalSmartContractError(t, "Init Error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)

	// InitExternalSCConfigERC20KycLimit(t, assetName, []string{"0", "1", "2", "3"}, []*big.Int{big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(1000000000000000000)})
}

func TestEurusERC20_Mint(t *testing.T) {

	ethClient := initEthClient(t)
	var targetAddr string = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}

	assetList, addressList, err := externalSC.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		t.Fatal("GetAssetAddress error: ", err)
	}

	for index, assetName := range assetList {
		if assetName != "DAI" {
			continue
		}
		addr := addressList[index]
		erc20SC, err := contract.NewEurusERC20(addr, ethClient.Client)
		if err != nil {
			t.Fatal("erc20 smart contract error: ", err)
		}

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey: ", err)
		}

		decimalPt, err := externalSC.GetAssetDecimal(&bind.CallOpts{}, assetName)
		if err != nil {
			fmt.Println("Unable to get decimal: ", err, " assetName: ", assetName)
			continue
		}
		transOpt.GasLimit = 1000000000
		amount := big.NewInt(10000)
		ex := int64(math.Pow10(int(decimalPt.Int64())))
		amount = amount.Mul(amount, big.NewInt(ex))
		fmt.Println(assetName, " mint: ", amount)
		tx, err := erc20SC.Mint(transOpt, common.HexToAddress(targetAddr), amount)
		if err != nil {
			fatalSmartContractError(t, "Mint Error: ", err)
		}
		fmt.Println(assetName)
		queryEthReceipt(t, &ethClient, tx)
	}

}

func TestEurusERC20_Transfer(t *testing.T) {
	ethClient := initEthClient(t)
	transferToAddr := common.HexToAddress("0x4DfB6d6790054F3EB68324BC230E3104137CA8Db")
	transferAmount := big.NewInt(100000000)
	addr := getAddressBySmartContractName("ExternalSmartContractConfig")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	erc20Addr, err := externalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, "USDT")
	if err != nil {
		t.Fatal("GetAssetAddress error: ", err)
	}

	erc20SC, err := contract.NewEurusERC20(erc20Addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	tx, err := erc20SC.Transfer(transOpt, transferToAddr, transferAmount)
	if err != nil {
		fatalSmartContractError(t, "Transfer Error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusERC20_GetDeciamls(t *testing.T) {
	ethClient := initEthClient(t)

	for _, currency := range erc20Currency {
		erc20Instance, err := contract.NewTestERC20(getAddressBySmartContractName("OwnedUpgradeabilityProxy<"+currency+">"), ethClient.Client)
		if err != nil {
			fmt.Printf("Unable to NewTestERC20 for asset %s. Error: %s\r\n", "OwnedUpgradeabilityProxy<"+currency+">", err.Error())
			return
		}
		decimal, _ := erc20Instance.Decimals(&bind.CallOpts{})
		fmt.Println(currency+" decimal: ", decimal)
	}

}

func EurusERC20_AddWriterImpl(t *testing.T, erc20AssetName string) {
	ethClient := initEthClient(t)

	erc20Instance, err := contract.NewEurusERC20(getAddressBySmartContractName("OwnedUpgradeabilityProxy<"+erc20AssetName+">"), ethClient.Client)
	if err != nil {
		fmt.Printf("Unable to NewTestERC20 for asset %s. Error: %s\r\n", erc20AssetName, err.Error())
		return
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Printf("Unable to GetNewTransactorFromPrivateKey for asset %s. Error: %s\r\n", erc20AssetName, err.Error())
		return
	}

	platformWalletAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	tx, err := erc20Instance.AddWriter(transOpt, platformWalletAddr)
	if err != nil {
		fmt.Println("Failed to add writer to ", erc20AssetName, ". Error: ", err.Error())
	} else {
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Printf("Failed to query receipt for asset: %s. Error: %s\r\n", erc20AssetName, err.Error())
		} else {
			fmt.Printf("Asset: %s Receipt: %v\r\n", erc20AssetName, receipt)
		}
	}

	transOpt, err = ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Printf("Unable to GetNewTransactorFromPrivateKey for asset %s. Error: %s\r\n", erc20AssetName, err.Error())
		return
	}
	withdrawSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	tx, err = erc20Instance.AddWriter(transOpt, withdrawSCAddr)
	if err != nil {
		fmt.Println("Failed to add WithdrawSmartContract address to ", erc20AssetName, ". Error: ", err.Error())
	} else {
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Printf("Failed to query receipt for adding WithdrawSmartContract address for asset: %s. Error: %s\r\n", erc20AssetName, err.Error())
		} else {
			fmt.Printf("Asset name: %s, Receipt: %v\r\n", erc20AssetName, receipt)
		}
	}
}

func TestEurusERC20_AddWriter(t *testing.T) {
	for _, assetName := range erc20Currency {
		EurusERC20_AddWriterImpl(t, assetName)
	}
}

func TestEurusERC20_QueryBalance(t *testing.T) {
	for _, currency := range erc20Currency {
		EurusERC20QueryBalance(t, currency)
	}
}

func TestERC20_SideChainQueryEUN(t *testing.T) {
	for _, currency := range erc20Currency {
		ERC20_SideChainQueryEUNBalance(t, currency)
	}
}

func ERC20_SideChainQueryEUNBalance(t *testing.T, assetName string) {
	ethClient := initEthClient(t)

	erc20Addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")

	balance, err := ethClient.GetBalance(erc20Addr)
	if err != nil {
		fmt.Println("Unable to query ERC20 EUN balance: ", err)
	} else {
		fmt.Println("Balance of ", assetName, ": ", balance.String())
	}
}

func EurusERC20QueryBalance(t *testing.T, assetName string) {
	ethClient := initEthClient(t)

	erc20Instance, err := contract.NewTestERC20(getAddressBySmartContractName("OwnedUpgradeabilityProxy<"+assetName+">"), ethClient.Client)
	if err != nil {
		fmt.Printf("Unable to NewTestERC20 for asset %s. Error: %s\r\n", assetName, err.Error())
		return
	}

	balance, err := erc20Instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0x4ffb69cc66fc6f38cac81762f548cbedf35d77c2"))
	if err != nil {
		fmt.Printf("Unable to get balance for asset %s. Error: %s\r\n", assetName, err.Error())
	} else {
		fmt.Println("Account: 0x4ffb69cc66fc6f38cac81762f548cbedf35d77c2")
		fmt.Printf("%s Balance: %d\r\n", assetName, balance.Uint64())
	}
}

// func TestEurusERC20_AddBlacklist(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<BAL>")
// 	erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
// 	if err != nil {
// 		t.Fatal("EUR20 ", err)
// 	}
// 	keyPair, err := ethereum.GetEthKeyPair(testOwnerPrivateKey)
// 	if err != nil {
// 		t.Fatal("Get Key Pair Error: ", err)
// 	}
// 	fmt.Println("Test Owner Address: " + keyPair.Address.Hex())
// 	owners, err := erc20.GetOwners(&bind.CallOpts{})
// 	if err != nil {
// 		t.Fatal("Get Owners error: ", err)
// 	}

// 	for index, owner := range owners {
// 		fmt.Println("Owner [" + strconv.Itoa(index) + "]: " + owner.Hex())
// 	}

// 	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
// 	//transOpt.GasLimit=100000
// 	if err != nil {
// 		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
// 	}

// 	tx, err := erc20.AddBlackListDestAddress(transOpt, common.HexToAddress(approvalObserverAddr[0]))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	queryEthReceipt(t, &ethClient, tx)

// }

func TestEurusERC20_GetBlacklist(t *testing.T) {
	ethClient := initEthClient(t)
	for _, assetName := range erc20Currency {
		eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")

		eurusDepositUser := getAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")

		erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
		if err != nil {
			fmt.Println("EUR20 ", err)
		}
		tx, err := erc20.BlackListDestAddressMap(&bind.CallOpts{}, eurusDepositUser)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Currency : ", assetName, " BlackListDestAddress: ", tx)
	}

}

func TestEurusERC20_RemoveBlacklist(t *testing.T) {
	ethClient := initEthClient(t)
	eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<BAL>")
	eurusDepositUser := getAddressBySmartContractName("OwnedUpgradeabilityProxy<BAL>")
	erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
	if err != nil {
		t.Fatal("Eurus EUR20 ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 100000
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	tx, err := erc20.RemoveBlackListDestAddress(transOpt, eurusDepositUser)
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusERC20_AddEurusUserDepositToAllEurusERC20SCBlacklist(t *testing.T) {
	ethClient := initEthClient(t)
	for _, currency := range erc20Currency {
		eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + currency + ">")
		eurusUserDepositAddr := common.HexToAddress("OwnedUpgradeabilityProxy<EurusUserDeposit>")
		erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
		if err != nil {
			t.Fatal("EUR20 ", currency, " ", err)
		}
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey: ", currency, " ", err)
		}

		tx, err := erc20.AddBlackListDestAddress(transOpt, eurusUserDepositAddr)
		if err != nil {
			fmt.Println("currency: ", currency, " error: ", err)
			continue
		}
		//fmt.Println("Currency: " , currency , " Hash : ", tx.Hash().Hex())
		fmt.Println(tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("currency: ", currency, " error: ", err)
			continue
		}
		data, _ := json.Marshal(receipt)
		fmt.Println("currency: ", currency, " receipt: ", string(data))
	}
}

func TestEurusERC20_TransferFundingToContract(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("100000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, data := range erc20Currency {
		eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + data + ">")
		fmt.Println(eurusERC20.String())
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, eurusERC20.String(), amount)
		if err != nil {
			fmt.Println(err)
			continue
		}

		queryEthReceipt(t, &ethClient, tx)

		balance, err := ethClient.GetBalance(eurusERC20)
		if err != nil {
			fmt.Println(data, " get balance error: ", err)
			continue
		}

		fmt.Println(data, ": ", balance.String())
	}
}

func TestEurusERC20_GetOwnerList(t *testing.T) {
	ethClient := initEthClient(t)
	for _, asset := range erc20Currency {
		eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + asset + ">")
		erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
		if err != nil {
			t.Fatal("EurusEUR20 ", err)
		}
		tx, err := erc20.GetOwners(&bind.CallOpts{})
		for _, addr := range tx {
			fmt.Println(asset, " ", addr.Hex())
		}
	}
}

func TestEurusERC20_RemoveOwner(t *testing.T) {
	ethClient := initEthClient(t)

	eurusERC20 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	erc20, err := contract.NewEurusERC20(eurusERC20, ethClient.Client)
	if err != nil {
		t.Fatal("EUR20 ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	tx, err := erc20.RemoveOwner(transOpt, common.HexToAddress("0x17c09ddD7466083D1D7B389Fe22D48a0e21856Ba"))
	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusERC20_UpdateProxy(t *testing.T) {
	ethClient := initEthClient(t)
	assetName := "USDT"
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 10000000
	implAddr := getAddressBySmartContractName(assetName)
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(assetName+": ", implAddr.Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusERC20_SetInternalSCConfigAddress(t *testing.T) {
	ethClient := initEthClient(t)
	for _, data := range erc20Currency {
		addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + data + ">")
		fmt.Println("Currency: ", data, " proxy addr: ", addr.Hex())
		curr, err := contract.NewEurusERC20(addr, ethClient.Client)
		if err != nil {
			fmt.Println(err)
		}
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		internalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
		tx, err := curr.SetInternalSCConfigAddress(transOpt, internalSCAddr)
		if err != nil {
			fmt.Println(err)
		}

		queryEthReceipt(t, &ethClient, tx)

	}
}

func TestEurusERC20_GetUserDailyWithdrawLimit(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	erc20, err := contract.NewEurusERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	withdrewAmount, err := erc20.DailyWithdrewAmount(&bind.CallOpts{}, common.HexToAddress("0x4ffb69cc66fc6f38cac81762f548cbedf35d77c2"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Withdrew amount: ", withdrewAmount.String())

}

func TestEurusERC20_GetWriterList(t *testing.T) {
	ethClient := initEthClient(t)
	for _, currency := range erc20Currency {
		addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + currency + ">")
		fmt.Println(currency, " proxy addr: ", addr.Hex())
		curr, err := contract.NewEurusERC20(addr, ethClient.Client)
		if err != nil {
			fmt.Println("err: ", err, " currency: ", currency)
			continue
		}
		writerList, err := curr.GetWriterList(&bind.CallOpts{})
		if err != nil {
			fmt.Println("err: ", err, " currency: ", currency)
			continue
		}
		fmt.Println("Currency: ", currency)
		for _, writer := range writerList {
			fmt.Println("Writer: ", writer.Hex())
		}
	}
}

func TestEurusERC20_GetSideChainBalance(t *testing.T) {
	ethClient := initEthClient(t)
	walletAddr := common.HexToAddress("0x99B32dAD54F630D9ED36E193Bc582bbed273d666")
	for _, currency := range erc20Currency {
		addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + currency + ">")
		fmt.Println(currency, " proxy addr: ", addr.Hex())
		curr, err := contract.NewEurusERC20(addr, ethClient.Client)
		if err != nil {
			fmt.Println(currency, " error: ", err)
			continue
		}
		balance, err := curr.BalanceOf(&bind.CallOpts{}, walletAddr)
		if err != nil {
			fmt.Println(currency, " error: ", err)
			continue
		}
		fmt.Println(currency, " balance: ", balance.String())
	}
}

func TestEurusERC20_GetInternalSmartContractAddress(t *testing.T) {
	ethClient := initEthClient(t)

	for _, currency := range erc20Currency {
		addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + currency + ">")
		fmt.Println(currency, " proxy addr: ", addr.Hex())
		curr, err := contract.NewEurusERC20(addr, ethClient.Client)
		if err != nil {
			fmt.Println(currency, " error: ", err)
			continue
		}
		internalSCAddr, err := curr.GetInternalSCConfigAddress(&bind.CallOpts{})
		if err != nil {
			fmt.Println(currency, " error: ", err)
			continue
		}
		fmt.Println(currency, " Internal SC address: ", internalSCAddr.String())
	}
}

func TestInitEurusERC20(t *testing.T) {
	TestEurusERC20_Init(t)
	TestEurusERC20_AddWriter(t)
	TestEurusERC20_AddEurusUserDepositToAllEurusERC20SCBlacklist(t)
	// TestERC20_SetInternalSCConfigAddress(t)
	//TestMintERC20_Sidechain(t)
}
