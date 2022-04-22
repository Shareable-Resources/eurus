package main

import (
	"bytes"
	"eurus-backend/data_patch/patch"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Patch3 struct {
	patch.EthereumPatch
}

func main() {
	patch3 := new(Patch3)
	config := new(patch.PatchConfigBase)
	err := patch.LoadConfig("Patch3Config.json", config)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("LoadConfig error: ", err)
	}

	err = patch3.InitPatch(config, "patch3")
	if err != nil {
		panic(err)
	}

	err = patch3.RunPatch(patch3, nil)
	if err != nil {
		panic(err)
	}

	nameList, addrList, err := patch3.EurusInternalConfig.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to check the asset address: ", err)
	} else {
		var outputList []string = make([]string, 0)
		for i, assetName := range nameList {
			output := fmt.Sprintf("%s: %s", assetName, addrList[i])
			outputList = append(outputList, output)
		}
		log.GetLogger(log.Name.Root).Infoln(outputList)
	}
	fmt.Println("Patch ended")
}

func (me *Patch3) SubmitTransaction(context *patch.EthereumPatch, priKey string, config interface{}) ([]*big.Int, error) {

	var transIdList []*big.Int = make([]*big.Int, 0)

	implAddr := patch.GetAddressBySmartContractName("EurusInternalConfig", me.EthClient.ChainID.Int64())

	proxyInstance, err := mainnet_contract.NewOwnedUpgradeabilityProxy(context.EurusInternalConfigAddr, me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create proxy instance: ", err)
		return nil, err
	}

	checkAddr, err := proxyInstance.Implementation(&bind.CallOpts{})
	if err != nil || !bytes.Equal(checkAddr.Bytes(), implAddr.Bytes()) {

		transOpt0, err := patch.CreateTransOptNoSigner(context.EthClient, priKey, config.(*patch.PatchConfigBase))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("CreateTransOptNoSigner failed: ", err)
			return nil, err
		}

		txToBeExecuted, err := proxyInstance.UpgradeTo(transOpt0, implAddr)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to upgrade to new EurusInternalConfig: ", err)
			return nil, err
		}

		log.GetLogger(log.Name.Root).Infoln("Going to submit transaction to upgradeTo new EurusInternalConfig")
		transId, err := me.SubmitTransactionByGeneralMultiSign(priKey, context.EurusInternalConfigAddr, txToBeExecuted.Data())
		if err != nil {
			return nil, errors.Wrap(err, "Submit failed")
		}

		transIdList = append(transIdList, transId)
	}

	//Update currency addresses
	var currencyList []string = []string{"USDM", "BTCM", "ETHM", "MST"}
	var currencyAddrList []common.Address = make([]common.Address, 0)

	transOpt, err := patch.CreateTransOptNoSigner(context.EthClient, priKey, config.(*patch.PatchConfigBase))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("CreateTransOptNoSigner failed: ", err)
		return nil, err
	}

	for _, currency := range currencyList {
		addr := patch.GetAddressBySmartContractName("OwnedUpgradeabilityProxy<"+currency+">", context.EthClient.ChainID.Int64())
		log.GetLogger(log.Name.Root).Debugln("OwnedUpgradeabilityProxy<"+currency+">: ", addr.Hex())

		currencyAddrList = append(currencyAddrList, addr)

	}

	updateTx, err := context.EurusInternalConfig.BatchUpdateAssetAddress(transOpt, currencyList, currencyAddrList)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("AddCurrencyInfo failed: ", err)
		return nil, err
	}

	log.GetLogger(log.Name.Root).Infoln("Going to submit transaction to update currency list ", currencyList)
	transId, err := me.SubmitTransactionByGeneralMultiSign(priKey, context.EurusInternalConfigAddr, updateTx.Data())
	if err != nil {
		return nil, errors.Wrap(err, "Submit failed")
	}

	transIdList = append(transIdList, transId)
	return transIdList, nil
}
