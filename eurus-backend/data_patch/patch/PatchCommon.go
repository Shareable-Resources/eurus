package patch

import (
	"encoding/json"
	"eurus-backend/env"
	"eurus-backend/foundation"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

var eurusAddressMap map[string]common.Address = make(map[string]common.Address)
var mainnetAddressMap map[string]common.Address = make(map[string]common.Address)

func InitLog() {
	programName := os.Args[0]

	dir, err := os.Getwd()
	if err != nil {
		panic("initLog failed: " + err.Error())
	}
	_, programName = path.Split(programName)
	var logPath string = path.Join(dir, "log", programName+"_log.log")
	logger, err := log.NewLogger(log.Name.Root, logPath, logrus.DebugLevel)
	if err != nil {
		panic("Unable to create ROOT log " + err.Error())
	}
	mw := io.MultiWriter(os.Stdout, logger.Out)
	logger.SetOutput(mw)

	logger.Infoln("Environment: ", env.Tag)
	logger.Infoln("Version: ", foundation.GitCommit)
	logger.Infoln("Build date: ", foundation.BuildDate)
}

func LoadConfig(configFilePath string, config interface{}) error {
	configByte, loadErr := ioutil.ReadFile(configFilePath)
	if loadErr != nil {
		logger := log.GetLogger(log.Name.Root)
		logger.Error(loadErr.Error())
		return errors.Wrap(loadErr, "Read config file error")
	}

	err := json.Unmarshal(configByte, config)
	if err != nil {
		logger := log.GetLogger(log.Name.Root)
		logger.Error("Log config error: ", err.Error())
		return err
	}
	return nil
}

func LoadSmartContractConfig(scConfigFileNameSuffix string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get current working dir: ", err)
		return errors.Wrap(err, "Unable to get current working dir")
	}

	var suffix string
	if scConfigFileNameSuffix != "" {
		suffix = "_" + scConfigFileNameSuffix
	}
	configPath := path.Join(dir, "SmartContractDeploy"+suffix+".json")
	file, err := os.OpenFile(configPath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to load smart contract address JSON: ", err.Error())
		return errors.Wrap(err, "Unable to load smart contract address JSON")
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to read file: ", err.Error())
		return errors.Wrap(err, "Unable to read file")
	}
	var jsonMap map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal(rawData, &jsonMap)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid JSON format: ", err.Error())
		return errors.Wrap(err, "Invalid JSON format")
	}

	if smartContractObj, ok := jsonMap[strconv.FormatInt(env.DefaultEurusChainId, 10)]; ok {

		smartContractMap := smartContractObj.(map[string]interface{})

		smartContractInnerMap := smartContractMap["smartContract"].(map[string]interface{})
		for key, value := range smartContractInnerMap {

			child := value.(map[string]interface{})
			if intf, ok := child["address"]; ok {
				if addr, ok := intf.(string); ok {
					addrObj := common.HexToAddress(addr)
					eurusAddressMap[key] = addrObj
				} else {
					log.GetLogger(log.Name.Root).Errorf("%s address invalid\r\n", key)
					return errors.New(key + " address invalid")
				}
			} else {
				log.GetLogger(log.Name.Root).Errorf("%s address not found\r\n", key)
				return errors.New(key + " address not found")
			}
		}
	}

	if mainnetSmartContractObj, ok := jsonMap[strconv.FormatInt(env.DefaultMainnetChainId, 10)]; ok {

		mainnetSmartContractMap := mainnetSmartContractObj.(map[string]interface{})

		mainnetSmartContractInnerMap := mainnetSmartContractMap["smartContract"].(map[string]interface{})
		for key, value := range mainnetSmartContractInnerMap {

			child := value.(map[string]interface{})
			if intf, ok := child["address"]; ok {
				if addr, ok := intf.(string); ok {
					addrObj := common.HexToAddress(addr)
					mainnetAddressMap[key] = addrObj
				} else {
					log.GetLogger(log.Name.Root).Errorf("%s address invalid\r\n", key)
					return errors.New(key + " address invalid")
				}
			} else {
				log.GetLogger(log.Name.Root).Errorf("%s address not found\r\n", key)
				return errors.New(key + " address not found")
			}
		}
	}
	return nil
}

func CreateEurusEthClient(config *server.ServerConfigBase) (*ethereum.EthClient, error) {
	var ethClient ethereum.EthClient = ethereum.EthClient{
		Protocol: config.EthClientProtocol,
		IP:       config.EthClientIP,
		Port:     config.EthClientPort,
		ChainID:  big.NewInt(int64(config.EthClientChainID)),
	}

	_, err := ethClient.Connect()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect Eurus RPC: ", err.Error())
		return nil, err
	}

	return &ethClient, nil
}

func ReadTerminalHiddenInput(prompt string) ([]byte, error) {
	var data []byte
	var err error
	for {
		fmt.Print(prompt + " (screen is hidden): ")
		data, err = terminal.ReadPassword(0)
		if err != nil {
			if sysErr, ok := err.(syscall.Errno); ok {
				if sysErr == syscall.EOPNOTSUPP {
					fmt.Println("Debugger attached, using debug private key instead")
					data = []byte("64f8ba795cf8f78e9c3c7a1b154326ba6e0e6f994e4853f0a551c15519fb438e")
				}
			} else {
				return nil, err
			}
		}

		if len(data) > 0 {
			return data, nil
		}
	}
}

func GetAddressBySmartContractName(name string, chainId int64) common.Address {
	switch chainId {

	case env.DefaultEurusChainId:
		return eurusAddressMap[name]

	case env.DefaultMainnetChainId:
		return mainnetAddressMap[name]

	default:
		log.GetLogger(log.Name.Root).Errorln("Chain id: ", chainId, " not available")
		return common.Address{}
	}

}

func CreateTransOptNoSigner(ethClient *ethereum.EthClient, priKey string, config *PatchConfigBase) (*bind.TransactOpts, error) {
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(priKey, ethClient.ChainID)
	if err != nil {
		return nil, err
	}

	transOpt.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return tx, nil
	}
	transOpt.NoSend = true
	transOpt.GasLimit = config.GasLimit

	return transOpt, nil
}
