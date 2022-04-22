package test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"

)

func TestAddressWithoutPrefix(t *testing.T){
	addressWithoutPrefix:="51e2DD0d7F8224413d9c1C6aa22582abeDf23a39"
	fmt.Println("address with prefix: ", common.HexToAddress(addressWithoutPrefix).Hex())
}
