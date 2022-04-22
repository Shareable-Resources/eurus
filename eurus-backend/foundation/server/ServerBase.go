package server

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"math/big"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/uds"
	"eurus-backend/foundation/ws/ws_message"
)

type ServerBase struct {
	DefaultDatabase        *database.Database
	SlaveDatabase          *database.ReadOnlyDatabase
	EthClient              *ethereum.EthClient
	EthWebSocketClient     *ethereum.EthClient
	MainNetEthClient       *ethereum.EthClient
	MainNetWebSocketClient *ethereum.EthClient
	ServerConfig           *ServerConfigBase
	HttpServer             *network.HttpServer
	AuthClient             auth_base.IAuth
	ActualServer           interface{} ///Store the child class pointer
	ControlServer          *uds.UDSMessageServer
	ControlMessageChan     chan (*uds.MessageFrame)
}

// The InitLog create a log with name from constant log.Name.Root. The filePath should be in UserObserverConfig.json (logFilePath)
func (me *ServerBase) InitLog(filePath string) {
	_, err := log.NewLogger(log.Name.Root, filePath, logrus.DebugLevel)
	if err != nil {
		log.GetDefaultLogger().Error("Unable to create ROOT log ", err.Error())
	}
}

func (me *ServerBase) InitDBFromConfig(config *ServerConfigBase) error {
	_, err := me.InitDB(config.DBUserName, config.DBPassword, config.DBServerIP, config.DBServerPort,
		config.DBDatabaseName, config.DBSchemaName, config.DBAESKey, config.DBIdleConns, config.DBMaxOpenConns)
	if err != nil {
		return errors.Wrap(err, "Init Default DB error")
	}

	if config.DBSlaveServerIP != "" && config.DBSlaveServerPort != 0 {
		_, err := me.InitSlaveDB(config.DBUserName, config.DBPassword, config.DBSlaveServerIP, config.DBSlaveServerPort,
			config.DBDatabaseName, config.DBSchemaName, config.DBAESKey, config.DBIdleConns, config.DBMaxOpenConns)
		if err != nil {
			return errors.Wrap(err, "Init Slave DB error")
		}
	}
	return nil
}

func (me *ServerBase) InitDB(userName string, password string, ip string, port int, dbName string, schemaName string,
	key string, idleConnCount int, maxConnCount int) (*database.Database, error) {
	defaultDatabase := &database.Database{ReadOnlyDatabase: database.ReadOnlyDatabase{IP: ip, Port: port, UserName: userName, Password: password, DBName: dbName, SchemaName: schemaName, IdleConns: idleConnCount, MaxOpenConns: maxConnCount}}
	defaultDatabase.Validate()
	defaultDatabase.SetAESKey(key)

	//Try connect DB
	_, err := defaultDatabase.GetConn()
	if err != nil {
		return nil, err
	}
	me.DefaultDatabase = defaultDatabase
	return defaultDatabase, nil
}

func (me *ServerBase) InitSlaveDB(userName string, password string, ip string, port int, dbName string, schemaName string,
	key string, idleConnCount int, maxConnCount int) (*database.ReadOnlyDatabase, error) {

	dbm := &database.ReadOnlyDatabase{IP: ip, Port: port, UserName: userName, Password: password, DBName: dbName, SchemaName: schemaName, IdleConns: idleConnCount, MaxOpenConns: maxConnCount}
	dbm.Validate()
	dbm.SetAESKey(key)

	//Try connect DB
	_, err := dbm.GetConn()
	if err != nil {
		return nil, err
	}
	me.SlaveDatabase = dbm
	return me.SlaveDatabase, nil
}

