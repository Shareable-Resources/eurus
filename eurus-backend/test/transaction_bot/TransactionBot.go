package main

import (
	"eurus-backend/foundation/ethereum"
	"fmt"
	"math/big"
)

var ownerAddr string = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
var privateKeyHex string = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
var account1Addr string = "0xf044d5768c970C09d6b4Def161b73d6515fA4b86"
var account1PrivateKey string = "236e7b7229a3f32ae8628064e1abcad334797d48ac51397f52979d60bf2e3e32"
var account2Addr string = "0x70C0a4CAb50FcFBF1A7bA0FCB3E8570D39Ab461d"
var account2PrivateKey string = "94a5e3bad57adb8bbd15d1415df7cde43250ffdf781cee923b21d480ca911948"
var account3Addr string = "0x7592bcb5C9574dcb96735B46252c747461AbCB95"
var account3PrivateKey string = "36c87b146c6917223e6a613ff93a0143a3319e7c86b9a942d06f3ce8e349f393"
var account4Addr string = "0x4DfB6d6790054F3EB68324BC230E3104137CA8Db"

const USDT string = "0x84B73aCceA0AC803B9882E2d6Ac557F201f3161b"

func initErc20Mint(client *ethereum.EthClient) {
	erc20, err := NewERC20(USDT, client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hash, err := MintToken(client, erc20, privateKeyHex, account1Addr, big.NewInt(100))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, err = MintToken(client, erc20, privateKeyHex, account2Addr, big.NewInt(100))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, err = MintToken(client, erc20, privateKeyHex, account3Addr, big.NewInt(100))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
}

func transferERC20(client *ethereum.EthClient) {
	erc20, err := NewERC20(USDT, client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	hash, err := TransferToken(client, erc20, account1PrivateKey, account2Addr, big.NewInt(1))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, err = TransferToken(client, erc20, account2PrivateKey, account3Addr, big.NewInt(1))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, err = TransferToken(client, erc20, account3PrivateKey, account1Addr, big.NewInt(1))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
}

func initTransferEth(client *ethereum.EthClient) {
	hash, _, err := client.TransferETH(privateKeyHex, account1Addr, big.NewInt(10000000000000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, _, err = client.TransferETH(privateKeyHex, account2Addr, big.NewInt(10000000000000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, _, err = client.TransferETH(privateKeyHex, account3Addr, big.NewInt(10000000000000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
}

func transferEth(client *ethereum.EthClient) {
	hash, _, err := client.TransferETH(account1PrivateKey, account2Addr, big.NewInt(10000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, _, err = client.TransferETH(account2PrivateKey, account3Addr, big.NewInt(10000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
	hash, _, err = client.TransferETH(account3PrivateKey, account1Addr, big.NewInt(10000))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(hash)
	}
}

func main() {
	ethClient := &ethereum.EthClient{
		Protocol: "http",
		IP:       "13.228.80.104",
		Port:     8545,
		ChainID:  big.NewInt(2018),
	}
	ethClient.Connect()
	//initTransferEth(ethClient)
	//initErc20Mint(ethClient)
	transferERC20(ethClient)
	//transferEth(ethClient)
}
