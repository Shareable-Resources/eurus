package crypto

import (
	go_crypto "crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/ethereum"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

//var key string = "JdVa0oOqQAr0ZMdtcTwHrb7N4SIhpKqXA2vhvCGl6jA="
var key string = "JdVa0oOqQAr0ZMdtcTwHrb7N4SIhpKqXA2vhvCGl6jA="

/*
func TestEncrypt(t *testing.T) {

	encodedStr := base64.StdEncoding.EncodeToString([]byte("HelloW"))
	en, err := crypto.EncryptFromBase64(encodedStr, key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(en)
	}

	de, err := crypto.DecryptFromBase64(en, key)
	if err != nil {
		fmt.Println(err)
	} else {
		str, _ := base64.StdEncoding.DecodeString(de)
		fmt.Println(string(str))
	}

}
*/
func TestEncrypt2(t *testing.T) {
	pubKey := "MIGJAoGBAO6w9z2oQZiuQT95gAq65HergUVAadDq5sg8CfOKXof8VOVToBo0f5HP3GKzlI1Cru35i1obs7aZQ2DrM5mdSD7/2OY4mO8WjtNbP7IYfkR4u5WxmIJQ5au3C1mDK+498bNZqBWeZwNjXKouKTTfDwwMc5k0toaLr3fIdlHP2oiRAgMBAAE="
	priKey := "MIICXAIBAAKBgQDusPc9qEGYrkE/eYAKuuR3q4FFQGnQ6ubIPAnzil6H/FTlU6AaNH+Rz9xis5SNQq7t+YtaG7O2mUNg6zOZnUg+/9jmOJjvFo7TWz+yGH5EeLuVsZiCUOWrtwtZgyvuPfGzWagVnmcDY1yqLik03w8MDHOZNLaGi693yHZRz9qIkQIDAQABAoGAMl657hMBtLyhHEoBkUIbUH2qy/hp3CKWDQ9Ockxy4nOHXtWk5aLKgPTCZznKUX0O+T0+AQfzhscVBvDbdMFSK0Dx71B8LdnmLkvyUGu9nVTTuGgelTF0WQCf86Rd6LknGwZQAK6RNzeQZEq6XG/yJOuHWIPFsy9YHY2+yNyQkVkCQQD3IsLc7OmsfuB/UvKXsUsCoajhqOG7poh3tAY+YJMlWyToDGOXw1JSCTph6XbzYmnk/kRyhioxY78r2R5IatZ3AkEA90CqcJkiXP5JOgVSvAiYa9O14EfvXOjCYfAklFLekyru5OjUJhJ2YACykB1m4O1w7BvZwRMFvci5bca6QqlzNwJAI45fxNNdJ1E10XvIpWR/q0hA+P6IQ6xJFBfVkiHo6cX8QFqP9aTHckAbozyovYmaPLUMegGtjl+QgKmDPt4ILwJACdPzMmiT2hhtdrXxdPHuhRK0Pwb897d0yonOGms016q0Njse+6huNiCw+FOC3Fvzyh7NSARmjQWmgTuN+cpcfwJBAM58p/u2zv1KDezdD4ynn+GGvhPEvbdW7b0uXnQIZYRaTUzPVLMyteTW8NskqS97mAw4uD7RpDeUH6i0qzY2zxc="

	password := "admin123456"
	encryptedPassword, err := crypto.EncryptRAFormat([]byte(password), pubKey, go_crypto.SHA1)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		fmt.Println(encryptedPassword)

	}
	decryptedPassword, err := crypto.DecryptRAFormat(encryptedPassword, priKey)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		str, _ := base64.StdEncoding.DecodeString(decryptedPassword)
		fmt.Println(string(str))
	}
}

func TestDecryptAES(t *testing.T) {
	aesKey := "SP3+mcXwtxGu80UA8g0FEVbo6VmbeodK4z6Xq6waIrg="
	data := "kEPtz2HKXaDQBlH84yn8ndZAjDrJScGimD+NzoitcGw="

	val, _ := crypto.DecryptAESFromBase64(data, aesKey)
	v, _ := base64.StdEncoding.DecodeString(val)
	fmt.Println(string(v))

}

