package test

import (
	"encoding/json"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var addrList = []string{
	"0x1815F0C6Ce4453663cBED42637A054CC2392d90c",
}

func TestOwnedUpgradeabilityProxy_SideChainUpgradeToUserWallet(t *testing.T) {
	ethClient := initEthClient(t)

	// addrList := []string{
	// 	"0x88b7c329fb4531592b18664d2ab02b0f2903bae9",
	// }

	for _, addr := range addrList {
		fmt.Println("Upgrade address: ", addr)
		proxyAddr := common.HexToAddress(addr)

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey failed for address: ", addr, " error: ", err)
			continue
		}

		userWalletAddr := getAddressBySmartContractName("UserWalletProxy")

		proxySC, err := contract.NewOwnedUpgradeabilityProxy(proxyAddr, ethClient.Client)
		if err != nil {
			fmt.Println("NewOwnedUpgradeabilityProxy failed for address: ", addr, " error: ", err)
			continue
		}

		tx, err := proxySC.UpgradeTo(transOpt, userWalletAddr)
		if err != nil {
			fmt.Println("UpgradeTo failed for address: ", addr, " error: ", err)
			continue
		}
		fmt.Println("Tx hash: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("QueryEthReceipt failed for address: ", addr, " error: ", err)
			continue
		}
		receiptJson, _ := json.Marshal(receipt)
		fmt.Println(string(receiptJson))
	}
}

func TestOwnedUpgradeabilityProxy_MainnetGetProxyOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)

	addr := common.HexToAddress("0xdd783c8Ca9E1335a0Ef0e38b6385258eC3236203")

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	impl, err := proxy.ProxyOwner(&bind.CallOpts{})
	if err != nil {
		t.Fatal("ProxyOwner error: ", err)

	}

	fmt.Println("ProxyOwner: ", impl.Hex())
}

func TestOwnedUpgradeabilityProxy_MainnetSetProxyOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)

	addr := common.HexToAddress("0xdd783c8Ca9E1335a0Ef0e38b6385258eC3236203")

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	newOwner := getTestNetAddressBySmartContractName("GeneralMultiSigWallet")
	fmt.Println("New owner: ", newOwner.Hex())
	fmt.Println("Wait for 10 seconds")
	time.Sleep(10 * time.Second)
	tx, err := proxy.TransferProxyOwnership(transOpt, newOwner)
	if err != nil {
		t.Fatal("ProxyOwner error: ", err)
	}
	fmt.Println("Tx hash: ", tx.Hash().String())
	queryEthReceipt(t, &ethClient, tx)

}

func TestOwnedUpgradeabilityProxy_MainnetGetImplementation(t *testing.T) {
	ethClient := initTestNetEthClient(t)

	addr := common.HexToAddress("0xDc322792e3a5481692a8D582E500F6588962993b")

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	impl, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal("Implementation error: ", err)

	}

	fmt.Println("Implementation: ", impl.Hex())
}

func TestOwnedUpgradeabilityProxy_SideChainGetImplementation(t *testing.T) {
	ethClient := initEthClient(t)

	addr := common.HexToAddress("0xc6583F3cBd271D0Cc99cdDeB7F74FA9879AC9b61")

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewOwnedUpgradeabilityProxy error: ", err)
	}
	impl, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal("Implementation error: ", err)

	}

	fmt.Println("Implementation: ", impl.Hex())
}

func TestOwnedUpgradeabilityProxy_SideChainUpgradeToUserWalletProxy(t *testing.T) {
	ethClient := initEthClient(t)
	// addrList := []string{
	// 	"0x88b7c329fb4531592b18664d2ab02b0f2903bae9",
	// }

	for _, addr := range addrList {

		proxySC, err := contract.NewOwnedUpgradeabilityProxy(common.HexToAddress(addr), ethClient.Client)
		if err != nil {
			fmt.Println("NewOwnedUpgradeabilityProxy failed for address: ", addr, " error: ", err)
			continue
		}

		userWalletProxyAddr := getAddressBySmartContractName("UserWalletProxy")

		transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey failed for address: ", addr, " error: ", err)
			continue
		}
		transOpt1.GasLimit = 10000000
		tx, err := proxySC.UpgradeTo(transOpt1, userWalletProxyAddr)
		if err != nil {
			fmt.Println("UpgradeTo failed for address: ", addr, " error: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			fmt.Println("UpgradeTo QueryEthReceipt failed for address: ", addr, " error: ", err)
			continue
		}

		receiptJson, _ := json.Marshal(receipt)
		fmt.Println("Upgrade to: ", string(receiptJson))

		userProxy, err := contract.NewUserWalletProxy(common.HexToAddress(addr), ethClient.Client)
		if err != nil {
			fmt.Println("NewOwnedUpgradeabilityProxy error: ", err, " user addr: ", addr)
			continue
		}

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " user addr: ", addr)
			continue
		}
		internalSC := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

		tx2, err := userProxy.SetInternalSCAddress(transOpt, internalSC)
		if err != nil {
			fmt.Println("SetInternalSCAddress error: ", err, " user addr: ", addr)
			continue
		}
		receipt2, err := ethClient.QueryEthReceipt(tx2)
		if err != nil {
			fmt.Println("Query receipt failed for user: ", addr, " error: ", err)
			continue
		}
		if receipt2.Status == 0 {
			fmt.Println("Query receipt failed for user: ", addr, " Status is 0")
		} else {
			fmt.Println("Set Internal Smart contract config address successfully. user: ", addr)
		}

		// transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		// if err != nil {
		// 	fmt.Println("GetNewTransactorFromPrivateKey failed for address: ", addr, " error: ", err)
		// 	continue
		// }

		// userWalletAddr := getAddressBySmartContractName("UserWallet")

		// userWalletProxy, err := contract.NewUserWalletProxy(common.HexToAddress(addr), ethClient.Client)
		// if err != nil {
		// 	fmt.Println("NewUserWalletProxy failed for address: ", addr, " error: ", err)
		// 	continue
		// }
		// transOpt.GasLimit = 10000000
		// tx, err = userWalletProxy.SetUserWalletImplementation(transOpt, userWalletAddr)
		// if err != nil {
		// 	fmt.Println("SetUserWalletImplementation failed for address: ", addr, " error: ", err)
		// 	continue
		// }
		// receipt, err = ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
		// if err != nil {
		// 	fmt.Println("SetUserWalletImplementation QueryEthReceipt failed for address: ", addr, " error: ", err)
		// 	continue
		// }

		// receiptJson, _ = json.Marshal(receipt)
		// fmt.Println("SetUserWalletImplementation: ", string(receiptJson))

	}
}

