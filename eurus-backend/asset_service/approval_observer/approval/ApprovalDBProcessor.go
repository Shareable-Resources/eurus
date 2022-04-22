package approval

import (
	"encoding/json"

	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ApprovalDBProcessor struct {
	db         *database.Database
	slaveDb    *database.ReadOnlyDatabase
	Config     *ApprovalObserverConfig
	LoggerName string
}

func NewApprovalDBProcessor(db *database.Database, slaveDb *database.ReadOnlyDatabase, config *ApprovalObserverConfig, loggerName string) *ApprovalDBProcessor {
	processor := new(ApprovalDBProcessor)
	processor.db = db
	processor.slaveDb = slaveDb
	processor.Config = config
	processor.LoggerName = loggerName
	return processor
}

func (me *ApprovalDBProcessor) DbGetDecentralizedUserByWalletAddress(address string) (*user.User, error) {
	dbConn, err := me.slaveDb.GetConn()
	if err != nil {
		log.GetLogger(me.LoggerName).Errorln("Database Network Error: " + err.Error())
		return nil, errors.New("Database Network Error ") //for function check reload or not
	}
	user := new(user.User)
	tx := dbConn.Where("wallet_address = ?", strings.ToLower(address)).Find(user)
	err = tx.Error
	if err != nil {
		log.GetLogger(me.LoggerName).Errorln("Failed to find in Database. " + err.Error())
		return nil, err
	}
	return user, err
}

func (me *ApprovalDBProcessor) getCentralizedWithdrawRequestArgument(ext *ExtractedTransaction) (*common.Address, *big.Int, error) {
	args, err, extractState := ethereum.DefaultABIDecoder.DecodeABIInputArgument(ext.OriginalTransaction.Data(), "UserWallet", "submitWithdrawV1")
	if err != nil {
		return nil, nil, err
	}

	if extractState != ethereum.ExtractSuccess {
		return nil, nil, errors.New("Invalid to extract ABI")
	}
	val, ok := args["dest"]
	if !ok {
		return nil, nil, errors.New("Unable to find dest address from input argument")
	}

	addr, ok := val.(common.Address)
	if !ok {
		return nil, nil, errors.New("Invalid dest address")
	}
	val, ok = args["withdrawAmount"]
	if !ok {
		return nil, nil, errors.New("Unable to find amount from input argument")
	}
	val, ok = args["amountWithFee"]
	if !ok {
		return nil, nil, errors.New("Unable to find amount from input argument")
	}
	amount, ok := val.(*big.Int)
	if !ok {
		return nil, nil, errors.New("Invalid amount")
	}

	return &addr, amount, nil
}

func (me *ApprovalDBProcessor) ApprovalFail(receiptData *PendingPrewithdraw) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		log.GetLogger(me.LoggerName).Error("Unable to connect db", err.Error(), "This transaction failed to insert. Transaction hash : ", receiptData.RequestTransHash)
		return errors.New("Database Network Error: " + err.Error())
	}

	tx := dbConn.Where("request_trans_hash = ? and status = ? or status = ?", receiptData.RequestTransHash, asset.StatusPendingApproval, asset.StatusApproved).Updates(
		PendingPrewithdraw{
			Status: asset.StatusError,
		})
	err = tx.Error
	if err != nil {
		log.GetLogger(me.LoggerName).Errorln("Failed to change the transaction status to error. DB data : ", receiptData, "The hash of the transaction :", receiptData.RequestTransHash)
	}
	return err
}

