package service_server

import (
	"eurus-backend/auth_service/auth"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/env"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ServiceServer struct {
	server.ServerBase
	loginHandler           func(auth_base.IAuth)
	configEventHandler     func(message *amqp.Delivery, topic string, contentType string, content []byte)
	mqConfigServerConsumer *network.MQConsumer
	isLoggedIn             bool
}

func (me *ServiceServer) InitUDSControlServer(commandLineArgs *server.CommandLineArguments, handler func(req *server.ControlRequestMessage) (bool, string, error)) {

	fmt.Printf("\033[1;33m%s\033[0m\r\n", "Environment setting:")
	fmt.Printf("\033[1;33m%s %s\033[0m\r\n", "Environment:", env.Tag)
	fmt.Printf("\033[1;33m%s %s\033[0m\r\n", "Encrypted config file:", boolToYesNo(env.IsConfigEncrypted))
	fmt.Printf("\033[1;33m%s %s\033[0m\r\n", "Delete config after used:", boolToYesNo(env.IsDeleteConfigAfterUsed))

	log.GetLogger(log.Name.Root).Infoln("Environment setting:")
	log.GetLogger(log.Name.Root).Infoln("Environment: ", env.Tag)
	log.GetLogger(log.Name.Root).Infoln("Encrypted config file: ", boolToYesNo(env.IsConfigEncrypted))
	log.GetLogger(log.Name.Root).Infoln("Delete config after used: ", boolToYesNo(env.IsDeleteConfigAfterUsed))

	me.ServerBase.InitUDSControlServer(commandLineArgs, func(req *server.ControlRequestMessage) (bool, string, error) {
		var handled bool
		var err error
		var output string
		if handler != nil {
			handled, output, err = handler(req)
		}

		if err == nil && !handled {

			switch req.MethodName {
			case "verboseblocksubscriber":
				args := req.Data
				if len(args) == 0 {

					return true, "", errors.New("Invalid argument data type")
				}
				var executed bool = me.runVerboseBlockSubscriberCommand(args[0])

				if !executed {
					return true, "", errors.New("Usage: verboseBlockSubscriber [0/1]")
				}
				if args[0] == "0" {
					output += "Disabled\r\n"
				} else {
					output += "Enabled\r\n"
				}

				return true, output, nil
			case "help":
				return false, output + "verboseBlockSubscriber - Enable verbose logging on Block subscriber [0/1]\r\n", nil
			}
			return false, "", nil
		}

		return handled, output, err
	})
}

func boolToYesNo(boolVal bool) string {
	if boolVal {
		return "Yes"
	}

	return "No"
}

func (me *ServiceServer) runVerboseBlockSubscriberCommand(command string) bool {

	if command == "0" {
		if me.ServerBase.EthClient != nil && me.ServerBase.EthClient.Subscriber != nil {
			me.ServerBase.EthClient.Subscriber.IsVerboseLog = false
		}
		if me.ServerBase.MainNetEthClient != nil && me.ServerBase.MainNetEthClient.Subscriber != nil {
			me.ServerBase.MainNetEthClient.Subscriber.IsVerboseLog = false
		}
		return true
	} else if command == "1" {
		if me.ServerBase.EthClient != nil && me.ServerBase.EthClient.Subscriber != nil {
			me.ServerBase.EthClient.Subscriber.IsVerboseLog = true
		}

		if me.ServerBase.MainNetEthClient != nil && me.ServerBase.MainNetEthClient.Subscriber != nil {
			me.ServerBase.MainNetEthClient.Subscriber.IsVerboseLog = true
		}
		return true
	}

	return false
}

func (me *ServiceServer) InitAuth(loginHandler func(auth_base.IAuth), configEventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)) {
	authClient := auth.NewAuthClient()
	me.loginHandler = loginHandler
	me.configEventHandler = configEventHandler
	me.ServerBase.InitAuth(authClient, me.ServerConfig, me.processPostLogin)
}

