package userObserver

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"math/big"
	"strconv"
	"time"

	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type UserObserver struct {
	service_server.ServiceServer
	Config                    *UserObserverConfig
	scProcessor               *UserObserverSCProcessor
	processor                 *UserObserverProcessor
	sideChainScanBlockCounter *eurus_ethereum.ScanBlockCounter
	//counter
	//sideChainLogSubscription go_ethereum.Subscription
	//sideChainLogChannel      chan (types.Log)
}

func NewUserObserver() *UserObserver {
	observer := new(UserObserver)
	observer.Config = NewUserObserverConfig()
	observer.ServerConfig = &observer.Config.ServerConfigBase

	return observer
}

func (me *UserObserver) LoadConfig(args *service.ServiceCommandLineArgs) error {
	err := me.ServiceServer.LoadConfig(args, me.Config)
	if err != nil {
		return err
	}
	return nil
}

func (me *UserObserver) InitConfig() {
	me.ServiceServer.ServerConfig = &me.Config.ServerConfigBase
}

func (me *UserObserver) InitAll() {

	me.ServiceServer.InitAuth(me.processInit, nil)
}

func (me *UserObserver) InitAndStartSideChainScanBlockCounter() {
	var err error
	var prefix = "UserObserver"
	var loggerName string = log.Name.Root //Normal routine would use root logger name
	var dbFileName string = prefix + "_SideChain_Server_" + strconv.FormatInt(me.Config.ServiceId, 10) + ".db"
	me.sideChainScanBlockCounter, err = eurus_ethereum.NewScanBlockCounter(eurus_ethereum.ScanBlockModeContinuous,
		uint64(me.Config.EthClientChainID), me.rescanSideChainBlock, dbFileName, loggerName)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to init side chain scan block counter: ", err.Error())
		panic("Unable to init side chain scan block counter: " + err.Error())
	}
	err = me.sideChainScanBlockCounter.Start()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Start side chain scan block counter error: ", err.Error())
		panic("Start side chain scan block counter error: " + err.Error())
	}
}

func (me *UserObserver) processInit(authClient auth_base.IAuth) {
	// 1. creates normal logger name
	var normalRoutineProcessorLoggerName string = log.Name.Root
	// 2. starts getting config from config server (This steps require config server running)
	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}
	// 3. Init Db From Config
	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic(err)
	} else {
		log.GetLogger(log.Name.Root).Infoln("Init DB successfully")
	}
	// 4. Generate wallet, account, hdWallerAddr, err
	wallet, account, hdWalletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.GetServiceId(), 10))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load wallet key failed: " + err.Error())
		panic("Load wallet key failed: " + err.Error())
	}
	me.Config.HdWalletAddress = hdWalletAddr
	me.Config.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*account)
	fmt.Println("User observer wallet address: " + hdWalletAddr)
	log.GetLogger(log.Name.Root).Infoln("User observer wallet address: " + hdWalletAddr)
	// 5. Get ethereum client, web socket
	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthClient: ", err.Error())
		panic(err)
	}
	var sideChainClient *eurus_ethereum.EthClient
	sideChainClient, err = me.InitEthereumWebSocketClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthWebSocketClient for side chain: ", err.Error())
		panic(err)
	}

	// 6. Init Smart Contract Observer Processor
	me.scProcessor = NewUserObserverSCProcessor(me.Config, normalRoutineProcessorLoggerName)
	err = me.scProcessor.Init()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init UserObserverSCProcessor: ", err.Error())
		panic(err)
	}
	// 7. Init User Observer Processor
	me.processor = NewUserObserverProcessor(me.DefaultDatabase, me.SlaveDatabase, me.Config, me.scProcessor, normalRoutineProcessorLoggerName)
	err = sideChainClient.SubscribeBlockHeader(0)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to subscribe side chain block head: ", err.Error())
		panic(err.Error())
	}
	//8. Init side chain scan block counter for rescanning and store in sqlite
	me.InitAndStartSideChainScanBlockCounter()
	//9. Use side chain scan block counter for rescanning and store in sqlite
	go me.RunBlockSubscriber()
}

func (me *UserObserver) RunBlockSubscriber() {
	fmt.Println("Running Block Scanner")

	log.GetLogger(log.Name.Root).Debug("Running Block Scanner")
	for {
		block, serverErr := me.EthWebSocketClient.Subscriber.GetLatestBlock(true)

		if serverErr != nil {
			if serverErr.GetReturnCode() != foundation.NetworkError {
				log.GetLogger(log.Name.Root).Error("Unable to get block by block numner: ", serverErr.Error())
			}
			continue
		}
		me.processor.BlockHandler(block)
		//Update sqllite block number by sideChainScanBlockCounter
		err := me.sideChainScanBlockCounter.UpdateLatestBlock(block.Number())
		if err != nil {
			log.GetLogger(me.sideChainScanBlockCounter.LoggerName).Errorln("UpdateLatestBlock failed on block: ", block.Number().String(), " error: ", err.Error())
		}
	}

}

// Call ServerConfigBase.QueryAuthServerInfo()
func (me *UserObserver) QueryAuthServerInfo() error {
	err := me.ServiceServer.QueryAuthServerInfo()
	return err
}

func (me *UserObserver) InitBlockSubscriber() error {
	err := me.EthWebSocketClient.SubscribeBlockHeader(0)
	if err != nil {
		return err
	}
	return nil

}

func (me *UserObserver) rescanSideChainBlock(sender *eurus_ethereum.ScanBlockCounter, eventId decimal.Decimal,
	chainId uint64, from decimal.Decimal, to decimal.Decimal) {
	loggerName := asset.ConstructRescanLoggerName(log.Name.SideChainRescan, me.Config.ServiceId, eventId)
	loggerFilePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, loggerFilePath, logrus.DebugLevel)
	if err != nil {
		fmt.Println("Unable to create new logger for rescan event : ", eventId.String(), ". Use default one instead")
		log.GetLogger(log.Name.Root).Errorln("Unable to create new logger for rescan event : ", eventId.String(), ". Use default one instead")
		loggerName = log.Name.Root
	}
	defer log.RemoveLogger(loggerName)
	fmt.Println("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())
	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())
	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64()
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)

	for i := int64(0); i <= total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting block number: ", currBlockNum.String())
		//fmt.Println("Start getting block number: ", currBlockNum.String())
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

		me.processor.BlockHandler(block)

		status, err := me.sideChainScanBlockCounter.UpdateRescanEvent(eventId, currBlockNum)
		if err != nil {
			log.GetLogger(loggerName).Errorln("UpdateRescanEvent failed for ", currBlockNum.String(), ". Error: ", err.Error())

		} else {
			if status == eurus_ethereum.ScanBlockActionFinished {
				break
			}
		}
	}
	log.GetLogger(loggerName).Infoln("Rescan event: ", eventId.String(), " finished")
}
