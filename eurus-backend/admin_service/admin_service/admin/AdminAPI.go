package admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation/api/request"
)

type AdminLoginRequest struct {
	request.RequestBase
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type AdminLoginNextAction int

const (
	LoginGA AdminLoginNextAction = iota
	BindGA
)

type AdminLoginResponse struct {
	NextAction     AdminLoginNextAction `json:"nextAction"`
	AccessToken    string               `json:"accessToken"`
	GASecretQRCode string               `json:"gaSecretQRCode"`
}

type AdminLoginGARequest struct {
	request.RequestBase
	AccessToken string `json:"accessToken"`
	GACode      string `json:"gaCode"`
}

type VerifyGARequest struct {
	request.RequestBase
	GACode  string `json:"gaCode"`
	LoginIp string
}

type VerifyGAResponse struct {
	AccessToken           string                                `json:"accessToken"`
	FeaturePermissionList []*admin_common.FeaturePermissionPair `json:"featurePermissionList"`
}

type QueryAccountListRequest admin_common.QueryAccountListRequest
type QueryAccountListResponse []*AdminUserEx
