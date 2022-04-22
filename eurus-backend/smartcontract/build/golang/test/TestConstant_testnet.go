//go:build unittest_testnet
// +build unittest_testnet

package test

func InitEnvironment() {

	environment = "TESTNET"

	invokerAddrList = []string{
		"0x5c28178B275657c0B71f1513A239705605988BD4",
	}

	mainnetEthClientIP = "13.213.74.69"
	mainnetEthClientPort = 8545
	mainnetEthClientProtocol = "http"

	withdrawobserverAddr = []string{
		// "0xec219497501db463efae4c6453353cf9d587339a", //depreciated
		// "0xc06408ec9fce92517d024fe5b393ee3cfa89a22e",
		// "0xc9ba618c82e52bee35498e7c3beb61fbe78b2dd4",
		// "0xb629a6264a1bea9511442c2b7eda5707af11396f",
		// "0xa77b628a2e256202ada303b9b7c51f552f29e7b6",
		// "0x9ac3882df2b3393b27740a9db4b910745f644814",
		// "0xed4f044a4bd134102f860036eae8b7e98d479a96",
		"0xe3983a340edd00554213bfa76bd50744c40b0f19",
		"0x67a156ad4cd7585b59ac69b01763204e53d43896",
		"0x82dd783215087a9e94fc3fe27a48d6239d4d755f",
		"0xdf2d4754180813552b168af1161c60bb381470ab",
		"0xd13a6bbde085948631ce37996578f5e3de920179",
		"0x20d2fdf4be93d7ccada10ca166c3196f3ab84e83",
		"0x1ee80ade7945cf2c34e22697f4a4765d472a8685",
	}
	depositObserverAddr = []string{
		// "0x12684b8f168d9382dd9febd79f15437251515e1c", //depreciated
		// "0x6f42191cb9e05dc59b40fac59dc95399cbb04d1a",
		// "0x87474410d8ca6090d646c1145405e785bdb00704",
		// "0x5584d5d256001ec72c04250d75cf7ad29fa1b7ad",
		// "0x0e064b7b7820d6b34978a5c6e1a8d85820021fed",
		// "0x0d989d8d32384ae1732275010c313c06e22c4855",
		// "0x1ed4fd1f54c0ceac65d10009ac769475f0ac753f",
		"0x590daaef5d18eefba4ec613334dbd24e1180ef72",
		"0xd8e1d86e80eb90e4d1979709605533f91334e8b1",
		"0xfe49a07ee1f78f203522f5dd0e9922b00264063c",
		"0x271aaa0bbd210c831daeecc053ef4a747837ba68",
		"0x7f8db2161222bd1c4453a53d84f6e5eee01c42d6",
		"0x00427714678500c9054957fc381b16857f480134",
		"0x34c2ceb31dcbea5deb0ee7b96a6b0c4952a51fd3",
	}
	sweepServiceAddr = []string{
		"0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93",
	}
	approvalObserverAddr = []string{
		// "0xb0b56d86ca25cbd1708e15fea5e4103e67502d7e", //depreciated
		// "0x786027c99efcacdca741ccd59e3273bd24d973e1", //depreciated
		"0x75f8ef5ceca61643937e361ed2f157f9a23e81fa",
		"0x91d51001a2b171ee65882cfed078ed1a91d9ca32",
	}

	userObserverAddr = []string{
		"0x3cde30cecc7224819cb931fd2d2e8481fd3693d6",
	}

	userServerHDWalletAddr = []string{
		"0x66c0cb3dece53221d06bf4ff1027d3b858df2980",
	}
	configServerHDWalletAddr = []string{
		"0xce16ca3fdbe1ccc79927ccc19d8cc190e9fbb05f",
	}

	gasFeeCollectWalletAddr = "0xe19458B01AD05b81Ad27f1241494c70B9EFDaac6"
	sweepInvokerAddr = "0x5c28178b275657c0b71f1513a239705605988bd4"
	kycServerAddr = []string{
		"0x7909f2d6508c1fc08e11e06726831133230abc7f",
	}

	blockChainIndexerHDWalletAddr = []string{
		"0x1789ee6413ce7be0ae6b10ad151c7ff49eace2e1",
	}

	ethClientIP = "testnet.eurus.network"
	ethClientPort = 443
	ethClientProtocol = "https"

	chainId = "1984"
	testNetChainId = "4"
	testOwnerPrivateKey = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
	testOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
	testWithdrawApproverPrivateKey = "NA"
	testWithdrawApproverAddr = "NA"
	testWithdrawObserverPrivateKey = "NA"
	testWithdrawObserverAddr = "NA"
	testTestNetPlatformWalletOwnerPrivateKey = "NA"
	testTestNetPlatformWalletOwnerAddr = "NA"
	testTestNetOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
	testTestNetOwnerAddr = "0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"
	// testWithdrawObserverPrivateKey2 = "NA"
	// testWithdrawObserverAddr2 = "NA"
	// testDepositObserverAddr = "NA"
	smartContractFileName = "SmartContractDeploy_testnet.json"
}
