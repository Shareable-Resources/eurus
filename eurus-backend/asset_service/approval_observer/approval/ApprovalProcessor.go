package approval

import (
	"bytes"
	"encoding/json"
	"eurus-backend/asset_service/asset"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/smartcontract/build/golang/contract"
	"eurus-backend/user_service/user_service/user"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type ApprovalProcessor struct {
	EthClient      *eurus_ethereum.EthClient
	Config         *ApprovalObserverConfig
	DbProcessor    *ApprovalDBProcessor
	ApprovalWallet *contract.ApprovalWallet
	LoggerName     string
}

var executedAlreadyErr = errors.New("Executed already")

func NewApprovalProcessor(ethClient *eurus_ethereum.EthClient, dbProcessor *ApprovalDBProcessor, config *ApprovalObserverConfig, loggerName string) *ApprovalProcessor {
	var err error
	processor := new(ApprovalProcessor)
	processor.EthClient = ethClient
	processor.DbProcessor = dbProcessor
	processor.LoggerName = loggerName
	processor.Config = config
	for {
		processor.ApprovalWallet, err = contract.NewApprovalWallet(processor.Config.ApprovalWalletAddress, processor.EthClient.Client)
		if err != nil {
			log.GetLogger(processor.LoggerName).Error("Enable to load ApprovalWallet Smart Contract ERROR : ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		log.GetLogger(processor.LoggerName).Println("Approval Address is : ", processor.Config.ApprovalWalletAddress.Hex())
		break
	}
	return processor
}

/*
func (me *ApprovalProcessor) Init(ethClient *eurus_ethereum.EthClient, dbProcessor *ApprovalDBProcessor, config *ApprovalObserverConfig, loggerName string) {
	var err error
	me.EthClient = ethClient
	me.DbProcessor = dbProcessor
	me.LoggerName = loggerName
	me.Config = config
	for {
		me.ApprovalWallet, err = contract.NewApprovalWallet(me.Config.ApprovalWalletAddress, me.EthClient.Client)
		if err != nil {
			log.GetLogger(me.LoggerName).Error("Enable to load ApprovalWallet Smart Contract ERROR : ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		log.GetLogger(me.LoggerName).Println("Approval Address is : ", me.Config.ApprovalWalletAddress.Hex())
		break
	}
}
*/
func (me *ApprovalProcessor) BlockHandler(block *types.Block) {
	if block.Transactions().Len() > 0 { //check if block has tx
		me.GetTransactionsFromBlock(block)
	}
}

func (me *ApprovalProcessor) GetTransactionsFromBlock(block *types.Block) {
	transactions := block.Transactions()
	//log.GetLogger(me.LoggerName).Debug("Transactions Received: ", block.Number(), "transaction is : ", block.Transactions())
	for _, transaction := range transactions {
		if len(transaction.Data()) < 1 {
			continue
		}

		if transaction.To() == nil {
			continue
		}

		//Check if the transaction send to ERC20 address
		var transToAddr = transaction.To().Hex()
		var destAddr common.Address
		var amount *big.Int
		var adminFee *big.Int = big.NewInt(0)

		assetName, err := me.GetCurrencyNameByAddrFromSC(transToAddr)
		if err != nil {
			log.GetLogger(me.LoggerName).Errorln("Failed to get currency name from address for transaction hash: ", transaction.Hash().Hex())
			continue
		}

		var senderUser *user.User = nil
		var senderAddress string

		if assetName == "" {
			// Check if the transaction send to UserWallet address
			senderUser, err = me.DbProcessor.DbGetCenteralizedUserByWalletAddress(transToAddr)
			senderAddress = transToAddr
			if err != nil {
				continue
			}
			if senderUser.Id == 0 {
				//User not found
				continue
			}
			//Centralized user withdraw logic flow
			//Parse the smart contract function call input arguments
			args, err, _ := eurus_ethereum.DefaultABIDecoder.DecodeABIInputArgument(transaction.Data(), "UserWallet", "submitWithdrawV1")
			if err != nil {
				//Not submitWithdraw function call
				continue
			}
			assetName = args["assetName"].(string)
			if assetName == "" {
				log.GetLogger(me.LoggerName).Errorln("Invalid assetName argument for transaction hash: ", transaction.Hash().Hex())
			}
			var ok bool
			destAddr, ok = args["dest"].(common.Address)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve dest address. trans hash: ", transaction.Hash().Hex())
				continue
			}

			amount, ok = args["withdrawAmount"].(*big.Int)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve withdraw amount. trans hash: ", transaction.Hash().Hex())
				continue
			}

			var amountWithFee *big.Int
			amountWithFee, ok = args["amountWithFee"].(*big.Int)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve amount with fee. trans hash: ", transaction.Hash().Hex())
				continue
			}
			adminFee = adminFee.Sub(amountWithFee, amount)

		} else {
			//Decentralized user logic flow
			args, err, _ := eurus_ethereum.DefaultABIDecoder.DecodeABIInputArgument(transaction.Data(), "EurusERC20", "submitWithdraw")
			if err != nil {
				//Not submitWithdraw function call
				continue
			}
			var ok bool
			destAddr, ok = args["dest"].(common.Address)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve dest address. trans hash: ", transaction.Hash().Hex())
				continue
			}

			amount, ok = args["withdrawAmount"].(*big.Int)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve withdraw amount. trans hash: ", transaction.Hash().Hex())
				continue
			}
			var amountWithFee *big.Int
			amountWithFee, ok = args["amountWithFee"].(*big.Int)
			if !ok {
				log.GetLogger(me.LoggerName).Errorln("Unable to retrieve amount with fee. trans hash: ", transaction.Hash().Hex())
				continue
			}
			adminFee = adminFee.Sub(amountWithFee, amount)

			sender, err := transaction.AsMessage(types.NewEIP155Signer(transaction.ChainId()), nil)
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Error to get the sender. The transaction hash :", transaction.Hash().Hex())
				continue
			}
			senderAddress = strings.ToLower(sender.From().Hex())
			senderUser, err = me.DbProcessor.DbGetDecentralizedUserByWalletAddress(sender.From().Hex())
			if err != nil {
				continue
			}

			if senderUser == nil || senderUser.Id == 0 {
				log.GetLogger(me.LoggerName).Errorln("Unable to get user from DB. sender address: ", sender.From().Hex(), " trans hash: ", transaction.Hash().Hex())
				continue
			}
		}

		ext := NewExtractedTransaction(transaction, assetName, senderUser, senderAddress, destAddr.Hex())
		if len(ext.OriginalTransaction.Data()) == 0 {
			continue
		}
		ext.Amount = amount
		ext.AdminFee = adminFee
		ext.TransDate = time.Unix(int64(block.Time()), 0)

		emptyByte := make([]byte, len(ext.OriginalTransaction.Data()))
		if bytes.Equal(ext.OriginalTransaction.Data(), emptyByte) {
			log.GetLogger(me.LoggerName).Errorln("Empty Byte received", emptyByte, ext.OriginalTransaction.Data())
			continue
		}

		log.GetLogger(me.LoggerName).Debugln("Transaction Hash : ", transaction.Hash().Hex(), "Transaction To address : ", transaction.To().Hex())

		if assetName != "" {
			receipt, err := me.EthClient.QueryEthReceiptWithSetting(transaction, 2, 30)
			if err != nil {
				log.GetLogger(me.LoggerName).Errorln("Query Receipt failed.", err.Error(), "Transaction hash is ", transaction.Hash().Hex(), "transaction to : ", transaction.To().Hex(), "Gas price : ", transaction.GasPrice(),
					"Original transaction : ", transaction)
				continue
			}

			var pendingPreWithdrawData *PendingPrewithdraw
			var duplicate = false
			var isInsertDb bool = false
			for i := 0; i < me.Config.GetRetryCount(); i++ {
				isFailed, remarks := me.IsWithdrawRequestFailed(receipt)

				if isFailed {
					err := me.DbProcessor.InsertWithdrawFail(ext, receipt, remarks)
					if err != nil {
						log.GetLogger(me.LoggerName).Errorln("Fail to insert fail data to DB", err.Error(), " trans hash: ", ext.TxHash)
						fmt.Println("Now retry ", i, " times")
						time.Sleep(me.Config.GetRetryInterval() * time.Second)
					} else {
						isInsertDb = true
						break
					}
				} else {

					pendingPreWithdrawData, err = me.DbProcessor.DbInsertPendingPrewithdraw(ext, receipt)
					if err != nil {
						isInsertDb = false
						if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
							log.GetLogger(me.LoggerName).Error(err.Error(), " The transaction hash is ", ext.TxHash)
							duplicate = true
							break
						} else {
							log.GetLogger(me.LoggerName).Errorln("Fail to pending the block. ", err.Error(), "Transaction hash is ", transaction.Hash().Hex(), "transaction to : ", transaction.To().Hex(), "Gas price : ", transaction.GasPrice(),
								"Original transaction : ", transaction)
							fmt.Println("Now retry ", i, " times")
							time.Sleep(me.Config.GetRetryInterval() * time.Second)
							duplicate = false
							continue
						}
					} else {
						isInsertDb = true
						break
					}
				}
			}

			if !isInsertDb {
				continue
			}

			if duplicate {
				continue
			}

			var comfirmTransaction *types.Transaction
			if pendingPreWithdrawData != nil {
				comfirmTransaction = me.AutoApprove(pendingPreWithdrawData)
				if comfirmTransaction != nil {
					err = me.getApproveReceipt(ext, pendingPreWithdrawData, comfirmTransaction)
					if err != nil {
						if err == executedAlreadyErr {
							continue
						}
						log.GetLogger(me.LoggerName).Errorln("Cannot get receipt: ", err.Error(), "Transaction hash is ", transaction.Hash().Hex(), "transaction to : ", transaction.To().Hex(), "Gas price : ", transaction.GasPrice(),
							"Original transaction : ", transaction)
					}
				}
			}
		}
	}
}