func TestOwnedUpgradeabilityProxy_MainnetChangeOwner(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	addr := getTestNetAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")

	multiOwnable, err := mainnet_contract.NewMultiOwnable(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	newOwnerAddr := getTestNetAddressBySmartContractName("GeneralMultiSigWallet")
	fmt.Println("New owner address: ", newOwnerAddr.Hex())

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := multiOwnable.AddOwner(transOpt, newOwnerAddr)
	if err != nil {
		t.Fatal(err)
	}

	queryEthReceipt(t, &ethClient, tx)

	transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(testTestNetOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	tx1, err := multiOwnable.RemoveOwner(transOpt1, common.HexToAddress(testTestNetOwnerAddr))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx1)
}

func TestOwnedUpgradeabilityProxy_MainnetGetOwnerList(t *testing.T) {
	ethClient := initTestNetEthClient(t)
	for name, addr := range testNetAddressMap {

		fmt.Println("Smart contract: ", name, " address: ", addr.Hex())

		multiOwnable, err := mainnet_contract.NewMultiOwnable(addr, ethClient.Client)
		if err != nil {
			fmt.Println("NewMultiOwnable error: ", err)
			continue
		}
		// newOwnerAddr := getTestNetAddressBySmartContractName("GeneralMultiSigWallet")
		// fmt.Println("GeneralMultiSigWallet address: ", newOwnerAddr.Hex())
		ownerList, err := multiOwnable.GetOwners(&bind.CallOpts{})
		if err != nil {
			fmt.Println("GetOwners error: ", err)
			continue
		}
		for _, addr := range ownerList {
			fmt.Println(addr.Hex())
		}
	}
}

func TestOwnedUpgradeabilityProxy_SideChainGetOwnerList(t *testing.T) {
	ethClient := initEthClient(t)
	// addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<EurusPlatformWallet>")
	addr := common.HexToAddress("0x040078fb3772e904cdae529413c635545f81b20a")
	multiOwnable, err := contract.NewMultiOwnable(addr, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}
	newOwnerAddr := getAddressBySmartContractName("GeneralMultiSigWallet")
	fmt.Println("GeneralMultiSigWallet address: ", newOwnerAddr.Hex())
	ownerList, err := multiOwnable.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	for _, addr := range ownerList {
		fmt.Println(addr.Hex())
	}
}

func TestOwnedUpgradeabilityProxy_GetAllOwnerListAtSideChain(t *testing.T) {
	ethClient := initEthClient(t)

	for name, addr := range addressMap {
		fmt.Println("Smart contract: ", name)
		multiOwnable, err := contract.NewMultiOwnable(addr, ethClient.Client)
		if err != nil {
			fmt.Println("Create object error: ", err)
			continue
		}
		ownerList, err := multiOwnable.GetOwners(&bind.CallOpts{})
		if err != nil {
			fmt.Println("Get owner failed: ", err)
			continue
		}
		for _, addr := range ownerList {
			fmt.Println(addr.Hex())
		}
		fmt.Println("")
	}
}

func TestOwnedUpgradeabilityProxy_SetAllOwnerListAtSideChain(t *testing.T) {
	ethClient := initEthClient(t)
	multiSignOwner := getAddressBySmartContractName("GeneralMultiSigWallet")

	for name, addr := range addressMap {
		if strings.HasPrefix(name, "OwnedUpgradeability") || strings.HasPrefix(name, "Old") {
			continue
		}

		if strings.Contains(name, "USDM") || strings.Contains(name, "ETHM") || strings.Contains(name, "BTCM") ||
			strings.Contains(name, "WEUN") || strings.Contains(name, "GeneralMultiSigWallet") {
			continue
		}

		fmt.Println("Smart contract: ", name)
		multiOwnable, err := contract.NewMultiOwnable(addr, ethClient.Client)
		if err != nil {
			fmt.Println("Create object error: ", err)
			continue
		}
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err)
			continue
		}
		tx, err := multiOwnable.AddOwner(transOpt, multiSignOwner)
		if err != nil {
			fmt.Println("Get owner failed: ", err)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Get ETH receipt error: ", err)
			continue
		}
		if receipt.Status == 0 {
			fmt.Println("Receipt status is 0")
			continue
		} else {
			fmt.Println("Add owner success")
		}

		fmt.Println("")
	}
}

