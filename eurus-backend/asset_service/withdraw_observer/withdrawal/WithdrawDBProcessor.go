package withdrawal

import (
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"gorm.io/gorm"
)

func DbQueryWithdrawTransaction(context *WithdrawProcessorContext, approvalWalletAddr string, approvalTransId uint64) (*asset.WithdrawTransaction, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return nil, err
	}
	transRecord := new(asset.WithdrawTransaction)
	conn.Where("approval_wallet_address = ? AND request_trans_id = ?", strings.ToLower(approvalWalletAddr), approvalTransId).First(transRecord)

	if transRecord.Id == 0 {
		return nil, nil
	}
	return transRecord, nil
}

func DbUpdateWithdrawTransactionToBurnConfirming(context *WithdrawProcessorContext, transRecordId uint64, burnTransId *big.Int) error {
	var err error
	var conn *gorm.DB
	for i := 0; i < context.retrySetting.GetRetryCount(); i++ {
		conn, err = context.db.GetConn()
		if err != nil {
			time.Sleep(context.retrySetting.GetRetryInterval() * time.Second)
			continue
		}
		withdrawData := &asset.WithdrawTransaction{
			Status:      asset.StatusBurnConfirming,
			BurnTransId: burnTransId.Uint64(),
			DbModel:     database.DbModel{LastModifiedDate: time.Now()},
		}

		conn = conn.Where("id = ? AND ABS(status) <= ?", transRecordId, asset.StatusApproved).Updates(withdrawData)
		err = conn.Error
		if err == nil {
			break
		}
		log.GetLogger(context.loggerName).Errorln("Update transaction status error: ", err.Error(), " DB ID: ", transRecordId)

		time.Sleep(context.retrySetting.GetRetryInterval() * time.Second)
	}

	return err

}

func DbUpdateWithdrawTransactionToBurn(context *WithdrawProcessorContext, transLog *types.Log, requestTransId *uint256.Int,
	transRecordId uint64, burnTransId uint64, burnHash *common.Hash, burnDate *time.Time, mainnetFromAddr *common.Address) int {
	var err error
	var dbConn *gorm.DB
	for {
		dbConn, err = context.db.GetConn()
		if err != nil {
			logError("GetConn error", err, transLog, requestTransId.Uint64(), burnTransId, "")
			time.Sleep(time.Second)
			continue
		}

		withdrawData := &asset.WithdrawTransaction{
			Status:        asset.StatusBurned,
			BurnDate:      burnDate,
			BurnTransHash: burnHash.String(),
			BurnTransId:   burnTransId,
			DbModel:       database.DbModel{LastModifiedDate: time.Now()},
		}

		//Confirm status maybe skipped if the withdrawSmartContract only require one concensus
		dbConn = dbConn.Where("id = ? AND ABS(status) < ? ", transRecordId, asset.StatusBurned).Updates(withdrawData)

		err = dbConn.Error
		if err != nil {
			logError("Update burn status error", err, transLog, requestTransId.Uint64(), burnTransId, "")
			time.Sleep(time.Second)
			continue
		}
		break

	}
	return int(dbConn.RowsAffected)
}

func DbUpdateWithdrawTransactionConfirmBurnToError(context *WithdrawProcessorContext, transRecordId uint64, transferTransHash *common.Hash, fromState asset.WithdrawStatus, remarks string) error {
	var err error
	var dbConn *gorm.DB

	for i := 0; i < context.retrySetting.GetRetryCount(); i++ {
		dbConn, err = context.db.GetConn()
		if err != nil {
			time.Sleep(context.retrySetting.GetRetryInterval() * time.Second)
			continue
		}
		var transferHash string
		var now *time.Time = nil
		if transferTransHash != nil {
			transferHash = strings.ToLower(transferTransHash.Hex())
			tempNow := time.Now()
			now = &tempNow
		}
		//Confirm status maybe skipped if the withdrawSmartContract only require one concensus
		//conn = conn.Model(WithdrawTransaction{}).Where("id = ? AND ABS(status) <= ?", transRecordId, fromState).Updates(
		//  map[string]interface{}{
		//    "status":             fromState * StatusError,
		//    "mainnet_trans_hash": transferHash,
		//    "mainnet_trans_date": now,
		//    "remarks":            remarks,
		//    "last_modified_date"  : &now,
		//  })

		withdrawData := &asset.WithdrawTransaction{
			Status:           fromState * asset.StatusError,
			MainnetTransHash: transferHash,
			MainnetTransDate: now,
			Remarks:          remarks,
			DbModel:          database.DbModel{LastModifiedDate: time.Now()},
		}

		dbConn = dbConn.Where("id = ? AND ABS(status) <= ?", transRecordId, fromState).Updates(withdrawData)
		err = dbConn.Error
		if err == nil {
			break
		}

		time.Sleep(context.retrySetting.GetRetryInterval() * time.Second)
	}
	return err
}

