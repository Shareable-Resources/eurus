package user

import (
	"bytes"
	"eurus-backend/foundation"
	"eurus-backend/foundation/log"
	"eurus-backend/sign_service/sign_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func DeployUserWallet(server *UserServer, userId uint64) (common.Address, error) {

	transOpt, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get transOpt error: ", err, " userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	}

	transOpt.GasLimit = uint64(server.Config.SideChainGasLimit)
	// somehow the address return by below function is fake, the follow real proxy address will be: receipt.ContractAddress
	proxyAddress, tx, _, err := contract.DeployOwnedUpgradeabilityProxy(transOpt, server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot deploy user wallet proxy. Error: ", err, " userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot deploy user wallet proxy")
	}

	log.GetLogger(log.Name.Root).Infoln("OwnedUpgradeabilityProxy<UserWalletProxy> address: ", proxyAddress.Hex(), " tx hash:", tx.Hash().Hex(), " userId: ", userId)

	receipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get deploy user wallet proxy receipt. Error:", err, " userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get deploy user wallet proxy receipt")
	}
	if receipt.Status != 1 {
		receiptData, _ := receipt.MarshalJSON()
		log.GetLogger(log.Name.Root).Error("cannot get deploy user wallet proxy receipt: ", string(receiptData), " userId: ", userId)
		return common.Address{}, errors.New("cannot get deploy user wallet proxy receipt: " + receipt.TxHash.Hex())
	}

	log.GetLogger(log.Name.Root).Infoln("Deploy Proxy hash : ", receipt.TxHash, "receipt status", receipt.Status, " userId: ", userId)
	log.GetLogger(log.Name.Root).Infoln("proxy address : ", receipt.ContractAddress, " userId: ", userId)

	internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.Config.InternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get internalSCConfig. userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get internalSCConfig")
	}
	userWalletProxyAddr, err := internalSCConfig.GetUserWalletProxyAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get GetUserWalletProxyAddress userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get GetUserWalletProxyAddress")
	}

	// transOpt2, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
	// //transOpt, err = server.EthClient.GetNewTransactorFromPrivateKey(server.Config.HdWalletPrivateKey, server.EthClient.ChainID)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Error("cannot get transOpt userId: ", userId)
	// 	return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	// }

	err = userWalletUpgradeTo(server, receipt.ContractAddress, userWalletProxyAddr, userId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		return common.Address{}, err
	}

	userWalletImpl, err := internalSCConfig.GetUserWalletAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get GetUserWalletAddress userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get GetUserWalletAddress")
	}

	err = setUserWalletImpl(server, receipt.ContractAddress, userWalletImpl, userId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		return common.Address{}, err
	}

	userWalletProxy, err := contract.NewUserWalletProxy(receipt.ContractAddress, server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user wallet proxy userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot get user wallet proxy")
	}

	//setInternalSCAddress
	log.GetLogger(log.Name.Root).Infoln("setInternalSCAddress userId: ", userId)
	err = retryFunction(server.Config, func() (bool, error) {

		internalSCAddr, err := userWalletProxy.GetInternalSCAddress(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Warnln("Unable to query internal smart contract address. checking ignored. userId: ", userId)
		} else {
			if bytes.Equal(internalSCAddr.Bytes(), common.HexToAddress(server.Config.InternalSCConfigAddress).Bytes()) {
				log.GetLogger(log.Name.Root).Infoln("internal smart contract address already set. userId: ", userId)
				return false, nil
			}
		}

		transOpt7, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("SetInternalSCAddress cannot get transOpt userId: ", userId)
			return false, wrapErrorWithUserId(err, "SetInternalSCAddress cannot get transOpt", userId)
		}

		transOpt7.GasLimit = uint64(server.Config.SideChainGasLimit)
		tx, err = userWalletProxy.SetInternalSCAddress(transOpt7, common.HexToAddress(server.Config.InternalSCConfigAddress))
		if err != nil {
			log.GetLogger(log.Name.Root).Error("SetInternalSCAddress failed. ", " userId: ", userId)
			return false, wrapErrorWithUserId(err, "SetInternalSCAddress failed", userId)
		}
		log.GetLogger(log.Name.Root).Infoln("SetInternalSCAddress tx: ", tx.Hash().Hex(), " userId: ", userId)

		receipt7, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get SetInternalSCAddress receipt. tx: ", tx.Hash().Hex(), " userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get SetInternalSCAddress receipt. tx: "+tx.Hash().Hex(), userId)
		}

		if receipt7.Status != 1 {
			receiptData, _ := receipt7.MarshalJSON()
			log.GetLogger(log.Name.Root).Error("SetInternalSCAddress receipt failed ", string(receiptData), " userId: ", userId)
			return false, wrapErrorWithUserId(errors.New("SetInternalSCAddress receipt failed"), "tx: "+tx.Hash().Hex(), userId)
		}
		return false, nil
	}, userId)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		return common.Address{}, err
	}

	userWallet, err := contract.NewUserWallet(receipt.ContractAddress, server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot new user wallet", err.Error(), " userId: ", userId)
		return common.Address{}, errors.Wrap(err, "cannot new user wallet")
	}

	err = retryFunction(server.Config, func() (bool, error) {

		requirementCount, err := userWallet.Required(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Warnln("Unable to get requirement: ", err, " checking ignored. userId: ", userId)
		} else {
			if requirementCount.Uint64() == 2 {
				log.GetLogger(log.Name.Root).Infoln("Requirement already set to correct value. userId: ", userId)
				return false, nil
			}
		}

		transOpt3, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("ChangeRequirement cannot get transOpt userId: ", userId)
			return false, wrapErrorWithUserId(err, "ChangeRequirement cannot get transOpt", userId)
		}

		transOpt3.GasLimit = uint64(server.Config.SideChainGasLimit)
		tx, err = userWallet.ChangeRequirement(transOpt3, big.NewInt(2))
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot change requirement", err.Error(), " userId: ", userId)
			return false, errors.Wrap(err, "cannot change requirement")
		}

		log.GetLogger(log.Name.Root).Infoln("change requirement tx: ", tx.Hash(), " userId: ", userId)
		receipt3, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get change requirement receipt tx: ", tx.Hash().Hex(), " error: ", err.Error(), " userId: ", userId)
			return false, errors.Wrap(err, "cannot change requirement")
		}

		if receipt3.Status != 1 {
			receiptData, _ := receipt3.MarshalJSON()
			log.GetLogger(log.Name.Root).Error("transaction fail: ", string(receiptData), " user Id: ", userId)
			return false, errors.New("cannot change user wallet requirement for tx: " + receipt3.TxHash.Hex())
		}
		return false, nil
	}, userId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		return common.Address{}, err
	}

	for _, serverDetail := range server.Config.UserObserverList {

		err = retryFunction(server.Config, func() (bool, error) {

			walletOperatorList, err := userWallet.GetWalletOperatorList(&bind.CallOpts{})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to get wallet operator list: ", err, " userId: ", userId)
				return false, wrapErrorWithUserId(err, "Unable to get wallet operator list", userId)
			}

			var found bool
			for _, walletOperatorAddr := range walletOperatorList {
				if bytes.Equal(walletOperatorAddr.Bytes(), common.HexToAddress(serverDetail.WalletAddress).Bytes()) {
					found = true
					break
				}
			}
			if found {
				log.GetLogger(log.Name.Root).Infoln("User observer ", serverDetail.WalletAddress, " already added. userId: ", userId)
				return false, nil
			}

			transOpt4, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("AddWalletOperator GetNewTransactorFromPrivateKey error: ", err, " user observer address: ", serverDetail.WalletAddress, " User wallet address: ", receipt.ContractAddress.Hex(), " userId: ", userId)
				return false, wrapErrorWithUserId(err, "AddWalletOperator GetNewTransactorFromPrivateKey", userId)
			}

			transOpt4.GasLimit = uint64(server.Config.SideChainGasLimit)
			tx, err := userWallet.AddWalletOperator(transOpt4, common.HexToAddress(serverDetail.WalletAddress))
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("AddWalletOperator error: ", err, " user observer address: ", serverDetail.WalletAddress, " User wallet address: ", receipt.ContractAddress.Hex(), " userId: ", userId)
				return false, wrapErrorWithUserId(err, "AddWalletOperator", userId)
			}
			log.GetLogger(log.Name.Root).Infoln("AddWalletOperator tx: ", tx.Hash())
			addOpReceipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("AddWalletOperator QueryEthReceiptWithSetting error: ", err, " user observer address: ", serverDetail.WalletAddress, " User wallet address: ", receipt.ContractAddress.Hex(), " userId: ", userId)
				return false, wrapErrorWithUserId(err, "AddWalletOperator QueryEthReceiptWithSetting", userId)
			}
			if addOpReceipt.Status == 0 {
				receiptByte, _ := addOpReceipt.MarshalJSON()
				log.GetLogger(log.Name.Root).Errorln("AddWalletOperator receipt status failed: ", string(receiptByte), " user observer address: ", serverDetail.WalletAddress, " User wallet address: ", receipt.ContractAddress.Hex(), " userId: ", userId)
				return false, wrapErrorWithUserId(err, "AddWalletOperator receipt status", userId)
			}
			return false, nil
		}, userId)

		if err != nil {
			log.GetLogger(log.Name.Root).Errorln(err)
			return common.Address{}, err
		}
	}

	for _, invokerAddr := range server.Config.InvokerAddressList {

		err = retryFunction(server.Config, func() (bool, error) {

			isWriterExists, err := userWallet.IsWriter(&bind.CallOpts{}, common.HexToAddress(invokerAddr))
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to query IsWriter: ", err, " invoker addr: ", invokerAddr, " userId: ", userId)
				return false, wrapErrorWithUserId(err, "Unable to query IsWriter invoker addr: "+invokerAddr, userId)
			}
			if isWriterExists {
				log.GetLogger(log.Name.Root).Infoln("Writer already exists, skip add this address: ", invokerAddr, " userId: ", userId)
				return false, nil
			}

			transOpt6, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Add writer - GetNewTransactorFromSignServer error: ", err)
				return false, wrapErrorWithUserId(err, "Add writer GetNewTransactorFromSignServer", userId)
			}
			log.GetLogger(log.Name.Root).Infoln("Add invoker as writer. userId: ", userId)

			transOpt6.GasLimit = uint64(server.Config.SideChainGasLimit)
			tx, err = userWallet.AddWriter(transOpt6, common.HexToAddress(invokerAddr))
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Add writer error: ", err, "invokerAddr: ", invokerAddr, " userId: ", userId)
				return false, err
			}
			receipt6, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
			if err != nil {
				log.GetLogger(log.Name.Root).Error("cannot get add writer receipt: ", "invokerAddr: ", invokerAddr, " error: ", err.Error(), " userId: ", userId)
				return false, err
			}
			if receipt6.Status != 1 {
				receiptData, _ := receipt6.MarshalJSON()
				log.GetLogger(log.Name.Root).Error("add writer failed fail: ", string(receiptData), "invokerAddr: ", invokerAddr, " userId: ", userId)
				return false, errors.New("add writer failed receipt: " + receipt6.TxHash.Hex())
			}
			return false, nil
		}, userId)

		if err != nil {
			log.GetLogger(log.Name.Root).Errorln(err)
			return common.Address{}, err
		}
	}

	log.GetLogger(log.Name.Root).Infoln("Deploy smart contract finished. userId: ", userId, " smart contract address: ", receipt.ContractAddress.Hex())

	return receipt.ContractAddress, nil
}

