package secret

import (
	"eurus-backend/foundation/crypto"
)

func VerifyPasswordServerSignature(message string, signature string) (bool, error) {
	return crypto.VerifyRSASignFromBase64(passwordServerPublicKey, message, signature)
}

func GeneratePasswordClientSignature(message string) (string, error) {
	return crypto.GenerateRSASignFromBase64(passwordClientPrivateKey, message)
}

func DecryptPasswordServerData(data string) (string, error) {
	return crypto.DecryptRAFormat(data, passwordClientPrivateKey)
}
