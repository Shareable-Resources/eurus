package withdrawal

import (
	"encoding/json"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type WithdrawObserverSCProcessor struct {
	WithdrawContractAbi *abi.ABI

	Config                       *WithdrawObserverConfig
	sidechainEthClient           *eurus_ethereum.EthClient
	mainnetEthClient             *eurus_ethereum.EthClient
	MainnetPlatformWalletAddress *common.Address
	loggerName                   string
}

type WithdrawEvent struct {
	ApprovalWallet *common.Address
	RequestTransId *big.Int
	SrcWallet      *common.Address
	DestWallet     *common.Address
	AssetName      string
	BurnTransId    *big.Int
	Amount         *big.Int
}

func NewWithdrawObserverSCProcessor(config *WithdrawObserverConfig, loggerName string) *WithdrawObserverSCProcessor {
	processor := new(WithdrawObserverSCProcessor)
	processor.WithdrawContractAbi = new(abi.ABI)
	*processor.WithdrawContractAbi, _ = abi.JSON(strings.NewReader(string(contract.WithdrawSmartContractABI)))
	processor.Config = config
	processor.loggerName = loggerName
	return processor
}

func (me *WithdrawObserverSCProcessor) Init() error {

	me.sidechainEthClient = &eurus_ethereum.EthClient{
		Protocol: me.Config.EthClientProtocol,
		IP:       me.Config.EthClientIP,
		Port:     me.Config.EthClientPort,
		ChainID:  big.NewInt(int64(me.Config.EthClientChainID)),
	}
	_, err := me.sidechainEthClient.Connect()
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Unable to connect to sidechain: ", err.Error())
		return err
	}

	me.mainnetEthClient = &eurus_ethereum.EthClient{
		Protocol: me.Config.MainnetEthClientProtocol,
		IP:       me.Config.MainnetEthClientIP,
		Port:     me.Config.MainnetEthClientPort,
		ChainID:  big.NewInt(int64(me.Config.MainnetEthClientChainID)),
	}
	_, err = me.mainnetEthClient.Connect()
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Unable to connect to mainnet: ", err.Error())
		return err
	}
	//Query mainnet platform wallet address
	mainnetSC, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(me.Config.EurusInternalConfigAddress), me.mainnetEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewEurusInternalConfig error")
	}
	addr1, err := mainnetSC.PlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "Getting mainnet platform wallet address from EurusInternalConfig")
	}
	me.MainnetPlatformWalletAddress = &addr1

	return err
}

func (me *WithdrawObserverSCProcessor) GetTransIdFromBurnCompleteEvent(transLog *types.Log) (*big.Int, error) {
	args, err := me.WithdrawContractAbi.Unpack("BurnCompletedEvent", transLog.Data)
	if err != nil {
		return nil, err
	}
	transId, ok := args[0].(*big.Int)
	if !ok {
		transId = big.NewInt(-1)
	}

	return transId, nil
}