func (me *ApprovalProcessor) IsWithdrawRequestFailed(receipt *eurus_ethereum.BesuReceipt) (bool, asset.RemarksJson) {
	//duncan debug
	if receipt.Status == 0 {
		receiptData, _ := json.Marshal(receipt)
		return true, asset.NewRemarksJsonFromReceipt(receiptData)
	}
	return false, nil

	// abi := eurus_ethereum.DefaultABIDecoder.GetABI("ApprovalWallet")
	// submitWithdrawEvent, ok := abi.Events["SubmitWithdraw"]
	// if !ok {
	// 	log.GetLogger(me.LoggerName).Errorln("ApprovalWallet SubmitWithdraw event not found in ABI. Program aborted")
	// }

	// userProxyAbi := eurus_ethereum.DefaultABIDecoder.GetABI("UserWalletProxy")
	// submitWithdrawFailedEvent, ok := userProxyAbi.Events["SubmitWithdrawFailed"]
	// if !ok {
	// 	log.GetLogger(me.LoggerName).Errorln("UserWalletProxy SubmitWithdrawFailed event not found in ABI. Program aborted")
	// }

	// if len(receipt.Logs) == 0 {
	// 	return true, asset.NewRemarksJsonFromString("Topic not found")
	// }

	// for _, topicLog := range receipt.Logs {
	// 	if topicLog.Topics[0] == submitWithdrawEvent.ID {
	// 		return false, nil
	// 	} else if topicLog.Topics[0] == submitWithdrawFailedEvent.ID {

	// 		args, err := eurus_ethereum.DefaultABIDecoder.DecodeABIEventData(topicLog.Data, "UserWalletProxy", "SubmitWithdrawFailed")
	// 		if err != nil {

	// 			return true, asset.NewRemarksJsonFromError(err)
	// 		}
	// 		reasonByte, ok := args[2].([]byte)
	// 		if ok {
	// 			reason := string(reasonByte)
	// 			return true, asset.NewRemarksJsonFromString(reason)
	// 		}
	// 		return true, asset.NewRemarksJsonFromString("Reverted")
	// 	}
	// }
	// return true, asset.NewRemarksJsonFromString("Topic not found")
}

