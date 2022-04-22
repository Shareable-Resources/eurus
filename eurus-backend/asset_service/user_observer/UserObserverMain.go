package main

import (
	userObserver "eurus-backend/asset_service/user_observer/user_observe"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func main() {
	loadFromCMD()
}

func loadFromCMD() {
	fmt.Println("Starting UserObserver", secret.Tag)
	// Init observer to watch side chain block
	observer := userObserver.NewUserObserver()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()
	// Load server config from config path
	var err error = observer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}
	//Init log file from server config log file path, path ROOT
	observer.InitLog(observer.ServerConfig.LogFilePath)
	//start server logs will be saved to log.Name.Root
	log.GetLogger(log.Name.Root).Infoln("Starting user observer, server ID: ", observer.Config.ServiceId)
	//Loader Auth Server config
	err = observer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}
	observer.InitAll()
	observer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)

}
