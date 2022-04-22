package request

import (
	"eurus-backend/foundation/auth_base"
)

type IRequest interface {
	SetLoginToken(token auth_base.ILoginToken)
	SetMethod(method string)
	GetMethod() string
	SetRequestPath(path string)
	GetRequestPath() string
	SetNonce(nonce string)
	GetNonce() string
	GetLoginToken() auth_base.ILoginToken
}
