package background

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"
	"fmt"
	"os"
	// "strconv"
)

type BackgroundServer struct {
	service_server.ServiceServer
	Config          *BackgoundServerConfig
	loginMiddleware request.LoginMiddleware
	errorMiddleware response.ErrorMiddleware
	corsMiddleware  response.CORSMiddleware
	isSuccess       chan bool
}

type CurrencyIdSet map[string]string //Currency ID to currency ID mapping

func NewBackgroundServer() *BackgroundServer {
	backgroundServer := new(BackgroundServer)
	backgroundServer.Config = NewBackgroundServerConfig()
	backgroundServer.ServerConfig = &backgroundServer.Config.ServerConfigBase
	backgroundServer.isSuccess = make(chan bool)
	return backgroundServer
}

func (me *BackgroundServer) LoadConfig(args *service.ServiceCommandLineArgs) error {
	return me.ServiceServer.LoadConfigWithSetting(args, me.Config, false, false)
}

func (me *BackgroundServer) InitAll() {

	me.PrintVersion()

	go func() { me.ServiceServer.InitAuth(me.authLoginHandler, nil) }()

	for {
		select {
		case <-me.isSuccess:
			fmt.Println("Success")
			os.Exit(0)
		}
	}

}

func (me *BackgroundServer) authLoginHandler(authClient auth_base.IAuth) {
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

func (me *BackgroundServer) processInit() error {

	_, err := me.QueryConfigServer(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get config from config server: ", err)
		panic(err)
	}

	fmt.Println("Going to query config server coingeck id list")
	assetInfoList, err := QueryConfigServerCoingeckoId(me)
	if err != nil {
		panic(err)
	}

	var currencyIdSet CurrencyIdSet = make(CurrencyIdSet) //Currency ID to currency ID mapping
	for _, assetInfo := range assetInfoList {
		currencyIdSet[assetInfo.CurrencyId] = assetInfo.CurrencyId
	}

	var ids []string = make([]string, 0)
	for currencyId, _ := range currencyIdSet {
		ids = append(ids, currencyId)
	}
	url := ConstructQueryCoingeckoUrl(ids, []string{"eth"})
	fmt.Println("Going to query coingecko exchange rate")
	rate, err := QueryCoingeckoExchangeRate(url)

	if err != nil {
		panic(err)
	}

	fmt.Println("Request config server to update exchange rate")
	CallConfigServerUpdateExchangeRate(me, rate, assetInfoList)
	me.isSuccess <- true

	return nil
}