func (me *ApprovalProcessor) AutoApprove(todata *PendingPrewithdraw) *types.Transaction {
	if todata.CustomerType == 0 {
		transactionID := big.NewInt(0)
		transactionID.SetUint64(*todata.RequestTransId)
		tx, err := me.EthClient.InvokeSmartContract(me.Config, me.Config.HdWalletPrivateKey,
			me.Config.SideChainGasLimit,
			func(ethClient *eurus_ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
				tx, err := me.ApprovalWallet.ConfirmTransaction(transOpt, transactionID)
				if err != nil {
					errByte, _ := json.Marshal(err)
					log.GetLogger(me.LoggerName).Errorln("Fail to confirm the transaction", err.Error(), "Transaction ID is :", transactionID, "Error : ", string(errByte))
					return tx, false, err
				} else {
					log.GetLogger(me.LoggerName).Println("Transaction approve confirmed ", tx)
					return tx, true, err
				}
			})
		if err != nil {
			log.GetLogger(me.LoggerName).Errorln("Invoke Smart Contract Error.", err.Error())
			return nil
		}
		return tx
	}
	return nil
}

func (me *ApprovalProcessor) GetCurrencyNameByAddrFromSC(address string) (string, error) {

	var err error
	var assetName string
	for i := 0; i < me.Config.RetryCount; i++ {
		var externalSC *contract.ExternalSmartContractConfig
		externalSC, err = contract.NewExternalSmartContractConfig(common.HexToAddress(me.Config.ExternalSCConfigAddress), me.EthClient.Client)
		if err != nil {
			log.GetLogger(me.LoggerName).Errorln("Get new External Smart Contract config failed. ", err.Error(),
				" input address: ", address, " external sc config address: ", me.Config.ExternalSCConfigAddress)
			time.Sleep(me.Config.GetRetryInterval() * time.Second)
			continue
		}
		assetName, err = externalSC.GetErc20SmartContractByAddr(&bind.CallOpts{}, common.HexToAddress(address))
		if err != nil {
			log.GetLogger(me.LoggerName).Errorln("Get asset address failed.", err.Error(), " input address: ", address)
			break
		}
		return assetName, nil
	}
	return "", err
}

