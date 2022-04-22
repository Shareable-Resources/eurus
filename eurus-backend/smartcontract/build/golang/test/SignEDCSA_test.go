package test

import (
	_ethereum "eurus-backend/foundation/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	//"encoding/hex"
	"fmt"

	"math/big"
	"testing"
)

//func TestSign(t *testing.T) {
//	ethClient := initEthClient(t)
//	ethClient2 := initTestNetEthClient(t)
//	amount := big.NewInt(0)
//	amount, _ = amount.SetString("0", 10)
//
//	ethKeyPair, err := _ethereum.GetEthKeyPair("e11303b135906488e771047814dbcb9b3794149c3a1c1cd7f6d6064f07515b94")
//	if err != nil {
//		fmt.Println("ETH gen key err : ", err)
//	}
//	fmt.Println("gggg : " , ethKeyPair.Address.Hex())
//
//
//	nonce, err := ethKeyPair.GetNonce(&ethClient)
//	if err != nil {
//		fmt.Println("Get nonce err : ", err)
//	}
//
//	gasPrice, err := ethClient.GetGasPrice()
//	if err != nil {
//		fmt.Println("gas price err : ", err)
//	}
//
//	toAddress := common.HexToAddress(testTestNetPlatformWalletOwnerAddr)
//	const ETHTransferStandardGasLimit uint64 = 21000
//
//	tx := types.NewTransaction(nonce, toAddress, amount, ETHTransferStandardGasLimit, gasPrice, nil)
//	s := types.HomesteadSigner{}
//	hash := s.Hash(tx)
//	sig, err := crypto.Sign(hash.Bytes(), ethKeyPair.PrivateKey)
//	if err != nil {
//		fmt.Println("Sign err : ", err)
//	}
//	fmt.Println("ADDR : ", ethKeyPair.Address.Hex())
//
//
//	addr := getAddressBySmartContractName("SignatureVerifier")
//	platformWallet, err := contract.NewSignatureVerifier(addr, ethClient.Client)
//
//
//	if err != nil {
//		fmt.Println("Sign verifier err : ", err)
//	}
//
//
//	oriaddr, err := platformWallet.RecoverSigner(&bind.CallOpts{}, hash, sig)
//	if err != nil {
//		fmt.Println("Recovery err : ", err)
//	}
//	//
//	isValid, _ := platformWallet.VerifySignature(&bind.CallOpts{}, hash, sig, oriaddr)
//
//
//	addr1 := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
//	eurusPlat ,_ := contract.NewEurusPlatformWallet(addr1,ethClient2.Client)
//	//isWriter,_ := eurusPlat.IsWriter(&bind.CallOpts{},oriaddr)
//	//fmt.Println("isWriter : ", isWriter)
//
//	getList,_ := eurusPlat.GetWriterList(&bind.CallOpts{})
//	for _,e := range(getList) {
//		fmt.Println("Writer is : ", e.Hex())
//	}
//
//
//	fmt.Println("isValid : ", isValid)
//	fmt.Println("Signed addr : " , oriaddr.Hex())
//
//}

//func TestExtractSignData(t *testing.T) {
//	ggtt := make([]byte, 32)
//
//	_, err := rand.Read(ggtt)
//	if err != nil {
//		// handle error here
//	}
//	fmt.Println(ggtt)
//	id := 2000
//	hashData := common.Hash{byte(id)}
//	byteSign := make([]byte, 32)
//
//
//	copy(byteSign[:32],hashData.Bytes())
//	fmt.Println(byteSign)
//
//	key,err := _ethereum.GetEthKeyPair(testOwnerPrivateKey)
//	if err != nil {
//		fmt.Println(err)
//	}
//	signData, err := secp256k1.Sign(byteSign, crypto.FromECDSA(key.PrivateKey))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	publicKeyECDSA, _ := key.PrivateKey.Public().(*ecdsa.PublicKey)
//	pubKey := crypto.FromECDSAPub(publicKeyECDSA)
//
//	fmt.Println("Signdata : ", signData)
//
//
//	fmt.Println(hashData.Bytes())
//
//
//	unSign := secp256k1.VerifySignature(pubKey,hashData.Bytes(),signData)
//	fmt.Println("Sign Result : " , unSign)
//
//}
//
//

func TestCleanPending(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	amount := big.NewInt(10000000)
	amount, _ = amount.SetString("1000000000", 10)

	fmt.Println("Amount : ", amount)

	ethKeyPair, err := _ethereum.GetEthKeyPair("e11303b135906488e771047814dbcb9b3794149c3a1c1cd7f6d6064f07515b94")
	if err != nil {
		fmt.Println("ETH gen key err : ", err)
	}

	//make transaction
	nonce, err := ethKeyPair.GetNonce(&ethClient)
	if err != nil {
		fmt.Println("get nonce err : ", err)
	}
	fmt.Println(nonce)
	gasPrice, err := ethClient.GetGasPrice()
	plex := big.NewInt(50)
	gasPrice.Mul(gasPrice, plex)
	fmt.Println("Gas price : ", gasPrice)

	if err != nil {
		fmt.Println("gas price err : ", err)
	}

	toAddress := common.HexToAddress(testTestNetPlatformWalletOwnerAddr)
	const ETHTransferStandardGasLimit uint64 = 10000000

	tx := types.NewTransaction(47, toAddress, amount, ETHTransferStandardGasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(4)), ethKeyPair.PrivateKey)
	if err != nil {
		fmt.Println("err :", err)
	}

	v, r, s := signedTx.RawSignatureValues()
	fmt.Printf("v: %s, r: %s, s:%s\r\n", v, r, s)
	// err = ethClient.SendTransaction(signedTx)
	// if err != nil {
	// 	fmt.Println("err :" , err)
	// }
	// fmt.Println("HASHHHH : ", signedTx.Hash().Hex())
	// queryEthReceipt(t, &ethClient, signedTx)
}
