package test

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

var counter *ethereum.ScanBlockCounter
var mode = ethereum.ScanBlockModeDiscret

func TestInitScanBlockCounter(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(mode, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}
	err = counter.Start()
	if err != nil {
		t.Fatal("Start error: ", err)
	}

	time.Sleep(1 * time.Minute)
}

func RescanBlockHandler(sender *ethereum.ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal) {
	fmt.Printf("Event Id: %v, chaindId: %v, from: %v, to: %v\r\n", eventId, chainId, from, to)

	three := big.NewInt(3)
	newBlockNum := big.NewInt(0)

	newBlockNum = newBlockNum.Add(from.BigInt(), three)

	status, err := sender.UpdateRescanEvent(eventId, newBlockNum)
	if err != nil {
		fmt.Println("Error at UpdateRescanEvent: ", err.Error())
	} else {
		fmt.Println("Status : ", status)
	}
}

func TestUpdateScanBlock(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(mode, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}
	err = counter.Start()
	if err != nil {
		t.Fatal("Start error: ", err)
	}
	progress, err := counter.GetLatestProgress()
	if err != nil {
		t.Fatal("GetLatestProgress error: ", err)
	}

	one := decimal.NewFromInt(1)
	nextBlock := progress.StartNum.Add(one)

	err = counter.UpdateLatestBlock(nextBlock.BigInt())
	if err != nil {
		t.Fatal("UpdateLatestBlock error: ", err)
	}
}

func TestTriggerRescanBlock(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(mode, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}
	err = counter.Start()
	if err != nil {
		t.Fatal("Start error: ", err)
	}

	progress, err := counter.GetLatestProgress()
	if err != nil {
		t.Fatal("GetLatestProgress error: ", err)
	}
	ten := decimal.NewFromInt(10)
	nextBlock := progress.StartNum.Add(ten)

	if mode == ethereum.ScanBlockModeContinuous {

		err = counter.UpdateLatestBlock(nextBlock.BigInt())
		if err != nil {
			t.Fatal("UpdateLatestBlock error: ", err)
		}
	} else {
		err = counter.TriggerNewRescanEvent(nextBlock.BigInt())
		if err != nil {
			t.Fatal("TriggerNewRescanEvent error: ", err)
		}
	}

	time.Sleep(1 * time.Minute)
}

func TestGetEventProgress(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(mode, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}

	eventId := decimal.NewFromInt(5)
	progress, err := counter.GetEventProgress(eventId)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("from: %v, to: %v, current: %v\r\n", progress.StartNum, progress.EndNum, progress.CurrentNum)
}

func TestGetAllEvent(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeDiscret, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}

	eventList, err := counter.GetAllEventProgress()
	if err != nil {
		t.Fatal(err)
	}

	for _, eventProgress := range eventList {
		fmt.Printf("Event id: %v, start num: %v, end num: %v, current num: %v, status: %d\r\n",
			eventProgress.EventId, eventProgress.StartNum, eventProgress.EndNum, eventProgress.CurrentNum, eventProgress.Status)
	}
}

func TestCreateNewRescanBlock(t *testing.T) {
	var chainId uint64 = 2021
	var err error
	counter, err = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeDiscret, chainId, RescanBlockHandler, "ScanBlockChainCounter_Server", log.Name.Root)
	if err != nil {
		t.Fatal(err)
	}
	from := big.NewInt(10)
	to := big.NewInt(20)
	err = counter.CreateNewRescanInRange(from, to)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Minute)
}
