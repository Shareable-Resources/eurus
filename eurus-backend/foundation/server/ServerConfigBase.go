package server

import (
	"errors"
	"eurus-backend/foundation/log"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerConfigBase struct { //implements IServerConfig, IAuthBaseConfig , implements IHttpConfig, implements IMQConsumerConfig
	LogFilePath string `json:"logFilePath"` //Mandatory field
	ServiceId   int64  `json:"serviceId"`   //Mandatory field
	GroupId     int64  `json:"groupId"`

	//Config server
	ConfigServerIP   string `json:"configServerIP"`   //Mandatory field
	ConfigServerPort int    `json:"configServerPort"` //Mandatory field

	//Database
	DBServerIP        string `json:"dbServerIP"`
	DBServerPort      int    `json:"dbServerPort"`
	DBSlaveServerIP   string `json:"dbSlaveServerIP"`
	DBSlaveServerPort int    `json:"dbSlaveServerPort"`
	DBUserName        string `json:"dbUserName"`
	DBPassword        string `json:"dbPassword" eurus_conf:"noPrint"`
	DBDatabaseName    string `json:"dbDatabaseName"`
	DBSchemaName      string `json:"schemaName"`
	DBIdleConns       int    `json:"dbIdleConns"`
	DBMaxOpenConns    int    `json:"dbMaxOpenConns"`
	DBAESKey          string `json:"dbAesKey" eurus_conf:"noPrint"` //use it for decrypt the db server pw

	//Ethereum client
	EthClientProtocol          string `json:"ethClientProtocol"`
	EthClientIP                string `json:"ethClientIP"`
	EthClientPort              int    `json:"ethClientPort"`
	EthClientWebSocketProtocol string `json:"ethClientWebSocketProtocol"`
	EthClientWebSocketPort     int    `json:"ethClientWebSocketPort"`
	EthClientWebSocketIP       string `json:"ethClientWebSocketIP"`
	EthClientChainID           int    `json:"ethClientChainID"`
	MnemonicPhase              string `json:"mnemonicPhase" eurus_conf:"noPrint"` //Exists at config file only
	HdWalletAddress            string
	HdWalletPrivateKey         string `eurus_conf:"noPrint"`

	//Http server
	HttpServerIP            string         `json:"httpServerIP"`
	HttpServerPort          uint           `json:"httpServerPort"`
	HttpUseDefaultAccessLog bool           `json:"httpUseDefaultAccessLog"`
	HttpAccessLogger        *logrus.Logger `json:"-"`
	HttpErrorLogger         *logrus.Logger `json:"-"`
	IsEnableAccessLog       bool           `json:"isEnableAccessLog"` //Exists at config file only
	AccessLogFilePath       string         `json:"accessLogFilePath"`

	//ElasticSearch server
	ElasticSearchPath string `json:"elasticSearchPath"`

	//Auth client
	PrivateKey     string `json:"privateKey" eurus_conf:"noPrint"`
	AuthServerIp   string `json:"authServerIp"`   //Exists at config file only
	AuthServerPort uint   `json:"authServerPort"` //Exists at config file only
	AuthPath       string `json:"authPath"`       //Exists at config file only

	//Config Smart Contract for Config Server

	InternalSCConfigAddress    string `json:"internalSCConfigAddress"`
	ExternalSCConfigAddress    string `json:"externalSCConfigAddress"`
	EurusInternalConfigAddress string `json:"eurusInternalConfigAddress"`
	WalletAddressAddress       string `json:"walletAddressAddress"`

	MQUrl         string `json:"mqUrl"`
	RetryCount    int    `json:"retryCount"`
	RetryInterval int    `json:"retryInterval"`

	//Put your inherited config object in this field
	ActualConfig interface{} `json:"-"`
}

func (me *ServerConfigBase) ValidateEthClientWebSocketField() error {
	var err error = nil
	if me.EthClientWebSocketProtocol == "" {
		err = errors.New("EthClientWebSocketProtocol should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.EthClientWebSocketIP == "" {
		err = errors.New("EthClientWebSocketIP should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.EthClientWebSocketPort == 0 {
		err = errors.New("EthClientWebSocketPort should be provided for InitEthereumWebSocketClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	}
	return err
}

func (me *ServerConfigBase) ValidateEthClientField() error {
	var err error = nil
	if me.EthClientProtocol == "" {
		err = errors.New("EthClientProtocol should be provided for InitEthereumClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.EthClientIP == "" {
		err = errors.New("EthClientIP should be provided for InitEthereumClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.EthClientPort == 0 {
		err = errors.New("EthClientPort should be provided for InitEthereumClient()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	}
	return err
}

func (me *ServerConfigBase) ValidateDBField() error {
	var err error = nil
	if me.DBServerIP == "" {
		err = errors.New("DBServerIP should be provided for initDB()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.DBServerPort == 0 {
		err = errors.New("DBServerPort should be provided for initDB()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.DBUserName == "" {
		err = errors.New("DBUserName should be provided for initDB()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.DBPassword == "" {
		err = errors.New("DBPassword should be provided for initDB()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	} else if me.DBIdleConns == 0 {
		me.DBIdleConns = 10
	} else if me.DBMaxOpenConns == 0 {
		me.DBMaxOpenConns = 100
	} else if me.DBAESKey == "" {
		err = errors.New("DBAESKey should be provided for initDB()!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
	}
	return err
}

func (me *ServerConfigBase) ValidateField() {
	var err error
	if me.ConfigServerIP == "" {
		err = errors.New("ConfigServerIP is a mandatory field for the server config file!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	} else if me.ConfigServerPort == 0 {
		err = errors.New("ConfigServerPort is a mandatory field for the server config file!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	} else if me.ServiceId == 0 {
		err = errors.New("ServiceId is a mandatory field for the server config file!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)

	} else if me.DBAESKey == "" {
		err = errors.New("DBAESKey is a mandatory field for the server config file!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	} else if me.PrivateKey == "" {
		err = errors.New("PrivateKey is a mandatory field for the server config file!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	}
}

//Mandatory field. Empty string to bind all network interface
func (me *ServerConfigBase) GetServerIp() string {
	return me.HttpServerIP
}

//set 0 will use default port 80
func (me *ServerConfigBase) GetServerPort() uint {
	return me.HttpServerPort
}

//If true, GetAccessLog() will be ignored
func (me *ServerConfigBase) IsUseDefaultAccessLog() bool {
	return me.HttpUseDefaultAccessLog
}

func (me *ServerConfigBase) GetAccessLog() *logrus.Logger {
	return me.HttpAccessLogger
}

//Mandatory field
func (me *ServerConfigBase) GetErrorLog() *logrus.Logger {
	return me.HttpErrorLogger
}

func (me *ServerConfigBase) GetAuthIp() string {
	return me.AuthServerIp
}
func (me *ServerConfigBase) GetAuthPort() uint {
	return me.AuthServerPort
}
func (me *ServerConfigBase) GetServiceId() int64 {
	return me.ServiceId
}
func (me *ServerConfigBase) GetPrivateKey() string {
	return me.PrivateKey
}

func (me *ServerConfigBase) GetAuthPath() string {
	return me.AuthPath
}

func (me *ServerConfigBase) GetRetryCount() int {
	return me.RetryCount
}

func (me *ServerConfigBase) GetRetryInterval() time.Duration {
	return time.Duration(me.RetryInterval)
}

func (me *ServerConfigBase) SetHttpErrorLogger(logger *logrus.Logger) {
	me.HttpErrorLogger = logger
}

func (me *ServerConfigBase) GetAccessLogFilePath() string {
	return me.AccessLogFilePath
}

func (me *ServerConfigBase) IsEnableHttpAccessLog() bool {
	return me.IsEnableAccessLog
}

func (me *ServerConfigBase) GetMqUrl() string {
	return me.MQUrl
}

// func (me *ServerConfigBase) GetQueueName() string {
// 	var queueName string
// 	if me.GroupId == 0 {
// 		queueName = "service_queue_" + strconv.FormatInt(me.ServiceId, 10)
// 	} else {
// 		queueName = "service_group_queue_" + strconv.FormatInt(me.GroupId, 10)
// 	}
// 	return queueName
// }

func (me *ServerConfigBase) GetConfigFileOnlyFieldList() []string {
	return []string{"serviceId", "isEnableAccessLog", "authServerIp", "authServerPort", "authPath", "dbAesKey", "mnemonicPhase"}
}

func (me *ServerConfigBase) GetServerConfigBase() *ServerConfigBase {
	return me
}

func (me *ServerConfigBase) GetParent() interface{} {
	return nil
}