func wrapErrorWithUserId(err error, message string, userId uint64) error {
	return errors.Wrap(err, message+". userId: "+strconv.FormatUint(userId, 10))
}

func retryFunction(retryConfig foundation.IRetrySetting, functor func() (bool, error), userId uint64) error {
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
		log.GetLogger(log.Name.Root).Errorln("Error after all trials: ", err, " userId: ", userId)
	}
	return err
}

func userWalletUpgradeTo(server *UserServer, ownedUpgradeProxyAddr common.Address, userWalletProxyAddr common.Address, userId uint64) error {

	return retryFunction(server.Config, func() (bool, error) {

		proxy, err := contract.NewOwnedUpgradeabilityProxy(ownedUpgradeProxyAddr, server.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get proxy userId: ", userId)
			return true, wrapErrorWithUserId(err, "cannot get proxy", userId)
		}

		implAddr, err := proxy.Implementation(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get proxy Implementation: ", err, " userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get proxy Implementation", userId)
		}

		if bytes.Equal(userWalletProxyAddr.Bytes(), implAddr.Bytes()) {
			log.GetLogger(log.Name.Root).Infoln("OwnedUpgradeabilityProxy implement already set to UserWalletProxy address. userId: ", userId)
			return false, nil
		}

		transOpt2, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("userWalletUpgradeTo cannot get transOpt userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get proxy Implementation", userId)
		}

		log.GetLogger(log.Name.Root).Infoln("user wallet proxy : ", ownedUpgradeProxyAddr.Hex(), " userId: ", userId)

		transOpt2.GasLimit = uint64(server.Config.SideChainGasLimit)
		tx, err := proxy.UpgradeTo(transOpt2, userWalletProxyAddr)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot create user wallet upgrade proxy transaction ", err.Error(), " userId: ", userId)
		}
		log.GetLogger(log.Name.Root).Infoln("upgrade tx: ", tx.Hash(), " userId: ", userId)
		receipt2, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get upgrade user wallet proxy receipt. userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get upgrade user wallet proxy receipt", userId)
		}

		if receipt2.Status != 1 {
			receiptData, _ := receipt2.MarshalJSON()
			log.GetLogger(log.Name.Root).Error("transaction fail ", string(receiptData), " userId: ", userId)
			return false, wrapErrorWithUserId(errors.New("Receipt status is 0"), "cannot upgrade user wallet proxy for tx: "+receipt2.TxHash.Hex(), userId)
		}

		return false, nil
	}, userId)
}

