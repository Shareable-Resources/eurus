package main

import (
	"encoding/base64"
	"eurus-backend/foundation/crypto"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	argsWithProg := os.Args

	if len(argsWithProg) < 2 {
		printUsage()
		os.Exit(1)
	}

	mode := argsWithProg[1]

	if mode != "rsa" && mode != "m" && mode != "p" && mode != "aes" {
		printUsage()
		os.Exit(1)
	}

	var length int = 0
	var err error
	var fileName string
	var file *os.File

	for i := 2; i < len(argsWithProg); i += 2 {
		arg := argsWithProg[i]
		if i+1 >= len(argsWithProg) {
			printUsage()
			os.Exit(1)
		}
		switch arg {
		case "-l":
			lengthStr := argsWithProg[i+1]
			length, err = strconv.Atoi(lengthStr)
		case "-o":
			fileName = argsWithProg[i+1]
		}

		if err != nil {
			printUsage()
			os.Exit(1)
		}
	}

	if mode == "m" || mode == "p" {
		if length == 0 {
			fmt.Println("Missing -l parameter or invalid -l parameter")
			os.Exit(1)
		}
	}

	if fileName != "" {
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Println("Error open file: ", err)
			os.Exit(1)
		}
	}

	if file != nil {
		defer file.Close()
	}

	var output string
	switch mode {
	case "rsa":
		output, err = GenerateRsaKey(file)
	case "aes":
		output, err = GenerateAesKey(file)
	case "m":
		output, err = GenerateMnenomicPhase(length)
	case "p":
		output, err = GeneratePassword(length)
	}

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if file != nil {
		file.Write([]byte(output))
		fmt.Println("Exported to " + fileName)
	} else {
		fmt.Println(output)
	}

	os.Exit(0)
}

func GenerateRsaKey(file *os.File) (string, error) {
	priKey, priKeyBase64, err := crypto.GeneratePrivateKey()
	if err != nil {
		err := errors.Wrap(err, "Generate private key failed")
		return "", err
	}

	pubKey := crypto.GeneratePublicKey(priKey)

	outputText := "[private]\r\n" + priKeyBase64 + "\r\n[public]\r\n" + pubKey + "\r\n"
	return outputText, nil
}

func GenerateAesKey(file *os.File) (string, error) {
	aesKeyByte := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(aesKeyByte)
	aesKeyStr := base64.StdEncoding.EncodeToString(aesKeyByte)

	return aesKeyStr, nil
}

func GenerateMnenomicPhase(length int) (string, error) {
	if length < 12 || length > 24 {
		return "", errors.New("Invalid length")
	}
	if length%3 != 0 {
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

func GeneratePassword(length int) (string, error) {
	rand.Seed(time.Now().UnixNano())

	var password string
	for {
		password = ""
		var hasCapital, hasLittle, hasDigit bool
		for i := 0; i < length; i++ {
			val := rand.Intn(62)
			if val >= 0 && val <= 9 {
				hasDigit = true
				val += 48
			} else if val > 9 && val < 36 {
				hasCapital = true
				val += 55
			} else {
				hasLittle = true
				val += 61
			}

			password += string(val)
		}
		if hasCapital && hasLittle && hasDigit {
			break
		}
	}
	return password, nil
}

func printUsage() {
	fmt.Println("Usage: keyGen [mode] -l [mnemonic phase length/ password length]  -o [output file name]")
	fmt.Println("mode:")
	fmt.Println("rsa - generate RSA key pair")
	fmt.Println("m: generate mnemonic phase. MUST input -l option")
	fmt.Println("p: generate a random password. MUST INPUT -c option")
	fmt.Println("aes - generate AES key")
	fmt.Println("Options:")
	fmt.Println("-l: only appear when mode is m or p. Mnemonic phase length (12, 15, 21, 24) or password length")

}
