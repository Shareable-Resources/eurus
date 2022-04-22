package main

import (
	"eurus-backend/foundation/log"
	walletBg "eurus-backend/report_service/wallet_background_indexer/wallet_bg"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"fmt"
)

func loadServerFromCMD() {

	fmt.Println("Starting Wallet Background Indexer", secret.Tag)
	var indexer *walletBg.WalletBackgroundIndexer = walletBg.NewWalletBackgroundServer()
	var commandLineArgs = &service.ServiceCommandLineArgs{}
	commandLineArgs.ParseCommandLineArgument()

	var err error
	err = indexer.LoadConfig(commandLineArgs)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to load config: ", err.Error())
		panic(err)
	}

	indexer.InitLog(indexer.ServerConfig.LogFilePath)

	log.GetLogger(log.Name.Root).Info("Starting WalletBackgroundIndexer")

	err = indexer.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}
	indexer.InitAll()
	log.GetLogger(log.Name.Root).Info("Finished WalletBackgroundIndexer")

}

func main() {
	loadServerFromCMD()
}
