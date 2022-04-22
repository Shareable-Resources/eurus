package user

import (
	"encoding/json"
	"errors"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/elastic"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/marketing/reward"
	"eurus-backend/secret"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	kyc_const "eurus-backend/user_service/kyc_service/kyc/const"
	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
)

type UserServer struct {
	service_server.ServiceServer
	Config          *UserServerConfig
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	corsMiddleware  response.CORSMiddleware
	store           UserStore
	kycMq           *network.MQPublisher
	rewardProcessor *reward.RewardProcessor
	blockCypherDb   *database.Database
	elasticLogger   elastic.ElasticSearch
}

func NewUserServer() *UserServer {
	userServer := new(UserServer)
	userServer.Config = NewUserServerConfig()
	userServer.ServerConfig = &userServer.Config.ServerConfigBase

	return userServer
}

//Load Json Config
func (me *UserServer) LoadConfig(args *service.ServiceCommandLineArgs) error {
	var err error = me.ServiceServer.LoadConfig(args, me.Config)
	me.Config.InitInitialFundExactAmount()
	return err
}

func (me *UserServer) InitHttpServer(httpConfig network.IHttpConfig) error {
	//Calling base class
	err := me.ServerBase.InitHttpServer(httpConfig)
	if err != nil {
		return err
	}
	me.loginMiddleware.AuthClient = me.AuthClient
	me.setupRouter()
	return err
}

func (me *UserServer) InitAll() {
	rand.Seed(time.Now().UnixNano())
	me.ServiceServer.InitAuth(me.authLoginHandler, me.configEventReceived)
}

//Store is the data from other server, which would be used through out user server scope
func (me *UserServer) InitStore() {
	go me.InitKYCStore()
}

