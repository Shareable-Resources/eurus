package test

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestEurusInternalConfig_GetImplementation(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	impl, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Implementation", err)
	}
	fmt.Println("Implementation: ", impl.String())
	owner, err := proxy.ProxyOwner(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Implementation", err)
	}
	fmt.Println("Owner: ", owner.String())
}

func AddCurrencyInfo(t *testing.T, ethClient ethereum.EthClient, erc20AssetName string) {
	fmt.Println("Currency: ", erc20AssetName)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	erc20Addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + erc20AssetName + ">")
	tx, err := eurusInternalConfig.AddCurrencyInfo(transOpt, erc20Addr, erc20AssetName)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TX Hash: ", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusInternalConfig_AddCustomCurrencyInfo(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")

	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	tx, err := eurusInternalConfig.AddCurrencyInfo(transOpt, common.HexToAddress("0x256e8aB0e9d2a2d43cd4e9ba63d93D3d3Ed6DAE2"), "USDT")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("TX Hash: %s", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusInternalConfig_AddCurrencyInfo(t *testing.T) {
	ethClient := initTestNetEthClient(t)

	// AddCurrencyInfo(t, ethClient, "USDM")
	// AddCurrencyInfo(t, ethClient, "BTCM")
	// AddCurrencyInfo(t, ethClient, "ETHM")
	// AddCurrencyInfo(t, ethClient, "USDT")
	// AddCurrencyInfo(t, ethClient, "USDC")
	// AddCurrencyInfo(t, ethClient, "LINK")
	// AddCurrencyInfo(t, ethClient, "UNI")
	// AddCurrencyInfo(t, ethClient, "BNB")
	// AddCurrencyInfo(t, ethClient, "BUSD")
	// AddCurrencyInfo(t, ethClient, "YFI")
	// AddCurrencyInfo(t, ethClient, "DAI")
	// AddCurrencyInfo(t, ethClient, "OMG")
	// AddCurrencyInfo(t, ethClient, "VEN")
	// AddCurrencyInfo(t, ethClient, "AAVE")
	// AddCurrencyInfo(t, ethClient, "HT")
	// AddCurrencyInfo(t, ethClient, "SUSHI")
	// AddCurrencyInfo(t, ethClient, "TUSD")
	// AddCurrencyInfo(t, ethClient, "cDAI")
	// AddCurrencyInfo(t, ethClient, "SXP")
	// AddCurrencyInfo(t, ethClient, "BAT")
	// AddCurrencyInfo(t, ethClient, "USDK")
	// AddCurrencyInfo(t, ethClient, "WBTC")
	// AddCurrencyInfo(t, ethClient, "ZIL")
	// AddCurrencyInfo(t, ethClient, "SNX")
	// AddCurrencyInfo(t, ethClient, "OKB")
	// AddCurrencyInfo(t, ethClient, "BAND")
	// AddCurrencyInfo(t, ethClient, "MKR")
	// AddCurrencyInfo(t, ethClient, "HUSD")
	// AddCurrencyInfo(t, ethClient, "ZRX")
	// AddCurrencyInfo(t, ethClient, "PAX")
	// AddCurrencyInfo(t, ethClient, "COMP")
	// AddCurrencyInfo(t, ethClient, "RSR")
	// AddCurrencyInfo(t, ethClient, "BAL")
	// AddCurrencyInfo(t, ethClient, "PLA")
	AddCurrencyInfo(t, ethClient, "MST")
}

func TestEurusInternalConfig_RemoveCurrencyInfo(t *testing.T) {
	EurusInternalConfig_RemoveCurrencyInfo(t, "MST")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "USDT")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "USDC")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "LINK")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "UNI")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "BNB")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "BUSD")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "YFI")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "DAI")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "OMG")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "VEN")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "AAVE")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "HT")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "SUSHI")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "TUSD")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "cDAI")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "SXP")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "BAT")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "USDK")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "WBTC")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "ZIL")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "SNX")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "OKB")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "BAND")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "MKR")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "HUSD")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "ZRX")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "PAX")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "COMP")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "RSR")
	// EurusInternalConfig_RemoveCurrencyInfo(t, "BAL")

}

func EurusInternalConfig_RemoveCurrencyInfo(t *testing.T, assetName string) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 500000
	tx, err := eurusInternalConfig.RemoveCurrencyInfo(transOpt, assetName)
	if err != nil {
		fatalSmartContractError(t, "Init Error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusInternalConfig_GetCurrencyInfo(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	addr, err = eurusInternalConfig.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, "USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Address Name: %s\r\n", addr.Hex())
}

func TestEurusInternalConfig_SetEurusUserDepositAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	depositAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 500000
	tx, err := eurusInternalConfig.SetEurusUserDepositAddress(transOpt, depositAddr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Tx hash: ", tx.Hash().Hex())

	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusInternalConfig_GetAssetAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	assetNameList, addressList, err := eurusInternalConfig.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		t.Fatal("GetAssetAddress: ", err)
	}

	for i, name := range assetNameList {
		fmt.Printf("\r\n%s: %s", name, addressList[i].Hex())
	}

}

func TestEurusInternalConfig_SetEurusPlatformWalletAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	platformWalletAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := eurusInternalConfig.SetPlatformWalletAddress(transOpt, platformWalletAddr)
	if err != nil {
		t.Fatal("SetPlatformWalletAddress error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusInternalConfig_GetEurusUserDepositAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	depositAddr, err := eurusInternalConfig.EurusUserDepositAddress(&bind.CallOpts{})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("User deposit address: ", depositAddr.Hex())
}

func TestEurusInternalConfig_GetEurusPlatformWalletAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	walletAddr, err := eurusInternalConfig.PlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Platform wallet address: ", walletAddr.Hex())
}

func TestEurusInternalConfig_GetOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	ownerList, err := eurusInternalConfig.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, owner := range ownerList {
		fmt.Println("Owner: ", owner)
	}
}

func TestEurusInternalConfig_AddOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	ownerAddr := getTestNetAddressBySmartContractName("GeneralMultiSigWallet")
	tx, err := eurusInternalConfig.AddOwner(transOpt, ownerAddr)
	if err != nil {
		t.Fatal("SetPlatformWalletAddress error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusInternalConfig_BatchUpdateCurrencyInfo(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>")
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	transOpt.GasLimit = 1000000

	var assetNameList []string = []string{"USDM", "BTCM", "ETHM", "MST"}
	var assetAddressList []common.Address = make([]common.Address, 0)

	for _, assetName := range assetNameList {
		assetAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
		log.GetLogger(log.Name.Root).Debugln("OwnedUpgradeabilityProxy<"+assetName+">: ", assetAddr.Hex())
		assetAddressList = append(assetAddressList, assetAddr)
	}

	tx, err := eurusInternalConfig.BatchUpdateAssetAddress(transOpt, assetNameList, assetAddressList)
	if err != nil {
		t.Fatal("BatchUpdateAssetAddress failed: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusInternalSC_Init(t *testing.T) {
	TestEurusInternalConfig_AddCurrencyInfo(t)
	TestEurusInternalConfig_SetEurusPlatformWalletAddress(t)
}
