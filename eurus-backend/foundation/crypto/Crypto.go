package crypto

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	go_ethereum_crypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts"
	secp256k1_ "github.com/ethereum/go-ethereum/crypto/secp256k1"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var rsaKeyBitSize int = 1024

var raFormatHeaderSHA256 = []byte{0xAC, 0xBA, 0x1, 0x0}
var raFormatHeaderSHA1 = []byte{0xAC, 0xBA, 0x2, 0x0}

// generatePrivateKey creates a RSA Private Key of specified byte size
// converts an RSA private key to PKCS #1, ASN.1 DER form
func GeneratePrivateKey() (*rsa.PrivateKey, string, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBitSize)
	if err != nil {
		return nil, "", err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, "", err
	}

	byteKey := x509.MarshalPKCS1PrivateKey(privateKey)
	base64Key := base64.StdEncoding.EncodeToString(byteKey)
	return privateKey, base64Key, err
}

func GetRSAPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	privateBytetKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(privateBytetKey)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}

func GetWalletAddressFromMnemonic(mnemonic string) {

}

// converts an RSA public key to PKCS #1, ASN.1 DER form.
func GeneratePublicKey(privatekey *rsa.PrivateKey) string {
	byteKey := x509.MarshalPKCS1PublicKey(&privatekey.PublicKey)
	base64Key := base64.StdEncoding.EncodeToString(byteKey)

	return base64Key
}

func EncryptAESFromBase64(base64Data string, base64Key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	return EncryptAES(data, key)
}

func EncryptAES(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	data = aesPad(data)
	dataSize := len(data)

	buffer := make([]byte, aes.BlockSize+dataSize)
	_, err = io.ReadFull(rand.Reader, buffer[:aes.BlockSize])
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCEncrypter(block, buffer[:aes.BlockSize])
	blockMode.CryptBlocks(buffer[aes.BlockSize:], data)

	output := base64.StdEncoding.EncodeToString(buffer)

	return output, nil
}

func DecryptAESFromBase64(base64Data string, base64Key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	return DecryptAES(data, key)
}

