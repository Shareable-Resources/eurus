package main

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

func main() {
	log.NewLogger(log.Name.Root, "DummyConfirm.log", logrus.DebugLevel)
	//submitMintRequest("7e83fd446e38a203220e1d1d5e54aeb40995e5798c44c729edfc5f34ca8edf34")
	time2DoConfirm()
}

func submitMintRequest(privateKeyStr string) {
	ethClient := ethereum.EthClient{Protocol: "http", IP: "13.228.80.104", Port: 8545, ChainID: big.NewInt(2018)}
	ethClient.Connect()
	dummyPlatformWalletSC, err := contract.NewDummyPlatformWallet(common.HexToAddress("0x0216e1e274250128FE9021a5702b260cAf9767f8"), ethClient.Client)

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(privateKeyStr, ethClient.ChainID)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	tx, err := dummyPlatformWalletSC.SubmitMintRequest(transOpt, common.HexToAddress("0x10F8610FA2BA56264f222f8F0932102C4DF79ABf"), "USDT", big.NewInt(100), "0x6695374d5ee76d3cd6c2681d5de27e66d8bb3546de3c6476106c09e016c36f3a")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	receipt, err := ethClient.QueryEthReceipt(tx)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	receiptByte, _ := receipt.MarshalJSON()
	fmt.Println(string(receiptByte))
}

func confirmMintTrransaction(privateKeyStr string, txId int64) {
	ethClient := ethereum.EthClient{Protocol: "http", IP: "13.228.80.104", Port: 8545, ChainID: big.NewInt(2018)}
	ethClient.Connect()
	dummyPlatformWalletSC, err := contract.NewDummyPlatformWallet(common.HexToAddress("0x0216e1e274250128FE9021a5702b260cAf9767f8"), ethClient.Client)

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(privateKeyStr, ethClient.ChainID)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	transOpt.GasLimit = 1000000
	tx, err := dummyPlatformWalletSC.ConfirmTransaction(transOpt, big.NewInt(txId))
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	log.GetLogger(log.Name.Root).Debug("tx hash: " + tx.Hash().Hex())
	go func() {
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln(err.Error())
			panic(err)
		}
		receiptByte, _ := receipt.MarshalJSON()
		log.GetLogger(log.Name.Root).Debug("receipt: " + string(receiptByte))
	}()
}

func doConfirm(timeIn time.Time, privateKeyStr string, txId int64) {
	for {
		if timeIn.Unix() <= time.Now().Unix() {
			confirmMintTrransaction(privateKeyStr, txId)
			break
		}
	}
}

func time2DoConfirm() {
	timein := time.Now().Add(time.Second * 0)
	var transId int64 = 10
	doConfirm(timein, "ed2c36b7782314d2827d06d46d90ee1560e187ce2712ca8eee714f715f2eb49f", transId)
	doConfirm(timein, "e6ba4ea0deef2163270aa87d80ccb14bc15e0e9f648be715a723db5833f5b591", transId)
	doConfirm(timein, "0427b59bc5968a8c69b6a704002fc0ae1d02aec5c41001a453821ff35eab6d20", transId)
	doConfirm(timein, "39fbf4942e8b2ff9b71ec32c8c99b10321e6a3235ecdc1039cc3599b7844ef69", transId)
	doConfirm(timein, "c0f61339e6599fcf6e25aca4fd0d2d9ad25910836e9d0e1f9ca2422429a74371", transId)
	doConfirm(timein, "04292810e8a35ac50f958d6076e60a8a865fb49314ba6135425ee8fa6426e70f", transId)

	time.Sleep(10 * time.Second)
}
