package sc_validator

import (
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var nullAddress common.Address = common.Address{}

type SmartContractValidator struct {
	config             *ValidatorConfig
	sideChainEthClient *ethereum.EthClient
	mainnetEthClient   *ethereum.EthClient
	UserSelect         *UserSelect
}

func NewSmartContractValidator() *SmartContractValidator {
	validator := new(SmartContractValidator)
	validator.config = new(ValidatorConfig)
	validator.UserSelect = new(UserSelect)
	return validator
}

func (me *SmartContractValidator) LoadConfig(configPath string, scListPath string) error {

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &me.config)
	if err != nil {
		return errors.WithMessage(err, "Deserialize JSON file failed")
	}

	err = me.config.Load()
	if err != nil {
		return errors.WithMessage(err, "Unable to load config")
	}

	file, err = os.OpenFile(scListPath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return errors.WithMessage(err, "Unable to load smart contract address JSON ")
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.WithMessage(err, "Unable to read file")

	}

	err = me.config.LoadSmartContractConfig(rawData)
	if err != nil {
		return errors.WithMessage(err, "Unable to load smart contract config")
	}

	return nil
}

func (me *SmartContractValidator) Connect() error {
	if me.UserSelect.Network == SideChain || me.UserSelect.Network == AllNetwork {
		if me.config.SideChainEthIp != "" {
			me.sideChainEthClient = &ethereum.EthClient{
				ChainID:  big.NewInt(int64(me.config.SideChainId)),
				IP:       me.config.SideChainEthIp,
				Port:     me.config.SideChainEthPort,
				Protocol: me.config.SideChainEthProtocol,
			}

			_, err := me.sideChainEthClient.Connect()
			if err != nil {
				return errors.WithMessage(err, "Unable to connect to side chain")
			}
		}
	}
	if me.UserSelect.Network == MainNet || me.UserSelect.Network == AllNetwork {
		if me.config.MainnetEthIp != "" {
			me.mainnetEthClient = &ethereum.EthClient{
				ChainID:  big.NewInt(int64(me.config.MainnetChainId)),
				IP:       me.config.MainnetEthIp,
				Port:     me.config.MainnetEthPort,
				Protocol: me.config.MainnetEthProtocol,
			}

			_, err := me.mainnetEthClient.Connect()
			if err != nil {
				return errors.WithMessage(err, "Unable to connect to mainnet")
			}
		}
	}
	return nil
}

func (me *SmartContractValidator) ValidatePlatformWallet() error {
	printValidationHeader("Validating PlatformWallet")
	var isSuccess bool = false

	defer func() {
		printValidationResult("Validate PlatformWallet", isSuccess)
	}()

	addr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<PlatformWallet>")
	if addr == nullAddress {
		return errors.New("OwnedUpgradeabilityProxy<PlatformWallet> address not found")
	}
	platformWalletSC, err := contract.NewPlatformWallet(addr, me.sideChainEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewPlatformWallet")
	}

	fmt.Println("Checking Wallet operators")
	err = me.checkWalletOperatorList(platformWalletSC, me.config.DepositObserverList.AddressList, "Deposit observer")
	if err != nil {
		return err
	}

	err = me.checkWalletOwner(platformWalletSC, me.config.SideChainOwner.AddressList)
	if err != nil {
		return err
	}

	err = me.checkOwnerList(platformWalletSC)
	if err != nil {
		return err
	}

	err = me.checkInternalSmartContractConfig(platformWalletSC)
	if err != nil {
		return err
	}

	isSuccess = true
	return nil
}