func TestDecrypt(t *testing.T) {
	//privateKey := "MIICXAIBAAKBgQDusPc9qEGYrkE/eYAKuuR3q4FFQGnQ6ubIPAnzil6H/FTlU6AaNH+Rz9xis5SNQq7t+YtaG7O2mUNg6zOZnUg+/9jmOJjvFo7TWz+yGH5EeLuVsZiCUOWrtwtZgyvuPfGzWagVnmcDY1yqLik03w8MDHOZNLaGi693yHZRz9qIkQIDAQABAoGAMl657hMBtLyhHEoBkUIbUH2qy/hp3CKWDQ9Ockxy4nOHXtWk5aLKgPTCZznKUX0O+T0+AQfzhscVBvDbdMFSK0Dx71B8LdnmLkvyUGu9nVTTuGgelTF0WQCf86Rd6LknGwZQAK6RNzeQZEq6XG/yJOuHWIPFsy9YHY2+yNyQkVkCQQD3IsLc7OmsfuB/UvKXsUsCoajhqOG7poh3tAY+YJMlWyToDGOXw1JSCTph6XbzYmnk/kRyhioxY78r2R5IatZ3AkEA90CqcJkiXP5JOgVSvAiYa9O14EfvXOjCYfAklFLekyru5OjUJhJ2YACykB1m4O1w7BvZwRMFvci5bca6QqlzNwJAI45fxNNdJ1E10XvIpWR/q0hA+P6IQ6xJFBfVkiHo6cX8QFqP9aTHckAbozyovYmaPLUMegGtjl+QgKmDPt4ILwJACdPzMmiT2hhtdrXxdPHuhRK0Pwb897d0yonOGms016q0Njse+6huNiCw+FOC3Fvzyh7NSARmjQWmgTuN+cpcfwJBAM58p/u2zv1KDezdD4ynn+GGvhPEvbdW7b0uXnQIZYRaTUzPVLMyteTW8NskqS97mAw4uD7RpDeUH6i0qzY2zxc="
	//privateKey := "MIICXQIBAAKBgQDU5YufylzObiWijgdmBfeZAyKrSOxq6Nrh+5Oh4crA/QQPkp8ZAcOoRvApgLLqQuUbfd7egWMOwczLID/yrA0Wi3k+Tk9Z7z4SfrPxaWu1+elHBktzLhQLkf4bj2PVAJH1zePJIy4PJCIO8gFelstEHFvWcso2ePA1nnZ6knunJQIDAQABAoGAYtUNNGjlHI/VuNjmZl5uywHBnnKEDj17H12C86u2TFEpCXGvmhRPmFcWNq4gYNAdO937EKBQNBGT2Nhn12g3yl4+4edPkf39URg3ZwCq4uxAWQ9z5rwPa1eOOAok2VNrJd1NX4WEUerXRgzT1MU449jaZkGz4m8LEiKYKEIW+CECQQDZ2hQy80HXM/gJVIaYN2khzW7Q5LigrVzGf1SJbXGb10XfSbm1r0rN6DSa2haSfChb/lj8eN1xjXMaxiZTsRfbAkEA+i1VRWtdHNMdTnXW/gU3cfbGxo7/nj7WUU4LaZ2AvZ12K/dRkwZYPEyxlId2HnJwnw8TRX8FIR8H3mkV5HXs/wJAcmOcL5Sjgch7+Qo1EkAmJ+WixnUSrOvaxy+cx/x7pwTGX5RquwesE6pV1Omm6Ivg9Uz8lLUyManAQtLA1Tkr+QJBAI5BzOUmgdHsMhP1agUTzk1dd/ZcRfoj3RZqfI7X4ubvbMzfW2FxECdprOi6hm4VwPiRR/ISokYNMRpFQw+gBt0CQQCzxWkD/jeGCI9F4qYI9qwgT8lxce702Y2D9huaQ3sxc9NiD/Bhm1Nw8P8W2tgXENh9BV9xcq/gEsHu7OaidSfF"
	privateKey := "MIICXQIBAAKBgQC8TPGDk1ERv3Cy75l1h2Ymu3HzcozNohxr3msxkmTlIVUIyLYS/3Wt23/Cspxa/jmZbITrlH3ySNAfgcgxWpL6etWvgsvWuAQAl9Y9XouVM7SUbzsWPXwjeNEuRrX471fi1RMlXmma3i+3ef/g3MwGKuPQg2vmJDkPcYNWvwqMtwIDAQABAoGARFm1rSoTWlrgYMNiM+xDIzoVQeZJv/bhVDJ3PJZEM3vm7vpURQ8fu00a0iSdZtYzOELQ+OZ81S6jBVzJu0iCYvPoDyVeQECUdFRr83jYZBhJa3OxhTukcIbaMIdiX9+/9GE2YlSlW8vwTl4I3nmPM0gjCuClGjNq0GpaDqFdLJkCQQDqgwLvufSZHP/RH+Z4u6lOUYcE+14xOwbc6nvzJP8SZ9XaV11XzRYW5zMWNBFS5abhiIOLTSB1EwJTNGhLKoJDAkEAzY3yWJP6/elUFKGzUVcgdSQqk64b4OPSB7txYv+/9Iy+hyYxmrjhcMs1pVnt7IUMh1lPRobuqFNpbUT8/JEmfQJBAJAn5MA0URVl7CNDNst0oEiSeSso0ef1Xuyh/3uhfGJ2DdFf7vMwEKoXiSHW8+c/YGw3HJWe3ZraCo3XI11BTj0CQQDBXPBG5x3K0T1pExuuPrmnuFfdCIHtMiKXZ3qoWFFoaSWp2QH8vEJGasKSQ2CPiDFa+Wabv76t1ZltUuJxbNG5AkB+4NwalsCkie2ina7fHzxeLce1DC9udvQJa884LnPmp0A/DBLt2CPKfevxF+6AlPFPrmPffmYWgRzIyZUHk2IQ"
	//encryptedPassword := "rLoBAKYeRNgjNiJ+4r3uoElscuEQ/MAjESM+amYIql3UWUueBOL84fgaVgdKopjyLLYg209PzCcUmbvQuA6u4JIa9QNCKxnZO57hdMIo+mgLcyyHywDW12gvG5oFzLmdOi2zks11mWN2ZXTsarfctTapDvqJK+sx+t4uWI9QV2GU5/d/UMfjhxM4Rp8yqYhpfE+SnHXaMLHLkZwQWKc0LXDlw9E="
	//encryptedPassword := "rLoCAAHTKqoY+uu4OwXHw9ybrJ/CG2vYQj3mWEzFqmhON8thS6lEroaOXZH0NNnGRx2bAOJU7SLxBxfyg/aKULgtaseOUMKpQ6GZFdmHgTs1kUapiwqLBcDjhQeIWXGnVa4rg9W60eKn09D+X5TrMrl5AXi7Bb+1IFD3OrvPm62zygbQLI45jBce/Hhmvq7Fc2fQroOgDxDUs/YaVLr4Y2pY4SrK58cXYKi266DOKbI3aaXY11A9oCt8aWf9scd2dHEenXWplocu2m60Q4NLfFbu3Q/HYWDYeYk2ZUFphiJATFqc"
	encryptedPassword := "rLoCACdBICvL1VvHctM6IonrGxr8Dhz2TJnlSLF4xv/oIfaFEgrN+t9cb2hFemx7SxouHCM/Taw70OKuoKg9V7+CCH39MH8grLiA1R/Oz8ODA7tbrOHhCnRfAoBz9E6LOoWNqwjF66ME5FPOJaPraO557MpX7Derx4knDhfeWP4nr6W3beRBgWIeN3Af04s6I3uGFiK1xEZeYeuJcmdHHNmyjiPpdae0mBBYB99U/AFy1zOCUioKMkWAzzXMyPttNt0FZaaDbL/rqry+fzQmUXbEXZQtaSD5U6BEsyg9os+z4R7h"
	decryptedPassword, err := crypto.DecryptRAFormat(encryptedPassword, privateKey)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		str, _ := base64.StdEncoding.DecodeString(decryptedPassword)
		fmt.Println(string(str))
	}
}

