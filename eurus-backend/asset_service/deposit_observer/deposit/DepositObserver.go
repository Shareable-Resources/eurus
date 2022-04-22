package deposit

import (
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DepositObserver struct {
	service_server.ServiceServer
	Config         *DepositObserverConfig
	scProcessor    *DepositSCProcessor
	transProcessor *DepositProcessor

	mainnetBlockCounter   *ethereum.ScanBlockCounter
	sideChainBlockCounter *ethereum.ScanBlockCounter
}

func NewDepositObserver() *DepositObserver {
	observer := new(DepositObserver)
	observer.Config = NewDepositObserverConfig()
	observer.ServerConfig = &observer.Config.ServerConfigBase
	return observer
}

const sweepInvokerAddressKeyName string = "sweepServiceInvokerAddress"

func (me *DepositObserver) LoadConfig(args *service.ServiceCommandLineArgs) error {
	err := me.ServiceServer.LoadConfig(args, me.Config)
	if err != nil {
		return err
	}
	return nil
}

func (me *DepositObserver) InitConfig() {
	me.ServiceServer.ServerConfig = &me.Config.ServerConfigBase
}

func (me *DepositObserver) InitAll() {
	me.ServiceServer.InitAuth(me.processInit, me.ReceivedConfigServiceEvent)
}

func (me *DepositObserver) processInit(authClient auth_base.IAuth) {
	var processorLoggerName string = log.Name.Root
	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}

	if me.Config.MainnetBlockConfirmCount == 0 {
		me.Config.MainnetBlockConfirmCount = 8
	}

	sweepInvokerAddrStr, err := me.QuerySystemConfig(sweepInvokerAddressKeyName)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		panic(err)
	} else if sweepInvokerAddrStr == "" {
		errStr := "Sweep service invoker address is empty"
		log.GetLogger(log.Name.Root).Errorln(errStr)
		panic(errStr)
	}
	me.Config.SweepServiceInvokerAddress = common.HexToAddress(sweepInvokerAddrStr)

	err = me.loadAssetsAmountRelatedSettings()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err)
		panic(err)
	}

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic(err)
	} else {
		log.GetLogger(log.Name.Root).Infoln("Init DB successfully")
	}

	wallet, account, hdWalletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.GetServiceId(), 10))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load wallet key failed: " + err.Error())
		panic("Load wallet key failed: " + err.Error())
	}
	me.Config.HdWalletAddress = hdWalletAddr
	me.Config.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*account)
	fmt.Println("Deposit observer wallet address: " + hdWalletAddr)
	log.GetLogger(log.Name.Root).Infoln("Deposit observer wallet address: " + hdWalletAddr)

	//Connect to mainnet ethereum
	mainnetServerConfig := new(server.ServerConfigBase)
	mainnetServerConfig.EthClientChainID = me.Config.MainnetEthClientChainID
	mainnetServerConfig.EthClientWebSocketPort = me.Config.MainnetEthClientWebSocketPort
	mainnetServerConfig.EthClientWebSocketProtocol = me.Config.MainnetEthClientWebSocketProtocol
	mainnetServerConfig.EthClientWebSocketIP = me.Config.MainnetEthClientWebSocketIP

	me.MainNetEthClient, err = me.InitEthereumWebSocketClientFromConfig(mainnetServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthWebSocketClient for mainnet: ", err.Error())
		panic(err)
	}

	me.EthClient, err = me.InitEthereumWebSocketClientFromConfig(me.ServerConfig)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthWebSocketClient for side chain: ", err.Error())
		panic(err)
	}

	me.mainnetBlockCounter, _ = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeContinuous, me.MainNetEthClient.ChainID.Uint64(), me.mainnetRescanHandler,
		"DepositObserver_Mainnet_Server_"+strconv.FormatInt(me.Config.ServiceId, 10)+".db", log.Name.Root)

	me.sideChainBlockCounter, _ = ethereum.NewScanBlockCounter(ethereum.ScanBlockModeContinuous, me.EthClient.ChainID.Uint64(), me.sideChainRescanHandler,
		"DepositObserver_SideChain_Server_"+strconv.FormatInt(me.Config.ServiceId, 10)+".db", log.Name.Root)

	processorContext := NewDepositProcessorContext(me.DefaultDatabase, me.SlaveDatabase, me.Config, processorLoggerName, me.mainnetBlockCounter, me.sideChainBlockCounter)
	me.scProcessor = NewDepositSCProcessor(me.Config, processorContext)
	err = me.scProcessor.Init()
	if err != nil {
		log.GetLogger(log.Name.Root).Fatalln("Unable to init DepositSCProcessor: ", err.Error())
		panic(err)
	}

	me.transProcessor = NewDepositProcessor(me.Config, me.scProcessor, processorContext)
	err = me.transProcessor.Init()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init DepositProcessor: ", err.Error())
		panic(err)
	}

	err = me.MainNetEthClient.SubscribeBlockHeader(me.Config.MainnetBlockConfirmCount)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to subscribe mainnet block head: ", err.Error())
		panic(err.Error())
	}

	err = me.EthClient.SubscribeBlockHeader(0)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to subscribe side chain block head: ", err.Error())
		panic(err.Error())
	}

	fmt.Println("Deposit Observer start")
	_ = me.sideChainBlockCounter.Start()
	_ = me.mainnetBlockCounter.Start()

	me.transProcessor.RunMainnetBlockSubscriberAsync(me.MainNetEthClient.Subscriber)
	me.transProcessor.RunSideChainBlockSubscriberAsync(me.EthClient.Subscriber)

}

