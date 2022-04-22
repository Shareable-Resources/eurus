package ws_message

import (
	"encoding/json"
	"eurus-backend/foundation/api/request"
)

type RequestMessage struct { // implements json.Unmarshaler
	request.RequestBase
	MethodName string      `json:"methodName"`
	Timestamp  int64       `json:"timestamp"`
	Data       interface{} `json:"data"`
}

var RequestDataFieldFactoryMap map[string]func(string) interface{} = make(map[string]func(string) interface{})

func (req *RequestMessage) UnmarshalJSON(data []byte) error {
	type cloneType RequestMessage
	rawMsg := json.RawMessage{}
	req.Data = &rawMsg

	if err := json.Unmarshal(data, (*cloneType)(req)); err != nil {
		return err
	}

	factoryFunc, ok := RequestDataFieldFactoryMap[req.MethodName]

	if !ok {
		return nil
	} else {
		var intf interface{} = factoryFunc(req.MethodName)
		if len(rawMsg) > 0 {
			err := json.Unmarshal(rawMsg, intf)
			if err != nil {
				return err
			}
		}
		req.Data = intf //The debugger shows intf is nil, however it is not!

	}
	return nil
}