func TestRSA(t *testing.T) {
	_, base64PriKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		t.Errorf(err.Error())
	}

	privateKeyDer, err := base64.StdEncoding.DecodeString(base64PriKey)
	if err != nil {
		t.Errorf(err.Error())
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyDer)
	if err != nil {
		t.Errorf(err.Error())
	}
	pubKeyBase64 := crypto.GeneratePublicKey(privateKey)

	sign, _ := crypto.GenerateRSASignFromBase64(base64PriKey, "HelloWorld")

	fmt.Println("Private key: ", base64PriKey)
	fmt.Println("Public key: ", pubKeyBase64)
	fmt.Println("Signature: ", sign)
}

//func TestRSASign(t *testing.T) {
//	privateKeyBase64 := "MIICXQIBAAKBgQCySctqeJcF7UZO80uKlHlqMVLYEc3TRzMDlEO/D0WzADvzt6YbxdqbRoHg5jArb1aWHowH7KPRIAZmzWXU9UNgnzMYtFBtrg1iwn6vAcVi4ML3e5YiDCU+3kfC/zlpoppmPJK+/chvcxX/8Bn2TWLp/WFxwr2QdBafWd00jtXOvwIDAQABAoGAHiiZfkMz2M7QS96f8qs29p2tU9t/I7QxjvNvmkwtECCEA7AXKy+G0SlWIwMNNuu4i+BS86w18dQnARDRnvp2EQPUhPcLOSWE5ErwDV5rr6wzE0c6g5PR/NQawtnfzGE13Sei0jVAOM59JsAYOBvjFrayP3wwy3aTcnGqi3y3gdECQQDTC/s4zeBpxwuAcY3vE4Qvj1DnLlbx3Ey2IgxvnNng1yCoJs7xqk+mc7GulKlAWO8/fJjDOutpBFzsXjZ1cPO9AkEA2EOL34p9u0qAdVqqahjzhiQBmiMZPSoIma2QxHpJcONoXj93OjMjBA5lCtDmbUJy9vSjddQHka2B/3ZIhO82KwJANpzQGqfz3Yt7Z3Z7ExXKy//bIKe2Kgom866l5zAAlZU1xy2buZ+n0tv3a8Tec19QMe6NHWAe0OMnJKdQkl1zRQJBAI98MV99aXJMOfKi1xJRv7EuTonP47z7GH2t9LEHWZa2qtZtUVOh/geA7XVBqMCxJ8nGAnX42ZTEA85U0kSw2cUCQQC+0erKC4kZv+GxfME1oAWy6lTsEUtuix8UMC2aTu4p5xgyYwPZbwQn7xqLjOYmT7bP5Qw1xevAFlA7KJQoK4qU"
//	content := "{\"methodName\":\"authenticate\",\"nonce\":\"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\"timestamp\":12345677,\"data\":{\"serviceId\":3}}"
//	sign, err := crypto.GenerateRSASign(privateKeyBase64, content)
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//
//	fmt.Println("Sign: ", sign)
//}

