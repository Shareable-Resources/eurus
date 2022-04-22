package test

import (
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/secret"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestEurusUserDeposit_GetImplementation(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	impl, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Implementation", err)
	}
	fmt.Println("Implementation: ", impl.String())
	owner, err := proxy.ProxyOwner(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Owner", err)
	}
	fmt.Println("Owner: ", owner.String())
}

const depositObserverDecryptedPrivateKey string = "MIICXAIBAAKBgQDBUTY0ncrHX7xS0hmWbmSgO6vV1aYMlVViRETF4U8fyaGbsezol8LuMvGPjSIiPAIQGH3K9iDC6Z242PdKdzpgFL7D4TB0+1rxiHn/afcwLA+Rw7U20ZiQFkskbuqKivaZ2+Ct0gscfV/gJwmbzIoIRtl2UySYxmLXsNs4H01yHwIDAQABAoGADUheOBlLcI2EbBBhj7OAGH5hHS0z3pN4LWtRamNDw4RgJUmBZK3gx+saK+nfmYiT/7UfF433zEyu+J93xlcQ3KCKBKjz+fQe9lE1R2APJxP6kHLNoM2E96R4r8fh2GqMJoaduJlyR3SfjKDy1ltfO0Z8PL+dIHXRAot43ufQAnECQQDPEcasRnPenAtqjG+baNbEv0zFFC7GUTx3aMQTSkZ2Wb3JhSK4jcXlutCjkeT9SQssXKgLlwJyTvNdci7CpQ73AkEA7v+GgOM+kAtWrBQS1ESbDxnlSyVnJbo2mgiCUo0XFsw2XUXtf8V8288PJOXfuvtAqJX6O9lO+9F4KvczDT3kGQJBAKPGEXaMOnSkwrrA3Dz0jHkMPLHbJquf8M0YxYvkQRq2G89ZR37kUtNCEGZuq8hQj0/E8PxJsZurKfyMpMM6PT0CQFSyNz9LyOsRKZj30ChrW6wBWFHGIoSrNhhmNZD9sRYCLq3lTyI9oV7gRRSlZiEEU0irRa+Z9jSlafmH+w6RRVkCQAhivlRyiH05HwJewXk/6nwlUE/4ZI879IZH6oguL7D0y2buqQrc1v2abL+5xMW7WdyYQTY2P8GHZWcQiuNXOrk="

func TestEurusUserDeposit_GetWriterList(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	// for _, data := range sweepServiceAddr {
	for _, data := range depositObserverAddr {
		fmt.Println("predefined deposit observer:", data)
	}

	fmt.Println("predefined sweep invoker:", sweepInvokerAddr)

	ls, err := eurusUserDeposit.GetWriterList(&bind.CallOpts{})
	fmt.Println("=== Writer List ===")
	for _, data := range ls {
		fmt.Println(data.Hex())
	}
}

func TestEurusUserDeposit_AddWriter(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	//duncan debug
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	// for _, data := range sweepServiceAddr {
	// for _, data := range depositObserverAddr {
	// 	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	tx, err := eurusUserDeposit.AddWriter(transOpt, common.HexToAddress(data))
	// 	if err != nil {
	// 		//fatalSmartContractError(t, "Set Writer", err)
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	t.Logf("Tx Hash: %s", tx.Hash().Hex())
	// 	queryEthReceipt(t, &ethClient, tx)
	// }

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := eurusUserDeposit.AddWriter(transOpt, common.HexToAddress(sweepInvokerAddr))
	if err != nil {
		//fatalSmartContractError(t, "Set Writer", err)
		fmt.Println(err)
		t.Fatal(err)
	}
	t.Logf("Tx Hash: %s", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)

	ls, err := eurusUserDeposit.GetWriterList(&bind.CallOpts{})
	for _, data := range ls {
		fmt.Println(data.Hex())
	}
}

func TestEurusUserDeposit_RemoveWriter(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := eurusUserDeposit.RemoveWriter(transOpt, common.HexToAddress(testTestNetOwnerAddr))
	if err != nil {
		fatalSmartContractError(t, "Remove Writer", err)
	}

	t.Logf("Tx Hash: %s", tx.Hash().Hex())
}

func TestEurusUserDeposit_SetEurusInternalConfigAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := eurusUserDeposit.SetEurusInternalConfigAddress(transOpt, getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusInternalConfig>"))
	if err != nil {
		fatalSmartContractError(t, "Set Eurus Internal Config Address", err)
	}
	t.Logf("Tx Hash: %s", tx.Hash().Hex())
}

func TestEurusUserDeposit_SetEurusPlatformAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := eurusUserDeposit.SetEurusPlatformAddress(transOpt, getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>"))
	if err != nil {
		fatalSmartContractError(t, "Set Eurus Platform Address", err)
	}
	t.Logf("Tx Hash: %s", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)

}

