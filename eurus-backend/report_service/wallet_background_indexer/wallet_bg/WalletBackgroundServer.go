package wallet_bg

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	wallet_bg_model "eurus-backend/report_service/wallet_background_indexer/wallet_bg/model"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	// "strconv"
)

type AssetAddressMap map[common.Address]string
type AssetNameMap map[string]common.Address

const OneDaySeconds int64 = 24 * 3600
const poolSize int = 60

type WalletBackgroundIndexer struct {
	service_server.ServiceServer
	Config                      *WalletBackgroundIndexerConfig
	mainnetAssetNameToAddress   AssetNameMap
	sideChainAssetNameToAddress AssetNameMap
	eunAddressMap               AssetNameMap
	ethAddressMap               AssetNameMap
	// mainnetHotWalletAddress     common.Address
	// adminFeeWalletAddress       common.Address
	// gasFeeWalletAddress         common.Address
	// sweepServerWalletAddress    common.Address // Sweep Server Wallet Address of mainnet
	// eurusUserDepositAddress     common.Address
	// marketRegWalletAddress      common.Address
	// observerServerAddressList   []conf_api.AuthService
	eurusInternalConfigContract *mainnet_contract.EurusInternalConfig // To Get Asset name and hot wallet address
	internalSCConfigContract    *contract.InternalSmartContractConfig // To Get ERC 20 contract address by asset name
	externalSCConfigContract    *contract.ExternalSmartContractConfig
	isInitAuthSuccess           chan bool
	noOfGoRoutines              chan int

	sideChainEthClientPool chan *ethereum.EthClient
	mainnetEthClientPool   chan *ethereum.EthClient

	noOfRecordsSaved int
	noOfRecordsFail  int
	closeInMinutes   int
	startTime        time.Time
	endTime          time.Time
}

func NewWalletBackgroundServer() *WalletBackgroundIndexer {

	walletBgServer := new(WalletBackgroundIndexer)

	walletBgServer.Config = NewWalletBackgroundIndexerConfig()
	walletBgServer.ServerConfig = &walletBgServer.Config.ServerConfigBase
	walletBgServer.isInitAuthSuccess = make(chan bool)
	walletBgServer.mainnetAssetNameToAddress = make(AssetNameMap)
	walletBgServer.eunAddressMap = make(AssetNameMap)
	walletBgServer.ethAddressMap = make(AssetNameMap)
	walletBgServer.sideChainAssetNameToAddress = make(AssetNameMap)

	walletBgServer.noOfGoRoutines = make(chan int, poolSize*2)

	walletBgServer.sideChainEthClientPool = make(chan *ethereum.EthClient, poolSize)
	walletBgServer.mainnetEthClientPool = make(chan *ethereum.EthClient, poolSize)

	walletBgServer.noOfRecordsSaved = 0 // Per addresss' s balances, not per asset balances, success save
	walletBgServer.noOfRecordsFail = 0  // Per addresss' s balances, not per asset balances, fail save
	walletBgServer.closeInMinutes = 3   // If no go routine running in 3 minutes, program exits
	return walletBgServer
}

func (me *WalletBackgroundIndexer) LoadConfig(args *service.ServiceCommandLineArgs) error {
	return me.ServiceServer.LoadConfigWithSetting(args, me.Config, false, false)
}

func (me *WalletBackgroundIndexer) InitAll() {

	me.PrintVersion()

	go func() { me.ServiceServer.InitAuth(me.authLoginHandler, nil) }()
	for {
		select {
		case <-me.isInitAuthSuccess:
			fmt.Println("Init Auth Success")
		}
	}
}

func (me *WalletBackgroundIndexer) InitSmartContractInstance() error {
	var err error
	//EurusInternalConfig
	me.eurusInternalConfigContract, err = mainnet_contract.NewEurusInternalConfig(common.HexToAddress(me.Config.EurusInternalConfigAddress), me.MainNetEthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init eurusInternalConfigContract: ", err.Error())
		return err
	}
	//InternalSmartContractConfig
	me.internalSCConfigContract, err = contract.NewInternalSmartContractConfig(common.HexToAddress(me.Config.InternalSCConfigAddress), me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init internalSCConfigContract: ", err.Error())
		return err
	}

	me.externalSCConfigContract, err = contract.NewExternalSmartContractConfig(common.HexToAddress(me.Config.ExternalSCConfigAddress), me.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init externalSCConfigContract: ", err.Error())
		return err
	}
	return nil
}

// func (me *WalletBackgroundIndexer) InitObserverServerWalletAddressMap() error {
// 	for {
// 		authServiceInfoReq := conf_api.NewQueryConfigAuthInfoRequest()
// 		authServiceInfoReq.ServiceId = me.ServerConfig.ServiceId

