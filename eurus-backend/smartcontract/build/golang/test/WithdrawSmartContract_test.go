package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestWithrawSmartContract_GetWalletOperatorList(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	PrintWalletOperatorList(t, withdrawSC)
}

func TestWithrawSmartContract_AddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	for _, addr := range withdrawobserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("Get test owner private key: ", err)
			continue
		}
		tx, err := withdrawSC.AddWalletOperator(transOpt, common.HexToAddress(addr))
		if err != nil {
			fmt.Println("AddWalletOperator error: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " withdraw observer address: ", addr)
		} else {
			jsonByte, err := receipt.MarshalJSON()
			if err != nil {
				fmt.Println("Receipt marshel error: ", err)
			} else {
				fmt.Println(string(jsonByte))

			}
		}
	}

	for _, addr := range approvalObserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("Get test owner private key: ", err)
			continue
		}
		tx, err := withdrawSC.AddWalletOperator(transOpt, common.HexToAddress(addr))
		if err != nil {
			fmt.Println("AddWalletOperator error: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " withdraw observer address: ", addr)
		} else {
			jsonByte, err := receipt.MarshalJSON()
			if err != nil {
				fmt.Println("Receipt marshel error: ", err)
			} else {
				fmt.Println(string(jsonByte))

			}
		}
	}

}

func TestWithrawSmartContract_GetOwners(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	PrintOwners(t, withdrawSC)
}
func TestWithrawSmartContract_RemoveOwners(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testWithdrawObserverPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	withdrawSC.RemoveOwner(transOpt, common.HexToAddress("0x69ee00BF2a3dFF7DD9f6eE7D04f44E87FB8b917D"))
}
func TestWithrawSmartContract_SetOwners(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testWithdrawObserverPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := withdrawSC.SetWalletOwner(transOpt, common.HexToAddress(testWithdrawObserverAddr))
	if err != nil {
		t.Fatal(err)

		queryEthReceipt(t, &ethClient, tx)
	}
}

func TestWithrawSmartContract_GetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	PrintWalletOwner(t, withdrawSC)
}
func TestWithrawSmartContract_GetWalletWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	list, err := withdrawSC.GetWriterList(&bind.CallOpts{})

	for _, addr := range list {
		fmt.Println(addr.Hex())
	}
}
func TestWithrawSmartContract_SetWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	approvalWalletAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	tx, err := withdrawSC.AddWriter(transOpt, approvalWalletAddr)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithrawSmartContract_SetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	for _, addr := range withdrawobserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerAddr, ethClient.ChainID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(transOpt, addr, withdrawSC)
		tx, err := withdrawSC.SetWalletOwner(transOpt, common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(tx.Hash())
	}

}

func TestWithrawSmartContract_SetInternalSCAddr(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		fmt.Println(err)

	}

	fmt.Println(transOpt, addr, withdrawSC)
	intc := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	count := big.NewInt(5)
	tx, err := withdrawSC.Init(transOpt, intc, count)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(tx.Hash())
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithrawSmartContract_ConfirmTransaction(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testWithdrawObserverPrivateKey, ethClient.ChainID)
	tx, err := withdrawSC.ConfirmTransaction(transOpt, big.NewInt(8))
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithrawSmartContract_ChangeNumOfRequirement(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	tx, err := withdrawSC.ChangeRequirement(transOpt, big.NewInt(5))
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithdrawSC_Init(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	internalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	required := big.NewInt(5)
	tx, err := withdrawSC.Init(transOpt, internalSCAddr, required)
	if err != nil {
		t.Fatal("Init error:", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithdrawSC_Upgrade(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	implAddr := getAddressBySmartContractName("WithdrawSmartContract")
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		t.Fatal("UpgradeTo error : ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestWithdrawSC_GetConfirmations(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	for _, withdrawObsAddr := range withdrawobserverAddr {

		isConfirmed, err := withdrawSC.Confirmations(&bind.CallOpts{}, big.NewInt(13), common.HexToAddress(withdrawObsAddr))
		if err != nil {
			fmt.Printf("Error on %s: %s\r\n", withdrawObsAddr, err.Error())
			continue
		}
		fmt.Printf("%s: %v\r\n", withdrawObsAddr, isConfirmed)
	}

}

func TestWithdrawSmartContract_GetInternalSmartContractAddr(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	scAddr, err := withdrawSC.InternalSCConfig(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("InternalSCAddr: ", scAddr.Hex())

}

func TestWithdrawSC_GetMisData(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	wall, err := withdrawSC.MiscellaneousData(&bind.CallOpts{}, big.NewInt(8), "approvalWallet")

	if err != nil {
		t.Fatal(err)
	}

	walletAddr := common.BytesToAddress([]byte(wall))
	fmt.Println(walletAddr.Hex())

}

func TestWithdrawSC_GetTransactionValue(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	withdrawSC, err := contract.NewWithdrawSmartContract(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	txValue, err := withdrawSC.Transactions(&bind.CallOpts{}, big.NewInt(8))

	if err != nil {
		t.Fatal(err)
	}

	val := txValue.Value
	fmt.Println(val.String())

}

func TestInitWithdrawSC(t *testing.T) {
	TestWithrawSmartContract_SetWalletOwner(t)
	TestWithrawSmartContract_AddWalletOperator(t) //This function needs withdraw observers address, so may need to run afterward
	TestWithrawSmartContract_SetWriter(t)
	TestWithrawSmartContract_ChangeNumOfRequirement(t)
	TestWithrawSmartContract_SetInternalSCAddr(t)
}
