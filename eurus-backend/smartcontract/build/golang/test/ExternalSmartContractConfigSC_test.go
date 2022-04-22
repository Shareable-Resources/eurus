package test

import (
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestExternalSCConfig_UpgradeProxy(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	implAddr := getAddressBySmartContractName("ExternalSmartContractConfig")
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		t.Fatal("UpgradeTo error : ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestExternalSCConfig_GetAssetAddress(t *testing.T) {
	fmt.Println("Test get address")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	fmt.Println(addr)
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)

	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	assetList, addressList, err := externalSC.GetAssetAddress(&bind.CallOpts{})

	for i, addr1 := range addressList {
		fmt.Println("Addr: ", addr1.String(), " currency: ", assetList[i])
	}
}

func TestExternalSCConfig_GetErc20AddressByAssetName(t *testing.T) {
	fmt.Println("TestGetErc20AddressByAssetName")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	addr1, err := externalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, "USDT")
	if err != nil {
		t.Fatal(" GetErc20SmartContractAddrByAssetName error: ", err)
	}

	fmt.Println("USDT addr: ", addr1.String())
}

func AddErc20Asset(t *testing.T, assetName string, decimal int64, currencyId string) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
		//t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Println(err)
	}
	assetAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
	tx, err := externalSC.RemoveCurrencyInfo(transOpt, assetName)
	if err != nil {
		fmt.Println("RemoveCurrencyInfo error: ", assetName)
		//fatalSmartContractError(t, "AddCurrencyInfo error: ", err)
		t.Fatal(err)
	} else {
		fmt.Println("RemoveCurrencyInfo success: ")
		queryEthReceipt(t, &ethClient, tx)
	}

	transOpt, err = ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("====AddCurrencyInfo for ", assetName, " started====")
	tx, err = externalSC.AddCurrencyInfo(transOpt, assetAddr, assetName, big.NewInt(decimal), currencyId)
	if err != nil {
		fmt.Println(err)

	}
	queryEthReceipt(t, &ethClient, tx)
	fmt.Println("====AddCurrencyInfo for ", assetName, " ended====")
}

func TestExternalSCConfig_AddErc20Asset(t *testing.T) {
	AddErc20Asset(t, "USDM", 6, "tether")
	AddErc20Asset(t, "BTCM", 18, "wrapped-bitcoin")
	AddErc20Asset(t, "ETHM", 18, "")
	// AddErc20Asset(t, "PLA", 6, "tether")
	// AddErc20Asset(t, "ETH", 18, "")
	// AddErc20Asset(t, "USDT", 6, "tether")
	// AddErc20Asset(t, "USDC", 6, "usd-coin")
	// AddErc20Asset(t, "LINK", 18, "chainlink")
	// AddErc20Asset(t, "UNI", 18, "uniswap")
	// AddErc20Asset(t, "BNB", 18, "binancecoin")
	// AddErc20Asset(t, "BUSD", 18, "binance-usd")
	// AddErc20Asset(t, "YFI", 18, "yearn-finance")
	// AddErc20Asset(t, "DAI", 18, "dai")
	// AddErc20Asset(t, "OMG", 18, "omisego")
	// AddErc20Asset(t, "VEN", 18, "impulseven")
	// AddErc20Asset(t, "AAVE", 18, "aave")
	// AddErc20Asset(t, "HT", 18, "huobi-token")
	// AddErc20Asset(t, "SUSHI", 18, "sushi")
	// AddErc20Asset(t, "TUSD", 18, "true-usd")
	// AddErc20Asset(t, "cDAI", 8, "cdai")
	// AddErc20Asset(t, "SXP", 18, "swipe")
	// AddErc20Asset(t, "BAT", 18, "basic-attention-token")
	// AddErc20Asset(t, "USDK", 18, "usdk")
	// AddErc20Asset(t, "WBTC", 8, "wrapped-bitcoin")
	// AddErc20Asset(t, "ZIL", 12, "zilliqa")
	// AddErc20Asset(t, "SNX", 18, "havven")
	// AddErc20Asset(t, "OKB", 18, "okb")
	// AddErc20Asset(t, "BAND", 18, "band-protocol")
	// AddErc20Asset(t, "HUSD", 8, "husd")
	// AddErc20Asset(t, "MKR", 18, "maker")
	// AddErc20Asset(t, "ZRX", 18, "0x")
	// AddErc20Asset(t, "PAX", 18, "paxos-standard")
	// AddErc20Asset(t, "COMP", 18, "compound-governance-token")
	// AddErc20Asset(t, "RSR", 18, "reserve-rights-token")
	// AddErc20Asset(t, "BAL", 18, "balancer")
	AddErc20Asset(t, "MST", 6, "tether")
}

