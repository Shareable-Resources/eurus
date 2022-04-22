package user

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
)

func init() {
	ethereum.DefaultABIDecoder.ImportABIJson("EurusERC20", contract.EurusERC20ABI)
	ethereum.DefaultABIDecoder.ImportABIJson("DAppStockBase", contract.DAppStockBaseABI)
}
