package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"eurus-backend/data_patch/patch"
	"eurus-backend/env"
	"eurus-backend/foundation"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	go_ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
)

type Patch1Config struct {
	server.ServerConfigBase
	CentralizedUserList  []string `json:"centralizedUserList"`
	UserObserverAddrList []string `json:"userObserverAddrList"`
	GasLimit             uint64   `json:"gasLimit"`
}

func main() {

	var failCount int
	var successCount int
	var failAccountList []string = make([]string, 0)
	//Centralized user smart contract owner address
	var cenUserOwnerAddr common.Address = common.HexToAddress("0x374839556766d8D582689C1A651511d16B60b7A8")

	patch.InitLog()
	config := new(Patch1Config)

	err := patch.LoadConfig("Patch1Config.json", config)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("LoadConfig error: ", err)
	}

	err = patch.LoadSmartContractConfig("")
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln(err)
	}

	ethClient, err := patch.CreateEurusEthClient(&config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("CreateEurusEthClient error: ", err)
	}

	priKeyStr, err := patch.ReadTerminalHiddenInput("Input sign server centralized user owner private key hex")
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("ReadTerminalHiddenInput failed: ", err)
	}

	priKey, err := go_ethereum_crypto.HexToECDSA(string(priKeyStr))
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("HexToECDSA error: ", err)
	}

	userOwnerAddress := go_ethereum_crypto.PubkeyToAddress(priKey.PublicKey)
	log.GetLogger(log.Name.Root).Infoln("Centralized user smart contract owner address: ", userOwnerAddress.Hex())
	//Checking input private key wallet address is the expected owner address
	if !bytes.Equal(userOwnerAddress.Bytes(), cenUserOwnerAddr.Bytes()) {
		log.GetLogger(log.Name.Root).Panicln("Input private key does not match with centralized user owner address")
	}

	for _, addr := range config.CentralizedUserList {
		log.GetLogger(log.Name.Root).Infoln("==============================================")
		log.GetLogger(log.Name.Root).Infoln("Upgrade address: ", addr)
		userWalletAddr := common.HexToAddress(addr)

		proxySC, err := contract.NewOwnedUpgradeabilityProxy(userWalletAddr, ethClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("NewOwnedUpgradeabilityProxy failed for address: ", addr, " error: ", err)
			failAccountList = append(failAccountList, addr)
			failCount++
			continue
		}

		err = retryFunction(config, func() (bool, error) {

			implAddr, err := proxySC.Implementation(&bind.CallOpts{})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to get implementation for address: ", addr, " error: ", err)
				return false, err
			}

			userWalletProxyAddr := patch.GetAddressBySmartContractName("UserWalletProxy", env.DefaultEurusChainId)
			if bytes.Equal(implAddr.Bytes(), userWalletProxyAddr.Bytes()) {
				log.GetLogger(log.Name.Root).Infoln("Already set the user wallet proxy for address: ", addr)
				return false, nil
			}

			transOpt, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed for address: ", addr, " error: ", err)
				return false, err
			}
			transOpt.GasLimit = config.GasLimit

			tx, err := proxySC.UpgradeTo(transOpt, userWalletProxyAddr)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("UpgradeTo failed for address: ", addr, " error: ", err)
				return false, err
			}
			log.GetLogger(log.Name.Root).Infoln("Upgrade to Tx hash: ", tx.Hash().Hex())
			receipt, err := ethClient.QueryEthReceipt(tx)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("QueryEthReceipt failed for address: ", addr, " error: ", err)
				return false, err
			}
			receiptJson, _ := json.Marshal(receipt)
			log.GetLogger(log.Name.Root).Infoln(string(receiptJson))
			if receipt.Status == 0 {
				log.GetLogger(log.Name.Root).Errorln("Receipt status is 0 for user address: ", addr)
				return true, errors.New("UpgradeTo receipt status is 0")
			}
			log.GetLogger(log.Name.Root).Infoln("Upgrade UserWalletProxy success for address: ", addr)
			return false, nil
		}, addr)

		if err != nil {
			failAccountList = append(failAccountList, addr)
			failCount++
			continue
		}

		userWalletProxy, err := contract.NewUserWalletProxy(userWalletAddr, ethClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to create UserWalletProxy instance: ", err, " user addr: ", addr)
			failAccountList = append(failAccountList, addr)
			failCount++
			continue
		}

		err = retryFunction(config, func() (bool, error) {

			checkInternalConfigAddr, err := userWalletProxy.GetInternalSCAddress(&bind.CallOpts{})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Checking internal SC config address failed: ", err, " user addr: ", addr)
				return false, err
			}
			internalSCAddr := patch.GetAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>", env.DefaultEurusChainId)
			if bytes.Equal(checkInternalConfigAddr.Bytes(), internalSCAddr.Bytes()) {
				log.GetLogger(log.Name.Root).Infoln("Internal smart contract config already set for address: ", addr)
				return false, nil
			}

			transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("SetInternalSCAddress GetNewTransactorFromPrivateKey error: ", err)
				return false, err
			}
			transOpt1.GasLimit = config.GasLimit

			tx1, err := userWalletProxy.SetInternalSCAddress(transOpt1, internalSCAddr)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("SetInternalSCAddress SetInternalSCAddress failed: ", err)
				return false, err
			}
			log.GetLogger(log.Name.Root).Infoln("SetInternalSCAddress Tx hash: ", tx1.Hash().Hex())

			receipt1, err := ethClient.QueryEthReceipt(tx1)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("SetInternalSCAddress QueryEthReceipt failed for address: ", addr, " error: ", err)
				return false, err
			}
			receiptJson1, _ := json.Marshal(receipt1)
			log.GetLogger(log.Name.Root).Infoln(string(receiptJson1))

			if receipt1.Status == 0 {
				log.GetLogger(log.Name.Root).Errorln("SetInternalSCAddress receipt status is 0. User wallet: ", addr)
				return true, errors.New("SetInternalSCAddress receipt status is 0")
			}
			return false, nil

		}, addr)
		if err != nil {
			failAccountList = append(failAccountList, addr)
			failCount++
			continue
		}

		err = addWalletOperator(config, addr, string(priKeyStr), ethClient)
		if err != nil {
			failAccountList = append(failAccountList, addr)
			failCount++
			continue
		}

		err = retryFunction(config, func() (bool, error) {

			impAddr, err := proxySC.Implementation(&bind.CallOpts{})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Verify Implementation address failed: ", err, " user addr: ", addr)
				return true, err
			} else {
				log.GetLogger(log.Name.Root).Infoln("Implementation: ", impAddr.Hex(), " user addr: ", addr)
			}

			queryInternalSCAddr, err := userWalletProxy.GetInternalSCAddress(&bind.CallOpts{})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Verify internal SC config address failed: ", err, " user addr: ", addr)
				return true, err
			} else {
				log.GetLogger(log.Name.Root).Infoln("Internal SC config address set: ", queryInternalSCAddr.Hex(), " user addr: ", addr)
			}

			return false, nil

		}, addr)

		if err != nil {
			failAccountList = append(failAccountList, addr)
			failCount++
		} else {
			successCount++
		}
	}
	log.GetLogger(log.Name.Root).Infof("Total: %d, Success: %d, Fail: %d", len(config.CentralizedUserList), successCount, failCount)
	var failAccountStr string
	for _, failAddr := range failAccountList {
		failAccountStr = fmt.Sprintf("%s%s,", failAccountStr, failAddr)
	}
	log.GetLogger(log.Name.Root).Infoln("Fail account list: ", failAccountStr)
	log.GetLogger(log.Name.Root).Infoln("Finished")
}

