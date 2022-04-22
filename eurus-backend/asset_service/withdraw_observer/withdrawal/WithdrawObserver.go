package withdrawal

import (
	"encoding/json"
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/ethereum"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type WithdrawObserver struct {
	service_server.ServiceServer
	Config *WithdrawObserverConfig

	logProcessor          *WithdrawProcessor
	sideChainBlockCounter *eurus_ethereum.ScanBlockCounter

	sideChainCurrentBlockNumber *big.Int
	sideChainCurrentBlockHash   common.Hash
	isVerboseLog                bool
}

func NewWithdrawObserver() *WithdrawObserver {
	observer := new(WithdrawObserver)
	observer.Config = NewWithdrawObserverConfig()
	observer.ServerConfig = &observer.Config.ServerConfigBase
	observer.logProcessor = NewWithdrawProcessor(observer.Config, log.Name.Root)

	return observer
}

func (me *WithdrawObserver) LoadConfig(args *service.ServiceCommandLineArgs) error {
	err := me.ServiceServer.LoadConfig(args, me.Config)
	return err
}

func (me *WithdrawObserver) InitAll() error {
	err := me.QueryAuthServerInfo()
	if err != nil {
		fmt.Println("Unable to get auth server IP and port")
		log.GetLogger(log.Name.Root).Fatal("Unable to get auth server IP and port")
	}

	me.ServiceServer.InitAuth(me.processInit, nil)
	return nil
}

func (me *WithdrawObserver) Close() {

}

func (me *WithdrawObserver) processInit(authClient auth_base.IAuth) {
	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	if me.Config.SideChainGasLimit <= 0 {
		me.Config.SideChainGasLimit = 10000000
	}

	wallet, account, hdWalletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.FormatInt(me.Config.GetServiceId(), 10))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load wallet key failed: ", err)
		panic("Load wallet key failed: " + err.Error())
	}
	me.Config.HdWalletAddress = hdWalletAddr
	me.Config.HdWalletPrivateKey, _ = wallet.PrivateKeyHex(*account)

	log.GetLogger(log.Name.Root).Infoln("Withdraw observe Wallet address: ", me.Config.HdWalletAddress)

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to connect to DB: ", err)
		panic("Unable to connect to DB: " + err.Error())
	} else {
		log.GetLogger(log.Name.Root).Infoln("Init DB successfully")
	}

	_, err = me.InitEthereumWebSocketClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthWebSocketClient: ", err.Error())
		panic("Unable to Init EthWebSocketClient: " + err.Error())
	} else {
		log.GetLogger(log.Name.Root).Infoln("Init ETH web socket successfully")
	}

	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to Init EthWebSocketClient: ", err.Error())
		panic("Unable to Init EthWebSocketClient: " + err.Error())
	} else {
		log.GetLogger(log.Name.Root).Infoln("Init ETH client successfully")
	}
	err = me.logProcessor.Init(me.DefaultDatabase, me.SlaveDatabase)
	if err != nil {
		panic(err)
	}

	for {
		internalSCConfig, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.Config.InternalSCConfigAddress), me.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to create internal smart contract instance: ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		addr, err := internalSCConfig.GetWithdrawSmartContract(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to GetWithdrawSmartContract: ", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		me.Config.WithdrawSmartContractAddr = addr
		log.GetLogger(log.Name.Root).Infoln("WithdrawSmartContract address: ", addr.Hex())

		break
	}

	me.sideChainBlockCounter, _ = eurus_ethereum.NewScanBlockCounter(eurus_ethereum.ScanBlockModeContinuous,
		me.EthClient.ChainID.Uint64(), me.rescanSideChain, "WithdrawObserver_SideChain_"+strconv.FormatInt(me.Config.ServiceId, 10)+".db", log.Name.Root)

	err = me.EthWebSocketClient.SubscribeBlockHeader(0)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create block subscriber. Error: ", err.Error())
		panic("Unable to create block subscriber. " + err.Error())
	}

	fmt.Println("Withdraw observer started. Wallet address: ", me.Config.HdWalletAddress)

	log.GetLogger(log.Name.Root).Infoln("Withdraw observer started")

	_ = me.sideChainBlockCounter.Start()

	go me.processBlockAsync()

}

func (me *WithdrawObserver) processBlockAsync() {
	var isErrorOccured bool
	for {
		block, err := me.EthWebSocketClient.Subscriber.GetLatestBlock(true)
		if err != nil {
			if !isErrorOccured || me.isVerboseLog {
				log.GetLogger(log.Name.Root).Errorln("Unable to get latest block. Error: ", err.Error(), " wait for ", int(me.Config.GetRetryInterval()), " second")
				isErrorOccured = true
			}
			time.Sleep(time.Second * me.Config.GetRetryInterval())
			continue
		}

		if isErrorOccured {
			isErrorOccured = false
			log.GetLogger(log.Name.Root).Infoln("Get latest block resumed")
		}

		me.sideChainCurrentBlockNumber = block.Number()
		me.sideChainCurrentBlockHash = block.Hash()
		if me.isVerboseLog {
			log.GetLogger(log.Name.Root).Debugln("Start processing block number: ", block.Number(), " block hash: ", block.Hash().String())
		}
		me.logProcessor.ProcessTransaction(block)
		if me.isVerboseLog {
			log.GetLogger(log.Name.Root).Debugln("End processing block number: ", block.Number(), " block hash: ", block.Hash().String())
		}
		updateErr := me.sideChainBlockCounter.UpdateLatestBlock(block.Number())
		if updateErr != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to update latest block number: ", updateErr.Error())
		}
	}
}

