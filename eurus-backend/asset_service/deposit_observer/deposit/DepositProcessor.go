package deposit

import (
	"bytes"
	"context"
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"math/big"
	"time"

	go_ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type DepositProcessor struct {
	context                     *DepositProcessorContext //context has replaced db
	scProcessor                 *DepositSCProcessor
	rewardProcessor             *DepositRewardProcessor
	config                      *DepositObserverConfig
	mainnetAssetAddressList     []common.Address
	mintCompletedEventId        common.Hash
	depositETHEventId           common.Hash
	TransferEventID             common.Hash
	TransferMethodID            []byte
	currentMainnetBlockNumber   *big.Int
	currentMainnetBlockHash     common.Hash
	currentSideChainBlockNumber *big.Int
	currentSideChainBlockHash   common.Hash
	isVerboseLog                bool
}

func NewDepositProcessor(config *DepositObserverConfig, scProcessor *DepositSCProcessor, context *DepositProcessorContext) *DepositProcessor {
	processor := new(DepositProcessor)
	processor.config = config
	processor.context = context
	processor.scProcessor = scProcessor
	processor.rewardProcessor = NewDepositRewardProcessor(common.HexToAddress(config.InternalSCConfigAddress),
		scProcessor.sideChainEthClient, context.db, context.slaveDb,
		processor.config.HdWalletPrivateKey, log.GetLogger(context.LoggerName))
	return processor
}

func (me *DepositProcessor) Init() error {

	erc20 := ethereum.DefaultABIDecoder.GetABI("ERC20")
	transferEvent, ok := erc20.Events["Transfer"]
	if !ok {
		return errors.New("ERC20 ABI error. Unable to find Sweep event")
	}
	me.TransferEventID = transferEvent.ID

	transferMethod, ok := erc20.Methods["transfer"]
	if !ok {
		return errors.New("ERC20 ABI error. Unable to find Sweep event")
	}
	me.TransferMethodID = transferMethod.ID

	platformWalletAbiObj := ethereum.DefaultABIDecoder.GetABI("PlatformWallet")

	mintCompletedEvent, ok := platformWalletAbiObj.Events["MintCompletedEvent"]
	if !ok {
		return errors.New("EurusUserDeposit ABI error. Unable to find MintCompletedEvent event")
	}
	me.mintCompletedEventId = mintCompletedEvent.ID
	proxyObj := ethereum.DefaultABIDecoder.GetABI("OwnedUpgradeabilityProxy")
	depositETHEvent, ok := proxyObj.Events["DepositETH"]
	if !ok {
		return errors.New("Unable to get depositETH event ID")
	}

	me.depositETHEventId = depositETHEvent.ID

	err := me.rewardProcessor.Init(me.config.RegistrationRewardSetting)
	if err != nil {
		return errors.Wrap(err, "Init Reward process error")
	}

	me.mainnetAssetAddressList = make([]common.Address, 0)
	for _, addr := range me.scProcessor.MainnetAssetNameToAddress {
		me.mainnetAssetAddressList = append(me.mainnetAssetAddressList, addr)
	}
	return nil
}

func (me *DepositProcessor) RunMainnetBlockSubscriberAsync(subscriber *ethereum.BlockSubscriber) {

	go func(subscriber *ethereum.BlockSubscriber) {
		for {
			block, serverErr := subscriber.GetLatestBlock(false)
			if serverErr != nil {
				if serverErr.GetReturnCode() == foundation.RecordNotFound {
					log.GetLogger(me.context.LoggerName).Error(serverErr.Message)
					continue
				}

				if serverErr.GetReturnCode() != foundation.NetworkError {
					log.GetLogger(me.context.LoggerName).Error("Unable to get block: ", serverErr.Error())
				}
				time.Sleep(2 * time.Second)
				continue
			}
			log.GetLogger(log.Name.Root).Infoln("Processing mainnet block number: ", block.Number().String(), " block hash: ", block.Hash().Hex(), " transaction count: ", len(block.Transactions()))
			me.currentMainnetBlockNumber = block.Number()
			me.currentMainnetBlockHash = block.Hash()

			if BlockHasTransaction(block) {
				me.processMainnetTransaction(block)
			}
			log.GetLogger(log.Name.Root).Infoln("Finished mainnet block number: ", block.Number().String(), " block hash: ", block.Hash().Hex())

			err := me.context.MainnetRescanCounter.UpdateLatestBlock(block.Number())
			if err != nil {
				log.GetLogger(me.context.LoggerName).Errorln("Unable to update mainnet latest block number: ", block.Number().String(), " error: ", err.Error())
			}

		}
	}(subscriber)
}

func (me *DepositProcessor) RunSideChainBlockSubscriberAsync(subscriber *ethereum.BlockSubscriber) {
	go func(subscriber *ethereum.BlockSubscriber) {
		for {
			block, serverErr := subscriber.GetLatestBlock(true)
			if serverErr != nil {
				if serverErr.GetReturnCode() != foundation.NetworkError {
					log.GetLogger(me.context.LoggerName).Error("Unable to get block by block number: ", serverErr.Error())
				}
				time.Sleep(2 * time.Second)
				continue
			}
			me.currentSideChainBlockNumber = block.Number()
			me.currentSideChainBlockHash = block.Hash()
			if me.isVerboseLog {
				log.GetLogger(me.context.LoggerName).Debugln("Processing side chain block number: ", block.Number(), " block hash: ", block.Hash(), " transaction count: ", len(block.Transactions()))
			}
			if BlockHasTransaction(block) {
				me.processSideChainTransaction(block)
			}
			if me.isVerboseLog {
				log.GetLogger(me.context.LoggerName).Debugln("Finished side chain block number: ", block.Number(), " block hash: ", block.Hash())
			}
			err := me.context.SideChainRescanCounter.UpdateLatestBlock(block.Number())
			if err != nil {
				log.GetLogger(me.context.LoggerName).Errorln("Unable to update side chain latest block number: ", block.Number().String(), " error: ", err.Error())
			}
		}
	}(subscriber)
}

func (me *DepositProcessor) processSideChainTransaction(block *types.Block) {

	for _, trans := range block.Transactions() {
		_, _, status := ethereum.DefaultABIDecoder.DecodeABIInputArgument(trans.Data(), "PlatformWallet", "submitMintRequest")
		var receipt *ethereum.BesuReceipt
		var err error
		if status == ethereum.ExtractSuccess {
			receipt, err = me.scProcessor.QuerySideChainEthReceiptWithSetting(trans, 1, -1)
			if err != nil {
				log.GetLogger(me.context.LoggerName).Errorln("Unable to query receipt for trans hash: ", trans.Hash(), " Error: ", err.Error())
				continue
			}

		} else {
			_, _, status = ethereum.DefaultABIDecoder.DecodeABIInputArgument(trans.Data(), "PlatformWallet", "confirmTransaction")
			if status == ethereum.ExtractSuccess {
				receipt, err = me.scProcessor.QuerySideChainEthReceiptWithSetting(trans, 1, -1)
				if err != nil {
					log.GetLogger(me.context.LoggerName).Errorln("Unable to query receipt for trans hash: ", trans.Hash(), " Error: ", err.Error())
					continue
				}
			}
		}

		if status != ethereum.ExtractSuccess {
			continue
		}

		for _, logMessage := range receipt.Logs {
			if bytes.Equal(logMessage.Address.Bytes(), me.scProcessor.SidechainPlatformWalletAddr.Bytes()) {
				me.processSideChainLog(logMessage)
			}
		}
	}
}

func (me *DepositProcessor) processMainnetTransaction(block *types.Block) {
	transactions := block.Transactions()
	zero := big.NewInt(0)

	var transLogMap map[common.Hash][]types.Log = make(map[common.Hash][]types.Log)

	if len(transactions) > 0 {
		blockHash := block.Hash()
		query := go_ethereum.FilterQuery{
			BlockHash: &blockHash,
			Addresses: me.mainnetAssetAddressList,
			Topics: [][]common.Hash{
				{me.TransferEventID},
			},
		}

		logList, _ := me.scProcessor.mainnetEthClient.Client.FilterLogs(context.Background(), query)
		for _, logObj := range logList {
			arr, ok := transLogMap[logObj.TxHash]
			if !ok {
				arr = make([]types.Log, 0)
			}
			arr = append(arr, logObj)
			transLogMap[logObj.TxHash] = arr
		}

		for _, trans := range transactions {
			//Matching tranfer event log to transaction object
			if logObjList, ok := transLogMap[trans.Hash()]; ok {
				//Transaction contains ERC20 transfer event
				me.processTransfer(block, trans, true, logObjList)
			} else {
				if trans.To() != nil && trans.Value() != nil && trans.Value().Cmp(zero) > 0 {
					me.processTransfer(block, trans, false, nil)
				}
			}
		}
	}
}

func (me *DepositProcessor) processSideChainLog(transLog *types.Log) {

	topicBytes := transLog.Topics[0].Bytes()
	if bytes.Equal(topicBytes, me.mintCompletedEventId.Bytes()) {
		//MintCompleted event
		me.processMintCompletedEvent(transLog)
	}
}

func (me *DepositProcessor) processTransfer(block *types.Block, tx *types.Transaction, isERC20 bool, transferLogList []types.Log) {

	var assetTrans *AssetTransferTransaction = NewAssetTransferTransaction(block, tx, transferLogList)

	if isERC20 {
		assetTrans = FilterERC20TransactionOnly(assetTrans, me.scProcessor.MainnetAssetInfo)
		if assetTrans == nil {
			return
		}

		var isFound bool
		var candidateTransferLog types.Log

		for _, transferLog := range assetTrans.TransferLog {

			mainnetReceipant := common.BytesToAddress(transferLog.Topics[2].Bytes())
			senderAddr := common.BytesToAddress(transferLog.Topics[1].Bytes())

			if bytes.Equal(mainnetReceipant.Bytes(), me.scProcessor.MainnetUserDepositAddr.Bytes()) {
				//When send ERC20 to EurusUserDeposit smart contract, it is assumed it is a decentralized user deposit
				assetTrans.Receiptant = ethereum.ToLowerAddressString(senderAddr.Hex())
				assetTrans.sender = assetTrans.Receiptant
				isFound = true
				candidateTransferLog = transferLog
				break
			} else {
				//Check if sending to centralized mainnet mini address
				sideChainDestAddr, err := dbGetUserSideChainWalletAddress(mainnetReceipant, me.context)
				if err != nil || sideChainDestAddr == nil {
					continue
				}
				assetTrans.Receiptant = ethereum.ToLowerAddressString(sideChainDestAddr.Hex())
				assetTrans.sender = ethereum.ToLowerAddressString(senderAddr.Hex())
				isFound = true
				candidateTransferLog = transferLog
				break
			}
		}
		if !isFound {
			return
		} else {
			//Filter to only candidate log object
			assetTrans.TransferLog = []types.Log{candidateTransferLog}
		}
	} else {
		senderAddr, _ := assetTrans.GetSender()
		assetTrans.AssetName = "ETH"
		if bytes.Equal(common.HexToAddress(senderAddr).Bytes(), me.config.SweepServiceInvokerAddress.Bytes()) {
			//Exclude sweep service invoker send ETH to other wallet
			receipt, err := me.scProcessor.QueryMainnetEthReceiptWithSetting(assetTrans.OriginalTransaction, 2, 10)
			//Ignore query receipt error, as it does not matter
			if err == nil && receipt.Status == 0 {
				log.GetLogger(me.context.LoggerName).Debugf("Receipt status is 0 for sweep trans: %s\r\n", assetTrans.Hash().Hex())
				return
			}

			DbInsertSweepTransaction(me.context, assetTrans)

			log.GetLogger(me.context.LoggerName).Debugln("Ignore sweep service invoker sending ETH to wallet transaction: ", assetTrans.Hash().Hex())
			return
		}

		if bytes.Equal(assetTrans.GetTo().Bytes(), me.scProcessor.MainnetUserDepositAddr.Bytes()) {
			//When sending ETH to EurusUserDeposit smart contract, it is assumed it is a decentralized user deposit
			assetTrans.Receiptant = senderAddr
		} else {
			//Check if sending to centralized mainnet mini address
			addr, err := dbGetUserSideChainWalletAddress(*assetTrans.GetTo(), me.context)
			if err != nil || addr == nil {
				return
			}
			assetTrans.Receiptant = ethereum.ToLowerAddressString(addr.Hex())
		}
	}

	receipt, err := me.scProcessor.QueryMainnetEthReceiptWithSetting(assetTrans.OriginalTransaction, 2, 10)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorf("Unable to get receipt for trans: %s Error: %s\r\n", assetTrans.Hash().Hex(), err.Error())
		return
	}

	if receipt.Status == 0 {
		log.GetLogger(me.context.LoggerName).Errorf("Receipt status is 0 for trans: %s\r\n", assetTrans.Hash().Hex())
		return
	}

	depositTx, rowAffected, err := DbInsertPendingDeposit(me.context, assetTrans, receipt, me.TransferEventID, isERC20)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Insert transaction failed after retry. Trans: ", assetTrans.Hash().Hex(), "  Error: ", err.Error())
		return
	}

	go func() {
		err = me.HandleMintEvent(assetTrans.Hash(), assetTrans.AssetName, common.HexToAddress(assetTrans.Receiptant), depositTx.Amount.BigInt())
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("Fail to handle mint event : ", err)
		}
	}()

	if rowAffected == 0 {
		//Another Deposit observer already processing, so skip this transaction
		log.GetLogger(me.context.LoggerName).Debugln("Another Deposit observer already processing, so skip sweep event handling for ", assetTrans.Hash().Hex())
		return
	}

	go func() {
		err = me.HandleSweepEvent(assetTrans.Hash())
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("Fail to handle sweep event", err)
		}
	}()
}

