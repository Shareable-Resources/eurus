package sc_validator

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ISmartContractInternalSCConfig interface {
	GetInternalSmartContractConfig(opts *bind.CallOpts) (common.Address, error)
}

type ISmartContractInternalSCConfig1 interface {
	InternalSmartContractConfig(opts *bind.CallOpts) (common.Address, error)
}