func (me *SmartContractValidator) ValidateWithdrawSmartContract() error {
	var isSuccess bool = false

	defer func() {
		printValidationResult("Validate WithdrawSmartContract", isSuccess)
	}()

	printValidationHeader("Validate WithdrawSmartContract")

	withdrawAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	if withdrawAddr == nullAddress {
		return errors.New("Unable to find WithdrawSmartContract address")
	}

	withdrawSC, err := contract.NewWithdrawSmartContract(withdrawAddr, me.sideChainEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewWithdrawSmartContract")
	}

	fmt.Println("Checking wallet operator list with withdraw observer list")
	err = me.checkWalletOperatorList(withdrawSC, me.config.WithdrawObserverList.AddressList, "Withdraw Observer")
	if err != nil {
		return err
	}

	fmt.Println("Checking wallet operator list with ApprovalWallet")
	userApprovalWalletAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<ApprovalWallet>")
	if userApprovalWalletAddr == nullAddress {
		return errors.New("ApprovalWallet address not found")
	}

	err = me.checkWalletOperatorList(withdrawSC, []string{userApprovalWalletAddr.Hex()}, "Approval wallet")
	if err != nil {
		return err
	}

	fmt.Println("Checking writer list")
	writerList, err := withdrawSC.GetWriterList(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWriterList")
	}
	printAddressList("Writer list", writerList, me.config)
	isFound, _ := isAddressInList(writerList, userApprovalWalletAddr)
	if !isFound {
		return errors.New("Approval wallet address is not in writer list")
	}

	err = me.checkWalletOwner(withdrawSC, me.config.SideChainOwner.AddressList)
	if err != nil {
		return nil
	}

	err = me.checkOwnerList(withdrawSC)
	if err != nil {
		return err
	}

	isSuccess = true
	return nil
}

func (me *SmartContractValidator) ValidateApprovalWalletForUser() error {
	printValidationHeader("Validate ApprovalWallet for user")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate ApprovalWallet", isSuccess)
	}()

	approvalWalletAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<ApprovalWallet>")
	if approvalWalletAddr == nullAddress {
		return errors.New("ApprovalWallet address not found")
	}
	approvalWalletSC, err := contract.NewApprovalWallet(approvalWalletAddr, me.sideChainEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewApprovalWallet")
	}
	err = me.checkWalletOperatorList(approvalWalletSC, me.config.ApprovalObserverList.AddressList, "Approval observers")
	if err != nil {
		return err
	}

	currencyMap := me.config.GetCurrencyAdressMap(me.config.SideChainId)
	var currencyAddrList []string
	for _, currencyAddr := range currencyMap {
		currencyAddrList = append(currencyAddrList, currencyAddr.Hex())
	}

	err = me.checkInternalSmartContractConfig(approvalWalletSC)
	if err != nil {
		return err
	}
	err = me.checkWriter(approvalWalletSC, currencyAddrList)
	if err != nil {
		return err
	}
	isSuccess = true
	return nil
}

func (me *SmartContractValidator) ValidateERC20() error {

	printValidationHeader("Validate Validate ERC20")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate Validate ERC20", isSuccess)
	}()

	fmt.Println("Checking ERC20 addresses")
	var errorFound = false
	currencyMap := me.config.GetCurrencyAdressMap(me.config.SideChainId)

	for currencyName, currencyAddr := range currencyMap {
		fmt.Printf("%s:\t%s\t", currencyName, currencyAddr.Hex())
		_, err := contract.NewEurusERC20(currencyAddr, me.sideChainEthClient.Client)
		if err != nil {
			fmt.Println(colorRed + "FAILED " + err.Error() + colorReset)
			errorFound = true
		} else {
			fmt.Println("PASS")
		}
	}
	if errorFound {
		return errors.New("ERC20 failed")
	}

	platformWalletAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<PlatformWallet>")
	if platformWalletAddr == nullAddress {
		return errors.New("PlatformWallet address not found")
	}

	withdrawSCAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	if withdrawSCAddr == nullAddress {
		return errors.New("WithdrawSmartContract address not found")
	}

	var addrList []string = []string{
		platformWalletAddr.Hex(),
		withdrawSCAddr.Hex(),
	}

	for currencyName, currencyAddr := range currencyMap {
		fmt.Println("Checking writer list for " + currencyName)
		erc20SC, _ := contract.NewEurusERC20(currencyAddr, me.sideChainEthClient.Client)
		err := me.checkWriter(erc20SC, addrList)
		if err != nil {
			fmt.Printf("%s check writer failed: %v", currencyName, err)
			errorFound = true
		}
	}

	if errorFound {
		return errors.New("Checking writer failed")
	}

	fmt.Println("Checking owner list")
	for currencyName, currencyAddr := range currencyMap {
		fmt.Println("Checking writer list for " + currencyName)
		erc20SC, _ := contract.NewEurusERC20(currencyAddr, me.sideChainEthClient.Client)
		err := me.checkOwnerList(erc20SC)
		if err != nil {
			fmt.Printf("%s check writer failed: %v", currencyName, err)
			errorFound = true
		}
	}

	if errorFound {
		return errors.New("Checking owner list failed")
	}

	internalSCAddrCompare := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

	for currencyName, currencyAddr := range currencyMap {
		fmt.Println("Checking internal smart contract config address for " + currencyName)
		erc20SC, _ := contract.NewEurusERC20(currencyAddr, me.sideChainEthClient.Client)
		internalSCAddr, err := erc20SC.GetInternalSCConfigAddress(&bind.CallOpts{})
		if err != nil {
			errorFound = true
			fmt.Printf("Get internal smart contract config failed on currency %s, error: %v\r\n", currencyName, err)
		}
		if internalSCAddrCompare != internalSCAddr {
			errorFound = true
			fmt.Printf("Unmatch internal smart contract config address for currency %s, Error: %v\r\n", currencyName, err)
		}
	}
	if errorFound {
		return errors.New("Checking internal smart contract address failed")
	}
	isSuccess = true

	return nil
}

