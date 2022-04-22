package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestAdminFeeWallet_SetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err := AdminFeeWalletSC.SetWalletOwner(transOpt, common.HexToAddress(testOwnerAddr))
	if err != nil {
		t.Fatal("SetWalletOwner error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestAdminFeeWallet_GetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	proxySC, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}

	implAddr, err := proxySC.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}

	fmt.Println("Addr: ", implAddr.Hex())

}

func TestAdminFeeWallet_TransferEUN(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	proxySC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	amount := big.NewInt(0)
	amount = amount.SetUint64(100)
	transOpt.GasLimit = 1000000
	tx, err := proxySC.TransferETH(transOpt, common.HexToAddress(testOwnerAddr), amount)
	if err != nil {
		t.Fatal("Transfer error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestAdminFeeWallet_SetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	proxySC, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	actualAddr := getAddressBySmartContractName("AdminFeeWallet")
	tx, err := proxySC.UpgradeTo(transOpt, actualAddr)
	if err != nil {
		t.Fatal("UpgradeTo error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestAdminFeeWallet_AddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	//RYAN DEBUG
	addressList := []string{
		"0x7EBcc140Ceac70e93fFe5B6a205A1a693AB51109",
		"0x81ccE1BF6DCfCeF814914E66d1aD6ACcFeB3936e",
		"0x80B5aD3d6370f9C4b6b84a222f42CD3FC5755D81",
		"0x96ad0B6DdD7CEb4e231dA09F7AE5044b98De5362",
		"0xA14Ab677452e5c038067f6dF57eE073328FbD336",
		"0xEECd77226192B0Cf9CE510A254F31f03c3CD2D85",
		"0x6612457C18AB8A7B1caF25bd316401359132B5cd",
		"0xdF86024284b12A1a274183A920cFf1419D524F17",
		"0x206803F71C5eB93C61a060eeCfc94C2CE6D9747D",
		"0x15C9894F2C32A892A9965Fa23b6026Abd7450Ac8",
		"0x7076cE3E9C32c6D9F625E142556E0944BBDAe6a0",
		"0xffa498003B1b1AfC2804533fF7Fb11a5ACC059C2",
		"0xa4fC063Be769a24886d0fe219d9F1A24b09dAbC2",
		"0x38455a54cC578505b6a57Aa0957d1BA9B400c9CA",
		"0xb393Edd0Bf261C49C8A5C3fb21B23DE8753BE967",
	}

	for _, address := range addressList {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
		}

		tx, err := AdminFeeWalletSC.AddWalletOperator(transOpt, common.HexToAddress(address))
		if err != nil {
			// t.Fatal("AddWalletOperator error: ", err)
			fmt.Println("address: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("receipt error: ", err)
			continue
		}

		receiptByte, _ := receipt.MarshalJSON()
		fmt.Println(string(receiptByte))

	}

	// queryEthReceipt(t, &ethClient, tx)

}

func TestAdminFeeWallet_GetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	ownerAddr, err := AdminFeeWalletSC.GetWalletOwner(&bind.CallOpts{})
	if err != nil {
		t.Fatal("AddWalletOperator error: ", err)
	}

	fmt.Println("Wallet owner : ", ownerAddr.Hex())

}

func TestAdminFeeWallet_GetWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	operatorList, err := AdminFeeWalletSC.GetWalletOperatorList(&bind.CallOpts{})
	if err != nil {
		t.Fatal("AddWalletOperator error: ", err)
	}

	for _, operatorAddr := range operatorList {
		fmt.Println("Wallet owner : ", operatorAddr.Hex())
	}

}

func TestAdminFeeWallet_ChangeRequirement(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")

	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	req := big.NewInt(3)
	tx, err := AdminFeeWalletSC.ChangeRequirement(transOpt, req)
	if err != nil {
		t.Fatal("ChangeRequirement error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

//func TestAdminFeeWallet_GetConfirmation(t *testing.T) {
//	ethClient := initEthClient(t)
//	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")
//
//	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
//	if err != nil {
//		t.Fatal("NewPlatform error: ", err)
//	}
//
//	addressList := []string{
//		"0x81ccE1BF6DCfCeF814914E66d1aD6ACcFeB3936e",
//		"0x80B5aD3d6370f9C4b6b84a222f42CD3FC5755D81",
//		"0x96ad0B6DdD7CEb4e231dA09F7AE5044b98De5362",
//		"0xA14Ab677452e5c038067f6dF57eE073328FbD336",
//		"0xEECd77226192B0Cf9CE510A254F31f03c3CD2D85",
//		"0x6612457C18AB8A7B1caF25bd316401359132B5cd",
//		"0xdF86024284b12A1a274183A920cFf1419D524F17",
//		"0x206803F71C5eB93C61a060eeCfc94C2CE6D9747D",
//		"0x15C9894F2C32A892A9965Fa23b6026Abd7450Ac8",
//		"0x7076cE3E9C32c6D9F625E142556E0944BBDAe6a0",
//		"0xffa498003B1b1AfC2804533fF7Fb11a5ACC059C2",
//		"0xa4fC063Be769a24886d0fe219d9F1A24b09dAbC2",
//		"0x38455a54cC578505b6a57Aa0957d1BA9B400c9CA",
//		"0xb393Edd0Bf261C49C8A5C3fb21B23DE8753BE967",
//	}
//
//	transId := big.NewInt(15)
//	for _, address := range addressList {
//		isConfirm, err := AdminFeeWalletSC.Confirmations(&bind.CallOpts{}, transId, common.HexToAddress(address))
//		if err != nil {
//			fmt.Printf("%s error: %s\r\n", err.Error(), address)
//		} else {
//			fmt.Printf("%s: %v\r\n", address, isConfirm)
//		}
//	}
//}

func TestAdminFeeWallet_SetInternalSmartContractAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")
	AdminFeeWalletSC, err := contract.NewAdminFeeWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	internalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

	fmt.Println("InternalSmartContractConfig address: ", internalSCAddr.Hex())
	tx, err := AdminFeeWalletSC.SetInternalSmartContractConfig(transOpt, internalSCAddr)
	if err != nil {
		t.Fatal("ChangeRequirement error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestAdminFeeWallet_GetETHBalance(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<AdminFeeWallet>")
	amount, err := ethClient.GetBalance(addr)
	if err != nil {
		t.Fatal("GetBalance: ", err)
	}
	fmt.Println("ETH Balance: ", amount.String())
}

func TestInitAdminFeeWallet(t *testing.T) {
	TestAdminFeeWallet_SetWalletOwner(t)
	//Not required
	// TestAdminFeeWallet_AddWalletOperator(t)
	// TestAdminFeeWallet_ChangeRequirement(t)
	TestAdminFeeWallet_SetInternalSmartContractAddress(t)
}
