package deposit

import (
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type AssetAddressMap map[common.Address]string

type AssetNameMap map[string]common.Address

type DepositSCProcessor struct {
	sideChainEthClient *ethereum.EthClient
	mainnetEthClient   *ethereum.EthClient

	MainnetAssetInfo            AssetAddressMap
	MainnetAssetNameToAddress   AssetNameMap
	MainnetUserDepositAddr      *common.Address
	SidechainPlatformWalletAddr *common.Address
	config                      *DepositObserverConfig
	context                     *DepositProcessorContext //context has replaced db

	SubmitMintRequestMethodId []byte
}

func NewDepositSCProcessor(config *DepositObserverConfig, context *DepositProcessorContext) *DepositSCProcessor {
	processor := new(DepositSCProcessor)
	processor.config = config
	processor.context = context
	return processor
}

func (me *DepositSCProcessor) Init() error {
	var err error
	//Connect sidechain
	me.sideChainEthClient = &ethereum.EthClient{
		Protocol: me.config.EthClientProtocol,
		IP:       me.config.EthClientIP,
		Port:     me.config.EthClientPort,
		ChainID:  big.NewInt(int64(me.config.EthClientChainID)),
	}
	_, err = me.sideChainEthClient.Connect()
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Connect sidechain failed: ", err)
		return errors.WithMessage(err, "Connect sidechain failed")
	}

	//Connect mainchain
	me.mainnetEthClient = &ethereum.EthClient{
		Protocol: me.config.MainnetEthClientProtocol,
		IP:       me.config.MainnetEthClientIP,
		Port:     me.config.MainnetEthClientPort,
		ChainID:  big.NewInt(int64(me.config.MainnetEthClientChainID)),
	}
	_, err = me.mainnetEthClient.Connect()
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Connect mainnet failed: ", err)
		return errors.WithMessage(err, "Connect mainnet failed")
	}

	me.MainnetAssetInfo, err = me.queryMainnetAssetContractInfo()
	if err != nil {
		return err
	}

	// Make the reverse lookup map
	me.MainnetAssetNameToAddress = make(AssetNameMap)
	for addr, name := range me.MainnetAssetInfo {
		me.MainnetAssetNameToAddress[name] = addr
	}

	me.MainnetUserDepositAddr, err = me.queryDecentralizedUserDepositAddress()
	if err != nil {
		return err
	}
	me.SidechainPlatformWalletAddr, err = me.querySideChainPlatformWalletAddress()
	if err != nil {
		return err
	}

	platformWalletAbi := ethereum.DefaultABIDecoder.GetABI("PlatformWallet")
	submitMintRequestMethod, ok := platformWalletAbi.Methods["submitMintRequest"]

	if !ok {
		return errors.New("PlatformWallet submitMintRequest ABI not found")
	}
	me.SubmitMintRequestMethodId = submitMintRequestMethod.ID

	log.GetLogger(me.context.LoggerName).Infoln("Mainnet user deposit address: ", me.MainnetUserDepositAddr.Hex())

	log.GetLogger(me.context.LoggerName).Infoln("Side chain Platform wallet address: ", me.SidechainPlatformWalletAddr.Hex())
	log.GetLogger(me.context.LoggerName).Infoln("Mainnet asset info")
	for assetAddr, assetName := range me.MainnetAssetInfo {
		log.GetLogger(me.context.LoggerName).Infof("%s: %s", assetName, assetAddr.Hex())
	}
	return nil
}

func (me *DepositSCProcessor) queryMainnetAssetContractInfo() (AssetAddressMap, error) {
	externalSCConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(me.config.EurusInternalConfigAddress), me.mainnetEthClient.Client)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("NewEurusInternalConfig failed: ", err)
		return nil, err
	}
	assetNameList, addrList, err := externalSCConfig.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("GetAssetAddress failed: ", err, " contract address: ", me.config.EurusInternalConfigAddress)
		return nil, err
	}

	var assetInfo AssetAddressMap = make(AssetAddressMap)

	for i, assetName := range assetNameList {
		assetInfo[addrList[i]] = assetName
	}
	return assetInfo, nil
}

func (me *DepositSCProcessor) queryDecentralizedUserDepositAddress() (*common.Address, error) {
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(me.config.EurusInternalConfigAddress), me.mainnetEthClient.Client)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Unable to create ExternalSmartContractConfig object: ", err.Error())
		err1 := errors.WithMessage(err, "Unable to create ExternalSmartContractConfig object")
		return nil, err1
	}

	depositAddr, err := eurusInternalConfig.EurusUserDepositAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Unable to query EurusUserDepositAddress: ", err.Error())
		err1 := errors.WithMessage(err, "Unable to query EurusUserDepositAddress")
		return nil, err1
	}

	return &depositAddr, err
}