func (me *DepositObserver) QueryAuthServerInfo() error {
	err := me.ServiceServer.QueryAuthServerInfo()
	return err
}

func (me *DepositObserver) ReceivedConfigServiceEvent(message *amqp.Delivery, topic string, contentType string, content []byte) {
	event := new(conf_api.MQConfigServiceEvent)
	err := json.Unmarshal(content, &event)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("ReceivedConfigServiceEvent unmarshal failed: ", err, " content: ", content)
		return
	}
	if topic == "config.system" {
		if event.Action == conf_api.PublishActionInsert || event.Action == conf_api.PublishActionUpdate {
			addr, err := me.QuerySystemConfig(sweepInvokerAddressKeyName)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to query config server for invoker address: ", err)
				return
			}
			me.Config.SweepServiceInvokerAddress = common.HexToAddress(addr)
		}
	}

}

func (me *DepositObserver) loadAssetsAmountRelatedSettings() error {
	// Database table assets and asset_settings have different meaning of asset_name
	// So need to combine both tables to get the sweep trigger settings of tokens
	// AssetName here is USDT, BNB, ....; CurrencyId here is tether, binancecoin, ...
	assets, err := me.QueryAssets()
	if err != nil {
		return err
	}

	assetMappings := make(map[string]string)
	for _, a := range assets {
		assetMappings[a.CurrencyId] = a.AssetName
	}

	// AssetName here is tether, binancecoin, ...
	assetSettings, err := me.QueryAssetSettings()
	if err != nil {
		return err
	}

	// To simplify later use, the final result map use ETH, USDT, ... as key
	for _, setting := range assetSettings {
		// var assetName string
		// found := true
		// if setting.AssetName == "ethereum" {
		// 	assetName = "ETH"
		// } else {
		// 	assetName, found = assetMappings[setting.AssetName]
		// }

		// // Just skip unknown type of tokens
		// if !found {
		// 	continue
		// }

		// me.Config.AssetSettings[assetName] = setting
		me.Config.AssetSettings[setting.AssetName] = setting
	}

	return nil
}

func (me *DepositObserver) mainnetRescanHandler(sender *ethereum.ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal) {
	loggerName := asset.ConstructRescanLoggerName(log.Name.MainnetRescan, me.Config.ServiceId, eventId)

	filePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, filePath, logrus.DebugLevel)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create log file for rescan mainnet event id: ", eventId, " error: ", err)
	}
	defer log.RemoveLogger(loggerName)

	rescanContext := NewDepositProcessorContext(me.DefaultDatabase, me.SlaveDatabase, me.Config, loggerName, me.mainnetBlockCounter, me.sideChainBlockCounter)

	rescanSCProcessor := NewDepositSCProcessor(me.Config, rescanContext)
	err = rescanSCProcessor.Init()
	if err != nil {
		log.GetLogger(loggerName).Errorln("Unable to init DepositSCProcessor. Error: ", err)
		return
	}

	rescanProcessor := NewDepositProcessor(me.Config, rescanSCProcessor, rescanContext)
	err = rescanProcessor.Init()
	if err != nil {
		log.GetLogger(loggerName).Errorln("Unable to init DepositProcessor. Error: ", err)
		return
	}
	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())

	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64() + 1
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)

	for i := int64(0); i < total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting mainnet block number: ", currBlockNum.String())
		for {
			block, err = me.MainNetEthClient.GetBlockByNumber(currBlockNum)
			if err != nil {
				log.GetLogger(loggerName).Errorln("Unable to get block by number. Block number: ", currBlockNum, " Error: ", err)
				time.Sleep(me.Config.GetRetryInterval() * time.Second)
				continue
			}
			break
		}
		log.GetLogger(loggerName).Debugln("Received mainnet block number: ", currBlockNum.String())
		rescanProcessor.processMainnetTransaction(block)
		status, err := me.mainnetBlockCounter.UpdateRescanEvent(eventId, currBlockNum)
		if err != nil {
			log.GetLogger(loggerName).Errorln("UpdateRescanEvent failed for ", currBlockNum.String(), ". Error: ", err.Error())

		} else {
			if status == ethereum.ScanBlockActionFinished {
				break
			}
		}
	}
	log.GetLogger(loggerName).Infoln("Rescan event: ", eventId.String(), " finished")
}

