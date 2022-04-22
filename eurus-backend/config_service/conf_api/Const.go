package conf_api

import "eurus-backend/foundation/network"

var ConfigExchangeMetaData = network.MQExchangeMetaData{
	ExchangeName: "config_server_exchange",
	IsAutoDelete: true,
}

type ServiceGroupId int32

const (
	ServiceGroupUnknown           ServiceGroupId = iota //0
	ServiceGroupDeposit                                 //1 - has wallet (EUN)
	ServiceGroupWithdraw                                //2 - has wallet (EUN)
	ServiceGroupApproval                                //3 - has wallet (EUN)
	ServiceGroupUser                                    //4 - has wallet (EUN)
	ServiceGroupBlockchainIndexer                       //5
	ServiceGroupUserObserver                            //6 - has wallet (EUN)
	ServiceGroupKYC                                     //7
	ServiceGroupConfig                                  //8 - has wallet (EUN)
	ServiceGroupSign                                    //9 - has wallet, but only monitor user wallet and invoker, no need to monitor the HD wallet
	ServiceGroupAdmin                                   //10

	//For table wallet balance
	WalletBalanceGroupMainHotWallet     = 90 //90
	WalletBalanceGroupMainColdWallet    = 91 //91
	WalletBalanceGroupSignInvoker       = 92 //92
	WalletBalanceGroupSignUserWallet    = 93 //93
	WalletBalanceGroupUserWallet        = 94 //94
	WalletBalanceGroupAdminFeeWallet    = 95 //95 Wallet to collect admin fee
	WalletBalanceGroupGasFeeWallet      = 96 //96 Wallet to collect centralized user EUN gas
	WalletBalanceGroupSweepServerWallet = 97 //97 Sweep server wallet in mainnet
	WalletBalanceMarketRegWallet        = 98
)