func addWalletOperator(config *Patch1Config, addr string, priKeyStr string, ethClient *ethereum.EthClient) error {

	userWalletAddr := common.HexToAddress(addr)
	userWallet, err := contract.NewUserWallet(userWalletAddr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create UserWallet instance: ", err, " user addr: ", addr)
		return err
	}
	var walletOperatorList []common.Address = make([]common.Address, 0)
	err = retryFunction(config, func() (bool, error) {
		walletOperatorList, err = userWallet.GetWalletOperatorList(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to get walletOperatorList for address: ", addr, " error: ", err)
			return false, err
		} else {
			for _, writer := range walletOperatorList {
				log.GetLogger(log.Name.Root).Debugln("Writer: ", writer.Hex())
			}
		}
		return false, nil
	}, addr)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get walletOperatorList (aborted) for address: ", addr, " error: ", err)
		return err
	}

	var userObsToBeInsertList []common.Address = make([]common.Address, 0)
	for _, userObsAddrStr := range config.UserObserverAddrList {
		userObsAddr := common.HexToAddress(userObsAddrStr)
		var obsFound bool = false
		for _, writerAddr := range walletOperatorList {
			if bytes.Equal(writerAddr.Bytes(), userObsAddr.Bytes()) {
				obsFound = true
				break
			}
		}
		if obsFound {
			continue
		}
		userObsToBeInsertList = append(userObsToBeInsertList, userObsAddr)
	}

	if len(userObsToBeInsertList) > 0 {

		for _, insertObsAddr := range userObsToBeInsertList {
			err1 := retryFunction(config, func() (bool, error) {
				transOpt2, err := ethClient.GetNewTransactorFromPrivateKey(string(priKeyStr), ethClient.ChainID)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("AddWalletOperator GetNewTransactorFromPrivateKey error: ", err)
					return false, err
				}
				transOpt2.GasLimit = config.GasLimit
				tx2, err := userWallet.AddWalletOperator(transOpt2, insertObsAddr)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("AddWalletOperator address failed: ", err, " user addr: ", addr)
					return true, err
				}
				receipt2, err := ethClient.QueryEthReceiptWithSetting(tx2, 1, config.RetryCount)

				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("AddWalletOperator QueryEthReceipt failed for address: ", addr, " error: ", err)
					return false, err
				}
				receiptJson2, _ := json.Marshal(receipt2)
				log.GetLogger(log.Name.Root).Infoln(string(receiptJson2))

				if receipt2.Status == 0 {
					log.GetLogger(log.Name.Root).Errorln("AddWalletOperator receipt status is 0. User wallet: ", addr)
					return true, errors.New("AddWalletOperator receipt status is 0")
				}
				return false, nil
			}, addr)

			if err1 != nil {
				log.GetLogger(log.Name.Root).Errorln("AddWriter aborted. User wallet: ", addr)
				return err1
			}
		}

	}

	return nil
}

func retryFunction(retryConfig foundation.IRetrySetting, functor func() (bool, error), walletAddress string) error {
	var err error
	var isFatalError bool
	for i := 0; i < retryConfig.GetRetryCount(); i++ {
		isFatalError, err = functor()
		if err != nil {
			if !isFatalError {
				time.Sleep(retryConfig.GetRetryInterval() * time.Second)
				continue
			} else {
				break
			}
		} else {
			return nil
		}
	}
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Error after all trials: ", err, " wallet address: ", walletAddress)
	}
	return err
}
