package main

import (
	go_crypto "crypto"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"eurus-backend/foundation/crypto"
	"eurus-backend/secret"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/tyler-smith/go-bip39"
)

func main() {
	argsWithProg := os.Args

	if len(argsWithProg) == 1 || argsWithProg[1] == "--help" {
		printUsage()
		return
	}
	if len(argsWithProg) < 2 {
		fmt.Println("too few Command-Line Arguments")
		printUsage()
		return
	}

	EncryptFromTerminal(argsWithProg)
}

func printUsage() {
	fmt.Println("usage: [mode] [-t|-i] [text|input filename] [-k] [AES key in Base64 format] -o [output file] ")

	fmt.Println("mode value: ")
	fmt.Println("aes: AES encryption - If -k parameter is not provided, the program will auto generate a AES key and corresponding 'encrypted' AES key to encrypt the text or input file")
	fmt.Println("rsa: RSA encryption")
	fmt.Println("pw: Password encryption -p is required")
	fmt.Println("mn: Mnemonic phase")
	fmt.Println("serverAuthKey: Generate server auth key pair. An encrypted base64 private key and a base64 public key")
	fmt.Println("-k : For aes mode only. Provide a given AES key in Base64 format")
	fmt.Println("-t : input the message you want to encrypt from terminal")
	fmt.Println("-i : input the message you want to encrypt from file")
	fmt.Println("-p : password to encrypt file")
	fmt.Println("-o : output file . If this argument is not given, the result will be send to standard output")
	fmt.Println("-w : mnemonic phase word count. Count must be 12, 15, 18, 21, 24")
	fmt.Println("Remarks: All output is in Base64 format")
}

// 1, encrypt the "cipher" by the AES key to produce an encrypted message byte array "cipherEncrypted"
// 2, encrypt the AES key by RSA public key to produce an encrypted AES key byte array  "AESKeyEncrypted"
// 3, encrypted AES key byte array (AESKeyEncrypted) should be in fixed length 128  bytes
// 4, "AESKeyEncrypted" array and then append "cipherEncrypted" array to form byte array "mixedByte"
// 5, Base64("mixedByte") and write the string "messageToFile" to file
func EncryptFromTerminal(splitText []string) {

	var cipher []byte
	mode := splitText[1]
	if mode != "aes" && mode != "rsa" && mode != "serverAuthKey" && mode != "pw" && mode != "mn" {
		fmt.Println("Invalid mode")
		printUsage()
		return
	}
	var aesKey string
	var outputFileName string
	var outputFile *os.File
	var password string
	var wordCount int

	for i := 2; i < len(splitText); i++ {
		switch splitText[i] {
		case "-t":
			cipher = []byte(splitText[i+1])
			i++
		case "-i":
			data, err := ioutil.ReadFile(splitText[i+1])
			if err != nil {
				fmt.Println("File reading error: ", err)
				return
			}
			cipher = data
			i++
		case "-k":
			aesKey = splitText[i+1]
			i++
		case "-o":
			outputFileName = splitText[i+1]
			i++
		case "-p":
			password = splitText[i+1]
			i++
		case "-w":
			var err error
			count := splitText[i+1]
			wordCount, err = strconv.Atoi(count)
			if err != nil {
				fmt.Println("Invalid word count")
				return
			}

			i++
		default:
			printUsage()
			return
		}
	}
	// put the base64ed message in the file
	var err1 error
	if outputFileName != "" {
		outputFile, err1 = os.Create(outputFileName)
		if err1 != nil {
			log.Fatal(err1)
		}
		defer outputFile.Close()
	}

	var outputMessage string
	var err error
	if mode == "rsa" {
		outputMessage, err = processRSA(cipher)
		if err != nil {
			return
		}
	} else if mode == "aes" {
		outputMessage, err = processAES(cipher, aesKey)
		if err != nil {
			return
		}
	} else if mode == "serverAuthKey" {
		outputMessage, err = processServerAuthKey()
		if err != nil {
			return
		}
	} else if mode == "pw" {
		if password == "" {
			printUsage()
			return
		}
		outputMessage, err = processPasswordEncrypt(cipher, password)
		if err != nil {
			return
		}

	} else if mode == "mn" {
		var mnemonicPhase string
		mnemonicPhase, err = generateMnenomicPhase(wordCount)
		if err != nil {
			return
		}
		var encrypted string
		encrypted, err = crypto.EncryptRAFormat([]byte(mnemonicPhase), secret.Base64PublicKey, go_crypto.SHA256)
		if err != nil {
			fmt.Println(err)
			return
		}
		outputMessage = fmt.Sprintf("{\r\n\"unEncryptedMnemonicPhase\":\"%s\",\r\n\"mnemonicPhase\":\"%s\"\r\n}", mnemonicPhase, encrypted)
	}

	if outputFile != nil {
		_, err2 := outputFile.WriteString(outputMessage)
		if err2 != nil {
			log.Fatal(err2)
		}
	} else {
		fmt.Println(outputMessage)
	}

}