func (me *SmartContractValidator) ValidateInternalSmartContractConfig() error {
	printValidationHeader("Validate Internal smart contract config")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate Internal smart contract config", isSuccess)
	}()

	addr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

	internalSC, err := contract.NewInternalSmartContractConfig(addr, me.sideChainEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewInternalSmartContractConfig")
	}

	approvalWalletAddr, err := internalSC.GetApprovalWalletAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetApprovalWalletAddress")
	}
	addr = me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<ApprovalWallet>")
	if addr != approvalWalletAddr {
		return errors.New("Unmatched approval wallet address")
	}

	extSCConfigAddr, err := internalSC.GetExternalSCConfigAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetExternalSCConfigAddress")
	}
	addr = me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	if addr != extSCConfigAddr {
		return errors.New("Unmatched external smart contract config address")
	}

	platformWalletAddr, err := internalSC.GetInnetPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetInnetPlatformWalletAddress")
	}
	addr = me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<PlatformWallet>")
	if addr != platformWalletAddr {
		fmt.Printf("addr: %s, platformWalletAddr: %s", addr.Hex(), platformWalletAddr.Hex())
		return errors.New("Unmatched platform wallet address")
	}

	walletAddrMapAddr, err := internalSC.GetWalletAddressMap(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWalletAddressMap")
	}
	addr = me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<WalletAddressMap>")
	if addr != walletAddrMapAddr {
		return errors.New("Unmatched wallet address map address")
	}

	withdrawSCAddr, err := internalSC.GetWithdrawSmartContract(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWithdrawSmartContract")
	}
	addr = me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<WithdrawSmartContract>")
	if addr != withdrawSCAddr {
		return errors.New("Unmatched withdraw smart contract address")
	}

	isSuccess = true
	return nil
}

func (me *SmartContractValidator) ValidateExternalSmartContractConfig() error {
	printValidationHeader("Validate external smart contract config")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate external smart contract config", isSuccess)
	}()

	externalSCConfigAddr := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<ExternalSmartContractConfig>")
	if externalSCConfigAddr == nullAddress {
		return errors.New("Unable to get ExternalSmartContractConfig address from config")
	}

	externalSC, err := contract.NewExternalSmartContractConfig(externalSCConfigAddr, me.sideChainEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "NewExternalSmartContractConfig")
	}

	var errorFound bool = false
	currencyAddrMap := me.config.GetCurrencyAdressMap(me.config.SideChainId)
	for currencyName, currencyAddr := range currencyAddrMap {
		addr, err := externalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, currencyName)
		if err != nil {
			errorFound = true
			fmt.Printf("Unable to get ERC20 address for currency: %s, Error: %s\r\n", currencyName, err.Error())
			continue
		}
		if currencyAddr != addr {
			errorFound = true
			fmt.Printf("Unmatch ERC20 currency address for currency: %s\r\n", currencyName)
			continue
		}
	}

	if errorFound {
		return errors.New("Checking ERC20 address failed")
	}

	isSuccess = true
	return nil
}

