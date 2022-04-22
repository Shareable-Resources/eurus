package user

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/server"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type UserServerConfig struct {
	server.ServerConfigBase

	// PlatformWalletAddress string  `json:"platformWalletAddress"`
	MarketRegWalletAddress common.Address
	InitialFundAmount      float64 `json:"initialFundAmount"`
	EmailFrom              string  `json:"emailFrom"`
	// EmailPassword            string  `json:"emailPassword"`
	// SmtpHost                 string  `json:"smtpHost"`
	// SmtpPort                 string  `json:"smtpPort"`
	VerificationDuration int `json:"verificationDuration"` // second

	UserEthClientProtocol string `json:"userEthClientProtocol"`
	UserEthClientIP       string `json:"userEthClientIP"`
	UserEthClientPort     int    `json:"userEthClientPort"`

	MainnetEthClientProtocol string `json:"mainnetEthClientProtocol"`
	MainnetEthClientIP       string `json:"mainnetEthClientIP"`
	MainnetEthClientPort     int    `json:"mainnetEthClientPort"`
	MainnetEthClientChainID  int    `json:"mainnetEthClientChainID"`

	UserMainnetEthClientProtocol string `json:"userMainnetEthClientProtocol"`
	UserMainnetEthClientIP       string `json:"userMainnetEthClientIP"`
	UserMainnetEthClientPort     int    `json:"userMainnetEthClientPort"`

	SignServerUrl             string   `json:"signServerUrl"`
	KYCServerUrl              string   `json:"kycServerUrl"`
	InvokerAddressJson        string   `json:"invokerAddress"`
	InvokerAddressList        []string //sign server invoker addresses
	InitialFundExactAmount    *big.Int
	UserObserverList          []*conf_api.ServerDetail
	EmailServiceZone          string `json:"awsEmailServiceZone"`
	AwsAccessKeyId            string `json:"awsAccessKeyId"`
	AwsAccessSecretAccessKey  string `json:"awsAccessSecretAccessKey"`
	FaucetGasLimit            int64  `json:"faucetGasLimit"`
	UserMnenomicAesKey        string `json:"userMnenomicAesKey" eurus_conf:"noPrint"`
	SideChainGasLimit         int64  `json:"sideChainGasLimit"`
	MarketingRewardSchemeJson string `json:"marketingRewardScheme"`
	BSCAssetAddressList       string `json:"bscAssetAddressList"`

	BlockCypherDBServerIP     string `json:"blockCypherDBServerIP"`
	BlockCypherDBPort         int    `json:"blockCypherDBPort"`
	BlockCypherDBUserName     string `json:"blockCypherDBUserName"`
	BlockCypherDBPassword     string `json:"blockCypherDBPassword"`
	BlockCypherDBDatabaseName string `json:"blockCypherDBDatabaseName"`
	BlockCypherDBSchemaName   string `json:"blockCypherSchemaName"`
	BlockCypherChain          string `json:"blockCypherChain"`

	//Elastic search
	ElasticLoginDataFilePath string `json:"elasticLoginDataFilePath"`
}

func NewUserServerConfig() *UserServerConfig {
	config := new(UserServerConfig)
	config.InitialFundExactAmount = big.NewInt(0)
	config.VerificationDuration = 300
	config.ActualConfig = config
	return config
}

func (me *UserServerConfig) GetUserServerConfig() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *UserServerConfig) InitInitialFundExactAmount() {
	me.InitialFundExactAmount.SetUint64(uint64(math.Pow10(18) * me.InitialFundAmount))
}

func (me *UserServerConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *UserServerConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
