package kyc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"fmt"
	"strings"
	"time"

	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"
	"eurus-backend/user_service/user_service/user"

	"gorm.io/gorm"
)

// All DB related function put in here

func DbGetKYCCountryList(db *database.ReadOnlyDatabase) (*[]kyc_model.KYCCountryCode, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var list []kyc_model.KYCCountryCode
	tx := dbConn.Find(&list)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func DbCreateKYCStatus(db *database.Database, newRecord *kyc_model.UserKYCStatus) (*kyc_model.UserKYCStatus, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	newRecord.KYCLevel = 1
	newRecord.KYCRetryCount = 0
	newRecord.KYCStatus = 0
	tx := dbConn.Create(&newRecord)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return newRecord, nil
}

func DbGetKYCStatusOfUser(db *database.ReadOnlyDatabase, userId uint64) (*[]kyc_model.UserKYCStatus, *user.User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, nil, err
	}
	var kycStatusList []kyc_model.UserKYCStatus = make([]kyc_model.UserKYCStatus, 0)
	var user user.User

	tx := dbConn.Session(&gorm.Session{}).Where("id = ? ", userId).FirstOrInit(&user)
	if tx.Error != nil {
		return nil, nil, err
	}
	if user.Id == 0 {
		return nil, nil, errors.New("User not found")
	}

	tx = dbConn.Session(&gorm.Session{}).Where("user_id = ?", userId).Find(&kycStatusList)
	if tx.Error != nil {
		return nil, nil, err
	}

	for i, kycStatus := range kycStatusList {
		var kycImageList []*kyc_model.UserKYCImage = make([]*kyc_model.UserKYCImage, 0)
		tx = dbConn.Session(&gorm.Session{}).Raw("SELECT t1.* FROM user_kyc_images as t1 INNER JOIN (SELECT user_kyc_id, doc_type, MAX(image_seq) as image_seq FROM user_kyc_images WHERE user_kyc_id = ? GROUP BY user_kyc_id, doc_type) as t2 ON t1.user_kyc_id = t2.user_kyc_id AND t1.doc_type = t2.doc_type and t1.image_seq=t2.image_seq WHERE t1.user_kyc_id = ?",
			kycStatus.Id, kycStatus.Id).Find(&kycImageList)
		if tx.Error != nil {
			return nil, nil, err
		}
		kycStatus.Images = kycImageList
		kycStatusList[i] = kycStatus
	}

	return &kycStatusList, &user, err
}

func DbGetKYCStatusList(db *database.ReadOnlyDatabase, statusId *kyc_model.KYCStatusType, email string, walletAddress string, page int, pageSize int) (*[]kyc_model.UserKYCStatusDetail, int64, error) {

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, 0, err
	}

	if email != "" {
		email = strings.ToLower(email)
	}
	if walletAddress != "" {
		walletAddress = ethereum.ToLowerAddressString(walletAddress)
	}

	var kycStatusList []kyc_model.UserKYCStatusDetail
	var totalRows int64
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		// return any error will rollback

		countQuery := tx.Model(&kycStatusList).Select("count(*)")

		query := tx.Debug().Select("t1.*, users.wallet_address, users.email")
		if statusId != nil {
			query = query.Where("kyc_status = ?", *statusId)
			countQuery = countQuery.Where("kyc_status = ?", *statusId)
		}
		email = strings.ToLower(email)
		if email != "" {
			query = query.Where("lower(users.email) = ?", email)
			countQuery = countQuery.Where("lower(users.email) = ?", email)
		}

		if walletAddress != "" {
			walletAddress = ethereum.ToLowerAddressString(walletAddress)
			query = query.Where("users.wallet_address = ?", walletAddress)
			countQuery = countQuery.Where("users.wallet_address = ?", walletAddress)
		}
		query.Joins("as t1 LEFT JOIN users ON users.id = t1.user_id").Order("last_modified_date DESC")
		countQuery.Joins("as t1 LEFT JOIN users ON users.id = t1.user_id")
		if err := query.Scopes(database.Paginate(page, pageSize)).Find(&kycStatusList).Error; err != nil {
			return err
		}
		if err := countQuery.Count(&totalRows).Error; err != nil {
			return err
		}
		return nil
	})
	return &kycStatusList, totalRows, err
}

