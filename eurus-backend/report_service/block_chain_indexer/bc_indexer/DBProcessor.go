package bc_indexer

import (
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func DbGetUserByWalletAddress(context *blockChainProcessorContext, address string) (*user.User, error) {
	var returnErr error
	var userObj *user.User
	for i := 0; i < context.Config.GetRetryCount(); i++ {

		dbConn, err := context.SlaveDb.GetConn()
		if err != nil {
			returnErr = err
			log.GetLogger(context.LoggerName).Errorln("DB connection error when query user by wallet address: ", address, " Error: ", err)
			time.Sleep(context.Config.GetRetryInterval() * time.Second)
			continue
		}
		userObj = new(user.User)
		tx := dbConn.Where("wallet_address = ?", ethereum.ToLowerAddressString(address)).Find(userObj)
		if tx.Error != nil {
			returnErr = tx.Error
			log.GetLogger(context.LoggerName).Errorln("Error query user by wallet address: ", address, " Error: ", err)
			time.Sleep(context.Config.GetRetryInterval() * time.Second)
			continue
		}
		break
	}
	if returnErr != nil {
		return nil, returnErr
	}
	return userObj, nil
}

func DbGetUserByMainnetWalletAddress(context *blockChainProcessorContext, address string) (*user.User, error) {
	var returnErr error
	var userObj *user.User
	for i := 0; i < context.Config.GetRetryCount(); i++ {

		dbConn, err := context.SlaveDb.GetConn()
		if err != nil {
			returnErr = err
			log.GetLogger(context.LoggerName).Errorln("DB connection error when query user by mainnet wallet address: ", address, " Error: ", err)
			time.Sleep(context.Config.GetRetryInterval() * time.Second)
			continue
		}
		userObj = new(user.User)
		tx := dbConn.Where("mainnet_wallet_address = ?", address).Find(userObj)
		if tx.Error != nil {
			returnErr = tx.Error
			log.GetLogger(context.LoggerName).Errorln("Error query user by mainnet wallet address: ", address, " Error: ", err)
			time.Sleep(context.Config.GetRetryInterval() * time.Second)
			continue
		}
		break
	}
	if returnErr != nil {
		return nil, returnErr
	}
	return userObj, nil
}

func AddTxIndexToDB(db *database.Database, ext *ExtractedTransaction, chainId int, isFrom bool) error {
	dbConn, err := db.GetConn()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Tx Hash: "+strings.ToLower(ext.TxHash)+" Unable to connect db", err.Error())
		return errors.New("Database Network Error: " + err.Error())
	}

	var userWalletAddr string
	var userId uint64
	if isFrom {
		if ext.FromUser == nil {
			return foundation.NewError(foundation.InvalidArgument)
		}

		if ext.IsMainnetTrans {
			userWalletAddr = strings.ToLower(ext.FromUser.MainnetWalletAddress)
			userId = ext.FromUser.Id
		} else {
			userWalletAddr = strings.ToLower(ext.FromUser.WalletAddress)
			userId = ext.FromUser.Id
		}

	} else {

		if ext.ToUser != nil {
			if ext.IsMainnetTrans {
				userWalletAddr = strings.ToLower(ext.ToUser.MainnetWalletAddress)
				userId = ext.ToUser.Id
			} else {
				userWalletAddr = strings.ToLower(ext.ToUser.WalletAddress)
				userId = ext.ToUser.Id
			}
		}
	}

	if ext.ConfirmTransHash != "" && ext.RequestTransId != nil {

		if ext.FromUser == nil {
			return foundation.NewErrorWithMessage(foundation.InvalidArgument, "Missing from user object")
		}
		err = dbConn.Transaction(func(dbTx *gorm.DB) error {
			fromWalletAddr := strings.ToLower(ext.FromUser.WalletAddress)
			transferTrans := new(asset.TransferTransaction)

			updateTransferTrans := new(asset.TransferTransaction)
			updateTransferTrans.InitDate()
			updateTransferTrans.Status = asset.TransferStatus(ext.Status)
			updateTransferTrans.ConfirmTransHash = ext.ConfirmTransHash

			ctxDB := dbTx.Session(&gorm.Session{})
			ctxDB = ctxDB.Select("status", "confirm_trans_hash", "last_modified_date").Where("from_address = ? AND request_trans_id = ? AND chain = ? AND status < ? AND is_send = ?",
				fromWalletAddr, ext.RequestTransId.Uint64(), chainId, asset.TransferStatusConfirmed, isFrom).Updates(updateTransferTrans)

			if ctxDB.Error != nil {
				return ctxDB.Error
			}

			if ctxDB.RowsAffected > 0 && ((asset.TransferStatus(ext.Status) == asset.TransferStatusConfirmed && !isFrom) || (isFrom)) {
				ctxDB = dbTx.Session(&gorm.Session{})
				ctxDB = ctxDB.Select("*").Where("from_address = ? AND request_trans_id = ? AND chain = ? AND is_send = ?",
					fromWalletAddr, ext.RequestTransId.Uint64(), chainId, isFrom).FirstOrInit(transferTrans)
				if ctxDB.Error != nil {
					return ctxDB.Error
				}
				var addr string
				if isFrom {
					addr = transferTrans.FromAddress
				} else {
					addr = transferTrans.ToAddress
				}
				if transferTrans.UserId != 0 {
					txIndex := TransactionIndex{
						TxHash:        transferTrans.TxHash,
						WalletAddress: addr,
						UserId:        transferTrans.UserId,
						CreatedDate:   transferTrans.CreatedDate,
						AssetName:     transferTrans.AssetName,
						Status:        transferTrans.Status > 0,
					}
					ctxDB = dbTx.Session(&gorm.Session{})
					ctxDB = ctxDB.Create(txIndex)
					if ctxDB.Error != nil {
						return err
					}
				}
			}
			return nil
		})
		return err

	}

	txIndex := TransactionIndex{
		TxHash:        strings.ToLower(ext.TxHash),
		WalletAddress: userWalletAddr, //This wallet address is referred to user ID wallet address
		UserId:        userId,
		CreatedDate:   ext.CreatedDate,
		AssetName:     ext.AssetName,
		Status:        ext.Status > 0,
	}

	var txAmount decimal.Decimal
	if ext.AssetName != asset.EurusTokenName || ext.RequestTransId != nil {
		txAmount = decimal.NewFromBigInt(ext.Amount, 0)
	} else {
		txAmount = decimal.NewFromBigInt(ext.OriginalTransaction.Value(), 0)
	}

	var walletAddr string
	if ext.IsMainnetTrans {
		if ext.FromUser != nil {
			walletAddr = strings.ToLower(ext.FromUser.MainnetWalletAddress)
		} else {
			walletAddr, _ = ext.GetSender()
		}
	} else {
		if ext.FromUser != nil {
			walletAddr = strings.ToLower(ext.FromUser.WalletAddress)
		} else {
			walletAddr, _ = ext.GetSender()
		}
	}

	var requestTransId *decimal.Decimal
	if ext.RequestTransId != nil {
		temp := decimal.NewFromBigInt(ext.RequestTransId, 0)
		requestTransId = &temp
	}

	var requestTransIdVal *uint64
	if requestTransId != nil {
		var val uint64 = requestTransId.BigInt().Uint64()
		requestTransIdVal = &val
	}

	var finalGasPrice decimal.Decimal
	if ext.IsMainnetTrans {
		finalGasPrice = decimal.NewFromBigInt(ext.EffectiveGasPrice, 0)
	} else {
		finalGasPrice = decimal.NewFromBigInt(ext.OriginalTransaction.GasPrice(), 0)
	}

	gasFee := big.NewInt(0)
	gasFee = gasFee.Mul(new(big.Int).SetUint64(ext.UserGasUsed), finalGasPrice.BigInt())

	transferTx := asset.TransferTransaction{
		UserId:          userId,
		AssetName:       ext.AssetName,
		FromAddress:     walletAddr, //this wallet address refers to sender wallet address
		ToAddress:       strings.ToLower(ext.GetTo()),
		IsSend:          isFrom,
		RequestTransId:  requestTransIdVal,
		Status:          asset.TransferStatus(ext.Status),
		TxHash:          strings.ToLower(ext.TxHash),
		Chain:           int(ext.OriginalTransaction.ChainId().Int64()),
		Amount:          txAmount,
		GasFee:          decimal.NewFromBigInt(gasFee, 0),
		GasPrice:        finalGasPrice,
		TransGasUsed:    &ext.TransGasUsed,
		UserGasUsed:     &ext.UserGasUsed,
		TransactionDate: ext.CreatedDate,
		Remarks:         ext.Remarks,
	}
	transferTx.InitDate()

	err = dbConn.Transaction(func(dbTx *gorm.DB) error {

		ctxDB := dbTx.Session(&gorm.Session{})

		if asset.TransferStatus(ext.Status) != asset.TransferStatusPending {
			ctxDB = ctxDB.Create(&txIndex)
			err = ctxDB.Error
			if err != nil {
				return err
			}
		}

		ctxDB = dbTx.Session(&gorm.Session{})
		ctxDB = ctxDB.Create(&transferTx)
		err = ctxDB.Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		log.GetLogger(log.Name.Root).Error("Tx Hash: "+strings.ToLower(ext.TxHash)+" Unable to add transaction index to db", err.Error())
		return nil
	} else if err != nil {
		log.GetLogger(log.Name.Root).Error("Tx Hash: "+strings.ToLower(ext.TxHash)+" Unable to add transaction index to db", err.Error())
		return err
	}

	return nil
}

func DBGetUserList(db *database.ReadOnlyDatabase) (*user.UserList, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	user := new(user.UserList)
	tx := dbConn.Find(&user.User)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func DBGetTransferTransaction(db *database.ReadOnlyDatabase, fromWalletAddr string, requestTransId *big.Int, chainId int) (*asset.TransferTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	transferTrans := new(asset.TransferTransaction)
	dbTx := dbConn.Where("wallet_address = ? AND request_trans_id = ? AND chain = ?", strings.ToLower(fromWalletAddr), requestTransId, chainId).FirstOrInit(transferTrans)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return transferTrans, nil
}

func DBInsertPurchaseTransaction(db *database.Database, ext *ExtractedTransaction) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	model := new(asset.PurchaseTransaction)
	model.Amount = decimal.NewFromBigInt(ext.Amount, 0)
	model.AssetName = ext.AssetName
	model.FromAddress = ext.FromUser.WalletAddress
	model.ToAddress = ext.GetTo()
	model.CreatedDate = ext.CreatedDate
	model.LastModifiedDate = time.Now()
	model.Chain = int(ext.OriginalTransaction.ChainId().Int64())
	gasUsed := big.NewInt(0)
	gasUsed.SetUint64(ext.UserGasUsed)
	model.GasFee = decimal.NewFromBigInt(big.NewInt(0).Mul(ext.EffectiveGasPrice, gasUsed), 0)
	model.GasPrice = decimal.NewFromBigInt(ext.EffectiveGasPrice, 0)
	if ext.ProductId != nil {
		productId := ext.ProductId.Uint64()
		model.ProductId = &productId
	}
	if ext.Quantity != nil {
		quantity := decimal.NewFromBigInt(ext.Quantity, 0)
		model.Quantity = &quantity
	}
	model.TransGasUsed = ext.TransGasUsed
	model.UserGasUsed = ext.UserGasUsed
	model.Remarks = ext.Remarks
	model.TxHash = ext.TxHash
	model.UserId = ext.FromUser.Id
	model.PurchaseType = ext.TransactionType
	model.Status = asset.PurchaseStatus(ext.Status)

	dbTx := dbConn.Create(&model)

	return dbTx.Error
}

func DbInsertTopUpTransaction(db *database.Database, ext *ExtractedTransaction) error {

	topUpExt, ok := ext.ChildObject.(*TopUpExtractedTransaction)
	if !ok {
		return errors.New("Cannot convert to TopUpExtractedTransaction")
	}
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	topUpTrans := new(asset.TopUpTransaction)
	topUpTrans.InitDate()
	topUpTrans.TransferGas = decimal.NewFromBigInt(ext.Amount, 0)
	topUpTrans.TargetGas = decimal.NewFromBigInt(topUpExt.TargetGas, 0)
	topUpTrans.CustomerId = topUpExt.FromUser.Id
	topUpTrans.CustomerType = asset.CustomerUser
	topUpTrans.ToAddress = ethereum.ToLowerAddressString(topUpExt.GetTo())
	topUpTrans.FromAddress = topUpExt.FromUser.WalletAddress
	topUpTrans.GasPrice = decimal.NewFromBigInt(topUpExt.EffectiveGasPrice, 0)
	topUpTrans.TransactionDate = topUpExt.CreatedDate
	topUpTrans.TxHash = topUpExt.TxHash
	topUpTrans.Status = asset.TopUpStatus(topUpExt.Status)
	topUpTrans.TransGasUsed = topUpExt.TransGasUsed
	topUpTrans.UserGasUsed = topUpExt.UserGasUsed
	topUpTrans.IsDirectTopUp = topUpExt.IsDirectTopUp
	topUpTrans.Remarks = topUpExt.Remarks
	tx := dbConn.Create(&topUpTrans)

	return tx.Error
}
