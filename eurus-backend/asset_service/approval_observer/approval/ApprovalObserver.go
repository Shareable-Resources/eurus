package approval

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service_server"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type ApprovalObserver struct {
	service_server.ServiceServer
	Config          *ApprovalObserverConfig
	logChannel      chan (types.Log)
	logSubscription ethereum.Subscription
	closeChannel    chan (int)
	LogFilePath     string

	sideChainScanBlockCounter *eurus_ethereum.ScanBlockCounter
	dbProcessor               *ApprovalDBProcessor
	processor                 *ApprovalProcessor
}

func NewApprovalObserver() *ApprovalObserver {
	observer := new(ApprovalObserver)
	observer.Config = NewApprovalObserverConfig()
	observer.ServerConfig = &observer.Config.ServerConfigBase
	observer.logChannel = make(chan types.Log)
	observer.closeChannel = make(chan int)
	observer.processor = new(ApprovalProcessor)
	return observer
}

func (me *ApprovalObserver) InitSideChainScanBlockCounter() {
	var err error
	var prefix = "ApprovalObserver"
	var loggerName string = log.Name.Root
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

func (me *ApprovalObserver) processInit(authClient auth_base.IAuth) {
	//1. creates normal mogger name
	var processorLoggerName = log.Name.Root
	//2. starts getting config from config server (This steps require config server running)
	_, err := me.QueryConfigServer(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}
	//3. Init Db From Config
	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic(err)
	}
	//4. Get ethereum client, web socket
	for {
		_, err = me.InitEthereumWebSocketClientFromConfig(me.ServerConfig)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to Init EthWebSocketClient: ", err.Error())
			time.Sleep(me.Config.GetRetryInterval() * time.Second)
			continue
		}
		_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to Init EthClientClient: ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	err = me.InitBlockSubscriber()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init BlockSubscriber: ", err.Error())
		panic(err)
	}

	me.dbProcessor = NewApprovalDBProcessor(me.DefaultDatabase, me.SlaveDatabase, me.Config, processorLoggerName)

	for {
		instance, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.ServerConfig.InternalSCConfigAddress), me.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to get internal smart contract config", err.Error(), " .Ready to retry")
			time.Sleep(5 * time.Second)
			continue
		}

		me.Config.ApprovalWalletAddress, err = instance.GetApprovalWalletAddress(&bind.CallOpts{}) //for db get data use
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Enable to get the smart contract address. ERROR : ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}
	me.processor = NewApprovalProcessor(me.EthClient, me.dbProcessor, me.Config, processorLoggerName)

	//5. Init side chain scan block counter for rescanning and store in sqlite
	me.InitSideChainScanBlockCounter()
	//6. Use side chain scan block counter for rescanning and store in sqlite
	go me.RunBlockSubscriber()
}

func (me *ApprovalObserver) InitAll() error {
	err := me.QueryAuthServerInfo()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get auth server IP and port")
		panic("Unable to get auth server IP and port")
	}

	wallet, account, hdWalletAddr, _ := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.GetServiceId(), 10))

	me.Config.HdWalletAddress = hdWalletAddr
	me.Config.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*account)

	fmt.Println("Generated Wallet address: ", hdWalletAddr)
	//fmt.Println("Wallet Private Key: ", me.Config.HdWalletPrivateKey)
	log.GetLogger(log.Name.Root).Infoln("Generated Wallet address: ", hdWalletAddr)

	me.ServiceServer.InitAuth(me.processInit, nil)

	return nil
}

func (me *ApprovalObserver) InitBlockSubscriber() error {
	err := me.EthWebSocketClient.SubscribeBlockHeader(0)
	if err != nil {
		return err
	}
	return nil

}

func (me *ApprovalObserver) RunBlockSubscriber() {
	log.GetLogger(log.Name.Root).Debug("Running Block Scanner")
	for {
		block, serverErr := me.EthWebSocketClient.Subscriber.GetLatestBlock(true)
		if serverErr != nil {
			if serverErr.ReturnCode != foundation.NetworkError {
				log.GetLogger(log.Name.Root).Error("Unable to get block by block number: ", serverErr.Error())
			}
			continue
		}
		// log.GetLogger(log.Name.Root).Debug("block scanned: ", block.Number(), "block hash : ", block.TxHash(), "transaction is : ", block.Transactions())

		go me.processor.BlockHandler(block)

		err := me.sideChainScanBlockCounter.UpdateLatestBlock(block.Number())
		if err != nil {
			log.GetLogger(me.sideChainScanBlockCounter.LoggerName).Errorln("UpdateLatestBlock failed on block: ", block.Number().String(), " error: ", err.Error())
		}
	}
}

func (me *ApprovalObserver) rescanSideChainBlock(sender *eurus_ethereum.ScanBlockCounter, eventId decimal.Decimal,
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
