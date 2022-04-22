package kyc_model

import (
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	kyc_const "eurus-backend/user_service/kyc_service/kyc/const"
	"net/http"
)

type KYCDocType int

const (
	KYCDocUnknown KYCDocType = iota
	KYCDocPassport
	KYCDocIDCard
)

type KYCImageType int

const (
	KYCImageUnknown KYCImageType = iota
	KYCImgePassport
	KYCIDCardFront
	KYCIDCardBack
	KYCSelfie
)

func (data KYCImageType) String() string {
	switch data {
	case KYCImageUnknown:
		return "KYCImageUnknown"
	case KYCImgePassport:
		return "KYCImgePassport"
	case KYCIDCardFront:
		return "KYCIDCardFront"
	case KYCIDCardBack:
		return "KYCIDCardBack"
	case KYCSelfie:
		return "KYCSelfie"
	default:
		return "NOTFOUND"
	}

}

//API - kyc/user/getKYCCountryList
//API - user/kyc/getKycCountryList
type RequestGetKYCCountryList struct {
	request.RequestBase
}

func NewRequestGetKYCCountryList() *RequestGetKYCCountryList {
	req := new(RequestGetKYCCountryList)
	req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.GetKYCCountryList
	req.Method = http.MethodPost
	return req
}

type ResponseGetKYCCountryList struct {
	response.ResponseBase
	Data []*KYCCountryCode `json:"data"`
}

func NewResponseGetKYCCountryList() *ResponseGetKYCCountryList {
	obj := new(ResponseGetKYCCountryList)
	return obj
}

//API - /kyc/user/createKYCStatus
type RequestCreateKYCStatus struct {
	request.RequestBase
	UserId         uint64     `gorm:"column:user_id" json:"userId"`
	KYCCountryCode string     `gorm:"column:kyc_country_code" json:"kycCountryCode"`
	KYCDoc         KYCDocType `gorm:"column:kyc_doc" json:"kycDoc"`
}

func NewRequestCreateKYCStatus() *RequestCreateKYCStatus {
	req := new(RequestCreateKYCStatus)
	req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.CreateKYCStatus
	req.Method = http.MethodPost
	return req
}

type ResponseCreateKYCStatus struct {
	Id uint64 `json:"id"`
}

func NewResponseCreateKYCStatus() *ResponseCreateKYCStatus {
	obj := new(ResponseCreateKYCStatus)
	return obj
}

type RequestGetKYCStatusOfUser struct {
	request.RequestBase
	UserId uint64 `gorm:"column:user_id" json:"userId,string"`
}

type ResponseGetKYCStatusOfUser struct {
	KYCLevel int              `json:"kycLevel"`
	Data     *[]UserKYCStatus `json:"data"`
}

func NewRequestGetKYCStatusOfUser() *RequestGetKYCStatusOfUser {
	req := new(RequestGetKYCStatusOfUser)
	//req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.GetKYCStatusOfUser
	req.Method = http.MethodGet
	return req
}

//API - /kyc/admin/getKYCStatusList
type RequestGetKYCStatusList struct {
	request.RequestBase
	KYCStatus     *KYCStatusType `gorm:"column:kyc_status" json:"kycStatus,string"`
	Email         string         `json:"email"`
	WalletAddress string         `json:"walletAddress"`
}

func NewRequestGetKYCStatusList() *RequestGetKYCStatusList {
	req := new(RequestGetKYCStatusList)
	req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.GetKYCStatusList
	req.Method = http.MethodGet
	return req
}

func NewResponseGetKYCStatusList(userId uint64) *RequestGetKYCStatusOfUser {
	obj := new(RequestGetKYCStatusOfUser)
	obj.UserId = userId
	return obj
}

type ResponseGetKYCStatusList struct {
	TotalRows int64                  `json:"totalRows"`
	Data      *[]UserKYCStatusDetail `json:"data"`
}

