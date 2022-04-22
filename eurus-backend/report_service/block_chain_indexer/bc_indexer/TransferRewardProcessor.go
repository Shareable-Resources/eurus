package bc_indexer

import (
	"encoding/json"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/marketing/reward"
	"eurus-backend/smartcontract/build/golang/contract"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type TransferRewardProcessor struct {
	reward.RewardProcessor
	MarketingRegWalletAddress *common.Address
	RegistrationCriteriaList  *reward.RegistrationRewardCriteriaList
}

func NewTransferRewardProcessor(internalSmartContractConfigAddr common.Address, sideChainEthClient *ethereum.EthClient,
	db *database.Database, slaveDb *database.ReadOnlyDatabase, marketingInvokerPrivateKey string, logger *logrus.Logger) *TransferRewardProcessor {
	processor := new(TransferRewardProcessor)
	processor.SideChainEthClient = sideChainEthClient
	processor.InternalSmartContractConfigAddr = internalSmartContractConfigAddr
	processor.Logger = logger
	processor.MarketingInvokerPrivateKey = marketingInvokerPrivateKey
	processor.DbProcessor = reward.NewRewardDBProcessor(db, slaveDb)
	return processor
}

func (me *TransferRewardProcessor) Init(registrationCriteriaListJsonStr string) error {
	internalSC, err := contract.NewInternalSmartContractConfig(me.InternalSmartContractConfigAddr, me.SideChainEthClient.Client)
	if err != nil {
		return err
	}

	addr, err := internalSC.GetMarketingRegWalletAddress(&bind.CallOpts{})
	if err != nil {
		return err
	}

	me.MarketingRegWalletAddress = &addr

	if registrationCriteriaListJsonStr != "" {
		err = json.Unmarshal([]byte(registrationCriteriaListJsonStr), &me.RegistrationCriteriaList)
		if err != nil {
			me.Logger.Errorln("Unable to Unmarshal registration criteria list: ", err)
			return err
		}

	}

	return nil
}

func (me *TransferRewardProcessor) TransferRegistrationRewardToUser(extTrans *ExtractedTransaction, user *user.User) error {
	if user == nil {
		return nil
	}
	me.Logger.Debugln("Checking user if he is able to gain reward user id: ", user.Id)
	if me.RegistrationCriteriaList == nil || me.RegistrationCriteriaList.CenteralizeUserSetting == nil || me.RegistrationCriteriaList.DecentralizedUserSetting == nil {
		return nil
	}

	var criteriaSetting *reward.RegistrationUserRewardSetting
	if user.IsMetamaskAddr {
		criteriaSetting = me.RegistrationCriteriaList.DecentralizedUserSetting
	} else {
		criteriaSetting = me.RegistrationCriteriaList.CenteralizeUserSetting
	}

	if len(criteriaSetting.ExcludeSenderMap) > 0 {
		sender, err := extTrans.GetSender()
		if err != nil {
			me.Logger.Errorln("Get sender from transaction error: ", err, " user Id: ", user.Id)
			return err
		}

		_, ok := criteriaSetting.ExcludeSenderMap[sender]
		if ok {
			return nil
		}
	}

	distributed, err := me.DbProcessor.DbGetDistributedToken(user.Id, criteriaSetting.Reward.AssetName, reward.DistributedRegistration)
	if err != nil {
		me.Logger.Error("DbGetDistributedToken failed error: ", err, " wallet address: ", extTrans.GetTo())
		return err
	}

	if distributed.Id != nil && *distributed.Id > 0 {
		//Reward already granted
		me.Logger.Debugln("DbGetDistributedToken record found for user id: ", user.Id)
		return nil
	}

	me.Logger.Debugln("Getting criteria for asset name: ", extTrans.AssetName)
	criteria, found := criteriaSetting.Criteria[extTrans.AssetName]
	if !found {
		//Search for ANY currency criteria
		criteria, found = criteriaSetting.Criteria["*"]
		if !found {
			//Check asset name is no legible for reward
			me.Logger.Debugln("Getting criteria for asset name not found: ", extTrans.AssetName)
			return nil
		}
	}

	if !criteria.CompareCriteria.Compare(extTrans.Amount, criteria.CompareAmount) {
		//Total asset amount does not
		me.Logger.Debugln("Transfer amount does not exceed the criteria : ", extTrans.Amount.String(), " criteria: ", criteria.CompareCriteria, " compare amount: ", criteria.CompareAmount.String())
		return nil
	}

	distributedToken := new(reward.DistributedToken)
	distributedToken.AssetName = criteriaSetting.Reward.AssetName

	distributedToken.Amount = decimal.NewFromBigInt(criteriaSetting.Reward.Amount, 0)
	chainId := me.SideChainEthClient.ChainID.Uint64()
	distributedToken.Chain = &chainId
	distributedToken.DistributedType = reward.DistributedRegistration
	distributedToken.TriggerType = reward.TriggerSideChainTransfer
	distributedToken.UserId = user.Id
	distributedToken.FromAddress = ethereum.ToLowerAddressString(me.MarketingRegWalletAddress.Hex())
	distributedToken.ToAddress = ethereum.ToLowerAddressString(extTrans.GetTo())
	distributedToken.Status = reward.DistributedStatusPending
	distributedToken.InitDate()

	err = me.DbProcessor.DbInsertDistributedToken(distributedToken)
	if err != nil {
		me.Logger.Errorln("Unable to insert distributed_token: ", err, " wallet address: ", extTrans.GetTo(), " for user id: ", user.Id)
		return err
	}

	var receipt *ethereum.BesuReceipt
	if criteriaSetting.Reward.AssetName == "EUN" {
		me.Logger.Debugln("Going to transfer EUN to user id: ", user.Id)
		_, receipt, err = me.transferEUN(extTrans.GetTo(), criteriaSetting.Reward.Amount)
		if err != nil {
			me.Logger.Errorln("Unable to transfer EUN to user id: ", user.Id, " error: ", err)
			err = me.DbProcessor.DbUpdateDistributedToken(*distributedToken.Id, "", nil, 0, nil, reward.DistributedStatusError)
			if err != nil {
				me.Logger.Errorln("Unable to update distributed token record. Id: ", distributedToken.Id, " for user id: ", user.Id, " Error: ", err)
			}
			return err
		}

		gasPriceDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0)
		bigGasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
		gasFeeDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0).Mul(decimal.NewFromBigInt(bigGasUsed, 0))

		if receipt.Status == 0 {

			receiptData, _ := json.Marshal(receipt)
			me.Logger.Error("Transfer EUN failed. Dest wallet: ", extTrans.GetTo(), " receipt: ", string(receiptData))
			err = me.DbProcessor.DbUpdateDistributedToken(*distributedToken.Id, strings.ToLower(receipt.TxHash.Hex()), &gasPriceDec, receipt.GasUsed, &gasFeeDec, reward.DistributedStatusError)
			if err != nil {
				me.Logger.Errorln("Unable to update distributed token record where status = 0. Id: ", distributedToken.Id, " for user id: ", user.Id, " Error: ", err)
			}
			return errors.New("Transfer EUN failed")
		}
		me.Logger.Debugln("Transfer EUN success for user id: ", user.Id)
		err = me.DbProcessor.DbUpdateDistributedToken(*distributedToken.Id, strings.ToLower(receipt.TxHash.Hex()), &gasPriceDec, receipt.GasUsed, &gasFeeDec, reward.DistributedStatusSuccess)
		if err != nil {
			me.Logger.Errorln("Unable to update distributed token record. Id: ", distributedToken.Id, " for user id: ", user.Id, " Error: ", err)
		}
	} else {
		//TODO
		return errors.New("Asset type transfer is not yet implemented")
	}

	return nil
}