func TestVerifyECDSA(t *testing.T) {
	const exampleMnemonic string = "install future tornado nerve belt bag impact foam gift melody prison swing protect afraid talent drive quote change final myself black define tongue stock"
	wallet, account, address, err := crypto.GenerateHdWalletKey(exampleMnemonic, "Hello", "0", "4")
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("address: ", address)
	priKeyBytes, err := wallet.PrivateKeyBytes(*account)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("PrivateKeyHex: ", hex.EncodeToString(priKeyBytes))
	priKey, _ := wallet.PrivateKey(*account)

	publicKeyECDSA, ok := priKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	compressedPubKey := eth_crypto.CompressPubkey(publicKeyECDSA)
	fmt.Println("Compressed public key: ", hex.EncodeToString(compressedPubKey))

	// output looks the same here: &{0xc0000aa810 111467676592638012114394423922557735151709477236068497041031680458777853727687 91044589275583075592054080804256995961966216084890183147914377743969771213580}

	publicKeyBytes := eth_crypto.FromECDSAPub(publicKeyECDSA)

	pubKeyHex := hex.EncodeToString(publicKeyBytes)
	fmt.Println("publicKeyHex: ", pubKeyHex)
	nowTime := time.Now().UnixNano()

	var sentense string = "deviceId=3&timestamp=" + strconv.FormatInt(nowTime, 10) + "&walletAddress=" + address
	fmt.Println("sentence: ", sentense)

	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write([]byte(sentense))
	if err != nil {
		t.Error(err.Error())
	}

	var sha3Hash []byte
	sha3Hash = hasher.Sum(sha3Hash)

	sign, err := secp256k1.Sign(sha3Hash, priKeyBytes)
	if err != nil {
		t.Error(err.Error())
	}
	sign = sign[:len(sign)-1]
	signHex := hex.EncodeToString(sign)
	signByte := []byte(signHex)
	signDecodedByte := make([]byte, hex.DecodedLen(len(signByte)+1))
	_, err = hex.Decode(signDecodedByte, []byte(signHex))
	if err != nil {
		t.Error(err.Error())
	}
	isValid, err := crypto.VerifyECDSASignatureByHexKey(pubKeyHex, []byte(sentense), signDecodedByte)
	if err != nil {
		t.Error(err.Error())
	}
	if !isValid {
		t.Error("Not valid signature")
	} else {
		fmt.Println("sign: ", signHex)
		fmt.Println("sign: ", hex.EncodeToString(sign))
	}
}

