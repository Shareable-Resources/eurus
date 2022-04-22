package main

import (
	"eurus-backend/admin_service/admin_service/admin"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting AdminServer", secret.Tag)
	var adminServer *admin.AdminServer = admin.NewAdminServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = adminServer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	adminServer.InitLog(adminServer.ServerConfig.LogFilePath)

	err = adminServer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}

	adminServer.InitAll()
	adminServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
