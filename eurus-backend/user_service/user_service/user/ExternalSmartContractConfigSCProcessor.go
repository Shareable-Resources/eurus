package user

import (
	"errors"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (me *UserServer) GetAdminFeeFromSC(assetName string) (*big.Int, error) {
	externalSC, err := contract.NewExternalSmartContractConfig(common.HexToAddress(me.Config.ExternalSCConfigAddress), me.EthClient.Client)
	if err != nil {
		return nil, err
	}

	fee, err := externalSC.GetAdminFee(&bind.CallOpts{}, assetName)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func (me *UserServer) GetDecimalPlaceFromSC(asset string) (int, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(me.ServerConfig.ExternalSCConfigAddress), me.EthClient.Client)
	if err != nil {
		return 0, err
	}
	address, err := instance.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, asset)
	if err != nil {
		return 0, err
	}
	if address.Hex() == "0x0000000000000000000000000000000000000000" {
		return 0, errors.New("No such asset")
	}

	erc20, err := contract.NewERC20(address, me.EthClient.Client)
	if err != nil {
		return 0, err
	}
	decimals, err := erc20.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}
	return int(decimals), err
}
