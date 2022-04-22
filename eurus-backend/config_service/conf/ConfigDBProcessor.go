package conf

import (
	"encoding/json"

	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/database"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var configServerId int = 0
var authServerId int = 1
var internalSmartContractConfigParam = "internalSCConfigAddress"
var eurusInternalConfigParam = "eurusInternalConfigAddress"

type AssetConfig struct {
	Id        uint64 `json:"id"`
	Key       string `json:"Key"`
	Value     string `json:"value"`
	IsService bool   `json:"is_service"`
}

func GetServerSettingFromDB(db *database.ReadOnlyDatabase) (*conf_api.ServerSetting, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var data conf_api.ConfigMap
	var tx *gorm.DB

	serverSetting := conf_api.ServerSetting{}

	tx = dbConn.Where("id = ? AND key = ?", configServerId, "ip").Find(&data)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	serverSetting.IP = data.Value

	tx = dbConn.Where("id = ? AND key = ?", configServerId, "port").Find(&data)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(data.Value)
	if err != nil {
		return nil, err
	}
	serverSetting.Port = port

	tx = dbConn.Where("id = ? AND key = ?", configServerId, "path").Find(&data)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	serverSetting.Path = data.Value
	return &serverSetting, tx.Error
}

func GetConfigFromDB(db *database.ReadOnlyDatabase, req *conf_api.QueryConfigRequest) (*[]conf_api.ConfigMap, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}
	var data []conf_api.ConfigMap
	var tx *gorm.DB

	if req.Key != "" {
		var hasId = true
		if req.Id == 0 {
			hasId = false
		}
		tx = dbConn.Raw(`SELECT * FROM config_maps WHERE id = ? and is_service = ? and key = ?
						UNION 
							SELECT * FROM config_maps 
							WHERE id = (SELECT COALESCE(( SELECT group_id FROM service_groups WHERE service_id = ? ) , -1 )) 
							AND key NOT IN (SELECT key FROM config_maps WHERE id = ? and is_service = ? AND key = ?)  AND is_service = ?
							AND key = ?
						UNION
							SELECT * FROM config_maps 
							WHERE id = 0
							AND key NOT IN (	SELECT key FROM config_maps WHERE id = ? and is_service = ? 
												UNION 
												SELECT key FROM config_maps 
												WHERE id = (SELECT COALESCE(( SELECT group_id FROM service_groups WHERE service_id = ? ) , -1 )) 
												AND key NOT IN (SELECT key FROM config_maps WHERE id = ? and is_service = ? AND key = ?)  AND is_service = ?
												AND key = ?
											) 
							AND is_service = ? AND key = ?`,
			req.Id, hasId, req.Key, req.Id, req.Id, true, req.Key, false, req.Key, req.Id, hasId, req.Id, req.Id, true, req.Key, false, req.Key, false, req.Key).Find(&data)
	} else {
		// Get config maps that are
		// 1. match req.Id and is_service is true
		// 2. all configs related to this req.Id service group (i.e. If req.Id under service_group group_id 2, than all configs related to service group_id 2 will be fetched)
		var hasId = true
		if req.Id == 0 {
			hasId = false
		}
		tx = dbConn.Raw(`SELECT * FROM config_maps WHERE id = ? and is_service = ? 
						UNION 
							SELECT * FROM config_maps 
							WHERE id = (SELECT COALESCE(( SELECT group_id FROM service_groups WHERE service_id = ? ) , -1 )) 
							AND key NOT IN (SELECT key FROM config_maps WHERE id = ? and is_service = ?)  AND is_service = ?
						UNION
							SELECT * FROM config_maps 
							WHERE id = 0
							AND key NOT IN (	SELECT key FROM config_maps WHERE id = ? and is_service = ? 
												UNION 
												SELECT key FROM config_maps 
												WHERE id = (SELECT COALESCE(( SELECT group_id FROM service_groups WHERE service_id = ? ) , -1 )) 
												AND key NOT IN (SELECT key FROM config_maps WHERE id = ? and is_service = ?)  AND is_service = ?
											) 
							AND is_service = ?`,
			req.Id, hasId, req.Id, req.Id, true, false, req.Id, hasId, req.Id, req.Id, true, false, false).Find(&data)

	}

	return &data, tx.Error
}

func GetConfigAuthFromDB(db *database.ReadOnlyDatabase, req *conf_api.QueryConfigAuthInfoRequest) (*conf_api.ConfigAuthInfo, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	data := conf_api.ConfigAuthInfo{}
	// tx := dbConn.Where("id = ? AND is_service = TRUE ", req.ServiceId).Find(&data.ConfigData)
	// err = tx.Error
	// if err != nil {
	// 	return nil, err
	// }
	queryConfigReq := new(conf_api.QueryConfigRequest)
	queryConfigReq.Id = req.ServiceId

	configMap, err := GetConfigFromDB(db, queryConfigReq)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to query config_maps")
	}
	data.ConfigData = *configMap
	tx := dbConn.Find(&data.AuthData)
	err = tx.Error
	if err != nil {
		return nil, err
	}
	return &data, err
}

