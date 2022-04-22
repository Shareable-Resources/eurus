package test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/sign_service/sign_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	go_crypto "github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

var walletAddrList = []string{
	"0x1815F0C6Ce4453663cBED42637A054CC2392d90c",
}

func TestDeployProxy(t *testing.T) {
	ethClient := initEthClient(t)
	authClient := auth.NewAuthClient()
	authClient.SetLoginToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJFdXJ1c0F1dGgiLCJuYmYiOjE2MjMzMDk1NTcsImNsaWVudEluZm8iOiJ7XCJDbGllbnRJbmZvXCI6XCJ7XFxcInNlcnZpY2VJZFxcXCI6IDExNiwgXFxcInNlc3Npb25JZFxcXCI6MTM1fVwiLFwiU2VydmljZUlkXCI6MH0ifQ.aB6IhS1eBabjH_K7VtI9WN8oV46vnO_stZmpO6QBZDE")
	transOpt, err := ethClient.GetNewTransactorFromSignServer(authClient, "http://127.0.0.1:8083", sign_api.WalletKeyUserWalletOwner)

	//transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	//transOpt.GasLimit = 999999999999
	if err != nil {
		t.Fatal(err)
	}
	newProxyAddress, tx, _, err := contract.DeployOwnedUpgradeabilityProxy(transOpt, ethClient.Client)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Deployed proxy proxy addr : ", newProxyAddress)
	queryEthReceipt(t, &ethClient, tx)

}

func TestUserWallet_GetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := common.HexToAddress("0x2004e93d1f948748e5d39d30af71d0cd8ae4db78") //Double Proxy

	proxySC, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}

	implAddr, err := proxySC.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}

	fmt.Println("Addr: ", implAddr.Hex())

}

func TestUserWallet_SetWalletUserWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)

	scAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<CentralizedUserDappTest>")
	newUserWallet, err := contract.NewUserWallet(scAddr, ethClient.Client)
	if err != nil {
		fmt.Println(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	owner, err := newUserWallet.SetWalletOwner(transOpt, common.HexToAddress("0x99B32dAD54F630D9ED36E193Bc582bbed273d666"))

	if err != nil {
		fmt.Println(err)
	}

	queryEthReceipt(t, &ethClient, owner)
}

// func TestUserWalletProxy_GetUserWalletImplementation(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	addr := common.HexToAddress("0x2004e93d1f948748e5d39d30af71d0cd8ae4db78")
// 	userProxy, _ := contract.NewUserWalletProxy(addr, ethClient.Client)
// 	userWalletAddr, _ := userProxy.GetUserWalletImplementation(&bind.CallOpts{})
// 	fmt.Println("User Wallet impl: ", userWalletAddr.Hex())

// }

func TestUserWallet_AddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)

	for _, userWalletAddr := range walletAddrList {
		userWalletSc, err := contract.NewUserWallet(common.HexToAddress(userWalletAddr), ethClient.Client)
		if err != nil {
			fmt.Println("NewUserWallet error on address: ", userWalletSc, " error: ", err)
			continue
		}
		for _, addr := range userObserverAddr {
			transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
			if err != nil {
				fmt.Println("GetNewTransactorFromPrivateKey error on address: ", userWalletSc, " error: ", err)
				continue
			}
			transOpt.GasLimit = 10000000000
			tx, err := userWalletSc.AddWalletOperator(transOpt, common.HexToAddress(addr))
			if err != nil {
				fmt.Println("AddWalletOperator error on address: ", userWalletSc, " error: ", err)
				continue
			}

			receipt, err := ethClient.QueryEthReceipt(tx)
			if err != nil {
				fmt.Println("QueryEthReceipt error on address: ", userWalletSc, " error: ", err)
				continue
			}

			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Receipt: ", string(receiptJson))
		}
	}
}

