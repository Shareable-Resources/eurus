package bc_indexer

import (
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service_server"
	"eurus-backend/user_service/user_service/user"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type BlockChainIndexer struct {
	service_server.ServiceServer
	Config *BlockChainIndexerConfigBase
	Asset
	MainNetAsset                 Asset
	mainnetLoggerName            string
	sideChainLoggerName          string
	mainnetRescanLoggerName      string
	sideChainRescanLoggerName    string
	sideChainScanBlockCounter    *ethereum.ScanBlockCounter
	mainnetScanBlockCounter      *ethereum.ScanBlockCounter
	EurusUserDepositAddress      *common.Address
	MainnetPlatformWalletAddress *common.Address
	SweepServiceInvokerAddress   *common.Address
	MainnetPlatformWalletUser    *user.User
	SweepServiceInvokerUser      *user.User
	transferRewardProcessor      *TransferRewardProcessor
}

type Asset struct {
	AssetList map[common.Address]string
}

func NewBlockChainIndexer() *BlockChainIndexer {
	indexer := new(BlockChainIndexer)
	indexer.Config = NewBlockChainIndexerConfigBase()
	indexer.ServerConfig = &indexer.Config.ServerConfigBase
	indexer.mainnetLoggerName = log.Name.Mainnet
	indexer.sideChainLoggerName = log.Name.Root
	indexer.mainnetRescanLoggerName = log.Name.MainnetRescan

	return indexer
}

// func (me *BlockChainIndexer) LoadConfig(args *service.ServiceCommandLineArgs) error {
// 	err := me.ServiceServer.LoadConfig(args, me.Config)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (me *BlockChainIndexer) InitConfig() {
	me.ServiceServer.ServerConfig = &me.Config.ServerConfigBase
}

func (me *BlockChainIndexer) InitAll() {

	mainnetFilePath := asset.ConstructLogFilePath(me.mainnetLoggerName, me.Config.LogFilePath)

	_, err := log.NewLogger(me.mainnetLoggerName, mainnetFilePath, logrus.DebugLevel)
	if err != nil {
		panic("Unable to create mainnet log file: " + err.Error())
	}

	me.ServiceServer.InitAuth(me.authLoginHandler, nil)
}