// 		resp := new(auth.QueryConfigAuthInfoResponse)
// 		reqRes := api.NewRequestResponse(authServiceInfoReq, resp)

// 		_, err := me.SendConfigApiRequest(reqRes)
// 		if err != nil {
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}

// 		if resp.GetReturnCode() < int64(foundation.Success) {
// 			log.GetLogger(log.Name.Root).Errorln(resp.GetMessage())
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}

// 		err = conf_api.ConfigMapListToServerConfig(resp.Data.ConfigData, me.ServerConfig)
// 		if err != nil {
// 			log.GetLogger(log.Name.Root).Errorln("Unable to deserialize server config: ", err.Error())
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}

// 		me.ConvertAuthDataToAddressMap(&resp.Data.AuthData)

// 		log.GetLogger(log.Name.Root).Infoln("Service info count: ", len(resp.Data.AuthData), " loaded")

// 		//Sweep Server wallet
// 		sweepServerKey := "sweepServiceInvokerAddress"
// 		systemConfigReq := conf_api.NewGetSystemConfigRequest(sweepServerKey)
// 		_, serviceToken, err := me.AuthClient.VerifyLoginToken(me.AuthClient.GetLoginToken(), auth_base.VerifyModeService)
// 		systemConfigReq.LoginToken = serviceToken
// 		resp2 := new(conf_api.GetSystemConfigFullResponse)
// 		reqRes = api.NewRequestResponse(systemConfigReq, resp2)
// 		_, err = me.SendConfigApiRequest(reqRes)
// 		if err != nil {
// 			log.GetLogger(log.Name.Root).Errorln("Unable to get system config from config server: ", err.Error())
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}

// 		if resp2.GetReturnCode() < int64(foundation.Success) {
// 			log.GetLogger(log.Name.Root).Errorln(resp.GetMessage())
// 			time.Sleep(5 * time.Second)
// 			continue
// 		}
// 		me.sweepServerWalletAddress = common.HexToAddress(resp2.Data.Value)
// 		break
// 	}
// 	return nil
// }

// Query 1(Deposit),2(Withdraw),3(Approval),4(User),6(UserObserver),8(Config)
// Only query those that are within the above group and has wallet address in table [auth_services]

// func (me *WalletBackgroundIndexer) ConvertAuthDataToAddressMap(data *[]conf_api.AuthService) {
// 	serviceGroupNeeded := []conf_api.ServiceGroupId{
// 		conf_api.ServiceGroupDeposit,      //1
// 		conf_api.ServiceGroupWithdraw,     //2
// 		conf_api.ServiceGroupApproval,     //3
// 		conf_api.ServiceGroupUser,         //4
// 		conf_api.ServiceGroupUserObserver, //6
// 		conf_api.ServiceGroupConfig}       //8
// 	for _, authInfo := range *data {
// 		if contains(serviceGroupNeeded, conf_api.ServiceGroupId(authInfo.ServiceGroupId)) && authInfo.WalletAddress != "" {
// 			me.observerServerAddressList = append(me.observerServerAddressList, authInfo)
// 		}
// 	}
// }

func (me *WalletBackgroundIndexer) InitAssetsMap() error {
	//AssetNameList
	assetNameList, addrList, err := me.eurusInternalConfigContract.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init assets map: ", err.Error())
		return err
	}
	me.mainnetAssetNameToAddress["ETH"] = common.Address{}
	for i, assetName := range assetNameList {
		me.mainnetAssetNameToAddress[assetName] = addrList[i]
	}
	// //Hot Wallet
	// me.mainnetHotWalletAddress, err = me.eurusInternalConfigContract.PlatformWalletAddress(&bind.CallOpts{})
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to platformWalletAddress: ", err.Error())
	// 	return err
	// }

	// //Eurus User Deposit address
	// me.eurusUserDepositAddress, err = me.eurusInternalConfigContract.EurusUserDepositAddress(&bind.CallOpts{})
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to get Eurus user deposit address: ", err.Error())
	// 	return err
	// }

	// me.marketRegWalletAddress, err = me.internalSCConfigContract.GetMarketingRegWalletAddress(&bind.CallOpts{})
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to get market reg wallet address: ", err.Error())
	// 	return err
	// }
	// //Internal address map for asset name smart contract

	// hexStr := hex.EncodeToString(me.mainnetHotWalletAddress.Bytes())
	// //Hot Wallet Hex Address in String
	// fmt.Println("Hot Wallet Address")
	// fmt.Printf("%s    - %s \n", hexStr, "Hot wallet address")
	//Asset Hex Address in String
	fmt.Println("Asset Address")
	fmt.Println("[Smart Contract Address (Mainnet)]            -     [Smart Contract Address (Side Chain)]")

	assetNameList, addressList, err := me.externalSCConfigContract.GetAssetAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Fail to Sidechain asset address list: ", err)
		return err
	}

	me.sideChainAssetNameToAddress[asset.EurusTokenName] = common.Address{}

	for index, assetName := range assetNameList {
		me.sideChainAssetNameToAddress[assetName] = addressList[index]
		fmt.Printf("%s - %s - %s\n", addressList[index].String(), assetName, me.sideChainAssetNameToAddress[assetName])
	}

	me.ethAddressMap["ETH"] = common.Address{}
	//EUN
	me.eunAddressMap[asset.EurusTokenName] = *new(common.Address)

	// internalSC, err := contract.NewInternalSmartContractConfig(common.HexToAddress(me.Config.InternalSCConfigAddress), me.EthClient.Client)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Error("Unable to create InternalSmartContractConfig: ", err)
	// 	return err
	// }

	// me.adminFeeWalletAddress, err = internalSC.GetAdminFeeWalletAddress(&bind.CallOpts{})
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Error("Unable to GetAdminFeeWalletAddress: ", err)
	// 	return err
	// }
	// fmt.Println("Admin fee wallet address: ", me.adminFeeWalletAddress.Hex())

	// me.gasFeeWalletAddress, err = internalSC.GetGasFeeWalletAddress(&bind.CallOpts{})
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Error("Unable to GetGasFeeWalletAddress: ", err)
	// 	return err
	// }
	// fmt.Println("Gas fee wallet address: ", me.gasFeeWalletAddress.Hex())

	fmt.Println("Init Asset Name Map Success")
	log.GetLogger(log.Name.Root).Info("Init assetname map success")
	return nil
}