func DbGetKYCStatusById(db *database.ReadOnlyDatabase, userKycId uint64) (*kyc_model.UserKYCStatus, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var kycStatus *kyc_model.UserKYCStatus = new(kyc_model.UserKYCStatus)
	tx := dbConn.Where("id = ?", userKycId).FirstOrInit(&kycStatus)
	err = tx.Error
	if err != nil {
		return nil, err
	}

	return kycStatus, nil

}

func DbResetKYCStatus(db *database.Database, userKycStatusId uint64, operatorId database.NullString) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		// return any error will rollback
		var kycStatus kyc_model.UserKYCStatus
		if err := tx.Session(&gorm.Session{}).First(&kycStatus, "id = ?", userKycStatusId).Error; err != nil {
			return err
		}
		// Because non-zero value cannot be updated with struct, so here we user map[string]interface
		if err := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCStatus{}).Where("id = ?", userKycStatusId).Updates(map[string]interface{}{"kyc_retry_count": 0, "kyc_status": kyc_model.KYCStatusPending, "operator_id": operatorId}).Error; err != nil {
			return err
		}

		if err := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCImage{}).Where("user_kyc_id = ?", userKycStatusId).Updates(map[string]interface{}{"status": kyc_model.KYCImageStatusVoided, "operator_id": operatorId}).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

