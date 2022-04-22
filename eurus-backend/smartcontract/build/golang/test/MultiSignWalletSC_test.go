package test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ITestMutliSignWalletOperator interface {
	GetWalletOperatorList(opts *bind.CallOpts) ([]common.Address, error)
	GetWalletOwner(opts *bind.CallOpts) (common.Address, error)
}

type ITestMultiOwnable interface {
	GetOwners(opts *bind.CallOpts) ([]common.Address, error)
}

func PrintOwners(t *testing.T, sc ITestMultiOwnable) {
	opts := &bind.CallOpts{
		From: common.HexToAddress(testOwnerAddr),
	}

	list, _ := sc.GetOwners(opts)
	fmt.Println("Owner count: ", len(list))
	for _, addr := range list {
		fmt.Println(addr.String())
	}
}

///Share function
func PrintWalletOperatorList(t *testing.T, sc ITestMutliSignWalletOperator) {
	opts := &bind.CallOpts{
		From: common.HexToAddress(testOwnerAddr),
	}

	list, _ := sc.GetWalletOperatorList(opts)
	fmt.Println("Wallet operator count: ", len(list))
	for _, addr := range list {
		fmt.Println(addr.String())
	}
}

func PrintWalletOwner(t *testing.T, sc ITestMutliSignWalletOperator) {
	opts := &bind.CallOpts{
		From: common.HexToAddress(testOwnerAddr),
	}

	addr, _ := sc.GetWalletOwner(opts)
	fmt.Println("Wallet owner: ", addr.String())
}
