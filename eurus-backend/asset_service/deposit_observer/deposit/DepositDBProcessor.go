package deposit

import (
	"bytes"
	"fmt"

	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func DbInsertPendingDeposit(context *DepositProcessorContext, tx *AssetTransferTransaction, receipt *ethereum.BesuReceipt, transferID common.Hash, isERC20 bool) (*asset.DepositTransaction, int, error) {
	var dbTrans *asset.DepositTransaction = new(asset.DepositTransaction)
	var depositContext = context

	if isERC20 {
		if bytes.Equal(tx.TransferLog[0].Topics[0].Bytes(), transferID.Bytes()) {
			if len(tx.TransferLog[0].Topics) < 3 {
				log.GetLogger(context.LoggerName).Errorln("Invalid transfer event. Expected 3 topics for trans ", tx.Hash().Hex())
				return nil, 0, errors.New("Invalid transfer event. Expected 3 topics for trans " + tx.Hash().Hex())
			}

			dbTrans.MainnetTransHash = strings.ToLower(receipt.TxHash.Hex())
			senderAddr, _ := tx.GetSender()
			dbTrans.MainnetFromAddress = strings.ToLower(senderAddr)
			dbTrans.MainnetToAddress = ethereum.ToLowerAddressString(common.BytesToAddress(tx.TransferLog[0].Topics[2].Bytes()).Hex())
			dbTrans.AssetName = tx.AssetName

			gasUsed := decimal.NewFromBigInt(new(big.Int).SetUint64(receipt.GasUsed), 0)
			dbTrans.MainnetGasUsed = gasUsed
			dbTrans.MainnetGasFee = gasUsed.Mul(decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0))

			transDate := time.Unix(int64(tx.Block.Time()), 0)
			dbTrans.MainnetTransDate = &transDate

			dbTrans.CustomerType = asset.CustomerUser

			userObj, err := dbGetUserByWalletAddress(depositContext, tx.Receiptant)
			if err != nil {
				log.GetLogger(context.LoggerName).Errorln("Unable to query user by To address: ", err.Error())
				return nil, 0, err
			}
			dbTrans.InnetToAddress = ethereum.ToLowerAddressString(tx.Receiptant)
			dbTrans.CustomerId = userObj.Id
			dbTrans.Status = asset.DepositReceiptCollected

			dbTrans.CustomerType = asset.CustomerUser

			argumentList, err := ethereum.DefaultABIDecoder.DecodeABIEventData(tx.TransferLog[0].Data, "ERC20", "Transfer")
			if err != nil {
				log.GetLogger(context.LoggerName).Errorln("Unable to decode event data for trans: ", tx.Hash().Hex(), " error: ", err.Error())
				return nil, 0, err
			}
			if len(argumentList) < 1 {
				log.GetLogger(context.LoggerName).Errorln("Invalid event data for trans: ", tx.Hash().Hex(), " error: ", err.Error())
				return nil, 0, err
			}

			var amount *big.Int
			var ok bool
			amount, ok = argumentList[0].(*big.Int)
			if !ok {
				log.GetLogger(context.LoggerName).Errorln("Invalid amount data type for trans: ", tx.Hash().Hex())
				return nil, 0, errors.New("Invalid amount data type for trans: " + tx.Hash().Hex())
			}
			dbTrans.Amount = decimal.NewFromBigInt(amount, 0)

		}

	} else {
		dbTrans.MainnetTransHash = strings.ToLower(receipt.TxHash.Hex())
		sender, err := tx.GetSender()
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Sender not found. ", err)
			return nil, 0, err
		}
		dbTrans.MainnetFromAddress = strings.ToLower(sender)
		dbTrans.MainnetToAddress = strings.ToLower(tx.GetTo().String())
		dbTrans.AssetName = "ETH"
		gasUsed := decimal.NewFromBigInt(new(big.Int).SetUint64(receipt.GasUsed), 0)
		dbTrans.MainnetGasUsed = gasUsed
		dbTrans.MainnetGasFee = gasUsed.Mul(decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0))
		transDate := time.Unix(int64(tx.Block.Time()), 0)
		dbTrans.MainnetTransDate = &transDate
		dbTrans.CustomerType = asset.CustomerUser
		dbTrans.InnetToAddress = ethereum.ToLowerAddressString(tx.Receiptant)

		userObj, err := dbGetUserByWalletAddress(context, dbTrans.InnetToAddress)
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query user by To address: ", err.Error())
			return nil, 0, err
		}
		dbTrans.CustomerId = userObj.Id
		dbTrans.Status = asset.DepositReceiptCollected
		dbTrans.Amount = decimal.NewFromBigInt(tx.OriginalTransaction.Value(), 0)

	}

	_, rowAffected, err := dbInsertPendingDepositTransaction(context, dbTrans)
	if err != nil {
		return nil, 0, err
	}

	return dbTrans, rowAffected, nil

}