func (me *UserServer) InitKYCStore() {
	for {
		var err error
		me.store.KYCCountryCodeList, err = GetKYCCountryListFromKYCServer(me)
		if me.store.KYCCountryCodeList != nil {
			log.GetLogger(log.Name.Root).Errorln("KYC Country List is fetched from KYC Server", err)
			//fmt.Println("KYC Country List is fetched from KYC Server")
			break
		}
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Retry in 5s. Unable to get kyc country code list from kyc server: ", err)
			//fmt.Println("Retry in 5s. Unable to get kyc country code list from kyc server: ")
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

func (me *UserServer) authLoginHandler(authClient auth_base.IAuth) {

	_, err := me.QueryConfigServer(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}
	if me.Config.EmailServiceZone == "" {
		panic("Email service zone is empty")
	}
	if me.Config.SideChainGasLimit <= 0 {
		me.Config.SideChainGasLimit = 10000000
	}

	if me.Config.ElasticLoginDataFilePath == "" {
		log.GetLogger(log.Name.Root).Panicln("elasticLoginDataFile path is empty")
	}

	me.elasticLogger = elastic.NewElasticSearchByLogFile(me.Config.ElasticLoginDataFilePath)
	if !me.elasticLogger.IsValid() {
		log.GetLogger(log.Name.Root).Panicln("Unable to init elastic search by log file: ", me.elasticLogger.GetError())
	}
	me.InitStore()

	me.Config.UserObserverList, err = me.QueryUserObserverAddress()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("QueryUserObserverAddress failed: ", err)
		panic("QueryUserObserverAddress failed: " + err.Error())
	}

	if me.Config.InvokerAddressJson == "" {
		log.GetLogger(log.Name.Root).Errorln("Invoker address does not set")
		panic("InvokerAddress does not set")
	}

	me.Config.InvokerAddressList = make([]string, 0)
	err = json.Unmarshal([]byte(me.Config.InvokerAddressJson), &me.Config.InvokerAddressList)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid JSON format for invoker address list")
		panic("Invalid JSON format for invoker address list")
	}

	me.Config.InitInitialFundExactAmount()
	log.GetLogger(log.Name.Root).Debugln("Initial funding EUN amount: ", me.Config.InitialFundExactAmount.String())
	wallet, acc, walletAddr, err := secret.GenerateMintBurnKey(me.Config.MnemonicPhase, me.Config.PrivateKey, strconv.Itoa(int(me.Config.ServiceId)))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load Hdwallet error: ", err.Error())
		panic("Load Hdwallet error: " + err.Error())
	}

	me.Config.HdWalletPrivateKey, err = wallet.PrivateKeyHex(*acc)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Load wallet key error: ", err.Error())
		panic("Load Hdwallet key error: " + err.Error())
	}
	me.Config.HdWalletAddress = walletAddr
	fmt.Println("User server Wallet address: ", walletAddr)
	log.GetLogger(log.Name.Root).Infoln("User server Wallet address: ", walletAddr)

	err = me.InitDBFromConfig(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to DB: ", err)
		panic("Unable to connect to DB: " + err.Error())
	}

	log.GetLogger(log.Name.Root).Debugln("EthClient endpoint: ", me.ServerConfig.EthClientProtocol, "://", me.ServerConfig.EthClientIP, ":", me.ServerConfig.EthClientPort, " chain ID: ", me.ServerConfig.EthClientChainID)

	_, err = me.InitEthereumClientFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Init EthClient: ", err.Error())
		panic("Unable to Init EthClient: " + err.Error())
	}

	addr, err := GetWalletAddressMapAddress(me)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get WalletAddressMap: ", err.Error())
		panic("Unable to get WalletAddressMap: " + err.Error())
	}
	me.ServerConfig.WalletAddressAddress = addr.Hex()

	// addr, err = GetPlatformWalletAddress(me)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("Unable to get Platform Wallet Address: ", err.Error())
	// 	panic("Unable to get Platform Wallet Address: " + err.Error())
	// }
	// me.Config.PlatformWalletAddress = addr.Hex()

	addr, err = GetMarketRegWalletAddress(me)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get MarketingRegWallet Address: ", err.Error())
		panic("Unable to get MarketingRegWallet Address: " + err.Error())
	}
	me.Config.MarketRegWalletAddress = *addr

	me.kycMq = network.NewMQPublisher(me.ServerConfig.MQUrl, network.MQModeTaskQueue, kyc_const.TaskQueueMetaData.QueueName, log.GetLogger(log.Name.Root))

	err = me.kycMq.InitPublisher(false)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("InitPublisher failed: ", err.Error())
		panic("InitPublisher failed: " + err.Error())
	}

	me.rewardProcessor = reward.NewRewardProcessor(me.DefaultDatabase, me.SlaveDatabase, me.ServerConfig.HdWalletPrivateKey, log.GetLogger(log.Name.Root))
	err = me.rewardProcessor.Init(common.HexToAddress(me.Config.InternalSCConfigAddress), me.EthClient)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Init RewardProcessor failed: ", err.Error())
		panic("Init RewardProcessor failed: " + err.Error())
	}

	go func() {
		err = me.InitHttpServer(nil)
		if err != nil {
			panic(err)
		}
		err = me.HttpServer.Listen()
		if err != nil {
			panic(err)
		}

	}()
}

func (me *UserServer) InitDBFromConfig(config *UserServerConfig) error {
	err := me.ServiceServer.InitDBFromConfig(&config.ServerConfigBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("ServiceServer.InitDBFromConfig error: ", err)
		return err
	}

	me.blockCypherDb, err = DbConnectBlockCypher(me.Config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect block cypher database: ", err)
		return err
	}
	return nil
}

func (me *UserServer) configEventReceived(message *amqp.Delivery, topic string, contentType string, content []byte) {

}

func (me *UserServer) QueryUserObserverAddress() ([]*conf_api.ServerDetail, error) {

	req := conf_api.NewQueryServiceGroupDetailRequest(uint64(conf_api.ServiceGroupUserObserver))
	res := &conf_api.QueryServiceGroupDetailFullResponse{}
	var err error
	reqRes := api.NewRequestResponse(req, res)
	reqRes, err = me.SendConfigApiRequest(reqRes)
	if err != nil {
		return nil, err
	}

	if len(res.Data.ServerList) == 0 {
		return nil, errors.New("Unable to get UserObserver details")
	}

	for _, detail := range res.Data.ServerList {
		if detail.WalletAddress == "" {
			return nil, errors.New("One of the User Observer wallet address is empty")
		}
	}
	return res.Data.ServerList, err
}

