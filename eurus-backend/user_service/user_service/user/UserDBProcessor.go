package user

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func DbGetUserById(userId uint64, db *database.ReadOnlyDatabase) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	user := new(User)
	tx := dbConn.Where("id = ?", userId).Find(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func DbGetUserByWalletAddress(address string, db *database.ReadOnlyDatabase) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	address = strings.ToLower(address)
	user := new(User)
	tx := dbConn.Where("wallet_address = ?", address).Find(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, foundation.NewError(foundation.UserNotFound)
	}

	return user, nil
}

func DbGetUserByLoginAddress(address string, db *database.ReadOnlyDatabase) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	address = strings.ToLower(address)
	user := new(User)
	tx := dbConn.Where("login_address = ?", address).FirstOrInit(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, foundation.NewError(foundation.UserNotFound)
	}

	return user, nil
}

func DbAddNewUser(address string, db *database.Database, isPending bool, email string, isMetaMaskUser bool) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	address = strings.ToLower(address)
	user := new(User)
	user.LoginAddress = address
	user.WalletAddress = address
	user.MainnetWalletAddress = address
	user.KycLevel = 0
	user.IsMetamaskAddr = isMetaMaskUser
	user.Email = email
	if isPending {
		user.Status = UserStatusNotVerify

	} else {
		user.Status = UserStatusNormal
	}
	user.LastModifiedDate = time.Now()
	user.CreatedDate = user.LastModifiedDate

	tx := dbConn.Create(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DbUpdateNewCentralizedUser(address string, db *database.Database, email string) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	user := new(User)
	user.LoginAddress = address
	user.WalletAddress = address
	user.MainnetWalletAddress = address
	user.KycLevel = 0
	user.IsMetamaskAddr = false
	user.Email = email
	user.Status = UserStatusNotVerify

	user.LastModifiedDate = time.Now()
	user.LastLoginTime = time.Time{}
	user.CreatedDate = user.LastModifiedDate
	tx := dbConn.Session(&gorm.Session{}).Model(&user).
		Where("email = ? AND status NOT IN ? ", email, []UserStatus{UserStatusNormal, UserStatusSuspended}).
		Select("*").Omit("id").
		Updates(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}

	tx = dbConn.Session(&gorm.Session{}).Where("email = ? AND status = ?", email, UserStatusNotVerify).Find(&user)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Unable to query updated user")
	}
	return user, nil
}

func DbCheckVerificationExist(userId uint64, db *database.ReadOnlyDatabase, verificationType VerificationType) (bool, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return true, err
	}

	verification := new(Verification)
	tx := dbConn.Table("verifications").Where("user_id = ? AND type = ?", userId, verificationType).Find(verification)

	if tx.Error != nil {
		return true, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func DbAddNewVerification(userId uint64, db *database.Database, duration int, verificationType VerificationType) (*Verification, *foundation.ServerError) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}

	verification := new(Verification)
	verification.UserId = userId
	verification.LastModifiedDate = time.Now()
	verification.CreatedDate = time.Now()
	verification.ExpiredTime = time.Now().Add(time.Second * time.Duration(duration))
	verification.Type = int(verificationType)
	verification.Count = 1

	code := ""
	for i := 0; i < 6; i++ {
		code = code + strconv.Itoa(rand.Intn(10))
	}

	verification.Code = code

	searchVerification := new(Verification)

	tx := dbConn.Session(&gorm.Session{}).Where("user_id = ? AND type = ?", userId, verificationType).First(&searchVerification)
	if tx.Error != nil {
		tx := dbConn.Session(&gorm.Session{}).Create(&verification)
		err = tx.Error
		if err != nil {
			return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
		}
		return verification, nil
	} else {
		if searchVerification.ExpiredTime.Before(time.Now()) {
			//Reset everything
			dbConn.Session(&gorm.Session{}).Where("user_id = ? AND type = ?", userId, verificationType).Updates(&verification)
		} else {
			return DbUpdateVerification(userId, db, duration, verificationType)
		}
	}

	return verification, nil
}

