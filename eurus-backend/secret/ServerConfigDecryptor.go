package secret

import (
	"encoding/base64"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"

	"github.com/pkg/errors"
)

func DecryptSensitiveConfig(config *server.ServerConfigBase) error {
	//config.DBAESKey
	if config.DBAESKey != "" {
		originMessage, err := crypto.DecryptRAFormat(config.DBAESKey, base64PrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to decrypt key 1: ", err.Error())
			return errors.Wrap(err, "Unable to decrypt key 1")
		}
		data, err := base64.StdEncoding.DecodeString(originMessage)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Decode base64 key 1 failed: ", err.Error())
			return errors.Wrap(err, "Decode base64 key 1 failed")
		}
		config.DBAESKey = string(data)
	}

	//config.HdWalletPrivateKey
	if config.HdWalletPrivateKey != "" {
		originMessage, err := crypto.DecryptRAFormat(config.HdWalletPrivateKey, base64PrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to decrypt key 2: ", err.Error())
			return errors.Wrap(err, "Unable to decrypt key 2")
		}
		data, err := base64.StdEncoding.DecodeString(originMessage)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Invalid key 2: ", err.Error())
			return errors.Wrap(err, "Invalid decrypt key 2")
		}
		config.HdWalletPrivateKey = string(data)
	}

	if config.PrivateKey != "" && config.PrivateKey != "N/A" {
		//config.PrivateKey
		originMessage, err := crypto.DecryptRAFormat(config.PrivateKey, base64PrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to decrypt key 3: ", err.Error())
			return errors.Wrap(err, "Unable to decrypt key 3")
		}
		data, err := base64.StdEncoding.DecodeString(originMessage)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Invalid key 3: ", err.Error())
			return errors.Wrap(err, "Invalid key 3")
		}
		config.PrivateKey = string(data)
	}

	if config.MnemonicPhase != "" {
		originMessage, err := crypto.DecryptRAFormat(config.MnemonicPhase, base64PrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to decrypt key 4: ", err.Error())
			return errors.Wrap(err, "Unable to decrypt key 4")
		}
		data, err := base64.StdEncoding.DecodeString(originMessage)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Invalid key 4: ", err.Error())
			return errors.Wrap(err, "Invalid key 4")
		}
		config.MnemonicPhase = string(data)
	}
	return nil
}

func DecryptConfigValue(value string) (string, error) {
	return crypto.DecryptRAFormat(value, base64PrivateKey)
}

// func GetDecryptedMnemonicAesKey() ([]byte, error) {
// 	// Change the AES key to plain text (base64 encoded) so code below is deprecated
// 	return base64.StdEncoding.DecodeString(mnemonicAesKey)

// Note that crypto.DecryptRAFormat returns b64encode(message), not message itself
// Message here is b64encoded AES key, that means in order to get the binary form, need to b64decode twice
/*b64encoded, err := crypto.DecryptRAFormat(mnemonicAesKey, base64PrivateKey)
if err != nil {
	log.GetLogger(log.Name.Root).Errorln("Unable to decrypt mnemonic AES key: ", err.Error())
	return nil, err
}

message, err := base64.StdEncoding.DecodeString(b64encoded)
if err != nil {
	log.GetLogger(log.Name.Root).Errorln("Unable to get mnemonic base64 AES key: ", err.Error())
	return nil, err
}

key, err := base64.StdEncoding.DecodeString(string(message))
if err != nil {
	log.GetLogger(log.Name.Root).Errorln("Unable to decode mnemonic AES key: ", err.Error())
	return nil, err
}

return key, nil*/
// }
