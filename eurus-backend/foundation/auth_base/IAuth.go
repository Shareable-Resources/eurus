package auth_base

import "eurus-backend/foundation"

type IAuth interface {
	IsLoggedIn() bool
	LoginAuthServer(config IAuthBaseConfig) error
	GetLoginToken() string
	GetSessionId() int64
	GenerateLoginToken(userId string) (ILoginToken, error)
	VerifyLoginToken(loginToken string, verifyMode VerifyMode) (bool, ILoginToken, *foundation.ServerError)
	RefreshLoginToken(loginToken string) (ILoginToken, *foundation.ServerError)
	RequestNonRefreshableLoginToken(userId string, duration int64, tokenType int16) (ILoginToken, *foundation.ServerError)
	GetAdditionalInfo() string
	RevokeLoginToken(loginToken string) (ILoginToken, *foundation.ServerError)
	VerifySignature(data []byte, sign string, serviceId int64) (bool, *foundation.ServerError)
	SetLoginHandler(loginHandler func(IAuth))
	GetConfig() IAuthBaseConfig
}

type VerifyMode int

const (
	VerifyModeService VerifyMode = 1 << iota
	VerifyModeUser
)

type ILoginToken interface {
	GetToken() string
	GetExpiredTime() int64
	GetCreatedDate() int64
	GetLastModifiedDate() int64
	GetUserId() string
	GetTokenType() int16
	SetToken(token string)
	SetUserId(userId string)
}
