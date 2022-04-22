package conf

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/database"
	"strings"
)

func UpdateAuthToDB(db *database.Database, req *conf_api.UpdateAuthRequest) error {
	dbConn, err := db.GetConn()

	tx := dbConn.Model(&conf_api.AuthService{}).Where("id=? ", req.Id).Updates(map[string]interface{}{"service_name": req.ServiceName, "pub_key": req.PubKey})
	err = tx.Error
	if err != nil {
		return err
	}
	return nil
}

func InsertAuthToDB(db *database.Database, req *conf_api.AddAuthRequest) (*conf_api.AuthService, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	tx := dbConn.Select("service_name", "pub_key").Create(&req.AuthService)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return &req.AuthService, nil
}

func DelAuthFromDB(db *database.Database, req *conf_api.DeleteAuthRequest) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	auth := conf_api.AuthService{}
	tx := dbConn.Where("id = ?", req.Id).Delete(&auth)
	err = tx.Error
	if err != nil {
		return err
	}
	return nil
}

func GetAuthFromDB(db *database.ReadOnlyDatabase) (*conf_api.AuthInfo, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	data := new(conf_api.AuthInfo)
	tx := dbConn.Find(&data.AuthData)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return data, err
}

func DBUpdateAuthServiceWalletAddress(db *database.Database, serviceId uint64, walletAddress string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	authService := new(conf_api.AuthService)
	authService.WalletAddress = strings.ToLower(walletAddress)
	dbTx := dbConn.Where("id = ? ", serviceId).Updates(authService)

	return dbTx.Error

}
