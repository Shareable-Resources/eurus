package sweep

import (
	"crypto/ecdsa"
	"encoding/base64"

	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	go_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SweepService struct {
	service_server.ServiceServer
	Config          *SweepServiceConfig
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	scProcessor     *SweepServiceSCProcessor
	processor       *SweepServiceProcessor
	context         *SweepServiceContext
	logger          *logrus.Logger
}

func NewSweepService() *SweepService {
	service := new(SweepService)
	service.Config = NewSweepServiceConfig()
	service.ServerConfig = &service.Config.ServerConfigBase
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 500
	return service
}

func (ss *SweepService) LoadConfig(args *service.ServiceCommandLineArgs) error {
	err := ss.ServiceServer.LoadConfig(args, ss.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Load config failed: ", err)
		return err
	}

	err = decryptConfig(&(ss.Config.CentralizedUserWalletMnemonicPhase))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Decrypt config failed on CentralizedUserWalletMnemonicPhase: ", err)
		return errors.Wrap(err, "Decrypt config failed on CentralizedUserWalletMnemonicPhase")
	}

	err = decryptConfig(&(ss.Config.InvokerPrivateKey))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Decrypt config failed on InvokerPrivateKey: ", err)
		return errors.Wrap(err, "Decrypt config failed on InvokerPrivateKey")
	}

	// More info on how to get the address from private key
	// https://hackernoon.com/how-to-generate-ethereum-addresses-technical-address-generation-explanation-25r3zqo
	// https://www.quicknode.com/guides/web3-sdks/how-to-generate-a-new-ethereum-address-in-go
	privateKeyECDSA, err := go_crypto.HexToECDSA(ss.Config.InvokerPrivateKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Loading InvokerPrivateKey failed: ", err)
		return errors.Wrap(err, "Loading InvokerPrivateKey failed")
	}

	publicKeyECDSA, ok := privateKeyECDSA.Public().(*ecdsa.PublicKey)
	if !ok {
		log.GetLogger(log.Name.Root).Error("Cannot parse the public key")
		return errors.New("Cannot parse the public key")
	}

	// publicKey := hex.EncodeToString(crypto.FromECDSAPub(publicKeyECDSA))
	address := go_crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	ss.Config.InvokerAddress = common.HexToAddress(address)

	return nil
}

func (ss *SweepService) InitLog(filePath string) {
	ss.ServiceServer.InitLog(filePath)
	ss.logger = log.GetLogger(log.Name.Root)
}

func (ss *SweepService) InitAll() {
	ss.ServiceServer.InitAuth(ss.processInit, nil)
}

func (ss *SweepService) processInit(authClient auth_base.IAuth) {
	_, err := ss.QueryConfigServer(ss.Config)
	if err != nil {
		ss.logger.Errorln("Unable to get config from config server: ", err)
		panic(err)
	}

	if ss.Config.QueryReceiptRetryCount == 0 {
		ss.Config.QueryReceiptRetryCount = 300
	}

	err = ss.InitDBFromConfig(ss.ServerConfig)
	if err != nil {
		ss.logger.Errorln("Unable to connect to DB: ", err)
		panic(err)
	} else {
		ss.logger.Infoln("Init DB successfully")
	}

	err = ss.initHTTPServer()
	if err != nil {
		ss.logger.Errorln("Unable to init http server:", err)
		panic(err)
	}

	err = ss.updateInvokerAddress()
	if err != nil {
		ss.logger.Errorln("Unable to process invoker private key:", err)
		panic(err)
	}

	err = ss.loadAssetsNameConfig()
	if err != nil {
		ss.logger.Errorln("Unable to load assets name config:", err)
		panic(err)
	}

	ss.context = NewSweepServiceContext(ss.DefaultDatabase, ss.Config, ss.logger)
	ss.scProcessor = NewSweepServiceSCProcessor(ss.Config, ss.context)
	err = ss.scProcessor.Init()
	if err != nil {
		ss.logger.Fatalln("Unable to init SweepServiceSCProcessor: ", err)
		panic(err)
	}

	ss.processor = NewSweepServiceProcessor(ss.Config, ss.context, ss.scProcessor)
	log.GetLogger(log.Name.Root).Infoln("Sweep service started")
	go ss.runDatabasePoller()
}

func (ss *SweepService) QueryAuthServerInfo() error {
	err := ss.ServiceServer.QueryAuthServerInfo()
	return err
}

func (ss *SweepService) updateInvokerAddress() error {
	const KeyName = "sweepServiceInvokerAddress"

	fmt.Println("Invoker address:", ss.Config.InvokerAddress)
	ss.logger.Infoln("Invoker address:", ss.Config.InvokerAddress)

	// Update invoker address to system config
	// Address is saved in lower case
	request := conf_api.NewAddOrUpdateSystemConfigRequest(KeyName)
	request.Key = KeyName
	request.Value = strings.ToLower(ss.Config.InvokerAddress.Hex())

	response := new(conf_api.AddOrUpdateSystemConfigFullResponse)
	reqRes := api.NewRequestResponse(request, response)

	_, err := ss.SendConfigApiRequest(reqRes)
	if err != nil {
		return err
	}

	if response.GetReturnCode() < int64(foundation.Success) {
		return errors.New(response.GetMessage())
	}

	ss.logger.Infoln("Successfully updated invoker address in system config")
	return nil
}