// func (me *DepositProcessor) processMintRequestEvent(transLog *types.Log) {
// 	//    event MintRequestEvent(string indexed mainnetDepositTransHash, uint256 indexed mintRequestTransId);
// 	log.GetLogger(me.context.LoggerName).Debugln("MintRequest event triggered")

// 	if transLog.Removed {
// 		log.GetLogger(me.context.LoggerName).Warnln("trans log removed. Mint request trans hash: ", transLog.TxHash.Hex())
// 		return
// 	}

// 	if len(transLog.Topics) < 2 {
// 		log.GetLogger(me.context.LoggerName).Errorln("Invalid topic count. Expected 2. For trans: ", transLog.TxHash.Hex())
// 		return
// 	}

// 	mintRequestIdByte := transLog.Topics[1].Bytes()
// 	mintRequestId := uint256.NewInt(0)
// 	mintRequestId = mintRequestId.SetBytes(mintRequestIdByte)

// 	args, err := ethereum.DefaultABIDecoder.DecodeABIEventData(transLog.Data, "PlatformWallet", "MintRequestEvent")
// 	if err != nil {
// 		log.GetLogger(me.context.LoggerName).Errorln("Unable to decode MintRequestEvent for mint request ID: ", mintRequestId.String())
// 		return
// 	}

