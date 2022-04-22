package sweep

import (
	"context"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type AssetAddressMap map[common.Address]string

type AssetNameMap map[string]common.Address

type SweepServiceSCProcessor struct {
	config                     *SweepServiceConfig
	context                    *SweepServiceContext
	mainnetEthClient           *ethereum.EthClient
	mainnetAssetInfo           AssetAddressMap
	mainnetAssetNameToAddress  AssetNameMap
	eurusUserDepositAddress    *common.Address
	eurusPlatformWalletAddress *common.Address
}

func NewSweepServiceSCProcessor(config *SweepServiceConfig, context *SweepServiceContext) *SweepServiceSCProcessor {
	processor := new(SweepServiceSCProcessor)
	processor.config = config
	processor.context = context
	return processor
}

func (scp *SweepServiceSCProcessor) Init() error {
	scp.mainnetEthClient = &ethereum.EthClient{
		Protocol: scp.config.MainnetEthClientProtocol,
		IP:       scp.config.MainnetEthClientIP,
		Port:     scp.config.MainnetEthClientPort,
		ChainID:  big.NewInt(int64(scp.config.MainnetEthClientChainID)),
	}
	_, err := scp.mainnetEthClient.Connect()
	if err != nil {
		return errors.WithMessage(err, "Connect mainnet failed")
	}

	// Asset contracts are needed when sweeping centralized user ERC20 tokens
	scp.mainnetAssetInfo, scp.mainnetAssetNameToAddress, err = scp.queryMainnetAssetContractInfo()
	if err != nil {
		return err
	}

	for assetName, assetAddr := range scp.mainnetAssetNameToAddress {
		fmt.Println("asset name:", assetName, "addr:", assetAddr.Hex())
		scp.context.logger.Infof("%s: %s", assetName, assetAddr.Hex())
	}

	scp.eurusUserDepositAddress, err = scp.queryDecentralizedUserDepositAddress()
	if err != nil {
		return err
	}

	fmt.Println("Address of EurusUserDeposit is", scp.eurusUserDepositAddress)
	scp.context.logger.Infoln("Address of EurusUserDeposit is", scp.eurusUserDepositAddress.Hex())

	scp.eurusPlatformWalletAddress, err = scp.queryPlatformWalletAddress()
	if err != nil {
		return err
	}

	fmt.Println("Address of EurusPlatformWallet is", scp.eurusPlatformWalletAddress)
	scp.context.logger.Infoln("Address of EurusPlatformWallet is", scp.eurusPlatformWalletAddress.Hex())

	return nil
}

func (scp *SweepServiceSCProcessor) queryMainnetAssetContractInfo() (AssetAddressMap, AssetNameMap, error) {
	scConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(scp.config.EurusInternalConfigAddress), scp.mainnetEthClient.Client)
	if err != nil {
		return nil, nil, err
	}

	assetNameList, addrList, err := scConfig.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		return nil, nil, err
	}

	assetInfo := make(AssetAddressMap)
	assetNameToAddress := make(AssetNameMap)

	for i, assetName := range assetNameList {
		assetInfo[addrList[i]] = assetName
		assetNameToAddress[assetName] = addrList[i]
	}

	return assetInfo, assetNameToAddress, nil
}

func (scp *SweepServiceSCProcessor) queryDecentralizedUserDepositAddress() (*common.Address, error) {
	scConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(scp.config.EurusInternalConfigAddress), scp.mainnetEthClient.Client)
	if err != nil {
		return nil, err
	}

	depositAddr, err := scConfig.EurusUserDepositAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &depositAddr, err
}

func (scp *SweepServiceSCProcessor) queryPlatformWalletAddress() (*common.Address, error) {
	scConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(scp.config.EurusInternalConfigAddress), scp.mainnetEthClient.Client)
	if err != nil {
		return nil, err
	}

	walletAddr, err := scConfig.PlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &walletAddr, err
}

func (scp *SweepServiceSCProcessor) getBalance(address common.Address, symbol string) (*big.Int, error) {
	if symbol == "ETH" {
		return scp.mainnetEthClient.GetBalance(address)
	}

	erc20, err := scp.getEurusERC20Contract(symbol)
	if err != nil {
		return nil, err
	}

	return erc20.BalanceOf(&bind.CallOpts{}, address)
}

