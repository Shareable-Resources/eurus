package main

import (
	"bufio"
	"encoding/json"
	"eurus-backend/data_patch/patch"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	go_ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
)

type Patch2Config struct {
	server.ServerConfigBase
	GasLimit uint64 `json:"gasLimit"`
}

func main() {

	patch.InitLog()
	err := patch.LoadSmartContractConfig("")
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln(err)
	}

	_ = ethereum.DefaultABIDecoder.ImportABIJson("GeneralMultiSigWallet", mainnet_contract.GeneralMultiSigWalletABI)

	config := new(Patch2Config)
	err = patch.LoadConfig("Patch2Config.json", config)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("LoadConfig error: ", err)
	}
	var option string
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Option: ")
		fmt.Println("1: Submit add currency transaction")
		fmt.Println("2: Confirm submitted transaction")
		fmt.Print("Input an option: ")
		option, _ = reader.ReadString('\n')
		option = strings.TrimSuffix(option, "\n")
		if option != "1" && option != "2" {
			continue
		}
		break
	}

	ethClient, err := patch.CreateEurusEthClient(&config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("CreateEurusEthClient error: ", err)
	}

	switch option {
	case "1":
		SubmitTransaction(config, ethClient)
	case "2":
		ConfirmTransaction(config, ethClient)
	}

}

func SubmitTransaction(config *Patch2Config, ethClient *ethereum.EthClient) {
	log.GetLogger(log.Name.Root).Infoln("Going to submit transaction")

	configAddr := patch.GetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>", ethClient.ChainID.Int64())

	priKeyStr := ReadPrivateKey()

	configSC, err := mainnet_contract.NewEurusInternalConfig(configAddr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("NewEurusInternalConfig error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("GetNewTransactorFromPrivateKey error: ", err)
	}

	transOpt.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return tx, nil
	}
	transOpt.NoSend = true
	transOpt.GasLimit = config.GasLimit
	//Pack function data
	mstAddr := patch.GetAddressBySmartContractName("OwnedUpgradeabilityProxy<MST>", ethClient.ChainID.Int64())
	tx, err := configSC.AddCurrencyInfo(transOpt, mstAddr, "MST")
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("Create AddCurrencyInfo data failed: ", err)
	}
	log.GetLogger(log.Name.Root).Infoln("transaction data: ", tx.Data())

	multiSignAddr := patch.GetAddressBySmartContractName("GeneralMultiSigWallet", ethClient.ChainID.Int64())
	log.GetLogger(log.Name.Root).Infoln("GeneralMultiSigWallet address: ", multiSignAddr.Hex())

	multiSign, err := mainnet_contract.NewGeneralMultiSigWallet(multiSignAddr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("NewGeneralMultiSigWallet error: ", err)
	}
	transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("GetNewTransactorFromPrivateKey 1 error: ", err)
	}
	//Submit transaction with AddCurrencyInfo data
	tx1, err := multiSign.SubmitTransaction(transOpt1, configAddr, big.NewInt(0), tx.Data())
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("SubmitTransaction  error: ", err)
	}
	log.GetLogger(log.Name.Root).Infoln("Tran hash: ", tx1.Hash().Hex())
	receipt, err := ethClient.QueryEthReceiptWithSetting(tx1, 1, 120)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("QueryEthReceiptWithSetting failed: ", err)
	}
	receiptData, _ := json.Marshal(receipt)
	log.GetLogger(log.Name.Root).Infoln("Receipt data: ", string(receiptData))

	if receipt.Status == 0 {
		log.GetLogger(log.Name.Root).Errorln("Receipt status is 0")
	} else {
		multiSignWalletAbi := ethereum.DefaultABIDecoder.GetABI("GeneralMultiSigWallet")
		var isFound bool = false
		for _, logData := range receipt.Logs {
			if logData.Topics[0] == multiSignWalletAbi.Events["Submission"].ID {
				transId := logData.Topics[1].Big()
				log.GetLogger(log.Name.Root).Infoln("Transaction ID: ", transId.String())
				isFound = true
				break
			}
		}
		if !isFound {
			log.GetLogger(log.Name.Root).Errorln("Unable to find transaction ID")
		}
	}
}

