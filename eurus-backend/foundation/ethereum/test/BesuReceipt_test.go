package test

import (
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"fmt"
	"math/big"
	"testing"
)

func TestRevertReason(t *testing.T) {
	ethClient := ethereum.EthClient{Protocol: "http",
		IP:      "13.228.80.104",
		Port:    8545,
		ChainID: big.NewInt(2018),
	}
	_, err := ethClient.Connect()
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := ethClient.GetConfirmedTransactionReceiptByHashStr("0x8c56027abe1d3c5e4f3da143fdf66bb4494ed05ca55c52e1255a5c422bef055c")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf(" Revert reason: %s\r\n", string(receipt.RevertReason))

	receiptByte, _ := json.Marshal(receipt)
	fmt.Printf("JSON: %s\r\n", string(receiptByte))

}