func (me *ApprovalDBProcessor) InsertWithdrawFail(ext *ExtractedTransaction, receipt *eurus_ethereum.BesuReceipt, remarks asset.RemarksJson) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		log.GetLogger(me.LoggerName).Error("Unable to connect db", err.Error(), "This transaction failed to insert. Transaction hash : ", ext.TxHash)
		return errors.New("Database Network Error: " + err.Error())
	}

	if ext.User == nil {
		return errors.New("User object is nil for trans hash: " + ext.TxHash)
	}
	if ext.User.Id == 0 {
		return errors.New("User object is 0 for trans hash: " + ext.TxHash)
	}

	userGasUsed := me.extractFailedGasUsed(&receipt.Receipt)

	currTime := time.Now()
	failGasUsed := big.NewInt(0)
	failGasUsed = failGasUsed.SetUint64(receipt.GasUsed)

	gasFee := big.NewInt(0)
	gasFee = gasFee.Mul(ext.OriginalTransaction.GasPrice(), failGasUsed)

	insertData := new(PendingPrewithdraw)
	approvalCheck := dbConn.Where("request_trans_hash = ?", strings.ToLower(ext.TxHash)).Find(insertData)
	err = approvalCheck.Error
	if err != nil {
		return err
	}
	if insertData.RequestTransHash != "" {
		err = errors.New("duplicate key value violates unique constraint")
		return err
	}

	insertData.InitDate()
	insertData.CustomerId = uint64(ext.User.Id)
	insertData.CustomerType = asset.CustomerUser
	insertData.ApprovalWalletAddress = strings.ToLower(me.Config.ApprovalWalletAddress.String())
	insertData.RequestTransId = nil
	insertData.RequestTransHash = strings.ToLower(ext.TxHash)
	insertData.AssetName = ext.AssetName
	insertData.InnetFromAddress = strings.ToLower(ext.sender)
	insertData.SidechainGasUsed = decimal.NewFromBigInt(failGasUsed, 0)
	insertData.SidechainGasFee = decimal.NewFromBigInt(gasFee, 0)

	insertData.MainnetToAddress = ext.To
	insertData.Amount = decimal.NewFromBigInt(ext.Amount, 0)
	insertData.AdminFee = decimal.NewFromBigInt(ext.AdminFee, 0)
	insertData.Status = -asset.StatusPendingApproval
	insertData.GasPrice = decimal.NewFromBigInt(ext.OriginalTransaction.GasPrice(), 0)
	if userGasUsed != nil {
		insertData.UserGasUsed = decimal.NewFromBigInt(userGasUsed, 0)
	} else {
		insertData.UserGasUsed = insertData.SidechainGasUsed
	}

	err = dbConn.Transaction(
		func(dbTx *gorm.DB) error {
			pending := new(PendingPrewithdraw)
			mergeData := new(asset.WithdrawTransaction)
			tx := dbTx.Create(insertData)
			err = tx.Error
			if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint.") {
				log.GetLogger(me.LoggerName).Error("Unable to add transaction index to db", err.Error())
				return err
			} else if err != nil {
				log.GetLogger(me.LoggerName).Error("Unable to add transaction index to db", err.Error())
				return err
			}

			withdrawTable := new(asset.WithdrawTransaction)
			withdrawalCheck := dbConn.Where("request_trans_hash = ?", strings.ToLower(ext.TxHash)).Find(withdrawTable)
			err = withdrawalCheck.Error
			if err != nil {
				return err
			}
			if withdrawTable.RequestTransHash != "" {
				err = errors.New("duplicate key value violates unique constraint")
				return err
			}

			mergeData.Id = insertData.Id
			mergeData.CustomerId = insertData.CustomerId
			mergeData.CustomerType = insertData.CustomerType
			mergeData.AssetName = insertData.AssetName
			mergeData.Amount = insertData.Amount
			mergeData.ApprovalWalletAddress = strings.ToLower(insertData.ApprovalWalletAddress)
			mergeData.RequestTransId = insertData.RequestTransId
			mergeData.RequestTransHash = insertData.RequestTransHash
			mergeData.RequestDate = ext.TransDate
			mergeData.ReviewDate = &currTime
			mergeData.InnetFromAddress = strings.ToLower(insertData.InnetFromAddress)
			mergeData.MainnetToAddress = strings.ToLower(insertData.MainnetToAddress)
			mergeData.Status = -asset.StatusPendingApproval
			mergeData.CreatedDate = insertData.CreatedDate
			mergeData.LastModifiedDate = time.Now()
			mergeData.SidechainGasFee = insertData.SidechainGasFee
			mergeData.SidechainGasUsed = insertData.SidechainGasUsed
			mergeData.AdminFee = insertData.AdminFee
			mergeData.GasPrice = insertData.GasPrice
			mergeData.UserGasUsed = insertData.UserGasUsed

			receiptByte, _ := json.Marshal(receipt)
			log.GetLogger(me.LoggerName).Errorln("ConfirmBurn Receipt status error: ", insertData.RequestTransHash, " tx receipt: ", string(receiptByte),
				" trigger by tx: ", insertData.RequestTransHash)
			remarksJson := asset.NewRemarksJsonFromReceipt(receiptByte)

			mergeData.Remarks = remarksJson.String()

			tx = dbTx.Where("request_trans_hash = ?", insertData.RequestTransHash).Delete(pending)
			err = tx.Error
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Failed to add the transaction to withdraw_transaction. DB data : ", insertData, "The hash of the transaction :", insertData.RequestTransHash, " error: ", err)
			}

			tx = dbTx.Create(mergeData)
			err = tx.Error
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Failed to add the transaction to withdraw_transaction. DB data : ", insertData, "The hash of the transaction :", insertData.RequestTransHash, " error: ", err)
			}
			return err
		})

	return err
}