func (me *WithdrawObserver) rescanSideChain(sender *eurus_ethereum.ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal) {

	loggerName := asset.ConstructRescanLoggerName(log.Name.SideChainRescan, me.Config.ServiceId, eventId)

	filePath := asset.ConstructLogFilePath(loggerName, me.Config.LogFilePath)
	_, err := log.NewLogger(loggerName, filePath, logrus.DebugLevel)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create log file for rescan side chain event id: ", eventId, " error: ", err)
	}
	defer log.RemoveLogger(loggerName)

	rescanProcessor := NewWithdrawProcessor(me.Config, loggerName)
	err = rescanProcessor.Init(me.DefaultDatabase, me.SlaveDatabase)
	if err != nil {
		log.GetLogger(loggerName).Errorln("Unable to init WithdrawProcessor. Error: ", err)
		return
	}
	log.GetLogger(loggerName).Infoln("Start rescan for event id: ", eventId.String(), " from block: ", from.String(), " to block: ", to.String())

	delta := to.Sub(from)
	var total int64 = delta.BigInt().Int64() + 1
	var currBlockNum *big.Int = from.BigInt()
	var one *big.Int = big.NewInt(1)
	var isErrorOccurs bool
	for i := int64(0); i < total; i, currBlockNum = i+1, currBlockNum.Add(currBlockNum, one) {
		var block *types.Block
		var err error
		log.GetLogger(loggerName).Debugln("Start getting block number: ", currBlockNum.String())

		for {
			block, err = me.EthClient.GetBlockByNumber(currBlockNum)
			if err != nil {
				if !isErrorOccurs || me.isVerboseLog {
					log.GetLogger(loggerName).Errorln("Unable to get block by number. Block number: ", currBlockNum, " Error: ", err, " wait for ", int(me.Config.GetRetryInterval()), " second")
					isErrorOccurs = true
				}
				time.Sleep(me.Config.GetRetryInterval() * time.Second)
				continue
			}
			break
		}

		if isErrorOccurs {
			isErrorOccurs = false
			log.GetLogger(loggerName).Infoln("Resume getting block")
		}
		log.GetLogger(loggerName).Debugln("Received block number: ", currBlockNum.String())
		rescanProcessor.ProcessTransaction(block)
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

func logError(message string, err error, transLog *types.Log, requestTransId uint64, burnTransId uint64, newTransHash string) {
	errByte, _ := json.Marshal(err)

	log.GetLogger(log.Name.Root).Errorf("%s - Error: %s. JSON: %s, Created transHash: %s, burn transId: %d, request transId: %d, request transHash: %s\r\n",
		message, err.Error(), string(errByte), newTransHash, burnTransId, requestTransId, transLog.TxHash.Hex())
}

func (me *WithdrawObserver) getCurrentSideChainBlock() (*big.Int, common.Hash) {
	return me.sideChainCurrentBlockNumber, me.sideChainCurrentBlockHash
}

func (me *WithdrawObserver) InitUDSControlServer(arg *server.CommandLineArguments, handler func(req *server.ControlRequestMessage) (bool, string, error)) {

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
			case "setmainnettransfergastipcap":
				if len(args) == 0 {
					return true, output, errors.New("Missing Gwei value")
				}

				val, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return true, output, errors.New("Argument is not an integer")
				}
				me.Config.MainnetTransferGasTipCap = val
				return true, "MainnetTransferGasTipCap set successfully", nil

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
					me.isVerboseLog = true
					output += fmt.Sprintln("Enabled verbose log")
				} else if isEnable == "0" {
					me.isVerboseLog = false
					output += fmt.Sprintln("Disabled verbose log")
				} else {
					return true, "", errors.New("Invalid argument. Should be 1/0")
				}
				return true, output, nil
			case "resumefromstate60":
				fallthrough
			case "resumefromstate70":
				if len(args) == 0 {
					return true, output, errors.New("Missing request trans hash")
				}
				err := me.logProcessor.processResumbitTransferTokenToMainnet(args[0])
				if err != nil {
					return true, "", err
				} else {
					output += "Request submit successfully, look up for log to see the progress\r\n"
				}
				return true, output, nil
			case "help":
				output += fmt.Sprintln("sidechainblock - Print currently processing side chain block number & hash")
				output += fmt.Sprintln("enableverboselog [1/0] - Enable/disable verbose log mode")
				output += fmt.Sprintln("setMainnetTransferGasTipCap - <Wei value>")
				output += fmt.Sprintln("ResumeFromState60 [request trans hash] - Resume transaction starts from state 60")
				output += fmt.Sprintln("ResumeFromState70 [request trans hash] - Resubmit transfer token request at mainnet")
			}
			return false, output, nil

		}
		return false, output, err
	})

}
