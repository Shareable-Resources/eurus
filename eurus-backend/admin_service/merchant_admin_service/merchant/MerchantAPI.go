package merchant_admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/admin_service/merchant_common"
	"eurus-backend/foundation/api/request"
	"net/http"
)

type MerchantAccountState int

const (
	MerchantAccountUnknown = iota
	MerchantAccountNormal
	MerchantAccountForceChangePassword
	MerchantAccountSuspended
)

type MerchantLoginRequest struct {
	request.RequestBase
	MerchantId   uint64 `json:"merchantId"`
	UserName     string `json:"userName"`
	PasswordHash string `json:"passwordHash"`
}

type MerchantLoginResponse struct {
	AccountState MerchantAccountState `json:"accountState"`
	LoginToken   string               `json:"loginToken"`
}

type MerchantRefreshLoginTokenRequest struct {
	request.RequestBase
}

type MerchantRefreshLoginTokenResponse struct {
	Token            string `json:"token"`
	ExpireTime       int64  `json:"expiryTime"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	CreatedDate      int64  `json:"createdDate"`
}

type MerchantChangePasswordRequest struct {
	request.RequestBase
	MerchantId      uint64 `json:"merchantId"`
	UserName        string `json:"username"`
	OldPasswordHash string `json:"oldPasswordHash"`
	NewPasswordHash string `json:"newPasswordHash"`
}

type GetMerchantRefundRequest struct {
	request.RequestBase
	admin_common.PaginateRequest
	Status int `json:"status"`
}

func NewGetMerchantRefundRequest() *GetMerchantRefundRequest {
	req := new(GetMerchantRefundRequest)
	req.Method = http.MethodGet
	req.RequestPath = RootPath + RefundRequestListPath
	req.Status = -1
	return req
}

type GetMerchantRefundResponse struct {
	RecordCount int64
	List        []merchant_common.MerchantRefundRequest
}

type RefundRequest struct {
	request.RequestBase
	RequestId       uint64
	Answer          merchant_common.RefundRequestStatus `json:"answer"`
	Comment         string                              `json:"comment"`
	RefundTransHash string                              `json:"refundTransHash"`
}

func NewRefundRequest(requestId uint64) *RefundRequest {
	req := new(RefundRequest)
	req.Method = http.MethodPost
	req.RequestPath = RootPath + RefundRequestPath
	req.RequestId = requestId
	return req
}
