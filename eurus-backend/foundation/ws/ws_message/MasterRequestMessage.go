package ws_message

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MasterRequestMessage struct {
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Sign      string `json:"sign"`

	Request *RequestMessage
}

func CreateMasterRequestMessage(requestData interface{}, methodName string) (*MasterRequestMessage, error) {
	reqMsg := new(RequestMessage)
	reqMsg.Data = requestData
	reqMsg.Timestamp = time.Now().Unix()
	reqMsg.Nonce = uuid.New().String()
	reqMsg.MethodName = methodName
	reqMsgData, err := json.Marshal(reqMsg)
	if err != nil {
		return nil, err
	}
	req := new(MasterRequestMessage)
	req.Message = string(reqMsgData)
	req.Timestamp = reqMsg.Timestamp
	req.Nonce = reqMsg.Nonce
	req.Request = reqMsg

	return req, nil
}