func TestCreateImportWalletSignature(t *testing.T) {
	var err error
	var walletAddress common.Address
	var deviceId = ""
	var privateKey string = "fa54e1ea4f97784c8b7747290d67484bcbaee1a6635d1f52f843dde9eae943d0"
	ethKeyPair, _ := ethereum.GetEthKeyPair(privateKey)

	publicKeyECDSA, _ := ethKeyPair.PrivateKey.Public().(*ecdsa.PublicKey)
	walletAddress = eth_crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("Wallet address: ", ethereum.ToLowerAddressString(walletAddress.Hex()))

	compressedPubKey := eth_crypto.CompressPubkey(publicKeyECDSA)
	fmt.Println("public key compressed: ", hex.EncodeToString(compressedPubKey))
	unixNow := time.Now().Unix()
	var sentense string = "deviceId=" + deviceId + "&timestamp=" + strconv.FormatInt(unixNow, 10) + "&walletAddress=" + ethereum.ToLowerAddressString(walletAddress.Hex())[2:]

	fmt.Println("sentence: ", sentense)
	fmt.Println("timestamp: ", unixNow)
	fmt.Println("deviceId: ", deviceId)
	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write([]byte(sentense))
	if err != nil {
		t.Error(err.Error())
	}

	var sha3Hash []byte
	sha3Hash = hasher.Sum(sha3Hash)
	testPrivKey := privateKey
	p, _ := ethereum.GetEthKeyPair(testPrivKey)

	sign, err := secp256k1.Sign(sha3Hash, eth_crypto.FromECDSA(p.PrivateKey))
	if err != nil {
		t.Error(err.Error())
	}
	sign = sign[:len(sign)-1]
	isValid, err := crypto.VerifyECDSASignatureByHexKey(hex.EncodeToString(compressedPubKey), []byte(sentense), sign)
	if err != nil {
		t.Error(err.Error())
	}
	if !isValid {
		t.Error("Not valid signature")
	} else {
		fmt.Println("sign: ", hex.EncodeToString(sign))
	}

	fmt.Println("API: /user/importWallet")
	var req map[string]interface{} = make(map[string]interface{})
	req["nonce"] = uuid.New().String()
	req["deviceId"] = deviceId
	req["timestamp"] = unixNow
	req["walletAddress"] = ethereum.ToLowerAddressString(walletAddress.Hex())[2:]
	req["publicKey"] = hex.EncodeToString(compressedPubKey)
	req["sign"] = hex.EncodeToString(sign)

	reqJson, _ := json.MarshalIndent(req, "", "\t")
	fmt.Println(string(reqJson))

}

func TestPubKeyStringToAddress(t *testing.T) {

	pubBytes, err := hex.DecodeString("03bbba49a934014049d99a7f5c809fd0da59b1cb47dcbd0b3fe097adc5eaa5ec42")
	if err != nil {
		fmt.Println(err, "--1")
	}
	x, y := secp256k1.DecompressPubkey(pubBytes)
	pubkey := elliptic.Marshal(secp256k1.S256(), x, y)

	fmt.Println(hex.EncodeToString(pubkey))

	pubStr, err := crypto.DecompressPubKey("03bbba49a934014049d99a7f5c809fd0da59b1cb47dcbd0b3fe097adc5eaa5ec42")

	fmt.Println(pubStr)

	addr, err := crypto.PubKeyStringToAddress(hex.EncodeToString(pubkey))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(addr.Hex())
}

func TestDecryptByPassword(t *testing.T) {
	encrypted := `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`
	encoded, err := crypto.DecryptByPassword("1234", encrypted)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Println(string(data))
}

func TestHashBCrypt(t *testing.T) {

	text := "abcd1234"
	salt := "nimdarepus"
	hash, err := crypto.GenerateBCryptHash([]byte(text), []byte(salt))
	if err != nil {
		t.Fatal(t)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(hash))
}
