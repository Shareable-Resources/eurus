package sign

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/smartcontract/build/golang/contract"
)

func init() {
	err := ethereum.DefaultABIDecoder.ImportABIJson("UserWallet", contract.UserWalletABI)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to import UserWallet ABI: ", err)
		panic("Unable to import UserWallet ABI: " + err.Error())
	}
}