func (me *WalletBackgroundIndexer) authLoginHandler(authClient auth_base.IAuth) {
	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	err = me.processInit()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

}

func (me *WalletBackgroundIndexer) GetBalance(ethClient *ethereum.EthClient, isMainnet bool, address common.Address, addressMap AssetNameMap, assetName string) (*big.Int, error) {
	if isMainnet {
		return me.GetMainnetBalance(ethClient, address, addressMap, assetName)
	} else {
		return me.GetSideChainBalance(ethClient, address, addressMap, assetName)
	}
}

//End the indexer if go routines did not run for (closeInMins) mins
func (me *WalletBackgroundIndexer) closeAllIfGoRoutinesHasStoppedFor(closeInMins int) {
	stoppedAfterCount := closeInMins
	count := 0
	for {
		time.Sleep(time.Minute * 1)

		if len(me.noOfGoRoutines) == 0 {
			count++
			stopMsg := "Go rountines stopped for " + strconv.Itoa(count) + " minutes"
			log.GetLogger(log.Name.Root).Infoln(stopMsg)
			fmt.Println(stopMsg)
		} else {
			resetMsg := "Reset counter"
			log.GetLogger(log.Name.Root).Infoln(resetMsg)
			fmt.Println(resetMsg)
			count = 0
		}

		if count >= stoppedAfterCount {
			close(me.noOfGoRoutines)
			close(me.mainnetEthClientPool)
			close(me.sideChainEthClientPool)
			exitMsg := "[Exit] WalletBackgroundIndexer exit after " + strconv.Itoa(stoppedAfterCount) + " minutes of no active go routines"
			failMsg := "Fail " + strconv.Itoa(me.noOfRecordsFail) + " records"
			successMsg := "Insert " + strconv.Itoa(me.noOfRecordsSaved) + " records"
			log.GetLogger(log.Name.Root).Infoln(exitMsg)
			log.GetLogger(log.Name.Root).Infoln(failMsg)
			log.GetLogger(log.Name.Root).Infoln(successMsg)
			fmt.Println(exitMsg)
			fmt.Println(failMsg)
			fmt.Println(successMsg)
			me.endTime = time.Now()
			diff := me.endTime.Sub(me.startTime)
			log.GetLogger(log.Name.Root).Infoln("End to get balance at ", me.endTime.Format(time.RFC1123))
			log.GetLogger(log.Name.Root).Infoln("Used Time [" + fmtDuration(diff) + "]")
			fmt.Println("End to get balance at : " + me.endTime.Format(time.RFC1123))
			fmt.Println("Used Time [" + fmtDuration(diff) + "]")
			os.Exit(0)
			break
		}
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (me *WalletBackgroundIndexer) GetMainnetBalance(ethClient *ethereum.EthClient, address common.Address, addressMap AssetNameMap, assetName string) (*big.Int, error) {
	if assetName == "ETH" {
		return ethClient.GetBalance(address)
	}

	contractAddress, found := addressMap[assetName]
	if !found {
		return nil, errors.Errorf("Cannot find the address for asset %v", assetName)
	}

	return me.getTokenBalance(address, contractAddress, ethClient.Client)
}

func (me *WalletBackgroundIndexer) GetSideChainBalance(ethClient *ethereum.EthClient, address common.Address, addressMap AssetNameMap, assetName string) (*big.Int, error) {
	if assetName == asset.EurusTokenName {
		return ethClient.GetBalance(address)
	}

	contractAddress, found := addressMap[assetName]
	if !found {
		return nil, errors.Errorf("Cannot find the address for asset %v", assetName)
	}
	return me.getTokenBalance(address, contractAddress, ethClient.Client)
}

func (me *WalletBackgroundIndexer) processInit() error {

	var err error

	if me.Config.InvokerWalletAddrJson != "" {
		err = json.Unmarshal([]byte(me.Config.InvokerWalletAddrJson), &me.Config.InvokerWalletAddrList)
		if err != nil {
			panic("InvokerWalletAddrJson invalid JSON format: " + err.Error())
		}
	}

	err = me.InitDBFromConfig(&me.Config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to DB: ", err)
		panic(err)
	}

	//Test DB connection
	_, err = me.DefaultDatabase.GetConn()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connection DB: ", err.Error())
	} else {
		log.GetLogger(log.Name.Root).Infoln("Login DB successfully")
	}

	me.isInitAuthSuccess <- true

	walletConfigList, err := DbQueryWalletBalanceConfig(me.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Panic("Query wallet balance config error: ", err)
	}

	req := conf_api.NewQueryConfigAuthInfoRequest()
	res := new(conf_api.QueryConfigAuthResponse)

	reqRes := api.NewRequestResponse(req, res)

	reqRes, err = me.SendConfigApiRequest(reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Panic("Unable to query auth service info from config server: ", err)
	}

	if res.ReturnCode != int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Panic("Query auth server response error code: ", res.ReturnCode, " message: ", res.Message)
	}

	walletList := me.prepareWalletConfigList(res.Data.AuthData, walletConfigList)
	//Connect to eurus ethereum
	me.EthClient, err = me.InitEurusEthereumClientFromConfig()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic(err)
	}
	//Connect to mainnet ethereum
	me.MainNetEthClient, err = me.InitMainNetEthereumClientFromConfig()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic(err)
	}

	log.GetLogger(log.Name.Root).Infoln("Connecting Side chain RPC ", poolSize, " times")
	err = me.connectEthClientPool(me.sideChainEthClientPool, me.Config.EthClientWebSocketProtocol, me.Config.EthClientWebSocketIP, me.Config.EthClientWebSocketPort, int64(me.Config.EthClientChainID))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Side chain Eth client pool connect error: ", err)
		panic("Side chain Eth client pool connect error: " + err.Error())
	}

	log.GetLogger(log.Name.Root).Infoln("Connecting Ethereum mainnet RPC ", poolSize, " times")
	err = me.connectEthClientPool(me.mainnetEthClientPool, me.Config.MainnetEthClientWebSocketProtocol, me.Config.MainnetEthClientWebSocketIP, me.Config.MainnetEthClientWebSocketPort, int64(me.Config.MainnetEthClientChainID))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("mainnet Eth client pool connect error: ", err)
		panic("mainnet Eth client pool connect error: " + err.Error())
	}

	//Init Smart Contract
	err = me.InitSmartContractInstance()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init smart contract: ", err.Error())
		panic(err)
	}

	//Init assets map
	err = me.InitAssetsMap()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to init assets map: ", err.Error())
		panic(err)
	}

	me.startTime = time.Now()
	fmt.Println("Start to get balance at : " + me.startTime.Format(time.RFC1123))
	log.GetLogger(log.Name.Root).Infoln("Start to get balance at : " + me.startTime.Format(time.RFC1123))

	currBlockNum, err := me.EthClient.Client.BlockNumber(context.Background())
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get current block number: ", err.Error())
		panic(err)
	}

	// Processing ERC20 smart contract total supply
	log.GetLogger(log.Name.Root).Infoln("Starting to get total supply from ERC20 smart contract")
	blockNum := big.NewInt(0)
	blockNum.SetUint64(currBlockNum)

	unixDay := int64(me.startTime.Unix()/OneDaySeconds) * OneDaySeconds
	yesterday := time.Unix(unixDay, 0).AddDate(0, 0, -1)
	me.ProcessTotalSupply(&me.sideChainAssetNameToAddress, blockNum, &yesterday, int(me.EthClient.ChainID.Int64()))

	for _, walletConfig := range walletList {
		var address common.Address = common.HexToAddress(walletConfig.WalletAddress)
		me.noOfGoRoutines <- 1
		fmt.Printf("Add go routines from (%d) ->(%d) for %s, wallet: %s, asset: %s\n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines), walletConfig.Description, walletConfig.WalletAddress, walletConfig.AssetName)
		var assetMap AssetNameMap
		if walletConfig.ChainId == me.Config.EthClientChainID {
			if walletConfig.AssetName == asset.EurusTokenName {
				assetMap = me.eunAddressMap
			} else {
				assetMap = me.sideChainAssetNameToAddress
			}

			go me.saveWalletAddress(me.sideChainEthClientPool, conf_api.ServiceGroupId(walletConfig.WalletType), &address, assetMap, me.Config.EthClientChainID, false, nil)

		} else {
			if walletConfig.AssetName == "ETH" {
				assetMap = me.ethAddressMap
			} else {
				assetMap = me.mainnetAssetNameToAddress
			}
			go me.saveWalletAddress(me.mainnetEthClientPool, conf_api.ServiceGroupId(walletConfig.WalletType), &address, assetMap, me.Config.MainnetEthClientChainID, true, nil)
		}

	}

	// 1. Save hot wallet
	// log.GetLogger(log.Name.Root).Infoln("Starting to get mainnet hot wallet address balance ----")
	// me.noOfGoRoutines <- 1
	// fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupMainHotWallet, &me.mainnetHotWalletAddress, me.mainnetAssetNameToAddress, me.Config.MainnetEthClientChainID, true, nil)

	// //2. Save Cold Wallet
	// log.GetLogger(log.Name.Root).Infoln("Starting to get mainnet EurusUserDeposit balance ----")
	// me.noOfGoRoutines <- 1
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupMainColdWallet, &me.eurusUserDepositAddress, me.mainnetAssetNameToAddress, me.Config.MainnetEthClientChainID, true, nil)

	log.GetLogger(log.Name.Root).Infoln("Starting to get mainnet centralized user cold wallet balance ----")
	cenUserList, err := DbGetCentralizedUserWalletAddressesFromDb(me.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get centralized user wallet addresses: ", err.Error())
		panic("Unable to get centralized user wallet addresses: " + err.Error())
	}
	for _, user := range *cenUserList {
		var address common.Address = common.HexToAddress(user.MainnetWalletAddress)
		var userId = user.Id
		me.noOfGoRoutines <- 1
		fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
		go me.saveWalletAddress(me.mainnetEthClientPool, conf_api.WalletBalanceGroupMainColdWallet, &address, me.mainnetAssetNameToAddress, me.Config.MainnetEthClientChainID, true, &userId)
	}

	//3. Side chain wallet
	log.GetLogger(log.Name.Root).Infoln("Starting to get sidechain wallet balance ----")
	userList, err := DbGetUserWalletAddressesFromDb(me.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get user wallet addresses: ", err.Error())
		panic(err)
	}

	for _, user := range *userList {
		var address common.Address = common.HexToAddress(user.WalletAddress)
		var userId = user.Id
		me.noOfGoRoutines <- 1
		fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
		go me.saveWalletAddress(me.sideChainEthClientPool, conf_api.WalletBalanceGroupUserWallet, &address, me.sideChainAssetNameToAddress, me.Config.EthClientChainID, false, &userId)
	}

	//Get server/observer address from config server, assign to observerServerAddressList
	// err = me.InitObserverServerWalletAddressMap()
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to obtain auth_servers from config server: ", err.Error())
	// 	panic(err)
	// }

	// //4. Save Observer' s and Server' s Wallet in mainnet
	// //Only query config server, withdrawObs, depositObs, approvalObs, userServer, userObserver
	// for _, authData := range me.observerServerAddressList {
	// 	var address common.Address = common.HexToAddress(authData.WalletAddress)
	// 	me.noOfGoRoutines <- 1
	// 	fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// 	go me.saveWalletAddress(conf_api.ServiceGroupId(authData.ServiceGroupId), &address, me.eunAddressMap, me.Config.EthClientChainID, false, nil)
	// }

	//Withdraw observer mainnet balance
	// var withdrawObserverList []conf_api.AuthService = make([]conf_api.AuthService, 0)
	// for _, authData := range me.observerServerAddressList {
	// 	if authData.ServiceGroupId == int(conf_api.ServiceGroupWithdraw) {
	// 		withdrawObserverList = append(withdrawObserverList, authData)
	// 	}
	// }

	// for _, withdrawObsData := range withdrawObserverList {
	// 	var address common.Address = common.HexToAddress(withdrawObsData.WalletAddress)
	// 	me.noOfGoRoutines <- 1
	// 	fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// 	go me.saveWalletAddress(conf_api.ServiceGroupWithdraw, &address, me.ethAddressMap, me.Config.MainnetEthClientChainID, true, nil)

	// }

	// //5. Get Admin fee wallet balance at side chain
	// log.GetLogger(log.Name.Root).Infoln("Starting to get admin fee wallet balance ----")
	// me.noOfGoRoutines <- 1
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupAdminFeeWallet, &me.adminFeeWalletAddress, me.sideChainAssetNameToAddress, me.Config.EthClientChainID, false, nil)

	// //6. Get gas fee wallet balance at side chain
	// log.GetLogger(log.Name.Root).Infoln("Starting to get gas fee wallet balance ----")
	// me.noOfGoRoutines <- 1
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupGasFeeWallet, &me.gasFeeWalletAddress, me.eunAddressMap, me.Config.EthClientChainID, false, nil)

	// //7. Save sign server user wallet balance in mainnet
	// log.GetLogger(log.Name.Root).Infoln("Starting to get sign server user wallet balance ----")
	// var userWalletAddress common.Address = common.HexToAddress(me.Config.UserWalletOwnerWalletAddr)
	// me.noOfGoRoutines <- 1
	// fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupSignUserWallet, &userWalletAddress, me.eunAddressMap, me.Config.EthClientChainID, false, nil)

	// //8. Save sign server invoker wallet balance in mainnet
	// log.GetLogger(log.Name.Root).Infoln("Starting to get sign server invoker wallet balance ---- ")
	// log.GetLogger(log.Name.Root).Infoln("Count: ", len(me.Config.InvokerWalletAddrList))
	// for _, invokerAddr := range me.Config.InvokerWalletAddrList {
	// 	var invokerWalletAddress common.Address = common.HexToAddress(invokerAddr)
	// 	me.noOfGoRoutines <- 1
	// 	fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// 	go me.saveWalletAddress(conf_api.WalletBalanceGroupSignInvoker, &invokerWalletAddress, me.eunAddressMap, me.Config.EthClientChainID, false, nil)
	// }

	// //9. Save sweep server balance in mainnet
	// log.GetLogger(log.Name.Root).Infoln("Starting to get sweepServiceInvoker wallet balance ----")
	// me.noOfGoRoutines <- 1
	// fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// go me.saveWalletAddress(conf_api.WalletBalanceGroupSweepServerWallet, &me.sweepServerWalletAddress, me.ethAddressMap, me.Config.MainnetEthClientChainID, true, nil)

	// //10. Save market registration wallet balance in sidechain
	// log.GetLogger(log.Name.Root).Infoln("Starting to get marketing registration wallet balance ----")
	// me.noOfGoRoutines <- 1
	// fmt.Printf("Add go routines from (%d) ->(%d) \n", len(me.noOfGoRoutines)-1, len(me.noOfGoRoutines))
	// go me.saveWalletAddress(conf_api.WalletBalanceMarketRegWallet, &me.marketRegWalletAddress, me.eunAddressMap, me.Config.EthClientChainID, false, nil)

	//End. End the indexer if go routines did not run for (me.closeInMinutes) mins
	go me.closeAllIfGoRoutinesHasStoppedFor(me.closeInMinutes)
	return nil
}

