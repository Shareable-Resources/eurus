package sweep

import (
	"eurus-backend/asset_service/asset"
	"fmt"
	"math/big"
	"strings"
	"time"

	"eurus-backend/user_service/user_service/user"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func DBGetPendingSweepWallets(context *SweepServiceContext) ([]asset.PendingSweepWallet, error) {
	dbConn, err := context.db.GetConn()
	if err != nil {
		return nil, err
	}

	var wallets []asset.PendingSweepWallet

	err = dbConn.Order("last_modified_date").Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func DBInsertPendingSweepWallet(context *SweepServiceContext, userID *uint64, mainnetWalletAddress string, currency string) error {
	dbConn, err := context.db.GetConn()
	if err != nil {
		return err
	}

	wallet := new(asset.PendingSweepWallet)
	wallet.InitDate()
	wallet.UserID = userID
	wallet.MainnetWalletAddress = strings.ToLower(mainnetWalletAddress)
	wallet.AssetName = currency

	err = dbConn.Create(wallet).Error
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			context.logger.Warnf("Record with same key exists: {%v, %v}\n", wallet.MainnetWalletAddress, wallet.AssetName)
			return nil
		}

		return err
	}

	return nil
}

func DBRequeuePendingSweepWallet(context *SweepServiceContext, id uint64, prevGasFeeCap *big.Int, prevGasTipCap *big.Int, prevGasLimit *uint64) error {
	dbConn, err := context.db.GetConn()
	if err != nil {
		return err
	}

	wallet := asset.PendingSweepWallet{}

	gfc := toDecimal(prevGasFeeCap)
	gtc := toDecimal(prevGasTipCap)

	err = dbConn.Model(&wallet).Where("id = ?", id).Updates(map[string]interface{}{"last_modified_date": time.Now(), "previous_gas_fee_cap": gfc, "previous_gas_tip_cap": gtc, "previous_gas_limit": prevGasLimit}).Error
	if err != nil {
		return err
	}

	return nil
}

func DBDeletePendingSweepWallet(context *SweepServiceContext, id uint64) error {
	dbConn, err := context.db.GetConn()
	if err != nil {
		return err
	}

	wallet := asset.PendingSweepWallet{}

	err = dbConn.Where("id = ?", id).Delete(&wallet).Error
	if err != nil {
		return err
	}

	return nil
}

func DBGetCentralizedUsers(context *SweepServiceContext) ([]user.User, error) {
	dbConn, err := context.db.GetConn()
	if err != nil {
		return nil, err
	}

	var users []user.User

	err = dbConn.Where("is_metamask_addr = ?", false).Order("id").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func DBInsertAssetAllocationCost(context *SweepServiceContext, transHash common.Hash, allocationType string, gasUsed uint64, gasPrice *big.Int) error {
	var err error

	for i := 0; i < context.retrySetting.GetRetryCount() || i == 0; i++ {
		var dbConn *gorm.DB
		dbConn, err = context.db.GetConn()
		if err != nil {
			continue
		}

		assetAllocationCost := new(asset.AssetAllocationCost)
		assetAllocationCost.InitDate()
		assetAllocationCost.TransHash = transHash.Hex()
		assetAllocationCost.AllocationType = allocationType
		assetAllocationCost.GasUsed = gasUsed
		assetAllocationCost.GasPrice = decimal.NewFromBigInt(gasPrice, 0)

		err = dbConn.Create(assetAllocationCost).Error
		if err == nil {
			return nil
		}

		context.logger.Errorln("Unable to insert asset allocation cost into DB for trans:", assetAllocationCost.TransHash, ", error:", err)
		fmt.Println("retry", i, "times")
	}

	return err
}

func toDecimal(bi *big.Int) *decimal.Decimal {
	if bi == nil {
		return nil
	}

	dec := decimal.NewFromBigInt(bi, 0)
	return &dec
}
