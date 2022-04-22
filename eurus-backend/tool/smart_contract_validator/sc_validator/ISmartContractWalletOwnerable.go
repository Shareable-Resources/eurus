package sc_validator

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ISmartContractWalletOwnerable interface {
	GetWalletOwner(opts *bind.CallOpts) (common.Address, error)
	GetWalletOperatorList(opts *bind.CallOpts) ([]common.Address, error)
}