func TestUserWallet_AddWriter(t *testing.T) {
	ethClient := initEthClient(t)

	for _, userWalletAddr := range walletAddrList {
		userWalletSc, err := contract.NewUserWallet(common.HexToAddress(userWalletAddr), ethClient.Client)
		if err != nil {
			fmt.Println("NewUserWallet error on address: ", userWalletSc, " error: ", err)
			continue
		}

		for _, invokerAddr := range invokerAddrList {
			transOpt, err := ethClient.GetNewTransactorFromPrivateKey(userWalletSCOwnerPrivateKey, ethClient.ChainID)
			if err != nil {
				fmt.Println("GetNewTransactorFromPrivateKey error on address: ", userWalletSc, " error: ", err)
				continue
			}
			transOpt.GasLimit = 10000000
			tx, err := userWalletSc.AddWriter(transOpt, common.HexToAddress(invokerAddr))
			if err != nil {
				fmt.Println("AddWalletOperator error on address: ", userWalletSc, " error: ", err)
				continue
			}

			receipt, err := ethClient.QueryEthReceipt(tx)
			if err != nil {
				fmt.Println("QueryEthReceipt error on address: ", userWalletSc, " error: ", err)
				continue
			}

			receiptJson, _ := json.Marshal(receipt)
			fmt.Println("Receipt: ", string(receiptJson))
		}

	}
}

func TestSetInternalSCAddr(t *testing.T) {
	ethclient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<CentralizedUserDappTest>")
	transOpt, err := ethclient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethclient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	userWalletInstance, err := contract.NewUserWallet(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	intcAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	tx, err := userWalletInstance.SetInternalSmartContractConfig(transOpt, intcAddr)
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethclient, tx)
}

func TestUserWallet_ChangeRequirement(t *testing.T) {
	ethClient := initEthClient(t)

	for _, userWalletAddr := range walletAddrList {
		userWalletSc, err := contract.NewUserWallet(common.HexToAddress(userWalletAddr), ethClient.Client)
		if err != nil {
			fmt.Println("NewUserWallet error on address: ", userWalletSc, " error: ", err)
			continue
		}

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error on address: ", userWalletSc, " error: ", err)
			continue
		}
		transOpt.GasLimit = 10000000000
		tx, err := userWalletSc.ChangeRequirement(transOpt, big.NewInt(2))

		if err != nil {
			fmt.Println("AddWalletOperator error on address: ", userWalletSc, " error: ", err)
			continue
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("QueryEthReceipt error on address: ", userWalletSc, " error: ", err)
			continue
		}

		receiptJson, _ := json.Marshal(receipt)
		fmt.Println("Receipt: ", string(receiptJson))

	}
}

