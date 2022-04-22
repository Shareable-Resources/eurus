package auth

import (
	"errors"
	"eurus-backend/foundation/database"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenType int16

const (
	RefreshableToken TokenType = iota
	NonRefreshableToken
	ResetLoginPasswordToken
	ResetPaymentPasswordToken
)

func DbInsertUserSession(db *database.Database, token string, serviceId int64, expiredTime int64, userId string, tokenType TokenType) error {
	expiredTimeObj := time.Unix(expiredTime, 0)
	var expiredTimePtr *time.Time = &expiredTimeObj

	if expiredTime == 0 {
		expiredTimePtr = nil
	}
	//var tokenType int
	//if refreshable {
	//	tokenType = 0
	//} else {
	//	tokenType = 1
	//}
	userSession := &UserSession{Token: token, ServiceId: serviceId, ExpiredTime: expiredTimePtr, UserId: userId, Type: int16(tokenType)}
	userSession.InitDate()
	dbConn, err := db.GetConn()

	if err != nil {
		return err
	}
	tx := dbConn.Create(userSession)

	return tx.Error
}

//func DBInsertLoginRequestTokenMap(db *database.Database, loginRequestToken string, expiredTime int64)error{
//	expiredTimeObj := time.Unix(expiredTime, 0)
//	var expiredTimePtr *time.Time = &expiredTimeObj
//	if expiredTime == 0 {
//		expiredTimePtr = nil
//	}
//	now := time.Now()
//	dbConn, err := db.GetConn()
//	loginRequestTokenMap:=new(LoginRequestTokenMap)
//	loginRequestTokenMap.LoginRequestToken = loginRequestToken
//	loginRequestTokenMap.ExpiredTime = expiredTimePtr
//	loginRequestTokenMap.CreatedDate = now
//	loginRequestTokenMap.LastModifiedDate = now
//	if err != nil {
//		return err
//	}
//	tx := dbConn.Create(loginRequestTokenMap)
//	if tx.Error != nil {
//		return tx.Error
//	}
//
//	return nil
//
//}

/// If record not found, both return value are nil
func DbQueryUserSession(db *database.ReadOnlyDatabase, token string, serviceIdList []int64) (*UserSession, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	userSession := new(UserSession)
	var tx *gorm.DB
	if len(serviceIdList) > 0 {
		tx = dbConn.Where("(service_id IN ? or service_id = 0 ) and token = ? and disabled <> TRUE", serviceIdList, token).First(&userSession)
	} else {
		tx = dbConn.Where(" service_id = 0 and token = ? and disabled <> TRUE", token).First(&userSession)
	}

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			//Record not found
			return nil, nil
		}
		return nil, tx.Error
	}

	return userSession, nil
}

/// If record not found, both return nil, error
func DbRevokeUserSession(db *database.Database, token string, serviceId int64) (*UserSession, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	userSession := new(UserSession)

	err = dbConn.Transaction(func(tx *gorm.DB) error {
		stx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token = ?", token).First(&userSession)

		if stx.Error != nil {
			return stx.Error
		}

		if userSession.ServiceId != serviceId {
			return errors.New("server id mismatch")
		}

		stx = tx.Model(&UserSession{}).Where("service_id = ? and token = ?", serviceId, token).Updates(UserSession{Disabled: true, DbModel: database.DbModel{LastModifiedDate: time.Now()}})
		if stx.Error != nil {
			return stx.Error
		}
		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return nil, err
	}
	return userSession, nil
}

/// If record not found, both return nil, error
func DbRefreshUserSession(db *database.Database, token string, serviceId int64) (*UserSession, error) {

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	userSession := new(UserSession)

	err = dbConn.Transaction(func(tx *gorm.DB) error {
		stx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token = ?", token).First(&userSession)

		if stx.Error != nil {
			return stx.Error
		}

		if userSession.ServiceId != serviceId {
			return errors.New("server id mismatch")
		}

		if userSession.ExpiredTime == nil {
			// the case where token will not be expired
			return nil
		}
		if userSession.ExpiredTime.Before(time.Now()) {
			return errors.New("token expired")
		}
		if userSession.Type != 0 {
			return errors.New("This token is a non-refreshable token!")
		}

		allowTime := userSession.LastModifiedDate.Add(time.Duration(tokenTimeoutSecond*6/10) * time.Second)
		if time.Now().Before(allowTime) {
			return errors.New("token update too frequently")
		}

		newExpireTime := userSession.ExpiredTime.Add(time.Duration(tokenTimeoutSecond) * time.Second)
		stx = tx.Model(&UserSession{}).Where("service_id = ? and token = ?", serviceId, token).Updates(UserSession{ExpiredTime: &newExpireTime, DbModel: database.DbModel{LastModifiedDate: time.Now()}})
		if stx.Error != nil {
			return stx.Error
		}

		userSession.ExpiredTime = &newExpireTime
		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return nil, err
	}
	return userSession, nil
}
