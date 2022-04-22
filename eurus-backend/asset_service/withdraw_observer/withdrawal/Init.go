package withdrawal

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
)

func init() {
	ethereum.DefaultABIDecoder.ImportABIJson("WithdrawSmartContract", contract.WithdrawSmartContractABI)

}
