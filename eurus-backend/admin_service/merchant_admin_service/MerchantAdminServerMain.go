package main

import (
	merchant_admin "eurus-backend/admin_service/merchant_admin_service/merchant"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting MerchantAdminServer", secret.Tag)
	merchantAdminServer := merchant_admin.NewMerchantAdminServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = merchantAdminServer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	merchantAdminServer.InitLog(merchantAdminServer.ServerConfig.LogFilePath)

	err = merchantAdminServer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}

	merchantAdminServer.InitAll()
	merchantAdminServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