func DbUpdateVerification(userId uint64, db *database.Database, duration int, verificationType VerificationType) (*Verification, *foundation.ServerError) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}

	code := ""
	for i := 0; i < 6; i++ {
		code = code + strconv.Itoa(rand.Intn(10))
	}

	verification := new(Verification)
	tx := dbConn.Where("user_id = ? AND type = ?", userId, verificationType).Find(&verification)

	if tx.Error != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, tx.Error.Error())
	}

	verification.Count++

	if verification.Count%4 > 0 {
		if time.Now().Unix()-verification.LastModifiedDate.Unix() < 60 {
			return nil, foundation.NewErrorWithMessage(foundation.RequestTooFrequenct, "Request too frequent, wait for 1 minute")
		}
	} else {
		if time.Now().Unix()-verification.LastModifiedDate.Unix() < 300 {
			return nil, foundation.NewErrorWithMessage(foundation.RequestTooFrequenct, "Request too frequent, wait for 5 minutes")
		}
	}

	tx = dbConn.Table("verifications").Where("user_id = ? AND type = ?", userId, verificationType).Updates(map[string]interface{}{
		"count":              verification.Count,
		"code":               code,
		"last_modified_date": time.Now(),
		"expired_time":       time.Now().Add(time.Second * time.Duration(duration)),
	}).Debug()
	if tx.Error != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, tx.Error.Error())
	}

	verification.Code = code

	return verification, nil

}

func DbUpdateLoginTime(address string, db *database.Database) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	address = strings.ToLower(address)
	user := new(User)
	err = dbConn.Transaction(func(tx *gorm.DB) error {
		var err error
		resultTx := tx.Session(&gorm.Session{}).Model(&user).Where("login_address = ?", address).Updates(map[string]interface{}{"last_login_time": time.Now()}).Debug()
		err = resultTx.Error
		if err != nil {
			return err
		}
		resultTx = tx.Session(&gorm.Session{}).Where("login_address = ?", address).FirstOrInit(&user)
		err = resultTx.Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DbGetNextUserId(db *database.Database) (uint64, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return 0, err
	}
	userId := new(UserId)
	tx := dbConn.Raw("SELECT nextval(pg_get_serial_sequence(?,?))", "users", "id").Scan(userId)
	err = tx.Error
	if err != nil {
		return 0, err
	}
	return userId.NextVal, nil
}

//TODO query config server instead
func DbGetAdminFeeDecimal(db *database.ReadOnlyDatabase, assetName string) (int64, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return 0, err
	}
	asset := new(conf_api.Asset)

	if assetName != "ETH" {

		tx := dbConn.Where("asset_name = ?", assetName).Find(&asset)

		if tx.Error != nil {
			err = tx.Error
			return 0, err
		}

		if asset.AssetName == "" {
			return 0, errors.New("asset not found")
		}
	} else {
		asset.Decimal = 18
	}

	return asset.Decimal, nil

}

type userFaucets struct {
	LastModifiedDate time.Time
	Status           int
	TransHash        string
}

func DbCheckIfUserAllowFaucet(db *database.ReadOnlyDatabase, assetName string, userId uint64) (faucetRecord *userFaucets, valid bool, pending bool, err error) {

	userFaucetsInstance := new(userFaucets)

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, false, false, err
	}

	tx := dbConn.Table("user_faucets").Where("user_id = ? AND key = ? AND status >= 0", userId, assetName).First(userFaucetsInstance)

	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			return nil, false, false, tx.Error
		}
	}

	if tx.RowsAffected == 0 {
		return nil, true, false, nil
	}

	if userFaucetsInstance.Status == 1 {
		return userFaucetsInstance, false, true, nil
	}

	if userFaucetsInstance.LastModifiedDate.Unix()+86400 > time.Now().Unix() {
		return userFaucetsInstance, false, false, nil
	}

	return userFaucetsInstance, true, false, nil
}

func DbAddUserFaucetRecord(db *database.Database, key string, userId uint64, txHash string, status int) error {

	type userFaucets struct {
		LastModifiedDate time.Time
		UserId           uint64
		Key              string
		TransHash        string
		Status           int
		CreatedDate      time.Time
	}
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	user := new(userFaucets)
	tx := dbConn.Table("user_faucets").Where("user_id = ? AND key = ?", userId, key).First(user)

	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			return tx.Error
		}

		user.UserId = userId
		user.Key = key
		user.LastModifiedDate = time.Now()
		user.CreatedDate = time.Now()
		user.Status = status
		user.TransHash = txHash

		tx = dbConn.Table("user_faucets").Create(user)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}

	user.LastModifiedDate = time.Now()
	user.TransHash = txHash
	user.Status = status
	tx = dbConn.Table("user_faucets").Where("user_id = ? AND key = ?", userId, key).Save(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func DbCheckEmailExist(db *database.Database, email string) (bool, *User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, nil, err
	}

	userObj := new(User)
	tx := dbConn.Where("email = ?", email).Find(&userObj)

	if tx.Error != nil {
		return false, nil, tx.Error
	}

	if tx.RowsAffected > 0 {
		return true, userObj, nil
	}

	return false, nil, nil

}

