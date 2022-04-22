package auth

import (
	"eurus-backend/foundation/ws/ws_message"
)

func init() {
	//Register request data field corresponding class type
	ws_message.RequestDataFieldFactoryMap[ApiAuth] = func(string) interface{} {
		return &AuthenticateRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiRequestPaymentLoginToken] = func(string) interface{} {
		return &NonRefreshableLoginTokenRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiRequestLoginToken] = func(string) interface{} {
		return &RequestLoginTokenRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiVerifyLoginToken] = func(string) interface{} {
		return &VerifyLoginTokenRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiRefreshLoginToken] = func(string) interface{} {
		return &RefreshLoginTokenRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiRevokeLoginToken] = func(string) interface{} {
		return &RevokeLoginTokenRequest{}
	}

	ws_message.RequestDataFieldFactoryMap[ApiVerifySign] = func(string) interface{} {
		return &VerifySignRequest{}
	}

	////////////////////RESPONSE///////////////////////////
	ws_message.ResponseDataFieldFactoryMap[ApiAuth] = func(string) interface{} {
		return &AuthenticateResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiRequestLoginToken] = func(string) interface{} {
		return &RequestLoginTokenResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiVerifyLoginToken] = func(string) interface{} {
		return &VerifyLoginTokenResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiRefreshLoginToken] = func(string) interface{} {
		return &RefreshLoginTokenResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiRevokeLoginToken] = func(string) interface{} {
		return &RevokeLoginTokenResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiRequestPaymentLoginToken] = func(string) interface{} {
		return &NonRefreshableLoginTokenResponse{}
	}

	ws_message.ResponseDataFieldFactoryMap[ApiVerifySign] = func(string) interface{} {
		return &VerifySignResponse{}
	}
	SetupJwt([]byte("Dummy"), 3600)
}
