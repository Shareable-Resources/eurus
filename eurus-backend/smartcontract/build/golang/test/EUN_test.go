package test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEUN_GetBalance(t *testing.T) {
	ethClient := initEthClient(t)

	platformWallet := getAddressBySmartContractName("OwnedUpgradeabilityProxy<MarketingRegWallet>")

	balance, err := ethClient.GetBalance(platformWallet)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Balance: ", EUN_FormatBalance(balance))
}

func EUN_FormatBalance(balance *big.Int) string {

	decimal := big.NewInt(int64(math.Pow10(18)))
	reminder := big.NewInt(0).Mod(balance, decimal)

	result := big.NewInt(0).Div(balance, decimal)

	return result.String() + "." + fmt.Sprintf("%018d", reminder.Uint64())
}

func TestEUN_GetWithdrawObserverBalance(t *testing.T) {
	ethClient := initEthClient(t)
	for _, addr := range withdrawobserverAddr {
		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(addr, EUN_FormatBalance(balance))
	}

}
