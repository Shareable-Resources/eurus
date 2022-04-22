package main

import (
	background "eurus-backend/asset_service/background/asset_bg"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {

	fmt.Println("Starting Background Observer", secret.Tag)
	var observer *background.BackgroundServer = background.NewBackgroundServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = observer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	observer.InitLog(observer.ServerConfig.LogFilePath)
	log.GetLogger(log.Name.Root).Info("Starting BackgroundServer")

	err = observer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}
	observer.InitAll()
	log.GetLogger(log.Name.Root).Info("Finished BackgroundServer")

}

func main() {
	loadServerFromCMD()
}