func ConfirmTransaction(config *Patch2Config, ethClient *ethereum.EthClient) {

	var transId uint64
	var err error

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Input transaction ID:")
		transIdStr, _ := reader.ReadString('\n')
		transIdStr = strings.TrimSuffix(transIdStr, "\n")
		transId, err = strconv.ParseUint(transIdStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid input: ", err)
			continue
		}
		break
	}

	priKeyStr := ReadPrivateKey()

	multiSignAddr := patch.GetAddressBySmartContractName("GeneralMultiSigWallet", ethClient.ChainID.Int64())
	log.GetLogger(log.Name.Root).Infoln("GeneralMultiSigWallet address: ", multiSignAddr.Hex())
	multiSign, err := mainnet_contract.NewGeneralMultiSigWallet(multiSignAddr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("NewGeneralMultiSigWallet error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("GetNewTransactorFromPrivateKey error: ", err)
	}
	transIdBig := big.NewInt(0)
	transIdBig.SetUint64(transId)
	transOpt.GasLimit = config.GasLimit
	tx, err := multiSign.ConfirmTransaction(transOpt, transIdBig)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("ConfirmTransaction error: ", err)
	}
	log.GetLogger(log.Name.Root).Infoln("Trans hash: ", tx.Hash().Hex())

	receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 120)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("QueryEthReceiptWithSetting error: ", err)
	}
	receiptData, _ := json.Marshal(receipt)
	log.GetLogger(log.Name.Root).Infoln("Receipt data: ", string(receiptData))
	if receipt.Status == 0 {
		log.GetLogger(log.Name.Root).Errorln("Receipt status is 0")
	} else {
		multiSignAbi := ethereum.DefaultABIDecoder.GetABI("GeneralMultiSigWallet")
		var isFound bool = false
		for _, logData := range receipt.Logs {
			if logData.Topics[0] == multiSignAbi.Events["Execution"].ID {
				isFound = true
				log.GetLogger(log.Name.Root).Infoln("Execution successful, trans ID: ", transId)

				PrintEurusInternalConfigAssetList(ethClient, "MST")

				break
			} else if logData.Topics[0] == multiSignAbi.Events["ExecutionFailure"].ID {
				isFound = true
				log.GetLogger(log.Name.Root).Infoln("Execution failed, trans ID: ", transId)
				break
			}
		}
		if !isFound {
			log.GetLogger(log.Name.Root).Infoln("More confirmation requried for trans ID: ", transId)
		}
	}
}

func ReadPrivateKey() string {
	priKeyStr, err := patch.ReadTerminalHiddenInput("Input private key hex")
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("ReadTerminalHiddenInput failed: ", err)
	}

	priKey, err := go_ethereum_crypto.HexToECDSA(string(priKeyStr))
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("HexToECDSA error: ", err)
	}

	ownerAddress := go_ethereum_crypto.PubkeyToAddress(priKey.PublicKey)

	log.GetLogger(log.Name.Root).Infoln("Owner address is :", ownerAddress.Hex())

	return string(priKeyStr)
}

func PrintEurusInternalConfigAssetList(ethClient *ethereum.EthClient, checkName string) {

	addr := patch.GetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>", ethClient.ChainID.Int64())
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(addr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("NewEurusInternalConfig error: ", err)
		return
	}
	assetNameList, addressList, err := eurusInternalConfig.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetAssetAddress error: ", err)
		return
	}

	var assetList []string = make([]string, 0)
	var isFound bool = false
	for i, name := range assetNameList {
		if name == checkName {
			isFound = true
		}
		assetList = append(assetList, fmt.Sprintf("%s: %s", name, addressList[i].Hex()))
	}
	log.GetLogger(log.Name.Root).Infoln(assetList)
	if !isFound {
		log.GetLogger(log.Name.Root).Errorln(checkName, " not found in EurusInternalConfig smart contract")
	}
}