func DbUpdateWithdrawTransactionToError(context *WithdrawProcessorContext, transRecordId uint64, transferTransHash *common.Hash, fromState asset.WithdrawStatus, remarks string, dBtx *gorm.DB) error {
	var err error

	var transferHash string
	var now *time.Time = nil
	if transferTransHash != nil {
		transferHash = strings.ToLower(transferTransHash.Hex())
		tempNow := time.Now()
		now = &tempNow
	}
	withdrawData := &asset.WithdrawTransaction{
		Status:           fromState * asset.StatusError,
		MainnetTransHash: transferHash,
		MainnetTransDate: now,
		Remarks:          remarks,
		DbModel:          database.DbModel{LastModifiedDate: time.Now()},
	}

	err = dBtx.Where("id = ? AND ABS(status) <= ?", transRecordId, fromState).Updates(withdrawData).Error
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to update withdraw transaction to error, error:", err.Error())
		return err
	}

	return err
}

func DbUpdateWithdrawTransactionToComplete(context *WithdrawProcessorContext, transRecordId uint64, transferHash *common.Hash, mainnetFromAddr *common.Address) int {

	var err error
	conn, err := context.db.GetConn()
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to Connect Db ", err)
		return 0
	}
	now := time.Now()
	//Confirm status maybe skipped if the withdrawSmartContract only require one concensus
	withdrawData := &asset.WithdrawTransaction{
		Status:             asset.StatusCompleted,
		MainnetTransHash:   strings.ToLower(transferHash.String()),
		MainnetTransDate:   &now,
		MainnetFromAddress: strings.ToLower(mainnetFromAddr.Hex()),
		DbModel:            database.DbModel{LastModifiedDate: time.Now()},
	}

	nDbtx := conn.Where("id = ? AND status = ? ", transRecordId, asset.StatusTransferProcessing).Updates(withdrawData)

	err = nDbtx.Error
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to update withdraw transation to Completed state error", err.Error())
		//dBtx.Rollback()
		return 0
	}
	//if err := dBtx.Commit().Error; err != nil {
	//	log.GetLogger(context.loggerName).Errorln("Fail to commit transation error", err.Error())
	//	return 0
	//}
	//log.GetLogger(context.loggerName).Errorln("Unable to set withdraw record to completed. Trans DB ID: ", transRecordId, " Error: ", err.Error())
	time.Sleep(2 * time.Second)
	return int(nDbtx.RowsAffected)

}