func (me *BlockChainIndexer) authLoginHandler(authClient auth_base.IAuth) {
	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}

	err = me.initHdWallet()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Init HDWallet error: ", err)
		panic("Init HDWallet error: " + err.Error())
	}

	fmt.Println("Blockchain indexer wallet address: ", me.Config.HdWalletAddress)

	err = me.InitDBFromConfig(&me.Config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic(err)
	}

	//Test DB connection
	_, err = me.DefaultDatabase.GetConn()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connection DB: ", err.Error())
	} else {
		log.GetLogger(log.Name.Root).Infoln("Login DB successfully")
	}
	_, err = me.InitEthereumWebSocketClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthWebSocketClient: ", err.Error())
		panic(err)
	}

	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthClient: ", err.Error())
		panic(err)
	}

	err = me.LoadAssetList()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Load Asset List: ", err.Error())
		panic(err)
	}
	err = me.InitBlockSubscriber()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init BlockSubscriber: ", err.Error())
		panic(err)
	}

	mainnetServerConfig := new(server.ServerConfigBase)
	mainnetServerConfig.EthClientChainID = me.Config.MainnetEthClientChainID
	mainnetServerConfig.EthClientIP = me.Config.MainnetEthClientIP
	mainnetServerConfig.EthClientPort = me.Config.MainnetEthClientPort
	mainnetServerConfig.EthClientProtocol = me.Config.MainnetEthClientProtocol
	mainnetServerConfig.EthClientWebSocketPort = me.Config.MainnetEthClientWebSocketPort
	mainnetServerConfig.EthClientWebSocketProtocol = me.Config.MainnetEthClientWebSocketProtocol
	mainnetServerConfig.EthClientWebSocketIP = me.Config.MainnetEthClientWebSocketIP

	_, err = me.InitMainNetEthereumClientFromConfig(mainnetServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthClient: ", err.Error())
		panic(err)
	}
	_, err = me.InitMainNetEthereumWebSocketClientFromConfig(mainnetServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthWebSocketClient: ", err.Error())
		panic(err)
	}

	err = me.InitMainNetBlockSubscriber()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init BlockSubscriber: ", err.Error())
		panic(err)
	}
	err = me.LoadMainNetAssetList()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Load Asset List: ", err.Error())
		panic(err)
	}

	addr, err := GetEurusUserDepositAddress(me.Config.EurusInternalConfigAddress, me.MainNetEthClient)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetEurusUserDepositAddress failed: ", err.Error())
		panic(err)
	}
	me.EurusUserDepositAddress = addr

	addr, err = GetMainnetPlatformWalletAddress(me.Config.EurusInternalConfigAddress, me.MainNetEthClient)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetMainnetPlatformWalletAddress failed: ", err.Error())
		panic(err)
	}
	me.MainnetPlatformWalletAddress = addr

	me.MainnetPlatformWalletUser = new(user.User)
	me.MainnetPlatformWalletUser.Id = 0
	me.MainnetPlatformWalletUser.IsMetamaskAddr = true
	me.MainnetPlatformWalletUser.MainnetWalletAddress = ethereum.ToLowerAddressString(me.MainnetPlatformWalletAddress.Hex())

	invokerAddrStr, err := me.QuerySystemConfig("sweepServiceInvokerAddress")
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err.Error())
		panic(err)
	}
	if invokerAddrStr == "" {
		log.GetLogger(log.Name.Root).Errorln("Sweep service invoker address is empty")
		panic("Sweep service invoker address is empty")
	}

	invokerAddr := common.HexToAddress(invokerAddrStr)
	me.SweepServiceInvokerAddress = &invokerAddr

	me.SweepServiceInvokerUser = new(user.User)
	me.SweepServiceInvokerUser.Id = 0
	me.SweepServiceInvokerUser.IsMetamaskAddr = true
	me.SweepServiceInvokerUser.MainnetWalletAddress = ethereum.ToLowerAddressString(invokerAddr.String())

	me.transferRewardProcessor = NewTransferRewardProcessor(common.HexToAddress(me.Config.InternalSCConfigAddress),
		me.EthClient, me.DefaultDatabase, me.SlaveDatabase,
		me.Config.HdWalletPrivateKey, log.GetLogger(log.Name.Root))
	err = me.transferRewardProcessor.Init(me.Config.RegistrationCriteriaListJsonStr)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init transferRewardProcessor: ", err.Error())
		panic("Unable to init transferRewardProcessor. Error: " + err.Error())
	}

	me.mainnetScanBlockCounter, err = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeContinuous,
		uint64(me.Config.MainnetEthClientChainID), me.rescanMainnetBlock, "BlockChainIndexer_Mainnet_Server_"+strconv.FormatInt(me.Config.ServiceId, 10)+".db", me.mainnetLoggerName)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init mainnet scan block counter: ", err.Error())
		panic("Unable to init mainnet scan block counter: " + err.Error())
	}
	log.GetLogger(log.Name.Root).Infoln("Mainnet local DB: BlockChainIndexer_Mainnet_Server_" + strconv.FormatInt(me.Config.ServiceId, 10) + ".db")

	me.sideChainScanBlockCounter, err = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeContinuous,
		uint64(me.Config.EthClientChainID), me.rescanSideChainBlock, "BlockChainIndexer_SideChain_Server_"+strconv.FormatInt(me.Config.ServiceId, 10)+".db", me.sideChainLoggerName)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init side chain scan block counter: ", err.Error())
		panic("Unable to init side chain scan block counter: " + err.Error())
	}

	log.GetLogger(log.Name.Root).Infoln("Side chain local DB: BlockChainIndexer_SideChain_Server_" + strconv.FormatInt(me.Config.ServiceId, 10) + ".db")

	err = me.mainnetScanBlockCounter.Start()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Start mainnet scan block counter error: ", err.Error())
		panic("Start mainnet scan block counter error: " + err.Error())
	}

	err = me.sideChainScanBlockCounter.Start()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Start side chain scan block counter error: ", err.Error())
		panic("Start side chain scan block counter error: " + err.Error())
	}

	fmt.Println("Blockchain indexer is ready to start")
	log.GetLogger(log.Name.Root).Infoln("Blockchain indexer is ready to start. Service ID: ", me.Config.ServiceId)
	go me.RunBlockSubscriber()
	go me.RunMainNetBlockSubscriber()
}