// 	depositTransHashStr, ok := args[0].(string)
// 	if !ok {
// 		log.GetLogger(me.context.LoggerName).Errorln("Unable to convert deposit trans hash at MintRequestEvent for mint request ID: ", mintRequestId.String())
// 		return
// 	}

// 	depositTransHash := common.HexToHash(depositTransHashStr)

// 	receipt, err := me.scProcessor.ConfirmMint(mintRequestId)
// 	if err != nil {
// 		log.GetLogger(me.context.LoggerName).Errorf("Cannot confirm mint for deposit trans hash: %s, mint request Id: %s, Error: %s\r\n", depositTransHash.Hex(), mintRequestId.String(), err.Error())
// 	} else if err == nil && receipt == nil {
// 		log.GetLogger(me.context.LoggerName).Infoln("Mint request already confirmed. Deposit trans hash: ", depositTransHash.Hex(), " mint request id:", mintRequestId.String())
// 	} else if receipt.Status != 1 {
// 		receiptByte, _ := json.Marshal(receipt)
// 		log.GetLogger(me.context.LoggerName).Errorf("Confirm mint result failed. Deposit trans hash: %s, mint request id: %s, receipt: %s", depositTransHash.Hex(), mintRequestId.String(), string(receiptByte))

// 		var remarks asset.RemarksJson
// 		if receiptByte != nil {
// 			remarks = asset.NewRemarksJsonFromReceipt(receiptByte)
// 		} else {
// 			remarks = asset.NewRemarksJsonFromString("Confirm mint result failed")
// 		}

