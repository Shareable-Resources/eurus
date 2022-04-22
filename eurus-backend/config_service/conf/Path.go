package conf

var RootPath string = "/config"

//Config Related
var GetServerSettingPath string = "/setting"
var GetConfigPath string = "/"
var UpdateConfigPath string = "/update"
var DelConfigPath string = "/{serverId}/{key}"
var DelGroupConfigPath string = "/group/{serverId}/{key}"
var QueryServiceGroupAddressPath = "/group/address/{groupId}"
var GetServiceGroupIdPath = "/group/getGroupId"

//Config and Auth
var GetConfigAuthPath string = "/auth/all"

var GetSystemConfigPath string = "/system/{key}"
var AddOrUpdateSystemConfigPath string = "/system/{key}"

//Auth Related
var GetAuthPath string = "/auth"
var AddAuthPath string = "/auth/insert"
var UpdateAuthPath string = "/auth/update"
var DelAuthPath string = "/auth/{id}"
var SetServerWalletAddressPath string = "/auth/setWalletAddress"

//Mint Burn Related
var AddMintConfigPath string = "/mint/insert"
var AddBurnConfigPath string = "/burn/insert"
var AddApprovalConfigPath string = "/approval/insert"
var AddUserConfigPath string = "/user/insert"

//SC Currency Related
var GetERC20SCPath string = "/asset/currency"
var AddAssetPath string = "/asset/currency"
var DelAssetPath string = "/asset/currency/{asset}"
var GetERC20NamePath string = "/asset/currency/{assetAddress}"

//SC Public Currency Related
var GetPublicERC20SCPath string = "/asset/public-currency"
var AddPublicAssetPath string = "/asset/public-currency"
var DelPublicAssetPath string = "/asset/public-currency/{asset}"
var GetPublicERC20NamePath string = "/asset/public-currency/{assetAddress}"

//SC Wallet Related
var SetInnetWalletPath string = "/asset/innet-wallet"
var GetInnetWalletPath string = "/asset/innet-wallet"
var SetMainnetWalletPath string = "/asset/mainnet-wallet"
var GetMainnetWalletPath string = "/asset/mainnet-wallet"
var SetColdWalletPath string = "/asset/cold-wallet"

//Asset part
var SetWithdrawFeeETH string = "/asset/withdraw-fee-eth"
var SetExchangeRate string = "/asset/exchange-rate"
var SetAsset string = "/asset/public-currency-db"
var GetAsset string = "/asset/public-currency-db"

//Setting Related
var GetFaucetConfigPath string = "/asset/faucet-config"

var GetAssetSettingsPath = "/asset/settings"

//KeyPair generation
var GetKeyPairPath string = "/key"
