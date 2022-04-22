package kyc

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"log"
)

func init() {
	err := ethereum.DefaultABIDecoder.ImportABIJson("EurusERC20", contract.EurusERC20ABI)
	if err != nil {
		fmt.Println("Unable to import ERC20 ABI: ", err, " program exit")
		log.Fatalln("Unable to import ERC20 ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("UserWallet", contract.UserWalletABI)

	if err != nil {
		fmt.Println("Unable to import UserWallet ABI: ", err, " program exit")
		log.Fatalln("Unable to import UserWallet ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("ApprovalWallet", contract.ApprovalWalletABI)

	if err != nil {
		fmt.Println("Unable to import ApprovalWallet ABI: ", err, " program exit")
		log.Fatalln("Unable to import ApprovalWallet ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("UserWalletProxy", contract.UserWalletProxyABI)
	if err != nil {
		fmt.Println("Unable to import UserWalletProxy ABI: ", err, " program exit")
		log.Fatalln("Unable to import UserWalletProxy ABI: ", err)
	}

}
