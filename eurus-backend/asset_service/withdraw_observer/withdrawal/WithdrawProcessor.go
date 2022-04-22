package withdrawal

import (
	"encoding/json"
	"strings"

	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"math/big"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

type WithdrawProcessor struct {
	context                *WithdrawProcessorContext
	smartContractProcessor *WithdrawObserverSCProcessor
	config                 *WithdrawObserverConfig
	loggerName             string
}

func NewWithdrawProcessor(config *WithdrawObserverConfig, loggerName string) *WithdrawProcessor {
	processor := new(WithdrawProcessor)

	// context.db will be set later
	processor.context = new(WithdrawProcessorContext)
	processor.context.retrySetting = config
	processor.context.loggerName = loggerName

	processor.config = config
	processor.smartContractProcessor = NewWithdrawObserverSCProcessor(config, loggerName)
	processor.loggerName = loggerName
	return processor
}

func (me *WithdrawProcessor) Init(db *database.Database, slaveDb *database.ReadOnlyDatabase) error {
	err := me.smartContractProcessor.Init()
	if err != nil {
		return err
	}
	me.context.db = db
	me.context.slaveDb = slaveDb
	return nil
}

func (me *WithdrawProcessor) ProcessTransaction(block *types.Block) {
	if block.Transactions().Len() == 0 {
		return
	}

	for _, trans := range block.Transactions() {

		args, _, status := eurus_ethereum.DefaultABIDecoder.DecodeABIInputArgument(trans.Data(), "WithdrawSmartContract", "confirmTransaction")
		if status == eurus_ethereum.ExtractSuccess {
			receipt, err := me.smartContractProcessor.sidechainEthClient.QueryEthReceiptWithSetting(trans, 1, 20)
			if err != nil {
				continue
			}
			if receipt.Status == 0 {
				transId := args["transactionId"].(*big.Int)
				log.GetLogger(me.loggerName).Errorln("Approval failed. Receipt status is 0. Trans hash: ", trans.Hash().Hex(),
					" request trans id: ", transId.String())
				continue
			}

			for _, logMessage := range receipt.Logs {
				me.processLogEvent(logMessage)
			}
		}
	}
}

func (me *WithdrawProcessor) processLogEvent(transLog *types.Log) {
	if transLog.Topics[0] == me.config.WithdrawEventTopic {
		isHandled := me.processWithdrawEventLog(transLog)

		if !isHandled {
			log.GetLogger(me.loggerName).Debugln("get block....... status change 40")
			me.retryProcessing(transLog, func(log *types.Log) bool {
				return me.processWithdrawEventLog(log)
			})
		}

	} else if transLog.Topics[0] == me.config.BurnCompletedEventTopic {
		log.GetLogger(me.loggerName).Debugln("Burn completed")
		isHandled := me.processBurnCompletedEventLog(transLog)

		if !isHandled {
			me.retryProcessing(transLog, func(log *types.Log) bool {
				return me.processBurnCompletedEventLog(log)
			})
		}
	}
}

func (me *WithdrawProcessor) processWithdrawEventLog(transLog *types.Log) bool {

	if len(transLog.Topics) >= 2 {
		topicHash1 := transLog.Topics[1]
		var approvalWallet *common.Address = eurus_ethereum.HashToAddress(&topicHash1)

		topicHash2 := transLog.Topics[2]
		var requestTransId *uint256.Int = eurus_ethereum.HashToInt256(&topicHash2)
		withdrawEvent, err := me.smartContractProcessor.ParseWithdrawEvent(transLog, requestTransId.Uint64())
		if err != nil {
			log.GetLogger(me.loggerName).Errorln("ParseWithdrawEvent Failed.", err.Error())
			return true
		}
		log.GetLogger(me.loggerName).Debugf("WithdrawEvent -  Approval wallet: %s, approve transId: %s, src wallet: %s, dest mainnet wallet: %s, asset name: %s, trans Id: %d, amount: %d\r\n",
			approvalWallet.String(), requestTransId.String(), withdrawEvent.SrcWallet.String(), withdrawEvent.DestWallet.String(), withdrawEvent.AssetName, withdrawEvent.BurnTransId.Uint64(), withdrawEvent.Amount.Uint64())

		//Query DB for the record
		dbTransRecord, err := DbQueryWithdrawTransaction(me.context, approvalWallet.String(), requestTransId.Uint64())
		if err != nil || dbTransRecord == nil {
			if err == nil {
				err = errors.New("Record not found")
			}
			log.GetLogger(me.loggerName).Errorln("WithdrawEvent - Query DB error: ", err.Error(), " trigger by tx: ", transLog.TxHash.Hex())
			return false
		}

		//Withdraw confirmation concensus
		log.GetLogger(me.loggerName).Debugln("WithdrawEvent - ConfirmTransaction ", withdrawEvent.BurnTransId.String())

		tx, receipt, err := me.smartContractProcessor.ConfirmBurn(withdrawEvent.BurnTransId, transLog)
		if err != nil {
			remarksJson := asset.NewRemarksJsonFromError(err)
			err := DbUpdateWithdrawTransactionConfirmBurnToError(me.context, dbTransRecord.Id, nil, asset.StatusApproved, remarksJson.String())
			if err != nil {
				log.GetLogger(me.loggerName).Errorln("DbUpdateWithdrawTransactionConfirmBurnToError on receipt status error:  ", tx.Hash().Hex(), " trigger by tx: ", transLog.TxHash.Hex())
			}
			return false
		}

		log.GetLogger(me.loggerName).Debugln("Receipt status : ", receipt.Status, receipt.TxHash.Hex())

		if receipt.Status != 1 {
			receiptByte, _ := json.Marshal(receipt)
			log.GetLogger(me.loggerName).Errorln("ConfirmBurn Receipt status error: ", tx.Hash().Hex(), " tx receipt: ", string(receiptByte),
				" trigger by tx: ", transLog.TxHash.Hex(), " burn trans Id:", withdrawEvent.BurnTransId.String())
			remarksJson := asset.NewRemarksJsonFromReceipt(receiptByte)
			err := DbUpdateWithdrawTransactionConfirmBurnToError(me.context, dbTransRecord.Id, nil, asset.StatusApproved, remarksJson.String())
			if err != nil {
				log.GetLogger(me.loggerName).Errorln("DbUpdateWithdrawTransactionConfirmBurnToError on receipt status error:  ", tx.Hash().Hex(), " tx receipt: ", string(receiptByte), " trigger by tx: ", transLog.TxHash.Hex())
			}
			return true
		} else {
			log.GetLogger(me.loggerName).Infof("Confirm withdraw success. burn trans Id: %s\r\n", withdrawEvent.BurnTransId.String())
			//Update DB status
			err := DbUpdateWithdrawTransactionToBurnConfirming(me.context, dbTransRecord.Id, withdrawEvent.BurnTransId)
			if err != nil {
				log.GetLogger(me.loggerName).Errorln("Update transaction status error: ", tx.Hash().Hex(),
					" error: ", err.Error(), " trigger by tx: ", transLog.TxHash.Hex(), " burn trans Id:", withdrawEvent.BurnTransId.String())
			}

		}
	} else {
		log.GetLogger(me.loggerName).Errorf("Invalid topic length for transaction: %s. Got length: %d\r\n ", transLog.TxHash.Hex(), len(transLog.Topics))

	}
	return true
}

func (me *WithdrawProcessor) processBurnCompletedEventLog(transLog *types.Log) bool {

	if len(transLog.Topics) < 3 {
		log.GetLogger(me.loggerName).Errorln("Invalid topic count: ", len(transLog.Topics), " tx hash: ", transLog.TxHash.Hex())
		return true
	}
	hashData := transLog.Topics[1]
	approvalWallet := eurus_ethereum.HashToAddress(&hashData)
	data := transLog.Topics[2]
	requestTransId := eurus_ethereum.HashToInt256(&data)
	transId, err := me.smartContractProcessor.GetTransIdFromBurnCompleteEvent(transLog)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Unable to decode log data: ", err.Error(), " transHash: ", transLog.TxHash.Hex())
		return true
	}

	log.GetLogger(me.loggerName).Debugf("Burn completed: approval transId: %d, approval wallet: %s, burn trans Id: %s\r\n", requestTransId.Uint64(), approvalWallet.String(), transId.String())

	dbTrans, err := DbQueryWithdrawTransaction(me.context, approvalWallet.String(), requestTransId.Uint64())
	if err != nil {
		log.GetLogger(me.loggerName).Debugf("Burn completed: QueryWithdrawTransaction error: %s, requestTransId: %s burnTransId: %s, approval wallet: %s, burn trans Id: %s\r\n", err.Error(), requestTransId.Uint64(), approvalWallet.String(), transId.String())
		return false
	} else if dbTrans == nil {
		log.GetLogger(me.loggerName).Debugf("Burn completed: QueryWithdrawTransaction record not found, approval transId: %s, approval wallet: %s, burn trans Id: %s\r\n", requestTransId.Uint64(), approvalWallet.String(), transId.String())
		return false
	}

	burnDate := me.smartContractProcessor.GetTransLogTimestamp(transLog)

	platformWalletAddr := me.smartContractProcessor.MainnetPlatformWalletAddress

	_ = DbUpdateWithdrawTransactionToBurn(me.context, transLog, requestTransId,
		dbTrans.Id, transId.Uint64(), &transLog.TxHash, burnDate, platformWalletAddr)

	//fmt.Println(rowAffected)
	//if rowAffected == 0 {
	//	log.GetLogger(me.loggerName).Debugf("Burn completed: Another withdraw observer handling the event. Jump to next step, approval transId: %s, approval wallet: %s, burn trans Id: %s\r\n", requestTransId.Uint64(), approvalWallet.String(), transId.String())
	//
	//}

	//Query burned Record

	//Query the data and insert to withdraw_request table ----- RYAN
	log.GetLogger(me.loggerName).Debugln("Ready To insert sign data and ori data")
	withdrawTx, rowCount, err := DbInsertWithdrawRequest(me, transLog, *approvalWallet, requestTransId)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Failed to handle transaction after burned. ", err, " request trans hash: ", dbTrans.RequestTransHash)
		return false
	}
	if rowCount == 0 {
		log.GetLogger(me.loggerName).Errorln("Sign request insert failed. Failed to handle transaction after burned. request trans hash: ", dbTrans.RequestTransHash)
		return false
	}

	rowCount, err = DbCountWithdrawRequests(me.context, withdrawTx.Id)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("DbCountWithdrawRequests error: ", err, " request trans hash: ", withdrawTx.RequestTransHash)
		return false
	}

	if rowCount < 5 {
		log.GetLogger(me.loggerName).Debugln("row count is : ", rowCount, ", skip", " request trans hash: ", withdrawTx.RequestTransHash)
		return true
	}

	log.GetLogger(me.loggerName).Debugln("Try to change status request trans hash: ", withdrawTx.RequestTransHash)

	rowAffected, err := DbUpdateWithdrawTransactionToConfirmingTransfer(me.context, withdrawTx)
	if err != nil {
		return false
	}

	if rowAffected == 0 {
		log.GetLogger(me.loggerName).Infoln("Transaction is already in state 60, skip tranfer token at mainnet. request trans hash: ", withdrawTx.RequestTransHash)
		return true
	}

	//
	//if rowAffected == 0 {
	//	log.GetLogger(me.loggerName).Debugf("Burn completed: Another withdraw observer handling the event. Procedure skipped, approval transId: %s, approval wallet: %s, burn trans Id: %s\r\n", requestTransId.Uint64(), approvalWallet.String(), transId.String())
	//	return true
	//}

	//Asynchorous call mainnet
	go me.TransferTokenToMainnet(withdrawTx)

	return true

}