func TestUserWallet_GetRequirement(t *testing.T) {
	ethclient := initEthClient(t)

	addr := common.HexToAddress("0x88b7c329fb4531592b18664d2ab02b0f2903bae9")

	userWalletInstance, err := contract.NewUserWallet(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	required, err := userWalletInstance.Required(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Required: ", required.String())
}

func TestComfirmTransactionCount(t *testing.T) {
	ethclient := initEthClient(t)
	addr := getAddressBySmartContractName("UserWallet")
	//transOpt, err := ethclient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethclient.ChainID)
	//if err != nil {
	//	t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	//}
	userWalletInstance, err := contract.NewUserWallet(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	tx, err := userWalletInstance.GetConfirmationCount(&bind.CallOpts{}, big.NewInt(3))
	if err != nil {
		fmt.Println(err)
	}
	//queryEthReceipt(t,&ethclient,tx)
	fmt.Println(tx)
}

func TestComfirmTransaction(t *testing.T) {
	ethclient := initEthClient(t)
	addr := getAddressBySmartContractName("UserWallet")
	transOpt, err := ethclient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethclient.ChainID)
	//if err != nil {
	//	t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	//}
	//transOpt.GasLimit = 1000000
	userWalletInstance, err := contract.NewUserWallet(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	tx, err := userWalletInstance.ConfirmTransaction(transOpt, big.NewInt(3))
	//tx, err := userWalletInstance.Required(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethclient, tx)
	fmt.Println(tx)
}

func TestTransactionCheck(t *testing.T) {
	ethclient := initEthClient(t)
	addr := getAddressBySmartContractName("UserWallet")

	userWalletInstance, err := contract.NewUserWallet(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	//transOpt, err := ethclient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethclient.ChainID)
	//if err != nil {
	//	t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	//}
	//tx, err := userWalletInstance.RejectTransaction(transOpt, big.NewInt(24))
	tx, err := userWalletInstance.Transactions(&bind.CallOpts{}, big.NewInt(3))
	if err != nil {
		fmt.Println(err)
	}
	//queryEthReceipt(t,&ethclient,tx)
	fmt.Println(tx.Value)
}

func TestUserWallet_GetWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0xacfee88888b0efa9997a80114225e2212b8ec90f"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}
	walletOperatorList, _ := userWalletInstance.GetWalletOperatorList(&bind.CallOpts{})
	for _, walletOperatorAddr := range walletOperatorList {
		fmt.Println(walletOperatorAddr)
	}
	fmt.Println("Owner:")
	ownerList, _ := userWalletInstance.GetOwners(&bind.CallOpts{})
	for _, ownerAddr := range ownerList {
		fmt.Println(ownerAddr)
	}

	walletOwner, _ := userWalletInstance.GetWalletOwner(&bind.CallOpts{})
	fmt.Println("Wallet owner: ", walletOwner.Hex())
}

func TestUserWallet_GetOwner(t *testing.T) {
	ethClient := initEthClient(t)
	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x25f4585ded01bda1cf644df857910abd711b7552"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	ownerList, err := userWalletInstance.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	for _, owner := range ownerList {
		fmt.Println("Owner: ", owner.Hex())
	}

}

func TestInitUserWallet(t *testing.T) {
	//TestUpgradeToNewProxy(t)
	TestSetInternalSCAddr(t)
	TestUserWallet_SetWalletUserWalletOwner(t)

}

type DummyEthClientCommunicator struct {
	ethereum.EthClientCommunicator
}

func (me *DummyEthClientCommunicator) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func TestUserWallet_RequestTransfer(t *testing.T) {
	ethClient := initEthClient(t)
	//communicator := new(DummyEthClientCommunicator)
	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x1f60190ac3aca7d3001c44bbebf5f85cdc684e6f"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.devabcde123456")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}
	privateKey, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		t.Fatal("Private key error: ", err)
	}
	fmt.Println("Payment Wallet address: ", account.Address.Hex())

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(hex.EncodeToString(privateKey), ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}

	transOpt.GasLimit = 300000
	// transOpt.NoSend = true
	amount := big.NewInt(1000000)
	// sign, err := generateTransferSignature("0xb13070bd35176fae45ff692d5279495547473974", "USDT", amount)
	// if err != nil {
	// 	t.Fatal("GenerateTransferSignature error: ", err)
	// }
	tx, err := userWalletInstance.RequestTransferV1(transOpt,
		common.HexToAddress("0xb13070bd35176fae45ff692d5279495547473974"),
		"USDT",
		amount,
		/*sign*/)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(hex.EncodeToString(tx.Data()))
	// fmt.Println("Sign:")
	// v, r, s := tx.RawSignatureValues()
	// fmt.Println(v, " ", r, " ", s)
	fmt.Println("Trans hash: ", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func generateTransferSignature(dest string, assetName string, amount *big.Int) ([]byte, error) {
	addrType, err := abi.NewType("address", "address", nil)
	if err != nil {
		return nil, err
	}

	strType, err := abi.NewType("string", "string", nil)
	if err != nil {
		return nil, err
	}

	intType, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return nil, err
	}
	arg := abi.Argument{
		Name: "dest",
		Type: addrType,
	}

	arg1 := abi.Argument{
		Name: "assetName",
		Type: strType,
	}

	arg2 := abi.Argument{
		Name: "amount",
		Type: intType,
	}
	argList := abi.Arguments{arg, arg1, arg2}

	packedByte, err := argList.Pack(common.HexToAddress(dest), assetName, amount)
	if err != nil {
		return nil, err
	}
	fmt.Println("PackedByte: ")
	fmt.Println(packedByte)

	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write(packedByte)
	if err != nil {
		return nil, err
	}

	var resultHash []byte
	resultHash = hasher.Sum(resultHash)
	fmt.Println("Hash: ", resultHash)
	priKey, err := go_crypto.HexToECDSA("fa0cf2b11a0c1bb8227395f4cf89733d620e96c04dca4f2bc95e2b839d7dcfd2")
	if err != nil {
		return nil, err
	}
	addr := go_crypto.PubkeyToAddress(priKey.PublicKey)
	fmt.Println("Signature sign by address: ", addr.Hex())
	sign, err := go_crypto.Sign(resultHash, priKey)
	if err != nil {
		return nil, err
	}
	return sign, nil

}