// func (me *ApprovalProcessor) TxDataHandler(ext *ExtractedTransaction) *ExtractedTransaction {
// 	txData := ext.OriginalTransaction.Data()
// 	ext = FilterTransfer(ext, txData)
// 	return ext
// }

// func (me *ApprovalProcessor) WalletAddressHandler(ext *ExtractedTransaction) error {
// 	senderUser, err := me.dbProcessor.ApprovalProcessorDBProcessor(ext.sender)
// 	if err != nil {
// 		return err
// 	}

// 	if senderUser.Id != 0 {
// 		ext.User = senderUser

// 	}

// 	if err != nil {
// 		return err
// 	}

// 	toUser, err := me.dbProcessor.DbGetCenteralizedUserByWalletAddress(ext.To)
// 	if err != nil {
// 		return err
// 	}

// 	if toUser.Id != 0 {
// 		ext.User = toUser

// 	}

// 	if err != nil {
// 		return err
// 	}
// 	ext.User = toUser
// 	return nil
// }

// func (me *ApprovalDBProcessor) ApprovalProcessorDBProcessor(address string) (*user.User, error) {
// 	const retryGetUserNum int = 5
// 	var err error
// 	var user *user.User
// 	for i := 0; i < retryGetUserNum; i++ {
// 		user, err = me.DbGetCenteralizedUserByWalletAddress(address)
// 		if err != nil && strings.Contains(err.Error(), "Database Network Error") {
// 			time.Sleep(me.Config.GetRetryInterval() * time.Second)
// 			continue
// 		} else if err != nil {
// 			return nil, err
// 		} else {
// 			break
// 		}
// 	}
// 	return user, err
// }

// func (me *ApprovalProcessor) extractTopic(topicLog interface{}) (topicIsTrue bool, err error) {
// 	arr := topicLog.([]common.Hash)
// 	eventTopicAddress := hex.EncodeToString(me.Config.ApprovalProcessorTopic.Bytes())

