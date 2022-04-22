package wallet_bg

import (
	"eurus-backend/foundation/database"
	wallet_bg_model "eurus-backend/report_service/wallet_background_indexer/wallet_bg/model"
	"eurus-backend/user_service/user_service/user"

	"gorm.io/gorm/clause"
)

// All DB related function put in here

func DbQueryWalletBalanceConfig(db *database.Database) ([]wallet_bg_model.WalletBalanceConfig, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var configList []wallet_bg_model.WalletBalanceConfig = make([]wallet_bg_model.WalletBalanceConfig, 0)
	tx := dbConn.Model(wallet_bg_model.WalletBalanceConfig{}).Scan(&configList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return configList, nil
}

func DbBatchInsertWalletBalance(db *database.Database, newRecords *[]wallet_bg_model.WalletBalance) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	for _, record := range *newRecords {
		tx := dbConn.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "mark_date"}, {Name: "wallet_type"}, {Name: "wallet_address"}, {Name: "asset_name"}, {Name: "chain_id"}}, // key colume
			DoUpdates: clause.AssignmentColumns([]string{"balance", "created_date"}),                                                                   // column needed to be updated
		}).Create(&record)
		err = tx.Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DbGetUserWalletAddressesFromDb(db *database.Database) (*[]user.User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var list []user.User
	statusList := []user.UserStatus{user.UserStatusNormal, user.UserStatusSuspended}

	tx := dbConn.Model(&user.User{}).Select("id, wallet_address,mainnet_wallet_address").Where("status IN (?)", statusList).Find(&list)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return &list, err
}

func DbInsertAssetTotalSupply(db *database.Database, assetTotalSupply *wallet_bg_model.AssetTotalSupply) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	tx := dbConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "asset_name"}, {Name: "chain_id"}, {Name: "mark_date"}},           // key column
		DoUpdates: clause.AssignmentColumns([]string{"total_supply", "last_modified_date", "block_number"}), // column needed to be updated
	}).Create(&assetTotalSupply)
	err = tx.Error
	if err != nil {
		return err
	}
	return nil
}

func DbGetCentralizedUserWalletAddressesFromDb(db *database.Database) (*[]user.User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var list []user.User
	statusList := []user.UserStatus{user.UserStatusNormal, user.UserStatusSuspended}

	tx := dbConn.Model(&user.User{}).Select("id, wallet_address,mainnet_wallet_address").Where("is_metamask_addr = ? AND status IN (?)", false, statusList).Find(&list)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return &list, err
}