func (me *WalletBackgroundIndexer) connectEthClientPool(pool chan *ethereum.EthClient, protocol, ip string, port int, chainId int64) error {

	for i := 0; i < poolSize; i++ {
		var connErr error
		var sideChainEthClient *ethereum.EthClient = &ethereum.EthClient{
			Protocol: protocol,
			IP:       ip,
			Port:     port,
			ChainID:  big.NewInt(chainId),
		}

		for j := 0; j < me.Config.RetryCount; j++ {
			_, connErr = sideChainEthClient.Connect()
			if connErr != nil {
				time.Sleep(time.Second * me.Config.GetRetryInterval())
			} else {
				break
			}
		}

		if connErr != nil {
			return connErr
		}
		pool <- sideChainEthClient
	}
	return nil
}

func (me *WalletBackgroundIndexer) prepareWalletConfigList(authServiceList []conf_api.AuthService, configList []wallet_bg_model.WalletBalanceConfig) []wallet_bg_model.WalletBalanceConfig {
	var authServiceGroupMap map[int][]conf_api.AuthService = make(map[int][]conf_api.AuthService)
	var authServiceMap map[uint64]conf_api.AuthService = make(map[uint64]conf_api.AuthService)

	for _, authService := range authServiceList {
		if authService.ServiceGroupId == 0 {
			continue
		}

		//Indexing by service group ID
		if _, ok := authServiceGroupMap[authService.ServiceGroupId]; !ok {
			authServiceGroupMap[authService.ServiceGroupId] = make([]conf_api.AuthService, 0)
		}
		authServiceGroupMap[authService.ServiceGroupId] = append(authServiceGroupMap[authService.ServiceGroupId], authService)

		//Indexing by service ID
		authServiceMap[uint64(authService.Id)] = authService
	}

	var outputConfigList []wallet_bg_model.WalletBalanceConfig = make([]wallet_bg_model.WalletBalanceConfig, 0)

	for _, config := range configList {
		switch config.ConfigType {
		case wallet_bg_model.WalletServiceGroupId:
			if serviceList, ok := authServiceGroupMap[config.ServiceGroupId]; !ok {
				log.GetLogger(log.Name.Root).Warnln("Unable to get service group id: ", config.ServiceGroupId, ". Skipped")
				continue
			} else {
				for _, service := range serviceList {
					outputConfig := wallet_bg_model.WalletBalanceConfig{}
					outputConfig.Id = config.Id
					outputConfig.ServiceGroupId = config.ServiceGroupId
					outputConfig.ServiceId = uint64(service.Id)
					if service.WalletAddress == "" {
						log.GetLogger(log.Name.Root).Warnln("Server id: ", service.Id, " empty wallet address. Skipped")
						continue
					}
					outputConfig.WalletAddress = service.WalletAddress
					outputConfig.ConfigType = wallet_bg_model.WalletConfigResolved
					outputConfig.WalletType = service.ServiceGroupId
					outputConfig.ChainId = config.ChainId
					outputConfig.AssetName = config.AssetName
					outputConfig.Description = config.Description

					outputConfigList = append(outputConfigList, outputConfig)
				}
			}
		case wallet_bg_model.WalletServiceId:
			if service, ok := authServiceMap[config.ServiceId]; !ok {
				log.GetLogger(log.Name.Root).Warnln("Unable to get service id: ", config.ServiceId, ". Skipped")
				continue
			} else {
				outputConfig := wallet_bg_model.WalletBalanceConfig{}
				outputConfig.Id = config.Id
				outputConfig.ServiceGroupId = config.ServiceGroupId
				outputConfig.ServiceId = uint64(service.Id)
				if service.WalletAddress == "" {
					log.GetLogger(log.Name.Root).Warnln("Server Id: ", service.Id, " empty wallet address. Skipped")
					continue
				}
				outputConfig.WalletAddress = service.WalletAddress
				outputConfig.ConfigType = wallet_bg_model.WalletConfigResolved
				outputConfig.WalletType = service.ServiceGroupId
				outputConfig.ChainId = config.ChainId
				outputConfig.AssetName = config.AssetName
				outputConfig.Description = config.Description

				outputConfigList = append(outputConfigList, outputConfig)
			}
		case wallet_bg_model.WalletAddress:
			config.ConfigType = wallet_bg_model.WalletConfigResolved
			outputConfigList = append(outputConfigList, config)
		default:
			log.GetLogger(log.Name.Root).Warnln("Invalid config type: ", config.ConfigType, ". Skipped")
			continue
		}
	}

	return outputConfigList

}