func TestOwndUpgradeabilityProxy_RedeployUserWalletOwnedProxy(t *testing.T) {
	log.NewLogger(log.Name.Root, "/tmp/RedeployUserWallet.log", logrus.DebugLevel)
	var userWalletAddrList []string = []string{
		"0x37de2417c9873857e714f245958303a1cc32f9a0",
		"0x1f07fb9227aca24f30a338b947ffbc806a2c82d8",
		"0x4261a22e731645844d609c6f672c472e7e85b907",
		"0xd21a1175922299433deec446500d49f0d9b25859",
		"0xf27b54b0781d6c3a2174b82a474a88079411ac10",
		"0xc60b78f8b0f9022761da6c48d53ffbad394ad150",
		"0x4ad55f52c23499053ebb61ec80d1bcb6fc9e094f",
		"0x32f7ce28cd38a843adfd7660f035f7688253288d",
		"0xacc795036071d57af26816ee3e142a10799e4b41",
		"0xab52dfa82c7cb904695a55ce279e9a76541437b9",
		"0x74604d9261fd531ac91bb4c264dfa47344ecac5e",
		"0x03f9833eaaf99a22b2d4336371a79e370c673b9c",
		"0x36c144209ff52c2fc92e289348bce10a6935bf15",
		"0xdae85eca3e47ff5285d720826b86cd86059e6bfd",
		"0xd4983f30b49a8b1408058bc8a0273c89a9306960",
		"0x7adac01c9d66496d954d488b807639ab51aa6667",
		"0x6d438dded7bf03e01b35a564fcaa30062f30e233",
		"0xc4c7a6bc0981032b858be0766f95d65d9b1b6ad0",
		"0xa6a21d46fcead0d06248b1ff2ed6c6b02f5e32dd",
		"0x0979eda5586d60e49c39783fbf87f123067f80a2",
		"0x8307c5bcfa15220efc04ac2e893b43cf8faf9556",
		"0xef04bc652311b6903d91d5245cab58e9ce04130c",
		"0xf6fa8b1715230d4b8099ba5c0c323f0fd88b8387",
		"0x8db9be1301ddca68b3770e69478aefbb73836f34",
		"0xd5f269447b3238e61aa4e80a886039383a2586d7",
		"0x8828ae592433d1a85a88841f4531f8834538b727",
		"0x05862a244d1d192194db3aedd991fbd6362e011e",
		"0x9e5d09343f8e48560c8a6b5d717bd047713977f1",
		"0x443a214af82e2363d11b3cd0e5f8a429d602e8d1",
		"0xee8fc6c1b85238a49b9677f0d755310d31c29697",
		"0xd826868372501bc2a4bc21cb4ac415286b197b60",
		"0xf60265ff331cd85b00f4b29f685d2ff951eb2078",
		"0x47b5a68ecbbb84d016e19ae73ca764717b17de55",
		"0xce6ea6f47c45d0c5e0e9828a18fb8f1c78ec00eb",
		"0xe95260e3dbac9164c0f643b3cde62d81c8e5e329",
		"0xf144ab9b04af3f0cd7aa859bf52a93555840c911",
		"0x512d1365abecae947c98f33e76b10be9a1baeb96",
		"0x9f86ed154575dcbe37a3d917dc69fb14358864af",
		"0xf727449e81c8e2a62fe6db2e61fe5f4ef8501cd4",
		"0x12bf04f42de6f1c26e98c97c3cefbcc00adc8c30",
		"0x534cd112e99d45c0ea888fb858acd455ceee626d",
		"0xe6712b4e2076319a48974466477a4b29a5a142da",
		"0x9b158f5937b97ae713e1b57bee700e397b79a91d",
		"0x7dabc6d8bef126092c64311e670a66fa98b2b679",
		"0x8f23f891f4b02afec2f71a36843e4978f468af88",
		"0x4ca0c1c53e40a491a709ab35de7d5d91f604602e",
		"0x1bf01e543fec2be2bf3a22bfd9e58bb37384e868",
		"0x80d3a6cf3c6e527e8defff22f678029d6a069494",
		"0xd806dc9430115591c2229fba119c1f6dee191b78",
		"0xb4a70777f1bd82de875aa49bfdd68e6913be51df",
		"0x32e5e02546b0d9e19ad76d84892d4d3a5cb26c4c",
		"0xe474421ccaac1bf9cbe182c1d25d64508fc7574f",
		"0x3431bfeb7ae31b05335a1d54c6e7cb7670d3088c",
		"0x6beacafc60c8f34bc452b5717c8ccc1dba132ab0",
		"0x2fbd110d6624d2bc295f6fd254d185f367bc1909",
		"0x0a604a2d0e9240d5f1e10bffbc49d807577a557c",
		"0xea99d7551c60f4cb0e7f0999de8829009fbe1a46",
		"0x8835b61bca7b8f4d140b1d8a93aaf9b3f3097fdf",
		"0xcfd4040302cfb7ef35699a60f894086090ada2ae",
		"0x7a0c44991fa3e7a62f83544d1c1d3b9b78814746",
		"0x4ec0e9d44b03224e4c8f00d351d65880c89f6f66",
		"0x1f85726917b4ddbc0d1fedffa7b34e5979088ee2",
		"0x3c28719baa2b05ebb2d0ae9d88adf8f8e9aba983",
		"0x7629ea34b8de6bf2249b7435535099a286bf6390",
	}

	var walletAddrEmailMap map[string]string = map[string]string{
		"0x37de2417c9873857e714f245958303a1cc32f9a0": "michael.reichstein@gmail.com",
		"0x1f07fb9227aca24f30a338b947ffbc806a2c82d8": "ruhaeli000@gmail.com",
		"0x4261a22e731645844d609c6f672c472e7e85b907": "abhijitmakal860@gmail.com",
		"0xd21a1175922299433deec446500d49f0d9b25859": "jkeu88@gmail.com",
		"0xf27b54b0781d6c3a2174b82a474a88079411ac10": "eu82@18m.dev",
		"0xc60b78f8b0f9022761da6c48d53ffbad394ad150": "duy.tranduc2002@gmail.com",
		"0x4ad55f52c23499053ebb61ec80d1bcb6fc9e094f": "nna6450oke@gmail.com",
		"0x32f7ce28cd38a843adfd7660f035f7688253288d": "eu31@18m.dev",
		"0xacc795036071d57af26816ee3e142a10799e4b41": "eu84@18m.dev",
		"0xab52dfa82c7cb904695a55ce279e9a76541437b9": "eu09@18m.dev",
		"0x74604d9261fd531ac91bb4c264dfa47344ecac5e": "eu36@18m.dev",
		"0x03f9833eaaf99a22b2d4336371a79e370c673b9c": "bassembelahyaa@gmail.com",
		"0x36c144209ff52c2fc92e289348bce10a6935bf15": "nglingluk@gmail.com",
		"0xdae85eca3e47ff5285d720826b86cd86059e6bfd": "eu08@18m.dev",
		"0xd4983f30b49a8b1408058bc8a0273c89a9306960": "mortelmarlon4@gmail.com",
		"0x7adac01c9d66496d954d488b807639ab51aa6667": "car1298999@gmail.com",
		"0x6d438dded7bf03e01b35a564fcaa30062f30e233": "eu83@18m.dev",
		"0xc4c7a6bc0981032b858be0766f95d65d9b1b6ad0": "eu35@18m.dev",
		"0xa6a21d46fcead0d06248b1ff2ed6c6b02f5e32dd": "alisiid664@gmail.com",
		"0x0979eda5586d60e49c39783fbf87f123067f80a2": "ragilh309@gmail.com",
		"0x8307c5bcfa15220efc04ac2e893b43cf8faf9556": "eu86@18m.dev",
		"0xef04bc652311b6903d91d5245cab58e9ce04130c": "goldhub85@gmail.com",
		"0xf6fa8b1715230d4b8099ba5c0c323f0fd88b8387": "umehndukwe84@gmail.com",
		"0x8db9be1301ddca68b3770e69478aefbb73836f34": "yeti.chang@yahoo.com.hk",
		"0xd5f269447b3238e61aa4e80a886039383a2586d7": "circle.yuen@goldhub.hk",
		"0x8828ae592433d1a85a88841f4531f8834538b727": "eu05@18m.dev",
		"0x05862a244d1d192194db3aedd991fbd6362e011e": "eu11@18m.dev",
		"0x9e5d09343f8e48560c8a6b5d717bd047713977f1": "ken.li@goldhub.hk",
		"0x443a214af82e2363d11b3cd0e5f8a429d602e8d1": "circleytest01@gmail.com",
		"0xee8fc6c1b85238a49b9677f0d755310d31c29697": "circleytest02@gmail.com",
		"0xd826868372501bc2a4bc21cb4ac415286b197b60": "circleytest03@gmail.com",
		"0xf60265ff331cd85b00f4b29f685d2ff951eb2078": "circleytest01@yahoo.com",
		"0x47b5a68ecbbb84d016e19ae73ca764717b17de55": "circleytest04@yopmail.com",
		"0xce6ea6f47c45d0c5e0e9828a18fb8f1c78ec00eb": "circleytest05@yopmail.com",
		"0xe95260e3dbac9164c0f643b3cde62d81c8e5e329": "jackyjc999@gmail.com",
		"0xf144ab9b04af3f0cd7aa859bf52a93555840c911": "alan.fung@goldhub.hk",
		"0x512d1365abecae947c98f33e76b10be9a1baeb96": "yetichang2018@gmail.com",
		"0x9f86ed154575dcbe37a3d917dc69fb14358864af": "jayyson741@gmail.com",
		"0xf727449e81c8e2a62fe6db2e61fe5f4ef8501cd4": "test111501@yopmail.com",
		"0x12bf04f42de6f1c26e98c97c3cefbcc00adc8c30": "lamxxxl1234@gmail.com",
		"0x534cd112e99d45c0ea888fb858acd455ceee626d": "eurustest999@yopmail.com",
		"0xe6712b4e2076319a48974466477a4b29a5a142da": "aamirabbaskhan792@gmail.com",
		"0x9b158f5937b97ae713e1b57bee700e397b79a91d": "divyanshupanday7374@gmail.com",
		"0x7dabc6d8bef126092c64311e670a66fa98b2b679": "test111502@yopmail.com",
		"0x8f23f891f4b02afec2f71a36843e4978f468af88": "invest167@outlook.com",
		"0x4ca0c1c53e40a491a709ab35de7d5d91f604602e": "eu81@18m.dev",
		"0x1bf01e543fec2be2bf3a22bfd9e58bb37384e868": "test111502cir@yopmail.com",
		"0x80d3a6cf3c6e527e8defff22f678029d6a069494": "hichou9za@gmail.com",
		"0xd806dc9430115591c2229fba119c1f6dee191b78": "jk2255669012@gmail.com",
		"0xb4a70777f1bd82de875aa49bfdd68e6913be51df": "cb587034@gmail.com",
		"0x32e5e02546b0d9e19ad76d84892d4d3a5cb26c4c": "cathy.lam@goldhub.hk",
		"0xe474421ccaac1bf9cbe182c1d25d64508fc7574f": "eu99@18m.dev",
		"0x3431bfeb7ae31b05335a1d54c6e7cb7670d3088c": "yeti.chang@goldhub.hk",
		"0x6beacafc60c8f34bc452b5717c8ccc1dba132ab0": "peichandev@gmail.com",
		"0x2fbd110d6624d2bc295f6fd254d185f367bc1909": "kyrtarty1@gmail.com",
		"0x0a604a2d0e9240d5f1e10bffbc49d807577a557c": "c689123456@gmail.com",
		"0xea99d7551c60f4cb0e7f0999de8829009fbe1a46": "eu00@18m.dev",
		"0x8835b61bca7b8f4d140b1d8a93aaf9b3f3097fdf": "yetiland@yahoo.com.hk",
		"0xcfd4040302cfb7ef35699a60f894086090ada2ae": "eddie.abcc@gmail.com",
		"0x7a0c44991fa3e7a62f83544d1c1d3b9b78814746": "shaynejia35@gmail.com",
		"0x4ec0e9d44b03224e4c8f00d351d65880c89f6f66": "owner@yopmail.com",
		"0x1f85726917b4ddbc0d1fedffa7b34e5979088ee2": "eurusowner@yopmail.com",
		"0x3c28719baa2b05ebb2d0ae9d88adf8f8e9aba983": "lamheiyee9@gmail.com",
		"0x7629ea34b8de6bf2249b7435535099a286bf6390": "marcusyucola@gmail.com",
	}

	var addressMap map[string]string = make(map[string]string)
	var errorUserWalletList []string = make([]string, 0)
	var notExistsList []string = make([]string, 0)

	userWalletProxyAddr := getAddressBySmartContractName("UserWalletProxy")
	internalSCAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

	log.GetLogger(log.Name.Root).Infoln("User wallet proxy: ", userWalletProxyAddr.Hex())
	log.GetLogger(log.Name.Root).Infoln("Internal Smart contract config: ", internalSCAddr.Hex())
	ethClient := initEthClient(t)
	for _, userWalletAddr := range userWalletAddrList {

		oldUserWallet, err := contract.NewUserWallet(common.HexToAddress(userWalletAddr), ethClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("user wallet not exists: ", userWalletAddr, " error: ", err)
			notExistsList = append(notExistsList, userWalletAddr)
			continue
		}

		newUserWalletAddr, err := DeployUserWallet(userWalletAddr, &ethClient)
		if err != nil {
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		}

		walletOwner, err := oldUserWallet.GetWalletOwner(&bind.CallOpts{})
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetWalletOwner Error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		//Patch to use bug fixed user wallet for the old user wallet
		// userProxy, err := contract.NewUserWalletProxy(common.HexToAddress(userWalletAddr), ethClient.Client)
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("NewUserWalletProxy: ", userWalletAddr, " error: ", err)
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }
		// transOpt9, _ := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		// transOpt9.GasLimit = 400000
		// patchedUserWalletAddr := getAddressBySmartContractName("UserWallet")
		// tx9, err := userProxy.SetUserWalletImplementation(transOpt9, patchedUserWalletAddr)
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("SetUserWalletImplementation Error: ", err)
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }
		// receipt9, err := ethClient.QueryEthReceipt(tx9)
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("SetUserWalletImplementation QueryEthReceipt Error: ", err)
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }
		// if receipt9.Status != 1 {
		// 	log.GetLogger(log.Name.Root).Errorln("SetUserWalletImplementation QueryEthReceipt status = 0. Tx hash: ", tx9.Hash().Hex())
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }

		userWallet, err := contract.NewUserWallet(newUserWalletAddr, ethClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("NewUserWallet failed for account: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		}

		transOpt3, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey on SetWalletOwner Error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
		//Set wallet owner to new user wallet
		tx3, err := userWallet.SetWalletOwner(transOpt3, walletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("SetWalletOwner failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx3.Hash().Hex())
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		receipt3, err := ethClient.QueryEthReceiptWithSetting(tx3, 1, 50)
		if receipt3.Status != 1 {
			log.GetLogger(log.Name.Root).Errorln("SetWalletOwner receipt failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx3.Hash().Hex())
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		transOpt5, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
		//Adding sign invoker as
		log.GetLogger(log.Name.Root).Debugln("Add owner for ", userWalletAddr)
		tx5, err := userWallet.AddOwner(transOpt5, common.HexToAddress(signServerAddr))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddOwner failed on add owner: ", userWalletAddr, " error: ", err)
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		receipt5, err := ethClient.QueryEthReceiptWithSetting(tx5, 1, 50)
		if receipt5.Status != 1 {
			log.GetLogger(log.Name.Root).Errorln("AddOwner receipt failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx5.Hash().Hex())
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}

		//Transfer EUN to new wallet
		// balance, err := ethClient.GetBalance(common.HexToAddress(userWalletAddr))
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("GetBalance failed on add owner: ", userWalletAddr, " error: ", err)
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// }
		// log.GetLogger(log.Name.Root).Infoln("Old Balance: ", balance, " user wallet: ", userWalletAddr)
		// if balance.Cmp(big.NewInt(0)) > 0 {
		// 	log.GetLogger(log.Name.Root).Infoln("require to transfer EUN to new wallet from ", userWalletAddr, " to ", newUserWalletAddr.Hex())
		// 	//Transfer EUN
		// 	estimateBalance := balance

		// 	walletOwner, err := oldUserWallet.GetWalletOwner(&bind.CallOpts{})
		// 	if err != nil {
		// 		log.GetLogger(log.Name.Root).Errorln("GetWalletOwner failed on add owner: ", userWalletAddr, " error: ", err)
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}

		// 	err = SetWalletOwner(userWalletAddr, oldUserWallet, &ethClient, common.HexToAddress(testOwnerAddr))
		// 	if err != nil {
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}

		// 	transOpt8, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		// 	if err != nil {
		// 		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed on transfer: ", userWalletAddr, " error: ", err)
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}
		// 	transOpt8.NoSend = true
		// 	tx8, err := oldUserWallet.DirectRequestTransfer(transOpt8, newUserWalletAddr, "EUN", estimateBalance)
		// 	if err != nil {
		// 		log.GetLogger(log.Name.Root).Errorln("Estimate DirectRequestTransfer failed: ", userWalletAddr, " error: ", err)
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}

		// 	gasLimit := big.NewInt(0)
		// 	gasLimit.SetUint64(tx8.Gas())

		// 	eunUsed := tx8.GasPrice().Mul(tx8.GasPrice(), gasLimit)
		// 	final := balance.Sub(balance, eunUsed)
		// 	log.GetLogger(log.Name.Root).Infoln("Final balance: ", final.String(), " user: ", userWalletAddr)
		// 	transOpt8.NoSend = false
		// 	tx8, err = oldUserWallet.DirectRequestTransfer(transOpt8, newUserWalletAddr, "EUN", estimateBalance)
		// 	if err != nil {
		// 		log.GetLogger(log.Name.Root).Errorln("DirectRequestTransfer failed: ", userWalletAddr, " error: ", err)
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}
		// 	log.GetLogger(log.Name.Root).Infoln("DirectRequestTransfer hash: ", tx8.Hash().Hex(), " user wallet: ", userWalletAddr)
		// 	receipt8, err := ethClient.QueryEthReceipt(tx8)
		// 	if err != nil {
		// 		log.GetLogger(log.Name.Root).Errorln("DirectRequestTransfer QueryEthReceipt failed: ", userWalletAddr, " error: ", err)
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}
		// 	if receipt8.Status != 1 {
		// 		log.GetLogger(log.Name.Root).Errorln("DirectRequestTransfer receipt status = 0: ", userWalletAddr, " tx hash: ", tx8.Hash().Hex())
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}

		// 	err = SetWalletOwner(userWalletAddr, oldUserWallet, &ethClient, walletOwner)
		// 	if err != nil {
		// 		errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 		continue
		// 	}

		// }
		walletInfoHash, err := AddUserAddressToWalletSC(common.HexToAddress(userWalletAddr), walletAddrEmailMap[userWalletAddr], &ethClient, false, false)
		if err != nil {
			errorUserWalletList = append(errorUserWalletList, userWalletAddr)
			continue
		}
		log.GetLogger(log.Name.Root).Errorln("AddUserAddressToWalletSC user: ", userWalletAddr, " tx hash: ", walletInfoHash)

		// transOpt6, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed on remove owner: ", userWalletAddr, " error: ", err)
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }

		// log.GetLogger(log.Name.Root).Debugln("Remove owner for ", userWalletAddr)
		// tx6, err := userWallet.RemoveOwner(transOpt6, common.HexToAddress(testOwnerAddr))
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Errorln("RemoveOwner failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx6.Hash().Hex())
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }
		// receipt6, err := ethClient.QueryEthReceiptWithSetting(tx5, 1, 50)
		// if receipt6.Status != 1 {
		// 	log.GetLogger(log.Name.Root).Errorln("RemoveOwner receipt failed on user: ", userWalletAddr, " error: ", err, " tx hash: ", tx6.Hash().Hex())
		// 	errorUserWalletList = append(errorUserWalletList, userWalletAddr)
		// 	continue
		// }

		addressMap[userWalletAddr] = newUserWalletAddr.Hex()

	}

	var logMsg string = "Error wallet: \r\n"

	for _, errorWallet := range errorUserWalletList {
		logMsg += errorWallet + "\r\n"
	}

	logMsg += "Not exists wallet: \r\n"
	for _, notExist := range notExistsList {
		logMsg += notExist + "\r\n"
	}

	logMsg += "Migrated user wallet: "
	for oldWallet, newWallet := range addressMap {
		logMsg += fmt.Sprintf("{\"%s\":\"%s\"},\r\n", strings.ToLower(oldWallet), strings.ToLower(newWallet))
	}

	log.GetLogger(log.Name.Root).Infoln(logMsg)

}