func (me *WithdrawObserverSCProcessor) ParseWithdrawEvent(transLog *types.Log, requestTransId uint64) (*WithdrawEvent, error) {
	withdrawEvent := new(WithdrawEvent)
	fields, err := me.WithdrawContractAbi.Unpack("WithdrawEvent", transLog.Data)
	if err != nil {
		logError("Unpack failed", err, transLog, requestTransId, 0, "")
		return nil, err
	}
	if len(fields) < 5 {
		logError("Data field count invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}
	srcWallet, ok := fields[0].(common.Address)
	if !ok {
		logError("srcWallet field invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}
	withdrawEvent.SrcWallet = &srcWallet
	destWallet, ok := fields[1].(common.Address)
	if !ok {
		logError("destWallet field invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}
	withdrawEvent.DestWallet = &destWallet

	withdrawEvent.AssetName, ok = fields[2].(string)
	if !ok {
		logError("assetName field invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}

	withdrawEvent.BurnTransId, ok = fields[3].(*big.Int)
	if !ok {
		logError("transId field invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}
	withdrawEvent.Amount, ok = fields[4].(*big.Int)
	if !ok {
		logError("amount field invalid", err, transLog, requestTransId, 0, "")
		return nil, err
	}

	return withdrawEvent, nil
}

func (me *WithdrawObserverSCProcessor) GetTransLogTimestamp(transLog *types.Log) *time.Time {
	block, err := me.sidechainEthClient.GetBlockByNumber(big.NewInt(int64(transLog.BlockNumber)))
	var burnDate *time.Time = nil
	if err == nil {
		burnDate = new(time.Time)
		*burnDate = time.Unix(int64(block.Time()), 0)
	}
	return burnDate
}

func (me *WithdrawObserverSCProcessor) ConfirmBurn(burnTransId *big.Int, transLog *types.Log) (*types.Transaction, *eurus_ethereum.BesuReceipt, error) {
	transOpt, err := me.sidechainEthClient.GetNewTransactorFromPrivateKey(me.Config.HdWalletPrivateKey, me.sidechainEthClient.ChainID)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("WithdrawEvent - Unable to create smart contract instance: ", err.Error())
		return nil, nil, err
	}

	withdrawSC, err := contract.NewWithdrawSmartContract(me.Config.WithdrawSmartContractAddr, me.sidechainEthClient.Client)
	if err != nil {
		me.logSmartContractError("Unable to create smart contract instance", transLog.TxHash.Hex(), err)
		return nil, nil, err
	}

	//Withdraw confirmation concensus
	log.GetLogger(me.loggerName).Debugln("WithdrawEvent - ConfirmTransaction ", burnTransId.String())

	transOpt.GasLimit = uint64(me.Config.SideChainGasLimit)
	tx, err := withdrawSC.ConfirmTransaction(transOpt, burnTransId)
	if err != nil {
		me.logSmartContractError("Unable to invoke ConfirmTransaction", transLog.TxHash.Hex(), err)
		return nil, nil, err
	}

	log.GetLogger(me.loggerName).Debugln("ConfirmTransaction trans hash: ", tx.Hash().Hex(), " to withrdraw burn transId: ", burnTransId.String())
	receipt, err := me.sidechainEthClient.QueryEthReceipt(tx)
	if err != nil {
		me.logSmartContractError("Unable to get receipt for transaction", transLog.TxHash.Hex(), err)
		return nil, nil, err
	}
	receiptByte, _ := json.Marshal(receipt)

	log.GetLogger(me.loggerName).Debugln("ConfirmTransaction receipt: ", string(receiptByte))
	return tx, receipt, nil
}

func (me *WithdrawObserverSCProcessor) TransferTokenAtMainnet(burnTransHash string, reqTransHash string, destAddr *common.Address, assetName string, amount *big.Int, signData [][]byte) (*types.Transaction, error) {

	log.GetLogger(me.loggerName).Debugln("Going to transfer token to mainnet. Request trans hash: ", reqTransHash)
	mainnetPlatformWallet, err := mainnet_contract.NewEurusPlatformWallet(*me.MainnetPlatformWalletAddress, me.mainnetEthClient.Client)
	if err != nil {
		me.logSmartContractError("NewEurusPlatformWallet error", burnTransHash, err)
		return nil, err
	}

	//Getting sidechain InternalSmartContractConfig
	internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.Config.InternalSCConfigAddress), me.sidechainEthClient.Client)
	if err != nil {
		me.logSmartContractError("NewInternalSmartContractConfig error: ", burnTransHash, err)
		return nil, err
	}
	//Getting sidechain ERC20 address
	erc20Addr, err := internalSCConfig.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, assetName)
	if err != nil {
		me.logSmartContractError("GetErc20SmartContractAddrByAssetName error: ", burnTransHash, err)
		return nil, err
	}

	erc20SC, err := contract.NewEurusERC20(erc20Addr, me.sidechainEthClient.Client)
	if err != nil {
		me.logSmartContractError("NewEurusERC20 error: ", burnTransHash, err)
		return nil, err
	}
	//Getting sidechain ERC20 decimal point
	ourDecimalPoint, err := erc20SC.Decimals(&bind.CallOpts{})
	if err != nil {
		me.logSmartContractError("EurusERC20 Decimals error: ", burnTransHash, err)
		return nil, err
	}

	//Transfer token at mainnet
	var tx *types.Transaction
	var isSuccess bool = false
	var transOpt *bind.TransactOpts

	transOpt, err = me.mainnetEthClient.GetNewTransactorFromPrivateKey(me.Config.HdWalletPrivateKey, me.mainnetEthClient.ChainID)
	if err != nil {
		me.logSmartContractError("Get private key error", burnTransHash, err)
		return nil, err
	}
	if me.Config.MainnetTransferGasLimit > 0 {
		transOpt.GasLimit = me.Config.MainnetTransferGasLimit
	}
	log.GetLogger(me.loggerName).Debug("gas limit: ", transOpt.GasLimit)

	if me.Config.MainnetTransferGasTipCap > 0 {
		transOpt.GasTipCap = big.NewInt(me.Config.MainnetTransferGasTipCap)
	}
	log.GetLogger(me.loggerName).Debug("gas tip cap: ", transOpt.GasTipCap)

	if me.Config.MainnetTransferGasFeeCap > 0 {
		transOpt.GasFeeCap = big.NewInt(me.Config.MainnetTransferGasFeeCap)
	} else {
		gasFeeCap, err := me.mainnetEthClient.SuggestGasFeeCap()
		if err != nil {
			log.GetLogger(me.loggerName).Errorln("Unable to suggest gas fee cap. Request trans hash: ", reqTransHash)
		} else {
			transOpt.GasFeeCap = gasFeeCap
		}
	}
	log.GetLogger(me.loggerName).Debug("gas fee cap: ", transOpt.GasFeeCap)

	tx, err = mainnetPlatformWallet.Transfer(transOpt, common.HexToHash(reqTransHash), *destAddr, assetName, amount, ourDecimalPoint, signData)
	if err != nil {
		if strings.Contains(err.Error(), "failed to estimate gas needed") || strings.Contains(err.Error(), "intrinsic gas too low") {
			if strings.Contains(err.Error(), "failed to estimate gas needed") {
				me.logSmartContractError("NewEurusPlatformWallet estimate gas error. Using try and error logic to guess the gas limit", burnTransHash, err)
			}
		}
		me.logSmartContractError("NewEurusPlatformWallet gas limit at: "+strconv.FormatUint(transOpt.GasLimit, 10)+" transfer error: ", burnTransHash, err)
	} else {
		isSuccess = true
	}

	if !isSuccess {
		return nil, err
	}

	log.GetLogger(me.loggerName).Infoln("Transfer token to mainnet transaction broadcasted. Request trans hash: ", reqTransHash, " new trans hash: ", tx.Hash().Hex())

	return tx, nil

}

func (me *WithdrawObserverSCProcessor) logSmartContractError(message string, triggerTransHash string, err error) {
	errByte, _ := json.Marshal(err)

	log.GetLogger(me.loggerName).Errorf("%s: triggered trans hash: %s, error message:%s, JSON: %s\r\n", message, triggerTransHash, err.Error(), string(errByte))
}