func (me *WalletBackgroundIndexer) saveWalletAddress(ethClientPool chan *ethereum.EthClient, walletType conf_api.ServiceGroupId, walletAddress *common.Address, addressMap AssetNameMap, chainId int, isMainnet bool, userId *uint64) {

	ethClient := <-ethClientPool

	defer func() {
		ethClientPool <- ethClient
		<-me.noOfGoRoutines
		fmt.Printf("\nGo rountines removed, active no. of  go rountines : (%d) \n", len(me.noOfGoRoutines))
	}()

	var listToBeInserted []wallet_bg_model.WalletBalance
	//fmt.Printf("Address of slice %p add of Arr \n", &listToBeInserted)
	walletAddressInHex := hex.EncodeToString(walletAddress.Bytes())
	var netStr = ""
	var userIdStr = ""
	if isMainnet {
		netStr = "Main"
	} else {
		netStr = "SideChain"
	}
	if userId != nil {
		userIdStr = strconv.FormatUint(*userId, 10)
	} else {
		userIdStr = ""
	}
	netStr += ", Chain Id :[" + strconv.Itoa(chainId) + "], WalletType:[" + strconv.FormatInt(int64(walletType), 10) + "]"
	netStr += ", WalletAddr :[" + walletAddressInHex + "], UserId:[" + userIdStr + "]"

	for assetName, _ := range addressMap {
		//fmt.Printf("%s  asset [%s] \n", netStr, assetName)
		nonCachedAssetName := assetName
		netStr += " Asset :[" + nonCachedAssetName + "]"
		balanceInBigInt, err := me.GetBalance(ethClient, isMainnet, *walletAddress, addressMap, nonCachedAssetName)
		if err != nil {
			netStr = "Fail to get balance: " + netStr
			fmt.Println(netStr)
			log.GetLogger(log.Name.Root).Errorln(netStr, err.Error())
		} else {
			var newRecord = wallet_bg_model.NewWalletBalance(strconv.FormatInt(int64(walletType), 10), walletAddressInHex, nonCachedAssetName, balanceInBigInt, userId, chainId)
			listToBeInserted = append(listToBeInserted, *newRecord)
		}
	}

	fmt.Printf("\n-----------------------------------------------------\nERC asset balance read as below:\nChain Id       : %d \nUser Id        : %s \nWallet address : %s \nAsset fetched  : %d \nSaving (%d) records to db....\n-----------------------------------------------------\n", chainId, userIdStr, walletAddressInHex, len(addressMap), len(listToBeInserted))
	err := BatchInsertWalletBalance(me, &listToBeInserted)
	if err != nil {
		logMsg := "Balance save fail for address (" + walletAddressInHex + ")"
		fmt.Println(logMsg)
		log.GetLogger(log.Name.Root).Errorln(logMsg, err.Error())
		me.noOfRecordsFail++
	} else {
		logMsg := "Balances save successfully for address (" + walletAddressInHex + ")"
		log.GetLogger(log.Name.Root).Info(logMsg)
		me.noOfRecordsSaved++
	}

}