func (me *ApprovalDBProcessor) PendingSuccess(pendingPreWithdrawData *PendingPrewithdraw, comfirmTransaction *types.Transaction, receipt *eurus_ethereum.BesuReceipt) error {
	currTime := time.Now()
	dbConn, err := me.db.GetConn()
	if err != nil {
		log.GetLogger(me.LoggerName).Error("Unable to connect db", err.Error())
		return errors.New("Database Network Error: " + err.Error())
	}

	mergeData := new(asset.WithdrawTransaction)
	mergeData.InitDate()
	mergeData.Id = pendingPreWithdrawData.Id
	mergeData.CustomerId = pendingPreWithdrawData.CustomerId
	mergeData.CustomerType = pendingPreWithdrawData.CustomerType
	mergeData.AssetName = pendingPreWithdrawData.AssetName
	mergeData.Amount = pendingPreWithdrawData.Amount
	mergeData.ApprovalWalletAddress = strings.ToLower(pendingPreWithdrawData.ApprovalWalletAddress)
	mergeData.RequestTransId = pendingPreWithdrawData.RequestTransId
	mergeData.RequestTransHash = pendingPreWithdrawData.RequestTransHash
	mergeData.RequestDate = currTime
	mergeData.ReviewDate = &currTime
	mergeData.ReviewTransHash = strings.ToLower(comfirmTransaction.Hash().Hex())
	mergeData.ReviewedBy = strings.ToLower(me.Config.HdWalletAddress)
	mergeData.InnetFromAddress = strings.ToLower(pendingPreWithdrawData.InnetFromAddress)
	mergeData.MainnetToAddress = strings.ToLower(pendingPreWithdrawData.MainnetToAddress)
	mergeData.Status = asset.StatusApproved
	mergeData.SidechainGasUsed = pendingPreWithdrawData.SidechainGasUsed
	mergeData.SidechainGasFee = pendingPreWithdrawData.SidechainGasFee
	mergeData.AdminFee = pendingPreWithdrawData.AdminFee
	mergeData.GasPrice = pendingPreWithdrawData.GasPrice
	mergeData.UserGasUsed = pendingPreWithdrawData.UserGasUsed

	receiptByte, _ := json.Marshal(receipt)

	remarksJson := asset.NewRemarksJsonFromReceipt(receiptByte)

	mergeData.Remarks = remarksJson.String()

	err = dbConn.Transaction(
		func(dbTx *gorm.DB) error {
			pending := new(PendingPrewithdraw)
			tx := dbTx.Create(mergeData)
			err = tx.Error
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Failed to add the transaction to withdraw_transaction. DB data : ", pendingPreWithdrawData.Id, "The hash of the transaction :", pendingPreWithdrawData.RequestTransHash, " error: ", err)
				return err
			}
			tx = dbTx.Where("request_trans_hash = ?", pendingPreWithdrawData.RequestTransHash).Delete(pending)
			err = tx.Error
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Failed to add the transaction to withdraw_transaction. DB data : ", pendingPreWithdrawData.Id, "The hash of the transaction :", pendingPreWithdrawData.RequestTransHash, " error: ", err)
			}
			return err
		})

	if err != nil {
		log.GetLogger(me.LoggerName).Errorln("Failed to add the transaction to withdraw_transaction. DB data : ", pendingPreWithdrawData.Id, "The hash of the transaction :", pendingPreWithdrawData.RequestTransHash, " error: ", err)
	}
	return err
}

