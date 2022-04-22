package user

import (
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/env"
	"eurus-backend/foundation/database"
	"eurus-backend/marketing/reward"

	"gorm.io/gorm"
)

func DBGetRecentTransactionIndex(req *QueryRecentTransactionDetailsRequest, db *database.ReadOnlyDatabase) ([]*asset.TransactionIndex, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	txIndices := []*asset.TransactionIndex{}

	tx := dbConn.Where("asset_name = ? AND created_date >now()- interval '30 day' AND user_id = ?", req.CurrencySymbol, req.UserId).Order("created_date DESC").Find(&txIndices)

	err = tx.Error
	if err != nil {
		return nil, err
	}
	return txIndices, nil
}

func DBGetRecentTransferTransaction(req *QueryRecentTransactionDetailsRequest, db *database.ReadOnlyDatabase) ([]*asset.TransferTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var transferTransList []*asset.TransferTransaction = make([]*asset.TransferTransaction, 0)

	transferSql := `SELECT * from (SELECT user_id, asset_name, from_address, to_address, tx_hash, chain, amount, is_send, gas_fee, transaction_date, created_date, user_gas_used, gas_price, status, remarks
		FROM transfer_transactions WHERE asset_name = ?  AND created_date >now()- interval '30 day' AND user_id = ? AND chain = ? `

	withdrawSql := ` UNION 
	SELECT customer_id, asset_name, mainnet_from_address, mainnet_to_address, mainnet_trans_hash, 1, amount, false, 0, mainnet_trans_date, w.created_date, user_gas_used, 0, ? , remarks
	FROM withdraw_transactions as w
	INNER JOIN users u ON u.mainnet_wallet_address = w.mainnet_to_address AND w.customer_id <> u.id
	WHERE w.asset_name = ? AND w.created_date >now()- interval '30 day' AND u.id = ?  AND w.customer_type = ? AND w.status = ?`

	sql := transferSql
	if req.ChainId != env.DefaultEurusChainId {
		sql += withdrawSql
	}
	sql += ` ) as r ORDER by created_date DESC`

	tx := dbConn.Session(&gorm.Session{}).Raw(sql, req.CurrencySymbol, req.UserId, req.ChainId,
		asset.TransferStatusConfirmed, req.CurrencySymbol, req.UserId, asset.CustomerUser, asset.StatusCompleted).Find(&transferTransList)

	err = tx.Error
	if err != nil {
		return nil, err
	}

	return transferTransList, nil
}

func DBGetRecentWithdrawTransaction(req *QueryRecentTransactionDetailsRequest, db *database.ReadOnlyDatabase) ([]*asset.WithdrawTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	withdrawTxs := []*asset.WithdrawTransaction{}
	tx := dbConn.Where("asset_name = ? AND request_date>now()- interval '30 day' AND customer_id = ? AND customer_type = 0", req.CurrencySymbol, req.UserId).Order("request_date DESC").Find(&withdrawTxs)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return withdrawTxs, nil
}

func DBGetRecentDepositTransaction(req *QueryRecentTransactionDetailsRequest, db *database.ReadOnlyDatabase) ([]*asset.DepositTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	depositTxs := []*asset.DepositTransaction{}
	tx := dbConn.Where("asset_name = ? AND created_date>now()- interval '30 day' AND customer_id = ? AND customer_type = 0", req.CurrencySymbol, req.UserId).Order("created_date DESC").Find(&depositTxs)

	err = tx.Error
	if err != nil {
		return nil, err
	}
	return depositTxs, nil
}

func DBGetRecentPurchaseTransaction(req *QueryRecentTransactionDetailsRequest, db *database.ReadOnlyDatabase) ([]*asset.PurchaseTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	purchaseTxs := make([]*asset.PurchaseTransaction, 0)
	dbTx := dbConn.Model(new(asset.PurchaseTransaction)).Where("asset_name = ? AND created_date>now()- interval '30 day' AND user_id = ?", req.CurrencySymbol, req.UserId).Order("created_date DESC").Find(&purchaseTxs)
	err = dbTx.Error
	if err != nil {
		return nil, err
	}
	return purchaseTxs, nil
}

func DBGetRecentDistributedTokenTransaction(userId uint64, assetName string, db *database.ReadOnlyDatabase) ([]*reward.DistributedToken, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var distributedToken []*reward.DistributedToken = make([]*reward.DistributedToken, 0)

	tx := dbConn.Where("user_id = ? AND asset_name = ? AND created_date>now()- interval '30 day'", userId, assetName).FirstOrInit(&distributedToken)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return distributedToken, nil
}

func DBGetRecentTopUpTransaction(userId uint64, db *database.ReadOnlyDatabase) ([]*asset.TopUpTransaction, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	var topUpTransactionList []*asset.TopUpTransaction = make([]*asset.TopUpTransaction, 0)

	tx := dbConn.Where("customer_id = ? AND customer_type = ? AND created_date>now()- interval '30 day'", userId, asset.CustomerUser).Find(&topUpTransactionList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return topUpTransactionList, nil

}

func DBGetWalletAddressByUserId(userId uint64, db *database.ReadOnlyDatabase) (string, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return "", errors.New("Database Network Error: " + err.Error())
	}
	user := new(User)
	tx := dbConn.Where("id = ?", userId).Find(user)
	err = tx.Error
	if err != nil {
		return "", err
	}
	return user.WalletAddress, err
}
