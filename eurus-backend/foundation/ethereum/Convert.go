package ethereum

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func HashToAddress(hash *common.Hash) *common.Address {
	addr := common.BytesToAddress(hash.Bytes())
	return &addr
}

func HashToInt256(hash *common.Hash) *uint256.Int {
	var returnInt *uint256.Int = uint256.NewInt(0)
	returnInt = returnInt.SetBytes(hash.Bytes())
	return returnInt
}

func ToLowerAddressString(addr string) string {
	addrObj := common.HexToAddress(addr)
	return strings.ToLower(addrObj.Hex())
}