func (me *WithdrawProcessor) TransferTokenToMainnet(withdrawTx *asset.WithdrawTransaction) {

	signData, _, err := DBQuerySignedData(me.context, withdrawTx)
	if err != nil {
		log.GetLogger(me.loggerName).Error("Fail to get the sign byte. ", err, " request trans hash: ", withdrawTx.RequestTransHash)
		return
	}
	if len(signData) < 5 {
		log.GetLogger(me.loggerName).Error("Not enough signature. Count: ", len(signData), " request trans hash: ", withdrawTx.RequestTransHash)
		return
	}

	log.GetLogger(me.loggerName).Infof("Going to transfer to mainnet. burn trans Id: %d,  Request transId: %d. Approval Wallet: %s.", withdrawTx.BurnTransId, withdrawTx.RequestTransId, withdrawTx.ApprovalWalletAddress)

	mainnetDestAddr := common.HexToAddress(withdrawTx.MainnetToAddress)

	conn, err := me.context.db.GetConn()
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Fail to get db conn", err.Error())
		return
	}

	dBtx := conn.Begin()

	withdrawTransaction := new(asset.WithdrawTransaction)
	nTx := dBtx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("approval_wallet_address = ? AND  request_trans_id = ? and status = ?", strings.ToLower(withdrawTx.ApprovalWalletAddress), withdrawTx.Id, asset.StatusConfirmingTransfer).Find(withdrawTransaction)
	err = nTx.Error
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Fail to lock table", err.Error())
		dBtx.Rollback()
		return
	}
	withdrawData := &asset.WithdrawTransaction{
		Status: asset.StatusTransferProcessing,
	}

	nDbtx := dBtx.Where("id = ? AND status = ? ", withdrawTx.Id, asset.StatusConfirmingTransfer).Updates(withdrawData)

	err = nDbtx.Error
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Fail to update withdraw transaction to Completed state error", err.Error())
		dBtx.Rollback()
		return
	}
	if err := dBtx.Commit().Error; err != nil {
		log.GetLogger(me.loggerName).Errorln("Fail to commit transaction error", err.Error())
		return
	}

	var receipt *ethereum.BesuReceipt
	var tx *types.Transaction
	for i := 0; i < me.config.RetryCount; i++ {
		tx, err = me.smartContractProcessor.TransferTokenAtMainnet(withdrawTx.BurnTransHash, withdrawTx.RequestTransHash, &mainnetDestAddr, withdrawTx.AssetName, withdrawTx.Amount.BigInt(), signData)
		if err == nil {
			break
		}
		log.GetLogger(me.loggerName).Errorln("TransferTokenAtMainnet error: ", err, "Request trans hash: ", withdrawTx.RequestTransHash, " retry count: ", i)
		if i < me.config.RetryCount {
			time.Sleep(me.config.GetRetryInterval() * 5)
		}
	}

	if err != nil {
		remarksJson := asset.NewRemarksJsonFromError(err)
		remarksByte, _ := json.Marshal(remarksJson)
		for {
			var isSuccess bool = false
			dbConn2, err := me.context.db.GetConn()
			if err != nil {
				log.GetLogger(me.loggerName).Errorln("Fail to get db conn", err.Error())
				return
			}

			err2 := dbConn2.Transaction(func(dbTx2 *gorm.DB) error {
				err1 := DbUpdateWithdrawTransactionToError(me.context, withdrawTx.Id, nil, asset.StatusTransferProcessing, string(remarksByte), dbTx2)
				if err1 != nil {
					log.GetLogger(me.loggerName).Errorln("DbUpdateWithdrawTransactionToError error: ", err1, " retrying")
					return err1
				}
				isSuccess = true
				return nil
			})

			if !isSuccess || err2 != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			break

		}
		log.GetLogger(me.loggerName).Debug("Insert DB for WithdrawTransactionToError return. Request trans hash: ", withdrawTx.RequestTransHash)
		return
	} else {

		_ = DbUpdateMainnetTransferTransHash(me.context, withdrawTx.RequestTransHash, withdrawTx.Id, tx.Hash())
		receipt, err = me.smartContractProcessor.mainnetEthClient.QueryEthReceiptWithSetting(tx, 3, -1)
		if err != nil {
			log.GetLogger(me.loggerName).Errorln("QueryEthReceipt error: ", err, " mainnet transfer transaction: ", tx.Hash().Hex(), " request trans hash: ", withdrawTx.RequestTransHash)
			dbConn, _ := me.context.db.GetConn()
			transferTransHash := tx.Hash()
			remarksJson := asset.NewRemarksJsonFromError(err)
			_ = DbUpdateWithdrawTransactionToError(me.context, withdrawTx.Id, &transferTransHash, asset.StatusTransferProcessing, remarksJson.String(), dbConn)
			return
		}
	}

	var logStatus string = "Completed"
	if receipt.Status != 1 {
		logStatus = "Error"
	}

	err = DbInsertAssetAllocationCost(me.context, receipt)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("Unable to insert withdrawal asset allocation cost into DB")
	}

	var err1 error
	receiptByte, _ := json.Marshal(receipt)
	log.GetLogger(me.loggerName).Infof("Transfer to mainnet %s. Request Trans Hash: %s. Trans Hash: %s. Request transId: %d. Approval Wallet: %s. Receipt: %s\r\n",
		logStatus, withdrawTx.RequestTransHash, tx.Hash().Hex(), withdrawTx.RequestTransId, withdrawTx.ApprovalWalletAddress, string(receiptByte))
	txHash := tx.Hash()
	for {
		log.GetLogger(me.loggerName).Debugln("Receipt Status isssss: ", receipt.Status, "ID is ", withdrawTx.RequestTransId, " request trans hash: ", withdrawTx.RequestTransHash)
		if receipt.Status == 1 {
			_ = DbUpdateWithdrawTransactionToComplete(me.context, withdrawTx.Id, &txHash, me.smartContractProcessor.MainnetPlatformWalletAddress)
			err1 = nil
		} else {
			log.GetLogger(me.loggerName).Debugln("Receipt Status faillllll: ", receipt.Status, "ID is ", withdrawTx.RequestTransId, " request trans hash: ", withdrawTx.RequestTransHash)
			remarksJson := asset.NewRemarksJsonFromString("Transaction status is 0x0. Trans hash: " + txHash.Hex())
			dbConn, _ := me.context.db.GetConn()
			err1 = DbUpdateWithdrawTransactionToError(me.context, withdrawTx.Id, &txHash, asset.StatusTransferProcessing, remarksJson.String(), dbConn)
		}

		if err1 == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

}

