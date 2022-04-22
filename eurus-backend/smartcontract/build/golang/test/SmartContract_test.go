package test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	go_ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var environment string // Either LOCAL or DEV or TESTNET

var testOwnerPrivateKey string
var testOwnerAddr string

var testUserWalletSCOwnerAddr string
var userWalletSCOwnerPrivateKey string

var invokerAddrList []string

var chainId string
var testNetChainId string

var testWithdrawApproverPrivateKey string
var testWithdrawApproverAddr string

var testTestNetPlatformWalletOwnerPrivateKey string
var testTestNetPlatformWalletOwnerAddr string
var testWithdrawObserverPrivateKey string
var testWithdrawObserverAddr string

var testTestNetOwnerPrivateKey string
var testTestNetOwnerAddr string

// var testDepositObserverAddr string

var approvalObserverHDWalletAddr []string
var approvalObserverAddr []string
var depositObserverAddr []string
var sweepServiceAddr []string
var withdrawobserverAddr []string
var erc20Currency []string
var userServerHDWalletAddr []string

var configServerHDWalletAddr []string

var blockChainIndexerHDWalletAddr []string

var addressMap map[string]common.Address = make(map[string]common.Address)
var testNetAddressMap map[string]common.Address = make(map[string]common.Address)

var kycServerAddr []string

var userObserverAddr []string

var signServerAddr string

var gasFeeCollectWalletAddr string

var sweepInvokerAddr string

var mainnetEthClientIP string
var mainnetEthClientPort int
var mainnetEthClientProtocol string

var ethClientIP string
var ethClientPort int
var ethClientProtocol string

var smartContractFileName string

func init() {

	InitEnvironment()

	erc20Currency = []string{
		"USDT",
		"USDC",
		"LINK",
		"UNI",
		"BNB",
		"BUSD",
		"YFI",
		"DAI",
		"OMG",
		"VEN",
		"AAVE",
		"HT",
		"SUSHI",
		"TUSD",
		"cDAI",
		"SXP",
		"BAT",
		"USDK",
		"WBTC",
		"ZIL",
		"SNX",
		"OKB",
		"BAND",
		"MKR",
		"HUSD",
		"ZRX",
		"PAX",
		"COMP",
		"RSR",
		"BAL",
		"ETH",
		"PLA",
		"USDM",
		"BTCM",
		"ETHM",
		"MST",
	}

	file, err := os.OpenFile("../../../"+smartContractFileName, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Unable to load smart contract address JSON ", err.Error())
		return
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Unable to read file: ", err.Error())
		return
	}
	var jsonMap map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal(rawData, &jsonMap)
	if err != nil {
		fmt.Println("Invalid JSON format: ", err.Error())
		return
	}

	if smartContractObj, ok := jsonMap[chainId]; ok {
		smartContractMap := smartContractObj.(map[string]interface{})

		smartContractInnerMap := smartContractMap["smartContract"].(map[string]interface{})
		for key, value := range smartContractInnerMap {

			child := value.(map[string]interface{})
			if intf, ok := child["address"]; ok {
				if addr, ok := intf.(string); ok {
					addrObj := common.HexToAddress(addr)
					addressMap[key] = addrObj
				} else {
					fmt.Printf("%s address invalid\r\n", key)
				}
			} else {
				fmt.Printf("%s address not found\r\n", key)
			}
		}
	}

	if testNetSmartContractObj, ok := jsonMap[testNetChainId]; ok {
		testNetSmartContractMap := testNetSmartContractObj.(map[string]interface{})

		testNetSmartContractInnerMap := testNetSmartContractMap["smartContract"].(map[string]interface{})
		for key, value := range testNetSmartContractInnerMap {

			child := value.(map[string]interface{})
			if intf, ok := child["address"]; ok {
				if addr, ok := intf.(string); ok {
					addrObj := common.HexToAddress(addr)
					testNetAddressMap[key] = addrObj
				} else {
					fmt.Printf("%s address invalid\r\n", key)
				}
			} else {
				fmt.Printf("%s address not found\r\n", key)
			}
		}
	}
}

