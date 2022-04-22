package test

import (
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/shopspring/decimal"
)

var lastTransId *big.Int

func TestApprovalWallet_QueryProxyImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	addr1, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Implementation", err)
	}
	fmt.Println("Implementation: ", addr1.String())

}
func TestApprovalWallet_QueryWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	writerList, err := approvalWallet.GetWriterList(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "GetWriterList", err)
	}

	for _, writer := range writerList {
		fmt.Println("Writer: " + writer.String())
	}
}
func TestApprovalWallet_AddWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	for _, assetName := range erc20Currency {
		erc20Addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Printf("Asset: %s: transOpt Error: %s\r\n ", assetName, err)
			continue
		}
		transOpt.GasLimit = 100000
		tx, err := approvalWallet.AddWriter(transOpt, erc20Addr)
		if err != nil {
			fmt.Printf("Asset: %s: Error: %s\r\n ", assetName, err.Error())
			continue
		} else {
			receipt, err := ethClient.QueryEthReceipt(tx)
			if err != nil {
				fmt.Printf("Query receipt error on asset: %s: Error: %s\r\n", assetName, err.Error())
				continue
			}
			receiptData, _ := json.Marshal(receipt)

			fmt.Printf("Asset: %s add write completed: %s\r\n", assetName, string(receiptData))
		}
	}

	// transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// tx1, err := approvalWallet.AddWriter(transOpt, common.HexToAddress(testOwnerAddr))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// queryEthReceipt(t, &ethClient, tx1)

}

func TestApprovalWallet_SubmitWithdraw(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")
	fmt.Println("USDT proxy addr: ", addr.Hex())
	erc20SC, err := contract.NewEurusERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	transOpt.GasLimit = 100000000
	amount := big.NewInt(1000000)

	totalAmount := big.NewInt(1340000)

	tx, err := erc20SC.SubmitWithdraw(transOpt, common.HexToAddress(testOwnerAddr), amount, totalAmount)
	if err != nil {
		fatalSmartContractError(t, "SubmitWithdraw", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

//func TestDeployWalletAddressMap(t *testing.T){
//	ethClient := initEthClient(t)
//	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
//	fmt.Println("Wallet Address Map proxy addr: ", addr.Hex())
//
//	walletAddressMapProxy, err := contract.NewOwnedUpgradeabilityProxy(addr,ethClient.Client)
//	if err != nil{
//		t.Fatal(err)
//	}
//
//	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
//
//	contract.DeployWalletAddressMap(transOpt, ethClient.Client)
//
//	walletAddressMapProxy.UpgradeTo()
//
//}

func TestApprovalWallet_QueryPendingWithdrawList(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWalletReader(addr, ethClient.Client)
	if err != nil {
		fatalSmartContractError(t, "NewApprovalWallet", err)
	}
	//transOpts, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	callOpts := bind.CallOpts{From: common.HexToAddress(testOwnerAddr)}

	pendingList, err := approvalWallet.GetPendingTransactionList(&callOpts)
	if err != nil {
		fatalSmartContractError(t, "GetPendingTransactionList", err)
	}
	// queryEthReceipt(t, &ethClient, pendingList)
	size := len(pendingList)

	var offset int = size
	if size > 0 {
		fmt.Println("buffer size: ", size)
		offset -= 32
		count := uint256.NewInt(0)
		count = count.SetBytes(pendingList[offset:size])

		fmt.Println("Count: ", count.String())
		for i := 0; i < int(count.Uint64()); i++ {
			transId := uint256.NewInt(0)
			transId = transId.SetBytes(pendingList[offset-32 : offset])
			offset -= 32
			lastTransId = transId.ToBig()
			fmt.Println("TransId: ", transId.String())

			srcAddr := uint256.NewInt(0)
			srcAddr = srcAddr.SetBytes(pendingList[offset-20 : offset])
			offset -= 20

			fmt.Println("Src addr:", srcAddr.Hex())

			destAddr := uint256.NewInt(0)
			destAddr = destAddr.SetBytes(pendingList[offset-20 : offset])
			offset -= 20

			fmt.Println("destAddr addr:", destAddr.Hex())

			strLen := uint256.NewInt(0)
			strLen = strLen.SetBytes(pendingList[offset-32 : offset])
			offset -= 32

			fmt.Println("String len: ", strLen)

			var assetName []byte
			byteLen := (int((strLen.Uint64())/32.0) + 1) * 32
			assetName = pendingList[offset-byteLen : offset]
			offset -= byteLen
			fmt.Println(assetName)
			fmt.Println(string(assetName))

			amount := uint256.NewInt(0)
			amount = amount.SetBytes(pendingList[offset-32 : offset])
			offset -= 32
			fmt.Println("Amount: ", amount.Uint64())

			timestamp := uint256.NewInt(0)
			timestamp = amount.SetBytes(pendingList[offset-32 : offset])
			offset -= 32
			fmt.Println("Timestamp: ", timestamp.Uint64())
		}
	}

	fmt.Println("Original bytes: ", pendingList)

}

func TestAprovalWallet_GetWalletOwner(t *testing.T) {
	fmt.Println("TestGetWalletOperatorList")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	a, err := approvalWallet.GetWalletOwner(&bind.CallOpts{})
	fmt.Println(a)

	//PrintWalletOwner(t, approvalWallet)
}

func TestAprovalWallet_SetWalletOwner(t *testing.T) {
	fmt.Println("TestGetWalletOperatorList")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	//platformAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
	}

	tx, err := approvalWallet.SetWalletOwner(transOpt, common.HexToAddress(testOwnerAddr))
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestApprovalWallet_GetWalletOperatorList(t *testing.T) {
	fmt.Println("TestGetWalletOperatorList")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	PrintWalletOperatorList(t, approvalWallet)

}

func TestApprovalWallet_GetOwners(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	PrintOwners(t, approvalWallet)
}

func TestApprovalWallet_AddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)

	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	for _, data := range approvalObserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println(err)
		}

		tx, err := approvalWallet.AddWalletOperator(transOpt, common.HexToAddress(data))
		if err != nil {
			fmt.Println(err)
		} else {
			receipt, err := ethClient.QueryEthReceipt(tx)
			if err != nil {
				fmt.Println("Query receipt error: ", err, " approval observer: ", data)
			} else {
				receiptJson, _ := json.Marshal(receipt)
				fmt.Println("Approval address: ", data, " receipt: ", string(receiptJson))
			}
		}

	}

}