func DbInsertSweepTransaction(context *DepositProcessorContext, tx *AssetTransferTransaction) {
	trans := new(asset.DepositTransaction)
	trans.InitDate()
	trans.MainnetTransHash = strings.ToLower(tx.Hash().String())
	trans.AssetName = tx.AssetName
	trans.Amount = decimal.NewFromBigInt(tx.OriginalTransaction.Value(), 0)
	trans.Status = asset.SweepTrans
	transDate := time.Unix(int64(tx.Block.Time()), 0)
	trans.MainnetTransDate = &transDate
	sender, _ := tx.GetSender()
	trans.MainnetFromAddress = ethereum.ToLowerAddressString(sender)
	trans.MainnetToAddress = ethereum.ToLowerAddressString(tx.OriginalTransaction.To().Hex())

	_, _, _ = dbInsertPendingDepositTransaction(context, trans)
}

func dbInsertPendingDepositTransaction(context *DepositProcessorContext, dbTrans *asset.DepositTransaction) (*asset.DepositTransaction, int, error) {
	var err error
	var rowAffected int
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		conn, err := context.db.GetConn()
		if err != nil {
			time.Sleep(time.Duration(context.retrySetting.GetRetryInterval()) * time.Second)
			continue
		}
		dbTrans.InitDate()
		dbConn := conn.Create(dbTrans)
		if dbConn.Error != nil {
			if strings.Contains(dbConn.Error.Error(), "duplicate key") {
				return nil, 0, nil
			}
			log.GetLogger(context.LoggerName).Errorln("Unable to insert into DB for trans: ", dbTrans.MainnetTransHash, " error: ", dbConn.Error.Error())
		} else {
			rowAffected = int(dbConn.RowsAffected)
			break
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return dbTrans, rowAffected, nil
}

func dbGetUserByWalletAddress(context *DepositProcessorContext, address string) (*user.User, error) {
	var err error
	var userObj *user.User

	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.slaveDb.GetConn()
		if err1 != nil {
			err = err1
			continue
		}

		userObj = new(user.User)
		tx := dbConn.Where("wallet_address = ?", ethereum.ToLowerAddressString(address)).FirstOrInit(userObj)
		err = tx.Error
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		userObj = nil
	}
	return userObj, err

}

func dbUpdateDepositTransToCollected(context *DepositProcessorContext,
	transDbId uint64, collectTransHash common.Hash, collectTransDate *time.Time) error {

	var err error

	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}
		dbConn = dbConn.Model(asset.DepositTransaction{}).Where("id = ? AND ABS(status) < ?",
			transDbId, asset.DepositAssetCollected).Updates(
			asset.DepositTransaction{
				MainnetCollectTransHash: strings.ToLower(collectTransHash.Hex()),
				MainnetCollectTransDate: collectTransDate,
				Status:                  asset.DepositAssetCollected,
			})

		err = dbConn.Error
		if err != nil {
			log.GetLogger(context.LoggerName).Errorf("Unable to update deposit transaction status for collect trans hash: %s, Error: %s\r\n", collectTransHash.Hex(), err.Error())
		} else {
			break
		}
	}
	return err
}

