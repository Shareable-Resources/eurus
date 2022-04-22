//go:build unittest_dev
// +build unittest_dev

package test

import (
	"fmt"
)

func InitEnvironment() {

	environment = "DEV"

	invokerAddrList = []string{
		"0x5c28178B275657c0B71f1513A239705605988BD4",
	}

	mainnetEthClientIP = "13.213.74.69"
	mainnetEthClientPort = 8545
	mainnetEthClientProtocol = "http"

	approvalObserverHDWalletAddr = []string{
		"0xa82f5f0b7ff17d333ff5653c5ac3f8f97787c2cf",
		"0x897ceb7e5ff76d0872fe22b87e13cf83393462da",
		"0x30548ad0f14af72a8a130b8900b3dbd4e7a4b5f5",
		"0x5633fad91f3096cd32a3d8fe5e5bd31921d96c94",
	}
	approvalObserverAddr = []string{
		"0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93",
		"0xC0aEa01F3C3FD62F8b099Ca23f3465E9a3273bb8",
		"0x3a12786e56FF5158959Efc89E14279D1E0876eF9",
		"0x8cd4C94ab2038629F70BfA5F0CAf82BfD006686A",
	}
	withdrawobserverAddr = []string{
		"0x8ecCA8D0b23182CACE1F22E10E9c4B9CC749073A",
		"0xf143230E42f810d22Cba5eEe7d548bd338510655",
		"0xB86e3B602B52A073c1577e46b589ab499e4fbF49",
		"0x7D23B1999d72C775Be8ac76cA31aCbFc7c0ce3EC",
		"0x0a87D461f3f128d6891d16936e4fABeCCB1933cF",
		"0x81df49B3aeDc47b2f94C881C8fA9981D840Ab235",
		"0x8D27c8fb371A6cF18c6Fdf85F28BAd5A172875fF",
		"0xd3dFF32BdC06BBFF6Db9E1c94A1c64e77D8cef20",
		"0x73DAec53c8669Ce5E8E08789b4BE1E58B804a333",
		"0xA345973Fa94E786Dc483E96Fa73c5fA232d11858",
		"0xEd0A0030d7C4786D760b5a917d7563aEA2F37a3e",
		"0xF598BB24D95F03157f11D4a84Eb06b5d640b77E2",
		"0x40427d2d3c7be8BB78FC1459dbC6b4d38Dab989f",
		"0x69ee00BF2a3dFF7DD9f6eE7D04f44E87FB8b917D",
	}
	depositObserverAddr = []string{
		"0x81ccE1BF6DCfCeF814914E66d1aD6ACcFeB3936e",
		"0x80B5aD3d6370f9C4b6b84a222f42CD3FC5755D81",
		"0x96ad0B6DdD7CEb4e231dA09F7AE5044b98De5362",
		"0xA14Ab677452e5c038067f6dF57eE073328FbD336",
		"0xEECd77226192B0Cf9CE510A254F31f03c3CD2D85",
		"0x6612457C18AB8A7B1caF25bd316401359132B5cd",
		"0xdF86024284b12A1a274183A920cFf1419D524F17",
		"0x206803F71C5eB93C61a060eeCfc94C2CE6D9747D",
		"0x15C9894F2C32A892A9965Fa23b6026Abd7450Ac8",
		"0x7076cE3E9C32c6D9F625E142556E0944BBDAe6a0",
		"0xffa498003B1b1AfC2804533fF7Fb11a5ACC059C2",
		"0xa4fC063Be769a24886d0fe219d9F1A24b09dAbC2",
		"0x38455a54cC578505b6a57Aa0957d1BA9B400c9CA",
		"0xb393Edd0Bf261C49C8A5C3fb21B23DE8753BE967",
	}
	sweepServiceAddr = []string{
		"0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93",
	}
	userServerHDWalletAddr = []string{
		"0x8c9e314f4bd8dacde1ca4e9fd63064f7cd3388fa",
		"0x7EBcc140Ceac70e93fFe5B6a205A1a693AB51109",
		"0x135d622c27b0a9f59a438133e1e0fb6debfff5ff",
	}
	configServerHDWalletAddr = []string{
		"0x8c9e314f4bd8dacde1ca4e9fd63064f7cd3388fa",
	}

	userObserverAddr = []string{
		"0x7317c5684e7c93dda6191305415a323aab363efb",
	}

	signServerAddr = "0x48bf9a25520aeea2d56e540215f6cf3c9323f384"

	gasFeeCollectWalletAddr = "0xe19458B01AD05b81Ad27f1241494c70B9EFDaac6"
	sweepInvokerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"

	kycServerAddr = []string{
		"0x94de047cbc5dfe218a8cc4f33ae8df698ac493fe",
	}

	blockChainIndexerHDWalletAddr = []string{
		"0x50e29bf9619c69f1ae417405a8fb4692cab50d41",
	}

	if environment == "STAGING" {
		ethClientIP = "13.228.80.104"
		ethClientPort = 8545
		ethClientProtocol = "http"

		chainId = "2018"
		testNetChainId = "4"
		testOwnerPrivateKey = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
		testOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"

		testWithdrawApproverPrivateKey = "d38e80c2b812d1a116a53acb9b9cf208b35e5bc5600ad7a2b8ac34b06228b5ec"
		testWithdrawApproverAddr = "0x959bc5245b1dc260daa37e5a3ecbfd2097176da0"
		testWithdrawObserverPrivateKey = "92bd179b40e3fa7464853e3e5d989d8aaf5a73d6767f241f46efe51795e664dc"
		testWithdrawObserverAddr = "0x8ecCA8D0b23182CACE1F22E10E9c4B9CC749073A"
		testTestNetPlatformWalletOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		testTestNetPlatformWalletOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
		testTestNetOwnerPrivateKey = "906b9fa7fab30a174fae300b9580906738b68b28afed1a47aaf0a6d7e420cb9d"
		testTestNetOwnerAddr = "0x11F892B0230abB6B8cB3d484f8D1CE6b31d7f31B"

		// testWithdrawObserverPrivateKey2 = "359655ab07e5c29f16ad8183aa1a381988203a4424e8744db90a4d486070f35c"
		// testWithdrawObserverAddr2 = "0xcd78F2486911A6a97b3c03f140acc1B70ffAA274"
		// testDepositObserverAddr = "TBC"
	} else if environment == "DEV" {
		ethClientIP = "13.228.169.25"
		ethClientPort = 10002
		ethClientProtocol = "http"

		chainId = "2021"
		testNetChainId = "4"
		testOwnerPrivateKey = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
		testOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"

		userWalletSCOwnerPrivateKey = "5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"
		testUserWalletSCOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"

		testWithdrawApproverPrivateKey = "d38e80c2b812d1a116a53acb9b9cf208b35e5bc5600ad7a2b8ac34b06228b5ec"
		testWithdrawApproverAddr = "0x959bc5245b1dc260daa37e5a3ecbfd2097176da0"
		testWithdrawObserverPrivateKey = "92bd179b40e3fa7464853e3e5d989d8aaf5a73d6767f241f46efe51795e664dc"
		testWithdrawObserverAddr = "0x8ecCA8D0b23182CACE1F22E10E9c4B9CC749073A"
		testTestNetPlatformWalletOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		testTestNetPlatformWalletOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
		testTestNetOwnerPrivateKey = "906b9fa7fab30a174fae300b9580906738b68b28afed1a47aaf0a6d7e420cb9d"
		testTestNetOwnerAddr = "0x11F892B0230abB6B8cB3d484f8D1CE6b31d7f31B"

		// testTestNetOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		// testTestNetOwnerAddr = "0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"

		// testWithdrawObserverPrivateKey2 = "359655ab07e5c29f16ad8183aa1a381988203a4424e8744db90a4d486070f35c"
		// testWithdrawObserverAddr2 = "0xcd78F2486911A6a97b3c03f140acc1B70ffAA274"
		// testDepositObserverAddr = "TBC"

		smartContractFileName = "SmartContractDeploy_dev.json"

	} else if environment == "LOCAL" {
		chainId = "5777"
		testNetChainId = "4"
		//duncan
		testOwnerPrivateKey = "a3aa9df0858438609310921dfeedc85b8e6b6d7334a45c804bd10af645f8b05e"
		testOwnerAddr = "0xc9469463843998b0a0173f539dcfc1e30c508989"
		testWithdrawApproverPrivateKey = "3a110a66352fc217f959965eab03ce97407db9c011457e55c327e0f130a2cdc5"
		testWithdrawApproverAddr = "0x1a1946ffea82a52cdd518cd4d16d56c288fb3b7d"
		testWithdrawObserverPrivateKey = "92bd179b40e3fa7464853e3e5d989d8aaf5a73d6767f241f46efe51795e664dc"
		testWithdrawObserverAddr = "0x8ecCA8D0b23182CACE1F22E10E9c4B9CC749073A"
		testTestNetPlatformWalletOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		testTestNetPlatformWalletOwnerAddr = "0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"
		testTestNetOwnerPrivateKey = "906b9fa7fab30a174fae300b9580906738b68b28afed1a47aaf0a6d7e420cb9d"
		testTestNetOwnerAddr = "0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"
		// testWithdrawObserverPrivateKey2 = "c6e4570315c659410fa52b0e976b27606db390e134ca8bf953679191ba53e335"
		// testWithdrawObserverAddr2 = "0x994c0f7ee9d027459300b494413944e8ec19293f"
		// testDepositObserverAddr = "0x6267B266D0895D8c9EAb798526F85D1308f59657"
		// testOwnerPrivateKey = "d4ff8434d368f2eecab60189779c285b651316f06be17e8dfbc7a9d964eb0959"
		// testOwnerAddr = "0xf073eb9c8da0bdb375e15ede7388346a41a1795d"
		// testWithdrawApproverPrivateKey = "ce22eddc6f67c94365bc3d886fae3a55541ec52f5aa339b3d1f5400fc166d9b3"
		// testWithdrawApproverAddr = "0x6b4c028f65b9c9c45429f91c3330180ff94964eb"
		// testWithdrawObserverPrivateKey = "f08604b474bb9adaeaa14d4f1aac735537921c6ce22d934bbbd0d619d857f7df"
		// testWithdrawObserverAddr = "0x20dcad9edd74d840907215b1f9351ef4896bb7cf"
		// testTestNetOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		// testTestNetOwnerAddr = "0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"
		// testWithdrawObserverPrivateKey2 = "d8fbe03c34e4b5fb0063e721249416dcb76409c70a9cb9230308b147a6e3e566"
		// testWithdrawObserverAddr2 = "0xecfb0f1d451ead7bfec7bd0e45ec849ea6e4c116"
		// testDepositObserverAddr = "0x6267B266D0895D8c9EAb798526F85D1308f59657"

		// testOwnerPrivateKey = "23e199cf3d7aa4baf73bc3047e8b34ba4b34e371b7acbdfc44af5a999918c38d"
		// testOwnerAddr = "0x2587a266c64a19aa90f11fecf90e8bf8b45bd2b9"
		// testWithdrawApproverPrivateKey = "c39df2e44891712106dd45050616162897d94c1b514823ac94958e5bef2db8e0"
		// testWithdrawApproverAddr = "0x959bc5245b1dc260daa37e5a3ecbfd2097176da0"
		// testWithdrawObserverPrivateKey = "f08604b474bb9adaeaa14d4f1aac735537921c6ce22d934bbbd0d619d857f7df"
		// testWithdrawObserverAddr = "0x20dcad9edd74d840907215b1f9351ef4896bb7cf"
		// testTestNetOwnerPrivateKey = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10"
		// testTestNetOwnerAddr = "0xa5bD66B90c9F4175F3baf3dD25155Fd31543eF81"
		// testWithdrawObserverPrivateKey2 = "c6e4570315c659410fa52b0e976b27606db390e134ca8bf953679191ba53e335"
		// testWithdrawObserverAddr2 = "0x994c0f7ee9d027459300b494413944e8ec19293f"
		// testDepositObserverAddr = "0x6267B266D0895D8c9EAb798526F85D1308f59657"
	} else {
		fmt.Println("Invalid environment value")
	}
}
