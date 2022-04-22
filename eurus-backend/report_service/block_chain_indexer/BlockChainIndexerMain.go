package main

import (
	"eurus-backend/foundation/log"
	"eurus-backend/report_service/block_chain_indexer/bc_indexer"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {
	fmt.Println("Starting blockChainIndexer", secret.Tag)
	indexer := bc_indexer.NewBlockChainIndexer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = indexer.LoadConfig(commandLineArgs, indexer.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}
	//indexer.InitConfig()
	indexer.InitLog(indexer.ServerConfig.LogFilePath)

	err = indexer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}
	indexer.InitAll()
	indexer.InitUDSControlServer(&commandLineArgs.CommandLineArguments, nil)
}

func main() {
	loadServerFromCMD()
}
