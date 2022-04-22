package test

import (
	"encoding/base64"
	"eurus-backend/foundation/crypto"
	"eurus-backend/secret"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	keyText := secret.AdminUserKey
	aesKey, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		t.Fatal(err)
	}

	cipher, err := crypto.EncryptAES([]byte("Hello"), aesKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cipher)

	data, _ := base64.StdEncoding.DecodeString("iF3+6y3btWkZQC2Os0TKixcT5PmHXIt7Tsh0CSimlY4=")
	base64Encoded, _ := crypto.DecryptAES(data, aesKey)
	result, _ := base64.StdEncoding.DecodeString(base64Encoded)
	fmt.Println(string(result))
}
