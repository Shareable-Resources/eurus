package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	defaultKey, _ := uuid.NewRandom()
	fmt.Printf("Input API key in plain text, or leave it blank to use default value [%v]:\n", defaultKey.String())
	scanner.Scan()
	key := strings.TrimSpace(scanner.Text())
	if key == "" {
		key = defaultKey.String()
	}

	defaultSalt, _ := uuid.NewRandom()
	fmt.Printf("Input salt, or leave it blank to use default value [%v]:\n", defaultSalt.String())
	scanner.Scan()
	salt := strings.TrimSpace(scanner.Text())
	if salt == "" {
		salt = defaultSalt.String()
	}

	var merchantCode string
	for {
		fmt.Printf("Input merchant code:\n")
		scanner.Scan()
		merchantCode = strings.TrimSpace(scanner.Text())
		if merchantCode != "" {
			break
		}
	}

	var version int = argon2.Version
	var memory uint32 = 65536
	var time uint32 = 1
	var parallelism uint8 = 8

	saltBytes := []byte(merchantCode + "." + salt)
	idKey := argon2.IDKey([]byte(key), saltBytes, time, memory, parallelism, 32)

	result := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version,
		memory,
		time,
		parallelism,
		base64.RawStdEncoding.EncodeToString(saltBytes),
		base64.RawStdEncoding.EncodeToString(idKey))

	fmt.Println("================================================================")
	fmt.Println("API key:")
	fmt.Println(key)
	fmt.Println("Salt:")
	fmt.Println(salt)
	fmt.Println("Merchant code:")
	fmt.Println(merchantCode)
	fmt.Println("PHC string format result:")
	fmt.Println(result)
}