func (scp *SweepServiceSCProcessor) getEurusERC20Contract(symbol string) (*contract.EurusERC20, error) {
	contractAddress, found := scp.mainnetAssetNameToAddress[symbol]
	if !found {
		return nil, errors.Errorf("Failed to find the contract addess of asset %v", symbol)
	}

	erc20, err := contract.NewEurusERC20(contractAddress, scp.mainnetEthClient.Client)
	if err != nil {
		return nil, err
	}

	return erc20, nil
}

func (scp *SweepServiceSCProcessor) topUpTransactionFee(address common.Address, topUpAmount *big.Int, designatedGasFeeCap *big.Int, designatedGasTipCap *big.Int) (gasFeeCap *big.Int, gasTipCap *big.Int, err error) {
	_, tx, err := scp.mainnetEthClient.TransferETHInEIP1559(scp.config.InvokerPrivateKey, address.Hex(), topUpAmount, false, false, false, designatedGasFeeCap, designatedGasTipCap, nil)
	if err != nil {
		scp.context.logger.Errorln("Failed to top-up address", address.Hex(), ":", err)
		return nil, nil, err
	}

	receipt, err := scp.mainnetEthClient.QueryEthReceiptWithSetting(tx, 1, scp.config.QueryReceiptRetryCount)
	if err != nil {
		scp.context.logger.Errorln("Failed to query receipt, trans hash:", tx.Hash().Hex(), ",", err)
		return tx.GasFeeCap(), tx.GasTipCap(), err
	}

	scp.databaseProcessAfterReceipt(tx, receipt, "Top-up", nil)

	if receipt.Status == 0 {
		scp.context.logger.Errorln("Receipt returns status", receipt.Status)
		return tx.GasFeeCap(), tx.GasTipCap(), errors.Errorf("Receipt returns status %v", receipt.Status)
	}

	return tx.GasFeeCap(), tx.GasTipCap(), nil
}

func (scp *SweepServiceSCProcessor) estimateSweepCost(address common.Address, privateKey string, symbol string) (gasFeeCap *big.Int, gasTipCap *big.Int, gasLimit uint64, err error) {
	// Run the transaction but without really send it to estimate gas usage
	tx, _, err := scp._sweepAddress(address, privateKey, symbol, true, nil, nil, nil, nil)
	if err == nil {
		return tx.GasFeeCap(), tx.GasTipCap(), tx.Gas(), nil
	}

	scp.context.logger.Warnln("Failed to estimate the cost of sweeping", address.Hex(), "try to estimate without dry-running transaction, asset:", symbol)

	baseFee, err := scp.mainnetEthClient.GetBaseFee()
	if err != nil {
		scp.context.logger.Errorln("Failed to get base fee, address:", address.Hex(), ", asset:", symbol, ",", err)
		return nil, nil, 0, err
	}

	gasTipCap, err = scp.mainnetEthClient.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		scp.context.logger.Errorln("Failed to get gas tip cap, address:", address.Hex(), ", asset:", symbol, ",", err)
		return nil, nil, 0, err
	}

	gasFeeCap = new(big.Int).Add(gasTipCap, new(big.Int).Mul(baseFee, big.NewInt(2)))

	return gasFeeCap, gasTipCap, ethereum.ETHTransferStandardGasLimit * 2, err
}

func (scp *SweepServiceSCProcessor) sweepCentralizedUserWalletAddress(address common.Address, privateKey string, symbol string, designatedGasFeeCap *big.Int, designatedGasTipCap *big.Int, designatedGasLimit *uint64, walletIDToDelete uint64) error {
	_, _, err := scp._sweepAddress(address, privateKey, symbol, false, designatedGasFeeCap, designatedGasTipCap, designatedGasLimit, &walletIDToDelete)
	return err
}