func TestInitDevData(t *testing.T) {

	fmt.Println("Init Smart Contract")
	TestInitInternalSC(t)
	TestInit_ExternalSC(t)
	TestInitERC20(t)
	TestInitEurusERC20(t)
	TestEurusInternalSC_Init(t)
	TestInitWithdrawSC(t)
	TestInitPlatformWallet(t)
	TestInitAdminFeeWallet(t)
	TestInit_ApprovalWawllet(t)

	if environment != "TESTNET" && environment != "MAINNET" {
		//Mint USDT to owner
		TestMint(t)
	}
	TestWalletAddressMap_AddWriter(t)
	//Adding owner to wallet address map
	//TestGetWalletOperatorList(t)
	//TestSetWalletOwnerToWithdrawSmartContract(t)

	TestEurusInternalConfig_SetEurusUserDepositAddress(t) //Set EurusUserDeposit address to EurusInternalConfig smart contract
	TestEurusUserDeposit_SetEurusInternalConfigAddress(t) //Set EurusInternalConfig address to EurusUserDeposit smart contract
	TestEurusUserDeposit_AddWriter(t)                     //Need to manual set the writer address before run
	TestEurusUserDeposit_SetEurusPlatformAddress(t)
	TestEurusPlatformWallet_SetEurusInternalConfigAddress(t)

	TestTransferFundingToWithdrawObserver(t)
	TestTransferFundingToDepositObserver(t)
	TestTransferFundingToApprovalObserver(t)
	TestTransferFundingToUserServer(t)
	TestTransferFundingToConfigServer(t)
	TestTransferFundingToKYCServer(t)
	TestTransferFundingToBlockChainIndexer(t)
	TestTransferRinkebyFundingToWithdrawServer(t)
	TestEurusUserDeposit_SetEtherForwardAddress(t)
	TestInit_MarketRegWallet(t)
	//Funding to Sign user wallet owner
	//Funding to Sign user invoker wallet
	//Funding to Sweep service invoker wallet

}

func TestFullWithdrawLogicFlow(t *testing.T) {
	fmt.Println("TestSubmitWithdraw")
	TestApprovalWallet_SubmitWithdraw(t)

	fmt.Println("TestQueryPendingWithdrawList")
	TestApprovalWallet_QueryPendingWithdrawList(t)
	fmt.Println("TestApproveTransId")
	ApproveTransId(t, lastTransId)
}

func TestTransferFundingToWithdrawObserver(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, addr := range withdrawobserverAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, addr, amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Tx: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, -1)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for withdraw observer: ", addr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Withdraw observer address: ", addr, " receipt: ", string(receiptJson))
		}

		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		fmt.Println("Withdraw observer: ", addr, "Balance: ", balance.String())
	}

}

func TestTransferFundingToDepositObserver(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, addr := range depositObserverAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, addr, amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for deposit observer: ", addr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("deposit observer address: ", addr, " receipt: ", string(receiptJson))
		}

		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(addr)
		fmt.Println(balance.String())
	}
	//queryEthReceipt(t, &ethClient, tx)
	// balance, err := ethClient.GetBalance(common.HexToAddress(testDepositObserverAddr))
	// if err != nil {
	// 	t.Fatal("GetBalance: ", err)
	// }
	// fmt.Println("Balance: ", balance.String())
}

func TestTransferFundingToApprovalObserver(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, addr := range approvalObserverAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, addr, amount)
		if err != nil {
			fmt.Println(err)
			continue
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for deposit observer: ", addr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("deposit observer address: ", addr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Approval observer: ", addr, " balance: ", balance.String())

	}

}

func TestTransferFundingToUserServer(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, addr := range userServerHDWalletAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, addr, amount)
		if err != nil {
			fmt.Println(err)
			continue
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for deposit observer: ", addr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("deposit observer address: ", addr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Approval observer: ", addr, " balance: ", balance.String())

	}

}

