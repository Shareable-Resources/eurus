package main

import (
	"eurus-backend/env"
	"eurus-backend/foundation/log"
	"eurus-backend/service_base/service"
	"eurus-backend/user_service/user_service/user"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting UserServer", env.Tag)
	var userServer *user.UserServer = user.NewUserServer()
	var commandLineArgs = new(service.ServiceCommandLineArgs)
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = userServer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	userServer.InitLog(userServer.ServerConfig.LogFilePath)

	err = userServer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}
	userServer.InitAll()
	userServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