func (me *BlockChainIndexer) QueryAuthServerInfo() error {
	err := me.ServiceServer.QueryAuthServerInfo()
	return err
}

func (me *BlockChainIndexer) initHdWallet() error {
	if me.Config.MnemonicPhase == "" {
		return errors.New("Missing Mnemonic phase")
	}

	wallet, acc, addr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.GetServiceId(), 10))
	if err != nil {
		return err
	}
	me.Config.HdWalletAddress = addr
	me.Config.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*acc)
	return nil
}
func (me *BlockChainIndexer) InitBlockSubscriber() error {
	err := me.EthWebSocketClient.SubscribeBlockHeader(0)
	if err != nil {
		return err
	}
	return nil

}
func (me *BlockChainIndexer) InitMainNetBlockSubscriber() error {
	err := me.MainNetWebSocketClient.SubscribeBlockHeader(8)
	if err != nil {
		return err
	}
	return nil

}

func (me *BlockChainIndexer) RunBlockSubscriber() {
	context := new(blockChainProcessorContext)
	context.Db = me.DefaultDatabase
	context.SlaveDb = me.SlaveDatabase
	context.EthClient = me.EthClient
	context.IsMainnet = false
	context.AssetAddressMap = me.Asset
	context.Config = me.Config
	context.LoggerName = me.sideChainLoggerName
	context.EurusUserDepositAddress = nil
	context.transferRewardProcessor = me.transferRewardProcessor

	for {
		block, serverErr := me.EthWebSocketClient.Subscriber.GetLatestBlock(true)
		if serverErr != nil {
			if serverErr.GetReturnCode() != foundation.NetworkError {
				log.GetLogger(log.Name.Root).Error("Unable to get latest block: ", serverErr.Error())
			}
		} else {

			ProcessBlock(context, block)

			err := me.sideChainScanBlockCounter.UpdateLatestBlock(block.Number())
			if err != nil {
				log.GetLogger(context.LoggerName).Errorln("UpdateLatestBlock failed on block: ", block.Number().String(), " error: ", err.Error())
			}
		}
	}
}

func (me *BlockChainIndexer) RunMainNetBlockSubscriber() {
	context := new(blockChainProcessorContext)
	context.Db = me.DefaultDatabase
	context.SlaveDb = me.SlaveDatabase
	context.EthClient = me.MainNetEthClient
	context.IsMainnet = true
	context.AssetAddressMap = me.MainNetAsset
	context.Config = me.Config
	context.LoggerName = me.mainnetLoggerName
	context.EurusUserDepositAddress = me.EurusUserDepositAddress
	context.MainnetPlatformWalletAddress = me.MainnetPlatformWalletAddress
	context.MainnetPlatformWalletUser = me.MainnetPlatformWalletUser
	context.SweepServiceInvokerAddress = me.SweepServiceInvokerAddress
	context.SweepServiceInvokerUser = me.SweepServiceInvokerUser

	for {
		block, serverErr := me.MainNetWebSocketClient.Subscriber.GetLatestBlock(false)
		if serverErr != nil {
			if serverErr.ReturnCode != foundation.NetworkError {
				log.GetLogger(log.Name.Root).Error("Unable to get latest block: ", serverErr.Error())
			}
		} else {
			ProcessMainNetBlock(context, block)
			err := me.mainnetScanBlockCounter.UpdateLatestBlock(block.Number())
			if err != nil {
				log.GetLogger(context.LoggerName).Errorln("UpdateLatestBlock failed on block: ", block.Number().String(), " error: ", err.Error())
			}
		}
	}
}