func (scp *SweepServiceSCProcessor) _sweepAddress(address common.Address, privateKey string, symbol string, estimateGasOnly bool, designatedGasFeeCap *big.Int, designatedGasTipCap *big.Int, designatedGasLimit *uint64, walletIDToDelete *uint64) (*types.Transaction, *ethereum.BesuReceipt, error) {
	var tx *types.Transaction
	var err error

	balance, err := scp.getBalance(address, symbol)
	if err != nil {
		if !estimateGasOnly {
			scp.context.logger.Errorln("Failed to get", symbol, "balance of address", address.Hex(), ",", err)
		}
		return nil, nil, err
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		if !estimateGasOnly {
			scp.context.logger.Errorln("Address", address.Hex(), "has 0", symbol, ", skip the sweeping process")
		}
		return nil, nil, errors.Errorf("Already no %v in address", symbol)
	}

	// ETH Amount to be swept will be fewer because some are used as transaction fee
	amount := new(big.Int).Set(balance)

	if symbol == "ETH" {
		if !estimateGasOnly {
			scp.context.logger.Infoln(amount, "(or fewer) from address", address.Hex(), "will be swept to EurusPlatformWallet, asset:", symbol)
		}

		var transfered *big.Int
		transfered, tx, err = scp.mainnetEthClient.TransferETHInEIP1559(privateKey, scp.eurusPlatformWalletAddress.Hex(), amount, true, true, estimateGasOnly, designatedGasFeeCap, designatedGasTipCap, designatedGasLimit)
		if err != nil {
			if !estimateGasOnly {
				scp.context.logger.Errorln("Failed to transfer ETH, address:", address.Hex(), "asset:", symbol, ",", err)
			}
			return nil, nil, err
		}

		// Actual transfer ETH is fewer, because some are spent in transaction fee
		amount.Set(transfered)
	} else {
		if !estimateGasOnly {
			scp.context.logger.Infoln(amount, "from address", address.Hex(), "will be swept to platform wallet, asset:", symbol)
		}

		erc20, err := scp.getEurusERC20Contract(symbol)
		if err != nil {
			if !estimateGasOnly {
				scp.context.logger.Errorln("Failed to init contract instance to perform actions, address:", address.Hex(), ", asset:", symbol, ",", err)
			}
			return nil, nil, err
		}

		txOps, err := scp.mainnetEthClient.GetNewTransactorFromPrivateKey(privateKey, scp.mainnetEthClient.ChainID)
		if err != nil {
			if !estimateGasOnly {
				scp.context.logger.Errorln("Failed to create transactor of address", address.Hex(), ", asset:", symbol, ",", err)
			}
			return nil, nil, err
		}

		// Caller can override gas settings, this is only effective when really doing the transaction
		if !estimateGasOnly {
			if designatedGasFeeCap != nil {
				txOps.GasFeeCap = designatedGasFeeCap
			}

			if designatedGasTipCap != nil {
				txOps.GasTipCap = designatedGasTipCap
			}

			if designatedGasLimit != nil {
				txOps.GasLimit = *designatedGasLimit
			}
		}

		txOps.NoSend = estimateGasOnly

		tx, err = erc20.Transfer(txOps, *scp.eurusPlatformWalletAddress, amount)
		if err != nil {
			if !estimateGasOnly {
				scp.context.logger.Errorln("Failed to transfer ERC20 token, address:", address.Hex(), "asset:", symbol, ",", err)
			}
			return nil, nil, err
		}
	}

	if estimateGasOnly {
		// Because transaction is actually not sent, query receipt must be fail, caller also just want to know gas estimate
		return tx, nil, nil
	}

	receipt, err := scp.mainnetEthClient.QueryEthReceiptWithSetting(tx, 1, scp.config.QueryReceiptRetryCount)
	if err != nil {
		scp.context.logger.Errorln("Failed to query receipt, trans hash:", tx.Hash().Hex(), ",", err)
		return nil, nil, err
	}

	scp.databaseProcessAfterReceipt(tx, receipt, "Sweep", walletIDToDelete)

	if receipt.Status == 0 {
		scp.context.logger.Errorln("Receipt returns status", receipt.Status)
		return nil, nil, errors.Errorf("Receipt returns status %v", receipt.Status)
	}

	scp.context.logger.Infoln("Transferred", amount, "from", address.Hex(), "to EurusPlatformWallet", scp.eurusPlatformWalletAddress, ", asset:", symbol)
	return tx, receipt, nil
}

