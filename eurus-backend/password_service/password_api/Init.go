package password_api

import "eurus-backend/foundation/ws/ws_message"

func init() {
	ws_message.RequestDataFieldFactoryMap[PasswordRequestMethodName] = func(string) interface{} {
		return new(PasswordRequest)
	}
	ws_message.ResponseDataFieldFactoryMap[PasswordRequestMethodName] = func(string) interface{} {
		return new(PasswordResponse)
	}
}
