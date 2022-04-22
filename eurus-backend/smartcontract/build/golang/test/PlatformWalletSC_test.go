package test

import (
	"encoding/base64"
	"eurus-backend/foundation/crypto"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestPlatformWallet_SetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err := platformWalletSC.SetWalletOwner(transOpt, common.HexToAddress(testOwnerAddr))
	if err != nil {
		t.Fatal("SetWalletOwner error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestPlatformWallet_GetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

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

func TestPlatformWallet_UpgradeTo(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	proxySC, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	platformAddr := getAddressBySmartContractName("PlatformWallet")
	implAddr, err := proxySC.UpgradeTo(transOpt, platformAddr)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}

	fmt.Println("Addr: ", implAddr)

}

func TestPlatformWallet_TransferEUN(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	proxySC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	amount := big.NewInt(0)
	amount = amount.SetUint64(10000000000000)
	transOpt.GasLimit = 1000000
	tx, err := proxySC.TransferETH(transOpt, common.HexToAddress("0x8c9e314f4bd8dacde1ca4e9fd63064f7cd3388fa"), amount)
	if err != nil {
		t.Fatal("Transfer error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestPlatformWallet_SetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	proxySC, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	actualAddr := getAddressBySmartContractName("PlatformWallet")
	tx, err := proxySC.UpgradeTo(transOpt, actualAddr)
	if err != nil {
		t.Fatal("UpgradeTo error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)

}
func TestPlatformWallet_RemoveWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := platformWalletSC.RemoveWalletOperator(transOpt, common.HexToAddress("0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"))
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)

}
func TestPlatformWallet_AddWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	for _, addr := range depositObserverAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
		}

		tx, err := platformWalletSC.AddWalletOperator(transOpt, common.HexToAddress(addr))
		if err != nil {
			// t.Fatal("AddWalletOperator error: ", err)
			fmt.Println("address: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("receipt error: ", err, " for deposit observer address: ", addr)
			continue
		}
		receiptByte, _ := receipt.MarshalJSON()
		fmt.Println(string(receiptByte))
	}

	// for _, address := range userServerHDWalletAddr {
	// 	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	// 	if err != nil {
	// 		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	// 	}

	// 	tx, err := platformWalletSC.AddWalletOperator(transOpt, common.HexToAddress(address))
	// 	if err != nil {
	// 		// t.Fatal("AddWalletOperator error: ", err)
	// 		fmt.Println("address: ", err)
	// 		continue
	// 	}
	// 	receipt, err := ethClient.QueryEthReceipt(tx)
	// 	if err != nil {
	// 		fmt.Println("receipt error: ", err, " for user server address: ", address)
	// 		continue
	// 	}

	// 	receiptByte, _ := receipt.MarshalJSON()
	// 	fmt.Println(string(receiptByte))

	// }

	// queryEthReceipt(t, &ethClient, tx)

}

func TestPlatformWallet_GetWalletOwner(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	ownerAddr, err := platformWalletSC.GetWalletOwner(&bind.CallOpts{})
	if err != nil {
		t.Fatal("AddWalletOperator error: ", err)
	}

	fmt.Println("Wallet owner : ", ownerAddr.Hex())

}

func TestPlatformWallet_GetWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	writerList, err := platformWalletSC.GetWriterList(&bind.CallOpts{})
	if err != nil {
		t.Fatal("AddWalletOperator error: ", err)
	}
	for _, writer := range writerList {

		fmt.Println("Wallet writer : ", writer.Hex())
	}

}

func TestPlatformWallet_GetWalletOperator(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	operatorList, err := platformWalletSC.GetWalletOperatorList(&bind.CallOpts{})
	if err != nil {
		t.Fatal("AddWalletOperator error: ", err)
	}

	for _, operatorAddr := range operatorList {
		fmt.Println("Wallet owner : ", operatorAddr.Hex())
	}

}

func TestPlatformWallet_GetRequirement(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	tx, err := platformWalletSC.Required(&bind.CallOpts{})
	if err != nil {
		t.Fatal("ChangeRequirement error: ", err)
	}

	fmt.Println("Required : ", tx)
}

func TestPlatformWallet_ChangeRequirement(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	req := big.NewInt(5)
	tx, err := platformWalletSC.ChangeRequirement(transOpt, req)
	if err != nil {
		t.Fatal("ChangeRequirement error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

// func TestPlatformWallet_GetConfirmationList(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

// 	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
// 	if err != nil {
// 		t.Fatal("NewPlatform error: ", err)
// 	}
// 	arrList, err := platformWalletSC.GetConfirmations(&bind.CallOpts{}, big.NewInt(122))
// 	if err != nil {
// 		t.Fatal("GetConfirmations error: ", err)
// 	}

// 	for _, addr := range arrList {
// 		fmt.Println(addr)
// 	}
// }

// func TestPlatformWallet_GetConfirmation(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")

// 	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
// 	if err != nil {
// 		t.Fatal("NewPlatform error: ", err)
// 	}

// 	addressList := []string{
// 		"0x81ccE1BF6DCfCeF814914E66d1aD6ACcFeB3936e",
// 		"0x80B5aD3d6370f9C4b6b84a222f42CD3FC5755D81",
// 		"0x96ad0B6DdD7CEb4e231dA09F7AE5044b98De5362",
// 		"0xA14Ab677452e5c038067f6dF57eE073328FbD336",
// 		"0xEECd77226192B0Cf9CE510A254F31f03c3CD2D85",
// 		"0x6612457C18AB8A7B1caF25bd316401359132B5cd",
// 		"0xdF86024284b12A1a274183A920cFf1419D524F17",
// 		"0x206803F71C5eB93C61a060eeCfc94C2CE6D9747D",
// 		"0x15C9894F2C32A892A9965Fa23b6026Abd7450Ac8",
// 		"0x7076cE3E9C32c6D9F625E142556E0944BBDAe6a0",
// 		"0xffa498003B1b1AfC2804533fF7Fb11a5ACC059C2",
// 		"0xa4fC063Be769a24886d0fe219d9F1A24b09dAbC2",
// 		"0x38455a54cC578505b6a57Aa0957d1BA9B400c9CA",
// 		"0xb393Edd0Bf261C49C8A5C3fb21B23DE8753BE967",
// 	}

// 	transId := big.NewInt(15)
// 	for _, address := range addressList {
// 		isConfirm, err := platformWalletSC.Confirmations(&bind.CallOpts{}, transId, common.HexToAddress(address))
// 		if err != nil {
// 			fmt.Printf("%s error: %s\r\n", err.Error(), address)
// 		} else {
// 			fmt.Printf("%s: %v\r\n", address, isConfirm)
// 		}
// 	}
// }

func TestPlatformWallet_SetInternalSmartContractAddress(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	internalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

	fmt.Println("InternalSmartContractConfig address: ", internalSCAddr.Hex())
	tx, err := platformWalletSC.SetInternalSmartContractConfig(transOpt, internalSCAddr)
	if err != nil {
		t.Fatal("ChangeRequirement error: ", err)
	}

	//addr2,err :=platformWalletSC.GetInternalSmartContractConfig(&bind.CallOpts{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(addr2)

	queryEthReceipt(t, &ethClient, tx)
}

func TestPlatformWallet_GetInternalSmartContractAddress(t *testing.T) {
	ethClient := initEthClient(t)
	callOpts := &bind.CallOpts{}
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	callOpts.From = common.HexToAddress(testOwnerAddr)
	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	INTCAddr, err := platformWalletSC.GetInternalSmartContractConfig(&bind.CallOpts{})
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	fmt.Println(INTCAddr.Hex())

}

func TestPlatformWallet_GetRequirements(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	required, err := platformWalletSC.Required(&bind.CallOpts{})
	if err != nil {
		t.Fatal("Required error: ", err)
	}
	fmt.Println("Required: ", required.String())
}

// func TestPlatformWallet_GetConfirmationCount(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
// 	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
// 	if err != nil {
// 		t.Fatal("NewPlatform error: ", err)
// 	}
// 	callOpts := &bind.CallOpts{}
// 	callOpts.From = common.HexToAddress(testOwnerAddr)

// 	INTCAddr, err := platformWalletSC.GetConfirmationCount(&bind.CallOpts{}, big.NewInt(122))
// 	if err != nil {
// 		t.Fatal("NewPlatform error: ", err)
// 	}
// 	fmt.Println(INTCAddr)

// }

func TestPlatformWallet_SendMintRequest(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewPlatform error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	var dest [32]byte

	tx, err := platformWalletSC.SubmitMintRequest(transOpt, addr, "USDT", big.NewInt(1000), dest)
	if err != nil {
		fmt.Println(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

// func TestConfirmMintRequest(t *testing.T) {
// 	ethClient := initEthClient(t)
// 	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
// 	platformWalletSC, err := contract.NewPlatformWallet(addr, ethClient.Client)
// 	if err != nil {
// 		t.Fatal("NewPlatform error: ", err)
// 	}
// 	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
// 	if err != nil {
// 		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
// 	}

// 	tx, err := platformWalletSC.ConfirmTransaction(transOpt, big.NewInt(63))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	queryEthReceipt(t, &ethClient, tx)

// 	req, err := platformWalletSC.Required(&bind.CallOpts{})
// 	if err != nil {
// 		t.Fatal("Required encounter error: ", err)
// 	} else {
// 		fmt.Println("Platform wallet required count: ", req.String())
// 	}

// }

func TestPlatformWallet_SignTransaction(t *testing.T) {
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<PlatformWallet>")
	ethClient := initEthClient(t)
	platformWallet, err := contract.NewPlatformWallet(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	transOpt, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	transOpt.NoSend = true
	// transOpt.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// 	return tx, nil
	// }

	tx, err := platformWallet.TransferETH(transOpt, common.HexToAddress("0"), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	txBin, err := tx.MarshalBinary()
	hexStr := base64.StdEncoding.EncodeToString(txBin)
	fmt.Println(hexStr)
	v, r, s := tx.RawSignatureValues()

	fmt.Printf("v: %s, r: %s, s:%s\r\n", v, r, s)

	sign, err := crypto.GenerateRSASignFromBase64("MIICXQIBAAKBgQCZqrBrEPT/YCnDOQ28j9pWuQotF8Ag1ymLu4TR2lMziZI4XtAeDaZIHVodbrYo4ps8EJwIdHp/v/ZLsVDqi4+OElRQAPWi6luzHjx2asdrU2kfaPTK5rwz4zZMzntUXpsPqOV73b7KVRYz8eUo6hnh+UJJtEikMk0iTk6xmjwLUwIDAQABAoGAN1cjZcubkATfxXDco5Xi+ex137t389vJGIuVz8WixSK8SBTNOqWLxtjbRVJDxAGoCj+pEdpw62UEnEnlbDKKpf7SxCMi3Cb/yvtucpfrJpldLUXdXsqkY2YjT7yJRYXklloxi2zCQzxvwgqUQeGMpfZtQhi3IXa7Bd65laN8nqECQQDEzH3/daA5hL5QMhh1ekfXERRZntXC75abr3MgvpHS2NgP2QaiKUQ6eXdKokaGCWnnmvnj8b3SIc0qdV/jkit7AkEAx+SXWU2ESWhde3uayBiBbs9TlxsTUJUqmqWcybBURHPLRp24C0x2LHRynHoAX3AefJDEWW0nsezZmcdEQJBMCQJBAKspC4CuJf9Ao2EOYNVz152GfkN/8HyNcljPXHsjI6LU8/28jJdm+q88y3K+9kVHVLOZxzLoImbq/Qyrbw13KJcCQFTynSsTWTOFCa0vYyDD4UWBECn4FKY7LgbYeJ/xsm4As5NH9W2/ybyspARBiKUGVb5kGz1RuPvRBsxmNWgmqlkCQQCvQ1baE7zUiP4xruz9dazh0YkuI4F8DDXiqS5H+sE3PLigf/rDjy2VlVDsvDUbzkuhfvuRslGYunII0wCbF2uB", string(txBin))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("signature:")
	fmt.Println(sign)
}

func TestInitPlatformWallet(t *testing.T) {
	TestPlatformWallet_SetWalletOwner(t)
	TestPlatformWallet_AddWalletOperator(t)
	// TestPlatformWallet_SetInternalSmartContractAddress(t)
}