func TestExternalSCConfig_InitKycLevel(t *testing.T) {

	var limit0 *big.Int = big.NewInt(0)
	var limit1 *big.Int = big.NewInt(0)
	// limit0.SetString("5000000000", 10)
	// limit1.SetString("4700000000000", 10)

	// InitExternalSCConfigERC20KycLimit(t, "USDT", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("176000000000000000000", 10)
	// limit1.SetString("164913000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "LINK", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BUSD", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000", 10)
	// limit1.SetString("4700000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "USDC", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("1000000000000000000", 10)
	// limit1.SetString("135000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "YFI", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "DAI", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("556000000000000000000", 10)
	// limit1.SetString("522223000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "OMG", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("12000000000000000000", 10)
	// limit1.SetString("10931000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BNB", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("209000000000000", 10)
	// limit1.SetString("195834000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "UNI", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("12500000000000000000000", 10)
	// limit1.SetString("11750000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "VEN", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("12500000000000000000000", 10)
	// limit1.SetString("11750000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "VEN", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("16000000000000000000", 10)
	// limit1.SetString("14243000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "AAVE", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("345000000000000000000", 10)
	// limit1.SetString("324138000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "HT", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "TUSD", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("250000000000000000000000", 10)
	// limit1.SetString("235000000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "cDAI", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("1667000000000000000000", 10)
	// limit1.SetString("1566667000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "SXP", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("6250000000000000000000", 10)
	// limit1.SetString("5875000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BAT", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "USDK", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("1000000000000000000", 10)
	// limit1.SetString("102000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "WBTC", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("45455000000000000000000", 10)
	// limit1.SetString("42727273000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "ZIL", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("45455000000000000000000", 10)
	// limit1.SetString("42727273000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "ZIL", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("435000000000000000000", 10)
	// limit1.SetString("408696000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "SNX", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("271000000000000000000", 10)
	// limit1.SetString("254055000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "OKB", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("556000000000000000000", 10)
	// limit1.SetString("522223000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BAND", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("2000000000000000000", 10)
	// limit1.SetString("1605000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "MKR", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "HUSD", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "HUSD", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("4630000000000000000000", 10)
	// limit1.SetString("4351852000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "ZRX", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000000000000000", 10)
	// limit1.SetString("4700000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "PAX", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("13000000000000000000", 10)
	// limit1.SetString("11326000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "COMP", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("125000000000000000000000", 10)
	// limit1.SetString("117500000000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "RSR", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("193000000000000000000", 10)
	// limit1.SetString("180770000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BAL", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("2000000000000000000", 10)
	// limit1.SetString("1356000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "ETH", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000", 10)
	// limit1.SetString("4700000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "PLA", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("5000000000", 10)
	// limit1.SetString("4700000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "USDM", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("2000000000000000000", 10)
	// limit1.SetString("1356000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "ETHM", []string{"0", "1"}, []*big.Int{limit0, limit1})

	// limit0.SetString("1000000000000000000", 10)
	// limit1.SetString("102000000000000000000", 10)
	// InitExternalSCConfigERC20KycLimit(t, "BTCM", []string{"0", "1"}, []*big.Int{limit0, limit1})

	limit0.SetString("5000000000000000000000", 10)
	limit1.SetString("4700000000000000000000000", 10)
	InitExternalSCConfigERC20KycLimit(t, "MST", []string{"0", "1"}, []*big.Int{limit0, limit1})
}

