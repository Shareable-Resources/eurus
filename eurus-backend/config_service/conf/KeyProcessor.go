package conf

import (
	"encoding/json"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"

	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
)

func GenerateKeyPair() (*conf_api.KeyPair, error) {
	privateKey, base64PrivateKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	base64PublicKey := crypto.GeneratePublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	keyPair := &conf_api.KeyPair{
		PrivateKey: base64PrivateKey,
		PublicKey:  base64PublicKey,
	}
	return keyPair, nil
}

func GetKeyPair(req *request.RequestBase) *response.ResponseBase {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to generate RSA KeyPair"+err.Error(), "\nRequest Params: "+string(reqStr))
		res := response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		return res
	}
	res := response.CreateSuccessResponse(req, keyPair)
	return res
}