func (scp *SweepServiceSCProcessor) sweepEurusUserDeposit(symbol string, walletIDToDelete uint64) error {
	if symbol == "ETH" {
		scp.context.logger.Infoln("Any ETH sent to EurusUserDeposit are auto forwarded to EurusPlatformWallet, nothing to do here")
		err := DBDeletePendingSweepWallet(scp.context, walletIDToDelete)
		if err != nil {
			scp.context.logger.Warnln("Failed to delete pending sweep wallet, id:", walletIDToDelete)
		}
		return nil
	}

	balance, err := scp.getBalance(*scp.eurusUserDepositAddress, symbol)
	if err != nil {
		scp.context.logger.Errorln("Failed to get", symbol, "balance of EurusUserDeposit,", err)
		return err
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		scp.context.logger.Errorln("EurusUserDeposit has 0", symbol, ", skip the sweeping process")
		return errors.Errorf("Already no %v in EurusUserDeposit", symbol)
	}

	scp.context.logger.Infoln(balance, "from EurusUserDeposit will be swept to platform wallet, asset:", symbol)

	sc, err := mainnet_contract.NewEurusUserDeposit(*scp.eurusUserDepositAddress, scp.mainnetEthClient.Client)
	if err != nil {
		scp.context.logger.Errorln("Cannot get the instance of EurusUserDeposit:", err)
		return err
	}

	// Before doing sweeping, make sure invoker is the writer of smart contract, just to prevent wasting transaction fee
	isWriter, err := sc.IsWriter(&bind.CallOpts{}, scp.config.InvokerAddress)
	if err != nil {
		scp.context.logger.Errorln("Cannot determine if invoker is one of the writer of EurusUserDeposit", err)
		return err
	}

	if !isWriter {
		scp.context.logger.Errorln("Invoker is not EurusUserDeposit writer, so cannot process further action")
		return errors.Errorf("Not EurusUserDeposit writer")
	}

	tx, err := scp.mainnetEthClient.InvokeSmartContract(scp.config, scp.config.InvokerPrivateKey, 0, func(ethClient *ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
		// tranbsactionHash is not used in smart contract currently, so just keep it empty
		tx, err := sc.Sweep(transOpt, [32]byte{}, scp.config.InvokerAddress, symbol, balance)
		return tx, false, err
	})

	if err != nil {
		scp.context.logger.Errorln("Failed to invoke smart contract sweeping", err)
		return err
	}

	receipt, err := scp.mainnetEthClient.QueryEthReceiptWithSetting(tx, 1, scp.config.QueryReceiptRetryCount)
	if err != nil {
		scp.context.logger.Errorln("Failed to query receipt, trans hash:", tx.Hash().Hex(), ",", err)
		return err
	}

	scp.databaseProcessAfterReceipt(tx, receipt, "Sweep", &walletIDToDelete)

	if receipt.Status == 0 {
		scp.context.logger.Errorln("Receipt returns status", receipt.Status)
		return errors.Errorf("Receipt returns status %v", receipt.Status)
	}

	return nil
}

func (scp *SweepServiceSCProcessor) databaseProcessAfterReceipt(tx *types.Transaction, receipt *ethereum.BesuReceipt, allocationType string, walletIDToDelete *uint64) {
	err := DBInsertAssetAllocationCost(scp.context, tx.Hash(), allocationType, receipt.GasUsed, receipt.EffectiveGasPrice)
	if err != nil {
		// Just give warning message, there are always ways to get back transaction details so don't let it block the process
		scp.context.logger.Warnln("Failed to insert asset allocation cost, txHash:", tx.Hash())
	}

	if walletIDToDelete == nil {
		return
	}

	// Failed sweeping should be requeued
	if receipt.Status == 0 {
		return
	}

	err = DBDeletePendingSweepWallet(scp.context, *walletIDToDelete)
	if err != nil {
		scp.context.logger.Warnln("Failed to delete pending sweep wallet, id:", walletIDToDelete)
	}
}