func TestTransferFundingToUserObserver(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, addr := range userObserverAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, addr, amount)
		if err != nil {
			fmt.Println(err)
			continue
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for deposit observer: ", addr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("deposit observer address: ", addr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(addr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Approval observer: ", addr, " balance: ", balance.String())

	}

}

func TestTransferFundingToConfigServer(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}
	for _, configAddr := range configServerHDWalletAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, configAddr, amount)
		if err != nil {
			t.Fatal(err)

		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for config server: ", configAddr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Config server address: ", configAddr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(configAddr))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Config server: ", configAddr, " balance: ", balance.String())
	}

}

func TestTransferRinkebyFundingToWithdrawServer(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("140000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, withdrawObsAddr := range withdrawobserverAddr {
		_, tx, err := ethClient.TransferETH(testTestNetOwnerPrivateKey, withdrawObsAddr, amount)
		if err != nil {
			fmt.Println("Transfer ETH failed: ", err, " address: ", withdrawObsAddr)
			continue
		}
		fmt.Println("Trans hash: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, -1)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for withdraw observer: ", withdrawObsAddr)
			continue
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("withdraw observer address: ", withdrawObsAddr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(withdrawObsAddr))
		if err != nil {
			fmt.Println("Unable to get balance: ", err)
			continue
		}

		fmt.Println("withdraw observer: ", withdrawObsAddr, " balance: ", balance.String())
	}
}

func TestTransferFundingToKYCServer(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}

	for _, kycAddr := range kycServerAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, kycAddr, amount)
		if err != nil {
			t.Fatal(err)
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for KYC server: ", kycAddr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("KYC server address: ", kycAddr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(kycAddr))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("KYC server: ", kycAddr, " balance: ", balance.String())
	}

}

func TestTransferFundingToBlockChainIndexer(t *testing.T) {

	ethClient := initEthClient(t)
	amount := big.NewInt(0)
	var ok bool
	amount, ok = amount.SetString("1000000000000000000", 10)
	if !ok {
		t.Fatal("Create amount failed")
	}
	for _, blockChainIndexerAddr := range blockChainIndexerHDWalletAddr {
		_, tx, err := ethClient.TransferETH(testOwnerPrivateKey, blockChainIndexerAddr, amount)
		if err != nil {
			t.Fatal(err)

		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Error query receipt ", err, " for block chain indexer: ", blockChainIndexerAddr)
		} else {
			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Block chain indexer address: ", blockChainIndexerAddr, " receipt: ", string(receiptJson))
		}
		balance, err := ethClient.GetBalance(common.HexToAddress(blockChainIndexerAddr))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Block chain indexer: ", blockChainIndexerAddr, " balance: ", balance.String())
	}

}

//Depreciated
// func TestTransferFundingToUserServerHDWalletAddr(t *testing.T) {

// 	ethClient := initEthClient(t)
// 	amount := big.NewInt(0)
// 	var ok bool
// 	amount, ok = amount.SetString("100000000000000000", 10)
// 	if !ok {
// 		t.Fatal("Create amount failed")
// 	}

// 	for _, data := range userServerHDWalletAddr {
// 		_, tx, err := ethClient.TransferETH(testWithdrawObserverPrivateKey, data, amount)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		receipt, err := ethClient.QueryEthReceipt(tx)

// 		if err != nil {
// 			fmt.Println("Failed to transfer EUN to user server address: ", data)
// 			continue
// 		} else {
// 			receiptJson, _ := receipt.MarshalJSON()
// 			fmt.Println("Transfer funding to user server address: ", data, " receipt: ", string(receiptJson))
// 		}

// 		balance, err := ethClient.GetBalance(common.HexToAddress(data))
// 		if err != nil {
// 			fmt.Println("Get balance error on user server address: ", data, " error: ", err)
// 			continue
// 		}
// 		fmt.Println(data)
// 		fmt.Println(balance.String())
// 	}
// }