func TestApprovalWallet_Approve(t *testing.T) {
	transId := big.NewInt(0)
	ApproveTransId(t, transId)
}

func ApproveTransId(t *testing.T, transId *big.Int) {
	fmt.Printf("TestApprove transId: %s\r\n", transId.String())
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testWithdrawApproverPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	transOpt.GasLimit = 1000000000
	tx, err := approvalWallet.ConfirmTransaction(transOpt, transId)
	if err != nil {
		fatalSmartContractError(t, "ConfirmTransaction", err)
	}

	queryEthReceipt(t, &ethClient, tx)
	curTime := time.Now()

	dbTrans := new(asset.WithdrawTransaction)
	dbTrans.Id = transId.Uint64() + 100
	dbTrans.Amount = decimal.NewFromInt32(1000000)
	dbTrans.AssetName = "USDT"
	dbTrans.ReviewedBy = strings.ToLower(testWithdrawApproverAddr)
	dbTrans.ReviewDate = &curTime
	dbTrans.RequestDate = curTime
	dbTrans.ApprovalWalletAddress = strings.ToLower(addr.Hex())
	reqTransId := transId.Uint64()
	dbTrans.RequestTransId = &reqTransId
	dbTrans.RequestTransHash = "Dummy"
	dbTrans.CustomerId = 1
	dbTrans.CustomerType = asset.CustomerUser
	dbTrans.Status = asset.StatusApproved
	dbTrans.InnetFromAddress = strings.ToLower(common.HexToAddress(testOwnerAddr).String())
	dbTrans.MainnetToAddress = strings.ToLower(common.HexToAddress("0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93").String())

	db := &database.Database{
		ReadOnlyDatabase: database.ReadOnlyDatabase{
			IP:         "18.163.252.201",
			Port:       9999,
			SchemaName: "public",
			DBName:     "postgres",
			UserName:   "admin",
			Password:   "kHSUOn/LW3vdKHIYX0tH8ovMuEHMhpj+VcoDHjd6xBM=",
		},
	}

	db.SetAESKey("XUFAKrxLKna5cZ2REBfFkii0btPBEehRApCbHPtQ6g8=")

	conn, err := db.GetConn()
	if err != nil {
		t.Fatal(err)
	}
	conn.Create(&dbTrans)
}