func TestGetWithdrawCurrencyDecimal(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	deciaml, err := externalSC.GetAssetDecimal(&bind.CallOpts{}, "ETH")
	if err != nil {
		t.Fatal("GetAssetDecimal error: ", err)
	}
	fmt.Println("ETH decimal :", deciaml)
}

func TestGetOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	owner, _ := externalSC.GetOwners(&bind.CallOpts{})
	fmt.Print(owner)

}

func TestExternalSCConfig_SetWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	for _, configAddr := range configServerHDWalletAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal(err)
		}

		tx, err := externalSC.AddWriter(transOpt, common.HexToAddress(configAddr))
		if err != nil {
			t.Fatal(err)
		}
		queryEthReceipt(t, &ethClient, tx)
	}
}

func TestExternalSCConfig_SetUSDTWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	getUSDTProxy := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	tx, err := externalSC.AddWriter(transOpt, getUSDTProxy)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestExternalSCConfig_GetID(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	tx, err := externalSC.GetAssetListID(&bind.CallOpts{}, "TESTRYAN")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Name : ", tx)

}

func TestExternalSCConfig_SetETHFee(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := externalSC.SetAdminFee(transOpt, "ETHM", big.NewInt(15000000000000000))
	//tx, err := externalSC.SetETHFee(transOpt, big.NewInt(15000000000000000), []string{}, []*big.Int{})
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestExternalSCConfig_GetAdminFee(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	tx, err := externalSC.GetAdminFee(&bind.CallOpts{}, "MST")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amount : ", tx)

}

func TestExternalSCConfig_GetWriterList(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")

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

func TestExternalSCConfig_GetKycLimit(t *testing.T) {
	ethClient := initEthClient(t)
	externalSCConfigInstance, err := contract.NewExternalSmartContractConfig(getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>"), ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}

	rst, err := externalSCConfigInstance.GetCurrencyKycLimit(&bind.CallOpts{}, "MST", "0")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rst)
}

func InitExternalSCConfigERC20KycLimit(t *testing.T, assetName string, kycLevels []string, limits []*big.Int) {
	ethClient := initEthClient(t)
	externalSCConfigInstance, err := contract.NewExternalSmartContractConfig(getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>"), ethClient.Client)
	if err != nil {
		t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	if len(kycLevels) != len(limits) {
		t.Fatal("number of kyc Level and limits not match", err)
	}
	for i, e := range kycLevels {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey: ", err)
		}
		tx, err := externalSCConfigInstance.SetKycLimit(transOpt, assetName, e, limits[i])
		if err != nil {
			fatalSmartContractError(t, "Init Error: ", err)
		}
		queryEthReceipt(t, &ethClient, tx)
	}

}

func TestExternalSCConfig_SetEurusGasPrice(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
		//t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Println(err)
	}

	tx, err := externalSC.SetEurusGasPrice(transOpt, big.NewInt(2400000000))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestExternalSCConfig_GetEurusGasPrice(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
		//t.Fatal("New ExternalSmartContractConfig error: ", err)
	}

	price, err := externalSC.GetEurusGasPrice(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Price: ", price.String())
}

func TestExternalSCConfig_SetMaxTopUpGasAmount(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
		//t.Fatal("New ExternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Println(err)
	}

	tx, err := externalSC.SetMaxTopUpGasAmount(transOpt, big.NewInt(100000000))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestExternalSCConfig_GetMaxTopUpGasAmount(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	externalSC, err := contract.NewExternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
		//t.Fatal("New ExternalSmartContractConfig error: ", err)
	}

	amount, err := externalSC.GetMaxTopUpGasAmount(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Max : ", amount.String())

}

func TestInit_ExternalSC(t *testing.T) {
	TestExternalSCConfig_SetWriter(t)
	TestExternalSCConfig_SetETHFee(t)
	TestExternalSCConfig_AddErc20Asset(t)
	TestExternalSCConfig_InitKycLevel(t)
	TestExternalSCConfig_SetEurusGasPrice(t)
	TestExternalSCConfig_SetMaxTopUpGasAmount(t)
}
