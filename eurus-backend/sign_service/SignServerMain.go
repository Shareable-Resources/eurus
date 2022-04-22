package main

import (
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/sign_service/sign"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting SignServer", secret.Tag)
	var signServer *sign.SignServer = sign.NewSignServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = signServer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	signServer.InitLog(signServer.ServerConfig.LogFilePath)

	err = signServer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}
	signServer.InitAll()
	signServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, signServer.Calibrate)
}

func main() {
	loadServerFromCMD()
}