//API - /kyc/admin/loginAdminUser
type RequestLoginAdminUser struct {
	request.RequestBase
	AdminUser
	//Username string `gorm:"column:username" json:"username"`
	//Password string `gorm:"column:password" json:"password"`
}

func NewRequestLoginAdminUser() *RequestLoginAdminUser {
	req := new(RequestLoginAdminUser)
	req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.LoginAdminUser
	req.Method = http.MethodGet
	return req
}

func NewResponseLoginAdminUser(userId uint64) *RequestGetKYCStatusOfUser {
	obj := new(RequestGetKYCStatusOfUser)
	obj.UserId = userId
	return obj
}

type ResponseLoginAdminUser struct {
	response.ResponseBase
	Data *[]UserKYCStatus `json:"data"`
}

type RequestSubmitKYCDocument struct {
	request.RequestBase
	UserKYCStatusId uint64       `json:"userKYCStatusId"`
	ImageType       KYCImageType `json:"imageType"`
	FileExtension   string       `json:"fileExtension"`
}

type RequestCreateAdminUser struct {
	request.RequestBase
	AdminUser
}

func NewRequestCreateAdminUser() *RequestCreateAdminUser {
	obj := new(RequestCreateAdminUser)
	obj.Method = http.MethodPost
	return obj
}

type ResponseCreateAdminUser struct {
	Username string `gorm:"column:username" json:"username"`
}

//kyc/user/submitKYCApproval
type RequestSubmitKYCApproval struct {
	request.RequestBase
	Id        uint64        `gorm:" json:"id"`
	KYCStatus KYCStatusType `gorm:" json:"kycStatus"`
}

func NewRequestSubmitKYCApproval() *RequestSubmitKYCApproval {
	obj := new(RequestSubmitKYCApproval)
	obj.Method = http.MethodPost
	return obj
}

type ResponseSubmitKYCApproval struct {
	response.ResponseBase
	AdminUser
}

//kyc/user/resetKYCStatus
type RequestResetKYCStatus struct {
	request.RequestBase
	AdminUser
	Id uint64 `gorm:" json:"id"`
}

func NewRequestResetKYCStatus() *RequestResetKYCStatus {
	obj := new(RequestResetKYCStatus)
	obj.Method = http.MethodPost
	return obj
}

type ResponseResetKYCStatus struct {
	response.ResponseBase
	AdminUser
}

//kyc/user/updateKYCStatus
type RequestUpdateKYCStatus struct {
	request.RequestBase
	UserKYCStatus
}

func NewRequestUpdateKYCStatus() *RequestUpdateKYCStatus {
	obj := new(RequestUpdateKYCStatus)
	obj.Method = http.MethodPost
	return obj
}

type ResponseUpdateKYCStatus struct {
	PromotedLevel int `json:"promotedLevel"`
}

type FullResponseUpdateKYCStatus struct {
	response.ResponseBase
	Data *ResponseUpdateKYCStatus `json:"data"`
}

type RequestChangeAdminPassword struct {
	request.RequestBase
	AdminUser
	NewPassword string `json:"newPassword"`
}

func NewRequestChangeAdminPassword() *RequestChangeAdminPassword {
	obj := new(RequestChangeAdminPassword)
	obj.Method = http.MethodPost
	return obj
}

type ResponseChangeAdminPassword struct {
	Username string `gorm:"column:username" json:"username"`
}

//kyc/admin/refres
type RefreshTokenRequest struct {
	request.RequestBase
}

type RefreshTokenResponse struct {
	Token            string `json:"token"`
	ExpireTime       int64  `json:"expiryTime"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	CreatedDate      int64  `json:"createdDate"`
}

func NewRefreshTokenRequest() *RefreshTokenRequest {
	req := new(RefreshTokenRequest)
	req.RequestPath = kyc_const.RootPath + kyc_const.EndPoint.RefreshToken
	req.Method = http.MethodPost
	return req
}
