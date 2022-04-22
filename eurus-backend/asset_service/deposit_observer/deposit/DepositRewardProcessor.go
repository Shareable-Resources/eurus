package deposit

import (
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/marketing/reward"
	"eurus-backend/smartcontract/build/golang/contract"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type DepositRewardProcessor struct {
	reward.RewardProcessor
	MarketingRegWalletAddress *common.Address
	RegistrationCriteriaList  *reward.RegistrationRewardCriteriaList
}

func NewDepositRewardProcessor(internalSmartContractConfigAddr common.Address, sideChainEthClient *ethereum.EthClient,
	db *database.Database, slaveDb *database.ReadOnlyDatabase, marketingInvokerPrivateKey string, logger *logrus.Logger) *DepositRewardProcessor {

	processor := new(DepositRewardProcessor)
	processor.SideChainEthClient = sideChainEthClient
	processor.InternalSmartContractConfigAddr = internalSmartContractConfigAddr
	processor.Logger = logger
	processor.MarketingInvokerPrivateKey = marketingInvokerPrivateKey
	processor.DbProcessor = reward.NewRewardDBProcessor(db, slaveDb)
	return processor
}

func (me *DepositRewardProcessor) Init(registrationCriteriaListJsonStr string) error {
	err := me.RewardProcessor.Init(me.InternalSmartContractConfigAddr, me.SideChainEthClient)
	if err != nil {
		return err
	}
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
			return err
		}
	}

	return nil
}

func (me *DepositRewardProcessor) TransferRegistrationRewardToUser(depositTrans *asset.DepositTransaction, user *user.User) error {
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

	distributed, err := me.DbProcessor.DbGetDistributedToken(user.Id, criteriaSetting.Reward.AssetName, reward.DistributedRegistration)
	if err != nil {
		me.Logger.Error("DbGetDistributedToken failed error: ", err, " wallet address: ", depositTrans.InnetToAddress)
		return err
	}

	if distributed.Id != nil && *distributed.Id > 0 {
		//Reward already granted
		me.Logger.Debugln("DbGetDistributedToken record found for user id: ", user.Id)
		return nil
	}

	me.Logger.Debugln("Getting criteria for asset name: ", depositTrans.AssetName)
	criteria, found := criteriaSetting.Criteria[depositTrans.AssetName]
	if !found {
		//Search for ANY currency criteria
		criteria, found = criteriaSetting.Criteria["*"]
		if !found {
			//Check asset name is no legible for reward
			me.Logger.Debugln("Getting criteria for asset name not found: ", depositTrans.AssetName)
			return nil
		}
	}

	if !criteria.CompareCriteria.Compare(depositTrans.Amount.BigInt(), criteria.CompareAmount) {
		//Total asset amount does not
		me.Logger.Debugln("Deposit amount does not exceed the criteria : ", depositTrans.Amount.String(), " criteria: ", criteria.CompareCriteria, " compare amount: ", criteria.CompareAmount.String())
		return nil
	}

	var receipt *ethereum.BesuReceipt
	if criteriaSetting.Reward.AssetName == "EUN" {
		me.Logger.Debugln("Going to transfer EUN to user id: ", user.Id)
		_, receipt, err = me.TransferEUN(depositTrans.InnetToAddress, *me.MarketingRegWalletAddress, criteriaSetting.Reward.Amount)
		if err != nil {
			return err
		}
		if receipt.Status == 0 {
			receiptData, _ := json.Marshal(receipt)
			me.Logger.Error("Transfer EUN failed. Dest wallet: ", depositTrans.InnetToAddress, " receipt: ", string(receiptData))
			err = me.InsertDistributedTokenToDb(receipt, criteriaSetting.Reward.AssetName, depositTrans.InnetToAddress, criteriaSetting.Reward.Amount, user.Id, reward.DistributedStatusError)
			if err != nil {
				me.Logger.Errorln("InsertDistributedTokenToDb failed for user: ", user.Id, " error: ", err)
			}
			return errors.New("Transfer EUN failed")
		}
		me.Logger.Debugln("Transfer EUN success for user id: ", user.Id)
	} else {
		//TODO
		return errors.New("Asset type transfer is not yet implemented")
	}

	err = me.InsertDistributedTokenToDb(receipt, criteriaSetting.Reward.AssetName, depositTrans.InnetToAddress, criteriaSetting.Reward.Amount, user.Id, reward.DistributedStatusSuccess)
	if err != nil {
		me.Logger.Errorln("InsertDistributedTokenToDb failed for user: ", user.Id, " error: ", err)
	}
	return err
}

func (me *DepositRewardProcessor) InsertDistributedTokenToDb(receipt *ethereum.BesuReceipt, assetName string,
	toAddr string, amount *big.Int, userId uint64, status reward.TokenDistributedStatus) error {
	distributedToken := new(reward.DistributedToken)
	distributedToken.AssetName = assetName
	rewardAmount := decimal.NewFromBigInt(amount, 0)
	distributedToken.Amount = rewardAmount
	chainId := me.SideChainEthClient.ChainID.Uint64()
	distributedToken.Chain = &chainId
	distributedToken.DistributedType = reward.DistributedRegistration
	distributedToken.TriggerType = reward.TriggerDeposit
	distributedToken.UserId = userId
	distributedToken.TxHash = strings.ToLower(receipt.TxHash.Hex())
	distributedToken.FromAddress = ethereum.ToLowerAddressString(me.MarketingRegWalletAddress.Hex())
	distributedToken.ToAddress = toAddr
	gasPriceDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0)
	distributedToken.GasPrice = &gasPriceDec
	distributedToken.GasUsed = receipt.GasUsed

	bigGasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
	gasFeeDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0).Mul(decimal.NewFromBigInt(bigGasUsed, 0))

	distributedToken.GasFee = &gasFeeDec
	distributedToken.Status = status
	distributedToken.InitDate()

	err := me.DbProcessor.DbInsertDistributedToken(distributedToken)
	if err != nil {
		me.Logger.Errorln("Unable to insert distributed_token: ", err, " wallet address: ", toAddr)

	}
	return err
}