func (me *SmartContractValidator) checkOwnerList(smartContract ISmartContractOwnerable) error {
	fmt.Println("Checking owner")
	ownerList, err := smartContract.GetOwners(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetOwners")
	}

	printAddressList("Owner list", ownerList, me.config)

	isFound, notFoundIndexList := isSubSet(ownerList, me.config.SideChainOwner.AddressList)
	if !isFound {
		message := itemToString(notFoundIndexList, me.config.SideChainOwner.AddressList)
		return errors.Errorf("Owner address not found. %s", message)
	}
	return nil
}

func (me *SmartContractValidator) checkTestNetOwnerList(smartContract ISmartContractOwnerable) error {
	fmt.Println("Checking owner")
	ownerList, err := smartContract.GetOwners(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetOwners")
	}

	printAddressList("Owner list", ownerList, me.config)

	isFound, notFoundIndexList := isSubSet(ownerList, me.config.MainNetOwner.AddressList)
	if !isFound {
		message := itemToString(notFoundIndexList, me.config.MainNetOwner.AddressList)
		return errors.Errorf("Owner address not found. %s", message)
	}
	return nil
}

func (me *SmartContractValidator) checkWalletOperatorList(smartContract ISmartContractWalletOwnerable,
	expectedAddrList []string, caption string) error {
	walletOperatorList, err := smartContract.GetWalletOperatorList(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWalletOperatorList")
	}

	printStringList(caption, expectedAddrList, me.config)
	printAddressList("Wallet operator list", walletOperatorList, me.config)

	isMatch, notFoundIndexList := isSubSet(walletOperatorList, expectedAddrList)
	if !isMatch {
		message := itemToString(notFoundIndexList, expectedAddrList)
		return errors.Errorf("%s address not found. %s", caption, message)
	}
	return nil
}

func (me *SmartContractValidator) checkWalletOwner(smartContract ISmartContractWalletOwnerable, expectedWalletOwnerList []string) error {
	fmt.Println("Checking wallet owner")
	walletOwner, err := smartContract.GetWalletOwner(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWalletOwner")
	}
	printAddressList("Smart contract Wallet owner", []common.Address{walletOwner}, me.config)

	isMatch, notFoundIndexList := isSubSet([]common.Address{walletOwner}, expectedWalletOwnerList)
	if !isMatch {
		message := itemToString(notFoundIndexList, expectedWalletOwnerList)
		return errors.Errorf("Wallet owner address not found. %s", message)
	}
	return nil
}

func (me *SmartContractValidator) checkWriter(smartContract ISmartContractReadWritePermissionable, expectedWriter []string) error {
	writerList, err := smartContract.GetWriterList(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetWriterList")
	}

	printAddressList("Smart contract writer list", writerList, me.config)
	isMatch, notFoundIndexList := isSubSet(writerList, expectedWriter)
	if !isMatch {
		message := itemToString(notFoundIndexList, expectedWriter)
		return errors.Errorf("Writer address not found. %s", message)
	}
	return nil
}

func (me *SmartContractValidator) checkInternalSmartContractConfig(smartContract ISmartContractInternalSCConfig) error {
	fmt.Println("Checking internal smart contract config address")
	internalSCAddr1 := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	if internalSCAddr1 == nullAddress {
		return errors.New("Unable to get internal smart contract config address from config")
	}
	ownerAddr := me.config.SideChainOwner.AddressList[0]

	internalSCAddr, err := smartContract.GetInternalSmartContractConfig(&bind.CallOpts{
		From: common.HexToAddress(ownerAddr),
	})
	if err != nil {
		return errors.WithMessage(err, "GetInternalSmartContractConfig failed")
	}
	if internalSCAddr != internalSCAddr1 {
		return errors.New("Unmatched internal smart contract config")
	}
	return nil
}

