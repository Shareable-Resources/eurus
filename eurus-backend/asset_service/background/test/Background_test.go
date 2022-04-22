package test

import (
	background "eurus-backend/asset_service/background/asset_bg"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	a := []string{"0-5x-long-bitcoin-cash-token", "1x-short-cosmos-token"}
	b := []string{"sats", "eth"}
	url := background.ConstructQueryCoingeckoUrl(a, b)
	c, _ := background.QueryCoingeckoExchangeRate(url)
	fmt.Println(c)
	fmt.Println(c["0-5x-long-bitcoin-cash-token"])

}

func TestMain(t *testing.T) {
	main.loadServerFromCMD()
}