func (me *BlockChainIndexer) LoadAssetList() error {
	me.AssetList = make(map[common.Address]string)
	assetName, assetAddr, err := GetAssetList(me.ServerConfig.ExternalSCConfigAddress, me.EthClient)
	var i int
	for i = 0; i < me.Config.GetRetryCount(); i++ {
		if err != nil {
			time.Sleep(time.Duration(me.Config.RetryInterval) * time.Second)
			continue
		}
	}
	if err != nil {
		return err
	}

	for i = 0; i < len(assetName); i++ {
		me.AssetList[assetAddr[i]] = assetName[i]
	}
	return nil
}

func (me *BlockChainIndexer) LoadMainNetAssetList() error {
	me.MainNetAsset.AssetList = make(map[common.Address]string)
	assetName, assetAddr, err := GetMainNetAssetList(me.ServerConfig.EurusInternalConfigAddress, me.MainNetEthClient)
	var i int
	for i = 0; i < me.Config.GetRetryCount(); i++ {
		if err != nil {
			time.Sleep(time.Duration(me.Config.RetryInterval) * time.Second)
			continue
		}
	}
	if err != nil {
		return err
	}

	for i = 0; i < len(assetName); i++ {
		me.MainNetAsset.AssetList[assetAddr[i]] = assetName[i]
	}
	return nil
}

func (me *BlockChainIndexer) rescanMainnetBlock(sender *ethereum.ScanBlockCounter, eventId decimal.Decimal,
	chainId uint64, from decimal.Decimal, to decimal.Decimal) {
	loggerName := asset.ConstructRescanLoggerName(log.Name.MainnetRescan, me.Config.ServiceId, eventId)
	//loggerName := me.mainnetRescanLoggerName + "_" + strconv.FormatInt(me.Config.ServiceId, 10) + "_" + eventId.String()
	loggerFilePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, loggerFilePath, logrus.DebugLevel)
	if err != nil {
		log.GetLogger(me.mainnetLoggerName).Errorln("Unable to create new logger for rescan event : ", eventId.String(), ". Use default one instead")
		loggerName = me.mainnetLoggerName
	}

	defer log.RemoveLogger(loggerName)
	fmt.Println("Start mainnet rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())

	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())
	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64()
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)

	context := new(blockChainProcessorContext)
	context.Db = me.DefaultDatabase
	context.SlaveDb = me.SlaveDatabase
	context.EthClient = me.MainNetEthClient
	context.IsMainnet = true
	context.AssetAddressMap = me.MainNetAsset
	context.Config = me.Config
	context.LoggerName = loggerName
	context.EurusUserDepositAddress = me.EurusUserDepositAddress
	context.MainnetPlatformWalletAddress = me.MainnetPlatformWalletAddress
	context.MainnetPlatformWalletUser = me.MainnetPlatformWalletUser
	context.SweepServiceInvokerAddress = me.SweepServiceInvokerAddress
	context.SweepServiceInvokerUser = me.SweepServiceInvokerUser

	for i := int64(0); i <= total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting block number: ", currBlockNum.String())
		for {
			block, err = me.MainNetEthClient.GetBlockByNumber(currBlockNum)
			if err != nil {
				log.GetLogger(loggerName).Errorln("Unable to get block by number. Block number: ", currBlockNum, " Error: ", err)
				time.Sleep(me.Config.GetRetryInterval() * time.Second)
				continue
			}
			break
		}

		log.GetLogger(loggerName).Debugln("Start processing block number: ", block.Number().String())

		ProcessMainNetBlock(context, block)

		status, err := me.mainnetScanBlockCounter.UpdateRescanEvent(eventId, currBlockNum)
		if err != nil {
			log.GetLogger(loggerName).Errorln("UpdateRescanEvent failed for ", currBlockNum.String(), ". Error: ", err.Error())

		} else {
			if status == ethereum.ScanBlockActionFinished {
				break
			}
		}
	}
	fmt.Println("Mainnet Rescan event: ", eventId.String(), " finished")
	log.GetLogger(loggerName).Infoln("Rescan event: ", eventId.String(), " finished")

}

