package conf

import (
	"encoding/json"
	"errors"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"
	"time"

	"github.com/shopspring/decimal"
	//"fmt"
)

func GetAssetListFromDB(db *database.ReadOnlyDatabase, req *conf_api.GetAssetRequest) *response.ResponseBase {
	var err error
	dbConn, err := db.GetConn()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get config and auth info"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	newAsset := make([]*conf_api.Asset, 0)
	tx := dbConn.Find(&newAsset)
	err = tx.Error
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get config and auth info"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, newAsset)
	return res
}

func AddAssetToDB(me *ConfigServer, req *conf_api.SetAssetRequest) error {
	dbConn, err := me.DefaultDatabase.GetConn()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		//res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return err
	}
	err = req.ValidateSetAsset()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		//res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return err
	}
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		//res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return err
	}

	tx := dbConn.Create(req.Asset)
	err = tx.Error
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Request Params Error: "+err.Error(), "\nRequest Params: "+string(reqStr))
		//res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return err
	}

	//res := response.CreateSuccessResponse(req, tx)
	return nil

}

func SetExchangeRateToDB(db *database.Database, rate decimal.Decimal, assetName string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	time := time.Now()

	exchangeRate := new(conf_api.ExchangeRates)

	tx := dbConn.Where("asset_name = ?", assetName).FirstOrInit(exchangeRate)
	err = tx.Error
	if err != nil {
		return err
	}

	affectedRow := tx.RowsAffected
	exchangeRate.AssetName = assetName
	exchangeRate.Rate = rate
	exchangeRate.LastModifiedDate = time

	if affectedRow == 0 {
		exchangeRate.AssetName = assetName
		exchangeRate.CreatedDate = time
		tx = dbConn.Create(exchangeRate)
	} else {
		tx = dbConn.Where("asset_name = ?", assetName).Updates(exchangeRate)
	}

	err = tx.Error
	return err
}

func GetAssetDecimalFromDB(db *database.ReadOnlyDatabase, assetName string) (int64, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return 0, err
	}

	getDecimal := new(conf_api.Asset)
	tx1 := dbConn.Where("asset_name = ?", assetName).FirstOrInit(getDecimal)
	err = tx1.Error
	if err != nil {
		return 0, err
	}
	if tx1.RowsAffected == 0 {
		err := errors.New("Cannot find this currency in asset table")
		return 0, err
	}
	return getDecimal.Decimal, nil
}

func SetETHWithdrawFeeToDB(db *database.Database, ethAdminFee string) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	configMap := new(conf_api.ConfigMap)
	const KeyName = "ETHWithdrawAdminFee"
	tx := dbConn.Where("key = ? AND id = ?", KeyName, -1).Find(configMap)
	err = tx.Error
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query asset config key = ", KeyName, ". Error: ", err)
		return err
	}

	var isService *bool = new(bool)
	*isService = false
	configMap.Value = ethAdminFee
	configMap.Key = KeyName
	configMap.IsService = isService
	configMap.Id = -1

	if tx.RowsAffected == 0 {
		tx = dbConn.Create(configMap)
	} else {
		tx = dbConn.Where("key = ? AND id = ? AND is_service = ? ", KeyName, -1, false).Updates(configMap)
	}

	err = tx.Error
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to update asset config key = ", KeyName, ". Error :", err)
		return err
	}
	return nil
}