func setUserWalletImpl(server *UserServer, ownedUpgradeProxyAddr common.Address, userWalletImplAddr common.Address, userId uint64) error {

	return retryFunction(server.Config, func() (bool, error) {

		proxy, err := contract.NewUserWalletProxy(ownedUpgradeProxyAddr, server.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get UserWalletProxy. userId: ", userId)
			return true, wrapErrorWithUserId(err, "cannot get UserWalletProxy", userId)
		}

		implAddr, err := proxy.GetUserWalletImplementation(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get UserWallet Implementation: ", err, " userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get UserWallet Implementation", userId)
		}

		if bytes.Equal(userWalletImplAddr.Bytes(), implAddr.Bytes()) {
			log.GetLogger(log.Name.Root).Infoln("UserWalletProxy UserWalletImplementation already set to UserWalletImpl address. userId: ", userId)
			return false, nil
		}

		transOpt2, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("SetUserWalletImplementation cannot get transOpt userId: ", userId)
			return false, wrapErrorWithUserId(err, "SetUserWalletImplementation cannot get transOpt", userId)
		}

		log.GetLogger(log.Name.Root).Infoln("user wallet proxy : ", ownedUpgradeProxyAddr.Hex(), " userId: ", userId)

		transOpt2.GasLimit = uint64(server.Config.SideChainGasLimit)
		tx, err := proxy.SetUserWalletImplementation(transOpt2, userWalletImplAddr)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot create user wallet upgrade proxy transaction ", err.Error(), " userId: ", userId)
		}
		log.GetLogger(log.Name.Root).Infoln("upgrade tx: ", tx.Hash(), " userId: ", userId)
		receipt2, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get SetUserWalletImplementation receipt. userId: ", userId)
			return false, wrapErrorWithUserId(err, "cannot get SetUserWalletImplementation receipt", userId)
		}

		if receipt2.Status != 1 {
			receiptData, _ := receipt2.MarshalJSON()
			log.GetLogger(log.Name.Root).Error("transaction fail ", string(receiptData), " userId: ", userId)
			return false, wrapErrorWithUserId(errors.New("Receipt status is 0"), "cannot SetUserWalletImplementation for tx: "+receipt2.TxHash.Hex(), userId)
		}

		return false, nil
	}, userId)
}