func TestApprovalWallet_GetInternalSmartContractConfig(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	internalAddr, err := approvalWallet.GetInternalSmartContractConfig(&bind.CallOpts{From: common.HexToAddress(testOwnerAddr)})
	if err != nil {
		t.Fatal("GetInternalSmartContractConfig: ", err)
	}

	fmt.Println("Internal Smart contract config address: ", internalAddr.Hex())
}

func TestApprovalWallet_SetInternalSmartContractConfig(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	getAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalAddr, err := approvalWallet.SetInternalSmartContractConfig(transOpt, getAddr)
	if err != nil {
		t.Fatal("GetInternalSmartContractConfig: ", err)
	}
	queryEthReceipt(t, &ethClient, internalAddr)
	//fmt.Println("Internal Smart contract config address: ", internalAddr.Hex())
}

func TestApprovalWallet_Upgrade(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	implAddr := getAddressBySmartContractName("ApprovalWallet")
	tx, err := proxy.UpgradeTo(transOpt, implAddr)
	if err != nil {
		t.Fatal("UpgradeTo error : ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestApprovalWallet_SetOwner(t *testing.T) {

	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}
	platformAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	tx, err := approvalWallet.AddOwner(transOpt, platformAddr)
	if err != nil {
		fmt.Print(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}
func TestApprovalWallet_ChangeRequirement(t *testing.T) {

	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")
	approvalWallet, err := contract.NewApprovalWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	tx, err := approvalWallet.ChangeRequirement(transOpt, big.NewInt(1))
	if err != nil {
		fmt.Print(err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

//func TestDeployUserWallet(t *testing.T) {
//	ethClient := initEthClient(t)
//
//	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	address, tx, userWallet, err := contract.DeployUserWallet(transOpt, ethClient.Client, big.NewInt(1))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(address)
//	fmt.Println(tx)
//	fmt.Println(userWallet)
//}

func TestGetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)

	addr := common.HexToAddress("0xed9d25411da94366e6d45bac84d652958b6f1ade")
	userWallet, err := contract.NewUserWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewApprovalWallet: ", err)
	}

	a, err := userWallet.GetWalletOwner(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(a.Hex())
}

func TestImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := common.HexToAddress("0xc50344f197abb4b204d8454d2dc2e99a4f4b79ab")

	p, e := contract.NewUpgradeabilityProxy(addr, ethClient.Client)

	if e != nil {
		t.Fatal("err: ", e)
	}

	a, e := p.Implementation(&bind.CallOpts{})
	if e != nil {
		t.Fatal("err: ", e)
	}
	fmt.Println(a.Hex())

	userWalletProxy, err := contract.NewUserWallet(addr, ethClient.Client)
	if e != nil {
		t.Fatal("get userWalletProxy err: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey("5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99", ethClient.ChainID)

	//transOpt, err := ethClient.GetNewTransactorFromSignServer("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJFdXJ1c0F1dGgiLCJuYmYiOjE2MTU3NzYxOTYsImNsaWVudEluZm8iOiJ7XCJDbGllbnRJbmZvXCI6XCJ7XFxcInNlcnZpY2VJZFxcXCI6IDkyLCBcXFwic2Vzc2lvbklkXFxcIjozMX1cIixcIlNlcnZpY2VJZFwiOjB9In0._4qs5uHPvjM7fRH4qwzEXUEebeoF-GtvCJXWDdr7ETI")
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	transOpt.GasLimit = 99999

	tx, err :=
		userWalletProxy.SetWalletOwner(transOpt, common.HexToAddress("0xc50344f197abb4b204d8454d2dc2e99a4f4b79ab"))
	if err != nil {
		t.Fatal("set wallet owner: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestApprovalWallet_SetFallbackAddress(t *testing.T) {
	ethClient := initEthClient(t)

	approvalWalletAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<ApprovalWallet>")

	approvalWallet, err := contract.NewApprovalWallet(approvalWalletAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	approvalReaderAddr := getAddressBySmartContractName("ApprovalWalletReader")
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx1, err := approvalWallet.SetFallbackAddress(transOpt, approvalReaderAddr)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx1)
}

func TestInit_ApprovalWawllet(t *testing.T) {
	TestAprovalWallet_SetWalletOwner(t)
	TestApprovalWallet_SetFallbackAddress(t)
	TestApprovalWallet_AddWriter(t)
	TestApprovalWallet_AddWalletOperator(t)
	TestApprovalWallet_SetInternalSmartContractConfig(t)
	TestApprovalWallet_ChangeRequirement(t)
}
