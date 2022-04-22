package patch

import (
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	go_ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type IEthereumSubmitPatch interface {
	SubmitTransaction(patcher *EthereumPatch, priKey string, config interface{}) ([]*big.Int, error)
}

type IEthereumConfirmPatch interface {
	ConfirmTransaction(patcher *EthereumPatch, priKey string, transId []*big.Int, config interface{}) error
}

type EthereumPatch struct {
	Config                  *PatchConfigBase
	EthClient               *ethereum.EthClient
	EurusInternalConfig     *mainnet_contract.EurusInternalConfig
	EurusInternalConfigAddr common.Address
}

func (me *EthereumPatch) InitPatch(config *PatchConfigBase, scConfigFileNameSuffix string) error {
	me.Config = config
	InitLog()
	err := LoadSmartContractConfig(scConfigFileNameSuffix)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		return err
	}

	_ = ethereum.DefaultABIDecoder.ImportABIJson("GeneralMultiSigWallet", mainnet_contract.GeneralMultiSigWalletABI)

	me.EthClient, err = CreateEurusEthClient(config.GetServerConfigBase())
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("CreateEurusEthClient error: ", err)
		return err
	}

	me.EurusInternalConfigAddr = GetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>", me.EthClient.ChainID.Int64())
	me.EurusInternalConfig, err = mainnet_contract.NewEurusInternalConfig(me.EurusInternalConfigAddr, me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create EurusInternalConfig object: ", err)
		return err
	}
	return nil
}

func (me *EthereumPatch) RunPatch(submitter IEthereumSubmitPatch, confirmer IEthereumConfirmPatch) error {

	priKeyStr := me.ReadPrivateKey("Input submit transations private key hex")
	confirmPriKeyStr := me.ReadPrivateKey("Input confirm transations private key hex")

	transIdList, err := submitter.SubmitTransaction(me, priKeyStr, me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("SubmitTransaction failed " + err.Error())
	}

	log.GetLogger(log.Name.Root).Infoln("Transaction ID list : ", transIdList)

	if confirmer != nil {
		err = confirmer.ConfirmTransaction(me, confirmPriKeyStr, transIdList, me.Config)
		if err != nil {
			log.GetLogger(log.Name.Root).Panicln("Confirm transaction failed: ", err)
		}
	} else {
		_ = me.ConfirmTransaction(confirmPriKeyStr, transIdList)
	}

	return nil
}

func (me *EthereumPatch) ReadPrivateKey(prompt string) string {
	priKeyStr, err := ReadTerminalHiddenInput(prompt)
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

func (me *EthereumPatch) SubmitTransactionByGeneralMultiSign(priKeyStr string, targetAddress common.Address, data []byte) (*big.Int, error) {

	multiSignAddr := GetAddressBySmartContractName("GeneralMultiSigWallet", me.EthClient.ChainID.Int64())
	log.GetLogger(log.Name.Root).Infoln("GeneralMultiSigWallet address: ", multiSignAddr.Hex())

	multiSign, err := mainnet_contract.NewGeneralMultiSigWallet(multiSignAddr, me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("NewGeneralMultiSigWallet error: ", err)
		return nil, err
	}
	transOpt1, err := me.EthClient.GetNewTransactorFromPrivateKey(priKeyStr, me.EthClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey 1 error: ", err)
		return nil, err
	}
	transOpt1.GasLimit = me.Config.GasLimit
	//Submit transaction with AddCurrencyInfo data
	tx1, err := multiSign.SubmitTransaction(transOpt1, targetAddress, big.NewInt(0), data)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("SubmitTransaction  error: ", err)
		return nil, err
	}
	log.GetLogger(log.Name.Root).Infoln("Submit Tran hash: ", tx1.Hash().Hex())
	receipt, err := me.EthClient.QueryEthReceiptWithSetting(tx1, 1, 120)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("QueryEthReceiptWithSetting failed: ", err)
		return nil, err
	}
	receiptData, _ := json.Marshal(receipt)
	log.GetLogger(log.Name.Root).Infoln("Receipt data: ", string(receiptData))

	if receipt.Status == 0 {
		log.GetLogger(log.Name.Root).Errorln("Receipt status is 0. Tx hash: ", tx1.Hash().String())
		return nil, errors.New("Receipt status is 0")
	} else {
		multiSignWalletAbi := ethereum.DefaultABIDecoder.GetABI("GeneralMultiSigWallet")

		for _, logData := range receipt.Logs {
			if logData.Topics[0] == multiSignWalletAbi.Events["Submission"].ID {
				transId := logData.Topics[1].Big()
				log.GetLogger(log.Name.Root).Infoln("Transaction ID: ", transId.String())
				return transId, nil
			}
		}

		log.GetLogger(log.Name.Root).Errorln("Unable to find transaction ID. Tx hash: ", tx1.Hash().String())

	}

	return nil, errors.New("Trans ID not found")
}

func (me *EthereumPatch) ConfirmTransaction(priKey string, transIdList []*big.Int) error {

	multiSignAbi := ethereum.DefaultABIDecoder.GetABI("GeneralMultiSigWallet")

	multiSignAddr := GetAddressBySmartContractName("GeneralMultiSigWallet", me.EthClient.ChainID.Int64())
	log.GetLogger(log.Name.Root).Infoln("GeneralMultiSigWallet address: ", multiSignAddr.Hex())

	multiSign, err := mainnet_contract.NewGeneralMultiSigWallet(multiSignAddr, me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("NewGeneralMultiSigWallet error: ", err)
		return err
	}

	for _, transId := range transIdList {
		transOpt1, err := me.EthClient.GetNewTransactorFromPrivateKey(priKey, me.EthClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey 1 error: ", err)
			return err
		}
		transOpt1.GasLimit = me.Config.GasLimit
		tx, err := multiSign.ConfirmTransaction(transOpt1, transId)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to confirm transaction ID : ", transId.String(), " error: ", err)
			continue
		}
		log.GetLogger(log.Name.Root).Infoln("Confirmation tx hash for transaction ID: ", transId.String(), " ", tx.Hash().String())
		receipt, err := me.EthClient.QueryEthReceipt(tx)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to query receipt for transaction ID: ", transId.String(), " error: ", err)
			continue
		}

		if receipt.Status == 0 {
			receiptData, _ := json.Marshal(receipt)
			log.GetLogger(log.Name.Root).Errorln("Receipt is 0 for transaction ID: ", transId.String(), " receipt: ", string(receiptData))
			continue
		} else {

			var isFound bool = false
			var isFailed bool = false

			for _, logData := range receipt.Logs {
				if logData.Topics[0] == multiSignAbi.Events["Execution"].ID {
					isFound = true
					log.GetLogger(log.Name.Root).Infoln("Execution successful, trans ID: ", transId)
					break
				} else if logData.Topics[0] == multiSignAbi.Events["ExecutionFailure"].ID {
					isFound = true
					isFailed = true
					log.GetLogger(log.Name.Root).Errorln("Execution failed, trans ID: ", transId)
					break
				}
			}
			if !isFound {
				log.GetLogger(log.Name.Root).Warnln("More confirmation requried for trans ID: ", transId)
			}

			if !isFailed {
				log.GetLogger(log.Name.Root).Infoln("Confirm successfully. transaction ID: ", transId.String(), " tx hash: ", tx.Hash().String())
			} else {
				log.GetLogger(log.Name.Root).Errorln("Confirm execution failed. transaction ID: ", transId.String(), " tx hash: ", tx.Hash().String())

			}
		}

	}

	log.GetLogger(log.Name.Root).Infoln("Confirmation finished")
	return nil

}