// 		err := dbUpdateDepositTransToError(me.context, depositTransHash, asset.DepositMintConfirming, remarks.String(), mintRequestId)
// 		if err != nil {
// 			log.GetLogger(me.context.LoggerName).Errorln("dbUpdateDepositTransToError failed on DepositMintConfirming error. Deposit trans hash: ", depositTransHash.Hex())
// 		}
// 	} else {
// 		dbUpdateDepositTransToMintConfirming(me.context, depositTransHash, *mintRequestId)

// 		log.GetLogger(me.context.LoggerName).Debugf("Confirm mint success. Deposit Trans Hash: %s, mint request id: %s\r\n", depositTransHash.Hex(), mintRequestId.String())
// 	}

// }

func (me *DepositProcessor) processMintCompletedEvent(transLog *types.Log) {

	log.GetLogger(me.context.LoggerName).Debugln("processMintCompleted event triggered")
	if len(transLog.Topics) < 2 {
		log.GetLogger(me.context.LoggerName).Errorln("Invalid topic count. Expected 2. For trans hash:", transLog.TxHash.Hex())
		return
	}
	if transLog.Removed {
		log.GetLogger(me.context.LoggerName).Warnln("trans log removed. Mint complete trans hash:", transLog.TxHash.Hex())
		return
	}
	var depositTransHash common.Hash = transLog.Topics[1]

	log.GetLogger(me.context.LoggerName).Debugln("processMintCompleted event triggered for deposit trans hash: ", depositTransHash.Hex())
	blockNum := big.NewInt(0)
	blockNum = blockNum.SetUint64(transLog.BlockNumber)

	blockTime, _ := me.scProcessor.GetSideChainBlockTimeFromBlockNumber(blockNum)

	rowAffected, err := dbUpdateDepositTransToCompleted(me.context, depositTransHash, transLog.TxHash, blockTime, transLog.Address)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("dbUpdateDepositTransToCompleted failed on deposit trans hash: ", depositTransHash.Hex())
		return
	}
	if rowAffected > 0 {
		user, depositTrans, err := dbGetUserDepositTransactionDetails(me.context, depositTransHash)
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("dbGetUserDepositTransactionDetails error: ", err, " user id: ", user.Id, " payment wallet address: ", user.WalletAddress)
			return
		}
		err = me.rewardProcessor.TransferRegistrationRewardToUser(depositTrans, user)
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("TransferRegistrationRewardToUser error: ", err, " user id: ", user.Id, " payment wallet address: ", user.WalletAddress)
			return
		}
	} else {
		log.GetLogger(me.context.LoggerName).Info("dbUpdateDepositTransToCompleted row affected is 0, another process already update the record. Deposit trans hash: ", depositTransHash.Hex())
	}

}