func SetUserWalletOwner(server *UserServer, ownerAddressString string, contractAddress common.Address, userId uint64) error {

	transOpt, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
	if err != nil {
		return err
	}

	userWalletSC, err := contract.NewUserWallet(contractAddress, server.EthClient.Client)

	transOpt.GasLimit = uint64(server.Config.SideChainGasLimit)
	tx, err := userWalletSC.SetWalletOwner(transOpt, common.HexToAddress(ownerAddressString))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("SetWalletOwner error: ", err, " userId: ", userId)
		return err
	}

	log.GetLogger(log.Name.Root).Infoln("set wallet owner tx: ", tx.Hash(), " userId: ", userId)

	receipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet owner. userId: ", userId)
		return err
	}
	if receipt.Status != 1 {
		log.GetLogger(log.Name.Root).Error("transaction failed, status is 0. TxHash: ", tx.Hash().Hex(), " userId: ", userId)
		return wrapErrorWithUserId(errors.New("set user wallet owner status is 0. TxHash: "+receipt.TxHash.Hex()), "", userId)
	}

	log.GetLogger(log.Name.Root).Infoln("SetWalletOwner success. TxHash: ", receipt.TxHash.Hex(), " userId: ", userId)
	return nil
}

func SetUserWalletInternalSmartContractConfig(server *UserServer, walletAddress common.Address, internalSCConfigAddr common.Address) (*types.Transaction, error) {
	transOpt, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
	if err != nil {
		return nil, err
	}

	userWalletSC, err := contract.NewUserWallet(walletAddress, server.EthClient.Client)

	if err != nil {
		return nil, err
	}

	transOpt.GasLimit = uint64(server.Config.SideChainGasLimit)
	tx, err := userWalletSC.SetInternalSmartContractConfig(transOpt, internalSCConfigAddr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	log.GetLogger(log.Name.Root).Infoln("set Internal Smart Contract Config tx: ", tx.Hash())

	receipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet owner")
		return nil, err
	}
	if receipt.Status != 1 {
		log.GetLogger(log.Name.Root).Error("transaction failed:")
		return nil, errors.New("set user wallet owner failed for hx: " + receipt.TxHash.Hex())
	}

	return tx, nil
}

func GetPurchaseTransaction(server *UserServer, transHash common.Hash) (*types.Transaction, error) {
	tx, isPending, err := server.EthClient.GetTransaction(transHash)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get transaction: ", err)
		return nil, err
	}

	if isPending {
		log.GetLogger(log.Name.Root).Error("Transaction is pending")
		return nil, errors.New("Transaction is pending")
	}

	return tx, nil
}

func GetEIP155TransactionSender(tx *types.Transaction) (*common.Address, error) {

	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
	if err != nil {
		return nil, err
	}

	addr := msg.From()
	return &addr, nil
}

func GetMaxTopUpGasAmount(server *UserServer) (*big.Int, error) {

	extSC, err := contract.NewExternalSmartContractConfig(common.HexToAddress(server.Config.ExternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create external smart contract instance: ", err)
		return nil, errors.Wrap(err, "Unable to create external smart contract instance")
	}
	return extSC.GetMaxTopUpGasAmount(&bind.CallOpts{})
}