func (me *TransferRewardProcessor) transferEUN(walletAddr string, amount *big.Int) (*types.Transaction, *ethereum.BesuReceipt, error) {

	marketingRegWallet, err := contract.NewMarketingWallet(*me.MarketingRegWalletAddress, me.SideChainEthClient.Client)
	if err != nil {
		me.Logger.Errorln("NewMarketingWallet error: ", err.Error(), " wallet address: ", walletAddr)
		return nil, nil, err
	}
	transOpt, err := me.SideChainEthClient.GetNewTransactorFromPrivateKey(me.MarketingInvokerPrivateKey, me.SideChainEthClient.ChainID)
	if err != nil {
		me.Logger.Errorln("GetNewTransactorFromPrivateKey error: ", err.Error(), " wallet address: ", walletAddr)
		return nil, nil, err
	}

	transOpt.GasLimit = 1000000

	tx, err := marketingRegWallet.TransferETH(transOpt, common.HexToAddress(walletAddr), amount)
	if err != nil {
		me.Logger.Errorln("TransferETH error: ", err.Error(), " wallet address: ", walletAddr)
		return nil, nil, err
	}

	me.Logger.Infoln("transfer EUN to user trans hash ", tx.Hash().Hex(), " for wallet address: ", walletAddr)
	receipt, err := me.SideChainEthClient.QueryEthReceipt(tx)
	if err != nil {
		me.Logger.Errorln("Unable to transfer EUN to user: ", err.Error(), " wallet address: ", walletAddr)
		return tx, nil, err
	}
	return tx, receipt, nil
}