//Query the data and insert to withdraw_request table ----- RYAN
func DbInsertWithdrawRequest(me *WithdrawProcessor, transLog *types.Log, approvalWalletAddr common.Address, requestTransId *uint256.Int) (*asset.WithdrawTransaction, int64, error) {
	conn, err := me.context.db.GetConn()
	if err != nil {
		return nil, 0, err
	}

	var isTxSuccess bool = false
	withdrawTransaction := new(asset.WithdrawTransaction)
	tx := conn.Begin()
	defer func() {
		if !isTxSuccess {
			tx.Rollback()
			log.GetLogger(me.context.loggerName).Debugln("Tranaction rollback called for request trans Id: ", requestTransId.Uint64())
		}
	}()

	nTx := tx.Session(&gorm.Session{}).Clauses(clause.Locking{Strength: "SHARE"}).Where("approval_wallet_address = ? AND request_trans_id = ? and abs(status) >= ?", ethereum.ToLowerAddressString(approvalWalletAddr.Hex()), requestTransId.Uint64(), asset.StatusBurned).Find(withdrawTransaction)
	err = nTx.Error
	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("Fail to lock table", err.Error())
		return nil, 0, err
	}

	if withdrawTransaction.RequestTransHash == "" {
		log.GetLogger(me.context.loggerName).Errorln("Withdraw transaction not found for abs(status) >= 50. Request trans Id: ", requestTransId.Uint64())
		return nil, 0, errors.New("Withdraw transaction not found for abs(status) >= 50")
	}

	ethKeyPair, err := _ethereum.GetEthKeyPair(me.config.HdWalletPrivateKey)
	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("ETH gen key err : ", err, " request trans id: ", requestTransId.Uint64())
		return nil, 0, err
	}

	requestTransHash := common.HexToHash(withdrawTransaction.RequestTransHash)
	toAddress := common.HexToAddress(withdrawTransaction.MainnetToAddress)
	hash, err := GetWithdrawTransactionHash(requestTransHash, toAddress, withdrawTransaction.AssetName, withdrawTransaction.Amount.BigInt())

	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("Hash err : ", err, " request trans id: ", requestTransId.Uint64())
		return nil, 0, err
	}

	signData, err := crypto.Sign(hash[:], ethKeyPair.PrivateKey)
	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("Sign err : ", err, " request trans id: ", requestTransId.Uint64())
		return nil, 0, err
	}

	withdrawRequest := new(WithdrawRequest)
	withdrawRequest.InitDate()
	withdrawRequest.WithdrawId = withdrawTransaction.Id
	withdrawRequest.ServiceId = uint64(me.config.ServiceId)
	withdrawRequest.RequestTransHash = withdrawTransaction.RequestTransHash
	withdrawRequest.Data = hash.Bytes()
	withdrawRequest.Signature = signData

	rowAffect := 0

	err = tx.Transaction(func(tx1 *gorm.DB) error {
		tx1.Session(&gorm.Session{}).Exec("LOCK TABLE withdraw_requests IN ACCESS EXCLUSIVE MODE")

		nTx = tx1.Session(&gorm.Session{}).Create(withdrawRequest)
		err = nTx.Error
		if err != nil {
			log.GetLogger(me.context.loggerName).Errorln("Fail to insert data to withdraw_request", err)
			return err
		}
		data := new(WithdrawRequestData)
		tx2 := tx1.Session(&gorm.Session{}).Where("withdraw_id = ?", withdrawRequest.WithdrawId).Find(&data.withdrawRequestList)
		err = tx2.Error
		if err != nil {
			log.GetLogger(me.context.loggerName).Errorln("Fail to query in withdraw_request", err)
			return err
		}
		rowAffect = int(tx2.RowsAffected)
		return nil
	})
	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("Fail to insert data to withdraw_request", err)
		return nil, 0, err
	}

	if err == nil {
		nTx = tx.Commit()
		err = nTx.Error
	}

	if err != nil {
		log.GetLogger(me.context.loggerName).Errorln("Fail to insert data to withdraw_request", err)
		return nil, 0, err
	}

	isTxSuccess = true
	return withdrawTransaction, int64(rowAffect), err
}

func DbCountWithdrawRequests(context *WithdrawProcessorContext, withdrawId uint64) (int64, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return 0, err
	}
	data := new(WithdrawRequestData)
	tx := conn.Where("withdraw_id = ?", withdrawId).Find(&data.withdrawRequestList)
	err = tx.Error
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to query in withdraw_request", err)
		return 0, err
	}
	log.GetLogger(context.loggerName).Debugln("Count of withdraw request data is : ", int64(len(data.withdrawRequestList)))
	return tx.RowsAffected, nil

}

func DbUpdateWithdrawTransactionToConfirmingTransfer(context *WithdrawProcessorContext, withdrawTx *asset.WithdrawTransaction) (int64, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return 0, err
	}
	data := asset.WithdrawTransaction{
		Status: asset.StatusConfirmingTransfer,
	}
	tx := conn.Where("id = ? and abs(status) = ?", withdrawTx.Id, asset.StatusBurned).Updates(data)
	err = tx.Error
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to Update to database.", err, "The withdraw id is ", withdrawTx.Id)
		return 0, err
	}

	return tx.RowsAffected, nil
}

func DbQueryWithdrawTransactionToConfirmingTransfer(context *WithdrawProcessorContext) (*asset.WithdrawTransactionData, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return nil, err
	}
	withdrawTransaction := new(asset.WithdrawTransactionData)
	tx := conn.Where("status = ?", asset.StatusConfirmingTransfer).Limit(5).Find(&withdrawTransaction)
	err = tx.Error
	if err != nil {
		log.GetLogger(context.loggerName).Errorln("Fail to query the confirming transfer data ", err)
		return nil, err
	}
	log.GetLogger(context.loggerName).Debugln("Row count :", tx.RowsAffected)

	return withdrawTransaction, nil
}