//Join verifications table on user.id to check whether the first record of the joined table[code] which is  order by created_date DESC is the same as code from the method
func DbVerifyCode(db *database.Database, userId uint64, code string, verificationType int) (bool, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, err
	}

	type result struct {
		UserId      uint64
		Code        string
		ExpiredTime time.Time
		Type        int16
	}

	results := new(result)
	tx := dbConn.Session(&gorm.Session{})
	tx = tx.Model(&User{}).Joins("LEFT JOIN verifications on verifications.user_id = users.id").Select("verifications.user_id", "verifications.code", "verifications.expired_time", "verifications.type").Where("users.id = ? AND verifications.type = ?", userId, verificationType).Order("verifications.created_date DESC").FirstOrInit(results)

	if tx.Error != nil {
		return false, err
	}

	if results.Code != code {
		return false, nil
	}

	if results.ExpiredTime.Before(time.Now()) {
		return false, errors.New("code expired")
	}
	tx = dbConn.Session(&gorm.Session{})
	tx = tx.Table("verifications").Where("user_id = ? AND code = ? AND type = ?", userId, results.Code, verificationType).Updates(map[string]interface{}{"expired_time": time.Now()})
	if tx.Error != nil {
		return false, tx.Error
	}

	return true, nil

}

func DbAddOrUpdateUserDevice(db *database.Database, customerId uint64, customerType int16, deviceId string, publicKey string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	ud := new(UserDevice)
	tx := dbConn.Where("customer_id = ? AND customer_type = ? AND device_id = ?", customerId, customerType, deviceId).First(ud)

	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			return tx.Error
		}

		// Add new user device
		ud.InitDate()
		ud.CustomerId = customerId
		ud.CustomerType = customerType
		ud.DeviceId = deviceId
		ud.PubKey = publicKey

		tx := dbConn.Create(ud)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}

	tx = dbConn.Where("customer_id = ? AND customer_type = ? AND device_id = ?", customerId, customerType, deviceId).Updates(UserDevice{PubKey: publicKey, DbModel: database.DbModel{LastModifiedDate: time.Now()}})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DbGetUserDevicePublicKey(db *database.ReadOnlyDatabase, customerId uint64, customerType int16, deviceId string) (string, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return "", err
	}

	ud := new(UserDevice)
	tx := dbConn.Where("customer_id = ? AND customer_type = ? AND device_id = ?", customerId, customerType, deviceId).FirstOrInit(ud)

	if tx.Error != nil {
		return "", tx.Error
	}

	return ud.PubKey, nil
}

func DbUpdateUserRegisterSuccessful(db *database.Database, userId uint64, mnemonic string, mainnetAddr string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	tx := dbConn.Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"status": UserStatusVerifiedNotSetPaymentAddress, "mnemonic": mnemonic, "last_modified_date": time.Now(), "mainnet_wallet_address": mainnetAddr})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DbUpdateUserOwnerWalletAddress(db *database.Database, userId uint64, ownerWalletAddress string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	tx := dbConn.Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"owner_wallet_address": ownerWalletAddress, "last_modified_date": time.Now()})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DbUpdateUserWalletAddress(db *database.Database, userId uint64, walletAddress string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	walletAddress = strings.ToLower(walletAddress)
	tx := dbConn.Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"wallet_address": walletAddress, "status": UserStatusNormal, "last_modified_date": time.Now(), "last_login_time": time.Now()})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DbCheckUnVerifiedUserExistById(db *database.ReadOnlyDatabase, userId uint64) (bool, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return true, err
	}
	var count int64
	tx := dbConn.Table("users").Where("id = ? AND status in ?", userId, []UserStatus{UserStatusNotVerify, UserStatusVerifiedNotSetPaymentAddress}).Count(&count)
	if tx.Error != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func DbInsertLoginRequestToken(db *database.Database) (*LoginRequestTokenMap, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	intString := "0123456789"
	code := ""
	for i := 0; i < 6; i++ {
		code = code + string(intString[rand.Intn(10)])
	}
	now := time.Now()
	loginRequestTokenMap := new(LoginRequestTokenMap)
	loginRequestTokenMap.LoginRequestToken = code
	loginRequestTokenMap.CreatedDate = now
	loginRequestTokenMap.LastModifiedDate = now
	expiredTime := now.Add(time.Second * time.Duration(300))
	loginRequestTokenMap.ExpiredTime = &expiredTime

	tx := dbConn.Table("login_request_token_maps").Create(loginRequestTokenMap)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return loginRequestTokenMap, nil

}

