package deposit

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"log"
)

func init() {
	err := ethereum.DefaultABIDecoder.ImportABIJson("ERC20", contract.ERC20ABI)
	if err != nil {
		fmt.Println("Unable to import ERC20 ABI: ", err, " program exit")
		log.Fatalln("Unable to import ERC20 ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("EurusUserDeposit", mainnet_contract.EurusUserDepositABI)
	if err != nil {
		fmt.Println("Unable to import EurusUserDeposit ABI: ", err, " program exit")
		log.Fatalln("Unable to import EurusUserDeposit ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("OwnedUpgradeabilityProxy", contract.OwnedUpgradeabilityProxyABI)
	if err != nil {
		fmt.Println("Unable to import OwnedUpgradeabilityProxy ABI: ", err, " program exit")
		log.Fatalln("Unable to import OwnedUpgradeabilityProxy ABI: ", err)
	}

	err = ethereum.DefaultABIDecoder.ImportABIJson("PlatformWallet", contract.PlatformWalletABI)
	if err != nil {
		fmt.Println("Unable to import PlatformWallet ABI: ", err, " program exit")
		log.Fatalln("Unable to import PlatformWallet ABI: ", err)
	}

}