func (me *ApprovalDBProcessor) DbGetCenteralizedUserByWalletAddress(addr string) (*user.User, error) {
	var errReturn error = nil

	for i := 0; i < me.Config.RetryCount; i++ {
		errReturn = nil
		conn, err := me.slaveDb.GetConn()
		if err != nil {
			errReturn = errors.Wrap(err, "Get DB conn error")
			log.GetLogger(me.LoggerName).Errorln("Get DB conn error: ", err.Error())
			continue
		}
		var searchUser *user.User = new(user.User)
		dbTrans := conn.Where("wallet_address = ? AND is_metamask_addr = ?", strings.ToLower(addr), false).Find(&searchUser)
		err = dbTrans.Error
		if err != nil {
			errReturn = errors.Wrap(err, "Query user table error")
			log.GetLogger(me.LoggerName).Errorln("Query centralized user failed: ", errReturn.Error(), " input address: ", addr)
			continue
		}
		return searchUser, nil
	}

	return nil, errReturn
}

func (me *ApprovalDBProcessor) DbInsertPendingPrewithdraw(ext *ExtractedTransaction, receipt *eurus_ethereum.BesuReceipt) (*PendingPrewithdraw, error) {

	extractTopic, err := me.extractSuccessReceiptTopic(&receipt.Receipt)
	if err != nil {
		//This case suppose NOT to occurs as the receipt checking already be done before calling DbInsertPendingPrewithdraw()

		// if strings.Contains(err.Error(), "Error receipt") {
		// 	err := me.InsertWithdrawFail(ext, receipt)
		// 	if err != nil {
		// 		log.GetLogger(me.LoggerName).Errorln("Fail to insert fail data to DB: ", err.Error(), " trans hash: ", ext.TxHash)
		// 	}
		// 	return nil, err
		// }
		log.GetLogger(me.LoggerName).Errorln("Extract DB failed: ", err.Error(), " trans hash: ", ext.TxHash)
		return nil, err
	}

	var returnErr error
	for i := 0; i < me.Config.RetryCount; i++ {
		returnErr = nil
		dbConn, err := me.db.GetConn()
		if err != nil {
			log.GetLogger(me.LoggerName).Error("Unable to connect db", err.Error(), "This transaction insert failed. Transaction hash :", strings.ToLower(ext.TxHash))
			returnErr = errors.Wrap(err, "Database Network Error")
			continue
		}
		withdrawTable := new(asset.WithdrawTransaction)
		todata := new(PendingPrewithdraw)

		withdrawalCheck := dbConn.Where("request_trans_hash = ?", strings.ToLower(ext.TxHash)).Find(withdrawTable)
		err = withdrawalCheck.Error
		if err != nil {
			returnErr = errors.Wrap(err, "Query withdraw_transaction failed")
			continue
		}
		if withdrawTable.RequestTransHash != "" {
			err = errors.New("duplicate key value violates unique constraint")
			return nil, err
		}

		approvalCheck := dbConn.Where("request_trans_hash = ?", strings.ToLower(ext.TxHash)).Find(todata)
		err = approvalCheck.Error
		if err != nil {
			returnErr = errors.Wrap(err, "Query pending_prewithdraw error")
			continue
		}
		if todata.RequestTransHash != "" {
			err = errors.New("duplicate key value violates unique constraint")
			return nil, err
		}

		userObj := new(user.User)
		tx := dbConn.Where("wallet_address = ?", strings.ToLower(ext.sender)).Find(userObj)
		err = tx.Error
		if err != nil {
			returnErr = errors.Wrap(err, "Query user error")
			continue
		}

		todata.InitDate()
		todata.CustomerId = uint64(userObj.Id)
		todata.CustomerType = asset.CustomerUser

		todata.InnetFromAddress = strings.ToLower(extractTopic.SrcWallet.Hex())
		todata.MainnetToAddress = strings.ToLower(extractTopic.DestWallet.Hex())
		todata.ApprovalWalletAddress = strings.ToLower(me.Config.ApprovalWalletAddress.String())
		requestTransId := extractTopic.TransId.Uint64()
		todata.RequestTransId = &requestTransId
		todata.RequestTransHash = strings.ToLower(ext.TxHash)
		todata.AssetName = extractTopic.AssetName

		todata.Amount = decimal.NewFromBigInt(extractTopic.Amount, 0)

		todata.Status = asset.StatusPendingApproval //to be confirm
		todata.AdminFee = decimal.NewFromBigInt(extractTopic.FeeAmount, 0)

		gasUsed := new(big.Int).SetInt64(int64(receipt.GasUsed))
		gasPrice := new(big.Int).Set(ext.OriginalTransaction.GasPrice())
		gasFee := new(big.Int)
		gasFee.Mul(gasUsed, gasPrice)
		todata.SidechainGasUsed = decimal.NewFromInt(int64(receipt.GasUsed))
		todata.SidechainGasFee = decimal.NewFromBigInt(gasFee, 0)
		todata.GasPrice = decimal.NewFromBigInt(gasPrice, 0)
		if extractTopic.UserGasUsed != nil {
			todata.UserGasUsed = decimal.NewFromBigInt(extractTopic.UserGasUsed, 0)
		} else {
			todata.UserGasUsed = todata.SidechainGasUsed
		}

		query := dbConn.Create(todata)
		err = query.Error
		if err != nil {
			//error may be due to duplicate key. As another approval observer already insert this transaction
			returnErr = err
			log.GetLogger(me.LoggerName).Errorln("Error to insert to DB ", err, "Transaction hash :", strings.ToLower(ext.TxHash))
			continue
		}
		return todata, nil
	}

	return nil, returnErr
}

