package ws_message

import (
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"time"
)

type ResponseMessage struct {
	response.ResponseBase
	MethodName string `json:"methodName"`
	Timestamp  int64  `json:"timestamp"`
	Sign       string `json:"sign"`
}

var ResponseDataFieldFactoryMap map[string]func(string) interface{} = make(map[string]func(string) interface{})

func (me *ResponseMessage) UnmarshalJSON(data []byte) error {
	type cloneType ResponseMessage
	rawMsg := json.RawMessage{}
	me.Data = &rawMsg

	err := json.Unmarshal(data, (*cloneType)(me))
	if err != nil {
		return err
	}

	factoryFunc, ok := ResponseDataFieldFactoryMap[me.MethodName]

	if !ok {
		return nil
	} else {
		var intf interface{} = factoryFunc(me.MethodName)
		if len(rawMsg) > 0 {
			err := json.Unmarshal(rawMsg, intf)
			if err != nil {
				return err
			}
		}
		me.Data = intf //The debugger shows intf is nil, however it is not!

	}
	return nil

}

func CreateSuccessResponseMessage(methodName string, data interface{}, nonce string) *ResponseMessage {
	res := new(ResponseMessage)
	res.Message = foundation.Success.String()
	res.ReturnCode = int64(foundation.Success)
	res.Nonce = nonce
	res.Timestamp = time.Now().Unix()
	res.MethodName = methodName
	res.Data = data
	return res
}

func CreateErrorResponseMessage(code foundation.ServerReturnCode, message string, nonce string) *ResponseMessage {
	res := new(ResponseMessage)
	res.Message = message
	res.ReturnCode = int64(code)
	res.Nonce = nonce
	res.Timestamp = time.Now().Unix()
	return res
}
