package ga

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

func EnableTwoFA(username string, appName string, suffix string) (string, string, error) {
	data := username + time.Now().Format(time.RFC3339Nano)
	secret := strings.ToUpper(base32.StdEncoding.EncodeToString(
		hmac.New(sha1.New, []byte(data)).Sum(nil)))

	qrCode := GetQRCode(username, suffix, appName, secret)
	return secret, qrCode, nil
}

func GetQRCode(username string, suffix string, appName string, secret string) string {
	qrurl := fmt.Sprintf("otpauth://totp/%s_%s@%s?secret=%s", username, suffix, appName, secret)
	png, _ := qrcode.Encode(qrurl, qrcode.Medium, 256)
	return base64.StdEncoding.EncodeToString(png)
}

func GenTwoFACode(secret string) (string, error) {

	secretData, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	// secretData := []byte(secret)

	h := hmac.New(sha1.New, secretData)
	err = binary.Write(h, binary.BigEndian, time.Now().Unix()/30)
	if err != nil {
		return "", err
	}
	data := h.Sum(nil)
	off := data[19] & 0xf
	it := data[off : off+4]
	it[0] = it[0] & 0x7f
	code := fmt.Sprintf("%06d", binary.BigEndian.Uint32(it)%1000000)
	return code, nil
}

func VerifyTwoFACode(secret, code string) bool {
	expect, _ := GenTwoFACode(secret)
	return expect == code
}