func DeployUserWallet(oldWalletAddr string, ethClient *ethereum.EthClient) (common.Address, error) {

	var fixGasLimit uint64 = 400000
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get transOpt error: ", err, " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	}

	// somehow the address return by below function is fake, the follow real proxy address will be: receipt.ContractAddress
	proxyAddress, tx, _, err := contract.DeployOwnedUpgradeabilityProxy(transOpt, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot deploy user wallet proxy. Error: ", err, " oldWalletAddr: ", oldWalletAddr)
		return common.HexToAddress("0x0"), errors.Wrap(err, "cannot deploy user wallet proxy")
	}

	log.GetLogger(log.Name.Root).Infoln("OwnedUpgradeabilityProxy<UserWalletProxy> address: ", proxyAddress.Hex(), " tx hash:", tx.Hash().Hex(), " oldWalletAddr: ", oldWalletAddr)

	receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get deploy user wallet proxy receipt. Error:", err, " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get deploy user wallet proxy receipt")
	}
	if receipt.Status != 1 {
		receiptData, _ := receipt.MarshalJSON()
		log.GetLogger(log.Name.Root).Errorln("cannot get deploy user wallet proxy receipt: ", string(receiptData), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.New("cannot get deploy user wallet proxy receipt: " + receipt.TxHash.Hex())
	}

	var proxyContractAddress common.Address = receipt.ContractAddress

	log.GetLogger(log.Name.Root).Infoln("Deploy Proxy hash : ", receipt.TxHash, "receipt status", receipt.Status, " oldWalletAddr: ", oldWalletAddr)
	log.GetLogger(log.Name.Root).Infoln("proxy address : ", proxyContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)

	internalSCConfigAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	internalSCConfig, err := contract.NewInternalSmartContractConfig(internalSCConfigAddr, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get internalSCConfig. oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get internalSCConfig")
	}
	userWalletProxyAddr, err := internalSCConfig.GetUserWalletProxyAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get GetUserWalletProxyAddress oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get GetUserWalletProxyAddress")
	}

	userWalletAddr, err := internalSCConfig.GetUserWalletAddress(&bind.CallOpts{})
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get GetUserWalletAddress oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get GetUserWalletAddress")
	}
	transOpt2, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	//transOpt, err = server.EthClient.GetNewTransactorFromPrivateKey(server.Config.HdWalletPrivateKey, server.EthClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get transOpt oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	}

	proxy, err := contract.NewOwnedUpgradeabilityProxy(proxyContractAddress, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get proxy oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get proxy")
	}

	log.GetLogger(log.Name.Root).Infoln("user wallet proxy : ", proxyContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)
	transOpt2.GasLimit = fixGasLimit
	tx, err = proxy.UpgradeTo(transOpt2, userWalletProxyAddr)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot create user wallet upgrade proxy transaction ", err.Error(), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot create user wallet upgrade proxy transaction")
	}

	log.GetLogger(log.Name.Root).Infoln("upgrade tx: ", tx.Hash(), " oldWalletAddr: ", oldWalletAddr)
	receipt2, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get upgrade user wallet proxy receipt. oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get upgrade user wallet proxy receipt")
	}

	if receipt2.Status != 1 {
		receiptData, _ := receipt2.MarshalJSON()
		log.GetLogger(log.Name.Root).Errorln("transaction fail ", string(receiptData), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.New("cannot upgrade user wallet proxy for tx: " + receipt2.TxHash.Hex() + " failed")
	}

	// userWalletProxy, err := contract.NewUserWalletProxy(receipt.ContractAddress, ethClient.Client)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("cannot get user wallet proxy oldWalletAddr: ", oldWalletAddr)
	// 	return common.Address{}, errors.Wrap(err, "cannot get user wallet proxy")
	// }

	// transOpt7, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	// //transOpt, err = server.EthClient.GetNewTransactorFromPrivateKey(server.Config.HdWalletPrivateKey, server.EthClient.ChainID)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("cannot get transOpt oldWalletAddr: ", oldWalletAddr)
	// 	return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	// }

	log.GetLogger(log.Name.Root).Infoln("SetUserWalletImplementation oldWalletAddr: ", oldWalletAddr)
	// transOpt7.GasLimit = fixGasLimit
	// tx, err = userWalletProxy.SetUserWalletImplementation(transOpt7, userWalletAddr)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("SetUserWalletImplementation failed. ", " oldWalletAddr: ", oldWalletAddr)
	// 	return common.Address{}, errors.Wrap(err, "SetUserWalletImplementation failed")
	// }
	// log.GetLogger(log.Name.Root).Infoln("Set User Wallet Implementation tx: ", tx.Hash().Hex(), " oldWalletAddr: ", oldWalletAddr)
	// receipt7, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
	// if err != nil {
	// 	log.GetLogger(log.Name.Root).Errorln("cannot get SetUserWalletImplementation receipt. tx: ", tx.Hash().Hex(), " oldWalletAddr: ", oldWalletAddr)
	// 	return common.Address{}, errors.Wrap(err, "cannot get SetUserWalletImplementation receipt")
	// }

	// if receipt7.Status != 1 {
	// 	receiptData, _ := receipt7.MarshalJSON()
	// 	log.GetLogger(log.Name.Root).Errorln("SetUserWalletImplementation receipt failed ", string(receiptData), " oldWalletAddr: ", oldWalletAddr)
	// 	return common.Address{}, errors.New("SetUserWalletImplementation receipt failed")
	// }

	transOpt9, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey on SetInternalSmartContractConfig Error: ", err)
		return common.Address{}, err
	}

	userWallet, err := contract.NewUserWallet(receipt.ContractAddress, ethClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot new user wallet", err.Error(), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot new user wallet")
	}
	transOpt9.GasLimit = fixGasLimit
	tx9, err := userWallet.SetInternalSmartContractConfig(transOpt9, internalSCConfigAddr)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Cannot SetInternalSmartContractConfig: ", err)
		return common.Address{}, err
	}
	receipt9, err := ethClient.QueryEthReceiptWithSetting(tx9, 1, 50)
	if receipt9.Status != 1 {
		receiptData, _ := receipt9.MarshalJSON()
		log.GetLogger(log.Name.Root).Errorln("SetInternalSmartContractConfig receipt failed on user: ", userWalletAddr, " receipt: ", string(receiptData))
		return common.Address{}, errors.New("Set internal SC address receipt error")
	}

	transOpt3, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get transOpt oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot get transOpt")
	}
	transOpt3.GasLimit = fixGasLimit
	tx, err = userWallet.ChangeRequirement(transOpt3, big.NewInt(2))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot change requirement", err.Error(), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot change requirement")
	}

	log.GetLogger(log.Name.Root).Infoln("change requirement tx: ", tx.Hash(), " oldWalletAddr: ", oldWalletAddr)
	receipt3, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get change requirement receipt tx: ", tx.Hash().Hex(), " error: ", err.Error(), " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.Wrap(err, "cannot change requirement")
	}

	if receipt3.Status != 1 {
		receiptData, _ := receipt3.MarshalJSON()
		log.GetLogger(log.Name.Root).Errorln("transaction fail: ", string(receiptData), " user Id: ", oldWalletAddr)
		return common.Address{}, errors.New("cannot change user wallet requirement for tx: " + receipt3.TxHash.Hex())
	}

	var atLeastObserverSuccess bool
	for _, userObsAddr := range userObserverAddr {
		transOpt4, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddWalletOperator GetNewTransactorFromPrivateKey error: ", err, " user observer address: ", userObsAddr, " User wallet address: ", receipt.ContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)
			continue
		}
		transOpt4.GasLimit = fixGasLimit
		tx, err := userWallet.AddWalletOperator(transOpt4, common.HexToAddress(userObsAddr))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddWalletOperator error: ", err, " user observer address: ", userObsAddr, " User wallet address: ", receipt.ContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)
			continue
		}
		log.GetLogger(log.Name.Root).Infoln("AddWalletOperator tx: ", tx.Hash())
		addOpReceipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddWalletOperator QueryEthReceiptWithSetting error: ", err, " user observer address: ", userObsAddr, " User wallet address: ", receipt.ContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)
			continue
		}
		if addOpReceipt.Status == 0 {
			receiptByte, _ := addOpReceipt.MarshalJSON()
			log.GetLogger(log.Name.Root).Errorln("AddWalletOperator receipt status failed: ", string(receiptByte), " user observer address: ", userObsAddr, " User wallet address: ", receipt.ContractAddress.Hex(), " oldWalletAddr: ", oldWalletAddr)
			continue
		}
		atLeastObserverSuccess = true
	}
	if !atLeastObserverSuccess {
		log.GetLogger(log.Name.Root).Errorln("None of the user observer able to add writer. ", err, " oldWalletAddr: ", oldWalletAddr)
		return common.Address{}, errors.New("UserWallet add observer address failed. " + receipt.ContractAddress.Hex())
	}

	log.GetLogger(log.Name.Root).Infoln("Add invoker as writer. oldWalletAddr: ", oldWalletAddr)
	for _, invokderAddr := range invokerAddrList {
		transOpt6, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Add writer - GetNewTransactorFromPrivateKey error: ", err)
			return common.Address{}, nil
		}
		transOpt6.GasLimit = fixGasLimit
		tx, err = userWallet.AddWriter(transOpt6, common.HexToAddress(invokderAddr))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Add writer error: ", err, " oldWalletAddr: ", oldWalletAddr)
			return common.Address{}, nil
		}
		receipt6, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("cannot get add writer receipt: ", err.Error(), " oldWalletAddr: ", oldWalletAddr)
			return common.Address{}, err
		}
		if receipt6.Status != 1 {
			receiptData, _ := receipt6.MarshalJSON()
			log.GetLogger(log.Name.Root).Errorln("add writer failed fail: ", string(receiptData), " oldWalletAddr: ", oldWalletAddr)
			return common.Address{}, errors.New("add writer failed receipt: " + receipt6.TxHash.Hex())
		}
	}

	log.GetLogger(log.Name.Root).Infoln("Deploy smart contract finished. oldWalletAddr: ", oldWalletAddr, " smart contract address: ", receipt.ContractAddress.Hex())

	return proxyContractAddress, nil
}