func (me *BlockChainIndexer) rescanSideChainBlock(sender *ethereum.ScanBlockCounter, eventId decimal.Decimal,
	chainId uint64, from decimal.Decimal, to decimal.Decimal) {

	fmt.Println("Start sidechain rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())

	loggerName := asset.ConstructRescanLoggerName(log.Name.SideChainRescan, me.Config.ServiceId, eventId)
	loggerFilePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, loggerFilePath, logrus.DebugLevel)
	if err != nil {
		fmt.Println("Unable to create new logger for rescan event : ", eventId.String(), ". Use default one instead")
		log.GetLogger(log.Name.Root).Errorln("Unable to create new logger for rescan event : ", eventId.String(), ". Use default one instead")
		loggerName = log.Name.Root
	}

	defer log.RemoveLogger(loggerName)

	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())
	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64()
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)

	context := new(blockChainProcessorContext)
	context.Db = me.DefaultDatabase
	context.SlaveDb = me.SlaveDatabase
	context.EthClient = me.EthClient
	context.IsMainnet = false
	context.AssetAddressMap = me.Asset
	context.Config = me.Config
	context.LoggerName = loggerName
	context.EurusUserDepositAddress = nil
	context.transferRewardProcessor = NewTransferRewardProcessor(common.HexToAddress(me.Config.InternalSCConfigAddress),
		me.EthClient, me.DefaultDatabase, me.SlaveDatabase,
		me.Config.HdWalletPrivateKey, log.GetLogger(loggerName))

	for {
		err = context.transferRewardProcessor.Init(me.Config.RegistrationCriteriaListJsonStr)
		if err != nil {
			log.GetLogger(loggerName).Error("Unable to init transferRewardProcessor: ", err.Error(), " Wait for ", me.Config.GetRetryInterval()+60*time.Second, " to retry")
			time.Sleep(me.Config.GetRetryInterval() + 60*time.Second)
			continue
		}
		break
	}

	for i := int64(0); i <= total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting block number: ", currBlockNum.String())
		for {
			block, err = me.EthClient.GetBlockByNumber(currBlockNum)
			if err != nil {
				log.GetLogger(loggerName).Errorln("Unable to get block by number. Block number: ", currBlockNum, " Error: ", err)
				time.Sleep(me.Config.GetRetryInterval() * time.Second)
				continue
			}
			break
		}

		log.GetLogger(loggerName).Debugln("Start processing block number: ", block.Number().String())

		ProcessBlock(context, block)

		log.GetLogger(loggerName).Debugln("Finished processing block number: ", block.Number().String())

		status, err := me.sideChainScanBlockCounter.UpdateRescanEvent(eventId, currBlockNum)
		if err != nil {
			log.GetLogger(loggerName).Errorln("UpdateRescanEvent failed for ", currBlockNum.String(), ". Error: ", err.Error())
		} else {
			if status == ethereum.ScanBlockActionFinished {
				break
			}
		}
	}
	fmt.Println("Sidechain Rescan event: ", eventId.String(), " finished")
	log.GetLogger(loggerName).Infoln("Rescan event: ", eventId.String(), " finished")

}
