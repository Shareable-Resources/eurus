package main

import (
	"eurus-backend/asset_service/approval_observer/approval"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting Approval Observer", secret.Tag)
	var observer *approval.ApprovalObserver = approval.NewApprovalObserver()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()
	var err error
	err = observer.LoadConfig(commandLineArgs, observer.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	observer.InitLog(observer.ServerConfig.LogFilePath)

	err = observer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}

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
