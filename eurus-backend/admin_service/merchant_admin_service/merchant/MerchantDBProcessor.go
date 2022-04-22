package merchant_admin

import (
	"errors"
	"eurus-backend/admin_service/merchant_common"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantAdminDBProcessor struct {
	db *database.Database
}

func NewMerchantAdminDBProcessor(db *database.Database) *MerchantAdminDBProcessor {
	merchantAdminDBProcessor := new(MerchantAdminDBProcessor)
	merchantAdminDBProcessor.db = db
	return merchantAdminDBProcessor
}

func (me *MerchantAdminDBProcessor) QueryRefundRequestList(merchantId uint64, status *merchant_common.RefundRequestStatus, pageNum int, pageSize int) (int64, []merchant_common.MerchantRefundRequest, error) {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return 0, nil, err
	}

	var list []merchant_common.MerchantRefundRequest = make([]merchant_common.MerchantRefundRequest, 0)
	dbTx := dbConn.Session(&gorm.Session{}).Model(new(merchant_common.MerchantRefundRequest)).Where("merchant_id = ? ", merchantId)
	if status != nil {
		dbTx = dbTx.Where("status = ? ", status)
	}
	var recordCount int64 = 0
	dbTx = dbTx.Count(&recordCount)

	dbTx = dbTx.Order("created_date DESC").Scopes(paginate(dbTx, pageNum, pageSize)).Find(&list)
	if dbTx.Error != nil {
		return 0, nil, dbTx.Error
	}

	return recordCount, list, nil
}

func (me *MerchantAdminDBProcessor) getMerchantByUserName(merchantId uint64, username string) (*MerchantAdmin, error) {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	admin := new(MerchantAdmin)
	txConn := dbConn.Where("merchant_id = ? AND username = ?", merchantId, username).FirstOrInit(&admin)

	if txConn.Error != nil {
		return nil, txConn.Error
	}
	return admin, nil
}

func (me *MerchantAdminDBProcessor) updateMerchantPassword(merchantId uint64, username string, password string) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return err
	}
	admin := new(MerchantAdmin)
	admin.PasswordHash = password
	admin.LastModifiedDate = time.Now()
	admin.Status = MerchantAccountNormal

	txConn := dbConn.Where("merchant_id = ? AND username = ?", merchantId, username).Updates(admin)
	return txConn.Error
}

func (me *MerchantAdminDBProcessor) UpdateRefundRequest(merchantId uint64, opearatorId uint64, data *RefundRequest) error {
	dbConn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	model := new(merchant_common.MerchantRefundRequest)
	model.LastModifiedDate = time.Now()
	model.MerchantOperatorId = &opearatorId
	model.Status = data.Answer
	model.OperatorComment = data.Comment
	if data.RefundTransHash != "" {
		model.RefundTransHash = ethereum.ToLowerAddressString(data.RefundTransHash)
	}

	err = dbConn.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {

		checkModel := new(merchant_common.MerchantRefundRequest)
		//SELECT FOR UPDATE
		dbTx := dbConn.Session(&gorm.Session{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", data.RequestId).Find(&checkModel)
		if dbTx.Error != nil {
			return dbTx.Error
		}

		if dbTx.RowsAffected == 0 {
			return errors.New("Request not found")
		}

		if checkModel.MerchantId != merchantId {
			return errors.New("Request is for another merchant")
		}

		if checkModel.Status != merchant_common.RefundPending {
			return errors.New("Request is not in pending state")
		}

		dbTx = dbConn.Session(&gorm.Session{}).Where("merchant_id = ? and id = ? and status = ? ", merchantId, data.RequestId, merchant_common.RefundPending).Updates(&model)
		if dbTx.Error != nil {
			return dbTx.Error
		}

		if dbTx.RowsAffected == 0 {
			return errors.New("Request not found")
		}
		return nil
	})

	return err
}

func paginate(db *gorm.DB, pageNum int, pageSize int) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNum - 1) * pageSize
		dbTx := db.Offset(offset).Limit(pageSize)
		return dbTx
	}

}