func (me *UserServer) setupRouter() {
	fmt.Println("3. func - setupRouter - start")
	me.HttpServer.Router.Use(me.corsMiddleware.Handler)
	me.HttpServer.Router.Get(RootPath+GetClientVersionPath, me.getClientVersion)
	me.HttpServer.Router.Post(RootPath+ImportWalletPath, me.importWallet)
	me.HttpServer.Router.Get(RootPath+GetServerConfigPath, me.getServerConfig)
	me.HttpServer.Router.Post(RootPath+RegisterPath, me.register)
	me.HttpServer.Router.Post(RootPath+VerificationPath, me.verification)
	me.HttpServer.Router.Post(RootPath+LoginBySignaturePath, me.loginBySignature)
	me.HttpServer.Router.Post(RootPath+ResendVerificationEmailPath, me.resendVerificationEmail)
	me.HttpServer.Router.Post(RootPath+ForgetLoginPasswordPath, me.forgetLoginPassword)
	me.HttpServer.Router.Post(RootPath+VerifyForgetLoginPasswordPath, me.verifyForgetLoginPasswordCode)
	me.HttpServer.Router.Get(RootPath+GetUserWalletAddressPath, me.getUserWalletAddress)
	me.HttpServer.Router.Get(RootPath+AssetAddressListPath, me.getAssetAddressList)
	me.HttpServer.Router.Get(RootPath+MarketingBanner, me.getMarketingBanner)
	me.HttpServer.Router.Mount(RootPath, me.userRouter())

	fmt.Println()
	fmt.Println("Routes registered------")
	fmt.Println()
	//Print registered route
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(me.HttpServer.Router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}

func (me *UserServer) userRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	// /user
	r.Post(RefreshTokenPath, me.refreshLoginToken)
	r.Get(GetUserDetailsPath, me.getUserDetails)
	r.Get(GetWithdrawAdminFeePath, me.getWithdrawAdminFee)
	r.Post(GetRecentTransactionPath, me.getRecentTransaction)
	r.Post(SetupPaymentWalletPath, me.verificationCallback)
	r.Get(RequestLoginRequestTokenPath, me.requestLoginRequestToken)
	r.Post(RequestLoginTokenByLoginRequestTokenPath, me.requestLoginToken)
	r.Get(RequestPaymentLoginTokenPath, me.requestPaymentLoginToken)
	r.Post(RequestChangePaymentPasswordPath, me.requestChangePaymentPassword)
	r.Post(RequestChangeLoginPasswordPath, me.requestChangeLoginPassword)
	r.Get(UserPreferenceStoragePath, me.getUserPreferenceStorage)
	r.Post(UserPreferenceStoragePath, me.setUserPreferenceStorage)
	r.Post(ResetLoginPasswordPath, me.resetLoginPassword)
	r.Post(ForgetPaymentPasswordPath, me.forgetPaymentPassword)
	r.Post(VerifyForgetPaymentPasswordPath, me.verifyForgetPaymentPassword)
	r.Post(ResetPaymentPasswordPath, me.resetPaymentPassword)
	r.Post(FindEmailWalletAddressPath, me.findEmailWalletAddress)
	r.Post(RegisterDevicePath, me.registerDevice)
	r.Post(VerifyDevicePath, me.verifyDevice)
	r.Post(FaucetPath, me.faucet)
	r.Get(GetFaucetConfigPath, me.getFaucetConfig)
	r.Post(SignTransactionPath, me.signTransaction)
	r.Post(TopUpPaymentTransactionPath, me.topUpPaymentWallet)

	// /user/kyc
	r.Mount(KYCServerPath, me.kycRouter())
	// /user/merchant
	r.Mount(MerchantServerPath, me.merchantRouter())
	// /user/marketing
	r.Mount(MarketingPath, me.marketingRouter())
	// /user/blockCypher
	r.Mount(BlockCypherPath, me.blockCypherRouter())
	return r
}