func SetWalletOwner(userWalletAddr string, userWallet *contract.UserWallet, ethClient *ethereum.EthClient, walletOwner common.Address) error {

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed: ", userWalletAddr, " error: ", err)
	}
	tx, err := userWallet.SetWalletOwner(transOpt, walletOwner)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("SetWalletOwner failed: ", userWalletAddr, " error: ", err)
		return err
	}
	receipt, err := ethClient.QueryEthReceipt(tx)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("SetWalletOwner QueryEthReceipt failed: ", userWalletAddr, " error: ", err)
		return err
	}
	if receipt.Status != 1 {
		log.GetLogger(log.Name.Root).Errorln("SetWalletOwner receipt status = 0: ", userWalletAddr)
		return err
	}
	return nil
}

func AddUserAddressToWalletSC(walletAddress common.Address, email string, ethClient *ethereum.EthClient, userIsMerchant bool, userIsMetaMask bool) (string, error) {

	var fixGasLimit uint64 = 400000

	walletAddressMapAddr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	instance, err := contract.NewWalletAddressMap(walletAddressMapAddr, ethClient.Client)
	txHash := ""
	if err != nil {
		err = errors.Wrap(err, "NewWalletAddressMap error")
		return txHash, err
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		err = errors.Wrap(err, "GetNewTransactorFromSignServer error")
		return txHash, err
	}
	transOpt.GasLimit = fixGasLimit
	log.GetLogger(log.Name.Root).Debugln("AddUserAddressToWalletSC gas limit: ", transOpt.GasLimit)
	tx, err := instance.AddWalletInfo(transOpt, walletAddress, email, userIsMerchant, userIsMetaMask)
	if err != nil {
		err = errors.Wrap(err, "AddWalletInfo error")
		return txHash, err
	}
	txHash = tx.Hash().Hex()

	receipt, err := ethClient.QueryEthReceipt(tx)
	if err != nil {
		err = errors.Wrap(err, "Query AddWalletInfo receipt error. Trans hash: "+txHash)
		return "", err
	}

	if receipt.Status == 0 {
		err = errors.New("Query AddWalletInfo receipt status is 0. Trans hash: " + txHash)
		return "", err
	}
	return txHash, nil
}