func DelConfigFromDB(db *database.Database, req *conf_api.DeleteConfigRequest) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}

	config := conf_api.ConfigMap{}
	tx := dbConn.Where("id = ? AND key = ? AND is_service = ?", req.Id, req.Key, req.IsGroup).Delete(&config)
	err = tx.Error
	if err != nil {
		return err
	}
	return nil
}

func InsertConfigToDB(db *database.Database, req *conf_api.AddConfigRequest) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	tx := dbConn.Create(req.ConfigMap)
	err = tx.Error
	if err != nil {
		return err
	}
	return nil
}

func ConfigIsExist(db *database.ReadOnlyDatabase, req *conf_api.AddConfigRequest) (bool, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, err
	}
	configMap := &conf_api.ConfigMap{}
	tx := dbConn.Where("id = ? AND key = ?", req.Id, req.Key).Find(configMap)
	err = tx.Error
	if err != nil {
		return false, err
	} else if configMap.Key == "" && configMap.Value == "" {
		return false, nil
	}
	return true, nil
}

func UpdateConfigToDB(db *database.Database, req *conf_api.AddConfigRequest) error {
	dbConn, err := db.GetConn()
	if err != nil {
		return err
	}
	var tx *gorm.DB

	var requestData = *req.ConfigMap

	tx = dbConn.Where("id = ? AND key = ?", req.Id, req.Key).Find(req.ConfigMap)
	err = tx.Error
	if err != nil {
		return err
	}

	if requestData.Value != req.Value && requestData.Value != "" {
		tx = dbConn.Model(&conf_api.ConfigMap{}).Where("id = ? AND key = ?", req.Id, req.Key).Update("value", requestData.Value)
		err = tx.Error
		if err != nil {
			return err
		}
	}
	if requestData.IsService != req.IsService {
		tx = dbConn.Model(&conf_api.ConfigMap{}).Where("id = ? AND key = ?", req.Id, req.Key).Update("is_service", requestData.IsService)
		err = tx.Error
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFaucetConfigFromDB(db *database.ReadOnlyDatabase) ([]*conf_api.FaucetConfig, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	configMap := new(conf_api.ConfigMap)

	tx := dbConn.Where("key = ? AND id = ? AND is_service = ?", "faucetConfig", conf_api.ServiceGroupUser, false).FirstOrInit(configMap)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var faucetConfigs []*conf_api.FaucetConfig
	err = json.Unmarshal([]byte(configMap.Value), &faucetConfigs)
	if err != nil {
		return nil, err
	}

	return faucetConfigs, nil
}

func GetServerDetailByGroupFromDB(db *database.ReadOnlyDatabase, groupId uint) ([]*conf_api.ServerDetail, error) {

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	var serverList []*conf_api.ServerDetail = make([]*conf_api.ServerDetail, 0)

	dbConn = dbConn.Model(&conf_api.AuthService{}).Joins("INNER JOIN service_groups ON service_groups.service_id = auth_services.id").Where("service_groups.group_id = ?", groupId).Select("id", "service_name", "wallet_address").Find(&serverList)
	if dbConn.Error != nil {
		return nil, dbConn.Error
	}
	return serverList, nil

}

func TryGetSystemConfigFromDB(db *database.ReadOnlyDatabase, key string) (found bool, sc *conf_api.SystemConfig, err error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, nil, err
	}

	sc = new(conf_api.SystemConfig)

	// Found
	err = dbConn.Where("key = ?", key).First(&sc).Error
	if err == nil {
		return true, sc, nil
	}

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}

	// Other errors
	return false, nil, err
}

func InsertSystemConfigToDB(db *database.Database, ownerID int64, isService bool, key string, value string) (*conf_api.SystemConfig, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	sc := new(conf_api.SystemConfig)
	sc.OwnerID = ownerID
	sc.IsService = isService
	sc.Key = key
	sc.Value = value

	err = dbConn.Create(sc).Error
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func UpdateSystemConfigInDB(db *database.Database, key string, value string) (*conf_api.SystemConfig, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	err = dbConn.Where("key = ?", key).Updates(conf_api.SystemConfig{Value: value}).Error
	if err != nil {
		return nil, err
	}

	sc := new(conf_api.SystemConfig)

	// Select again to return
	err = dbConn.Where("key = ?", key).First(&sc).Error
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func TryGetGroupIDOfService(db *database.ReadOnlyDatabase, serviceID int64) (found bool, groupID int64, err error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return false, 0, err
	}

	sg := new(conf_api.ServiceGroup)

	// Table does not restrict ServiceID can only belong to 1 GroupID, but currently no such case
	// In future if this is changed, logic here also need to change
	err = dbConn.Where("service_id = ?", serviceID).First(&sg).Error
	if err == nil {
		return true, int64(sg.GroupId), nil
	}

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, nil
	}

	// Other errors
	return false, 0, err
}

func DBGetAssetSettings(db *database.Database) ([]conf_api.AssetSetting, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	var settings []conf_api.AssetSetting

	err = dbConn.Find(&settings).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func DBGetGroupIdByServiceId(db *database.ReadOnlyDatabase, serviceId int64) (int64, error) {
	dbConn, err := db.GetConn()
	if err != nil {
		return 0, err
	}
	serviceGroup := new(conf_api.ServiceGroup)
	tx := dbConn.Where("service_id = ?", serviceId).FirstOrInit(&serviceGroup)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int64(serviceGroup.GroupId), nil
}