func (me *DepositSCProcessor) querySideChainPlatformWalletAddress() (*common.Address, error) {
	internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.config.InternalSCConfigAddress), me.sideChainEthClient.Client)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Unable to create InternalSmartContractConfig instance: ", err.Error())
		return nil, errors.WithMessage(err, "Unable to create InternalSmartContractConfig instance")
	}

	addr, err := internalSCConfig.GetInnetPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorln("Unable to GetInnetPlatformWalletAddress: ", err.Error())
		return nil, errors.WithMessage(err, "Unable to GetInnetPlatformWalletAddress")
	}

	return &addr, nil
}

func (me *DepositSCProcessor) SweepTokenToPlatformWallet(transHash common.Hash, assetName string, sender common.Address, amount *big.Int) (*ethereum.BesuReceipt, *big.Int, error) {

	userDepositSC, err := mainnet_contract.NewEurusUserDeposit(*me.MainnetUserDepositAddr, me.mainnetEthClient.Client)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorf("Cannot NewEurusUserDeposit. Trans hash: %s, Error: %s\r\n ", transHash.Hex(), err.Error())
		return nil, nil, errors.WithMessagef(err, "Cannot NewEurusUserDeposit. Trans hash: %s", transHash.Hex())
	}

	tx, err := me.mainnetEthClient.InvokeSmartContract(me.config, me.config.HdWalletPrivateKey,
		me.config.SideChainGasLimit,
		func(ethClient *ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
			tx, err := userDepositSC.Sweep(transOpt, transHash, sender, assetName, amount)
			if err != nil {
				log.GetLogger(me.context.LoggerName).Debugf("Sweep failed. Trans Hash: %s. Error: %s", transHash.Hex(), err.Error())
			}
			return tx, false, err
		})

	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorf("Sweep error. Deposit Trans hash: %s, Error: %s\r\n", transHash.Hex(), err.Error())
		return nil, nil, err
	}

	log.GetLogger(me.context.LoggerName).Infof("Sweep tx hash: %s. Deposit trans hash: %s\r\n", tx.Hash().Hex(), transHash.Hex())

	receipt, err := me.mainnetEthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(me.context.LoggerName).Errorf("Sweep query receipt error. Deposit Trans hash: %s, Sweep Trans Hash %s, Error: %s\r\n", transHash.Hex(), tx.Hash().Hex(), err.Error())
		return nil, nil, err
	}

	if receipt.Status != 1 {
		receiptByte, _ := json.Marshal(receipt)
		log.GetLogger(me.context.LoggerName).Errorf("Receipt status is 0. Trans hash: %s. Deposit Trans hash: %s, Receipt: %s\r\n", tx.Hash().Hex(), transHash.Hex(), string(receiptByte))
		return nil, nil, errors.New("Receipt status is 0")
	}

	return receipt, tx.GasPrice(), nil
}

func (me *DepositSCProcessor) SubmitMintRequest(depositTransHash common.Hash, assetName string,
	sidechainDestAddr common.Address, amount *big.Int) (*ethereum.BesuReceipt, error) {

	platformWallet, err := contract.NewPlatformWallet(*me.SidechainPlatformWalletAddr, me.sideChainEthClient.Client)
	if err != nil {
		return nil, errors.WithMessage(err, "NewPlatformWallet error")
	}

	tx, err := me.sideChainEthClient.InvokeSmartContract(me.config, me.config.HdWalletPrivateKey,
		me.config.SideChainGasLimit,
		func(ethClient *ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
			tx, err1 := platformWallet.SubmitMintRequest(transOpt, sidechainDestAddr, assetName, amount, depositTransHash)
			if err1 != nil {
				log.GetLogger(me.context.LoggerName).Debugf("SubmitMintRequest failed. Deposit trans Hash: %s. Error: %s", depositTransHash.Hex(), err1.Error())
			}
			return tx, false, err1
		})

	if err != nil {
		err = errors.WithMessagef(err, "SubmitMintRequest error. Deposit trans hash: %s", depositTransHash.Hex())
		return nil, err
	}

	log.GetLogger(me.context.LoggerName).Infoln("SubmitMintRequest broadcast success. tx hash: ", tx.Hash().Hex(), " deposit trans hash: ", depositTransHash.Hex())

	receipt, err := me.sideChainEthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		return nil, errors.WithMessagef(err, "SubmitMinRequest query receipt failed. tx hash: %s,  deposit trans hash: %s ", tx.Hash().Hex(), depositTransHash.Hex())
	}

	return receipt, nil
}

// func (me *DepositSCProcessor) ConfirmMint(mintRequestId *uint256.Int) (*ethereum.BesuReceipt, error) {
// 	platformWallet, err1 := contract.NewPlatformWallet(*me.SidechainPlatformWalletAddr, me.sideChainEthClient.Client)
// 	if err1 != nil {
// 		return nil, errors.WithMessage(err1, "NewPlatformWallet error")
// 	}