func DBQuerySignedData(context *WithdrawProcessorContext, withdrawTx *asset.WithdrawTransaction) ([][]byte, [][32]byte, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return nil, nil, err
	}
	data := new(WithdrawRequestData)
	tx := conn.Select("data , signature").Where("withdraw_id = ?", withdrawTx.Id).Find(&data.withdrawRequestList)
	err = tx.Error
	if err != nil {
		return nil, nil, err
	}

	if (len(data.withdrawRequestList)) < 5 {
		err = errors.New("The signature not enough 5")
		return nil, nil, err
	}
	signData := make([][]byte, len(data.withdrawRequestList))
	for i, e := range data.withdrawRequestList {
		signData[i] = e.Signature
	}
	hashData := make([][32]byte, len(data.withdrawRequestList))
	for i, e := range data.withdrawRequestList {
		hashData[i] = common.BytesToHash(e.Data)
	}
	return signData, hashData, nil

}

func DbInsertAssetAllocationCost(context *WithdrawProcessorContext, receipt *ethereum.BesuReceipt) error {
	var err error
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}
		now := time.Now()
		assetAllocationCost := new(asset.AssetAllocationCost)
		assetAllocationCost.TransHash = strings.ToLower(receipt.TxHash.Hex())
		assetAllocationCost.GasPrice = decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0)
		assetAllocationCost.GasUsed = receipt.GasUsed
		assetAllocationCost.AllocationType = "Withdrawal"
		assetAllocationCost.CreatedDate = now
		assetAllocationCost.LastModifiedDate = now
		dbConn = dbConn.Create(assetAllocationCost)

		if dbConn.Error != nil {
			log.GetLogger(context.loggerName).Errorln("Unable to insert withdrawal asset allocation cost into DB", assetAllocationCost.TransHash, " error: ", dbConn.Error.Error())
			fmt.Println("retry ", i, " times")
			continue
		} else {
			break
		}

	}
	return err
}

func DbUpdateMainnetTransferTransHash(context *WithdrawProcessorContext, requestTransHash string, withdrawId uint64, transferTransHash common.Hash) error {
	var err error
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}
		withdrawData := &asset.WithdrawTransaction{
			MainnetTransHash: strings.ToLower(transferTransHash.String()),
		}

		dbTx := dbConn.Where("id = ? AND status = ? ", withdrawId, asset.StatusTransferProcessing).Updates(withdrawData)
		err = dbTx.Error
		if dbTx.Error != nil {
			log.GetLogger(context.loggerName).Errorln("Unable to update mainnet transfer tran hash to DB. Request trans hash: ", requestTransHash, " transfer trans hash: ", transferTransHash.String(), " Error:", err)
			fmt.Println("retry ", i, " times")
			continue
		} else {
			break
		}
	}
	return err
}
func DbQueryBurnedTransaction(context *WithdrawProcessorContext, requestTransHash string) (*asset.WithdrawTransaction, error) {
	conn, err := context.db.GetConn()
	if err != nil {
		return nil, err
	}
	withdrawTx := new(asset.WithdrawTransaction)

	tx := conn.Where("request_trans_hash = ? AND abs(status) IN ? ", strings.ToLower(requestTransHash), []asset.WithdrawStatus{asset.StatusConfirmingTransfer, asset.StatusTransferProcessing}).Find(&withdrawTx)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return withdrawTx, nil
}

func DbUpdateErrorTransactionToConfirmingTransfer(context *WithdrawProcessorContext, requestTransHash string) error {
	conn, err := context.db.GetConn()
	if err != nil {
		return err
	}

	tx := conn.Model(&asset.WithdrawTransaction{}).Where("request_trans_hash = ? AND abs(status) in ?",
		strings.ToLower(requestTransHash),
		[]asset.WithdrawStatus{asset.StatusConfirmingTransfer, asset.StatusTransferProcessing}).Update("status", asset.StatusConfirmingTransfer)

	return tx.Error
}

func GetWithdrawTransactionHash(requestTransHash common.Hash, toAddr common.Address, assetName string, amount *big.Int) (common.Hash, error) {
	bytesType, _ := abi.NewType("bytes32", "bytes32", nil)
	addrType, _ := abi.NewType("address", "address", nil)
	strType, _ := abi.NewType("string", "string", nil)
	intType, _ := abi.NewType("uint256", "uint256", nil)

	arguments := abi.Arguments{
		abi.Argument{Name: "requestTransHash", Type: bytesType},
		abi.Argument{Name: "toAddr", Type: addrType},
		abi.Argument{Name: "assetName", Type: strType},
		abi.Argument{Name: "amount", Type: intType},
	}

	packedByte, err := arguments.Pack(requestTransHash, toAddr, assetName, amount)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(crypto.Keccak256(packedByte)), nil
}