func (me *ServerBase) LoadConfig(args *CommandLineArguments, config IServerConfig, fileLoader func(configFilePath string) ([]byte, error)) error {

	me.PrintVersion()

	configBase := config.GetServerConfigBase()
	if configBase == nil {
		fmt.Println("IServerConfig returns null ServerConfigBase")
		log.GetLogger(log.Name.Root).Errorln("IServerConfig returns null ServerConfigBase")
		return errors.New("IServerConfig returns null ServerConfigBase")
	}

	pathByte := []byte(args.ConfigFilePath)
	if args.ConfigFilePath == "" || (len(pathByte) == 2 && pathByte[0] == 34 && pathByte[1] == 34) {
		currExecPath, err := os.Getwd()
		if err != nil {
			logger := log.GetLogger(log.Name.Root)
			fmt.Println("Getting working directory error: ", err)
			logger.Error(err.Error())
			return err
		}
		args.ConfigFilePath = path.Join(currExecPath, "config", "ServerConfig.json")
	}

	fmt.Println("Waiting for config file at path: ", args.ConfigFilePath)
	for {
		if _, err := os.Stat(args.ConfigFilePath); err == nil {
			fmt.Println("Config file found, loading...")
			break
		} else {
			time.Sleep(time.Second)
		}
	}
	var configByte []byte
	var loadErr error
	if fileLoader == nil {
		configByte, loadErr = ioutil.ReadFile(args.ConfigFilePath)
		if loadErr != nil {
			fmt.Println("[ServerBase] Load file failed: ", loadErr)
			logger := log.GetLogger(log.Name.Root)
			logger.Error(loadErr.Error())
			return loadErr
		}
	} else {
		configByte, loadErr = fileLoader(args.ConfigFilePath)
		if loadErr != nil {
			logger := log.GetLogger(log.Name.Root)
			fmt.Println("[ServerBase] fileLoader returns error: ", loadErr)
			logger.Error(loadErr.Error())
			return loadErr
		}
	}

	err := json.Unmarshal(configByte, config)
	if err != nil {
		logger := log.GetLogger(log.Name.Root)
		fmt.Println("[ServerBase] config file unmarshal error: ", err)
		logger.Error("Log config error: ", err.Error())
		return err
	}
	config.ValidateField()
	if configBase.RetryCount <= 0 {
		configBase.RetryCount = 20
	}

	if configBase.RetryInterval <= 0 {
		configBase.RetryInterval = 1
	}

	configBase.HdWalletAddress = strings.ToLower(configBase.HdWalletAddress)

	config.SetHttpErrorLogger(log.GetLogger(log.Name.Root))

	fmt.Println("Service ID: ", config.GetServerConfigBase().ServiceId)

	return nil
}

func (me *ServerBase) WriteConfig(configFilePath string, newServerConfig *ServerConfigBase) error {
	pathByte := []byte(configFilePath)
	if configFilePath == "" || (len(pathByte) == 2 && pathByte[0] == 34 && pathByte[1] == 34) {
		currExecPath, err := os.Getwd()
		if err != nil {
			logger := log.GetLogger(log.Name.Root)
			logger.Error(err.Error())
			return err
		}
		configFilePath = path.Join(currExecPath, "config", "ServerConfig.json")
	}
	file, _ := json.MarshalIndent(newServerConfig, "", " ")
	_ = ioutil.WriteFile(configFilePath, file, 0644)
	return nil
}

func (me *ServerBase) InitHttpServer(httpConfig network.IHttpConfig) error {
	if httpConfig == nil {
		httpConfig = me.ServerConfig
	}
	var err error
	me.HttpServer, err = network.NewServer(httpConfig)
	if err == nil {
		api.SetApiLogger(me.HttpServer.Logger)
	}
	return err
}

func (me *ServerBase) InitAuth(authClient auth_base.IAuth, config auth_base.IAuthBaseConfig, loginHandler func(auth_base.IAuth)) {
	me.AuthClient = authClient
	me.AuthClient.SetLoginHandler(loginHandler)
	me.AuthClient.LoginAuthServer(config)
}

func (me *ServerBase) InitEthereumClient(protocol string, ip string, port int, chainID int64) (*ethereum.EthClient, error) {
	ethClient := ethereum.EthClient{Protocol: protocol, IP: ip, Port: port, ChainID: big.NewInt(chainID)}
	_, err := ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.EthClient = &ethClient
	return me.EthClient, nil
}

func (me *ServerBase) InitEthereumWebSocketClient(protocol string, ip string, port int, chainID int64) (*ethereum.EthClient, error) {
	ethClient := ethereum.EthClient{Protocol: protocol, IP: ip, Port: port, ChainID: big.NewInt(chainID)}
	_, err := ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.EthWebSocketClient = &ethClient
	return me.EthClient, nil
}

func (me *ServerBase) InitEthereumClientFromConfig(config *ServerConfigBase) (*ethereum.EthClient, error) {
	err := config.ValidateEthClientField()
	if err != nil {
		return nil, err
	}
	chainID := big.NewInt(int64(config.EthClientChainID))
	ethClient := ethereum.EthClient{Protocol: config.EthClientProtocol, IP: config.EthClientIP, Port: config.EthClientPort, ChainID: chainID}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.EthClient = &ethClient
	return me.EthClient, nil
}

func (me *ServerBase) InitMainNetEthereumClientFromConfig(config *ServerConfigBase) (*ethereum.EthClient, error) {
	err := config.ValidateEthClientField()
	if err != nil {
		return nil, err
	}
	chainID := big.NewInt(int64(config.EthClientChainID))
	ethClient := ethereum.EthClient{Protocol: config.EthClientProtocol, IP: config.EthClientIP, Port: config.EthClientPort, ChainID: chainID}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.MainNetEthClient = &ethClient
	me.MainNetEthClient.Logger = log.GetLogger(log.Name.Root)
	return me.MainNetEthClient, nil
}