func (me *DepositProcessor) HandleMintEvent(depositTransHash common.Hash, assetName string, sideChainDestAddr common.Address, amount *big.Int) error {
	rowAffected, err := dbUpdateDepositTransToMintRequesting(me.context, depositTransHash)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("dbUpdateDepositTransToMintConfirming aborted: ", err.Error(), ". Deposit trans hash: ", depositTransHash.Hex())
		return err
	} else if rowAffected == 0 {
		log.GetLogger(me.context.LoggerName).Errorln("dbUpdateDepositTransToMintConfirming skip due to transaction already in next state. Deposit trans hash: ", depositTransHash.Hex())
		return nil
	}

	receipt, err := me.scProcessor.SubmitMintRequest(depositTransHash, assetName, sideChainDestAddr, amount)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Submit mint request failed: ", err.Error(), " Deposit trans hash: ", depositTransHash.Hex())
		remarks := asset.NewRemarksJsonFromError(err)
		_ = dbUpdateDepositTransToError(me.context, depositTransHash, asset.DepositMintRequesting, remarks.String())

		return err
	}
	if receipt.Status != 1 {
		receiptByte, _ := json.Marshal(receipt)
		log.GetLogger(me.context.LoggerName).Errorln("Submit mint receipt failed.  Deposit trans hash: ", depositTransHash.Hex(), " Receipt: ", string(receiptByte))
		remarks := asset.NewRemarksJsonFromReceipt(receiptByte)
		err := dbUpdateDepositTransToError(me.context, depositTransHash, asset.DepositMintRequesting, remarks.String())
		if err != nil {
			log.GetLogger(me.context.LoggerName).Errorln("dbUpdateDepositTransToError aborted. Deposit trans hash: ", depositTransHash.Hex())
		}
		return nil
	}
	return nil
}

