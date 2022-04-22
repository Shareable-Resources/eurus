package sc_validator

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ISmartContractOwnerable interface {
	GetOwners(opts *bind.CallOpts) ([]common.Address, error)
}