func TestUserWallet_SubmitWithdraw(t *testing.T) {
	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x1f60190ac3aca7d3001c44bbebf5f85cdc684e6f"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	amount := big.NewInt(1000000)
	amountWithFee := big.NewInt(15000000)
	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.devabcde123456")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}
	privateKey, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		t.Fatal("Private key error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(hex.EncodeToString(privateKey), ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey: ", err)
	}
	transOpt.GasLimit = 1000000
	// sign, err := generateSubmitWitdrawSignature("0xb13070bd35176fae45ff692d5279495547473974", "USDT", amount, amountWithFee)
	// if err != nil {
	// 	t.Fatal("GenerateTransferSignature error: ", err)
	// }
	tx, err := userWalletInstance.SubmitWithdrawV1(transOpt,
		common.HexToAddress("0x5b40936e78bb80d2f71c8899076b7b7a636c6541"),
		amount,
		amountWithFee,
		"USDT",
		/*sign*/)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Trans hash: ", tx.Hash().Hex())

	queryEthReceipt(t, &ethClient, tx)
}

func generateSubmitWitdrawSignature(dest string, assetName string, amount *big.Int, amountWithFee *big.Int) ([]byte, error) {
	addrType, err := abi.NewType("address", "address", nil)
	if err != nil {
		return nil, err
	}

	strType, err := abi.NewType("string", "string", nil)
	if err != nil {
		return nil, err
	}

	intType, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return nil, err
	}
	arg := abi.Argument{
		Name: "dest",
		Type: addrType,
	}

	arg1 := abi.Argument{
		Name: "withdrawAmount",
		Type: intType,
	}

	arg2 := abi.Argument{
		Name: "amountWithFee",
		Type: intType,
	}

	arg3 := abi.Argument{
		Name: "assetName",
		Type: strType,
	}
	argList := abi.Arguments{arg, arg1, arg2, arg3}

	packedByte, err := argList.Pack(common.HexToAddress(dest), amount, amountWithFee, assetName)
	if err != nil {
		return nil, err
	}

	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write(packedByte)
	if err != nil {
		return nil, err
	}

	var resultHash []byte
	resultHash = hasher.Sum(resultHash)

	priKey, err := go_crypto.HexToECDSA("fa0cf2b11a0c1bb8227395f4cf89733d620e96c04dca4f2bc95e2b839d7dcfd2")
	if err != nil {
		return nil, err
	}
	addr := go_crypto.PubkeyToAddress(priKey.PublicKey)
	fmt.Println("Signature sign by address: ", addr.Hex())
	sign, err := go_crypto.Sign(resultHash, priKey)
	if err != nil {
		return nil, err
	}
	return sign, nil

}

func TestUserWallet_IsWriter(t *testing.T) {

	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x88b7c329fb4531592b18664d2ab02b0f2903bae9"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	isWriter, err := userWalletInstance.IsWriter(&bind.CallOpts{}, common.HexToAddress(invokerAddrList[0]))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("IsWriter: ", isWriter)
}

func TestUserWallet_GetGasFeeWalletAddress(t *testing.T) {
	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x88b7c329fb4531592b18664d2ab02b0f2903bae9"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	addr, _ := userWalletInstance.GetGasFeeWalletAddress(&bind.CallOpts{})
	fmt.Println("Gas fee wallet address: ", addr.Hex())

}

