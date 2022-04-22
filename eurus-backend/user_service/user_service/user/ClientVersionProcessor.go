package user

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/log"

	"github.com/mitchellh/mapstructure"
)

func ConfigMapListToClientVersion(configMapList []conf_api.ConfigMap, clientVersion *ClientVersion) error {
	keyValMap := conf_api.ConfigMapListToMap(configMapList, nil)

	decodeConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           clientVersion,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(decodeConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(keyValMap)
}

func QueryClientVersion(userServer *UserServer) (*ClientVersion, *foundation.ServerError) {
	queryReq := conf_api.NewQueryConfigRequest()
	queryReq.Id = userServer.ServerConfig.ServiceId

	resp := new(conf_api.QueryConfigResponse)
	reqRes := api.NewRequestResponse(queryReq, resp)

	_, err := userServer.SendConfigApiRequest(reqRes)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, resp.GetMessage())
	}

	if resp.GetReturnCode() < int64(foundation.Success) {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, resp.GetMessage())
	}

	clientVersion := new(ClientVersion)
	err = ConfigMapListToClientVersion(resp.Data.ConfigData, clientVersion)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
		return nil, foundation.NewErrorWithMessage(foundation.InternalServerError, resp.GetMessage())
	}
	return clientVersion, nil

}