func DecryptAES(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("data too short")
	}

	iv := data[:aes.BlockSize]

	if len(data)%aes.BlockSize != 0 {
		return "", errors.New("Invalid data size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	outputData := make([]byte, len(data)-aes.BlockSize)

	mode.CryptBlocks(outputData, data[aes.BlockSize:])
	outputData, err = aesUnpad(outputData)
	if err != nil {
		return "", err
	}
	output := base64.StdEncoding.EncodeToString(outputData)
	return output, nil
}

func aesPad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func aesUnpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func GenerateRSASignFromBase64(privateBase64Key string, message string) (string, error) {
	privateBytetKey, err := base64.StdEncoding.DecodeString(privateBase64Key)
	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBytetKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	return GenerateRSASign(privateKey, []byte(message))
}

func GenerateRSASignWithRawData(privateBase64Key string, message []byte) (string, error) {
	privateBytetKey, err := base64.StdEncoding.DecodeString(privateBase64Key)
	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBytetKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	return GenerateRSASign(privateKey, message)
}

/// RSA Signature
func GenerateRSASign(privateKey *rsa.PrivateKey, msg []byte) (string, error) {

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

/// Verify RSA Signature
func VerifyRSASignFromBase64(publicBase64Key string, message string, sign string) (bool, error) {

	publicByteKey, err := base64.StdEncoding.DecodeString(publicBase64Key)
	if err != nil {
		return false, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicByteKey)
	if err != nil {
		return false, err
	}
	return VerifyRSASign(publicKey, message, sign)
}

/// Verify RSA Signature
func VerifyRSASign(publicKey *rsa.PublicKey, message string, sign string) (bool, error) {

	msg := []byte(message)
	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return false, err
	}
	msgHashSum := msgHash.Sum(nil)

	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum, signByte, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func EncryptRSAFromBase64Key(publicBase64Key string, hashType crypto.Hash, message string) ([]byte, error) {

	publicByteKey, err := base64.StdEncoding.DecodeString(publicBase64Key)
	if err != nil {
		return []byte{}, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicByteKey)
	if err != nil {
		return []byte{}, err
	}
	return EncryptRSA(publicKey, hashType, message)
}

// the output : encrypted message is base64 ed
func EncryptRSA(publicKey *rsa.PublicKey, hashType crypto.Hash, message string) ([]byte, error) {
	var hashObj hash.Hash = nil
	if hashType == crypto.SHA256 {
		hashObj = sha256.New()
	} else if hashType == crypto.SHA1 {
		hashObj = sha1.New()
	} else {
		return nil, errors.New("Hash type not supported")
	}
	encryptedBytes, err := rsa.EncryptOAEP(
		hashObj,
		rand.Reader,
		publicKey,
		[]byte(message),
		nil)
	if err != nil {
		return []byte{}, err
	}

	return encryptedBytes, nil
}

func DecryptRSAFromBase64Key(privateBase64Key string, hashType crypto.Hash, encryptedBytes []byte) ([]byte, error) {
	privateBytetKey, err := base64.StdEncoding.DecodeString(privateBase64Key)
	if err != nil {
		return []byte{}, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBytetKey)
	if err != nil {
		return []byte{}, err
	}

	return DecryptRSAWithHashType(privateKey, encryptedBytes, hashType)
}

func DecryptRSAWithHashType(privateKey *rsa.PrivateKey, encryptedBytes []byte, hashType crypto.Hash) ([]byte, error) {
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: hashType})
	if err != nil {
		return []byte{}, err
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	//fmt.Println("decrypted message: ", string(decryptedBytes))
	return decryptedBytes, nil
}

func DecryptRSA(privateKey *rsa.PrivateKey, encryptedBytes []byte) ([]byte, error) {

	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return []byte{}, err
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	//fmt.Println("decrypted message: ", string(decryptedBytes))
	return decryptedBytes, nil
}

func GenerateWalletAddressFromMnemonic(mnemonic string) string {
	seed := bip39.NewSeed(mnemonic, "")

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {

	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {

	}
	return account.Address.Hex()
}

func GenerateHdWalletKey(mnemonic string, password string, walletIndex string, suffix string) (*hdwallet.Wallet, *accounts.Account, string, error) {
	if mnemonic == "" {
		return nil, nil, "", errors.New("Empty mnemonic phase")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, nil, "", err
	}
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, nil, "", err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/" + string(walletIndex) + "'/0/" + suffix)

	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, nil, "", err
	}

	addr, err := wallet.Address(account)
	if err != nil {
		return nil, nil, "", err
	}

	return wallet, &account, strings.ToLower(addr.Hex()), nil
}

func VerifyECDSASignatureByHexKey(hexPubKey string, rawData []byte, sign []byte) (bool, error) {
	pubKey, err := hex.DecodeString(hexPubKey)
	if err != nil {
		return false, err
	}
	return VerifyECDSASignatureByCompressedPubKey(pubKey, rawData, sign)
}

func VerifyECDSASignature(x, y *big.Int, rawData []byte, sign []byte) (bool, error) {
	pubKey := secp256k1_.CompressPubkey(x, y)
	return VerifyECDSASignatureByCompressedPubKey(pubKey, rawData, sign)
}

func VerifyECDSASignatureByCompressedPubKey(compressedPubKey []byte, rawData []byte, sign []byte) (bool, error) {

	hasher := sha3.NewLegacyKeccak256()
	_, err := hasher.Write(rawData)
	if err != nil {
		return false, err
	}

	var sha3Hash []byte
	sha3Hash = hasher.Sum(sha3Hash)

	isValid := secp256k1_.VerifySignature(compressedPubKey, sha3Hash, sign)
	return isValid, nil
}

func GenerateECDSSignature(hashData []byte, privateKeyHex string) ([]byte, error) {
	privateKey, err := go_ethereum_crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	return go_ethereum_crypto.Sign(hashData, privateKey)
}

// 1, encrypt the "cipher" by the AES key to produce an encrypted message byte array "encryptedCipher"
// 2, encrypt the AES key by RSA public key to produce an encrypted AES key byte array  "encryptedAESKey"
// 3, encrypted AES key byte array (encryptedAESKey) should be in fixed length 128  bytes
// 4, "encryptedAESKey" array and then append "cipherEncrypted" array to form byte array "mixedByte"
// 5, Base64("mixedByte") and write the string "messageToFile" to file
func EncryptRAFormat(message []byte, publicBase64Key string, hashType crypto.Hash) (string, error) {

	temp := make([]byte, 32)
	_, err := rand.Read(temp)
	if err != nil {
		return "", err
	}
	aesBase64Key := base64.StdEncoding.EncodeToString(temp)

	encryptedCipher, err := EncryptAES(message, temp)
	if err != nil {
		return "", err
	}
	cipherEncryptedByte, _ := base64.StdEncoding.DecodeString(encryptedCipher)

	encryptedAESKey, err := EncryptRSAFromBase64Key(publicBase64Key, hashType, aesBase64Key)
	if err != nil {
		return "", err
	}

	var mixByte []byte
	var header []byte
	switch hashType {
	case crypto.SHA256:
		header = raFormatHeaderSHA256

	case crypto.SHA1:
		header = raFormatHeaderSHA1

	}

	mixByte = append(mixByte, header...)
	mixByte = append(mixByte, encryptedAESKey...)
	mixByte = append(mixByte, cipherEncryptedByte...)

	return base64.StdEncoding.EncodeToString(mixByte), nil
}

// 1, Decode Base64 string "message" to procedure a byte array "mixedByte"
// 2, slice the first  128 bytes from array "mixedByte" to form byte array  "encryptedAESKey"
// 3, the remaining byte after slice is encrypted message byte array "encryptedCipher"
// 4. Decrypt byte array "encryptedAESKey" to produce AESKey
// 5. Use the AESKey to decrypt byte array "encryptedCipher" to produce the original cipher
func DecryptRAFormat(encryptedMessage string, privateBase64Key string) (string, error) {
	if encryptedMessage == "" {
		return "", nil
	}

	mixedByte, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	offset := 4
	var hashType crypto.Hash
	switch mixedByte[2] {
	case 1:
		hashType = crypto.SHA256
	case 2:
		hashType = crypto.SHA1
	default:
		return "", errors.New("Unsupported format")
	}
	encryptedAESKey := mixedByte[0+offset : 128+offset]
	encryptedCipher := mixedByte[128+offset:]

	aesKeyBytes, err := DecryptRSAFromBase64Key(privateBase64Key, hashType, encryptedAESKey)
	if err != nil {
		return "", err
	}

	aesKeyBase64 := string(aesKeyBytes)
	aesKey, err := base64.StdEncoding.DecodeString(aesKeyBase64)
	if err != nil {
		return "", err
	}
	originMessage, err := DecryptAES(encryptedCipher, aesKey)
	if err != nil {
		return "", err
	}
	return originMessage, nil
}

func Generate128BitsEntropyMnemonicPhrase() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func PubKeyStringToAddress(pubKeyString string) (*common.Address, error) {
	decompressedPubKey, _ := hex.DecodeString(pubKeyString)
	pubKey, err := go_ethereum_crypto.UnmarshalPubkey(decompressedPubKey)
	if err != nil {
		return nil, err
	}

	address := go_ethereum_crypto.PubkeyToAddress(*pubKey)
	return &address, nil
}

func DecompressPubKey(pubKeyString string) (string, error) {
	pubBytes, err := hex.DecodeString(pubKeyString)
	if err != nil {
		return "", err
	}
	x, y := secp256k1_.DecompressPubkey(pubBytes)
	if x == nil || y == nil {
		return "", errors.New("Invalid public key data")
	}
	pubkey := elliptic.Marshal(secp256k1_.S256(), x, y)

	return hex.EncodeToString(pubkey), nil
}

func EncryptByPassword(password string, data []byte) (string, error) {

	hasher := sha256.New()
	hasher.Write([]byte(password))
	aesKey := hasher.Sum(nil)

	return EncryptAES(data, aesKey)
}

func DecryptByPassword(password string, encrypted string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	aesKey := hasher.Sum(nil)

	dataEncrypted, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	return DecryptAES(dataEncrypted, aesKey)
}

func GenerateBCryptHash(plainContent []byte, salt []byte) ([]byte, error) {
	final := append(plainContent, salt...)
	return bcrypt.GenerateFromPassword(final, 10)
}

func VerifyBCryptHash(plainContent []byte, salt []byte, hash []byte) error {
	final := append(plainContent, salt...)
	return bcrypt.CompareHashAndPassword(hash, final)
}
