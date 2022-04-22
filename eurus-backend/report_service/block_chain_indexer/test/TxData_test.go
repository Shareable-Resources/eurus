package test

import (
	"encoding/hex"
	"eurus-backend/report_service/block_chain_indexer/bc_indexer"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

var testData1 string = "0xa9059cbb0000000000000000000000007592bcb5c9574dcb96735b46252c747461abcb950000000000000000000000000000000000000000000000000000000000000001"
var testData2 string = "0xa90591230000000000000000000000007592bcb5c9574dcb96735b46252c747461abcb950000000000000000000000000000000000000000000000000000000000000001"

func TestFilterTransfer(t *testing.T){
	dataByte,err:=hex.DecodeString(testData2[2:])

	if(err!=nil){
		t.Error(err.Error())
	}
	dataMap,err,state:=bc_indexer.ExtractTransferData(dataByte)
	if(err!=nil){
		t.Error(err.Error())
	}
	fmt.Println("state: ",state)
	fmt.Println("Recipient: ",dataMap["recipient"].(common.Address).Hex())
	fmt.Printf("Amount: %v\n",dataMap["amount"].(*big.Int))
}