package server

import (
	"encoding/json"
	"eurus-backend/foundation/ws/ws_message"
)

type ControlRequestMessage struct {
	Data []string
	ws_message.RequestMessage
}

func (me *ControlRequestMessage) UnmarshalJSON(data []byte) error {
	type cloneType ControlRequestMessage

	if err := json.Unmarshal(data, (*cloneType)(me)); err != nil {
		return err
	}

	byteDataPtr := me.RequestMessage.Data.(*json.RawMessage)
	err := json.Unmarshal(*byteDataPtr, &me.Data)

	return err
}