func (me *ServiceServer) processPostLogin(authClient auth_base.IAuth) {
	if !me.isLoggedIn {
		if me.loginHandler != nil {
			me.loginHandler(authClient)
		}
		if me.configEventHandler != nil {
			err := me.SubscribeConfigServerEvent(me.ServerConfig, me.configEventHandler)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to subscribe config server MQ: ", err)
				panic(err)
			}
		}
		me.isLoggedIn = true

		if me.ServerConfig.HdWalletAddress != "" {
			setServerWalletAddrReq := conf_api.NewSetServerWalletAddressRequest()
			setServerWalletAddrReq.WalletAddress = me.ServerConfig.HdWalletAddress
			setServerWalletAddrRes := new(conf_api.SetServerWalletAddressFullResponse)
			reqRes := api.NewRequestResponse(setServerWalletAddrReq, setServerWalletAddrRes)
			_, err := me.SendConfigApiRequest(reqRes)
			if err != nil {
				log.GetLogger(log.Name.Root).Warn("Unable to set wallet address to config server: ", err)
			} else if setServerWalletAddrRes.ReturnCode != int64(foundation.Success) {
				log.GetLogger(log.Name.Root).Warn("Unable to set wallet address to config server, server return code: ", setServerWalletAddrRes.ReturnCode, " message: ", setServerWalletAddrRes.Message)
			}
		}
	}
}

// QueryAuthServerInfo will send a request and assign to ServerConfig properties : AuthServerIp, AuthServerPort, AuthPath
func (me *ServiceServer) QueryAuthServerInfo() error {
	queryReq := conf_api.NewQueryServerSettingRequest()

	resp := new(conf_api.QueryServerSettingResponse)
	reqRes := api.NewRequestResponse(queryReq, resp)

	_, err := me.SendConfigApiRequest(reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to send /config/setting request: ", err)
		return err
	}

	if resp.GetReturnCode() < int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("/config/setting response error code: ", resp.ReturnCode, " message: ", resp.Message)
		return errors.New(resp.GetMessage())
	}

	me.ServerConfig.AuthServerIp = resp.Data.IP
	me.ServerConfig.AuthServerPort = uint(resp.Data.Port)
	me.ServerConfig.AuthPath = resp.Data.Path
	return nil
}

//After QueryConfigServer,
func (me *ServiceServer) QueryConfigServer(config server.IServerConfig) (*[]conf_api.ConfigMap, error) {

	queryReq := conf_api.NewQueryConfigRequest()
	queryReq.Id = me.ServerConfig.ServiceId

	resp := new(conf_api.QueryConfigResponse)
	reqRes := api.NewRequestResponse(queryReq, resp)

	_, err := me.SendConfigApiRequest(reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to send /config request: ", err)
		return nil, err
	}

	if resp.GetReturnCode() < int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("/config response error code: ", resp.ReturnCode, " message: ", resp.Message)
		return nil, errors.New(resp.GetMessage())
	}

	var parseConfig interface{} = config

	for parseConfig != nil {
		err = conf_api.ConfigMapListToServerConfig(resp.Data.ConfigData, parseConfig.(server.IServerConfig))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
			return nil, err
		}
		parseConfig = parseConfig.(server.IServerConfig).GetParent()
	}

	config.GetServerConfigBase().HdWalletAddress = strings.ToLower(config.GetServerConfigBase().HdWalletAddress)

	serviceIdReq := conf_api.NewGetServiceGroupIdRequest()
	serviceIdRes := new(conf_api.GetServiceGroupIdFullResponse)
	servicerIdreqRes := api.NewRequestResponse(serviceIdReq, serviceIdRes)

	_, err = me.SendConfigApiRequest(servicerIdreqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query service group Id: ", err.Error())
		return nil, err
	}
	if serviceIdRes.ReturnCode != int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("Query service group Id response error code: ", serviceIdRes.ReturnCode, " message: ", serviceIdRes.Message)
		return nil, errors.New(serviceIdRes.Message)
	}
	me.ServerConfig.GroupId = serviceIdRes.Data

	log.GetLogger(log.Name.Root).Infoln("Config loaded from server successfully")
	return &resp.Data.ConfigData, nil
}

