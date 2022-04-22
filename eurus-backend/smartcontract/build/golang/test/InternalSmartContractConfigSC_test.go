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

func TestInternalSCConfig_UpgradeProxy(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
	}

	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	implAddr := getAddressBySmartContractName("InternalSmartContractConfig")
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		fmt.Println("UpgradeTo error : ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetWalletAddressMap(t *testing.T) {
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	ethClient := initEthClient(t)
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.SetWalletAddressMap(transOpt, common.HexToAddress("0xA453e5615708e8D9149E3AFf75842628ee75d1Fb"))
	if err != nil {
		t.Fatal("SetWalletAddressMap error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestInternalSCConfig_SetApprovalWallet(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	addr1 := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	tx, err := internalSC.SetApprovalWalletAddress(transOpt, addr1)
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetUserWalletProxy(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	userWalletProxy := getAddressBySmartContractName("UserWalletProxy")
	tx, err := internalSC.SetUserWalletProxyAddress(transOpt, userWalletProxy)
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestInternalSCConfig_SetMarketRegWallet(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	marketRegWallet := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.SetMarketingRegWalletAddress(transOpt, marketRegWallet)
	if err != nil {
		t.Fatal("SetMarketingRegWalletAddress", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetUserWalletAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	userWalletAddr := getAddressBySmartContractName("UserWallet")
	tx, err := internalSC.SetUserWalletAddress(transOpt, userWalletAddr)
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetUserWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestInternalSCConfig_SetInnetWalletAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	platformWalletAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	tx, err := internalSC.SetInnetWalletAddress(transOpt, platformWalletAddr)
	if err != nil {
		t.Fatal("SetInnetWalletAddress error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetUserWallet(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	userWalletAddr := getAddressBySmartContractName("UserWallet")

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := internalSC.SetUserWalletAddress(transOpt, userWalletAddr)
	queryEthReceipt(t, &ethClient, tx)

}

func TestInternalSCConfig_GetConfig(t *testing.T) {
	fmt.Println("TestGetConfig")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)

	fmt.Println(addr)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	approvalWallet, err := internalSC.GetApprovalWalletAddress(&bind.CallOpts{})
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("GetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}

	fmt.Println("Approval wallet address: ", approvalWallet.String())

	scAddr, _ := internalSC.GetWithdrawSmartContract(&bind.CallOpts{})
	fmt.Println("Withdraw Smart contract address: ", scAddr.String())

	platformWalletAddr, _ := internalSC.GetInnetPlatformWalletAddress(&bind.CallOpts{})
	fmt.Println("PlatformWallet address: ", platformWalletAddr.Hex())

	walletMapAddr, _ := internalSC.GetWalletAddressMap(&bind.CallOpts{})
	fmt.Println("Wallet address map address: ", walletMapAddr.Hex())

	extAddr, _ := internalSC.GetExternalSCConfigAddress(&bind.CallOpts{})
	fmt.Println("GetExternalSCConfigAddress address: ", extAddr.Hex())

	adminfeeWalletAddr, _ := internalSC.GetAdminFeeWalletAddress(&bind.CallOpts{})
	fmt.Println("GetAdminFeeWallet address: ", adminfeeWalletAddr.Hex())

	usdtAddr, _ := internalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, "USDT")
	fmt.Println("USDT address: ", usdtAddr.Hex())

	userAddr, _ := internalSC.GetUserWalletAddress(&bind.CallOpts{})
	fmt.Println("UserWallet address: ", userAddr)

	gasFeeWalletAddr, _ := internalSC.GetGasFeeWalletAddress(&bind.CallOpts{})
	fmt.Println("Gas fee wallet address: ", gasFeeWalletAddr.Hex())

	marketingRegAddr, _ := internalSC.GetMarketingRegWalletAddress(&bind.CallOpts{})
	fmt.Println("Marketing reg wallet address: ", marketingRegAddr.Hex())

	userProxy, _ := internalSC.GetUserWalletProxyAddress(&bind.CallOpts{})
	fmt.Println("User wallet proxy address: ", userProxy.Hex())

	fee, err := internalSC.GetCentralizedGasFeeAdjustment(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("GetCentralizedGasFeeAdjustment: ", fee.String())

}

func TestInternalSCConfig_SetWalletAddressMapAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	walletMapAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")

	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	//transOpt.GasPrice = big.NewInt(100000000)
	tx, err := internalSC.SetWalletAddressMap(transOpt, walletMapAddr)
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetWalletAddressMap error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)
}
func TestInternalSCConfig_SetExternalSmartContractConfigAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.SetExternalSCConfigAddress(transOpt, getAddressBySmartContractName("OwnedUpgradeabilityProxy<ExternalSmartContractConfig>"))
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetWithdrawSmartContractConfigAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.SetWithdrawSmartContract(transOpt, getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>"))
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_SetAdminFeeWalletAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.SetAdminFeeWalletAddress(transOpt, getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>"))
	if err != nil {

		jsonData, _ := json.Marshal(err)

		t.Fatal("SetApprovalWalletAddress error: ", err, "\r\nJSON: ", string(jsonData))
	}
	queryEthReceipt(t, &ethClient, tx)
}

//func TestSetWithdrawFeeMap(t *testing.T) {
//	ethClient := initEthClient(t)
//	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
//	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
//	if err != nil {
//		t.Fatal("NewInternalSmartContractConfig error: ", err)
//	}
//
//	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
//	if err != nil {
//		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
//	}
//	tx, err := internalSC.SetwithdrawFeeETH(transOpt,"ETH", big.NewInt(10000000000000000))
//	if err != nil {
//		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
//	}
//	queryEthReceipt(t, &ethClient, tx)
//}

func TestAddOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.AddOwner(transOpt, common.HexToAddress(testOwnerAddr))
	if err != nil {
		fmt.Print(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}
func TestInternalSCConfig_AddWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := internalSC.AddWriter(transOpt, common.HexToAddress("0x8c9e314f4bd8dacde1ca4e9fd63064f7cd3388fa"))
	if err != nil {
		fmt.Print(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

//func TestGetWithdrawFeeMap(t *testing.T) {
//	ethClient := initEthClient(t)
//	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
//	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
//	if err != nil {
//		t.Fatal("NewInternalSmartContractConfig error: ", err)
//	}
//
//
//	tx, err := internalSC.WithdrawFeeMap(&bind.CallOpts{},"ETH")
//	if err != nil {
//		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
//	}
//	fmt.Println(tx)
//}

func TestInternalSCConfig_SetGasFeeWalletAddress(t *testing.T) {

	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err := internalSC.SetGasFeeWalletAddress(transOpt, common.HexToAddress(gasFeeCollectWalletAddr))
	if err != nil {
		t.Fatal("SetGasFeeWalletAddress error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestInternalSCConfig_SetCenteralizedGasFeeAdjustment(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	gasFee := big.NewInt(49000)
	tx, err := internalSC.SetCentralizedGasFeeAdjustment(transOpt, gasFee)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestInternalSCConfig_GetAssetByAssetName(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewInternalSmartContractConfig error: ", err)
	}

	assetAddr, err := internalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, "USDM")
	if err != nil {
		t.Fatal("GetErc20SmartContractAddrByAssetName: ", err)
	}

	fmt.Println(assetAddr)
}

func TestInitInternalSC(t *testing.T) {
	//TestInternalSCConfig_UpgradeProxy(t)
	TestInternalSCConfig_SetApprovalWallet(t)
	TestInternalSCConfig_SetInnetWalletAddress(t)
	TestInternalSCConfig_SetWalletAddressMapAddress(t)
	TestInternalSCConfig_SetAdminFeeWalletAddress(t)
	TestInternalSCConfig_SetExternalSmartContractConfigAddress(t)
	TestInternalSCConfig_SetWithdrawSmartContractConfigAddress(t)
	TestInternalSCConfig_SetGasFeeWalletAddress(t)
	TestInternalSCConfig_SetUserWalletProxy(t)
	TestInternalSCConfig_SetUserWallet(t)
	TestInternalSCConfig_GetConfig(t)
}
