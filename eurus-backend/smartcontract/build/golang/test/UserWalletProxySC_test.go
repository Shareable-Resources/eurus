package test

import (
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestUserWalletProxy_SetInternalSCAddress(t *testing.T) {
	var cenUserAddrList []string = []string{}

	ethClient := initEthClient(t)

	for _, cenAddr := range cenUserAddrList {
		userProxy, err := contract.NewUserWalletProxy(common.HexToAddress(cenAddr), ethClient.Client)
		if err != nil {
			fmt.Println("NewOwnedUpgradeabilityProxy error: ", err, " user addr: ", cenAddr)
			continue
		}

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " user addr: ", cenAddr)
			continue
		}
		internalSC := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

		tx, err := userProxy.SetInternalSCAddress(transOpt, internalSC)
		if err != nil {
			fmt.Println("SetInternalSCAddress error: ", err, " user addr: ", cenAddr)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt failed for user: ", cenAddr, " error: ", err)
			continue
		}
		if receipt.Status == 0 {
			fmt.Println("Query receipt failed for user: ", cenAddr, " Status is 0")
		} else {
			fmt.Println("Set Internal Smart contract config address successfully. user: ", cenAddr)
		}

	}
}

func TestUserWalletProxy_GetInternalSCAddress(t *testing.T) {

	ethClient := initEthClient(t)

	cenAddr := "0xb13070bd35176fae45ff692d5279495547473974"
	userProxy, err := contract.NewUserWalletProxy(common.HexToAddress(cenAddr), ethClient.Client)
	if err != nil {
		fmt.Println("NewOwnedUpgradeabilityProxy error: ", err, " user addr: ", cenAddr)
		t.Fatal(err)
	}

	internalSCAddr, err := userProxy.GetInternalSCAddress(&bind.CallOpts{})
	if err != nil {
		fmt.Println("SetInternalSCAddress error: ", err, " user addr: ", cenAddr)
		t.Fatal(err)
	}
	fmt.Println("internalSCAddr: ", internalSCAddr.Hex())

}
