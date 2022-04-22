package userObserver

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type AssetAddressMap map[common.Address]string

type UserObserverSCProcessor struct {
	UserWalletAbi           *abi.ABI
	sideChainEthClient      *ethereum.EthClient
	SideChainUserWalletAddr []common.Address
	context                 *UserObserverContext
	config                  *UserObserverConfig
}

func NewUserObserverSCProcessor(config *UserObserverConfig, loggerName string) *UserObserverSCProcessor {
	processor := new(UserObserverSCProcessor)
	processor.config = config
	processor.UserWalletAbi = ethereum.DefaultABIDecoder.GetABI("UserWallet")
	processor.context = NewUserObserverContext(nil, nil, config, loggerName)
	return processor
}

func (me *UserObserverSCProcessor) Init() error {
	var err error
	//Connect sidechain
	me.sideChainEthClient = &ethereum.EthClient{
		Protocol: me.config.EthClientProtocol,
		IP:       me.config.EthClientIP,
		Port:     me.config.EthClientPort,
		ChainID:  big.NewInt(int64(me.config.EthClientChainID)),
	}
	_, err = me.sideChainEthClient.Connect()
	if err != nil {
		return errors.WithMessage(err, "Connect sidechain failed")
	}

	//me.SideChainUserWalletAddr,err = GetUserWalletAddrFormDB(me)
	//if err != nil {
	//	return errors.WithMessage(err, "Can not get UserWallet address")
	//}

	return nil
}

func (me *UserObserverSCProcessor) GetTransIdFromTransferRequestEvent(transLog *types.Log, tx *types.Transaction) (*big.Int, error) {

	if transLog.Topics[0] != ethereum.DefaultABIDecoder.GetABI("UserWallet").Events["TransferRequestEvent"].ID {
		return nil, errors.New("Topic is not TransferRequestEvent")
	}
	if len(transLog.Topics) < 2 {
		return nil, errors.New("Trans ID not found")
	}
	transId := transLog.Topics[1].Big()

	return transId, nil
}

func (me *UserObserverSCProcessor) ConfirmTransactionWithTransID(transID *big.Int, transLog *types.Log, proxyAddr common.Address) (*types.Transaction, error) {
	transOpt, err := me.sideChainEthClient.GetNewTransactorFromPrivateKey(me.config.HdWalletPrivateKey, me.sideChainEthClient.ChainID)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("WithdrawEvent - Unable to create smart contract instance: ", err.Error())
		return nil, err
	}

	if me.config.SideChainGasLimit > 0 {
		transOpt.GasLimit = me.config.SideChainGasLimit
	}

	userWallet, err := contract.NewUserWallet(proxyAddr, me.sideChainEthClient.Client)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Can not load user wallet instance", err)
		return nil, err
	}
	tx, err := userWallet.ConfirmTransaction(transOpt, transID)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Fail to confirm the transaction", err)
		return nil, err
	}
	return tx, nil

}