func dbUpdateDepositTransToMintRequesting(context *DepositProcessorContext,
	depositTransHash common.Hash) (int, error) {
	var err error
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}

		var updateMap map[string]interface{} = map[string]interface{}{
			"status":             asset.DepositMintRequesting,
			"last_modified_date": time.Now(),
		}

		dbConn = dbConn.Model(asset.DepositTransaction{}).Where("mainnet_trans_hash = ? AND ABS(status) <= ? ",
			strings.ToLower(depositTransHash.Hex()), asset.DepositMintRequesting).Updates(updateMap)

		err = dbConn.Error
		if err == nil {
			return int(dbConn.RowsAffected), nil
		}
	}

	log.GetLogger(context.LoggerName).Errorf("dbUpdateDepositTransToMintRequesting - Unable to update deposit transaction status for deposit trans hash: %s, Error: %s\r\n", depositTransHash.Hex(), err.Error())
	return 0, err
}

// func dbUpdateDepositTransToMintConfirming(context *DepositProcessorContext,
// 	depositTransHash common.Hash, mintTransId uint256.Int) (int, error) {
// 	var err error
// 	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
// 		dbConn, err1 := context.db.GetConn()
// 		if err1 != nil {
// 			err = err1
// 			continue
// 		}

// 		var mintTransBigId *uint64 = new(uint64)
// 		*mintTransBigId = mintTransId.ToBig().Uint64()

// 		var updateMap map[string]interface{} = map[string]interface{}{
// 			"status":             asset.DepositMintConfirming,
// 			"last_modified_date": time.Now(),
// 			"mint_trans_id":      mintTransBigId,
// 		}

// 		dbConn = dbConn.Model(asset.DepositTransaction{}).Where("mainnet_trans_hash = ? AND ABS(status) <= ? ",
// 			strings.ToLower(depositTransHash.Hex()), asset.DepositMintConfirming).Updates(updateMap)

// 		err = dbConn.Error
// 		if err == nil {
// 			return int(dbConn.RowsAffected), nil
// 		}
// 	}
// 	log.GetLogger(context.LoggerName).Errorf("dbUpdateDepositTransToMintConfirming - Unable to update deposit transaction status for deposit trans hash: %s, Error: %s\r\n", depositTransHash.Hex(), err.Error())
// 	return 0, err
// }

func dbUpdateDepositTransToError(context *DepositProcessorContext, depositTransHash common.Hash,
	fromState asset.DepositStatus, remarks string) error {
	var err error
	now := time.Now()

	var updateMap map[string]interface{} = map[string]interface{}{
		"status":             fromState * asset.DepositError,
		"remarks":            remarks,
		"last_modified_date": now,
	}

	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}

		dbConn = dbConn.Model(asset.DepositTransaction{}).Where("mainnet_trans_hash = ? AND ABS(status) <=  ?",
			strings.ToLower(depositTransHash.Hex()), fromState).Updates(updateMap)

		err = dbConn.Error
		if err != nil {
			log.GetLogger(context.LoggerName).Errorf("dbUpdateDepositTransToError - Unable to update deposit transaction status for deposit trans hash: %s, Error: %s\r\n", depositTransHash.Hex(), err.Error())
		} else {
			return nil
		}
	}
	return err
}

func dbUpdateDepositTransToCompleted(context *DepositProcessorContext, depositTransHash common.Hash,
	mintTransHash common.Hash, mintDate *time.Time, innetFromAddress common.Address) (int, error) {
	var err error
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}
		now := time.Now()
		dbConn = dbConn.Model(asset.DepositTransaction{}).Where("mainnet_trans_hash = ? AND ABS(status) < ?",
			strings.ToLower(depositTransHash.Hex()), asset.DepositCompleted).Updates(
			map[string]interface{}{
				"status":             asset.DepositCompleted,
				"mint_trans_hash":    strings.ToLower(mintTransHash.Hex()),
				"mint_date":          mintDate,
				"innet_from_address": strings.ToLower(innetFromAddress.Hex()),
				"remarks":            "",
				"last_modified_date": now,
			})

		err = dbConn.Error
		if err != nil {
			log.GetLogger(context.LoggerName).Errorf("dbUpdateDepositTransToCompleted - Unable to update deposit transaction status for deposit trans hash: %s, Error: %s\r\n", depositTransHash.Hex(), err.Error())
		} else {
			return int(dbConn.RowsAffected), nil
		}
	}
	return 0, err
}