func (me *DepositObserver) sideChainRescanHandler(sender *ethereum.ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal) {
	loggerName := asset.ConstructRescanLoggerName(log.Name.SideChainRescan, me.Config.ServiceId, eventId)
	//loggerName := "SideChainRescan_" + strconv.FormatInt(me.Config.ServiceId, 10) + "_" + eventId.String()
	filePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, filePath, logrus.DebugLevel)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create log file for rescan side chain event id: ", eventId, " error: ", err)
	}
	defer log.RemoveLogger(loggerName)

	rescanContext := NewDepositProcessorContext(me.DefaultDatabase, me.SlaveDatabase, me.Config, loggerName, me.mainnetBlockCounter, me.sideChainBlockCounter)

	rescanSCProcessor := NewDepositSCProcessor(me.Config, rescanContext)
	err = rescanSCProcessor.Init()
	if err != nil {
		log.GetLogger(loggerName).Errorln("Unable to init DepositSCProcessor. Error: ", err)
		return
	}

	rescanProcessor := NewDepositProcessor(me.Config, rescanSCProcessor, rescanContext)
	err = rescanProcessor.Init()
	if err != nil {
		log.GetLogger(loggerName).Errorln("Unable to init DepositProcessor. Error: ", err)
		return
	}

	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())

	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64() + 1
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)

	for i := int64(0); i < total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting sidechain block number: ", currBlockNum.String())
		for {
			block, err = me.EthClient.GetBlockByNumber(currBlockNum)
			if err != nil {
				log.GetLogger(loggerName).Errorln("Unable to get block by number. Block number: ", currBlockNum, " Error: ", err)
				time.Sleep(me.Config.GetRetryInterval() * time.Second)
				continue
			}
			break
		}
		log.GetLogger(loggerName).Debugln("Received sidechain block number: ", currBlockNum.String())
		rescanProcessor.processSideChainTransaction(block)
		status, err := me.sideChainBlockCounter.UpdateRescanEvent(eventId, currBlockNum)
		if err != nil {
			log.GetLogger(loggerName).Errorln("UpdateRescanEvent failed for ", currBlockNum.String(), ". Error: ", err.Error())

		} else {
			if status == ethereum.ScanBlockActionFinished {
				break
			}
		}
	}
	log.GetLogger(loggerName).Infoln("Rescan event: ", eventId.String(), " finished")
}

func (me *DepositObserver) getCurrentMainnetBlock() (*big.Int, common.Hash) {
	return me.transProcessor.currentMainnetBlockNumber, me.transProcessor.currentMainnetBlockHash
}

func (me *DepositObserver) getCurrentSideChainBlock() (*big.Int, common.Hash) {
	return me.transProcessor.currentSideChainBlockNumber, me.transProcessor.currentSideChainBlockHash
}

//TODO Make it a generic class when Golang supports generic
func (me *DepositObserver) InitUDSControlServer(arg *server.CommandLineArguments, handler func(*server.ControlRequestMessage) (bool, string, error)) {

	me.ServiceServer.InitUDSControlServer(arg, func(req *server.ControlRequestMessage) (bool, string, error) {
		var handled bool
		var err error
		var output string
		if handler != nil {
			handled, output, err = handler(req)
		}

		if !handled {
			args := req.Data

			switch req.MethodName {
			case "mainnetblock":
				blockNum, hash := me.getCurrentMainnetBlock()
				if blockNum != nil {
					output += fmt.Sprintln(blockNum.String(), " ", hash.Hex())
				} else {
					output += fmt.Sprintln("No block processing at the moment")
				}
				return true, output, nil
			case "sidechainblock":
				blockNum, hash := me.getCurrentSideChainBlock()
				if blockNum != nil {
					output += fmt.Sprintln(blockNum.String(), " ", hash.Hex())
				} else {
					output += fmt.Sprintln("No block processing at the moment")
				}
				return true, output, nil
			case "enableverboselog":
				if len(args) == 0 {
					return true, output, errors.New("Invalid argument")
				}
				isEnable := args[0]
				if isEnable == "1" {
					me.transProcessor.isVerboseLog = true
					output += fmt.Sprintln("Enabled verbose log")
				} else if isEnable == "0" {
					me.transProcessor.isVerboseLog = false
					output += fmt.Sprintln("Disabled verbose log")
				} else {
					return true, "", errors.New("Invalid argument. Should be 1/0")
				}
				return true, output, nil
			case "help":
				output += fmt.Sprintln("mainnetblock - Print currently processing mainnet block number & hash")
				output += fmt.Sprintln("sidechainblock - Print currently processing side chain block number & hash")
				output += fmt.Sprintln("enableverboselog [1/0] - Enable/disable verbose log mode")
			}
			return false, output, nil

		}
		return false, output, err
	})

}
