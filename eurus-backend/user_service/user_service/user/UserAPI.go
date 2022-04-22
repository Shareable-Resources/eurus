package user

import (
	"errors"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/marketing/banner"
	"eurus-backend/marketing/reward"
	kyc_const "eurus-backend/user_service/kyc_service/kyc/const"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type PlatformType int8

const (
	Mobile PlatformType = iota
	Website
)

var replayAttackTimeDuration = "1m"

func checkTimestamp(timestamp int64) error {
	now := time.Now()
	t1 := now.Add(time.Second * 60)
	t2 := now.Add(time.Second * -60)
	reqTimestamp := time.Unix(timestamp, 0)
	if reqTimestamp.After(t1) && reqTimestamp.Before(t2) {
		reqTimestampStr := strconv.FormatInt(reqTimestamp.UnixNano(), 10)
		log.GetLogger(log.Name.Root).Error("Invalid Request Timestamp: ", reqTimestampStr)
		return errors.New("Invalid time!")
	}
	return nil
}

type WalletBaseRequest struct {
	WalletAddress  string `json:"walletAddress"`
	Timestamp      int64  `json:"timestamp"`
	PublicKey      string `json:"publicKey"`
	DeviceId       string `json:"deviceId"`
	IsPersonalSign bool   `json:"isPersonalSign"`
	LoginLogDetail
}
type LoginLogDetail struct {
	Sign               string `json:"sign"`
	AppVersion         string `json:"appVersion"`
	Os                 string `json:"os"`
	RegistrationSource string `json:"registrationSource"`
}

type ImportWalletRequest struct {
	request.RequestBase
	WalletBaseRequest
}

type ImportWalletResponse struct {
	Token         string `json:"token"`
	ExpiryTime    int64  `json:"expiryTime"`
	LastLoginTime string `json:"lastLoginTime"`
	Status        int    `json:"status"`
	TxHash        string `json:"txHash"`
}

func NewImportWalletRequest() *ImportWalletRequest {
	req := new(ImportWalletRequest)
	req.RequestPath = RootPath + ImportWalletPath
	req.Method = http.MethodPost
	return req
}

func NewImportWalletResponse(loginToken auth_base.ILoginToken, txHash string, user *User, req request.IRequest) *ImportWalletResponse {
	resp := new(ImportWalletResponse)
	resp.Token = loginToken.GetToken()
	unixTimeUTC := time.Unix(loginToken.GetExpiredTime(), 0).UTC() //gives unix time stamp in utc
	//unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)    // converts utc time to RFC3339 format
	resp.ExpiryTime = unixTimeUTC.Unix()
	if !user.LastLoginTime.IsZero() {
		unitTimeInRFC3339 := user.LastLoginTime.Format(time.RFC3339)
		resp.LastLoginTime = unitTimeInRFC3339
	} else {
		resp.LastLoginTime = ""
	}
	resp.Status = 0
	resp.TxHash = txHash
	return resp
}

func (me *ImportWalletRequest) CheckTimestamp() error {
	err := checkTimestamp(me.Timestamp)
	return err
}

type UserLoginBySignatureRequest struct {
	request.RequestBase
	WalletBaseRequest
}

func NewUserLoginBySignatureRequest() (me *UserLoginBySignatureRequest) {
	req := new(UserLoginBySignatureRequest)
	req.RequestPath = RootPath + LoginBySignaturePath
	req.Method = http.MethodPost
	return req
}

type UserLoginBySignatureResponse struct {
	Token                string `json:"token"`
	ExpiryTime           int64  `json:"expiryTime"`
	LastLoginTime        int64  `json:"lastLoginTime"`
	Status               int    `json:"status"`
	LastModifiedDate     int64  `json:"lastModifiedDate"`
	Mnemonic             string `json:"mnemonic"`
	WalletAddress        string `json:"walletAddress"`
	MainnetWalletAddress string `json:"mainnetWalletAddress"`
	OwnerWalletAddress   string `json:"ownerWalletAddress"`
	IsMetaMaskUser       bool   `json:"isMetaMaskUser"`
}

type QueryUserDetailsRequest struct {
	request.RequestBase
	Token string `json:"token"`
}

func NewQueryUserDetailsRequest(token string) (me *QueryUserDetailsRequest) {
	req := new(QueryUserDetailsRequest)
	req.RequestPath = RootPath + GetUserDetailsPath
	req.Method = http.MethodPost
	req.Token = token
	return req
}

type UserLoginId struct {
	LoginAddress string `json:"loginAddress"`
	UserId       uint64 `json:"userId"`
}

type QueryWithdrawAdminFee struct {
	request.RequestBase
	Token string `json:"token"`
}

func NewQueryWithdrawAdminFeeRequest(token string) (me *QueryWithdrawAdminFee) {
	req := new(QueryWithdrawAdminFee)
	req.RequestPath = RootPath + GetWithdrawAdminFeePath
	req.Method = http.MethodPost
	req.Token = token
	return req
}

type FaucetResponse struct {
	TxHash string `json:"txHash"`
	Status int    `json:"status"`
}

type RegistrationRequest struct {
	request.RequestBase
	Email          string `json:"email"`
	PublicKey      string `json:"publicKey"`
	LoginAddress   string `json:"loginAddress"`
	Signature      string `json:"signature"`
	Timestamp      int64  `json:"timestamp"`
	IsPersonalSign bool   `json:"isPersonalSign"`
	DeviceId       string `json:"deviceId"`
}

func NewRegistrationRequest() *RegistrationRequest {
	req := new(RegistrationRequest)
	req.RequestPath = RootPath + RegisterPath
	req.Method = http.MethodPost
	return req
}

type VerificationRequest struct {
	request.RequestBase
	Email     string `json:"email"`
	Code      string `json:"code"`
	DeviceId  string `json:"deviceId"`
	PublicKey string `json:"publicKey"`
}

type SetupPaymentWalletRequest struct {
	request.RequestBase
	UserWalletOwnerAddress string `json:"address"`
	LoginLogDetail
}

func NewSetupPaymentWalletRequest() *SetupPaymentWalletRequest {
	req := new(SetupPaymentWalletRequest)
	req.RequestPath = RootPath + SetupPaymentWalletPath
	req.Method = http.MethodPost
	return req
}

func NewVerificationRequest() *VerificationRequest {
	req := new(VerificationRequest)
	req.RequestPath = RootPath + VerificationPath
	req.Method = http.MethodPost
	return req
}

type VerificationResponse struct {
	UserId   uint64 `json:"userId"`
	Email    string `json:"email"`
	Mnemonic string `json:"mnemonic"`
	//WalletAddress string `json:"walletAddress"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}

type RegistrationResponse struct {
	UserId uint64 `json:"userId"`
	Code   string `json:"code"`
}

type ResendVerificationEmailRequest struct {
	request.RequestBase
	UserId uint64 `json:"userId"`
}

type RequestLoginRequestTokenResponse struct {
	LoginRequestToken string    `json:"loginRequestToken"`
	ExpiredTime       time.Time `json:"expiredTime"`
}

type RequestLoginTokenResponse struct {
	LoginToken string `json:"loginToken"`
}

func NewResendVerificationEmailRequest() *ResendVerificationEmailRequest {
	req := new(ResendVerificationEmailRequest)
	req.RequestPath = RootPath + ResendVerificationEmailPath
	req.Method = http.MethodPost
	return req
}

func NewRequestLoginRequestTokenRequest() *request.RequestBase {
	req := new(request.RequestBase)
	req.RequestPath = RootPath + RequestLoginRequestTokenPath
	req.Method = http.MethodPost
	return req
}

type RequestLoginTokenRequest struct {
	request.RequestBase
	LoginRequestToken string `json:"loginRequestToken"`
}

type RequestChangePasswordRequest struct {
	request.RequestBase
	OwnerWalletAddress    string `json:"ownerWalletAddress"`
	OldOwnerWalletAddress string `json:"oldOwnerWalletAddress"`
	DeviceId              string `json:"deviceId"`
	Timestamp             int64  `json:"timestamp"`
	PublicKey             string `json:"publicKey"`
	OldPublicKey          string `json:"oldPublicKey"`
	Sign                  string `json:"sign"`
	OldSign               string `json:"oldSign"`
	IsPersonalSign        bool   `json:"isPersonalSign"`
}

type RequestChangeLoginPasswordRequest struct {
	request.RequestBase
	OldLoginAddress string `json:"oldLoginAddress"`
	LoginAddress    string `json:"loginAddress"` //New login address
	DeviceId        string `json:"deviceId"`
	Timestamp       int64  `json:"timestamp"`
	PublicKey       string `json:"publicKey"` //New public key
	OldPublicKey    string `json:"oldPublicKey"`
	Sign            string `json:"sign"` //Signature by new private key
	OldSign         string `json:"oldSign"`
	IsPersonalSign  bool   `json:"isPersonalSign"`
}

type RequestUserStorage struct {
	request.RequestBase
	LoginLogDetail
	GetUserStorageResponse
}

type GetUserStorageResponse struct {
	UserId   uint64 `json:"userId"`
	Platform int    `json:"platform"`
	Storage  string `json:"storage"`
	Sequence int    `json:"sequence"`
}
type UserStorageSequenceResponse struct {
	Sequence int `json:"sequence"`
}

func NewRequestChangeLoginPasswordRequest() *RequestChangeLoginPasswordRequest {
	req := new(RequestChangeLoginPasswordRequest)
	req.RequestPath = RootPath + RequestChangeLoginPasswordPath
	req.Method = http.MethodPost
	return req
}

func NewRequestChangePaymentPasswordRequest() *RequestChangePasswordRequest {
	req := new(RequestChangePasswordRequest)
	req.RequestPath = RootPath + RequestChangePaymentPasswordPath
	req.Method = http.MethodPost
	return req
}

func NewRequestLoginTokenRequest() *RequestLoginTokenRequest {
	req := new(RequestLoginTokenRequest)
	req.RequestPath = RootPath + RequestLoginTokenByLoginRequestTokenPath
	req.Method = http.MethodPost
	return req
}

func NewRequestPaymentLoginTokenRequest() *request.RequestBase {
	req := new(request.RequestBase)
	req.RequestPath = RootPath + RequestPaymentLoginTokenPath
	req.Method = http.MethodGet
	return req
}

func NewGetUserPreferenceStorageRequest() *RequestUserStorage {
	req := new(RequestUserStorage)
	req.RequestPath = RootPath + UserPreferenceStoragePath
	req.Method = http.MethodGet
	return req
}

func NewSetUserPreferenceStorageRequest() *RequestUserStorage {
	req := new(RequestUserStorage)
	req.RequestPath = RootPath + UserPreferenceStoragePath
	req.Method = http.MethodPost
	return req
}

type ForgetLoginPasswordRequest struct {
	Email string `json:"email"`
	request.RequestBase
}

func NewForgetLoginPassword() *ForgetLoginPasswordRequest {
	req := new(ForgetLoginPasswordRequest)
	req.RequestPath = RootPath + ForgetLoginPasswordPath
	req.Method = http.MethodPost
	return req
}

type VerifyForgetPaymentPasswordRequest struct {
	request.RequestBase
	Code     string `json:"code"`
	DeviceId string `json:"deviceId"`
}

func NewVerifyForgetPaymentPasswordReqeust() *VerifyForgetPaymentPasswordRequest {
	req := new(VerifyForgetPaymentPasswordRequest)
	req.RequestPath = RootPath + VerifyForgetPaymentPasswordPath
	req.Method = http.MethodPost
	return req
}

type VerifyForgetLoginPasswordRequest struct {
	request.RequestBase
	Code  string `json:"code"`
	Email string `json:"email"`
}

func NewVerifyForgetLoginPasswordReqeust() *VerifyForgetLoginPasswordRequest {
	req := new(VerifyForgetLoginPasswordRequest)
	req.RequestPath = RootPath + VerifyForgetLoginPasswordPath
	req.Method = http.MethodPost
	return req
}

type ResetLoginPasswordReqeust struct {
	request.RequestBase
	DeviceId       string `json:"deviceId"`
	Timestamp      int64  `json:"timestamp"`
	LoginAddress   string `json:"loginAddress"`
	PublicKey      string `json:"publicKey"`
	Sign           string `json:"sign"`
	IsPersonalSign bool   `json:"isPersonalSign"`
}

func NewResetLoginPasswordReqeust() *ResetLoginPasswordReqeust {
	req := new(ResetLoginPasswordReqeust)
	req.RequestPath = RootPath + ResetLoginPasswordPath
	req.Method = http.MethodPost
	return req
}

type ResetPaymentPasswordReqeust struct {
	request.RequestBase
	DeviceId           string `json:"deviceId"`
	Timestamp          int64  `json:"timestamp"`
	OwnerWalletAddress string `json:"ownerWalletAddress"`
	PublicKey          string `json:"publicKey"`
	Sign               string `json:"sign"`
	IsPersonalSign     bool   `json:"isPersonalSign"`
}

func NewResetPaymentPasswordReqeust() *ResetPaymentPasswordReqeust {
	req := new(ResetPaymentPasswordReqeust)
	req.RequestPath = RootPath + ResetPaymentPasswordPath
	req.Method = http.MethodPost
	return req
}

type VerifyForgetPasswordResponse struct {
	Token    string `json:"token"`
	Mnemonic string `json:"mnemonic"`
}

type SetupUserWalletResponse struct {
	Token                string `json:"token"`
	WalletAddress        string `json:"walletAddress"`
	MainnetWalletAddress string `json:"mainnetWalletAddress"`
	IsMetamaskAddr       bool   `json:"isMetamaskAddr"`
}

type EmailWalletAddressRequest struct {
	request.RequestBase
	Email         string `json:"email"`
	WalletAddress string `json:"walletAddress"`
}

func NewEmailWalletAddressRequest() *EmailWalletAddressRequest {
	req := new(EmailWalletAddressRequest)
	req.RequestPath = RootPath + FindEmailWalletAddressPath
	req.Method = http.MethodPost
	return req
}

type EmailWalletAddressResponse struct {
	Email         string `json:"email"`
	WalletAddress string `json:"walletAddress"`
	UserType      int    `json:"userType"`
}

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
	req.RequestPath = RootPath + RefreshTokenPath
	req.Method = http.MethodPost
	return req
}

type SignTransactionRequest struct {
	request.RequestBase
	Value         string `json:"value"`
	GasPrice      uint64 `json:"gasPrice"`
	InputFunction string `json:"inputFunction"`
}

func NewSignTransactionRequest() *SignTransactionRequest {
	req := new(SignTransactionRequest)
	req.RequestPath = RootPath + SignTransactionPath
	req.Method = http.MethodPost
	return req
}

type SignTransactionResponse struct {
	SignedTx string `json:"signedTx"`
}

type RecentTransaction struct {
	DecimalPlace int           `json:"decimalPlace"`
	TransList    []interface{} `json:"transList"`
}

type QueryRecentTransactionDetailsRequest struct {
	UserId         uint64 `json:"userId"`
	LoginToken     string `json:"loginToken"`
	CurrencySymbol string `json:"currencySymbol"`
	ChainId        int64  `json:"chainId"`
	request.RequestBase
}

func NewQueryRecentTransactionRequest(loginToken string) (me *QueryRecentTransactionDetailsRequest) {
	req := new(QueryRecentTransactionDetailsRequest)
	req.LoginToken = loginToken
	req.RequestPath = RootPath + GetRecentTransactionPath
	req.Method = http.MethodGet
	return req
}

type ClientVersionRequest struct {
	request.RequestBase
}

type ClientVersion struct {
	IPhoneMinimumVersion  string `json:"iPhoneMinimumVersion"`
	AndroidMinimumVersion string `json:"androidMinimumVersion"`
}

func NewQueryClientVersionRequest() *ClientVersionRequest {
	req := new(ClientVersionRequest)
	req.RequestPath = "/user/clientVersion"
	req.Method = http.MethodGet
	return req
}

type ServerConfig struct {
	EurusRPCDomain                     string `json:"eurusRPCDomain"`
	EurusRPCPort                       int    `json:"eurusRPCPort"`
	EurusRPCProtocol                   string `json:"eurusPRCProtocol"`
	EurusChainId                       int    `json:"eurusChainId"`
	MainnetRPCDomain                   string `json:"mainnetRPCDomain"`
	MainnetRPCPort                     int    `json:"mainnetRPCPort"`
	MainnetRPCProtocol                 string `json:"mainnetRPCProtocol"`
	MainnetChainId                     int    `json:"mainnetChainId"`
	ExternalSmartContractConfigAddress string `json:"externalSmartContractConfigAddress"`
	EurusInternalConfigAddress         string `json:"eurusInternalConfigAddress"`
}

type QueryServerConfigRequest struct {
	request.RequestBase
}

func NewQueryServerConfigRequest() (me *QueryServerConfigRequest) {
	req := new(QueryServerConfigRequest)
	req.RequestPath = RootPath + GetServerConfigPath
	req.Method = http.MethodGet
	return req
}

type RegisterDeviceRequest struct {
	request.RequestBase
}

func NewRegisterDeviceRequest() *RegisterDeviceRequest {
	req := new(RegisterDeviceRequest)
	req.RequestPath = RootPath + RegisterDevicePath
	req.Method = http.MethodPost
	return req
}

type RegisterDeviceResponse struct {
	Code string `json:"code"`
}

type VerifyDeviceRequest struct {
	request.RequestBase
	Code      string `json:"code"`
	DeviceId  string `json:"deviceId"`
	PublicKey string `json:"publicKey"`
}

func NewVerifyDeviceRequest() *VerifyDeviceRequest {
	req := new(VerifyDeviceRequest)
	req.RequestPath = RootPath + VerifyDevicePath
	req.Method = http.MethodPost
	return req
}

type VerifyDeviceResponse struct {
	Mnemonic string `json:"mnemonic"`
}

type KYCCountryCodeRequest struct {
	request.RequestBase
}

func NewKYCCountryCodeRequest() *KYCCountryCodeRequest {
	req := new(KYCCountryCodeRequest)
	req.RequestPath = kyc_const.RootPath + kyc_const.UserServerPath + kyc_const.EndPoint.GetKYCCountryList
	req.Method = http.MethodGet
	return req
}

type KYCRequest struct {
	request.RequestBase
}

func NewKYCCRequest() *KYCCountryCodeRequest {
	req := new(KYCCountryCodeRequest)
	req.RequestPath = kyc_const.RootPath + kyc_const.UserServerPath
	req.Method = http.MethodGet
	return req
}

type RequestMerchantRefundRequest struct {
	request.RequestBase
	MerchantId        uint64 `json:"merchantId"`
	AssetName         string `json:"assetName"`
	Amount            uint64 `json:"amount"`
	PurchaseTransHash string `json:"purchaseTransHash"`
	Reason            string `json:"reason"`
}

func NewRequestMerchantRefundRequest() *RequestMerchantRefundRequest {
	req := new(RequestMerchantRefundRequest)
	req.RequestPath = RootPath + MerchantServerPath + RequestMerchantRefundPath
	req.Method = http.MethodPost
	return req
}

type RequestMerchantRefundResponse struct {
	RequestId uint64 `json:"refundRequestId"`
}

type QueryMerchantRefundStatusRequest struct {
	request.RequestBase
}

func NewQueryMerchantRefundStatusRequest() *QueryMerchantRefundStatusRequest {
	req := new(QueryMerchantRefundStatusRequest)
	req.RequestPath = RootPath + MerchantServerPath + MerchantRefundStatusPath
	req.Method = http.MethodGet
	return req
}

type UserReward struct {
	RewardType reward.TokenDistributedType `json:"rewardType"`
	AssetName  string                      `json:"assetName"`
	Amount     *big.Int                    `json:"amount"`
	TxHash     string                      `json:"txHash"`
	CreateDate time.Time                   `json:"createdDate"`
}

type QueryRewardListFullResponse struct {
	response.ResponseBase
	Data []*UserReward `json:"data"`
}

type RewardSchemeRequest struct {
	request.RequestBase
}

type RewardSchemeResponse struct {
	Data string `json:"data"`
}

func (me *RewardSchemeResponse) MarshalJSON() ([]byte, error) {
	return []byte(me.Data), nil
}

type GetWalletAddressRequest struct {
	request.RequestBase
	LoginAddress string `json:"loginAddress" validate:"required,hexadecimal"`
}

type GetWaleltAddressResponse string

type TopUpPaymentWalletRequest struct {
	request.RequestBase
	Signature       string `json:"signature"`
	TargetGasAmount uint64 `json:"targetGasAmount"`
}

func NewTopUpPaymentWalletRequest() *TopUpPaymentWalletRequest {
	req := new(TopUpPaymentWalletRequest)
	req.Method = http.MethodPost
	req.RequestPath = "/user/topUpPaymentWallet"
	return req
}

type TopUpPaymentWalletResponse struct {
	Tx               string `json:"tx"`
	EstimatedGasUsed uint64 `json:"estimatedGasUsed"`
}

type AssetAddressListResponse struct {
	Data string
}

func (me *AssetAddressListResponse) MarshalJSON() ([]byte, error) {
	return []byte(me.Data), nil
}

type BlockCypherAccessTokenRequest struct {
	request.RequestBase
	Coin string `json:"coin"`
}

func NewBlockCypherAccessTokenRequest() *BlockCypherAccessTokenRequest {
	req := new(BlockCypherAccessTokenRequest)
	req.Method = http.MethodGet
	req.RequestPath = "/user/blockCypher/accessToken"
	return req
}

type BlockCypherAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type MarketingBannerRequest struct {
	request.RequestBase
	Position uint64 `json:"position,string" validate:"required"`
	LangCode string `json:"langCode"`
}

type QueryMarketingBannerResponse []QueryMarketingBanner

type QueryMarketingBanner struct {
	Seq         uint64      `json:"seq"`
	BannerType  uint64      `json:"bannerType"`
	Icon        Icon        `json:"icon"`
	Content     string      `json:"content"`
	BannerImage BannerImage `json:"bannerImage"`
	Link        Link        `json:"link"`
}

type Icon struct {
	Mobile string `json:"mobile"`
}

type BannerImage struct {
	Mobile string `json:"mobile"`
}

type Link struct {
	Mobile string `json:"mobile"`
}

type QueryMarketingBannerList struct {
	Id              uint64                       `json:"id" gorm:"primaryKey"`
	Position        uint64                       `json:"position"`
	Seq             uint64                       `json:"seq"`
	BannerType      uint64                       `json:"bannerType"`
	IconUrlMobile   string                       `json:"iconUrlMobile"`
	LinkMobile      string                       `json:"linkMobile"`
	Status          banner.MarketingBannerStatus `json:"status"`
	StartDate       *time.Time                   `json:"startDate"`
	EndDate         *time.Time                   `json:"endDate"`
	IsDefault       bool                         `json:"isDefault"`
	BannerId        uint64                       `json:"bannerId"`
	LangCode        string                       `json:"langCode"`
	BannerUrlMobile string                       `json:"bannerUrlMobile"`
	Content         string                       `json:"content"`
}
