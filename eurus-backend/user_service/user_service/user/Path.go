package user

var RootPath = "/user"

var GetClientVersionPath = "/clientVersion"
var GetServerConfigPath = "/serverConfig"
var KYCServerPath = "/kyc"
var MerchantServerPath = "/merchant"
var MarketingPath = "/marketing"
var BlockCypherPath = "/blockCypher"

//account related
var LoginBySignaturePath = "/loginBySignature"
var ImportWalletPath = "/importWallet"
var GetUserDetailsPath = "/details"
var RefreshTokenPath = "/refreshToken"
var RegisterPath = "/registerByEmail"
var VerificationPath = "/verification"
var SetupPaymentWalletPath = "/setupPaymentWallet"
var ResendVerificationEmailPath = "/resendVerificationEmail"
var RequestLoginRequestTokenPath = "/requestLoginRequestToken"
var RequestLoginTokenByLoginRequestTokenPath = "/requestLoginToken"
var RequestPaymentLoginTokenPath = "/requestPaymentLoginToken"

var RequestChangePaymentPasswordPath = "/changePaymentPassword"
var RequestChangeLoginPasswordPath = "/changeLoginPassword"

var ForgetLoginPasswordPath = "/forgetLoginPassword"
var VerifyForgetLoginPasswordPath = "/verifyForgetLoginPassword"
var ResetLoginPasswordPath = "/resetLoginPassword"
var ForgetPaymentPasswordPath = "/forgetPaymentPassword"
var VerifyForgetPaymentPasswordPath = "/verifyForgetPaymentPassword"
var ResetPaymentPasswordPath = "/resetPaymentPassword"

var FindEmailWalletAddressPath = "/findEmailWalletAddress"

var RegisterDevicePath = "/registerDevice"
var VerifyDevicePath = "/verifyDevice"
var GetUserWalletAddressPath = "/userWalletAddress"

//transaction related
var GetRecentTransactionPath = "/recentTransaction"
var SignTransactionPath = "/signTransaction"
var TopUpPaymentTransactionPath = "/topUpPaymentWallet"

//withdraw admin fee
var GetWithdrawAdminFeePath = "/withdrawAdminFee/{curruncySymbol}"

var UserPreferenceStoragePath = "/storage"

//Faucet
var FaucetPath = "/testnet/asset/faucet/{curruncySymbol}"
var GetFaucetConfigPath = "/testnet/asset/faucet"

//kyc related
var GetKYCCountryListPath = "/getKYCCountryList"
var CreateKYCStatusPath = "/createKYCStatus"
var SubmitKYCApprovalPath = "/submitKYCApproval"
var GetKYCStatusByTokenPath = "/getKYCStatusByToken"
var SubmitKYCDocumentPath = "/submitKYCDocument"

//merchant related
var RequestMerchantRefundPath = "/requestRefund"
var MerchantRefundStatusPath = "/refundStatus"

//Marketing
var RewardListPath = "/rewardList"
var RewardSchemePath = "/rewardScheme"
var MarketingBanner = "/marketing/banner"

//Asset list
var AssetAddressListPath = "/assetAddressList/{chain}"

//Block cypher
var BlockCypherAccessTokenPath = "/accessToken"
