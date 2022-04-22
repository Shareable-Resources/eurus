package password

import (
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/ws/ws_message"
	"eurus-backend/password_service/password_api"
	"math"
	"time"
)

func VerifySignature(config *PasswordServerConfig, masterReq *ws_message.MasterRequestMessage, req *ws_message.RequestMessage) error {

	isValid, err := crypto.VerifyRSASignFromBase64(passwordClientPublicKey, masterReq.Message, masterReq.Sign)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to verify signature for request nonce: ", req.Nonce)
		return err
	}

	if !isValid {
		return errors.New("Signature mismatch")
	}

	if masterReq.Timestamp != req.Timestamp {
		log.GetLogger(log.Name.Root).Errorln("Master request timestamp does not match with request timestamp. Nonce: ", masterReq.Nonce)
		return errors.New("Master request timestamp does not match with request timestamp")
	}

	currTime := time.Now().Unix()
	delta := currTime - req.Timestamp
	if math.Abs(float64(delta)) > 10 {
		log.GetLogger(log.Name.Root).Errorln("Timestamp is not within acceptable tolerant range. Nonce: ", masterReq.Nonce)
		return errors.New("Timestamp is not within acceptable tolerant range")
	}

	return nil
}

func createErrorResponse(req *ws_message.RequestMessage, returnCode int64, message string) *ws_message.ResponseMessage {
	res := new(ws_message.ResponseMessage)

	res.MethodName = req.MethodName
	res.ReturnCode = returnCode
	res.Nonce = req.Nonce
	res.Message = message
	res.Timestamp = time.Now().Unix()
	return res
}

func createSuccessResponse(req *ws_message.RequestMessage, innerRes interface{}) *ws_message.ResponseMessage {
	res := new(ws_message.ResponseMessage)

	res.MethodName = req.MethodName
	res.ReturnCode = int64(foundation.Success)
	res.Nonce = req.Nonce
	res.Message = foundation.Success.String()
	res.Timestamp = time.Now().Unix()
	res.Data = innerRes

	return res
}

func processPasswordRequest(config *PasswordServerConfig, req *ws_message.RequestMessage) *ws_message.ResponseMessage {
	if config.Password == "" {
		return createErrorResponse(req, int64(foundation.RecordNotFound), "Password not ready yet")
	}

	res := new(password_api.PasswordResponse)
	res.EncryptedPassword = config.Password

	baseRes := createSuccessResponse(req, res)
	return baseRes
}