/*
func TestRevertMessage(t *testing.T) {
	ethClient := initEthClient(t)

	unitTestSC, _ := contract.NewUnitTest(common.HexToAddress("0xdCf90b79D55d6F12aB268993B7B83e6809fFDA04"), ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	//transOpt.GasLimit = 4700000000
	// abiMetaData, _ := abi.JSON(strings.NewReader(string(contract.UnitTestABI)))

	// inputData, err := abiMetaData.Methods["SetValue"].Inputs.Pack(big.NewInt(1))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// inputData = append(abiMetaData.Methods["SetValue"].ID, inputData...)

	// fromAddr := common.HexToAddress(testOwnerAddr)
	// toAddr := common.HexToAddress("0xdCf90b79D55d6F12aB268993B7B83e6809fFDA04")
	// msg := go_ethereum.CallMsg{
	// 	From:     fromAddr,
	// 	To:       &toAddr,
	// 	GasPrice: big.NewInt(13570),
	// 	Data:     inputData,
	// }
	// gasLimit, err := ethClient.Client.EstimateGas(context.Background(), msg)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// transOpt.GasLimit = gasLimit
	tx, err := unitTestSC.SetValue(transOpt, big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}
*/

func TestOwnedUpgradeabilityProxy(t *testing.T) {

	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDC>")
	sc, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)

	if err != nil {
		t.Fatal(err)
	}

	addr, err = sc.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Implementation address: ", addr.Hex())

}

func getAddressBySmartContractName(name string) common.Address {
	if addr, ok := addressMap[name]; ok {
		return addr
	}
	fmt.Printf("Unable to find smart contract name: %s\r\n", name)
	return common.HexToAddress("0x00")
}

func getTestNetAddressBySmartContractName(name string) common.Address {
	if addr, ok := testNetAddressMap[name]; ok {
		fmt.Println(name+" address: ", addr.Hex())
		return addr
	}
	fmt.Printf("Unable to find smart contract name: %s\r\n", name)
	return common.HexToAddress("0x00")

}