func (me *ApprovalDBProcessor) extractSuccessReceiptTopic(receipt *types.Receipt) (*ApprovalWalletTopic, error) {
	approvalWalletTopic := new(ApprovalWalletTopic)
	abi := eurus_ethereum.DefaultABIDecoder.GetABI("ApprovalWallet")
	submitWithdrawEvent, ok := abi.Events["SubmitWithdraw"]
	if !ok {
		log.GetLogger(me.LoggerName).Errorln("ApprovalWallet SubmitWithdraw event not found in ABI. Program aborted")
	}

	userWalletAbi := eurus_ethereum.DefaultABIDecoder.GetABI("UserWallet")
	withdrawRequestEvent, ok := userWalletAbi.Events["WithdrawRequestEvent"]
	if !ok {
		log.GetLogger(me.LoggerName).Errorln("UserWallet WithdrawRequestEvent event not found in ABI. Program aborted")
	}

	if len(receipt.Logs) == 0 {
		return nil, errors.New("Topic not found")
	}
	var isFound bool
	for _, topicLog := range receipt.Logs {
		if topicLog.Topics[0] != submitWithdrawEvent.ID {
			continue
		}

		if receipt.Status == 0 {
			return nil, errors.New("Error receipt")
		}

		args, err := eurus_ethereum.DefaultABIDecoder.DecodeABIEventData(topicLog.Data, "ApprovalWallet", "SubmitWithdraw")
		if err == nil {
			approvalWalletTopic.SrcWallet = *ethereum.HashToAddress(&topicLog.Topics[1])
			approvalWalletTopic.DestWallet = *ethereum.HashToAddress(&topicLog.Topics[2])
			approvalWalletTopic.SubmitterAddress = *ethereum.HashToAddress(&topicLog.Topics[3])
			approvalWalletTopic.AssetName = args[0].(string)
			approvalWalletTopic.TransId = args[1].(*big.Int)
			approvalWalletTopic.Amount = args[2].(*big.Int)
			approvalWalletTopic.FeeAmount = args[3].(*big.Int)
			isFound = true
			break
		}
	}

	if isFound {
		for _, topicLog := range receipt.Logs {
			if topicLog.Topics[0] == withdrawRequestEvent.ID {
				//Only centralized user have this event
				userGasUsed := topicLog.Topics[2].Big()
				approvalWalletTopic.UserGasUsed = userGasUsed
			}
		}
		return approvalWalletTopic, nil
	}
	return nil, errors.New("event not match")

}

func (me *ApprovalDBProcessor) extractFailedGasUsed(receipt *types.Receipt) *big.Int {
	gasUsed := big.NewInt(0)
	gasUsed.SetUint64(receipt.GasUsed)
	return gasUsed

	// userWalletProxyAbi := eurus_ethereum.DefaultABIDecoder.GetABI("UserWalletProxy")
	// withdrawFailedEvent, ok := userWalletProxyAbi.Events["SubmitWithdrawFailed"]
	// if !ok {
	// 	log.GetLogger(me.LoggerName).Fatalln("UserWalletProxy SubmitWithdrawFailed event not found in ABI. Program aborted")
	// }

	// for _, topicLog := range receipt.Logs {
	// 	if topicLog.Topics[0] == withdrawFailedEvent.ID {
	// 		userGasUsed := topicLog.Topics[2].Big()
	// 		return userGasUsed
	// 	}
	// }
	// return nil
}
