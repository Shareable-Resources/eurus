package main

import (
	"eurus-backend/env"
	"eurus-backend/password_service/password"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting Password Server ", env.Tag)
	var passwordServer *password.PasswordServer = password.NewPasswordServer()
	var commandLineArgs = new(service.ServiceCommandLineArgs)
	commandLineArgs.ParseCommandLineArgument()

	err := passwordServer.LoadConfig(&commandLineArgs.CommandLineArguments, passwordServer.Config, nil)
	if err != nil {
		panic("Load config error: " + err.Error())
	}
	passwordServer.InitLog(passwordServer.ServerConfig.LogFilePath)

	passwordServer.InitAll()

	fmt.Println("Password server start listening at: ", passwordServer.Config.UDSPath)
	passwordServer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
