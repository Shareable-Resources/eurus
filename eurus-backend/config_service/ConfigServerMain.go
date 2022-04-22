package main

import (
	"eurus-backend/config_service/conf"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting ConfigServer ", secret.Tag)
	var configServer *conf.ConfigServer = conf.NewConfigServer()
	var commandLineArgs = new(service.ServiceCommandLineArgs)
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = configServer.LoadConfig(commandLineArgs, configServer.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	configServer.InitLog(configServer.ServerConfig.LogFilePath)
	log.GetLogger(log.Name.Root).Info("Starting Config Server")

	err = configServer.InitDBFromConfig(configServer.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load DB: ", err.Error())
		panic(err)
	}

	err = configServer.LoadConfigFromDB()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config from DB: ", err.Error())
		panic(err)
	}

	err = configServer.LoadConfigServerConfigFromDB()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config server config from DB: ", err.Error())
		panic(err)
	}

	configServer.InitHDWallet()
	fmt.Println("Wallet address: ", configServer.ServerConfig.HdWalletAddress)
	log.GetLogger(log.Name.Root).Infoln("Wallet address: ", configServer.ServerConfig.HdWalletAddress)

	configServer.InitAuth(nil, nil)

	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthClient: ", err.Error())
		panic(err)
	}

	err = configServer.InitConfigMQ()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init MQ: ", err)
		panic(err)
	}
	go func() {
		err = configServer.InitHttpServer(nil)
		if err != nil {
			panic(err)
		}
		err = configServer.HttpServer.Listen()
		if err != nil {
			panic(err)
		}
	}()
	configServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)

}

func main() {
	loadServerFromCMD()
}