func (me *WithdrawProcessor) processResumbitTransferTokenToMainnet(requestTransHash string) error {

	withdrawTx, err := DbQueryBurnedTransaction(me.context, requestTransHash)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("DbQueryBurnedTransaction failed: ", err, " request trans hash: ", requestTransHash)
		return errors.Wrap(err, "DbQueryBurnedTransaction error")
	}
	if withdrawTx.Id == 0 {
		log.GetLogger(me.loggerName).Errorln("Cannot find the transaction or the transaction is not in state 70. request trans hash: ", requestTransHash)
		return errors.New("Cannot find the transaction or the transaction is not in state 70")
	}

	err = DbUpdateErrorTransactionToConfirmingTransfer(me.context, requestTransHash)
	if err != nil {
		log.GetLogger(me.loggerName).Errorln("DbUpdateErrorTransactionToConfirmingTransfer failed: ", err, " request trans hash: ", requestTransHash)
		return errors.Wrap(err, "DbUpdateErrorTransactionToConfirmingTransfer failed")
	}

	log.GetLogger(me.loggerName).Infoln("Manually trigger resubmit transfer request to mainnet. Request trans hash: ", requestTransHash)
	go me.TransferTokenToMainnet(withdrawTx)
	return nil
}

func (me *WithdrawProcessor) retryProcessing(transLog *types.Log, functor func(*types.Log) bool) {
	for i := 0; i < 10; i++ {
		log.GetLogger(me.loggerName).Infof("Wait 2 seconds to retry tran hash: %s\r\n", transLog.TxHash.Hex())
		time.Sleep(2 * time.Second)
		isHandled := functor(transLog)
		if isHandled {
			break
		}
	}
}
