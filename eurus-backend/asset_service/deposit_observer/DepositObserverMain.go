package main

import (
	"eurus-backend/asset_service/deposit_observer/deposit"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func main() {
	loadFromCMD()
}

func loadFromCMD() {
	fmt.Println("Starting depositObserver", secret.Tag)
	observer := deposit.NewDepositObserver()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = observer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	observer.InitLog(observer.ServerConfig.LogFilePath)
	log.GetLogger(log.Name.Root).Infoln("Starting deposit observer server ID: ", observer.Config.ServiceId)
	err = observer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}
	observer.InitAll()
	observer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}
