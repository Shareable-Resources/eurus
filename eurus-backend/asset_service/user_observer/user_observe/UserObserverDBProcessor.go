package userObserver

import (
	"eurus-backend/foundation/log"
	user "eurus-backend/user_service/user_service/user"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

//db *database.Database

func DbGetCentralizedUser(context *UserObserverContext, tx *types.Transaction) (*user.User, error) {

	dbConn, err := context.slaveDb.GetConn()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Fail to connect db.", err)
		return nil, err
	}
	centralizedUser := new(user.User)
	if tx.To() == nil {
		return nil, nil
	}
	toAddr := strings.ToLower(tx.To().Hex())

	dbTx := dbConn.Where("is_metamask_addr = ? and wallet_address = ?", false, toAddr).FirstOrInit(centralizedUser)
	err = dbTx.Error
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Fail to get user:", err, " for trans hash: ", tx.Hash().Hex())
		return nil, err
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	} else {
		return centralizedUser, nil
	}
}
