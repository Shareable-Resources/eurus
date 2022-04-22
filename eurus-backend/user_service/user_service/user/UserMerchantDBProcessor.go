package user

import (
	"eurus-backend/admin_service/merchant_common"
	"eurus-backend/foundation/database"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func DbInsertMerchantRefundRequest(db *database.Database, req *RequestMerchantRefundRequest, user *User) (*uint64, error) {
	dbConn, err := db.GetConn()

	if err != nil {
		return nil, err
	}
	model := new(merchant_common.MerchantRefundRequest)
	model.InitDate()
	model.Amount = decimal.NewFromBigInt(big.NewInt(0).SetUint64(req.Amount), 0)
	model.AssetName = req.AssetName
	model.MerchantId = req.MerchantId
	model.RefundReason = req.Reason
	model.Status = merchant_common.RefundPending
	model.UserId = &user.Id
	model.DestAddress = user.WalletAddress

	dbTx := dbConn.Create(&model)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return &model.Id, nil
}

func DbGetMerchantIdByWalletAddress(db *database.Database, walletAddress common.Address) (uint64, error) {

	//TODO
	return 1, nil
}

func DbGetMerchantRefundStatus(db *database.ReadOnlyDatabase, userId uint64) ([]*merchant_common.MerchantRefundRequest, error) {

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	refundList := make([]*merchant_common.MerchantRefundRequest, 0)
	dbTx := dbConn.Model(new(merchant_common.MerchantRefundRequest)).Where("user_id = ?", userId).Order("last_modified_date DESC").Find(&refundList)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return refundList, nil

}