// 	for i := 0; i < len(arr); i++ {
// 		num := hex.EncodeToString(arr[i].Bytes())
// 		if num == eventTopicAddress {
// 			return true, nil
// 		}
// 	}
// 	return false, err
// }

// func (me *ApprovalProcessor) getTopic(key int, value *types.Log) (filterKey int, err error) {
// 	//for key, value := range topicLog {
// 	newTopic := value.Topics
// 	topicChecker, err := me.extractTopic(newTopic)
// 	if err != nil {
// 		log.GetLogger(me.LoggerName).Errorln("Topic is not exist in receipt : ", err.Error())
// 	}
// 	if topicChecker == true {
// 		return key, nil
// 	} else {
// 		return key, errors.New("The Topic not find")
// 	}
// 	//}
// }

func (me *ApprovalProcessor) getApproveReceipt(ext *ExtractedTransaction, pendingPreWithdrawData *PendingPrewithdraw, comfirmTransaction *types.Transaction) error {
	var receipt *eurus_ethereum.BesuReceipt
	var err error

	receipt, err = me.EthClient.QueryEthReceiptWithSetting(comfirmTransaction, me.Config.RetryInterval, me.Config.RetryCount)
	if err != nil {
		log.GetLogger(me.LoggerName).Errorln("Can not waiting the transaction receipt. Request Transaction : ", ext.OriginalTransaction, " Request Hash : ", ext.OriginalTransaction.Hash().Hex(), " comfirmTransaction hash: ", comfirmTransaction.Hash().Hex(), " Error msg : ", err.Error())
		err = me.DbProcessor.ApprovalFail(pendingPreWithdrawData)
		if err != nil {
			log.GetLogger(me.LoggerName).Errorln("Failed to update the transaction status to failed in database. err msg :", err.Error(), "Request Transaction ", ext.TxHash, " comfirmTransaction hash: ", comfirmTransaction.Hash().Hex())
		}

	} else {
		log.GetLogger(me.LoggerName).Infoln("Received Receipt. Status : ", receipt.Status, " Hash : ", receipt.TxHash.Hex())
		if receipt.Status != 1 {

			revertReasonHex := strings.TrimPrefix(receipt.RevertReason, "0x")
			revertReason := common.Hex2Bytes(revertReasonHex)

			if !strings.Contains(string(revertReason), "MultiSigWallet: Transaction already executed") {
				log.GetLogger(log.Name.Root).Errorln("Revert reason hex: ", receipt.RevertReason, " string: ", string(revertReason), " Request Transaction ", ext.TxHash, " tx hash: ", receipt.TxHash.Hex())
				err = me.DbProcessor.ApprovalFail(pendingPreWithdrawData)
				if err != nil {
					log.GetLogger(me.LoggerName).Errorln("Failed to update the transaction status to failed in database. err msg :", err.Error(), "Transaction: ", ext, "Hash : ", ext.TxHash)
					return err
				}
				return errors.New("Receipt status is 0")
			} else {
				log.GetLogger(log.Name.Root).Warnln("Request Transaction ", ext.TxHash, " approval executed already. tx hash: ", receipt.TxHash.Hex())
				return executedAlreadyErr
			}
		} else {
			if pendingPreWithdrawData != nil && comfirmTransaction != nil {
				err = me.DbProcessor.PendingSuccess(pendingPreWithdrawData, comfirmTransaction, receipt)
				if err != nil {
					log.GetLogger(me.LoggerName).Errorln("Fail to insert data to db", err.Error(), "hash : ", pendingPreWithdrawData.RequestTransHash)
				}
				receiptByte, _ := json.Marshal(receipt)
				log.GetLogger(me.LoggerName).Debugln("Success to receive the receive receipt.", string(receiptByte))
			} else {
				log.GetLogger(me.LoggerName).Debugln("The receipt data is nil. It maybe already handle by another observer. The Transaction hash is : ", comfirmTransaction.Hash().Hex())
				return err
			}
		}

	}

	return err
}
