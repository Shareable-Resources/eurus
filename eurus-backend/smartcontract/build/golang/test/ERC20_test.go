package test

import (
	"encoding/hex"
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const transferTargetAddress = "0x4dfb6d6790054f3eb68324bc230e3104137ca8db"

func InitERC20_Rinkeby(t *testing.T, assetName string, supply *big.Int, decimals uint8) {
	ethClient := initTestNetEthClient(t)
	erc20Instance, err := contract.NewTestERC20(getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<"+assetName+">"), ethClient.Client)
	if err != nil {
		t.Fatal("New ERRC20 error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}

	tx, err := erc20Instance.Init(transOpt, assetName, assetName, supply, decimals)
	if err != nil {
		fatalSmartContractError(t, "Init Error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestERC20_InitRinkeby(t *testing.T) {
	amount := big.NewInt(0)

	amount = amount.SetUint64(0)

	// InitERC20_Rinkeby(t, "PLA", amount, 6)

	InitERC20_Rinkeby(t, "USDT", amount, 6)

	// amount = amount.SetUint64(15310845932149550)
	// InitERC20_Rinkeby(t, "USDC", amount, 6)

	// amount = amount.SetUint64(uint64(1000000000.0 * math.Pow10(18)))
	// InitERC20_Rinkeby(t, "LINK", amount, 18)

	// amount = amount.SetUint64(uint64(1000000000 * math.Pow10(18)))
	// InitERC20_Rinkeby(t, "UNI", amount, 18)

	// var isSuccess bool
	// amount, isSuccess = amount.SetString("16579517055253348798759097", 10)
	// if !isSuccess {
	// 	t.Fatal("BNB max supply failed")
	// }
	// InitERC20_Rinkeby(t, "BNB", amount, 18)

	// amount, isSuccess = amount.SetString("7980917476250000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("BUSD max supply failed")
	// }
	// InitERC20_Rinkeby(t, "BUSD", amount, 18)

	// amount, isSuccess = amount.SetString("36666000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("YFI max supply failed")
	// }
	// InitERC20_Rinkeby(t, "YFI", amount, 18)

	// amount, isSuccess = amount.SetString("4582063439963584546011659849", 10)
	// if !isSuccess {
	// 	t.Fatal("DAI max supply failed")
	// }
	// InitERC20_Rinkeby(t, "DAI", amount, 18)

	// amount, isSuccess = amount.SetString("140245398245132780789239631", 10)
	// if !isSuccess {
	// 	t.Fatal("OMG max supply failed")
	// }
	// InitERC20_Rinkeby(t, "OMG", amount, 18)

	// amount, isSuccess = amount.SetString("1000000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("VEN max supply failed")
	// }
	// InitERC20_Rinkeby(t, "VEN", amount, 18)

	// amount, isSuccess = amount.SetString("16000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("AAVE max supply failed")
	// }
	// InitERC20_Rinkeby(t, "AAVE", amount, 18)

	// amount, isSuccess = amount.SetString("500000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("HT max supply failed")
	// }
	// InitERC20_Rinkeby(t, "HT", amount, 18)

	// amount, isSuccess = amount.SetString("219608746618504864068370941", 10)
	// if !isSuccess {
	// 	t.Fatal("SUSHI max supply failed")
	// }
	// InitERC20_Rinkeby(t, "SUSHI", amount, 18)

	// amount, isSuccess = amount.SetString("1139038697380000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("TUSD max supply failed")
	// }
	// InitERC20_Rinkeby(t, "TUSD", amount, 18)

	// amount, isSuccess = amount.SetString("23294018011305724555", 10)
	// if !isSuccess {
	// 	t.Fatal("cDAI max supply failed")
	// }
	// InitERC20_Rinkeby(t, "cDAI", amount, 8)

	// amount, isSuccess = amount.SetString("285368788739134951000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("SXP max supply failed")
	// }
	// InitERC20_Rinkeby(t, "SXP", amount, 18)

	// amount, isSuccess = amount.SetString("1500000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("BAT max supply failed")
	// }
	// InitERC20_Rinkeby(t, "BAT", amount, 18)

	// amount, isSuccess = amount.SetString("32478711000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("USDK max supply failed")
	// }
	// InitERC20_Rinkeby(t, "USDK", amount, 18)

	// amount, isSuccess = amount.SetString("213272942183558446345", 10)
	// if !isSuccess {
	// 	t.Fatal("WBTC max supply failed")
	// }
	// InitERC20_Rinkeby(t, "WBTC", amount, 8)

	// amount, isSuccess = amount.SetString("213272942183558446345", 10)
	// if !isSuccess {
	// 	t.Fatal("ZIL max supply failed")
	// }
	// InitERC20_Rinkeby(t, "ZIL", amount, 12)

	// amount, isSuccess = amount.SetString("225834728820530628030769135", 10)
	// if !isSuccess {
	// 	t.Fatal("SNX max supply failed")
	// }
	// InitERC20_Rinkeby(t, "SNX", amount, 18)

	// amount, isSuccess = amount.SetString("300000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("OKB max supply failed")
	// }
	// InitERC20_Rinkeby(t, "OKB", amount, 18)

	// amount, isSuccess = amount.SetString("100000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("BAND max supply failed")
	// }
	// InitERC20_Rinkeby(t, "BAND", amount, 18)

	// amount, isSuccess = amount.SetString("79542198807227860", 10)
	// if !isSuccess {
	// 	t.Fatal("HUSD max supply failed")
	// }
	// InitERC20_Rinkeby(t, "HUSD", amount, 8)

	// amount, isSuccess = amount.SetString("994709090054240387057335", 10)
	// if !isSuccess {
	// 	t.Fatal("MKR max supply failed")
	// }
	// InitERC20_Rinkeby(t, "MKR", amount, 18)

	// amount, isSuccess = amount.SetString("1000000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("ZRX max supply failed")
	// }
	// InitERC20_Rinkeby(t, "ZRX", amount, 18)

	// amount, isSuccess = amount.SetString("1177158554520000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("PAX max supply failed")
	// }
	// InitERC20_Rinkeby(t, "PAX", amount, 18)

	// amount, isSuccess = amount.SetString("10000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("COMP max supply failed")
	// }
	// InitERC20_Rinkeby(t, "COMP", amount, 18)

	// amount, isSuccess = amount.SetString("100000000000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("RSR max supply failed")
	// }
	// InitERC20_Rinkeby(t, "RSR", amount, 18)

	// amount, isSuccess = amount.SetString("42105000000000000000000000", 10)
	// if !isSuccess {
	// 	t.Fatal("BAL max supply failed")
	// }
	// InitERC20_Rinkeby(t, "BAL", amount, 18)

}

func TestERC20_GetWriterList(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	for _, currency := range erc20Currency {
		addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + currency + ">")
		fmt.Println(currency, " proxy addr: ", addr.Hex())
		curr, err := contract.NewTestERC20(addr, ethClient.Client)
		if err != nil {
			fmt.Println("err: ", err, " currency: ", currency)
			continue
		}
		writerList, err := curr.GetWriterList(&bind.CallOpts{})
		if err != nil {
			fmt.Println("err: ", err, " currency: ", currency)
			continue
		}
		fmt.Println("Currency: ", currency)
		for _, writer := range writerList {
			fmt.Println("Writer: ", writer.Hex())
		}
	}
}

func MintERC20CoinAtMainnet(t *testing.T, ethClient ethereum.EthClient, erc20AssetName string) {

}

func TestERC20_MintMainnet(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	erc20AssetName := "PLA"
	erc20Instance, err := contract.NewTestERC20(getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<"+erc20AssetName+">"), ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	mintTargetAddressList := []string{"0x5e25326C287aa5617c15fd22cfE442bAA0b66113"}
	for _, account := range mintTargetAddressList {

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal(err)
		}

		transOpt.GasLimit = 100000
		amount := big.NewInt(0)
		amount, _ = amount.SetString("5000000000000", 10)

		tx, err := erc20Instance.Mint(transOpt, common.HexToAddress(account), amount)
		if err != nil {
			fmt.Println("Mint ERC20 to ", account, " error: ", err)
			continue
		}
		fmt.Println("Tx Hash: ", tx.Hash().Hex())
		queryEthReceipt(t, &ethClient, tx)
	}
}

func TestERC20_GetDeciamls(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<BTCM>")
	fmt.Println("Address: ", addr.String())
	erc20Instance, err := contract.NewTestERC20(addr, ethClient.Client)
	if err != nil {
		fmt.Printf("Unable to NewTestERC20 for asset %s. Error: %s\r\n", "OwnedUpgradeabilityProxy<PLA>", err.Error())
		return
	}
	decimal, err := erc20Instance.Decimals(&bind.CallOpts{})
	if err != nil {
		fmt.Println("error: ", err.Error())
	} else {
		fmt.Println("decimal: ", decimal)
	}

	sym, err := erc20Instance.Symbol(&bind.CallOpts{})
	if err != nil {
		fmt.Println("get symbol error: ", err.Error())
	} else {
		fmt.Println("symbol: ", sym)
	}
}

func TestERC20_AddWriter(t *testing.T) {
	testnetEthClient := initTestNetEthClient(t)
	for _, ERCName := range erc20Currency {
		if ERCName == "ETH" {
			continue
		}
		addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + ERCName + ">")
		erc20, err := contract.NewTestERC20(addr, testnetEthClient.Client)
		if err != nil {
			t.Fatal("NewApprovalWallet: ", err)
		}
		for _, userServerAddr := range userServerHDWalletAddr {
			transOpt, err := testnetEthClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, testnetEthClient.ChainID)
			if err != nil {
				fmt.Println("GetNewTransactorFromPrivateKey error :", err)
				continue
			}

			tx, err := erc20.AddWriter(transOpt, common.HexToAddress(userServerAddr))
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Tx hash: ", tx.Hash().Hex())
			receipt, err := testnetEthClient.QueryEthReceiptWithSetting(tx, 1, 20)
			if err != nil {
				fmt.Println(err)
				continue
			}
			data, _ := json.Marshal(receipt)
			fmt.Println(string(data))
		}
	}
}

func TestERC20_DecodeABIInputArg(t *testing.T) {

	ethereum.DefaultABIDecoder.ImportABIJson("EurusERC20", contract.EurusERC20ABI)
	data, _ := hex.DecodeString("a9059cbb0000000000000000000000008157d3d61ec41ec890611ffdddf8341591a52a4d00000000000000000000000000000000000000000000000000000000001e8480")
	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(data, "EurusERC20", "transfer")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(args)
}

func TestERC20_GetBalanceMainnet(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<BTCM>")
	testErc20, err := contract.NewTestERC20(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	// platformWallet := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	balance, err := testErc20.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0x809eecf74ee563d819ccb9761410e3541f18e747"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Balance: ", balance.String())
}

//Demo DApp smart contract function
func TestERC20_InitDAppSampleTokenSideChain(t *testing.T) {
	ethClient := initEthClient(t)
	var assetName string = "DAppSampleToken"
	dapp, err := contract.NewERC20(getAddressBySmartContractName("OwnedUpgradeabilityProxy<"+assetName+">"), ethClient.Client)
	if err != nil {
		t.Fatal("New ERRC20 error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	tx, err := dapp.Init(transOpt, assetName, assetName, big.NewInt(0), 6)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

	usdtAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>")

	transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey 2 : ", err)
	}

	tx1, err := dapp.AddWriter(transOpt1, usdtAddr)
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx1)
}

func TestInitERC20(t *testing.T) {
	// TestERC20_SetInternalSCConfigAddress(t)
	//TestMintERC20_Sidechain(t)
}
