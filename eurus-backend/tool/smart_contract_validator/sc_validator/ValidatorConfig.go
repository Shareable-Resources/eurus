package sc_validator

import (
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ValidatorConfig struct {
	SideChainEthIp       string `json:"sideChainEthIp"`
	SideChainEthPort     int    `json:"sideChainEthPort"`
	SideChainId          int    `json:"sideChainId"`
	SideChainEthProtocol string `json:"sideChainEthProtocol"`

	MainnetEthIp       string `json:"mainnetEthIp"`
	MainnetEthPort     int    `json:"mainnetEthPort"`
	MainnetChainId     int    `json:"mainnetChainId"`
	MainnetEthProtocol string `json:"mainnetEthProtocol"`

	WithdrawObserverList  *WalletInfo `json:"withdrawObserver"`
	DepositObserverList   *WalletInfo `json:"depositObserver"`
	ApprovalObserverList  *WalletInfo `json:"approvalObserver"`
	BlockChainIndexerList *WalletInfo `json:"blockChainIndexer"`
	SideChainOwner        *WalletInfo `json:"sideChainOwner"`
	MainNetOwner          *WalletInfo `json:"mainnetOwner"`
	CurrencyList          []string    `json:"currencyList"`

	SmartContractDictionary map[string]*SmartContractDictionary

	addressRepository map[common.Address]string
}

func (me *ValidatorConfig) LoadSmartContractConfig(data []byte) error {
	var smartContractRepository map[string]*SmartContractDictionary = make(map[string]*SmartContractDictionary)

	err := json.Unmarshal(data, &smartContractRepository)
	if err != nil {
		return err
	}
	me.SmartContractDictionary = smartContractRepository

	for chainId, scDict := range me.SmartContractDictionary {
		for scName, info := range scDict.SmartContractMap {
			me.addressRepository[common.HexToAddress(info.Address)] = scName + " chain ID: " + chainId
		}
	}

	return nil
}

func (me *ValidatorConfig) Load() error {

	me.addressRepository = make(map[common.Address]string)
	if me.ApprovalObserverList == nil {
		return errors.New("Missing approvalObserver config")
	}
	for i, approvalAddrStr := range me.ApprovalObserverList.AddressList {
		me.addressRepository[common.HexToAddress(approvalAddrStr)] = "Approval Observer " + strconv.Itoa(i)
	}

	if me.WithdrawObserverList == nil {
		return errors.New("Missing withdrawObserver config")
	}
	for i, withdrawAddrStr := range me.WithdrawObserverList.AddressList {
		me.addressRepository[common.HexToAddress(withdrawAddrStr)] = "Withdraw Observer " + strconv.Itoa(i)
	}

	if me.DepositObserverList == nil {
		return errors.New("Missing depositObserver config")
	}
	for i, depositObserverAddr := range me.DepositObserverList.AddressList {
		me.addressRepository[common.HexToAddress(depositObserverAddr)] = "Deposit Observer " + strconv.Itoa(i)
	}

	if me.BlockChainIndexerList == nil {
		return errors.New("Missing blockChainIndexer config")
	}
	for i, blockChainIndexerAddr := range me.BlockChainIndexerList.AddressList {
		me.addressRepository[common.HexToAddress(blockChainIndexerAddr)] = "Block Chain Indexer " + strconv.Itoa(i)
	}

	if me.SideChainOwner == nil {
		return errors.New("Missing sideChainOwner config")
	}
	for i, sideChainOwnerAddr := range me.SideChainOwner.AddressList {
		me.addressRepository[common.HexToAddress(sideChainOwnerAddr)] = "Side Chain Owner " + strconv.Itoa(i)
	}

	if me.MainNetOwner == nil {
		return errors.New("Missing mainnetOwner config")
	}
	for i, TestNetOwnerAddr := range me.MainNetOwner.AddressList {
		me.addressRepository[common.HexToAddress(TestNetOwnerAddr)] = "Mainnet Owner " + strconv.Itoa(i)
	}

	return nil
}

func (me *ValidatorConfig) GetSmartContractAddress(chainId int, contractName string) common.Address {
	dict, ok := me.SmartContractDictionary[strconv.Itoa(chainId)]
	if !ok {
		return nullAddress

	}
	contractInfo, ok := dict.SmartContractMap[contractName]
	if !ok {
		return nullAddress
	}
	return common.HexToAddress(contractInfo.Address)
}

func (me *ValidatorConfig) GetCurrencyAdressMap(chainId int) map[string]common.Address {
	var currencyMap map[string]common.Address = make(map[string]common.Address)
	for _, currencyName := range me.CurrencyList {
		addr := me.GetSmartContractAddress(chainId, "OwnedUpgradeabilityProxy<"+currencyName+">")
		currencyMap[currencyName] = addr
	}
	return currencyMap
}

func (me *ValidatorConfig) GetNameByAddress(addr common.Address) string {
	agent, ok := me.addressRepository[addr]
	if !ok {
		return "Unknown"
	}
	return agent
}
