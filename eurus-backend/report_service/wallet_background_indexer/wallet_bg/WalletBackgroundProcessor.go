package wallet_bg

import (
	"eurus-backend/foundation/log"
	wallet_bg_model "eurus-backend/report_service/wallet_background_indexer/wallet_bg/model"
	"eurus-backend/user_service/user_service/user"
)

// 1. Batch Insert Wallet Balance, this will insert batch of wallet balances
func BatchInsertWalletBalance(server *WalletBackgroundIndexer, list *[]wallet_bg_model.WalletBalance) error {

	err := DbBatchInsertWalletBalance(server.DefaultDatabase, list)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot insert wallet balance")
	}
	return err
}

// 2. Get All User cold wallet address in mainnet and wallet address in side chain
func GetUserWalletAddressesFromDb(server *WalletBackgroundIndexer) (*[]user.User, error) {

	list, err := DbGetUserWalletAddressesFromDb(server.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot insert wallet balance")
	}
	return list, err
}
