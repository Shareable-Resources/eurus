package main

import (
	"eurus-backend/asset_service/withdraw_observer/withdrawal"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting Withdraw observer", secret.Tag)
	observer := withdrawal.NewWithdrawObserver()

	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = observer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	observer.InitLog(observer.ServerConfig.LogFilePath)

	err = observer.InitAll()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init: ", err.Error())
		panic(err)
	}
	observer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
