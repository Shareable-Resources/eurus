package reward

import (
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type RewardProcessor struct {
	InternalSmartContractConfigAddr common.Address
	SideChainEthClient              *ethereum.EthClient
	DbProcessor                     *RewardDBProcessor
	MarketingInvokerPrivateKey      string
	Logger                          *logrus.Logger
	SideChainTransferGasLimit       uint64
	InternalSCConfig                *contract.InternalSmartContractConfig
}

func NewRewardProcessor(db *database.Database, slaveDb *database.ReadOnlyDatabase, marketInvokerPrivateKey string, logger *logrus.Logger) *RewardProcessor {
	processor := new(RewardProcessor)
	processor.DbProcessor = NewRewardDBProcessor(db, slaveDb)
	processor.Logger = logger
	processor.MarketingInvokerPrivateKey = marketInvokerPrivateKey
	return processor
}

func (me *RewardProcessor) Init(internalSCAddr common.Address, sideChainEthClient *ethereum.EthClient) error {
	var err error

	me.InternalSmartContractConfigAddr = internalSCAddr
	me.SideChainEthClient = sideChainEthClient

	me.InternalSCConfig, err = contract.NewInternalSmartContractConfig(me.InternalSmartContractConfigAddr, me.SideChainEthClient.Client)
	if err != nil {
		return err
	}

	return nil
}

func (me *RewardProcessor) TransferEUN(destWalletAddr string, srcWalletAddr common.Address, amount *big.Int) (*types.Transaction, *ethereum.BesuReceipt, error) {

	marketingRegWallet, err := contract.NewMarketingWallet(srcWalletAddr, me.SideChainEthClient.Client)
	if err != nil {
		me.Logger.Errorln("NewMarketingWallet error: ", err.Error(), " wallet address: ", destWalletAddr)
		return nil, nil, err
	}
	transOpt, err := me.SideChainEthClient.GetNewTransactorFromPrivateKey(me.MarketingInvokerPrivateKey, me.SideChainEthClient.ChainID)
	if err != nil {
		me.Logger.Errorln("GetNewTransactorFromPrivateKey error: ", err.Error(), " wallet address: ", destWalletAddr)
		return nil, nil, err
	}

	if me.SideChainTransferGasLimit > 0 {
		transOpt.GasLimit = me.SideChainTransferGasLimit
	} else {
		transOpt.GasLimit = 1000000
	}

	tx, err := marketingRegWallet.TransferETH(transOpt, common.HexToAddress(destWalletAddr), amount)
	if err != nil {
		me.Logger.Errorln("TransferETH error: ", err.Error(), " wallet address: ", destWalletAddr)
		return nil, nil, err
	}

	me.Logger.Infoln("transfer EUN to user trans hash ", tx.Hash().Hex(), " for wallet address: ", destWalletAddr)
	receipt, err := me.SideChainEthClient.QueryEthReceipt(tx)
	if err != nil {
		me.Logger.Errorln("Unable to transfer EUN to user: ", err.Error(), " wallet address: ", destWalletAddr)
		return tx, nil, err
	}
	return tx, receipt, nil
}
