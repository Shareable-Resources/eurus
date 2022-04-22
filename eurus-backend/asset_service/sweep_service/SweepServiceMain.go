package main

import (
	sweepService "eurus-backend/asset_service/sweep_service/sweep"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

const KeyStringLen = 64

func main() {
	fmt.Println("Starting SweepService", secret.Tag)

	sweepService := sweepService.NewSweepService()
	commandLineArgs := &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	err := sweepService.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	sweepService.InitLog(sweepService.ServerConfig.LogFilePath)
	log.GetLogger(log.Name.Root).Infoln("Starting Sweep Service server ID: ", sweepService.Config.ServiceId)

	err = sweepService.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}

	sweepService.InitAll()
	sweepService.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}
