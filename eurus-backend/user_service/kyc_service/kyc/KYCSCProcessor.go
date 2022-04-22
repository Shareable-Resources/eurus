package kyc

import (
	"errors"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type KYCSCProcessor struct {
	EthClient               *ethereum.EthClient
	WalletAddressMapAddress common.Address
	Logger                  *logrus.Logger
	Config                  *KYCConfig
}

func (me *KYCSCProcessor) Init(ethClient *ethereum.EthClient, config *KYCConfig, logger *logrus.Logger) error {
	var err error
	me.EthClient = ethClient
	me.Config = config
	me.Logger = logger
	internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(config.InternalSCConfigAddress), me.EthClient.Client)
	if err != nil {
		me.Logger.Error("Unable to new internal smart contract config. Error: ", err)
		return err
	}

	for i := 0; i < config.RetryCount; i++ {
		me.WalletAddressMapAddress, err = internalSCConfig.GetWalletAddressMap(&bind.CallOpts{})
		if err == nil {
			break
		}
	}
	if err != nil {
		me.Logger.Error("Get wallet address map error: ", err)
		return err
	}

	return nil
}

func (me *KYCSCProcessor) SetUserKYCLevel(walletAddress common.Address, level int) error {
	walletAddressMap, err := contract.NewWalletAddressMap(me.WalletAddressMapAddress, me.EthClient.Client)
	if err != nil {
		me.Logger.Error("Unable to create wallet address map instance. Wallet address: ", walletAddress.Hex(), " Error: ", err)
		return err
	}

	tx, err := me.EthClient.InvokeSmartContract(me.Config, me.Config.HdWalletPrivateKey, 0,
		func(ethClient *ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
			tx, err := walletAddressMap.SetWalletInfo(transOpt, walletAddress, "kycLevel", strconv.Itoa(level))
			return tx, false, err
		})

	if err != nil {
		me.Logger.Error("Unable to set wallet info. Wallet address: ", walletAddress.Hex(), " Error: ", err)
		return err
	}
	receipt, err := me.EthClient.QueryEthReceipt(tx)
	if err != nil {
		me.Logger.Error("Unable to query receipt. Wallet address: ", walletAddress.Hex(), " Error: ", err)
		return err
	}
	if receipt.Status == 0 {
		me.Logger.Error("Set wallet info failed to Wallet address: ", walletAddress.Hex())
		return errors.New("Set wallet info failed to wallet address: " + walletAddress.Hex())
	}

	return nil
}
