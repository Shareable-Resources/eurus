package userObserver

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"log"
)

func init() {
	err := ethereum.DefaultABIDecoder.ImportABIJson("UserWallet", contract.UserWalletABI)
	if err != nil {
		fmt.Println("Unable to import UserWallet ABI: ", err, " program exit")
		log.Fatalln("Unable to import UserWallet ABI: ", err)
	}

}
