package userObserver

import (
	//"eurus-backend/asset_service/asset"
	"encoding/json"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"

	//"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	//"fmt"
	//"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UserObserverProcessor struct {
	config      *UserObserverConfig
	context     *UserObserverContext //context has replaced db
	scProcessor *UserObserverSCProcessor
	// transferRequestEvent common.Hash
}

func NewUserObserverProcessor(db *database.Database, slaveDb *database.ReadOnlyDatabase, config *UserObserverConfig, scProcessor *UserObserverSCProcessor, loggerName string) *UserObserverProcessor {
	processor := new(UserObserverProcessor)
	processor.scProcessor = scProcessor
	processor.config = config
	processor.context = NewUserObserverContext(db, slaveDb, config, loggerName)
	return processor
}

func (me *UserObserverProcessor) BlockHandler(block *types.Block) {
	if block.Transactions().Len() > 0 { //check if block has tx
		me.ProcessBlock(block)
	}

}

func (me *UserObserverProcessor) ProcessBlock(block *types.Block) {
	transactions := block.Transactions()
	for _, transaction := range transactions {

		centralizedUser, err := DbGetCentralizedUser(me.context, transaction)
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("Unable to get user address: ", err.Error())
			continue
		}

		if centralizedUser == nil {
			continue
		}

		_, err, _ = ethereum.DefaultABIDecoder.DecodeABIInputArgument(transaction.Data(), "UserWallet", "requestTransferV1")
		if err != nil {

			log.GetLogger(me.context.LoggerName).Errorln("Unable to decode ABI for trans hash: ", transaction.Hash().Hex(), " Error: ", err)
			continue

		}

		receipt, err := me.scProcessor.sideChainEthClient.QueryEthReceiptWithSetting(transaction, 1, 20)
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("Query ETH Receipt failed.", err.Error(), "Trans hash is ", transaction.Hash().Hex(), "transaction to : ", transaction.To().Hex(), "Gas price : ", transaction.GasPrice(),
				"Original transaction : ", transaction)
			continue
		}

		if receipt.Status == 0 {
			continue
		}

		eventAbi := ethereum.DefaultABIDecoder.GetABI("UserWallet").Events["TransferRequestEvent"]
		for _, transLog := range receipt.Logs {

			if transLog.Topics[0] == eventAbi.ID {
				me.processLog(transLog, transaction, *transaction.To())
				break
			}
		}
	}
}

func (me *UserObserverProcessor) processLog(transLog *types.Log, requestTx *types.Transaction, userWalletAddr common.Address) {

	ethClient := me.scProcessor.sideChainEthClient
	tranId, err := me.scProcessor.GetTransIdFromTransferRequestEvent(transLog, requestTx)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Failed to get the transId", err)
	}

	if tranId != nil {
		tx, err := me.scProcessor.ConfirmTransactionWithTransID(tranId, transLog, userWalletAddr)
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("Fail to confirm the transaction.", err)
		}
		if tx != nil {
			confirmTxReceipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
			if err != nil {
				log.GetLogger(me.context.LoggerName).Errorln("Failed to get the receipt", err)
			}
			if confirmTxReceipt.Status == 0 {
				receiptData, _ := json.Marshal(confirmTxReceipt)
				log.GetLogger(me.context.LoggerName).Errorln("Fail to confirm the tx", err, " Hash: ", transLog.TxHash.String(),
					"Receipt: ", string(receiptData), " request trans hash: ", requestTx.Hash().Hex())
			} else {
				log.GetLogger(me.context.LoggerName).Infoln("Confirmation trans hash: ", tx.Hash().Hex(), " for request trans Id: ", requestTx.Hash().Hex())
			}

		}
	}

}
