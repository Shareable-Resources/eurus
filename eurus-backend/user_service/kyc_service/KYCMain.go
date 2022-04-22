package main

import (
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/user_service/kyc_service/kyc"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting KYC Server", secret.Tag)
	//1. Load Config.json from local files and assigned to server.ServerConfig
	var kycServer *kyc.KYCServer = kyc.NewKYCServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()
	var err error
	err = kycServer.LoadConfig(commandLineArgs, kycServer.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}
	kycServer.InitLog(kycServer.ServerConfig.LogFilePath)
	/*
		err = kycServer.QueryAuthServerInfo()
		if err != nil {
			log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
		}*/

	err = kycServer.InitAll()

	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init: ", err.Error())
		panic(err)
	}

	kycServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