// 	tx, err := me.sideChainEthClient.InvokeSmartContract(me.config, me.config.HdWalletPrivateKey,
// 		func(ethClient *ethereum.EthClient, transOpt *bind.TransactOpts) (*types.Transaction, bool, error) {
// 			log.GetLogger(me.context.LoggerName).Debugln("Before ConfirmTransaction: ", mintRequestId.String())
// 			if me.config.SideChainGasLimit > 0 {
// 				transOpt.GasLimit = me.config.SideChainGasLimit
// 			}
// 			tx1, err1 := platformWallet.ConfirmTransaction(transOpt, mintRequestId.ToBig())
// 			if err1 != nil {
// 				if strings.Contains(err1.Error(), "Transaction is already confirmed") {
// 					return nil, true, err1
// 				} else {
// 					log.GetLogger(me.context.LoggerName).Debugf("ConfirmTransaction failed. Mint request ID: %s, error: %s", mintRequestId.String(), err1.Error())
// 					err1 = errors.WithMessage(err1, " failed to estimate gas needed ")
// 				}
// 			}
// 			return tx1, false, err1
// 		})

// 	if err != nil {
// 		if strings.Contains(err.Error(), "Transaction is already confirmed") {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	log.GetLogger(me.context.LoggerName).Debugln("Before me.sideChainEthClient.QueryEthReceiptWithSetting")
// 	receipt, err1 := me.sideChainEthClient.QueryEthReceiptWithSetting(tx, 1, 20)

// 	if receipt.RevertReason != "" {
// 		checkRevertReason, err := hex.DecodeString(receipt.RevertReason[2:])
// 		if err != nil {
// 			log.GetLogger(me.context.LoggerName).Errorln("Fail to decode error", err1)
// 		}
// 		if strings.Contains(unix.ByteSliceToString(checkRevertReason), "Transaction is already confirmed") {
// 			return nil, nil
// 		}
// 	}

// 	if err1 != nil {
// 		return nil, errors.WithMessagef(err1, "ConfirmTransaction query receipt failed. Mint Request Id: %s", mintRequestId.String())
// 	}

// 	return receipt, nil
// }

func (me *DepositSCProcessor) QueryMainnetEthReceiptWithSetting(tx *types.Transaction, waitSecond int, retryCount int) (*ethereum.BesuReceipt, error) {
	return me.mainnetEthClient.QueryEthReceiptWithSetting(tx, waitSecond, retryCount)
}

func (me *DepositSCProcessor) QuerySideChainEthReceiptWithSetting(tx *types.Transaction, waitSecond int, retryCount int) (*ethereum.BesuReceipt, error) {
	return me.sideChainEthClient.QueryEthReceiptWithSetting(tx, waitSecond, retryCount)
}

func (me *DepositSCProcessor) GetSideChainBlockTimeFromBlockNumber(blockNum *big.Int) (*time.Time, error) {
	var err error
	var block *types.Block
	for i := 0; i < me.config.GetRetryCount() || i == 0; i++ {
		block, err = me.sideChainEthClient.GetBlockByNumber(blockNum)
		if err != nil {
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	blockTime := time.Unix(int64(block.Time()), 0)
	return &blockTime, nil
}

func (me *DepositSCProcessor) GetMainnetBlockTimeFromBlockNumber(blockNum *big.Int) (*time.Time, error) {
	var err error
	var block *types.Block
	for i := 0; i < me.config.GetRetryCount() || i == 0; i++ {
		block, err = me.mainnetEthClient.GetBlockByNumber(blockNum)
		if err != nil {
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	blockTime := time.Unix(int64(block.Time()), 0)
	return &blockTime, nil
}

func (me *DepositSCProcessor) GetMainnetBalance(address common.Address, assetName string) (*big.Int, error) {
	if assetName == "ETH" {
		return me.mainnetEthClient.GetBalance(address)
	}

	contractAddress, found := me.MainnetAssetNameToAddress[assetName]
	if !found {
		return nil, errors.Errorf("Cannot find the address for asset %v", assetName)
	}

	return me.getTokenBalance(address, contractAddress, me.mainnetEthClient.Client)
}

func (me *DepositSCProcessor) GetSideChainBalance(address common.Address, assetName string) (*big.Int, error) {
	if assetName == asset.EurusTokenName {
		return me.sideChainEthClient.GetBalance(address)
	}

	internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.config.InternalSCConfigAddress), me.sideChainEthClient.Client)
	if err != nil {
		return nil, err
	}

	contractAddress, err := internalSCConfig.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, assetName)
	if err != nil {
		return nil, err
	}

	return me.getTokenBalance(address, contractAddress, me.sideChainEthClient.Client)
}

func (me *DepositSCProcessor) getTokenBalance(address common.Address, contractAddress common.Address, client bind.ContractBackend) (*big.Int, error) {
	inst, err := contract.NewEurusERC20(contractAddress, client)
	if err != nil {
		return nil, err
	}

	return inst.BalanceOf(&bind.CallOpts{}, address)
}