func (me *ServerBase) InitEthereumWebSocketClientFromConfig(config *ServerConfigBase) (*ethereum.EthClient, error) {
	err := config.ValidateEthClientWebSocketField()
	if err != nil {
		return nil, err
	}
	chainID := big.NewInt(int64(config.EthClientChainID))
	ethClient := ethereum.EthClient{Protocol: config.EthClientWebSocketProtocol, IP: config.EthClientWebSocketIP, Port: config.EthClientWebSocketPort, ChainID: chainID}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.EthWebSocketClient = &ethClient
	me.EthWebSocketClient.Logger = log.GetLogger(log.Name.Root)
	return me.EthWebSocketClient, nil
}

func (me *ServerBase) InitMainNetEthereumWebSocketClientFromConfig(config *ServerConfigBase) (*ethereum.EthClient, error) {
	err := config.ValidateEthClientWebSocketField()
	if err != nil {
		return nil, err
	}
	chainID := big.NewInt(int64(config.EthClientChainID))
	ethClient := ethereum.EthClient{Protocol: config.EthClientWebSocketProtocol, IP: config.EthClientWebSocketIP, Port: config.EthClientWebSocketPort, ChainID: chainID}
	_, err = ethClient.Connect()
	if err != nil {
		return nil, err
	}
	me.MainNetWebSocketClient = &ethClient
	return me.MainNetWebSocketClient, nil
}

func (me *ServerBase) PostRequest(req api.RequestResponse) {}

func (me *ServerBase) Shutdown() {}

func (me *ServerBase) InitUDSControlServer(commandLineArgs *CommandLineArguments, handler func(req *ControlRequestMessage) (bool, string, error)) {
	if commandLineArgs.IsNoUDSControl {
		return
	}
	if commandLineArgs.UDSPath == "" {
		currExecPath, err := os.Getwd()
		if err != nil {
			logger := log.GetLogger(log.Name.Root)
			logger.Error("Getting working directory error: ", err)
			return
		}
		commandLineArgs.UDSPath = path.Join(currExecPath, "/sock/service_"+strconv.FormatInt(me.ServerConfig.ServiceId, 10))
	}
	me.ControlMessageChan = make(chan (*uds.MessageFrame))
	log.GetLogger(log.Name.Root).Infoln("Control socket path: ", commandLineArgs.UDSPath)
	fmt.Println("Control socket path: ", commandLineArgs.UDSPath)
	var err error
	me.ControlServer, err = uds.NewUDSMessageServer(commandLineArgs.UDSPath, log.GetLogger(log.Name.Root), me.ControlMessageChan)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Control socket create failed ", err)
		return
	}

	go me.ProcessControlMessage(handler)

	me.ControlServer.Listen()

}

func (me *ServerBase) ProcessControlMessage(handler func(*ControlRequestMessage) (bool, string, error)) {
	for {
		messageFrame := <-me.ControlMessageChan
		if messageFrame.Error != nil {
			log.GetLogger(log.Name.Root).Errorln("Control server received error: ", messageFrame.Error)
			continue
		}
		if messageFrame.IsConnectionEstablishedEvent {
			res := ws_message.CreateSuccessResponseMessage("welcomeMessage", "Service ID: "+strconv.FormatInt(me.ServerConfig.ServiceId, 10), "")
			me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
			continue
		}

		if messageFrame.MasterRequestMessage != nil {
			req := new(ControlRequestMessage)
			err := req.UnmarshalJSON([]byte(messageFrame.MasterRequestMessage.Message))
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Control server message error: ", err, " message: ", messageFrame.MasterRequestMessage.Message)
				res := ws_message.CreateErrorResponseMessage(foundation.InvalidArgument, err.Error(), messageFrame.MasterRequestMessage.Nonce)
				me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
				continue
			}

			log.GetLogger(log.Name.Root).Debugln("Control server message received: ", messageFrame.MasterRequestMessage.Message)
			fmt.Println("Command received: ", messageFrame.MasterRequestMessage.Message)
			var output string
			var isHandled bool
			if handler != nil {
				isHandled, output, err = handler(req)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Control server process command with error: ", err)
					res := ws_message.CreateErrorResponseMessage(foundation.InvalidArgument, err.Error(), messageFrame.MasterRequestMessage.Nonce)
					me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
					continue
				}

				if isHandled {
					res := ws_message.CreateSuccessResponseMessage(req.MethodName, output, req.Nonce)
					me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
					continue
				}
			}
			_, output, err = me.TerminalFunction(req, output)
			if err != nil {
				res := ws_message.CreateErrorResponseMessage(foundation.InvalidArgument, err.Error(), messageFrame.MasterRequestMessage.Nonce)
				me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
			} else {
				res := ws_message.CreateSuccessResponseMessage(req.MethodName, output, req.Nonce)
				me.ControlServer.SendJsonMessage(messageFrame.Conn, res)
			}
		}
	}
}