func TestEurusUserDeposit_GetEurusPlatformAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	platformWalletAddress, err := eurusUserDeposit.EurusPlatformWalletAddress(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Get Eurus Platform Address", err)
	}
	fmt.Printf("Platform Wallet Address: %s\r\n", platformWalletAddress.Hex())
}

func TestEurusUserDeposit_GetEtherForwardAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	proxy, err := mainnet_contract.NewEtherForwardOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	etherForwardAddress, err := proxy.EtherForwardAddress(&bind.CallOpts{})
	if err != nil {
		fatalSmartContractError(t, "Get Ether Forward Address", err)
	}
	fmt.Printf("Ether Forward Address: %s\r\n", etherForwardAddress.Hex())
}

func TestEurusUserDeposit_SetEtherForwardAddress(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	proxy, err := mainnet_contract.NewEtherForwardOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	// This may not be the proxy owner private key so change it yourself
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	transOpt.GasLimit = 100000

	tx, err := proxy.SetEtherForwardAddress(transOpt, getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>"))
	if err != nil {
		fatalSmartContractError(t, "Set Ether Forward Address", err)
	}

	t.Logf("Tx Hash: %s", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusUserDeposit_TransferReceiveFallback(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")

	txHash, _, err := ethClient.TransferETHToSmartContract(testTestNetOwnerPrivateKey, addr.Hex(), big.NewInt(100000000), false, nil, nil)
	if err != nil {
		fatalSmartContractError(t, "Transfer ETH to OwnedUpgradeabilityProxy<EurusUserDeposit>", err)
	}

	t.Logf("Tx Hash: %s", txHash)
}

func TestEurusUserDeposit_SweepETH(t *testing.T) {
	mnenomicPhase := ""
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")

	balance, err := ethClient.GetBalance(addr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Balance: ", balance.String())

	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	wallet, account, address, err := secret.GenerateMintBurnKey(mnenomicPhase, depositObserverDecryptedPrivateKey, "61")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Deposit Observer Wallet Address: %s", address)

	privateKeyHex, err := wallet.PrivateKeyHex(*account)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(privateKeyHex, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	var hashByte [32]byte
	copy(hashByte[:], common.HexToHash("0x533b2cf0bf836d923ca96d188080a4c1ef9620dac002526d6187d732633480b0").Bytes()[:32])
	tx, err := eurusUserDeposit.Sweep(transOpt, hashByte, common.HexToAddress(testTestNetOwnerAddr), "ETH", big.NewInt(100))
	if err != nil {
		fatalSmartContractError(t, "Sweep ETH Error", err)
	}
	t.Logf("Tx Hash: %s", tx.Hash().Hex())
}

func TestEurusUserDeposit_SweepERC20(t *testing.T) {
	mnenomicPhase := ""
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")

	erc20Instance, err := contract.NewTestERC20(getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<USDT>"), ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	balance, err := erc20Instance.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Balance: ", balance.String())

	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	wallet, account, address, err := secret.GenerateMintBurnKey(mnenomicPhase, depositObserverDecryptedPrivateKey, "61")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Deposit Observer Wallet Address: %s", address)

	privateKeyHex, err := wallet.PrivateKeyHex(*account)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Private key hex: %s", privateKeyHex)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(privateKeyHex, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	var hashByte [32]byte
	copy(hashByte[:], common.HexToHash("0x06abaffd129f1efef676f6b1234c3d94ad8a7833b5108e291259d8d63a0d51ca").Bytes()[:32])
	tx, err := eurusUserDeposit.Sweep(transOpt, hashByte, common.HexToAddress(testTestNetOwnerAddr), "USDT", big.NewInt(100))
	if err != nil {
		//		fmt.Printf("Tx Hash: %s", tx.Hash().Hex())
		fatalSmartContractError(t, "Transfer ERC20 to OwnedUpgradeabilityProxy<EurusUserDeposit>", err)
	}
	fmt.Printf("Tx Hash: %s", tx.Hash().Hex())
	queryEthReceipt(t, &ethClient, tx)
}

func TestEurusUserDeposit_GetOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")
	eurusUserDeposit, err := mainnet_contract.NewEurusUserDeposit(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	ownerList, err := eurusUserDeposit.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, owner := range ownerList {
		fmt.Println("Owner: ", owner)
	}
}

func TestEurusUserDeposit_GetAllBalance(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusUserDeposit>")

	for _, assetName := range erc20Currency {
		if assetName == "ETH" {
			continue
		}
		tokenAddr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<" + assetName + ">")
		erc20, err := mainnet_contract.NewERC20(tokenAddr, ethClient.Client)
		if err != nil {
			t.Fatal(assetName + ": " + err.Error())
		}

		balance, err := erc20.BalanceOf(&bind.CallOpts{}, addr)
		if err != nil {
			t.Fatal(assetName + ": " + err.Error())
		}

		fmt.Println("Asset: ", assetName, ": ", balance.String())
	}
}