func (me *DepositProcessor) HandleSweepEvent(depositTransHash common.Hash) error {
	user, depositTrans, err := dbGetUserDepositTransactionDetails(me.context, depositTransHash)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Cannot get back user info to start handle sweep event:", err, "trans hash:", depositTransHash.Hex())
		return err
	}

	addr := common.HexToAddress(depositTrans.MainnetToAddress)
	assetName := depositTrans.AssetName
	status := depositTrans.Status

	// Double check sweep transaction should not trigger sweep again
	if status == asset.SweepTrans {
		log.GetLogger(me.context.LoggerName).Errorln("Invalid transaction status:", status, "trans hash:", depositTransHash.Hex())
		return errors.Errorf("Attempting to handle sweep event on sweep transaction. trans hash: %v", depositTransHash.Hex())
	}

	assetSetting, found := me.config.AssetSettings[assetName]
	if !found {
		log.GetLogger(me.context.LoggerName).Errorln("Cannot do sweep token checking because of invalid asset:", assetName, "trans hash:", depositTransHash.Hex())
		return errors.Errorf("No sweep token setting for this asset: %v, trans hash: %v", assetName, depositTransHash.Hex())
	}

	balance, err := me.scProcessor.GetMainnetBalance(addr, assetName)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorf("Fail to query mainnet address %v balance of %v, user id: %v, decentralized user: %v, trans hash: %v\n", addr, assetName, user.Id, user.IsMetamaskAddr, depositTransHash.Hex())
		return err
	}

	if balance.Cmp(assetSetting.SweepTriggerAmount.BigInt()) < 0 {
		return nil
	}

	log.GetLogger(me.context.LoggerName).Infof("Mainnet address %v balance of %v is %v, which exceeds sweep token threshold %v\n", addr, assetName, balance, assetSetting.SweepTriggerAmount)

	var userID *uint64
	if user.IsMetamaskAddr {
		// No need to include decentralized user's user id
		userID = nil
	} else {
		userID = &user.Id
	}

	err = dbInsertPendingSweepWallet(me.context, userID, addr.Hex(), assetSetting.AssetName)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Error when inserting pending sweep wallet record:", err)
		return err
	}

	return nil
}