func (me *ServiceServer) LoadConfig(commandLineArgs *service.ServiceCommandLineArgs, config server.IServerConfig) error {
	return me.LoadConfigWithSetting(commandLineArgs, config, env.IsConfigEncrypted, env.IsDeleteConfigAfterUsed)
}

//Load File Based Config (XXXXXXXXServerConfig.json)
func (me *ServiceServer) LoadConfigWithSetting(commandLineArgs *service.ServiceCommandLineArgs, config server.IServerConfig,
	isEncrypted bool, isDeleteAfterUsed bool) error {

	err := me.ServerBase.LoadConfig(&commandLineArgs.CommandLineArguments,
		config, func(configPath string) ([]byte, error) {
			return service.LoadConfigFile(commandLineArgs.ConfigFilePath, isEncrypted, isDeleteAfterUsed, commandLineArgs.PasswordServerPath, commandLineArgs.UDSPath)
		})
	if err != nil {
		return err
	}

	return secret.DecryptSensitiveConfig(config.GetServerConfigBase())
}

func (me *ServiceServer) QuerySystemConfig(configName string) (string, error) {
	req := conf_api.NewGetSystemConfigRequest(configName)
	res := new(conf_api.GetSystemConfigFullResponse)
	reqRes := api.NewRequestResponse(req, res)
	_, err := me.SendConfigApiRequest(reqRes)
	if err != nil {
		return "", err
	}

	if res.ReturnCode != int64(foundation.Success) {
		return "", errors.New("Server code: " + strconv.FormatInt(res.ReturnCode, 10) + res.Message)
	}

	return res.Data.Value, nil

}

func (me *ServiceServer) QueryAssets() ([]*conf_api.Asset, error) {
	req := conf_api.NewGetAssetRequest()
	res := new(conf_api.GetAssetFullResponse)
	reqRes := api.NewRequestResponse(req, res)
	_, err := me.SendConfigApiRequest(reqRes)
	if err != nil {
		return nil, err
	}

	if res.ReturnCode != int64(foundation.Success) {
		return nil, errors.New("Server code: " + strconv.FormatInt(res.ReturnCode, 10) + res.Message)
	}

	return res.Data, nil
}

func (me *ServiceServer) QueryAssetSettings() ([]conf_api.AssetSetting, error) {
	req := conf_api.NewGetAssetSettingsRequest()
	res := new(conf_api.GetAssetSettingsFullResponse)
	reqRes := api.NewRequestResponse(req, res)
	_, err := me.SendConfigApiRequest(reqRes)
	if err != nil {
		return nil, err
	}

	if res.ReturnCode != int64(foundation.Success) {
		return nil, errors.New("Server code: " + strconv.FormatInt(res.ReturnCode, 10) + res.Message)
	}

	return res.Data, nil
}

func (me *ServiceServer) GetServiceIdFromServerLoginToken(loginToken auth_base.ILoginToken) (int64, error) {
	return me.AuthClient.(*auth.AuthClient).GetServiceIdFromServerLoginToken(loginToken)
}

func (me *ServiceServer) SubscribeConfigServerEvent(config *server.ServerConfigBase, configEventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)) error {
	me.mqConfigServerConsumer = new(network.MQConsumer)
	me.mqConfigServerConsumer.Logger = log.GetLogger(log.Name.Root)

	metaData := new(network.MQTaskQueueMetaData)
	metaData.IsAutoAck = true
	metaData.IsExclusive = false
	metaData.IsAutoDelete = conf_api.ConfigExchangeMetaData.IsAutoDelete
	if config.GroupId == 0 {
		metaData.QueueName = "general_queue_" + strconv.Itoa(int(config.ServiceId))
	} else {
		metaData.QueueName = "general_queue_group_" + strconv.Itoa(int(config.GroupId))
	}

	return me.mqConfigServerConsumer.SubscribeTopic(me.ServerConfig.GetMqUrl(), "config.*", &conf_api.ConfigExchangeMetaData,
		metaData, me.configEventHandler)
}