func DbUpdateKYCStatusToApproval(db *database.Database, userKycStatusId uint64, statusId kyc_model.KYCStatusType) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		// return any error will rollback
		var kycStatusInDb kyc_model.UserKYCStatus
		if err := tx.Session(&gorm.Session{}).First(&kycStatusInDb, "id = ?", userKycStatusId).Error; err != nil {
			return err
		}
		if kycStatusInDb.KYCStatus == kyc_model.KYCStatusWaitingForApproval {
			return errors.New("You have already submit the approval request, please wait for CS Admin for further handling")
		} else if kycStatusInDb.KYCStatus == kyc_model.KYCStatusApproved || kycStatusInDb.KYCStatus == kyc_model.KYCStatusRejected {
			return errors.New("Cannot approve kyc that is approved or rejected already")
		}
		if kycStatusInDb.KYCRetryCount >= 3 {
			return errors.New("You can only submit KYC request Approval for max. 3 times. Please contact CS for more information")
		}
		if err := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCStatus{}).Where("id = ?", userKycStatusId).Updates(kyc_model.UserKYCStatus{KYCStatus: kyc_model.KYCStatusWaitingForApproval, KYCRetryCount: kycStatusInDb.KYCRetryCount + 1}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func DbInsertKYCImage(db *database.Database, userKycId uint64, imageType kyc_model.KYCImageType) (*kyc_model.UserKYCImage, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	type maxImageStatus struct {
		MaxSeq uint64                       `gorm:"column:image_seq"`
		Status kyc_model.KYCImageStatusType `gorm:"column:status"`
	}

	var result *maxImageStatus = new(maxImageStatus)
	dbSession := dbConn.Session(&gorm.Session{})

	subQuery := dbSession.Model(&kyc_model.UserKYCImage{}).Select("MAX(image_seq)").Where("user_kyc_id = ? AND doc_type = ?", userKycId, imageType).Group("doc_type").Limit(1)
	tx := dbSession.Model(&kyc_model.UserKYCImage{}).Select("image_seq", "status").Where("user_kyc_id = ? AND doc_type = ? AND image_seq = (?)", userKycId, imageType, subQuery).Find(&result)

	if tx.Error != nil {
		return nil, tx.Error
	}

	pendingImage := kyc_model.NewUserKYCImage(userKycId, imageType, int(result.MaxSeq+1), kyc_model.KYCImageStatusReceived, "", "", database.NullString{sql.NullString{Valid: false}})

	if result.MaxSeq > 0 && result.Status <= kyc_model.KYCImageStatusUploaded {
		dbConn.Session(&gorm.Session{}).Transaction(func(tx *gorm.DB) error {

			updateMap := map[string]interface{}{
				"status": kyc_model.KYCImageStatusApproved,
			}
			tx.Model(&kyc_model.UserKYCImage{}).Where("user_kyc_id = ? AND image_seq = ? ", userKycId, result.MaxSeq).Updates(updateMap)

			tx1 := tx.Session(&gorm.Session{}).Create(&pendingImage)

			return tx1.Error

		})
	} else if result.MaxSeq > 0 && result.Status == kyc_model.KYCImageStatusUploaded {
		errMsg := fmt.Sprintf("Document type already approved for user KYC id: %d image type: %d", userKycId, imageType)
		log.GetLogger(log.Name.Root).Infoln(errMsg)
		return nil, nil
	} else {
		tx := dbConn.Session(&gorm.Session{}).Create(&pendingImage)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return pendingImage, nil
	}

	return pendingImage, nil
}

func DbGetAdminUser(db *database.ReadOnlyDatabase, adminUser *kyc_model.AdminUser) (*kyc_model.AdminUser, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	tx := dbConn.Where("username = ? AND password = ?", adminUser.Username, adminUser.Password).First(&adminUser)
	err = tx.Error
	if err != nil {
		return nil, err
	}

	return adminUser, nil
}
func DbGetAdminUserByUsername(db *database.ReadOnlyDatabase, username string) (*kyc_model.AdminUser, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var adminUser *kyc_model.AdminUser = &kyc_model.AdminUser{}
	tx := dbConn.Where("username = ?", username).First(&adminUser)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return adminUser, nil
}

func DbUpdateKYCStatusAndRelatedImages(db *database.Database, kycStatus *kyc_model.UserKYCStatus) (int, *user.User, error) {
	dbConn, err := db.GetConn()
	operatorId := kycStatus.OperatorId
	var kycStatusInDb = &kyc_model.UserKYCStatus{}

	if err != nil {
		return 0, nil, err
	}

	for _, image := range kycStatus.Images {
		if image.UserKYCId != kycStatus.Id {
			return 0, nil, errors.New("Image user kyc id does not match with user kyc status ID")
		}
	}
	var isLevelPromoted bool = false
	var userInDb *user.User
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		resultTx := tx.Session(&gorm.Session{}).First(&kycStatusInDb, "id = ?", kycStatus.Id)
		err = resultTx.Error
		if err != nil {
			return err
		}
		if kycStatusInDb.KYCStatus != kyc_model.KYCStatusWaitingForApproval {
			return errors.New("KYC status is not pending for approval")
		}

		//1. updated each updated images
		for _, image := range kycStatus.Images {
			if image.Status == kyc_model.KYCImageStatusApproved {
				resultTx = tx.Session(&gorm.Session{}).Where("user_kyc_id = ? AND doc_type = ? AND image_seq = ? AND status = ?", image.UserKYCId, image.DocType, image.ImageSeq, kyc_model.KYCImageStatusUploaded).Updates(kyc_model.UserKYCImage{Status: image.Status, OperatorId: operatorId})
				err = resultTx.Error
				if err != nil {
					return err
				}
			} else if image.Status == kyc_model.KYCImageStatusWaitingForResubmit {
				resultTx = tx.Session(&gorm.Session{}).Where("user_kyc_id = ? AND doc_type = ? AND image_seq = ? AND status <> ?", image.UserKYCId, image.DocType, image.ImageSeq, kyc_model.KYCImageStatusApproved).Updates(kyc_model.UserKYCImage{Status: image.Status, OperatorId: operatorId, RejectReason: image.RejectReason})
				err = resultTx.Error
				if err != nil {
					return err
				}
			}
		}
		userKycStatusData, _ := json.Marshal(kycStatusInDb)

		log.GetLogger(log.Name.Root).Debug("User KYC status: ", string(userKycStatusData))
		if kycStatus.KYCStatus == kyc_model.KYCStatusApproved {
			// Check if all images are approved
			var kycImagesInDb = []kyc_model.UserKYCImage{}
			subQuery := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCImage{}).Where("user_kyc_id = ?", kycStatus.Id).Select("user_kyc_id,doc_type,MAX(image_seq) as image_seq").Group("user_kyc_id,doc_type")
			resultTx = tx.Session(&gorm.Session{}).Table("user_kyc_images as t1").Joins("inner join (?) as t2 ON t1.user_kyc_id = t2.user_kyc_id AND t1.doc_type = t2.doc_type and t1.image_seq=t2.image_seq", subQuery).Find(&kycImagesInDb)
			err = resultTx.Error
			if err != nil {
				return err
			}
			var requiredTypes []kyc_model.KYCImageType
			switch kycStatusInDb.KYCLevel {
			// For KYC level 1
			// Approval For Passport = Passport and selfie
			// Approval For ID Card	 = Id Card Front, ID Card Back, and selfie
			case 1:
				{
					if kycStatusInDb.KYCDoc == kyc_model.KYCDocIDCard {
						requiredTypes = []kyc_model.KYCImageType{kyc_model.KYCIDCardFront, kyc_model.KYCIDCardBack, kyc_model.KYCSelfie}
					} else if kycStatusInDb.KYCDoc == kyc_model.KYCDocPassport {
						requiredTypes = []kyc_model.KYCImageType{kyc_model.KYCImgePassport, kyc_model.KYCSelfie}
					}
					break
				}
			}

			// If any required images is not approved after this update, rollback
			imageIsApproved := false
			for t := range requiredTypes {
				imageIsApproved = false
				for i := range kycImagesInDb {
					if kycImagesInDb[i].DocType == requiredTypes[t] && kycImagesInDb[i].Status == kyc_model.KYCImageStatusApproved {
						imageIsApproved = true
						break
					}
				}
				if !imageIsApproved {
					return errors.New("[" + requiredTypes[t].String() + "] is not yet approved")

				}
			}

			resultTx = tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCStatus{}).Where("id = ?", kycStatus.Id).Updates(kyc_model.UserKYCStatus{KYCStatus: kycStatus.KYCStatus, OperatorId: operatorId, ApprovalDate: database.NullTime{sql.NullTime{Time: time.Now(), Valid: true}}})
			if resultTx.Error != nil {
				return err
			}
			if resultTx.RowsAffected == 0 {
				return errors.New("No record is approved")
			}

			if kycStatusInDb.UserId == 0 {
				log.GetLogger(log.Name.Root).Error("Unable to query user id is 0")
				return errors.New("Unable to get user")
			}

			log.GetLogger(log.Name.Root).Debug("Advance user: ", kycStatusInDb.UserId, " KYC level to ", kycStatusInDb.KYCLevel)
			resultTx = tx.Session(&gorm.Session{}).Table("users").Where("id = ?", kycStatusInDb.UserId).Update("kyc_level", kycStatusInDb.KYCLevel)
			err = resultTx.Error
			if err != nil {
				return err
			}
			userInDb = new(user.User)
			resultTx = tx.Session(&gorm.Session{}).Where("id = ? ", kycStatusInDb.UserId).FirstOrInit(&userInDb)
			err = resultTx.Error
			if err != nil {
				return err
			}

			isLevelPromoted = true
		} else if kycStatus.KYCStatus == kyc_model.KYCStatusWaitingForResubmit {
			if err := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCStatus{}).Where("id = ?", kycStatus.Id).Updates(kyc_model.UserKYCStatus{KYCStatus: kycStatus.KYCStatus, OperatorId: operatorId}).Error; err != nil {
				return err
			}
			//If the client submitted for 3 times, and the CS admin reject the submission on the 3rd try
			//Update the kyc status to rejected, which indicate this kyc status is rejected for any submission afterwards
			if kycStatusInDb.KYCRetryCount >= 3 {
				if err := tx.Session(&gorm.Session{}).Model(&kyc_model.UserKYCStatus{}).Where("id = ?", kycStatus.Id).Updates(kyc_model.UserKYCStatus{KYCStatus: kyc_model.KYCStatusRejected}).Error; err != nil {
					return err
				}
			}

		}

		return nil
	})

	if isLevelPromoted {
		return kycStatusInDb.KYCLevel, userInDb, err
	}
	return 0, nil, err
}

func DbUpdateKYCImageStatusToUploaded(db *database.Database, updateImageKey string, kycImage *kyc_model.UserKYCImage) error {
	conn, err := db.GetConn()
	if err != nil {
		return err
	}
	kycImage.ImagePath = updateImageKey
	kycImage.Status = kyc_model.KYCImageStatusUploaded

	tx := conn.Where("user_kyc_id = ? AND image_seq = ? AND doc_type = ? AND status < ?",
		kycImage.UserKYCId, kycImage.ImageSeq, kycImage.DocType, kyc_model.KYCImageStatusUploaded).Updates(&kycImage)

	return tx.Error
}

func DbCreateAdminUser(db *database.Database, reqObj *kyc_model.AdminUser) (*kyc_model.AdminUser, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	tx := dbConn.Create(&reqObj)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return reqObj, nil
}

func DbChangeAdminPassword(db *database.Database, reqObj *kyc_model.AdminUser) (*kyc_model.AdminUser, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	tx := dbConn.Model(&kyc_model.AdminUser{}).Where("username = ?", reqObj.Username).Update("password", reqObj.Password)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return reqObj, nil
}