func (me *WalletBackgroundIndexer) getTokenBalance(address common.Address, contractAddress common.Address, client bind.ContractBackend) (*big.Int, error) {
	inst, err := contract.NewEurusERC20(contractAddress, client)
	if err != nil {
		return nil, err
	}
	return inst.BalanceOf(&bind.CallOpts{}, address)
}

// Connect mainnet ethereum
func (me *WalletBackgroundIndexer) InitMainNetEthereumClientFromConfig() (*ethereum.EthClient, error) {

	config := me.Config
	err := config.ValidateEthClientField()
	if err != nil {
		return nil, err
	}
	ethClient := ethereum.EthClient{
		Protocol: me.Config.MainnetEthClientProtocol,
		IP:       me.Config.MainnetEthClientIP,
		Port:     me.Config.MainnetEthClientPort,
		ChainID:  big.NewInt(int64(me.Config.MainnetEthClientChainID)),
	}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.MainNetEthClient = &ethClient //Rinkyby Chain - 4
	me.MainNetEthClient.Logger = log.GetLogger(log.Name.Root)
	return me.MainNetEthClient, nil
}

// Connect side chain ethereum
func (me *WalletBackgroundIndexer) InitEurusEthereumClientFromConfig() (*ethereum.EthClient, error) {

	config := me.Config
	err := config.ValidateEthClientField()
	if err != nil {
		return nil, err
	}
	ethClient := ethereum.EthClient{
		Protocol: me.Config.EthClientProtocol,
		IP:       me.Config.EthClientIP,
		Port:     me.Config.EthClientPort,
		ChainID:  big.NewInt(int64(me.Config.EthClientChainID)),
	}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.EthClient = &ethClient //Side Chain - 2021
	me.EthClient.Logger = log.GetLogger(log.Name.Root)
	return me.EthClient, nil
}

