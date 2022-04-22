package bc_indexer

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"log"
)

func init() {
	err := ethereum.DefaultABIDecoder.ImportABIJson("EurusERC20", contract.EurusERC20ABI)
	if err != nil {
		fmt.Println("Unable to import EurusERC20 ABI: ", err, " program exit")
		log.Fatalln("Unable to import EurusERC20 ABI: ", err)
	}
	err = ethereum.DefaultABIDecoder.ImportABIJson("ERC20", contract.ERC20ABI)
	if err != nil {
		fmt.Println("Unable to import ERC20 ABI: ", err, " program exit")
		log.Fatalln("Unable to import ERC20 ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("UserWallet", contract.UserWalletABI)
	if err != nil {
		fmt.Println("Unable to import UserWallet ABI: ", err, " program exit")
		log.Fatalln("Unable to import UserWallet ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("UserWalletProxy", contract.UserWalletProxyABI)
	if err != nil {
		fmt.Println("Unable to import UserWalletProxy ABI: ", err, " program exit")
		log.Fatalln("Unable to import UserWalletProxy ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("EurusUserDeposit", mainnet_contract.EurusUserDepositABI)
	if err != nil {
		fmt.Println("Unable to import EurusUserDepsoit ABI: ", err, " program exit")
		log.Fatalln("Unable to import EurusUserDepsoit ABI: ", err)
	}

}
