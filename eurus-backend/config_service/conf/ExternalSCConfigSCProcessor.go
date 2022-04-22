package conf

import (
	"errors"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func AddCurrencyInfoFromSC(server *ConfigServer, addr string, asset string, decimal int64, id string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	txHash := ""
	if err != nil {
		return txHash, err
	}
	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return txHash, err
	}
	tx, err := instance.AddCurrencyInfo(auth, common.HexToAddress(addr), asset, big.NewInt(decimal), id)
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()

	return txHash, nil
}

func DelCurrencyInfoFromSC(server *ConfigServer, asset string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	txHash := ""
	if err != nil {
		return txHash, err
	}
	auth, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return txHash, err
	}
	tx, err := instance.RemoveCurrencyInfo(auth, asset)
	if err != nil {
		return txHash, err
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

func GetCurrencySCAddrFromSC(server *ConfigServer, asset string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	address, err := instance.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, asset)
	if err != nil {
		return "", err
	}
	if address.Hex() == "0x0000000000000000000000000000000000000000" {
		return "", errors.New("No such asset")
	}
	return address.Hex(), nil
}

func GetCurrencyNameByAddrFromSC(server *ConfigServer, address string) (string, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	asset, err := instance.GetErc20SmartContractByAddr(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return "", err
	}
	if asset == "" {
		return "", errors.New("No such asset")
	}
	return asset, nil
}

func AddAdminFeeToSC(server *ConfigServer, asset string, amount *big.Int) (*types.Transaction, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)

	if err != nil {
		return nil, err
	}
	transOps, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return nil, err
	}

	transOps.GasLimit = server.Config.SideChainGasLimit
	tx, err := instance.SetAdminFee(transOps, asset, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func SetWithdrawFeeFromSC(server *ConfigServer, req *conf_api.SetWithdrawFeeETH) (*types.Transaction, error) {
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)

	if err != nil {
		return nil, err
	}
	transOps, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		return nil, err
	}
	transOps.GasLimit = server.Config.SideChainGasLimit
	formatValue, err := strconv.ParseInt(req.Value, 10, 64)
	if err != nil {
		return nil, err
	}

	assetArray, amountArray := getCurrencyInfoInArrayFromDB(server, big.NewInt(formatValue))

	tx, err := instance.SetETHFee(transOps, big.NewInt(formatValue), assetArray, amountArray)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func GetWithdrawFeeFromSC(server *ConfigServer) (*big.Int, error) {
	withdrawFee := new(big.Int).SetInt64(0)
	instance, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.ServerConfig.ExternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		return withdrawFee, err
	}
	withdrawFee, err = instance.GetAdminFee(&bind.CallOpts{}, "ETH")
	if err != nil {
		return withdrawFee, err
	}
	return withdrawFee, nil
}