func DbUpdateLoginRequestToken(db *database.Database, token string, loginRequestToken string, userId string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	tx := dbConn.Table("login_request_token_maps").Where("login_request_token = ? AND expired_time >= ? AND token IS NULL", loginRequestToken, time.Now()).Updates(map[string]interface{}{"token": token, "last_modified_date": time.Now(), "user_id": userId})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("login request token record not found")
	}

	return nil
}

func DbCheckIfLoginRequestTokenValid(db *database.ReadOnlyDatabase, loginRequestToken string) (bool, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, err
	}

	type loginRequest struct {
		Id int
	}

	result := new(loginRequest)

	tx := dbConn.Table("login_request_token_maps").Where("login_request_token = ? AND expired_time >= ? AND token IS NULL", loginRequestToken, time.Now()).Scan(result)

	if tx.Error != nil {
		return false, err
	}

	if tx.RowsAffected != 1 {
		return false, nil
	}
	return true, nil
}

func DbChangeUserOwnerWalletAddress(db *database.Database, userId uint64, ownerWalletAddress string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	ownerWalletAddress = ethereum.ToLowerAddressString(ownerWalletAddress)
	user, err := DbGetUserById(userId, &db.ReadOnlyDatabase)
	if err != nil {
		return err
	}

	tx := dbConn.Table("users").Where("id = ? AND status = ?", userId, UserStatusNormal).Updates(map[string]interface{}{"last_modified_date": time.Now(), "owner_wallet_address": ownerWalletAddress, "change_payment_password_count": user.ChangePaymentPasswordCount + 1})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil

}

func DbUpdateUserLoginAddress(db *database.Database, userId uint64, loginAddress string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	loginAddress = ethereum.ToLowerAddressString(loginAddress)
	user, err := DbGetUserById(userId, &db.ReadOnlyDatabase)
	if err != nil {
		return err
	}

	tx := dbConn.Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"login_address": loginAddress, "last_modified_date": time.Now(), "change_login_password_count": user.ChangeLoginPasswordCount + 1})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func DbGetUserStorageByID(reqObj *RequestUserStorage, db *database.ReadOnlyDatabase) (*UserPreferenceStorage, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	userStorage := new(UserPreferenceStorage)
	tx := dbConn.Where("user_id = ? and platform = ?", reqObj.UserId, reqObj.Platform).Find(&userStorage)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return userStorage, nil
}

func DbUpdateUserStorage(reqObj *RequestUserStorage, db *database.Database) (*UserPreferenceStorage, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	userStorage, err := DbGetUserStorageByID(reqObj, &db.ReadOnlyDatabase)
	if err != nil {
		return nil, err
	}
	userStorage.Storage = reqObj.Storage
	if userStorage.UserId == 0 && userStorage.Sequence == 0 {
		userStorage.Sequence = 0
		userStorage.UserId = uint64(reqObj.UserId)
		userStorage.Platform = uint8(reqObj.Platform)
		tx := dbConn.Create(userStorage)
		err = tx.Error
	} else {
		tx := dbConn.Model(&userStorage).Where("user_id = ? and platform = ?", reqObj.UserId, reqObj.Platform).Updates(map[string]interface{}{
			"sequence": gorm.Expr("sequence + ?", 1), "storage": reqObj.Storage})
		err = tx.Error

		userStorage.Sequence += 1 //for return value
	}
	if err != nil {
		return nil, err
	}

	return userStorage, nil
}

func DbGetUserByEmail(db *database.ReadOnlyDatabase, email string) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	user := new(User)
	tx := dbConn.Where("lower(email) = ?", strings.ToLower(email)).Find(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil

}

func DbGetUserByOwnerWalletAddress(db *database.ReadOnlyDatabase, ownerWalletAddress string) (*User, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	user := new(User)
	ownerWalletAddress = strings.ToLower(ownerWalletAddress)
	tx := dbConn.Where("owner_wallet_address = ?", ownerWalletAddress).Find(&user)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil

}

func DbUpdateUserMnemonicPhase(db *database.Database, userId uint64, mnemonic string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	tx := dbConn.Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"mnemonic": mnemonic, "last_modified_date": time.Now()})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