func processRSA(cipher []byte) (string, error) {
	var err error
	var outputMessage string
	outputMessage, err = crypto.EncryptRAFormat(cipher, secret.Base64PublicKey, go_crypto.SHA256)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return outputMessage, nil
}

func processAES(cipher []byte, aesKey string) (string, error) {
	var aesKeyByte []byte
	var err error
	var outputMessage string
	if aesKey == "" {
		aesKeyByte = make([]byte, 32)
		_, err := rand.Read(aesKeyByte)
		if err != nil {
			fmt.Println("Error generating AES key: ", err.Error())
			return "", err
		}
	} else {
		aesKeyByte, err = base64.StdEncoding.DecodeString(aesKey)
		if err != nil {
			fmt.Println("Error decoding AES key: ", err.Error())
			return "", err
		}
	}
	var result string
	if len(cipher) > 0 {
		result, err = crypto.EncryptAES(cipher, aesKeyByte)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
		outputMessage = fmt.Sprintf("[AES encrypted base64 string]\r\n%s", result)

	} else {
		encryptedAESKey, _ := crypto.EncryptRAFormat([]byte(base64.StdEncoding.EncodeToString(aesKeyByte)), secret.Base64PublicKey, go_crypto.SHA256)
		outputMessage = fmt.Sprintf("[Base64 AES key]\r\n%s\r\n[Encrypted Base64 AES key]\r\n%s", base64.StdEncoding.EncodeToString(aesKeyByte), encryptedAESKey)
	}
	return outputMessage, nil
}

func processServerAuthKey() (string, error) {
	priKey, priKeyStr, err := crypto.GeneratePrivateKey()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	pubKey := crypto.GeneratePublicKey(priKey)
	encryptedPriKey, err := crypto.EncryptRAFormat([]byte(priKeyStr), secret.Base64PublicKey, go_crypto.SHA256)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	outputMessage := fmt.Sprintf("{\r\n\"privateKey\":\"%s\",\r\n\"publicKey\":\"%s\"\r\n}", encryptedPriKey, pubKey)
	return outputMessage, nil
}

func processPasswordEncrypt(cipher []byte, password string) (string, error) {
	outputMessage, err := crypto.EncryptByPassword(password, cipher)
	if err != nil {
		fmt.Println("Encrypt failed: ", err)
		return "", err
	}
	return outputMessage, nil
}

func generateMnenomicPhase(length int) (string, error) {
	if length < 12 || length > 24 {
		fmt.Println("Invalid length")
		return "", errors.New("Invalid length")
	}
	if length%3 != 0 {
		fmt.Println("Length must be 12, 15, 18, 21, 24")
		return "", errors.New("Length must be 12, 15, 18, 21, 24")
	}
	bitLength := length / 3.0 * 32

	entropy, err := bip39.NewEntropy(bitLength)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}
