package conf_api

import (
	"eurus-backend/foundation/server"

	"github.com/mitchellh/mapstructure"
)

type ConfigMap struct {
	Id        int    `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsService *bool  `json:"is_service"`
}

func ConfigMapListToMap(configMapList []ConfigMap, excludeFieldList []string) map[string]interface{} {
	var excludeMap map[string]bool = make(map[string]bool)
	if excludeFieldList != nil {
		for _, field := range excludeFieldList {
			excludeMap[field] = true
		}
	}

	keyValueMap := make(map[string]interface{}, 0)
	for _, config := range configMapList {
		if _, found := excludeMap[config.Key]; !found {
			keyValueMap[config.Key] = config.Value
		}
	}

	return keyValueMap
}

//Assign ConfigMap key and value to server config base data model
func ConfigMapListToServerConfig(configMapList []ConfigMap, config server.IServerConfig) error {
	keyValMap := ConfigMapListToMap(configMapList, config.GetConfigFileOnlyFieldList())
	decodeConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           config, 
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(decodeConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(keyValMap)
}