func (me *UserServer) kycRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Get(GetKYCCountryListPath, me.getKYCCountryListFromKYCServer)
	r.Get(GetKYCStatusByTokenPath, me.getKYCStatusOfUser)
	r.Post(SubmitKYCApprovalPath, me.submitKYCApproval)
	r.Post(CreateKYCStatusPath, me.createKYCStatus)
	r.Post(SubmitKYCDocumentPath, me.submitKYCDocument)
	return r
}

func (me *UserServer) merchantRouter() chi.Router {
	r := chi.NewRouter()
	r.Post(RequestMerchantRefundPath, me.requestMerchantRefund)
	r.Get(MerchantRefundStatusPath, me.queryMerchantRefundStatus)
	return r
}

func (me *UserServer) marketingRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Get(RewardListPath, me.getRewardList)
	r.Get(RewardSchemePath, me.getRewardScheme)

	return r
}

func (me *UserServer) blockCypherRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(me.loginMiddleware.VerifyUserLoginToken)
	r.Use(me.errorMiddleware.ErrorHandler)
	r.Get(BlockCypherAccessTokenPath, me.getBlockCypherAccessToken)
	return r
}

func (me *UserServer) getClientVersion(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewQueryClientVersionRequest()
	clientVersion, err := QueryClientVersion(me)
	var res *response.ResponseBase
	if err != nil {
		res = response.CreateErrorResponse(reqObj, foundation.DatabaseError, err.Error())
	} else {
		res = response.CreateSuccessResponse(reqObj, clientVersion)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) importWallet(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewImportWalletRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	err := reqObj.CheckTimestamp()
	// var isValid bool
	if err != nil {
		res = response.CreateErrorResponse(reqObj, foundation.TimestampError, err.Error())
	} else {
		res, _ = ImportWallet(me, reqObj, req.RemoteAddr)
	}

	api.HttpWriteResponse(writer, reqObj, res)
	// if !res.IsInterfaceNil() || !isValid {
	// 	api.HttpWriteResponse(writer, reqObj, res)
	// } else if isValid {
	// 	if res.GetReturnCode() == int64(foundation.Success) {
	// 		user, err := DbGetUserByWalletAddress(reqObj.WalletAddress, me.DefaultDatabase)
	// 		if err != nil {
	// 			log.GetLogger(log.Name.Root).Errorln("Fail to get user address. Error: ", err)
	// 		} else {
	// 			importWalletData := NewLoginDataFromLoginBySignature(reqObj.LoginLogDetail, user.Id, req.RemoteAddr)
	// 			go func() {
	// 				elastic.InsertLoginLog(me.ServerConfig.ElasticSearchPath, importWalletData)
	// 			}()
	// 		}
	// 	}
	// 	api.HttpWriteResponse(writer, reqObj, res)
	// }

}

func (me *UserServer) getServerConfig(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewQueryServerConfigRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = GetServerConfig(reqObj, me)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) refreshLoginToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRefreshTokenRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	if res == nil {
		res = RefreshLoginToken(me.AuthClient, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) loginBySignature(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewUserLoginBySignatureRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = LoginBySignature(me, reqObj, req.RemoteAddr)

	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *UserServer) getUserDetails(writer http.ResponseWriter, req *http.Request) {
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	reqObj := NewQueryUserDetailsRequest(loginToken.GetToken())
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = GetUserDetails(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getRecentTransaction(writer http.ResponseWriter, req *http.Request) {
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	reqObj := NewQueryRecentTransactionRequest(loginToken.GetUserId())
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = GetRecentTransaction(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getWithdrawAdminFee(writer http.ResponseWriter, req *http.Request) {
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	reqObj := NewQueryWithdrawAdminFeeRequest(loginToken.GetToken())
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	symbol := chi.URLParam(req, "curruncySymbol")
	res = GetWithdrawAdminFee(me, reqObj, symbol)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) register(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRegistrationRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = Register(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) verification(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewVerificationRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = EmailVerification(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) verificationCallback(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewSetupPaymentWalletRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = SetupPaymentWallet(me, reqObj, req.RemoteAddr)

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) resendVerificationEmail(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewResendVerificationEmailRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = ResendVerificationEmail(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) requestLoginRequestToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestLoginRequestTokenRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = RequestLoginRequestToken(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) requestLoginToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestLoginTokenRequest()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = RequestLoginTokenFromLoginRequestToken(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) requestPaymentLoginToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestPaymentLoginTokenRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	me.AuthClient.VerifyLoginToken(me.AuthClient.GetLoginToken(), 0)
	res = RequestPaymentLoginToken(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) requestChangePaymentPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestChangePaymentPasswordRequest()
	res := api.RequestToModel(req, reqObj)
	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = ChangePaymentPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) requestChangeLoginPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestChangeLoginPasswordRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = ChangeLoginPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getUserPreferenceStorage(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewGetUserPreferenceStorageRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	me.AuthClient.VerifyLoginToken(me.AuthClient.GetLoginToken(), 0)
	res = RequestGetUserStorage(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) setUserPreferenceStorage(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewSetUserPreferenceStorageRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	me.AuthClient.VerifyLoginToken(me.AuthClient.GetLoginToken(), 0)
	res = RequestSetUserStorage(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) forgetLoginPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewForgetLoginPassword()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = ForgetLoginPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) verifyForgetLoginPasswordCode(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewVerifyForgetLoginPasswordReqeust()
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = VerifyForgetLoginPasswordCode(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) resetLoginPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewResetLoginPasswordReqeust()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = ResetLoginPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) forgetPaymentPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = ForgetPaymentPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) verifyForgetPaymentPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewVerifyForgetPaymentPasswordReqeust()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = VerifyForgetPaymentPasswordCode(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) resetPaymentPassword(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewResetPaymentPasswordReqeust()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = ResetPaymentPassword(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}
func (me *UserServer) findEmailWalletAddress(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewEmailWalletAddressRequest()
	res := api.RequestToModel(req, reqObj)
	res = FindEmailWalletAddress(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) registerDevice(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRegisterDeviceRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = RegisterDevice(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) verifyDevice(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewVerifyDeviceRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if loginToken == nil {
		res = response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = VerifyDevice(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) faucet(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	symbol := chi.URLParam(req, "curruncySymbol")
	res = Faucet(me, reqObj, symbol)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getFaucetConfig(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res = GetFaucetConfig(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) signTransaction(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewSignTransactionRequest()
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = SignTransaction(me, reqObj)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getKYCCountryListFromKYCServer(writer http.ResponseWriter, req *http.Request) {

	reqObj := new(request.RequestBase)
	res := response.CreateSuccessResponse(reqObj, me.store.KYCCountryCodeList)
	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *UserServer) createKYCStatus(writer http.ResponseWriter, req *http.Request) {

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	userIdStrJson := loginToken.GetUserId()
	userIdObj := new(UserLoginId)
	json.Unmarshal([]byte(userIdStrJson), userIdObj)
	userId := userIdObj.UserId

	reqObj := kyc_model.NewRequestCreateKYCStatus()
	reqObj.UserId = userId
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res, err := CreateKYCStatusFromKYCServer(me, reqObj)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to CreateKYCStatusFromKYCServer: " + err.Error())
		res = response.CreateErrorResponse(reqObj, foundation.DatabaseError, err.Error())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	//Convert model to jsonString
	jsonString, _ := json.Marshal(res.GetData())
	fmt.Println(string(jsonString))
	// convert json to struct
	parsedModel := kyc_model.ResponseCreateKYCStatus{}
	json.Unmarshal(jsonString, &parsedModel)
	res = response.CreateSuccessResponse(reqObj, parsedModel)
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) submitKYCApproval(writer http.ResponseWriter, req *http.Request) {

	reqObj := kyc_model.NewRequestSubmitKYCApproval()
	reqObj.Method = http.MethodPost
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res, err := SubmitKYCApproval(me, reqObj)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to submitKYCApproval: " + err.Error())
		res = response.CreateErrorResponse(reqObj, foundation.ServerReturnCode(res.GetReturnCode()), err.Error())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = response.CreateSuccessResponse(reqObj, res.GetData())
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getKYCStatusOfUser(writer http.ResponseWriter, req *http.Request) {

	loginToken := req.Context().Value("loginToken").(auth_base.ILoginToken)
	userIdStrJson := loginToken.GetUserId()
	userIdObj := new(UserLoginId)
	json.Unmarshal([]byte(userIdStrJson), userIdObj)
	userId := userIdObj.UserId

	reqObj := kyc_model.NewRequestGetKYCStatusOfUser()
	reqObj.RequestPath = kyc_const.RootPath + RootPath + kyc_const.EndPoint.GetKYCStatusList
	reqObj.UserId = userId
	res := api.RequestToModel(req, reqObj)
	if res != nil {
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}
	res, err := GetKYCStatusOfUser(me, reqObj)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to getKYCStatusOfUser: " + err.Error())
		res = response.CreateErrorResponse(reqObj, foundation.ServerReturnCode(res.GetReturnCode()), err.Error())
		api.HttpWriteResponse(writer, reqObj, res)
		return
	}

	res = response.CreateSuccessResponse(reqObj, res.GetData())
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) submitKYCDocument(writer http.ResponseWriter, req *http.Request) {
	log.GetLogger(log.Name.Root).Debug("submitKYCDocument called")
	reader, err := req.MultipartReader()
	if err != nil {
		emptyReqObj := new(request.RequestBase)
		res := response.CreateErrorResponse(emptyReqObj, foundation.BadRequest, "Not a valid multipart form data")
		api.HttpWriteResponse(writer, emptyReqObj, res)
		return
	}

	loginTokenObj, ok := req.Context().Value("loginToken").(auth_base.ILoginToken)
	if ok {
		reqObj, res := ProcessKYCDocument(me, reader, loginTokenObj)
		api.HttpWriteResponse(writer, reqObj, res)
	} else {
		emptyReqObj := new(request.RequestBase)
		res := response.CreateErrorResponse(emptyReqObj, foundation.LoginTokenInvalid, "Invalid Login token format")
		api.HttpWriteResponse(writer, emptyReqObj, res)
	}
}

func (me *UserServer) requestMerchantRefund(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewRequestMerchantRefundRequest()
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = ProcessRequestMerchantRefund(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) queryMerchantRefundStatus(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewQueryMerchantRefundStatusRequest()
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = ProcessQueryMerchantRefundStatus(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}
func (me *UserServer) getRewardList(writer http.ResponseWriter, req *http.Request) {

	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = ProcessQueryRewardList(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getRewardScheme(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(request.RequestBase)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = ProcessQueryRewardScheme(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getUserWalletAddress(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(GetWalletAddressRequest)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		res = processGetWalletAddressByLoginAddress(me.SlaveDatabase, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) topUpPaymentWallet(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(TopUpPaymentWalletRequest)
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = processTopUpPaymentWallet(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)

}

func (me *UserServer) getAssetAddressList(writer http.ResponseWriter, req *http.Request) {
	chainName := chi.URLParam(req, "chain")
	reqObj := new(request.RequestBase)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		res = processGetAssetAddressList(me.Config, chainName, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getBlockCypherAccessToken(writer http.ResponseWriter, req *http.Request) {
	reqObj := NewBlockCypherAccessTokenRequest()
	res := api.RequestToModel(req, reqObj)
	if res == nil {
		res = processGetBlockCypherAccessToken(me, reqObj)
	}
	api.HttpWriteResponse(writer, reqObj, res)
}

func (me *UserServer) getMarketingBanner(writer http.ResponseWriter, req *http.Request) {
	reqObj := new(MarketingBannerRequest)
	res := api.RequestToModelNoLoginToken(req, reqObj)
	if res == nil {
		res = ProcessQueryMarketingBanner(me, reqObj)
	}

	api.HttpWriteResponse(writer, reqObj, res)
}