func TestUSDT(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	fmt.Println("Proxy addr: ", addr.String())
	sc, err := contract.NewERC20(addr,
		ethClient.Client)

	if err != nil {
		t.Fatal(err)
	}

	// transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	// amount := big.NewInt(1234567876543567)
	// tx, err := sc.Init(transOpt, "USDT", "USDT", amount, 6)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// queryEthReceipt(t, &ethClient, tx)

	arrList, err := sc.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, addr := range arrList {
		fmt.Println(addr.Hex())
	}

	totalSupply, err := sc.TotalSupply(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	dec, err := sc.Decimals(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	sym, err := sc.Symbol(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Total supply: %s, Decimal: %d, Symbol: %s\r\n", totalSupply.String(), dec, sym)
}

func TestMint(t *testing.T) {
	fmt.Println("TestMint")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	usdt, err := contract.NewEurusERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(9000000000000000000)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := usdt.Mint(transOpt, common.HexToAddress(testOwnerAddr), amount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Trans Hash: ", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestGetUSDTOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	usdt, err := contract.NewERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	addresses, err := usdt.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Owner count: ", len(addresses))
	for _, addr := range addresses {
		fmt.Println(addr.Hex())
	}
}

func TestGetUSDTBalance(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	usdt, err := contract.NewERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := usdt.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Balance is: %s\r\n", balance.String())
}

func TestGetEUNBallance(t *testing.T) {
	ethClient := initEthClient(t)

	balance, err := ethClient.GetBalance(common.HexToAddress("0x7909f2d6508c1fc08e11e06726831133230abc7f"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Balance: ", balance.String())
}

func TestUSDTTransfer(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	usdt, err := contract.NewERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(5000000000)

	tx, err := usdt.Transfer(transOpt, common.HexToAddress("0x57106f35330a8532647170dB272C8498bFE35132"), amount)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestETHTransfer(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ETH>")
	eth, err := contract.NewERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(10000000)

	tx, err := eth.Transfer(transOpt, common.HexToAddress("0x6a44C65064067B21BE8F67f8461536a4b51f1C03"), amount)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestGetUsdtAddrToTestWallet(t *testing.T) {
	ethClient := initEthClient(t)
	testWallet, err := contract.NewTestWallet(common.HexToAddress("0x81fA2f6A452cD25a95Ced161888215D55297ef38"), ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	// transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	addr, err := testWallet.UsdtAddress(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Usdt address: ", addr.Hex())
}

func TestTransferFromTestWallet(t *testing.T) {
	ethClient := initEthClient(t)
	testWallet, err := contract.NewTestWallet(common.HexToAddress("0x81fA2f6A452cD25a95Ced161888215D55297ef38"), ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(900000)
	tx, err := testWallet.Confirm(transOpt, common.HexToAddress(testOwnerAddr), amount)
	if err != nil {
		fatalSmartContractError(t, "Confirm", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func initEthClient(t *testing.T) ethereum.EthClient {
	//DEV
	chainIdNum, err := strconv.ParseInt(chainId, 10, 32)
	if err != nil {
		t.Fatal("Invalid chain id: " + err.Error())
	}
	var ethClient ethereum.EthClient
	ethClient = ethereum.EthClient{Protocol: ethClientProtocol, IP: ethClientIP, Port: ethClientPort, ChainID: big.NewInt(chainIdNum)}
	// if environment == "DEV" {
	// 	ethClient = ethereum.EthClient{Protocol: "http", IP: "13.228.169.25", Port: 10002, ChainID: big.NewInt(chainIdNum)}
	// } else if environment == "STAGING" {
	// 	//Local
	// 	ethClient = ethereum.EthClient{Protocol: "http", IP: "13.228.80.104", Port: 8545, ChainID: big.NewInt(chainIdNum)}
	// } else if environment == "LOCAL" {
	// 	//Local
	// 	ethClient = ethereum.EthClient{Protocol: "http", IP: "127.0.0.1", Port: 9545, ChainID: big.NewInt(chainIdNum)}
	// } else if environment == "TESTNET" {
	// 	ethClient = ethereum.EthClient{Protocol: "https", IP: "testnet.eurus.network", Port: 443, ChainID: big.NewInt(chainIdNum)}
	// } else if environment == "RINKEBY" {
	// 	ethClient = ethereum.EthClient{Protocol: "http", IP: "13.212.253.186", Port: 8545, ChainID: big.NewInt(chainIdNum)}
	// } else if environment == "PROD" {
	// 	ethClient = ethereum.EthClient{Protocol: "http", IP: "13.213.109.88", Port: 8545, ChainID: big.NewInt(chainIdNum)}
	// }
	_, err = ethClient.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return ethClient
}

func initTestNetEthClient(t *testing.T) ethereum.EthClient {

	testNetChainIdNum, err := strconv.ParseInt(testNetChainId, 10, 32)
	if err != nil {
		t.Fatal("Invalid chain id: " + err.Error())
	}

	var ethClient ethereum.EthClient = ethereum.EthClient{
		Protocol: mainnetEthClientProtocol,
		IP:       mainnetEthClientIP,
		Port:     mainnetEthClientPort,
		ChainID:  big.NewInt(testNetChainIdNum),
	}

	_, err = ethClient.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return ethClient
}

func queryEthReceipt(t *testing.T, ethClient *ethereum.EthClient, tx *types.Transaction) {
	var receipt *ethereum.BesuReceipt
	var err error
	for {
		receipt, err = ethClient.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			if err == go_ethereum.NotFound {
				time.Sleep(time.Second)
				continue
			} else {
				log.Fatal(err)
			}
		} else {
			break
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction hash: %s\r\nStatus: %d\r\n", tx.Hash().Hex(), receipt.Status)

	receiptJson, _ := json.Marshal(receipt)
	fmt.Printf("Receipt content: %s\r\n", string(receiptJson))
}

func fatalSmartContractError(t *testing.T, message string, err error) {
	jsonData, _ := json.Marshal(err)
	t.Fatal(message, ": ", err, "\r\nJSON: ", string(jsonData))
}

type ParentConfig struct {
	FieldOne int `json:"fieldOne"`
}

func (me ParentConfig) GetFieldOne() int {
	return me.FieldOne
}

type ChildConfig struct {
	ParentConfig
	FieldTwo int `json:"fieldTwo"`
}

func (me ChildConfig) GetFieldOne() int {
	return 18
}

type IConfig interface {
	GetFieldOne() int
}

func TestLoadConfig(t *testing.T) {
	var config *ChildConfig = new(ChildConfig)
	var configIntf IConfig = config
	Load(configIntf)
	fmt.Println("GetFieldOne: ", config.GetFieldOne())
}

func Load(configIntf IConfig) {
	var configStr string = `{"fieldOne":1, "fieldTwo":2}`
	err := json.Unmarshal([]byte(configStr), configIntf)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println("GetFieldOne: ", configIntf.GetFieldOne())
}

func TestSHA256(t *testing.T) {
	sum := sha256.Sum256([]byte("hello world\n"))
	fmt.Printf("%x", sum)
}

func TestGetDummyRequirment(t *testing.T) {
	ethClient := initEthClient(t)

	sc, _ := contract.NewDummyPlatformWallet(common.HexToAddress("0x0216e1e274250128FE9021a5702b260cAf9767f8"), ethClient.Client)
	req, _ := sc.Required(&bind.CallOpts{})
	fmt.Println("Required: ", req.String())
}

func TestRawTrans(t *testing.T) {
	// rawData, err := hex.DecodeString("f9015381b38302db5b830186a0941c874a3ce03f373e8b18196ff3bc2ae3fd1557a1880de0b6b3a7640000b8e42b92c6b90000000000000000000000000000000000000000000000000000000000017ac6000000000000000000000000000000000000000000000000000000000000008000000000000000000000000044f426bc9ac7a83521ea140aeb70523c0a85945a00000000000000000000000000000000000000000000000000000000613ae1a1000000000000000000000000000000000000000000000000000000000000000200000000000000000000000017deba6e45745d6b72f684a150566001283e4424000000000000000000000000a54dee79c3bb34251debf86c1ba7d21898ffb7ac820feea0d885a407d0581743890ec1d315c9c7e8be989663025700214d4ed03bb979a01ba0745ec0f7785cc18cd55303769d12f248834dab1a100e3ab5be1c759f622ea154")
	rawData, err := hex.DecodeString("f901b082027685010c388d008404d3f6409488b7c329fb4531592b18664d2ab02b0f2903bae980b901443c4cd911000000000000000000000000b13070bd35176fae45ff692d5279495547473974000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000004414243440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000411af0dc5cf2f1cda05eaf4df03ce1dcb46f3e3fbb7f5c3996bd7b7a757a1f29ed244f1cde067a0ec118f5cbaacb7d9d16367d7077221aa40bc3c97cdf251a10340000000000000000000000000000000000000000000000000000000000000000820feda0b5eb1ff2dd6e57d5733d1c3f09034472f653389a67d9a9d4bf53466ef7e3ded9a060d2d1eb0125c44e8684df0847ea1d4ccee9a2ad5709cebead55a86ae3a3429d")
	if err != nil {
		t.Fatal(err)
	}

	stream := rlp.NewStream(bytes.NewReader(rawData), 0)
	tx := new(types.Transaction)
	err = tx.DecodeRLP(stream)
	if err != nil {
		t.Fatal(err)
	}

	txData := hex.EncodeToString(tx.Data())

	v, r, s := tx.RawSignatureValues()
	fmt.Println(txData)
	fmt.Printf("v: %s, r: %s, s: %s\r\n", v, r, s)
}
