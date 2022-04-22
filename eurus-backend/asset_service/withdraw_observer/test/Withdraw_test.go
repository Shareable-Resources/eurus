package test

import (
	"eurus-backend/asset_service/withdraw_observer/withdrawal"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetWithdrawTransactionHash(t *testing.T) {
	requestTransHash := common.HexToHash("0xd583a7cc4762ed99219a6b7e84e6df25080e2749237463b1ec48917e1db7d417")
	toAddr := common.HexToAddress("0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93")
	assetName := "USDT"
	amount := big.NewInt(1000000)

	hash, err := withdrawal.GetWithdrawTransactionHash(requestTransHash, toAddr, assetName, amount)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Println("requestTransHash:", requestTransHash)
	fmt.Println("toAddr:", toAddr)
	fmt.Println("assetName:", assetName)
	fmt.Println("amount:", amount)
	fmt.Println("result:", hash)
}
