package bc_indexer

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type blockChainProcessorContext struct {
	AssetAddressMap              Asset
	Db                           *database.Database
	SlaveDb                      *database.ReadOnlyDatabase
	Config                       *BlockChainIndexerConfigBase
	EthClient                    *ethereum.EthClient
	EurusUserDepositAddress      *common.Address
	MainnetPlatformWalletAddress *common.Address
	SweepServiceInvokerAddress   *common.Address
	MainnetPlatformWalletUser    *user.User
	SweepServiceInvokerUser      *user.User
	transferRewardProcessor      *TransferRewardProcessor
	IsMainnet                    bool
	LoggerName                   string
}

func TransactionsHandler(context *blockChainProcessorContext, block *types.Block) {
	transactions := block.Transactions()
	for _, transaction := range transactions {
		if context.IsMainnet {
			var err error
			isTransferTx, ext, _ := FilterERC20TransferTransaction(context, transaction, block)
			if !isTransferTx {
				ext = FilterSweepDecentralizedTransaction(context, transaction, block)
				if ext == nil {
					///Check if the transaction transfering ETH
					ext = FilterMainNetETHTransaction(context, transaction, block)
					if ext == nil {
						continue
					} else {
						//ETH transfer at mainnet
						err = ProcessExtractedTransaction(context, ext)
						if err != nil {
							log.GetLogger(context.LoggerName).Errorln("ProcessExtractedTransaction error: ", err)
						}
					}
				} else {
					err = ProcessExtractedTransaction(context, ext)
					if err != nil {
						log.GetLogger(context.LoggerName).Errorln("ProcessExtractedTransaction error: ", err)
					}
				}
			} else {
				//ERC20 transfer at mainnet
				ext = FilterEurusUserDepositTransaction(context, ext)
				if ext == nil {
					continue
				}
				err = ProcessExtractedTransaction(context, ext)
				if err != nil {
					log.GetLogger(context.LoggerName).Errorln("ProcessExtractedTransaction error: ", err)
				}
			}

		} else {
			isErc20Trans, ext, err := FilterEurusERC20TransferTransaction(context, transaction, block)
			if err != nil {
				log.GetLogger(context.LoggerName).Errorln("Process transcation failed: ", transaction.Hash().Hex(), " Error: ", err.Error())
				continue
			}

			if err == nil && !isErc20Trans {

				eunExt, err := FilterEUNTransaction(context, transaction, block)
				if eunExt != nil && err == nil {
					//Decentralized transfer EUN transaction
					err := ProcessExtractedTransaction(context, eunExt)
					if err != nil {
						log.GetLogger(context.LoggerName).Errorln(asset.EurusTokenName+" Tx Hash: "+transaction.Hash().Hex()+" Unable to extract details from transaction", err.Error())
					}
				} else if err == nil {
					ext, err := FilterCentralizedUserTransferTransaction(context, transaction, block)
					if err != nil {
						log.GetLogger(context.LoggerName).Errorln("Process FilterCentralizedUserTransferTransaction failed: ", transaction.Hash().Hex(), " Error: ", err.Error())
						continue
					}

					if ext == nil {
						ext, err = FilterCentralizedUserConfirmTransferTransaction(context, transaction, block)
						if err != nil {
							log.GetLogger(context.LoggerName).Errorln("Process FilterCentralizedUserConfirmTransferTransaction failed: ", transaction.Hash().Hex(), " Error: ", err.Error())
							continue
						}
						if ext == nil {
							ext, err = FilterTopUpTransaction(context, transaction, block)
							if err != nil {
								log.GetLogger(context.LoggerName).Errorln("Process FilterTopUpTransaction failed: ", transaction.Hash().Hex(), " Error: ", err.Error())
								continue
							}

							if ext == nil {

								purchaseExt, err := FilterEurusERC20MerchantTransaction(context, transaction, block)
								if err != nil {
									log.GetLogger(context.LoggerName).Errorln("Process FilterEurusERC20MerchantTransaction failed. Tx Hash: "+transaction.Hash().Hex()+" Error: ", err.Error())
									continue
								}
								if purchaseExt != nil {
									err := ProcessExtractedTransaction(context, purchaseExt)
									if err != nil {
										log.GetLogger(context.LoggerName).Errorln(asset.EurusTokenName+" Tx Hash: "+transaction.Hash().Hex()+" Unable to extract details from transaction", err.Error())
										continue
									}
								} else {
									//Not interested transaction
									continue
								}

							}
						}
					}
					//Centralized transferRequest transaction
					err = ProcessExtractedTransaction(context, ext)
					if err != nil {
						log.GetLogger(context.LoggerName).Errorln("Tx Hash: "+transaction.Hash().Hex()+" Unable to extract details from transaction", err.Error())
					}
				}

			} else if err == nil && isErc20Trans {

				//Decentralized ERC20 transaction
				err := ProcessExtractedTransaction(context, ext)
				if err != nil {
					log.GetLogger(context.LoggerName).Errorln("Tx Hash: "+transaction.Hash().Hex()+" Unable to extract details from transaction", err.Error())
				}
			}
		}
	}
}

func ProcessExtractedTransaction(context *blockChainProcessorContext, ext *ExtractedTransaction) error {
	var err error
	if context.IsMainnet && ext.AssetName == "ETH" {
		ext, _ = FilterMainNetSender(context, ext)
		ext, _ = FilterMainNetReceiver(context, ext)
	}

	switch ext.TransactionType {
	case asset.Transfer:
		if ext.FromUser != nil {
			//Exclude insert into DB if the transfer is sent to mainnet wallet address at mainnet transaction for the FromUser
			if !context.IsMainnet || ext.ToUser == nil || ext.ToUser.MainnetWalletAddress != context.MainnetPlatformWalletUser.MainnetWalletAddress {
				err = DbAddTxIndex(context, ext, true)
				if err != nil {
					return err
				}
			}
		}

		if (ext.ToUser != nil && ext.Status > 0) || ext.RequestTransId != nil {
			//Exclude invoker send ETH to user wallet transactions for the ToUser
			if context.IsMainnet && ext.FromUser != nil && ext.FromUser.MainnetWalletAddress == context.SweepServiceInvokerUser.MainnetWalletAddress {
				return nil
			}
			//Not require to insert DB for to user if the status is 0
			err = DbAddTxIndex(context, ext, false)

			if !context.IsMainnet {
				err1 := context.transferRewardProcessor.TransferRegistrationRewardToUser(ext, ext.ToUser)
				if err1 != nil {
					log.GetLogger(context.LoggerName).Errorln("Unable to process TransferRegistrationRewardToUser: ", err1)
				}
			}

			if err != nil {
				return err
			}

		}
	case asset.TopUp:
		if ext.FromUser == nil {
			log.GetLogger(context.LoggerName).Warnln("DbInsertTopUpTransaction null from user. Insert DB skipped. Trans hash: ", ext.TxHash)
			return nil
		}
		err = DbInsertTopUpTransaction(context.Db, ext)
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("DbInsertTopUpTransaction error: ", err, " Trans hash: ", ext.TxHash)
			return err
		}

	case asset.MerchantDeposit, asset.Purchase:
		err = DBInsertPurchaseTransaction(context.Db, ext)
		if err != nil {
			return err
		}
	}
	return nil
}
