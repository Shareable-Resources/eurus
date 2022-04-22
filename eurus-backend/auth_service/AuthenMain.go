package main

import (
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func main() {
	fmt.Println("Starting AuthServer", secret.Tag)
	var authServer *auth.AuthServer = auth.NewAuthServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}

	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = authServer.LoadConfig(commandLineArgs, authServer.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	authServer.InitLog(authServer.Config.LogFilePath)
	log.GetLogger(log.Name.Root).Info("Starting AuthServer")

	authServer.QueryConfigServer()
	err = authServer.InitDBFromConfig(&authServer.Config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load DB: ", err.Error())
		panic(err)
	}
	authServer.InitWebSocketServer(authServer.Config, "/ws")
	authServer.InitHttpServer(nil)

	go authServer.HttpServerListen()

	authServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)

}
