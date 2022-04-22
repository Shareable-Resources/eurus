package sweep

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/secret"
	"eurus-backend/user_service/user_service/user"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type SweepServiceProcessor struct {
	config      *SweepServiceConfig
	context     *SweepServiceContext
	scProcessor *SweepServiceSCProcessor
}

type sweepItem struct {
	isFinish bool
	userId   uint64
	address  common.Address
	symbol   string
	wg       *sync.WaitGroup
}

func NewSweepServiceProcessor(config *SweepServiceConfig, context *SweepServiceContext, scProcessor *SweepServiceSCProcessor) *SweepServiceProcessor {
	processor := new(SweepServiceProcessor)
	processor.config = config
	processor.context = context
	processor.scProcessor = scProcessor
	return processor
}

func (ssp *SweepServiceProcessor) PollDatabase() {
	wallets, err := DBGetPendingSweepWallets(ssp.context)
	if err != nil {
		ssp.context.logger.Errorln("Fail to poll pending sweep wallets from database:", err)
		return
	}

	for _, wallet := range wallets {
		info := fmt.Sprintf("address: %v, asset: %v", wallet.MainnetWalletAddress, wallet.AssetName)

		ssp.context.logger.Infoln("Discovered pending sweep wallet,", info)

		prevGasFeeCap, prevGasTipCap, prevGasLimit, err := ssp.processPendingSweepWallet(&wallet)
		if err == nil {
			ssp.context.logger.Infoln("Successfully processed wallet,", info)
			continue
		}

		balance, err := ssp.scProcessor.getBalance(common.HexToAddress(wallet.MainnetWalletAddress), wallet.AssetName)
		if err != nil {
			ssp.context.logger.Error("Failed to get", wallet.AssetName, "balance of address", wallet.MainnetWalletAddress, " error: ", err)
			if wallet.UserID != nil {
				ssp.context.logger.Errorln(", userID:", wallet.UserID)
			} else {
				ssp.context.logger.Errorln("")
			}
		}

		if err != nil || balance.Cmp(big.NewInt(0)) > 0 {
			ssp.context.logger.Infoln("Failed to process wallet, will requeue it and retry:", err, ",", info)

			err = DBRequeuePendingSweepWallet(ssp.context, wallet.ID, prevGasFeeCap, prevGasTipCap, prevGasLimit)
			if err == nil {
				continue
			}
		}

		ssp.context.logger.Infoln("Failed to requeue wallet:", err, ",", info)
	}
}