func (ss *SweepService) initHTTPServer() error {
	var err error
	err = ss.ServerBase.InitHttpServer(ss.Config)
	if err != nil {
		return err
	}
	ss.loginMiddleware.AuthClient = ss.AuthClient
	ss.setupRouter()

	// This is blocking call so run it goroutine
	go func() {
		err = ss.HttpServer.Listen()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (ss *SweepService) setupRouter() {
	ss.HttpServer.Router.Mount(RootPath, ss.configRouter())
}

func (ss *SweepService) configRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(ss.loginMiddleware.VerifyServiceLoginToken)
	r.Use(ss.errorMiddleware.JsonErrorHandler)
	r.Post(SweepAllPath, ss.sweepAll)

	return r
}

func (ss *SweepService) loadAssetsNameConfig() error {
	// Different database table have different meaning for asset_name
	// So need some maps to do the look up
	// AssetName here is USDT, BNB, ....; CurrencyId here is tether, binancecoin, ...
	assets, err := ss.QueryAssets()
	if err != nil {
		return err
	}

	// The words currency and symbol are from wiki
	// https://en.wikipedia.org/wiki/List_of_cryptocurrencies
	for _, asset := range assets {
		if assetNameList, ok := ss.Config.CurrencyToSymbol[asset.CurrencyId]; !ok {
			assetNameList = make([]string, 0)
			assetNameList = append(assetNameList, asset.AssetName)
			ss.Config.CurrencyToSymbol[asset.CurrencyId] = assetNameList
		} else {
			assetNameList = append(assetNameList, asset.AssetName)
			ss.Config.CurrencyToSymbol[asset.CurrencyId] = assetNameList
		}
		ss.Config.SymbolToCurrency[asset.AssetName] = asset.CurrencyId
	}

	ethList, ok := ss.Config.CurrencyToSymbol["ethereum"]
	if !ok {
		ethList = make([]string, 0)
		ethList = append(ethList, "ETH")
		ss.Config.CurrencyToSymbol["ethereum"] = ethList
	} else {
		ethList = append(ethList, "ETH")
		ss.Config.CurrencyToSymbol["ethereum"] = ethList
	}
	ss.Config.SymbolToCurrency["ETH"] = "ethereum"

	// AssetName here is tether, binancecoin, ...
	assetSettings, err := ss.QueryAssetSettings()
	if err != nil {
		return err
	}

	// To simplify later use, the final result map use ETH, USDT, ... as key
	for _, setting := range assetSettings {
		ss.Config.AssetSettings[setting.AssetName] = setting
		// var assetName string
		// found := true
		// if setting.AssetName == "ethereum" {
		// 	assetName = "ETH"
		// } else {
		// 	assetName, found = ss.Config.CurrencyToSymbol[setting.AssetName]
		// }

		// // Just skip unknown type of tokens
		// if !found {
		// 	continue
		// }

		// ss.Config.AssetSettings[assetName] = setting
	}

	return nil
}

func (ss *SweepService) runDatabasePoller() {
	const DefaultPollingInterval = 60

	// Polling interval setting in DB, if it is invalid, use default value
	interval := ss.Config.DBPollingInterval
	if interval <= 0 {
		ss.logger.Warnln("Invalid DBPollingInterval in config, will use default value instead")
		interval = DefaultPollingInterval
	}

	ss.logger.Infoln("Start polling database every", interval, "second(s)")

	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	for ; true; <-t.C {
		ss.processor.PollDatabase()
	}
}

func (ss *SweepService) sweepAll(writer http.ResponseWriter, req *http.Request) {
	reqObj := request.RequestBase{}

	go ss.runCheckNeedForSweep()

	err := api.HttpWriteResponseWithStatusCode(writer, &reqObj, nil, http.StatusNoContent)
	if err != nil {
		ss.context.logger.Errorln(err)
	}
}

func (ss *SweepService) runCheckNeedForSweep() {
	users, err := DBGetCentralizedUsers(ss.context)
	if err != nil {
		ss.context.logger.Errorln("Failed to get centralized users to perform sweep checking,", err)
		return
	}

	ss.context.logger.Infoln("Got", len(users), "centralized user(s) from database, will run sweep checking for them")

	err = ss.processor.CheckNeedForSweep(users)
	if err != nil {
		ss.context.logger.Errorln("Error returned from CheckNeedForSweep process,", err)
	}
}

func decryptConfig(config *string) error {
	decrypted, err := secret.DecryptConfigValue(*config)
	if err != nil {
		return err
	}

	data, err := base64.StdEncoding.DecodeString(decrypted)
	if err != nil {
		return err
	}

	*config = string(data)

	return nil
}
