package reward

import (
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RewardDBProcessor struct {
	db      *database.Database
	slaveDb *database.ReadOnlyDatabase
}

func NewRewardDBProcessor(db *database.Database, slaveDb *database.ReadOnlyDatabase) *RewardDBProcessor {
	processor := new(RewardDBProcessor)
	processor.db = db
	processor.slaveDb = slaveDb
	return processor
}

func (me *RewardDBProcessor) DbInsertDistributedToken(distributedToken *DistributedToken) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return err
	}
	if distributedToken.Status != DistributedStatusError {
		tx := dbConn.Create(&distributedToken)
		if tx.Error != nil {
			return tx.Error
		}
	} else {

		err = dbConn.Transaction(func(dbTx *gorm.DB) error {

			tx := dbConn.Create(&distributedToken)
			if tx.Error != nil {
				return tx.Error
			}

			tx = dbTx.Session(&gorm.Session{}).Create(&distributedToken.DistributedTokenError)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Unable to insert distributed_token_error record")
			}

			tx = dbTx.Session(&gorm.Session{}).Delete(&DistributedToken{}, "id = ?", distributedToken.Id)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Unable to delete distributed_tokens record")
			}

			return nil
		})

		return err

	}
	return nil
}

func (me *RewardDBProcessor) DbUpdateDistributedToken(id uint64, txHash string, gasPrice *decimal.Decimal, gasUsed uint64, gasFee *decimal.Decimal, status TokenDistributedStatus) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return err
	}
	var updateMap map[string]interface{} = make(map[string]interface{})
	if status != DistributedStatusError {
		updateMap["tx_hash"] = txHash
		updateMap["gas_price"] = gasPrice
		updateMap["gas_used"] = gasUsed
		updateMap["gas_fee"] = gasFee
		updateMap["status"] = status
		tx := dbConn.Model(DistributedToken{}).Where("id = ?", id).Updates(updateMap)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		err = dbConn.Transaction(func(dbTx *gorm.DB) error {
			model := new(DistributedToken)
			tx := dbTx.Session(&gorm.Session{}).Model(&DistributedToken{}).Where("id = ?", id).First(&model)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Unable to query distributed_tokens record")
			}
			model.TxHash = txHash
			model.GasPrice = gasPrice
			model.GasUsed = gasUsed
			model.GasFee = gasFee
			model.LastModifiedDate = time.Now()

			tx = dbTx.Session(&gorm.Session{}).Create(&model.DistributedTokenError)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Unable to insert distributed_token_error record")
			}

			tx = dbTx.Session(&gorm.Session{}).Delete(&DistributedToken{}, "id = ?", id)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Unable to delete distributed_tokens record")
			}

			return nil
		})

		return err
	}

	return nil
}

func (me *RewardDBProcessor) DbGetDistributedToken(userId uint64, assetName string, distributedType TokenDistributedType) (*DistributedToken, error) {
	dbConn, err := me.slaveDb.GetConn()
	if err != nil {
		return nil, err
	}
	distributedToken := new(DistributedToken)

	var tx *gorm.DB

	tx = dbConn.Where("user_id = ? AND asset_name = ? AND status >= ?", userId, assetName, DistributedStatusPending)
	if distributedType != DistributedUnknown {
		tx = tx.Where("distributed_type = ?", distributedType)
	}
	tx = tx.Order("created_date DESC").FirstOrInit(&distributedToken)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return distributedToken, nil
}

func (me *RewardDBProcessor) DbGetUserIdByPaymentWalletAddress(paymentWalletAddress string) (uint64, error) {
	dbConn, err := me.slaveDb.GetConn()
	if err != nil {
		return 0, err
	}
	var userId uint64 = 0
	tx := dbConn.Table("users").Select("id").Where("payment_wallet_address = ?", ethereum.ToLowerAddressString(paymentWalletAddress)).Find(&userId)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return userId, nil
}
