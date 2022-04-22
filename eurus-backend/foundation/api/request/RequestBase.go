package request

import "eurus-backend/foundation/auth_base"

type RequestBase struct { //implements IRequest
	LoginToken  auth_base.ILoginToken `json:"-"`
	Method      string                `json:"-"`
	RequestPath string                `json:"-"`
	Nonce       string                `json:"nonce"`
}

func (me *RequestBase) SetNonce(nonce string) {
	me.Nonce = nonce
}

func (me *RequestBase) GetNonce() string {
	return me.Nonce
}

func (me *RequestBase) SetRequestPath(path string) {
	me.RequestPath = path
}

func (me *RequestBase) GetRequestPath() string {
	return me.RequestPath
}

func (me *RequestBase) SetMethod(method string) {
	me.Method = method
}

func (me *RequestBase) GetMethod() string {
	return me.Method
}

func (me *RequestBase) SetLoginToken(loginToken auth_base.ILoginToken) {
	me.LoginToken = loginToken
}

func (me *RequestBase) GetLoginToken() auth_base.ILoginToken {
	return me.LoginToken
}

func InitRequestBase() *RequestBase {
	requestBase := &RequestBase{}
	return requestBase
}
