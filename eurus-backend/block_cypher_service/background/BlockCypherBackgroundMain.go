package main

import (
	"eurus-backend/block_cypher_service/background/block_cypher_bg"
	"eurus-backend/env"
	"eurus-backend/foundation/log"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {

	fmt.Println("Starting Block cypher background job", env.Tag)
	var bgJob *block_cypher_bg.BlockCypherBackground = block_cypher_bg.NewBlockCyhpherBackground()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = bgJob.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	bgJob.InitLog(bgJob.ServerConfig.LogFilePath)

	log.GetLogger(log.Name.Root).Info("Starting Block cypher background job")

	err = bgJob.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}

	<-bgJob.InitAll()

	log.GetLogger(log.Name.Root).Infoln("Finished")
}

func main() {
	loadServerFromCMD()
}