func TestUserWallet_InvokeAbitarySCCall(t *testing.T) {

	userWalletOwnerPrivateKey := "64f8ba795cf8f78e9c3c7a1b154326ba6e0e6f994e4853f0a551c15519fb438e"

	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x1815F0C6Ce4453663cBED42637A054CC2392d90c"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(userWalletOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("transopt: ", err)
	}
	unitTestAddr := getAddressBySmartContractName("UnitTest")

	unitTestInstance, err := contract.NewUnitTest(unitTestAddr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt.NoSend = true
	transOpt.GasLimit = 10000000
	value := big.NewInt(12)
	tx, err := unitTestInstance.SetValue(transOpt, value)
	if err != nil {
		t.Fatal("InvokeRequireMessage transaction: ", err)
	}

	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write(tx.Data())

	walletOwnerPrivateKey := "64f8ba795cf8f78e9c3c7a1b154326ba6e0e6f994e4853f0a551c15519fb438e"
	priKey, err := go_crypto.HexToECDSA(walletOwnerPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := go_crypto.Sign(hasher.Sum(nil), priKey)
	if err != nil {
		t.Fatal(err)
	}
	invokerPrivateKey := "92bd179b40e3fa7464853e3e5d989d8aaf5a73d6767f241f46efe51795e664dc"
	transOptInvoker, err := ethClient.GetNewTransactorFromPrivateKey(invokerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("transOptInvoker: ", err)
	}
	transOptInvoker.GasLimit = 10000000
	fmt.Println("Signature: ", hex.EncodeToString(signature), " Input args: ", hex.EncodeToString(tx.Data()))
	eun := big.NewInt(0)
	actualTx, err := userWalletInstance.InvokeSmartContract(transOptInvoker, unitTestAddr, eun, tx.Data())
	queryEthReceipt(t, &ethClient, actualTx)
}

func TestUserWallet_TopupPaymentAddress(t *testing.T) {
	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x1f60190ac3aca7d3001c44bbebf5f85cdc684e6f"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	invokerPrivateKey := "92bd179b40e3fa7464853e3e5d989d8aaf5a73d6767f241f46efe51795e664dc"
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(invokerPrivateKey, ethClient.ChainID)
	amount := big.NewInt(0)
	amount, _ = amount.SetString("90000000", 10)
	paymentWalletAddr := common.HexToAddress("0x24994b04f9c1eafd3df432edd777ec3874c8d421")

	hasher := sha3.NewLegacyKeccak256()

	addrType, _ := abi.NewType("address", "address", nil)

	arg := abi.Argument{
		Name: "paymentWalletAddr",
		Type: addrType,
	}
	argList := abi.Arguments{arg}
	packedData, err := argList.Pack(paymentWalletAddr)
	_, err = hasher.Write(packedData)

	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.devabcde123456")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		t.Fatal("Private key error: ", err)
	}
	fmt.Println("Payment Wallet address: ", account.Address.Hex())
	signature, err := go_crypto.Sign(hasher.Sum(nil), privateKey)
	if err != nil {
		t.Fatal(err)
	}
	transOpt.GasLimit = 10000000
	tx, err := userWalletInstance.TopUpPaymentWallet(transOpt, paymentWalletAddr, amount, signature)
	if err != nil {
		t.Fatal("TopupPaymentWallet error: ", err)
	}
	fmt.Println("Tx hash: ", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestUserWallet_DirectTopupPaymentAddress(t *testing.T) {
	ethClient := initEthClient(t)

	userWalletInstance, err := contract.NewUserWallet(common.HexToAddress("0x1f60190ac3aca7d3001c44bbebf5f85cdc684e6f"), ethClient.Client)
	if err != nil {
		t.Fatal("NewUserWallet: ", err)
	}

	amount := big.NewInt(0)
	amount, _ = amount.SetString("100000000", 10)

	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.dev")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}
	privateKey, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		t.Fatal("Private key error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(hex.EncodeToString(privateKey), ethClient.ChainID)
	transOpt.GasLimit = 10000000
	tx, err := userWalletInstance.DirectTopUpPaymentWallet(transOpt, amount, big.NewInt(10000000))
	if err != nil {
		t.Fatal("TopupPaymentWallet error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestUserWallet_TransferBackToUserWallet(t *testing.T) {
	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.devabcde123456")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}
	privateKey, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		t.Fatal("Private key error: ", err)
	}
	fmt.Println("Payment Wallet address: ", account.Address.Hex())

	ethClient := initEthClient(t)

	_, tx, err := ethClient.TransferETH(hex.EncodeToString(privateKey), "0x1f60190ac3aca7d3001c44bbebf5f85cdc684e6f", big.NewInt(49831461100000000))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestUserWallet_GenerateTopUpWalletSignature(t *testing.T) {
	decryptedBase64, err := crypto.DecryptAESFromBase64("72cQxlGd7NzSkXBCVFCIQMZlEy6Wat+wPc/AP+4HzoVN/s9KKOV13dVvgMDQYqrMqVOUhq4GDV2+LG+A23dnFhEY2mzYkRqyPxggc/kxBTQDGrrQeG2EEtgpE13PNl+W", "DKdzPniHV8yO/FYtqI4wt0RWu5vo9C6HvUl0MfqDAFw=")
	if err != nil {
		t.Fatal("DecryptAESFromBase64 error: ", err)
	}
	data, err := base64.StdEncoding.DecodeString(decryptedBase64)
	if err != nil {
		t.Fatal("DecodeString error: ", err)
	}
	mnenonicPhase := string(data)
	seed, err := bip39.NewSeedWithErrorChecking(mnenonicPhase, "eu21-9@18m.devabcde123456")
	if err != nil {
		t.Fatal("NewSeedWithErrorChecking error: ", err)
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		t.Fatal("NewFromSeed error: ", err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/1'/0/0")

	account, err := wallet.Derive(path, true)
	if err != nil {
		t.Fatal("Account error: ", err)
	}

	fmt.Println("Payment Wallet address: ", account.Address.Hex())
	privateKey, _ := wallet.PrivateKeyHex(account)

	addrType, err := abi.NewType("address", "address", nil)
	if err != nil {
		t.Fatal(err)
	}
	arg := abi.Argument{
		Name: "walletOwnerAddress",
		Type: addrType,
	}

	argList := abi.Arguments{arg}
	packed, err := argList.Pack(account.Address)
	if err != nil {
		t.Fatal(err)
	}

	var resultHash []byte
	resultHash = go_crypto.Keccak256(packed)

	sign, err := crypto.GenerateECDSSignature(resultHash, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Signature: ", hex.EncodeToString(sign))
}

func TestUserWallet_GetInternalSCAddr(t *testing.T) {
	ethclient := initEthClient(t)
	addr := common.HexToAddress("0x677f4D2930DF5484fFb5a248Fa014083e02a0460")

	userWalletInstance, err := contract.NewUserWalletProxy(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	internalSCAddr, err := userWalletInstance.GetInternalSCAddress(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(internalSCAddr.Hex())
}

func TestUserWallet_SetInteranlSCAddr(t *testing.T) {
	ethclient := initEthClient(t)
	addr := common.HexToAddress("0x677f4D2930DF5484fFb5a248Fa014083e02a0460")

	userWalletInstance, err := contract.NewUserWalletProxy(addr, ethclient.Client)
	if err != nil {
		fmt.Println(err)
	}
	transOpt, err := ethclient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethclient.ChainID)
	if err != nil {
		fmt.Println(err)
	}

	scAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	tx, err := userWalletInstance.SetInternalSCAddress(transOpt, scAddr)
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethclient, tx)
}

func TestUserWallet_SetOwner(t *testing.T) {

	var userWalletAddrList []string = []string{
		"0x2004e93d1f948748e5d39d30af71d0cd8ae4db78",
		"0x1dad7e04be1649e5cda02ea7e582f7be2f9d2103",
		"0x6dd67ee385b2f1a25c28fa27482e5138b639bd47",
		"0x94479487c057bb42597bd7e7d8265b2a08dd16ab",
		"0x3340bffccb1077dc2de052ddc3a9c3d8f8c9eafd",
		"0x1daf70fb148e0d99273706cafeea00be179f3b74",
		"0x1b1f42b0d6f949bd1b427c86a053d6e3a1d19d08",
		"0x42579d2a440de27c326aec057d79c71861c73c0f",
		"0x6e89d4552b64037227b4f57238ab601c63e682c0",
		"0xc5d7fcb4b6503d78156271081c4213dc5b4632d5",
		"0x7da27989f5b810faebbbb1e8069ebb56e96800a9",
		"0xe2b7c797ecba3605dc7ffb1c587f573e6c64357c",
		"0x0752fbab174591950ef4f89ec6d1b362ce57dca5",
		"0x55d9f1b16f4827b6034dbd3206e08bbed38eeb07",
		"0x54d5387253f10df0419bd6eb74cd97c69c128f87",
		"0x97df2101247a6794c21d093833f6fb4848798422",
		"0xad5b7d6a62a9476697082346ce6c9cdb38d16294",
		"0x85e431426d4fddf407dd575e73dfeef735c96f8d",
		"0xfa61d26dee3598e4bb267d502dc947e100983425",
		"0x9f932a4530bce0d37bdbc86524e5272abddb21a0",
		"0x9db7481baea4b847fa8b238cfc288067187b3130",
		"0x1815cdaa36d03fbebf6f98c9a95bdb340f76f475",
		"0x5ce4a04a264850d6786911288b3bee1f448fc89e",
		"0x2722d9f757254e55e2a8aaee89e8254cc1bf533e",
		"0x9dea7d2be2f6d0ea9a9c4146fb9080f61a9c1189",
		"0x1d60b58a917ebe3fccf14f0fa3e0c073a8feaed5",
		"0xe66f6d6b7961678765634b379360a17c691f19ec",
		"0x8e13239b205b0363681c01beb5e53a6cf729bf4e",
		"0x49210c570fd48b9aaa0dc8cf1ed70480a7c384b2",
		"0xe85b30e8db717d61c047095e161b6ab642014982",
		"0x9b004790614222e09aaa946cda51c1d0924ba2c4",
		"0x837794da1c3700bedc62f50f4aa372bd98687e19",
		"0xfc0cd27314ed4e38c08b474afc5027070d38f33a",
		"0x85226e8e1dde72fb54d93a2761fe46587cc729ea",
		"0x2c1f05acfe28d2522e9c5e1669b84471d4e74661",
		"0x7b333488c58cc79ac6b0edc5f54e477ec220f5fe",
		"0xc16ffe0da5dae266168df31f9891702277195c0f",
		"0x8d77691037c9fad03e3df1d638bd9c7bbf737e2d",
		"0x9c4b9cbc2af83e6d719a386731085c6bcc4f96fe",
		"0xeb38c2255244dc4fbba425cef34f3cf7a03e6eb6",
		"0xf377fd45471ef47fbb6aecc17a586d4dd3f48cd4",
		"0xdda51cf03b0624444a8683780d91927923e98191",
		"0x482ab6f4f73493cd3b22bddcf9b136cf1367ad60",
		"0x36a522ea9e93d68457e64b7409c4b8afff15b375",
		"0x56c45aa963cb87a8d474f176694c755e5b51d092",
		"0xf256761705f9f15b498f1954b10b6efe63d61c49",
		"0x7c488240c9100f4cda451354f49704781c07471a",
		"0x7caa8b87b6b27dcede99a185e0b7ecaf95425069",
		"0x4f4245fc801686af475a7501bd9c745020593f86",
		"0xe85f54653269d9fcecb15b5f089c97c3c2985020",
		"0x722981f4bcf5b464dcbecbe7b678e960afe13190",
		"0xa9f15eaf6aa6292070509ad40e4852b381e6e48c",
		"0xaefc260cc91713ad3643af3cc1fa9f4bcb5e133b",
		"0xae9a23c0ac3e8dc4c579197d83a4ff8605182989",
		"0x70e6ca28e28b9c154ca72a4a7f51ff403624be47",
		"0xeb52fede23825152a9a1b880f5188917237b1898",
		"0x303e189e85e234f3da02bc33716b6c7b288f9f5b",
		"0x2a65dd702f4895c53ddcfcf11ea1d13d1fb17836",
		"0x39d13aa326649cf3d38ebb38faf89551119ebb26",
		"0x4ca0786c42db02144de4af04f49354f8b34c908b",
		"0x1bc9ae462d3d5ece278e6bdd24f7d88f8468980a",
		"0x4c259d3c72d1f1acd763079624e74d42cc3024f5",
		"0xde55fef718aa7c316b3948f277918e3897d9b0a3",
		"0xa5fdb1d83483f4dc78d35cee4f3a8c9852f39285",
		"0x040078fb3772e904cdae529413c635545f81b20a",
	}
	log.NewLogger(log.Name.Root, "/tmp/UserWalletSetOwner.log", logrus.DebugLevel)
	ethClient := initEthClient(t)
	var errorUserWalletList []string = make([]string, 0)

	for _, userWalletAddr := range userWalletAddrList {

		userWallet, err := contract.NewUserWallet(common.HexToAddress(userWalletAddr), ethClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("cannot new user wallet", err.Error(), " WalletAddr: ", userWalletAddr)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		transOpt5, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
		//Adding sign invoker as
		log.GetLogger(log.Name.Root).Debugln("Add owner for ", userWalletAddr)
		tx5, err := userWallet.AddOwner(transOpt5, common.HexToAddress(testUserWalletSCOwnerAddr))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddOwner failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		receipt5, err := ethClient.QueryEthReceiptWithSetting(tx5, 1, 50)
		if receipt5.Status != 1 {
			log.GetLogger(log.Name.Root).Errorln("AddOwner receipt failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx5.Hash().Hex())
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		transOpt6, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
		//Adding sign invoker as
		log.GetLogger(log.Name.Root).Debugln("Remove owner for ", userWalletAddr)
		tx6, err := userWallet.RemoveOwner(transOpt6, common.HexToAddress(signServerAddr))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddOwner failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		receipt6, err := ethClient.QueryEthReceiptWithSetting(tx6, 1, 50)
		if receipt6.Status != 1 {
			log.GetLogger(log.Name.Root).Errorln("AddOwner receipt failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx6.Hash().Hex())
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
	}
}