func (me *WalletBackgroundIndexer) ProcessTotalSupply(assetMap *AssetNameMap, blockNumber *big.Int, markDate *time.Time, chainId int) {
	zeroAddress := common.Address{}
	for asset, address := range *assetMap {
		if bytes.Equal(address.Bytes(), zeroAddress.Bytes()) {
			continue
		}
		log.GetLogger(log.Name.Root).Debugln("Query asset: ", asset, " total supply. Address: ", address.Hex())
		erc20Contract, err := contract.NewEurusERC20(address, me.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Create EurusERC20 instance failed on asset: ", asset, " Error: ", err)
			continue
		}

		totalSupply, err := erc20Contract.TotalSupply(&bind.CallOpts{BlockNumber: blockNumber})
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Query total supply failed. Asset: ", asset, " Error: ", err)
			continue
		}

		assetTotalSupply := new(wallet_bg_model.AssetTotalSupply)
		if blockNumber == nil {
			blockNum, err := me.EthClient.Client.BlockNumber(context.Background())
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Get current block number failed. Asset: ", asset, " Error: ", err)
			} else {
				blockNumber = big.NewInt(0)
				blockNumber.SetUint64(blockNum)
			}

		}
		if blockNumber != nil {
			b := decimal.NewFromBigInt(blockNumber, 0)
			assetTotalSupply.BlockNumber = &b
		}
		assetTotalSupply.AssetName = asset
		assetTotalSupply.AssetAddress = ethereum.ToLowerAddressString(address.Hex())
		assetTotalSupply.ChainId = chainId
		total := decimal.NewFromBigInt(totalSupply, 0)
		assetTotalSupply.TotalSupply = &total
		assetTotalSupply.InitDate()
		assetTotalSupply.MarkDate = *markDate

		err = DbInsertAssetTotalSupply(me.DefaultDatabase, assetTotalSupply)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to insert database for asset: ", asset, " Error: ", err)
			continue
		}

	}
}

func contains(s []conf_api.ServiceGroupId, e conf_api.ServiceGroupId) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