func (me *SmartContractValidator) checkInternalSmartContractConfig1(smartContract ISmartContractInternalSCConfig1) error {
	fmt.Println("Checking internal smart contract config address")
	internalSCAddr1 := me.config.GetSmartContractAddress(me.config.SideChainId, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	if internalSCAddr1 == nullAddress {
		return errors.New("Unable to get internal smart contract config address from config")
	}
	internalSCAddr, err := smartContract.InternalSmartContractConfig(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "GetInternalSmartContractConfig failed")
	}
	if internalSCAddr != internalSCAddr1 {
		return errors.New("Unmatched internal smart contract config")
	}
	return nil
}

func (me *SmartContractValidator) ValidateMainnetEurusInternalConfig() error {
	printValidationHeader("Validate Internal smart contract config")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate Internal smart contract config", isSuccess)
	}()
	addr := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<EurusInternalConfig>")

	eurusInternalSC, err := contract.NewEurusInternalConfig(addr, me.mainnetEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "EurusPlatformWallet")
	}
	platformWalletAddr, err := eurusInternalSC.PlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "Get EurusPlatformWalletAddress")
	}
	configPlatformWalletAddr := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	if configPlatformWalletAddr != platformWalletAddr {
		fmt.Printf("addr in config: %s, Eurus platformWalletAddr in SC: %s", configPlatformWalletAddr.Hex(), platformWalletAddr.Hex())
		return errors.New("Unmatched Eurus platform wallet address")
	}

	fmt.Println(" Eurus platformwallet address PASS")

	userDepositAddress, err := eurusInternalSC.EurusUserDepositAddress(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "Get EurusUserDepositAddress")
	}

	configUserDepositAddress := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<EurusUserDeposit>")

	if userDepositAddress != configUserDepositAddress {
		fmt.Printf("addr in config: %s,  Eurus user Deposit Address in SC: %s\n", configUserDepositAddress.Hex(), userDepositAddress.Hex())
		return errors.New("Unmatched Eurus user Deposit Address")
	}
	fmt.Println("Eurus user Deposit Address PASS")

	err = me.checkTestNetOwnerList(eurusInternalSC)
	if err != nil {
		return err
	}

	var countCorrect = 0
	for _, currency := range me.config.CurrencyList {
		tempAddr, err := eurusInternalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, currency)
		if err != nil {
			return errors.WithMessage(err, "Fail to get "+currency+" address")
		}
		configCurrencyAddr := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<"+currency+">")
		if tempAddr != configCurrencyAddr {
			fmt.Printf("%s check address failed: %s . config addr : %s , currency SC addr : %s", currency, configCurrencyAddr.String(), tempAddr.String())
		} else {
			countCorrect += 1
			fmt.Printf("%s  :  %s\n", currency, "PASS")
		}

	}
	if countCorrect == len(me.config.CurrencyList) {
		isSuccess = true
	}

	return nil
}

func (me *SmartContractValidator) ValidateMainnetEurusPlatformWalletConfig() error {
	printValidationHeader("Validate Eurus Platform Wallet config")

	var isSuccess bool = false
	defer func() {
		printValidationResult("Validate Eurus Platform Wallet config", isSuccess)
	}()

	addr := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	eurusPlatformWallet, err := contract.NewEurusPlatformWallet(addr, me.mainnetEthClient.Client)
	if err != nil {
		return errors.WithMessage(err, "EurusPlatformWallet")
	}

	eurusInternalConfig, err := eurusPlatformWallet.GetEurusInternalConfig(&bind.CallOpts{})
	if err != nil {
		return errors.WithMessage(err, "Can not get Eurus Internal Config")
	}

	configEurusInternalConfig := me.config.GetSmartContractAddress(me.config.MainnetChainId, "OwnedUpgradeabilityProxy<EurusInternalConfig>")

	if eurusInternalConfig != configEurusInternalConfig {
		fmt.Printf("addr in config: %s,  Eurus Internal Config Address in SC: %s\n", configEurusInternalConfig.Hex(), eurusInternalConfig.Hex())
		return errors.New("Unmatched Eurus Internal Config Address")
	}

	fmt.Println(" Eurus Internal Config address PASS")

	err = me.checkTestNetOwnerList(eurusPlatformWallet)
	if err != nil {
		return err
	}

	err = me.checkWriter(eurusPlatformWallet, me.config.WithdrawObserverList.AddressList)
	if err != nil {
		return err
	}

	isSuccess = true
	return nil
}