func (me *ServerBase) TerminalFunction(req *ControlRequestMessage, preMessage string) (bool, string, error) {
	var output string = preMessage
	var isHandled bool = true
	var err error

	switch req.MethodName {
	case "date":
		temp := time.Now()
		output += fmt.Sprintln("RFC3339: ", temp.Format(time.RFC3339))
		output += fmt.Sprintln("UnixTimeStamp: ", temp.Unix())
	case "help":
		output += fmt.Sprintln("date - Display the server time")
		output += fmt.Sprintln("enableAccessLog - enable writing to access log file")
		output += fmt.Sprintln("disableAccessLog - disable writing to access log file")
		output += fmt.Sprintln("printConfig - Print config (without sensitive information)")
		output += fmt.Sprintln("version - Print version")
	case "enableaccesslog":
		if me.HttpServer != nil {
			me.HttpServer.EnableAccessLog()
			output = "Access log enabled"
		} else {
			output = "HTTP server is not enabled"
		}
	case "disableaccesslog":
		if me.HttpServer != nil {
			me.HttpServer.DisableAccessLog()
			output = "Access log disabled"
		} else {
			output = "HTTP server is not enabled"
		}
	case "printconfig":
		output, err = me.PrintConfig()
		if err != nil {
			return true, output, err
		}
	case "version":
		output = me.PrintVersion()
	default:
		isHandled = false
		output = fmt.Sprintln("Command not found. Type \"help\" to see command list")
	}

	return isHandled, output, nil
}

func (me *ServerBase) SendConfigApiRequest(reqRes *api.RequestResponse) (*api.RequestResponse, error) {

	configUrl := url.URL{
		Scheme: "http",
		Host:   me.ServerConfig.ConfigServerIP + ":" + strconv.Itoa(me.ServerConfig.ConfigServerPort),
		Path:   reqRes.Req.GetRequestPath(),
	}
	if reqRes.RetrySetting == nil {
		reqRes.RetrySetting = me.ServerConfig
	}
	return api.SendApiRequest(configUrl, reqRes, me.AuthClient)
}

func (me *ServerBase) SendApiRequest(urlStr string, reqRes *api.RequestResponse) (*api.RequestResponse, error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return api.SendApiRequest(*urlObj, reqRes, me.AuthClient)
}

func (me *ServerBase) HttpServerListen() error {
	//do something
	return me.HttpServer.Listen()
}

func (me *ServerBase) PrintConfig() (string, error) {
	var output string
	var mapStruct map[string]interface{}
	deConf := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &mapStruct,
		Squash:  true,
	}
	if me.ServerConfig.ActualConfig == nil {
		output += fmt.Sprintln("Warning: ServerBase.ServerConfig.ActualConfig is nil, please set it at the ServerConfigBase child class")
	}
	var err error
	var serializedInterface interface{}
	decoder, _ := mapstructure.NewDecoder(deConf)
	if me.ServerConfig.ActualConfig != nil {
		serializedInterface = me.ServerConfig.ActualConfig

	} else {
		serializedInterface = me.ServerConfig
	}

	err = decoder.Decode(serializedInterface)
	if err != nil {
		return "", err
	}
	t := reflect.TypeOf(serializedInterface).Elem()
	me.trimNoPrintField(t, mapStruct)
	configData, err := json.MarshalIndent(mapStruct, "", "	")
	if err != nil {
		return "", err
	}
	output += fmt.Sprintln(string(configData))
	return output, nil
}

func (me *ServerBase) PrintVersion() string {
	var output string
	output += fmt.Sprintf("\033[1;33m%s %s\033[0m\r\n", "Build version:", foundation.GitCommit)
	output += fmt.Sprintf("\033[1;33m%s %s\033[0m\r\n", "Build date:", foundation.BuildDate)

	log.GetLogger(log.Name.Root).Infoln("Build version: ", foundation.GitCommit)
	log.GetLogger(log.Name.Root).Infoln("Build date: ", foundation.BuildDate)

	fmt.Println(output)
	return output
}

func (me *ServerBase) trimNoPrintField(t reflect.Type, mapStruct map[string]interface{}) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			me.trimNoPrintField(field.Type, mapStruct)
		}
		tagStr := field.Tag.Get("eurus_conf")
		if tagStr == "noPrint" {
			var fieldName string
			jsonFieldName := field.Tag.Get("json")
			if jsonFieldName != "" && jsonFieldName != "-" {
				fieldName = jsonFieldName
			} else {
				fieldName = field.Name
			}
			delete(mapStruct, fieldName)
		}
	}
}