func DbInsertAssetAllocationCost(context *DepositProcessorContext, receipt *ethereum.BesuReceipt, gasPrice *big.Int) error {
	var err error
	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		dbConn, err1 := context.db.GetConn()
		if err1 != nil {
			err = err1
			continue
		}
		now := time.Now()
		assetAllocationCost := new(asset.AssetAllocationCost)
		assetAllocationCost.TransHash = receipt.TxHash.Hex()
		assetAllocationCost.GasPrice = decimal.NewFromBigInt(gasPrice, 0)
		assetAllocationCost.GasUsed = receipt.GasUsed
		assetAllocationCost.AllocationType = "Deposit"
		assetAllocationCost.CreatedDate = now
		assetAllocationCost.LastModifiedDate = now
		dbConn = dbConn.Create(assetAllocationCost)

		if dbConn.Error != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to insert asset allocation cost into DB for trans: ", assetAllocationCost.TransHash, " error: ", dbConn.Error.Error())
			fmt.Println("retry ", i, " times")
			continue
		} else {
			break
		}

	}
	return err
}

func dbGetUserSideChainWalletAddress(mainNetAddress common.Address, context *DepositProcessorContext) (*common.Address, error) {
	dbConn, err := context.slaveDb.GetConn()
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to connect DB. ", err)
		return nil, err
	}
	userTable := new(user.User)

	tx := dbConn.Where("mainnet_wallet_address = ? and status = ? and is_metamask_addr = ?", strings.ToLower(mainNetAddress.String()), user.UserStatusNormal, false).Find(&userTable)
	err = tx.Error
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to Find Address: ", mainNetAddress.String(), " Error:", err)
		return nil, err
	}
	if userTable.Id == 0 {
		return nil, nil
	}

	addr := common.HexToAddress(userTable.WalletAddress)
	return &addr, nil
}

func dbGetUserDepositTransactionDetails(context *DepositProcessorContext, depositTransHash common.Hash) (*user.User, *asset.DepositTransaction, error) {
	dbConn, err := context.slaveDb.GetConn()
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to connect DB. ", err)
		return nil, nil, err
	}

	dt := new(asset.DepositTransaction)

	err = dbConn.Where("mainnet_trans_hash = ?", strings.ToLower(depositTransHash.Hex())).First(dt).Error
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to get deposit transaction with given mainnet_trans_hash:", depositTransHash, err)
		return nil, nil, err
	}

	u := new(user.User)

	err = dbConn.Where("id = ?", dt.CustomerId).First(u).Error
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to get user info with given customer_id:", dt.CustomerId, err)
		return nil, nil, err
	}

	return u, dt, nil
}

func dbInsertPendingSweepWallet(context *DepositProcessorContext, userID *uint64, mainnetWalletAddress string, assetName string) error {
	dbConn, err := context.db.GetConn()
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Fail to connect DB. ", err)
		return err
	}

	wallet := new(asset.PendingSweepWallet)
	wallet.InitDate()
	wallet.UserID = userID
	wallet.MainnetWalletAddress = strings.ToLower(mainnetWalletAddress)
	wallet.AssetName = assetName

	err = dbConn.Create(wallet).Error
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			log.GetLogger(context.LoggerName).Warnf("Record with same key exists: {%v, %v}\n", wallet.MainnetWalletAddress, wallet.AssetName)
			return nil
		}

		return err
	}

	return nil
}