func TestRegisterFlow(t *testing.T) {
	ethClient := initEthClient(t)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey("5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99", ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	fakeAddr, tx, _, err := contract.DeployOwnedUpgradeabilityProxy(transOpt, ethClient.Client)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("fake address: ", fakeAddr.Hex())
	fmt.Println("tx: ", tx.Hash())

	receipt, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)

	if err != nil {
		t.Fatal("cannot get deploy receipt", err)
	}

	if receipt.Status != 1 {
		t.Fatal("transaction fail")
	}

	fmt.Println("real address: ", receipt.ContractAddress)

	proxy, err := contract.NewOwnedUpgradeabilityProxy(receipt.ContractAddress, ethClient.Client)

	if err != nil {
		t.Fatal(err)
	}

	transOpt, err = ethClient.GetNewTransactorFromPrivateKey("5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99", ethClient.ChainID)

	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err = proxy.UpgradeTo(transOpt, common.HexToAddress("0x016442a58cAD7a7110d11D4156e7B22B51069676"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("upgrade to tx: ", tx.Hash())

	receipt1, err := ethClient.QueryEthReceiptWithSetting(tx, 1, 20)

	if err != nil {
		t.Fatal("cannot get deploy receipt", err)
	}

	if receipt1.Status != 1 {
		t.Fatal("transaction fail")
	}

	userWallet, err := contract.NewUserWallet(receipt.ContractAddress, ethClient.Client)

	if err != nil {
		t.Fatal(err)
	}

	transOpt, err = ethClient.GetNewTransactorFromPrivateKey("5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99", ethClient.ChainID)

	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err = userWallet.SetWalletOwner(transOpt, common.HexToAddress("0x016442a58cAD7a7110d11D4156e7B22B51069676"))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("set wallet owner tx: ", tx.Hash())

}