func (ssp *SweepServiceProcessor) CheckNeedForSweep(users []user.User) error {
	const DefaultERC20Workers = 32
	const DefaultETHWorkers = 16

	erc20Workers := ssp.config.SweepERC20Workers
	if erc20Workers <= 0 {
		ssp.context.logger.Warnln("Invalid SweepERC20Workers in config, will use default value instead")
		erc20Workers = DefaultERC20Workers
	}

	ethWorkers := ssp.config.SweepETHWorkers
	if ethWorkers <= 0 {
		ssp.context.logger.Warnln("Invalid ETHWorkers in config, will use default value instead")
		ethWorkers = DefaultETHWorkers
	}

	// Because maximum capacity is known, it is easier and faster to use buffered channel instead of slice
	// Also no need to take care thread safety
	erc20Queue := make(chan *sweepItem, len(users)*len(ssp.config.AssetSettings)+erc20Workers)
	ethQueue := make(chan *sweepItem, len(users)+ethWorkers)

	ssp.context.logger.Infoln(len(users), "user(s) of", len(ssp.config.AssetSettings), "asset(s) will be checked")
	//Adding Eurus user deposit address to sweep user list
	users = append(users, user.User{Id: 0, MainnetWalletAddress: ethereum.ToLowerAddressString(ssp.scProcessor.eurusUserDepositAddress.Hex())})
	for _, u := range users {
		// Double check to make sure user mainnet wallet address is valid
		if !common.IsHexAddress(u.MainnetWalletAddress) {
			ssp.context.logger.Errorln("UserID", u.Id, "has invalid mainnet wallet address", u.MainnetWalletAddress, ", skip checking")
			continue
		}

		// This wait group is to force ETH checking must be after other tokens' checking finished
		address := common.HexToAddress(u.MainnetWalletAddress)
		wg := sync.WaitGroup{}

		for symbol := range ssp.config.AssetSettings {
			item := sweepItem{isFinish: false, userId: u.Id, address: address, symbol: symbol, wg: &wg}

			// Check balance of assets can be done parallelly, except ETH itself, this is to prevent wasting extra transaction fee
			// ETH checking need to be done after other tokens' checking are finished
			if symbol == "ETH" {
				ethQueue <- &item
			} else {
				wg.Add(1)
				erc20Queue <- &item
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(erc20Workers + ethWorkers)

	for i := 0; i < erc20Workers; i++ {
		erc20Queue <- &sweepItem{isFinish: true}
		go ssp.processNeedForSweepQueue(erc20Queue, &wg)
	}

	for i := 0; i < ethWorkers; i++ {
		ethQueue <- &sweepItem{isFinish: true}
		go ssp.processNeedForSweepQueue(ethQueue, &wg)
	}

	wg.Wait()
	ssp.context.logger.Infoln("All workers for sweep checking are finished, the whole process is completed")
	return nil
}

func (ssp *SweepServiceProcessor) processPendingSweepWallet(wallet *asset.PendingSweepWallet) (gasFeeCap *big.Int, gasTipCap *big.Int, gasLimit *uint64, err error) {
	gasFeeCap = big.NewInt(0)
	gasTipCap = big.NewInt(0)
	l := uint64(0)
	gasLimit = &l

	// AssetName in asset.PendingSweepWallet is ethereum, tether, etc
	// It will be more convenient to use symbol (ETH, USDT, etc) later, so look up the symbol first
	// symbol, found := ssp.config.CurrencyToSymbol[wallet.AssetName]
	// if !found {
	// 	ssp.context.logger.Errorln("Unknown asset:", wallet.AssetName)
	// 	return nil, nil, nil, fmt.Errorf("Unknown asset: %v", wallet.AssetName)
	// }
	symbol := wallet.AssetName

	address := common.HexToAddress(wallet.MainnetWalletAddress)

	if wallet.UserID != nil {
		// Just no need to dereference every time
		userID := *wallet.UserID

		// If user id is given, this is centralized user address
		// Try to get the user private key first, because if failed to get it, then can prevent loss of sending top-up transaction fee
		privateKeyECDSA, _, err := secret.GenerateCentralizedUserMainnetPrivateKey(ssp.config.CentralizedUserWalletMnemonicPhase, userID)
		if err != nil {
			ssp.context.logger.Errorln("Unable to get user private key, stop any further actions. User id:", userID, ",", err)
			return nil, nil, nil, err
		}

		userPrivateKey := common.Bytes2Hex(crypto.FromECDSA(privateKeyECDSA))

		// For a retrying wallet, previous gas cost may be given
		// In this case, use the previous values as base, add extra cap and limit and try again
		if wallet.PreviousGasFeeCap != nil && wallet.PreviousGasTipCap != nil && wallet.PreviousGasLimit != nil {
			gasFeeCap.Set(wallet.PreviousGasFeeCap.BigInt())
			gasFeeCap.Add(gasFeeCap, big.NewInt(ssp.config.SweepExtraGasFee))

			gasTipCap.Set(wallet.PreviousGasTipCap.BigInt())
			gasTipCap.Add(gasTipCap, big.NewInt(ssp.config.SweepExtraGasFee))

			*gasLimit = *wallet.PreviousGasLimit
			*gasLimit += ssp.config.SweepExtraGasLimit
		} else {
			gfc, gtc, gl, err := ssp.scProcessor.estimateSweepCost(address, userPrivateKey, symbol)
			if err != nil {
				ssp.context.logger.Errorln("Unable to estimate sweep cost. User id:", userID, ",", err)
				return nil, nil, nil, err
			}

			gasFeeCap.Set(gfc)
			gasTipCap.Set(gtc)

			// The gas estimate maybe too accurate, reserve more
			// Here take the closest number which is larger than the estimated and divisible by 1000
			n := uint64(0)
			if gl%1000 == 0 {
				n = 1
			} else {
				n = 0
			}
			n = (n + 1 + gl/1000) * 1000
			*gasLimit = n
		}

		// Address may not have sufficient ETH for transaction fee, in this case, top-up some ETH first
		// If sweeping ETH, suppose balance already covered the transaction fee so just skip this step
		if symbol != "ETH" {
			userETHBalance, err := ssp.scProcessor.getBalance(address, "ETH")
			if err != nil {
				// The error is not related to gas estimation so don't return them
				ssp.context.logger.Errorln("Unable to get ETH balance of address", address.Hex(), ", cannot determine if top-up transaction fee is needed. User id:", userID, ",", err)
				return nil, nil, nil, err
			}

			txFee := new(big.Int).Mul(gasFeeCap, big.NewInt(int64(*gasLimit)))

			if txFee.Cmp(userETHBalance) > 0 {
				ssp.context.logger.Infoln("ETH balance of address", address.Hex(), "is", inUnitEther(userETHBalance), "ETH, but transaction fee is estimated to be", inUnitEther(txFee), "ETH, need to top-up before sweeping. User id:", userID)

				topUpAmount := new(big.Int).Sub(txFee, userETHBalance)

				var topUpFeeCap *big.Int = nil
				var topUpTipCap *big.Int = nil
				const retry = 3

				// Do retry and increase gas cost every time, gas limit can be fixed to 21000
				for i := 0; i < retry; i++ {
					tufc, tutc, err := ssp.scProcessor.topUpTransactionFee(address, topUpAmount, topUpFeeCap, topUpTipCap)
					if err == nil {
						break
					}

					if tufc != nil {
						topUpFeeCap = tufc.Add(tufc, big.NewInt(ssp.config.SweepExtraGasFee))
					}

					if tutc != nil {
						topUpTipCap = tutc.Add(tutc, big.NewInt(ssp.config.SweepExtraGasFee))
					}
				}

				if err != nil {
					ssp.context.logger.Errorln("Unable to top-up address", address.Hex(), "after", retry, "time(s) retry, abort the process. User id:", userID, ",", err)
					return nil, nil, nil, err
				}

				ssp.context.logger.Infoln("Successfully top-up", inUnitGwei(topUpAmount), "Gwei to address", address.Hex(), ". User id:", userID)
			} else {
				ssp.context.logger.Infoln("ETH balance of address", address.Hex(), "is", inUnitEther(userETHBalance), "ETH, sufficient for sweeping transaction fee. User id:", userID)
			}
		}

		err = ssp.scProcessor.sweepCentralizedUserWalletAddress(address, userPrivateKey, symbol, gasFeeCap, gasTipCap, gasLimit, wallet.ID)
		if err != nil {
			// Error may be due to incorrect gas estimation so return it here
			ssp.context.logger.Errorln("Unable to sweep address", address.Hex(), "of asset", symbol, ". User id:", userID, ",", err)
			return gasFeeCap, gasTipCap, gasLimit, err
		}

		ssp.context.logger.Infoln("Successfully swept address", address.Hex(), "of", symbol)
	} else {
		if address != *ssp.scProcessor.eurusUserDepositAddress {
			ssp.context.logger.Warnln("Decentralized user's deposit address does not match the actual one, given:", address.Hex(), "actual:", ssp.scProcessor.eurusUserDepositAddress)
		}

		err := ssp.scProcessor.sweepEurusUserDeposit(symbol, wallet.ID)
		if err != nil {
			ssp.context.logger.Errorln("Failed to sweep EurusUserDeposit, asset:", symbol, ",", err)
			return nil, nil, nil, err
		}

		ssp.context.logger.Infoln("Successfully swept EurusUserDeposit of", symbol)
	}

	return gasFeeCap, gasTipCap, gasLimit, nil
}

func (ssp *SweepServiceProcessor) processNeedForSweepQueue(ch chan *sweepItem, wg *sync.WaitGroup) {
	for item := range ch {
		if item.isFinish {
			wg.Done()
			return
		}

		// ETH should wait for other token finish
		if item.symbol == "ETH" {
			item.wg.Wait()
		}

		ssp.processNeedForSweep(item.userId, item.address, item.symbol)

		if item.symbol != "ETH" {
			item.wg.Done()
		}
	}
}

func (ssp *SweepServiceProcessor) processNeedForSweep(userID uint64, address common.Address, symbol string) {
	balance, err := ssp.scProcessor.getBalance(address, symbol)
	if err != nil {
		ssp.context.logger.Errorln("Failed to get", symbol, "balance of address", address, ", userID:", userID, ",", err)
		return
	}

	assetSetting := ssp.config.AssetSettings[symbol]

	if balance.Cmp(assetSetting.SweepTriggerAmount.BigInt()) < 0 {
		return
	}

	ssp.context.logger.Infoln(symbol, "balance of address", address, "is", balance, ", which exceeds sweep token threshold", assetSetting.SweepTriggerAmount, ", will put it to pending")

	var userIdOrNull *uint64
	if userID > 0 {
		userIdOrNull = &userID
	}
	err = DBInsertPendingSweepWallet(ssp.context, userIdOrNull, address.Hex(), assetSetting.AssetName)
	if err != nil {
		ssp.context.logger.Errorln("Error when inserting pending sweep wallet record:", err)
	}
}

func inUnitGwei(wei *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(wei, -9)
}

func inUnitEther(wei *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(wei, -18)
}
