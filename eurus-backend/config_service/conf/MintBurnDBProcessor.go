package conf

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/database"
	"strconv"

	"gorm.io/gorm"
)

func AddMintBurnInfoToDB(configServer *ConfigServer, db *database.Database, pubKey string, req *conf_api.MintBurnServerSetting, groupID conf_api.ServiceGroupId) (*conf_api.ConfigMap, error) {

	dbConn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	var serviceName string

	switch groupID {
	case conf_api.ServiceGroupDeposit:
		serviceName = "mint (" + req.Name + ")"
	case conf_api.ServiceGroupWithdraw:
		serviceName = "burn (" + req.Name + ")"
	case conf_api.ServiceGroupApproval:
		serviceName = "Approval (" + req.Name + ")"
	case conf_api.ServiceGroupUser:
		serviceName = "User (" + req.Name + ")"
	case conf_api.ServiceGroupBlockchainIndexer:
		serviceName = "Blockchain Indexer (" + req.Name + ")"
	}

	isServiceState := true
	authService := &conf_api.AuthService{ServiceName: serviceName, PubKey: pubKey}
	config := &conf_api.ConfigMap{IsService: &isServiceState}
	serviceGroup := &conf_api.ServiceGroup{}

	err = dbConn.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Create(authService).Error; err != nil {
			return err
		}
		config.Id = authService.Id
		config.Key = "ethClientIP"
		if req.EthClientIP == "" {
			config.Value = configServer.ServerConfig.EthClientIP
		} else {
			config.Value = req.EthClientIP
		}
		if err = tx.Create(&config).Error; err != nil {
			return err
		}
		config.Key = "ethClientPort"
		if req.EthClientPort == 0 {
			config.Value = strconv.Itoa(configServer.ServerConfig.EthClientPort)
		} else {
			config.Value = strconv.Itoa(req.EthClientPort)
		}
		if err = tx.Create(&config).Error; err != nil {
			return err
		}
		config.Key = "name"
		config.Value = serviceName
		if err = tx.Create(&config).Error; err != nil {
			return err
		}

		for i := 0; i < len(req.MiscellaneousData); i++ {
			config.Key = req.MiscellaneousData[i].Key
			config.Value = req.MiscellaneousData[i].Value
			if err = tx.Create(&config).Error; err != nil {
				return err
			}
		}

		// for serviceGroup
		serviceGroup.GroupId = int(groupID)

		serviceGroup.ServiceId = authService.Id
		if err = tx.Create(&serviceGroup).Error; err != nil {
			return err
		}

		return nil
	})
	return config, nil
}
