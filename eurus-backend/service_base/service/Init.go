package service

import "eurus-backend/foundation/ws/ws_message"

func init() {
	ws_message.RequestDataFieldFactoryMap["passwordRequest"] = func(string) interface{} {
		return new(PasswordRequest)
	}
}
