package test

import (
	"encoding/json"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func TestWalletAddressMap_AddWallet(t *testing.T) {
	fmt.Println("TestAddWalletAddressMap")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)

	}
	transOpt.GasLimit = 10000000
	testAddr := common.HexToAddress("0x35CD8522CFDA982e90AB6859a7FfE8C0b350B3D1")
	tx, err := walletAddressMap.AddWalletInfo(transOpt, testAddr, "eu21@18m.dev", false, true)
	if err != nil {
		errByte, _ := json.Marshal(err)
		fmt.Println("AddWalletInfo error: ", err, " details: ", string(errByte))
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestWalletAddressMap_GetUser(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	///
	isExists, err := walletAddressMap.IsWalletAddressExist(&bind.CallOpts{}, common.HexToAddress("0x35cd8522cfda982e90ab6859a7ffe8c0b350b3d1"))
	if err != nil {
		t.Fatal("IsWalletAddressExist error: ", err)
	}

	fmt.Println("IsExists :", isExists)

}

func TestWalletAddressMap_RemoveWalletInfo(t *testing.T) {
	fmt.Println("TestWalletAddressMap_RemoveWalletInfo")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}

	tx, err := walletAddressMap.RemoveWalletInfo(transOpt, common.HexToAddress("3d40e7e7173b1725035f804d4b646405a91bd39128b6b43fdf0e10aeaa768b9d"))
	if err != nil {
		errByte, _ := json.Marshal(err)
		t.Fatal("RemoveWalletInfo error: ", err, " details: ", string(errByte))
	}

	queryEthReceipt(t, &ethClient, tx)

}

func TestWalletAddressMap_GetWalletList(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	walletList, err := walletAddressMap.GetWalletInfoList(&bind.CallOpts{})
	if err != nil {
		t.Fatal("GetWalletInfoList error: ", err)
	}
	for _, wallet := range walletList {
		fmt.Println(wallet.Hex())
	}

}

func TestWalletAddressMap_AddWriter(t *testing.T) {

	fmt.Println("testWalletAddressMap_AddWriter")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		fmt.Println("NewWalletAddressMap error: ", err)
	}

	for _, userWalletAddr := range userServerHDWalletAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
		}
		tx, err := walletAddressMap.AddWriter(transOpt, common.HexToAddress(userWalletAddr))
		if err != nil {
			t.Fatal("walletAddressMap.AddWriter error: ", err)
		}
		fmt.Println("tx hash: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := json.Marshal(receipt)

		fmt.Println(string(data))
	}
	for _, e := range erc20Currency {
		erc20Addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<" + e + ">")
		fmt.Println(erc20Addr.Hex())

		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err)
		}
		tx, err := walletAddressMap.AddWriter(transOpt, erc20Addr)
		if err != nil {
			fmt.Println("walletAddressMap.AddWriter error: ", err)
			continue
		}

		queryEthReceipt(t, &ethClient, tx)
		//fmt.Println(tx)
	}

	for _, kycAddr := range kycServerAddr {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
		}
		tx, err := walletAddressMap.AddWriter(transOpt, common.HexToAddress(kycAddr))
		if err != nil {
			t.Fatal("walletAddressMap.AddWriter error: ", err, " KYC addr: ", kycAddr)
		}
		fmt.Println("KYC addr: ", kycAddr, "tx hash: ", tx.Hash().Hex())
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := json.Marshal(receipt)

		fmt.Println("KYC addr: ", kycAddr, " receipt: ", string(data))
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := walletAddressMap.AddWriter(transOpt, common.HexToAddress(testUserWalletSCOwnerAddr))
	if err != nil {
		t.Fatal("testUserWalletSCOwnerAddr walletAddressMap.AddWriter error: ", err)
	}

	receipt, err := ethClient.QueryEthReceipt(tx)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(receipt)

	fmt.Println("testUserWalletSCOwnerAddr addr: ", testUserWalletSCOwnerAddr, " receipt: ", string(data))
}

func TestWalletAddressMap_SetLastUpdateTime(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 100000000
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := walletAddressMap.SetLastUpdateTime(transOpt, common.HexToAddress("0x01a6d1dD2171A45E6A3D3dc52952B40BE413fA93"), big.NewInt(1618999150))
	if err != nil {
		t.Fatal("walletAddressMap.AddWriter error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestWalletAddressMap_AddAddress(t *testing.T) {

	fmt.Println("testWalletAddressMap_AddWriter")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}

	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	transOpt.GasLimit = 100000000
	if err != nil {
		t.Fatal("GetNewTransactorFromPrivateKey error: ", err)
	}
	tx, err := walletAddressMap.AddWriter(transOpt, common.HexToAddress(testOwnerAddr))
	if err != nil {
		t.Fatal("walletAddressMap.AddWriter error: ", err)
	}

	queryEthReceipt(t, &ethClient, tx)
}

func TestWalletAddressMap_GetKycLevel(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<InternalSmartContractConfig>")
	// addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	internalSC, err := contract.NewInternalSmartContractConfig(addr, ethClient.Client)
	// walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	walletAddr, _ := internalSC.GetWalletAddressMap(&bind.CallOpts{})
	walletAddressMap, err := contract.NewWalletAddressMap(walletAddr, ethClient.Client)
	val, err := walletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, common.HexToAddress("0x82da583cc679f4c56f26caba15d9132d102d21fa"),
		"kycLevel")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Val: ", val)

}

func TestWalletAddressMap_SetKycLevel(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	addrList := []string{
		"0x6c41750c879659375dfaea4cc18df76560d303bd",
		"0x55bbe0e5bba8e01dddac6be0c8edec4155096291",
		"0xfc32402667182d11b29fab5c5e323e80483e7800",
		"0xbe365605591a843818d866101f3a600543ef2c0d",
		"0xbb8c0e68c35fc71558e95e98ef442059d72f0065",
		"0xd081a230357a54d4a93130b5785ada893a9a2e3c",
		"0x137478012b8bc67aa5a66cfa2f22a39d53e75a55",
		"0x0500bfc566b86d548389c3690efe1c49fc72c5dd",
		"0xa5bd66b90c9f4175f3baf3dd25155fd31543ef81",
		"0x4c47c6ae2757d09c42e1f7854270545be8516d2c",
		"0x4f0d0ed3755381f703f385a2c2ec5508e0b7e5da",
		"0x6d1bddef92f6e5a64ad8a1a09d251a4bcea86bce",
		"0x11406a7317579bd4527f21f9ac5b5c46745a9706",
		"0x8157d3d61ec41ec890611ffdddf8341591a52a4d",
		"0xf7db24fce09fd841ea4e774ef191463897e06b37",
		"0xcf0ead10a8a64c2dbb34f4a3b09ac1857bfabf1d",
		"0x1c1035cdd40c587bb0c61ae92a38768712477c7e",
		"0xa76a5aabf09c56cd8cb5d5c62803d5b282d1754c",
		"0x99b32dad54f630d9ed36e193bc582bbed273d666",
		"0x44f426bc9ac7a83521ea140aeb70523c0a85945a",
		"0x268fede3ee04c1ff53bebea07d8e062d4bcb52bf",
		"0x4dfb6d6790054f3eb68324bc230e3104137ca8db",
		"0x2ee5cacbee6b75a031f76266fb432d29eef7b291",
		"0xf420db1f8b3c8285969bcc3c39658fba0509349c",
		"0xaf540ae89bbd2bd81b5d56c2c8583f37a2ea731c",
		"0x975b6e38758e0b57a1f0628cd45ae30b42fccbf4",
		"0xda341464d6539adab4b8b0a61310c98dd8812642",
		"0xba1bf5d06f0c77e81f461e0ec88c1f36dd481e13",
		"0xb1950d6672df936261e74ed73db3098ccad2e3a9",
		"0xc9db5f3df334abaca986785a3994c031afd8e70b",
		"0x51e2dd0d7f8224413d9c1c6aa22582abedf23a39",
		"0x2619aec08c298948c6bb8dcf53e6c9b3b9bed034",
		"0x2aca65e851585014fe63458471fa395fdcba8f18",
		"0x353fe071dc5708fadef173ac19d041a563c2d99e",
		"0x6715647ec6cc3e54161391ea9bc67edc764e2fb6",
		"0x401d919875836c47785828a56d47a3d9394bb3d1",
		"0x77f6ae5af08d3f854e49301b432536f390fe4541",
		"0x631dbafd2192355aad7a6ce63594e905956e1c9b",
		"0x4ffb69cc66fc6f38cac81762f548cbedf35d77c2",
		"0xacc8cf9c80f294248875a9cfc800e267f694b1a5",
		"0x8d3ff907778fcdba50024e87d1c693366365c1f6",
		"0x0107c9a9ac5e0bf392de4d6f90371687db5b454d",
		"0x93ad54ceaec106dc91ce3b0b8d40b80998beab4e",
		"0xf22b3a4baf2bd721a6a49cda46e8f38958ed360b",
		"0xd17c9f1e141ceca1a7195ba9092a8931751a1814",
		"0x0c9c6f6507b7cf05a060adb86670f0cae17de4be",
		"0x227f3dffd185aa32ed883e90a1d747b2eccd723d",
		"0xcf23e11627bf0185546cf95c7110ee53dfe37c45",
		"0xd2ed863f008fc58418a8126ee1dd417e106d072b",
		"0x620692e6498ea9ddc8bb6ec40117f6e72ce06620",
		"0x3f849950e44d9b3772b331aefb8956f7e4b54abe",
		"0x6769322b0104f3495b6568a49c1c5974bd9c7613",
		"0x24cd5bcb80654537a51eeb48f0193e104d5c1102",
		"0xb8ce0a069f51642f4517ba49e2627db2fa853fbd",
		"0xe33daaa29b344f2712a545c725b5ba163c580c90",
		"0x8f60c0879531dbfeffda451802f8fc6909bb8a78",
		"0x056db572b2f2061320e0059695ede19049953906",
		"0xb4c7e0f4afda486dcd0784c545dc32be46ccdeed",
		"0xe7294f89dad62e4f2708f4f7673b8903c4f73e3d",
		"0x439edc73e1cbba5f5ec6271b64a00d4da7c01ad4",
		"0x307c2f2478f4a2582a8bf741d6305568b3d85375",
		"0x7728932d458588e66cd07f1c3699da4528589355",
		"0xd5ffbdd47cbb919a978415073c37eb36c5defcfc",
		"0x04f75b75b7e5722691735007bacdc9f43c0c8815",
		"0xcf34e1fdae8a85741e633d0672eefbd90961d34b",
		"0x97c70ce958785af288a8aee890fd7af565a4d7ca",
		"0xbd81b6c82b8cc545f7d36b2bbfdd57820a61dc75",
		"0xdd0f95f91a44a0a0ae7bd8012037c2d0100257b7",
		"0xb2cadb9d52cccce841e35d589f6c832240226c71",
		"0x16a39a1c1c9ca95bd8cc6f0432ddd97e1c9bdf77",
		"0x22981bf1ee93353ad8f570a53ccb01c4da107ddf",
		"0x7e5b76bf9492914b7c25784c9590047c50237552",
		"0x4cd099665d96f94f046e452d54c640d572e53081",
		"0xa1da03e8fa18ac257b1f93bf482ff815cd8e10a9",
		"0xff15a6aed2c0aaf60c84f4eb7880e91cfc25464f",
		"0x53792f3520a544e29da8484f02576f7c19c19866",
		"0x9b9912b0e61331a3225f98ee12f1df32a2dfff8a",
		"0x775602b43fabc723bd323a3f03f6508edf907a2d",
		"0xd80e27f1f377a7119d5a44a3366037f3f4edbce4",
		"0x8599016f92636b57b4193e568f66dfda91b74ae8",
		"0x29d7bfb71d25ffae360b35b6633ed1c9596f3ec3",
		"0xf60461d578badd7292ffc166dcdcedfcb389ca64",
		"0xef74831b378aecaab51a05d400165ff0f9c71b84",
		"0x1a6db76156e0511ea4d3be38c9a8e03a9018053f",
		"0xf8fd9b5cbb28103e0e0f19133d341a560fdedc9e",
		"0x6cf09ca68fa47f78d756fb509a4410642a6efed9",
		"0xd534e523c97a863dd6acdcd77b1cb47c56439a26",
		"0x4760d9c10b9ca3523494abe32de758b82380df33",
		"0x42fc4f29c5dea28b1e52186e0d04f6a48a96f8da",
		"0x99ed82dfbdd2a4cadcd77e22ba0223d8a869e9d5",
		"0xa18ec23dc6d4a7c7af53356bd31ee838bf349e44",
		"0xecdea3853c3b63c2d405b2486bddc09db1a25f0b",
		"0x69e2534c32925994c1fc89f549f228fefce46598",
		"0xa6c8a3a0e633d086647b07b78ce4903a5f97cf80",
		"0xa07f7beb4ac3ae2ae6d61595300ce3c722ce9a78",
		"0xf934b7065a12b501d83d5885a7dc28e55428d713",
		"0x2fcef38fee078374a7c37ddbe1370742c259f130",
	}
	for _, addr := range addrList {
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey error: ", err, " address: ", addr)
		}
		transOpt.GasLimit = 300000
		tx, err := walletAddressMap.SetWalletInfo(transOpt, common.HexToAddress(addr), "kycLevel", "0")
		if err != nil {
			fmt.Println("SetWalletInfo error: ", err, " address: ", addr)
			continue
		}
		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("Query receipt error: ", err, " address: ", addr)
			continue
		}

		receiptByte, _ := json.Marshal(receipt)
		fmt.Println("Address: ", addr, " ", string(receiptByte))

	}
}
func TestWalletAddressMap_GetOwner(t *testing.T) {
	fmt.Println("testWalletAddressMap_GetOwner")
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	fmt.Println("wallet address map proxy address: ", addr.Hex())

	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	addrs, err := walletAddressMap.GetOwners(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(addrs)
}

func TestWalletAddressMap_GetWriterList(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	fmt.Println("wallet address map proxy address: ", addr.Hex())

	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}

	writerList, err := walletAddressMap.GetWriterList(&bind.CallOpts{})
	if err != nil {
		t.Fatal("WriterList error: ", err)
	}
	fmt.Println("Writer List:")
	for _, writerAddr := range writerList {
		fmt.Println(writerAddr)
	}

}

func TestWalletAddressMap_RemoveWriter(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	fmt.Println("wallet address map proxy address: ", addr.Hex())

	walletAddressMap, err := contract.NewWalletAddressMap(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := walletAddressMap.RemoveWriter(transOpt, common.HexToAddress("0x0000000000000000000000000000000000000000"))
	if err != nil {
		t.Fatal("WriterList error: ", err)
	}
	queryEthReceipt(t, &ethClient, tx)
}

func TestWalletAddressMap_GetImplementation(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	fmt.Println("wallet address map proxy address: ", addr.Hex())

	proxy, err := contract.NewOwnedUpgradeabilityProxy(addr, ethClient.Client)
	if err != nil {
		t.Fatal("NewWalletAddressMap error: ", err)
	}

	implAddr, err := proxy.Implementation(&bind.CallOpts{})
	if err != nil {
		t.Fatal("Implementation error: ", err)
	}
	fmt.Println("Address: ", implAddr.Hex())

}

func TestWalletAddressMap_DataMigration(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")
	oldAddr := getAddressBySmartContractName("Old_OwnedUpgradeabilityProxy<WalletAddressMap>")
	fmt.Println("Old wallet address map: ", oldAddr.Hex())
	newWalletAddressMap, _ := contract.NewWalletAddressMap(addr, ethClient.Client)
	oldWalletAddressMap, _ := contract.NewWalletAddressMap(oldAddr, ethClient.Client)

	oldAddrList, err := oldWalletAddressMap.GetWalletInfoList(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}

	var accountMapping map[string]string = map[string]string{
		// "0xcff3e8123d876fca6cf17ada15731714879a202d": "0x2004e93d1f948748e5d39d30af71d0cd8ae4db78",
		"0xdae85eca3e47ff5285d720826b86cd86059e6bfd": "0x54d5387253f10df0419bd6eb74cd97c69c128f87",
		"0x512d1365abecae947c98f33e76b10be9a1baeb96": "0x8d77691037c9fad03e3df1d638bd9c7bbf737e2d",
		"0x12bf04f42de6f1c26e98c97c3cefbcc00adc8c30": "0xf377fd45471ef47fbb6aecc17a586d4dd3f48cd4",
		"0xef04bc652311b6903d91d5245cab58e9ce04130c": "0x5ce4a04a264850d6786911288b3bee1f448fc89e",
		"0xf6fa8b1715230d4b8099ba5c0c323f0fd88b8387": "0x2722d9f757254e55e2a8aaee89e8254cc1bf533e",
		"0xf144ab9b04af3f0cd7aa859bf52a93555840c911": "0xc16ffe0da5dae266168df31f9891702277195c0f",
		"0xee8fc6c1b85238a49b9677f0d755310d31c29697": "0x9b004790614222e09aaa946cda51c1d0924ba2c4",
		"0x32e5e02546b0d9e19ad76d84892d4d3a5cb26c4c": "0xa9f15eaf6aa6292070509ad40e4852b381e6e48c",
		"0x0a604a2d0e9240d5f1e10bffbc49d807577a557c": "0x303e189e85e234f3da02bc33716b6c7b288f9f5b",
		"0x03f9833eaaf99a22b2d4336371a79e370c673b9c": "0x0752fbab174591950ef4f89ec6d1b362ce57dca5",
		"0x8db9be1301ddca68b3770e69478aefbb73836f34": "0x9dea7d2be2f6d0ea9a9c4146fb9080f61a9c1189",
		"0xf727449e81c8e2a62fe6db2e61fe5f4ef8501cd4": "0xeb38c2255244dc4fbba425cef34f3cf7a03e6eb6",
		"0xab52dfa82c7cb904695a55ce279e9a76541437b9": "0x7da27989f5b810faebbbb1e8069ebb56e96800a9",
		"0xf60265ff331cd85b00f4b29f685d2ff951eb2078": "0xfc0cd27314ed4e38c08b474afc5027070d38f33a",
		"0x9b158f5937b97ae713e1b57bee700e397b79a91d": "0x36a522ea9e93d68457e64b7409c4b8afff15b375",
		"0x7a0c44991fa3e7a62f83544d1c1d3b9b78814746": "0x1bc9ae462d3d5ece278e6bdd24f7d88f8468980a",
		"0xcfd4040302cfb7ef35699a60f894086090ada2ae": "0x4ca0786c42db02144de4af04f49354f8b34c908b",
		"0xd21a1175922299433deec446500d49f0d9b25859": "0x3340bffccb1077dc2de052ddc3a9c3d8f8c9eafd",
		"0x7adac01c9d66496d954d488b807639ab51aa6667": "0xad5b7d6a62a9476697082346ce6c9cdb38d16294",
		"0x0979eda5586d60e49c39783fbf87f123067f80a2": "0x9db7481baea4b847fa8b238cfc288067187b3130",
		"0xe474421ccaac1bf9cbe182c1d25d64508fc7574f": "0xaefc260cc91713ad3643af3cc1fa9f4bcb5e133b",
		"0x4ad55f52c23499053ebb61ec80d1bcb6fc9e094f": "0x42579d2a440de27c326aec057d79c71861c73c0f",
		"0xe95260e3dbac9164c0f643b3cde62d81c8e5e329": "0x7b333488c58cc79ac6b0edc5f54e477ec220f5fe",
		"0xe6712b4e2076319a48974466477a4b29a5a142da": "0x482ab6f4f73493cd3b22bddcf9b136cf1367ad60",
		"0xb4a70777f1bd82de875aa49bfdd68e6913be51df": "0x722981f4bcf5b464dcbecbe7b678e960afe13190",
		"0xea99d7551c60f4cb0e7f0999de8829009fbe1a46": "0x2a65dd702f4895c53ddcfcf11ea1d13d1fb17836",
		"0x32f7ce28cd38a843adfd7660f035f7688253288d": "0x6e89d4552b64037227b4f57238ab601c63e682c0",
		"0x6d438dded7bf03e01b35a564fcaa30062f30e233": "0x85e431426d4fddf407dd575e73dfeef735c96f8d",
		"0x8307c5bcfa15220efc04ac2e893b43cf8faf9556": "0x1815cdaa36d03fbebf6f98c9a95bdb340f76f475",
		"0x8828ae592433d1a85a88841f4531f8834538b727": "0xe66f6d6b7961678765634b379360a17c691f19ec",
		"0x3c28719baa2b05ebb2d0ae9d88adf8f8e9aba983": "0xa5fdb1d83483f4dc78d35cee4f3a8c9852f39285",
		"0xc4c7a6bc0981032b858be0766f95d65d9b1b6ad0": "0xfa61d26dee3598e4bb267d502dc947e100983425",
		"0x9e5d09343f8e48560c8a6b5d717bd047713977f1": "0x49210c570fd48b9aaa0dc8cf1ed70480a7c384b2",
		"0xd826868372501bc2a4bc21cb4ac415286b197b60": "0x837794da1c3700bedc62f50f4aa372bd98687e19",
		"0xd5f269447b3238e61aa4e80a886039383a2586d7": "0x1d60b58a917ebe3fccf14f0fa3e0c073a8feaed5",
		"0x443a214af82e2363d11b3cd0e5f8a429d602e8d1": "0xe85b30e8db717d61c047095e161b6ab642014982",
		"0x9f86ed154575dcbe37a3d917dc69fb14358864af": "0x9c4b9cbc2af83e6d719a386731085c6bcc4f96fe",
		"0x8f23f891f4b02afec2f71a36843e4978f468af88": "0xf256761705f9f15b498f1954b10b6efe63d61c49",
		"0xd806dc9430115591c2229fba119c1f6dee191b78": "0xe85f54653269d9fcecb15b5f089c97c3c2985020",
		"0x8835b61bca7b8f4d140b1d8a93aaf9b3f3097fdf": "0x39d13aa326649cf3d38ebb38faf89551119ebb26",
		"0x4ec0e9d44b03224e4c8f00d351d65880c89f6f66": "0x4c259d3c72d1f1acd763079624e74d42cc3024f5",
		"0x1f07fb9227aca24f30a338b947ffbc806a2c82d8": "0x6dd67ee385b2f1a25c28fa27482e5138b639bd47",
		"0xacc795036071d57af26816ee3e142a10799e4b41": "0xc5d7fcb4b6503d78156271081c4213dc5b4632d5",
		"0xa6a21d46fcead0d06248b1ff2ed6c6b02f5e32dd": "0x9f932a4530bce0d37bdbc86524e5272abddb21a0",
		"0x534cd112e99d45c0ea888fb858acd455ceee626d": "0xdda51cf03b0624444a8683780d91927923e98191",
		"0x1f85726917b4ddbc0d1fedffa7b34e5979088ee2": "0xde55fef718aa7c316b3948f277918e3897d9b0a3",
		"0x37de2417c9873857e714f245958303a1cc32f9a0": "0x1dad7e04be1649e5cda02ea7e582f7be2f9d2103",
		"0xf27b54b0781d6c3a2174b82a474a88079411ac10": "0x1daf70fb148e0d99273706cafeea00be179f3b74",
		"0xce6ea6f47c45d0c5e0e9828a18fb8f1c78ec00eb": "0x2c1f05acfe28d2522e9c5e1669b84471d4e74661",
		"0x80d3a6cf3c6e527e8defff22f678029d6a069494": "0x4f4245fc801686af475a7501bd9c745020593f86",
		"0xc60b78f8b0f9022761da6c48d53ffbad394ad150": "0x1b1f42b0d6f949bd1b427c86a053d6e3a1d19d08",
		"0x36c144209ff52c2fc92e289348bce10a6935bf15": "0x55d9f1b16f4827b6034dbd3206e08bbed38eeb07",
		"0x7dabc6d8bef126092c64311e670a66fa98b2b679": "0x56c45aa963cb87a8d474f176694c755e5b51d092",
		"0x1bf01e543fec2be2bf3a22bfd9e58bb37384e868": "0x7caa8b87b6b27dcede99a185e0b7ecaf95425069",
		"0x4261a22e731645844d609c6f672c472e7e85b907": "0x94479487c057bb42597bd7e7d8265b2a08dd16ab",
		"0xd4983f30b49a8b1408058bc8a0273c89a9306960": "0x97df2101247a6794c21d093833f6fb4848798422",
		"0x47b5a68ecbbb84d016e19ae73ca764717b17de55": "0x85226e8e1dde72fb54d93a2761fe46587cc729ea",
		"0x05862a244d1d192194db3aedd991fbd6362e011e": "0x8e13239b205b0363681c01beb5e53a6cf729bf4e",
		"0x6beacafc60c8f34bc452b5717c8ccc1dba132ab0": "0x70e6ca28e28b9c154ca72a4a7f51ff403624be47",
		"0x7629ea34b8de6bf2249b7435535099a286bf6390": "0x040078fb3772e904cdae529413c635545f81b20a",
		"0x74604d9261fd531ac91bb4c264dfa47344ecac5e": "0xe2b7c797ecba3605dc7ffb1c587f573e6c64357c",
		"0x4ca0c1c53e40a491a709ab35de7d5d91f604602e": "0x7c488240c9100f4cda451354f49704781c07471a",
		"0x3431bfeb7ae31b05335a1d54c6e7cb7670d3088c": "0xae9a23c0ac3e8dc4c579197d83a4ff8605182989",
		"0x2fbd110d6624d2bc295f6fd254d185f367bc1909": "0xeb52fede23825152a9a1b880f5188917237b1898",
		"0x99b32dad54f630d9ed36e193bc582bbed273d666": "0x99b32dad54f630d9ed36e193bc582bbed273d666",
		"0x809eecf74ee563d819ccb9761410e3541f18e747": "0x809eecf74ee563d819ccb9761410e3541f18e747",
		"0xeaf888a633e9d6641b35364764f4934786c14851": "0xeaf888a633e9d6641b35364764f4934786c14851",
		"0xd02076408d7edef9b31e41bc109c70118241353a": "0xd02076408d7edef9b31e41bc109c70118241353a",
		"0xd7cdfadc13a6d0fe8e9b6474ea722f2592220dda": "0xd7cdfadc13a6d0fe8e9b6474ea722f2592220dda",
		"0xe7dffd84ebcba3e43774de810ea341a4620485c9": "0xe7dffd84ebcba3e43774de810ea341a4620485c9",
		"0x12986db5d522476d1dffc68305a7eff224df55d9": "0x12986db5d522476d1dffc68305a7eff224df55d9",
		"0x7e5b76bf9492914b7c25784c9590047c50237552": "0x7e5b76bf9492914b7c25784c9590047c50237552",
		"0x6826e319562823cd527646265c5f81d4f51a0899": "0x6826e319562823cd527646265c5f81d4f51a0899",
		"0x5c28178b275657c0b71f1513a239705605988bd4": "0x5c28178b275657c0b71f1513a239705605988bd4",
		"0x44f426bc9ac7a83521ea140aeb70523c0a85945a": "0x44f426bc9ac7a83521ea140aeb70523c0a85945a",
		"0x1866be0c998504d06bb71919187c2d589410548f": "0x1866be0c998504d06bb71919187c2d589410548f",
		"0x853d359a46cdd31ba44bfbc18dfea1e83c5e0630": "0x853d359a46cdd31ba44bfbc18dfea1e83c5e0630",
		"0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93": "0x01a6d1dd2171a45e6a3d3dc52952b40be413fa93",
		"0x2d1366c71e86f20de3eecc3f00f270d78a8cefe5": "0x2d1366c71e86f20de3eecc3f00f270d78a8cefe5",
		"0x4dfb6d6790054f3eb68324bc230e3104137ca8db": "0x4dfb6d6790054f3eb68324bc230e3104137ca8db",
		"0xd17c9f1e141ceca1a7195ba9092a8931751a1814": "0xd17c9f1e141ceca1a7195ba9092a8931751a1814",
		"0x52cb9ae015bc3e699bf2af05536f8b1e384e0234": "0x52cb9ae015bc3e699bf2af05536f8b1e384e0234",
		"0x42f0b59d0082785bffafdce0daa3d5c89f5e93e8": "0x42f0b59d0082785bffafdce0daa3d5c89f5e93e8",
		"0xa6bca7be25618b12617757032c290a56833056a5": "0xa6bca7be25618b12617757032c290a56833056a5",
		"0x245382ab1172bd8b641b4688d6eaea54b4ff24dc": "0x245382ab1172bd8b641b4688d6eaea54b4ff24dc",
		"0xf1689e88033d26b8ad3a414aec02cb2aee207331": "0xf1689e88033d26b8ad3a414aec02cb2aee207331",
		"0xc9db5f3df334abaca986785a3994c031afd8e70b": "0xc9db5f3df334abaca986785a3994c031afd8e70b",
		"0x0a506f423ef95f31fafd3b8f81907afc4507a1b6": "0x0a506f423ef95f31fafd3b8f81907afc4507a1b6",
		"0xfe945f0b546882df9157398de0e9b3ed07942cad": "0xfe945f0b546882df9157398de0e9b3ed07942cad",
		"0x81c0317c04b38f8bf60be2f4d8b0d2c3ac840c0c": "0x81c0317c04b38f8bf60be2f4d8b0d2c3ac840c0c",
		"0x288e7838994f2d07b16210dbed6184eb01ff6ddb": "0x288e7838994f2d07b16210dbed6184eb01ff6ddb",
		"0xa76a5aabf09c56cd8cb5d5c62803d5b282d1754c": "0xa76a5aabf09c56cd8cb5d5c62803d5b282d1754c",
		"0xc737f1481cc50e813bb91c93326c988f7d4d37c8": "0xc737f1481cc50e813bb91c93326c988f7d4d37c8",
		"0x9ebce94684c7a2c5790824a21fe14ac02db6c443": "0x9ebce94684c7a2c5790824a21fe14ac02db6c443",
		"0xd601ba6cc195bbdfe91493c0fe0c9d379ed2ab1b": "0xd601ba6cc195bbdfe91493c0fe0c9d379ed2ab1b",
		"0x39dccf5ff2d41018302ee549705cfbf06a4be920": "0x39dccf5ff2d41018302ee549705cfbf06a4be920",
		"0xaa72c944b3d6c9a54b90e6a68812ec676e44515c": "0xaa72c944b3d6c9a54b90e6a68812ec676e44515c",
		"0x02215c39def3991c683438703af8e7dbd67c7fb9": "0x02215c39def3991c683438703af8e7dbd67c7fb9",
		"0xfba595d7fbb1f967727f23381c130513982840f4": "0xfba595d7fbb1f967727f23381c130513982840f4",
		"0x558935f5caf98a3e4d1db2ec0da098fe746c889c": "0x558935f5caf98a3e4d1db2ec0da098fe746c889c",
		"0xebf5ffe82a703200c212fd499e4bee659491343e": "0xebf5ffe82a703200c212fd499e4bee659491343e",
		"0xf18a12236052f046dab9800354da7172185a438b": "0xf18a12236052f046dab9800354da7172185a438b",
		"0xc00f6d107f242acb151a93f1f249c6e3e56a69de": "0xc00f6d107f242acb151a93f1f249c6e3e56a69de",
		"0xba98075d9e2864787db2681740646a4fedad0ded": "0xba98075d9e2864787db2681740646a4fedad0ded",
		"0x163638be9ddad5cdce838cc721a806af2dd843d9": "0x163638be9ddad5cdce838cc721a806af2dd843d9",
		"0x510020a41f56d7ad2b8e682d9ae31ab115e9deb6": "0x510020a41f56d7ad2b8e682d9ae31ab115e9deb6",
		"0x9cf79b7aeadd471cd476e75e9f531afea05bc748": "0x9cf79b7aeadd471cd476e75e9f531afea05bc748",
		"0x4f2d70ecddf9c475748a403c45eea36062e2fa14": "0x4f2d70ecddf9c475748a403c45eea36062e2fa14",
		"0x3ec664828b6ccc22658f4470d9b143f40b1857de": "0x3ec664828b6ccc22658f4470d9b143f40b1857de",
		"0xfbff014e4882223b61253982cef74951a4eaca4e": "0xfbff014e4882223b61253982cef74951a4eaca4e",
		"0xd004c35029fb3f4f27ac918e80e04e41c5b96374": "0xd004c35029fb3f4f27ac918e80e04e41c5b96374",
		"0x962daa25a9cfd0ee828a189d35ed8cac9026c043": "0x962daa25a9cfd0ee828a189d35ed8cac9026c043",
		"0x45b010331a56ef10d4a7a6b9ad186e5c2db0c867": "0x45b010331a56ef10d4a7a6b9ad186e5c2db0c867",
		"0x02358868da51036c47059a23bae51c60eeec7644": "0x02358868da51036c47059a23bae51c60eeec7644",
		"0xb91f6b42977baadcb8f5fbb04eccb5cd3f014a73": "0xb91f6b42977baadcb8f5fbb04eccb5cd3f014a73",
		"0x3b5fdca72cbac221f07a6ce204617f86a2f06df0": "0x3b5fdca72cbac221f07a6ce204617f86a2f06df0",
		"0xe9e1690490cb154d5dbbf4c7cab5eec7e9318296": "0xe9e1690490cb154d5dbbf4c7cab5eec7e9318296",
		"0x4291b8c93aea0a2a0c0920dc9f9ba1581d8b51f8": "0x4291b8c93aea0a2a0c0920dc9f9ba1581d8b51f8",
		"0xe679dd6ae269f1bda511ab244c1386ea4b019625": "0xe679dd6ae269f1bda511ab244c1386ea4b019625",
		"0xa57a4360769ff865160e907f8c0bc5623a2b1b98": "0xa57a4360769ff865160e907f8c0bc5623a2b1b98",
		"0x68f0318fe06458821db8754e21027e4e2bf52be8": "0x68f0318fe06458821db8754e21027e4e2bf52be8",
		"0xe5d69e1c8dc0bb5f35c44650405774c6e0147384": "0xe5d69e1c8dc0bb5f35c44650405774c6e0147384",
		"0x352b4866d9e31966960118fafc3d9346044f8ca5": "0x352b4866d9e31966960118fafc3d9346044f8ca5",
		"0x5369d9800688660adf52f2b0b641a20cdcb1052d": "0x5369d9800688660adf52f2b0b641a20cdcb1052d",
		"0x2167c15f3667be33d616e2f213acc2072df2d9cb": "0x2167c15f3667be33d616e2f213acc2072df2d9cb",
		"0xa445d785b331c7d64c3b93b2d6085037b6693aa1": "0xa445d785b331c7d64c3b93b2d6085037b6693aa1",
		"0x48b91efeae1b5b131cc19f6abcadf43d92c5bb77": "0x48b91efeae1b5b131cc19f6abcadf43d92c5bb77",
		"0x84643323975c83c0052cf04abe7c0b90d562934b": "0x84643323975c83c0052cf04abe7c0b90d562934b",
		"0x473ea8cdb4b226a4045a76c80cd0e95fadf89379": "0x473ea8cdb4b226a4045a76c80cd0e95fadf89379",
		"0x74e60aa73079af67a98fa94d0dc370e2a65c4094": "0x74e60aa73079af67a98fa94d0dc370e2a65c4094",
		"0xfb689b4dd58fbe341d8ce7ef3d6861a937e273d4": "0xfb689b4dd58fbe341d8ce7ef3d6861a937e273d4",
		"0x772af0faa890e792b5dda2ce8bfca47c7478e350": "0x772af0faa890e792b5dda2ce8bfca47c7478e350",
		"0x7b4722a310b677796bf59c199d2110994f5825f0": "0x7b4722a310b677796bf59c199d2110994f5825f0",
		"0xc53afdb97997da7bf5e462ac26c3d78ff27d9967": "0xc53afdb97997da7bf5e462ac26c3d78ff27d9967",
		"0x8e4a817693c31849312549e58ea4f83f49d4df5a": "0x8e4a817693c31849312549e58ea4f83f49d4df5a",
		"0x5f782de89cd3afd2997476f74bfffc06fac25895": "0x5f782de89cd3afd2997476f74bfffc06fac25895",
		"0x468e9f1c3b8dcf6ca033e1c44404e40f9df94e9c": "0x468e9f1c3b8dcf6ca033e1c44404e40f9df94e9c",
		"0x889b11f5cc382003200dabc903d5d61acee194f9": "0x889b11f5cc382003200dabc903d5d61acee194f9",
		"0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a": "0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a",
		"0x59f60eaa816cd1389284ae96d121bec4262f629e": "0x59f60eaa816cd1389284ae96d121bec4262f629e",
		"0xef92a508db8975e756b119167b2977163a558c50": "0xef92a508db8975e756b119167b2977163a558c50",
		"0x29fbd775f4180fdd4f90b815dfa2596880647c86": "0x29fbd775f4180fdd4f90b815dfa2596880647c86",
		"0xb176712aa696beba96c2f2d3fc32c0a7adcb05e2": "0xb176712aa696beba96c2f2d3fc32c0a7adcb05e2",
		"0x1563915e194d8cfba1943570603f7606a3115508": "0x1563915e194d8cfba1943570603f7606a3115508",
		"0xb1ff5a4bd9b79de9a3719047d181fe9ae8a913d6": "0xb1ff5a4bd9b79de9a3719047d181fe9ae8a913d6",
		"0xf7016d4220b02b16b47c926abfa1f5f564af90b2": "0xf7016d4220b02b16b47c926abfa1f5f564af90b2",
		"0xb349e541cbc4f78aa20f3c4d627ec77bb5791e23": "0xb349e541cbc4f78aa20f3c4d627ec77bb5791e23",
		"0x930b91708c121a50b25e0143971141ed83299412": "0x930b91708c121a50b25e0143971141ed83299412",
		"0xe7e07c8399c5f60c3c0bc2f73c11a8379c8484f6": "0xe7e07c8399c5f60c3c0bc2f73c11a8379c8484f6",
		"0xae9c19e0ae50493912c2db9a2958be6340e2412c": "0xae9c19e0ae50493912c2db9a2958be6340e2412c",
		"0x6474c312f9613482eaa2601cdb6f8e300a2e6514": "0x6474c312f9613482eaa2601cdb6f8e300a2e6514",
		"0x000a5aa2b6f3cc277a3b6ba31eeaf2b798e059c1": "0x000a5aa2b6f3cc277a3b6ba31eeaf2b798e059c1",
		"0xcfbe777d1aef139c518327ab8bc668323a30e7fe": "0xcfbe777d1aef139c518327ab8bc668323a30e7fe",
		"0xdff6e7dddcd20751a129a459cc18dfdded825207": "0xdff6e7dddcd20751a129a459cc18dfdded825207",
		"0x120c055a664dc962c452e56f59ceea416c8a3505": "0x120c055a664dc962c452e56f59ceea416c8a3505",
		"0xd9aee902679e5da25a3c6908b5d2edcca086c9aa": "0xd9aee902679e5da25a3c6908b5d2edcca086c9aa",
		"0x50d707d6a43d9185a1f34411b7b190c889a9534b": "0x50d707d6a43d9185a1f34411b7b190c889a9534b",
		"0x8fcfaf17d94c0988b7d704619ee5466da3225f97": "0x8fcfaf17d94c0988b7d704619ee5466da3225f97",
		"0x7276bf9ec06b62a4ef4952298599b04abd6171c6": "0x7276bf9ec06b62a4ef4952298599b04abd6171c6",
		"0x7f82c1ab7b50f3fdd1087eff73d59338af461704": "0x7f82c1ab7b50f3fdd1087eff73d59338af461704",
		"0x75d4cfe0490a89d4092794f2d5fd002dd9092fbd": "0x75d4cfe0490a89d4092794f2d5fd002dd9092fbd",
		"0xd900f39164bf9112dc06c5ba058c3d3d68dada6f": "0xd900f39164bf9112dc06c5ba058c3d3d68dada6f",
		"0xb40679e1b3416545837dabc8fae0e5aea11796d3": "0xb40679e1b3416545837dabc8fae0e5aea11796d3",
		"0x679ae40350451d467c72879ec16633a1d977133b": "0x679ae40350451d467c72879ec16633a1d977133b",
		"0xf4297f00a3c8b54abf3087832bbe8bdb53efa7f2": "0xf4297f00a3c8b54abf3087832bbe8bdb53efa7f2",
		"0xb315ab8db5417519eb15aafe6888a4384412cdae": "0xb315ab8db5417519eb15aafe6888a4384412cdae",
		"0xa60ad6b48f4657c7152988171325e3c5457b72c2": "0xa60ad6b48f4657c7152988171325e3c5457b72c2",
		"0x563c0b9a16d0df25b3ca4a6bffd784a047b81bf3": "0x563c0b9a16d0df25b3ca4a6bffd784a047b81bf3",
		"0xbc154092a4e2e4952da1c77dca657dabf25162f1": "0xbc154092a4e2e4952da1c77dca657dabf25162f1",
		"0xd7e2e45f36f0876271b2d7d8369bc6940eadfa1b": "0xd7e2e45f36f0876271b2d7d8369bc6940eadfa1b",
		"0xcbcc12c5b448890ffd8d8e9342985db5faa29855": "0xcbcc12c5b448890ffd8d8e9342985db5faa29855",
		"0x1a0b1611d67102f66cd6369f285bd81cf0e49d96": "0x1a0b1611d67102f66cd6369f285bd81cf0e49d96",
		"0x1ea8dc232a269f5c2bd3ad4a95cf10426c612772": "0x1ea8dc232a269f5c2bd3ad4a95cf10426c612772",
		"0x942185dffcd20b33af5bdee0b970f9ebd222c3c2": "0x942185dffcd20b33af5bdee0b970f9ebd222c3c2",
		"0x5ba2de7a10d1285c4252192504696c57abe13ff3": "0x5ba2de7a10d1285c4252192504696c57abe13ff3",
		"0x2612650150b970449626bb966a94fd11ffd3d16c": "0x2612650150b970449626bb966a94fd11ffd3d16c",
		"0x22b06aea6f55712b5343144c85729110aac1b3cc": "0x22b06aea6f55712b5343144c85729110aac1b3cc",
		"0x864f91cc7913c728205cefadcd3b4be124e05218": "0x864f91cc7913c728205cefadcd3b4be124e05218",
		"0x332ca68868be43499b5c7a41fcc4fe17046d36cf": "0x332ca68868be43499b5c7a41fcc4fe17046d36cf",
		"0x9ea8224e1e961b15412082c22b9f93b375cfb44c": "0x9ea8224e1e961b15412082c22b9f93b375cfb44c",
		"0x8302b48d5203268a7f229f621803abb1e10af21b": "0x8302b48d5203268a7f229f621803abb1e10af21b",
		"0x1325c8659d7d67ddfe89520f48ed46c0cb228760": "0x1325c8659d7d67ddfe89520f48ed46c0cb228760",
		"0xe26c38c6499f5974055c941035302fc7482ec4be": "0xe26c38c6499f5974055c941035302fc7482ec4be",
		"0xa86cdb522482ed8c595a1a3fb64dc1b9b1f03152": "0xa86cdb522482ed8c595a1a3fb64dc1b9b1f03152",
		"0x9f933630e4dc61b0aad9cc2cc0531e9fa59639d8": "0x9f933630e4dc61b0aad9cc2cc0531e9fa59639d8",
		"0xd60f1a993266de726bfb80a74b7a20510c832e61": "0xd60f1a993266de726bfb80a74b7a20510c832e61",
		"0xe906be29f0735fd35235aefd9c85015ee2b5efc1": "0xe906be29f0735fd35235aefd9c85015ee2b5efc1",
		"0x97b8cb5e2bfe040a4637e22aaeb4cabc69aed9a2": "0x97b8cb5e2bfe040a4637e22aaeb4cabc69aed9a2",
		"0x394dacde3e5e74baae2345de6748bea235b54247": "0x394dacde3e5e74baae2345de6748bea235b54247",
		"0x386231623c415f4526e4088d553947d155050e1b": "0x386231623c415f4526e4088d553947d155050e1b",
		"0xefe9b633d5ce4a443eb8114959a5c3710b4553d7": "0xefe9b633d5ce4a443eb8114959a5c3710b4553d7",
		"0xf30d2ffff87c23b2edebe03ed67d986374ed78f4": "0xf30d2ffff87c23b2edebe03ed67d986374ed78f4",
		"0x1efbd83c67eb73c3bdf0be6a00054af17d06258b": "0x1efbd83c67eb73c3bdf0be6a00054af17d06258b",
		"0x574248771b81c762aaa0e10ead1506e943cd5324": "0x574248771b81c762aaa0e10ead1506e943cd5324",
		"0xebba19c96984d3a9c888d6af9191919ec18d1dd0": "0xebba19c96984d3a9c888d6af9191919ec18d1dd0",
		"0x6b1730b24208dadd0efb985823c2217fc1bb7368": "0x6b1730b24208dadd0efb985823c2217fc1bb7368",
		"0xeb1ca40c9c747e67d6ba25fdb34f948fbd73caa6": "0xeb1ca40c9c747e67d6ba25fdb34f948fbd73caa6",
		"0x8d3ff907778fcdba50024e87d1c693366365c1f6": "0x8d3ff907778fcdba50024e87d1c693366365c1f6",
		"0xa5a576403cdc40ff009e736d0fb2f8c0cbfebdf7": "0xa5a576403cdc40ff009e736d0fb2f8c0cbfebdf7",
		"0xda1a5f5b7d25691d310c37d1afe7faac7059c586": "0xda1a5f5b7d25691d310c37d1afe7faac7059c586",
		"0x289b1dad244dc4ef75d5bfefe48f17675e6357f4": "0x289b1dad244dc4ef75d5bfefe48f17675e6357f4",
		"0x5cf4fdb84492f2e8dabb483d7312a3ac2e162763": "0x5cf4fdb84492f2e8dabb483d7312a3ac2e162763",
		"0x3ef067f921d58993081cf25093ffb0d8e4947219": "0x3ef067f921d58993081cf25093ffb0d8e4947219",
		"0xda4c0958df705498fac04ab7ee2060b2fc64419e": "0xda4c0958df705498fac04ab7ee2060b2fc64419e",
		"0x19e9490f2d0efbcce0c41676df82c2dfd880bd1c": "0x19e9490f2d0efbcce0c41676df82c2dfd880bd1c",
		"0x946be3e01341521088c33ac34a7efaa3c04ef4d5": "0x946be3e01341521088c33ac34a7efaa3c04ef4d5",
		"0x233ddacb786df23c865cf49afbbd567206cb0fed": "0x233ddacb786df23c865cf49afbbd567206cb0fed",
		"0x69bfef5dfb7beb0c8a4795b7bce20ed1e6f5fb8e": "0x69bfef5dfb7beb0c8a4795b7bce20ed1e6f5fb8e",
		"0x25199e2d2af2ecb87fa72b3c83546f46d023205e": "0x25199e2d2af2ecb87fa72b3c83546f46d023205e",
		"0x6b5600446eb6aecb5b11b12b1728d1d125c5021f": "0x6b5600446eb6aecb5b11b12b1728d1d125c5021f",
		"0xc1dbceef7ea506402b4fe833c5f13c402ad1a5f5": "0xc1dbceef7ea506402b4fe833c5f13c402ad1a5f5",
		"0x6fc2caf7ca8e045ae6c50c2d1edc798b29b02fef": "0x6fc2caf7ca8e045ae6c50c2d1edc798b29b02fef",
		"0x5ba418d7884fe46df8e6a176c1dbfcb9b2dbe820": "0x5ba418d7884fe46df8e6a176c1dbfcb9b2dbe820",
		"0x6dbeb7ceb6ec6885be62f48996bcd8342107d9c4": "0x6dbeb7ceb6ec6885be62f48996bcd8342107d9c4",
		"0x229751ca0c76dbdc33a2d869965c085a84595ee3": "0x229751ca0c76dbdc33a2d869965c085a84595ee3",
		"0xe5f12c99f16315ec04759cdf14718d3549887861": "0xe5f12c99f16315ec04759cdf14718d3549887861",
		"0x3d5f40468cc2da6e2bd55afd35bd090eb1592c1c": "0x3d5f40468cc2da6e2bd55afd35bd090eb1592c1c",
		"0xe03afb82e0f1f9a1382b8780c25d7117d2508644": "0xe03afb82e0f1f9a1382b8780c25d7117d2508644",
		"0xbe000f445b1a3ac3abb5d76b04407ab09fdfb363": "0xbe000f445b1a3ac3abb5d76b04407ab09fdfb363",
		"0x323fd3b5c6e55088ab7ccbf26135a143532b78ea": "0x323fd3b5c6e55088ab7ccbf26135a143532b78ea",
		"0x2f62efee3f915dc9c5e07e9c3d0dc93c2f4ce5f0": "0x2f62efee3f915dc9c5e07e9c3d0dc93c2f4ce5f0",
		"0xda59b0bfdf4112fdc7168935c4c7ad1e76e76ceb": "0xda59b0bfdf4112fdc7168935c4c7ad1e76e76ceb",
		"0xc8c5b36c4a2d4d3b72379d135ee64be0f7d46cad": "0xc8c5b36c4a2d4d3b72379d135ee64be0f7d46cad",
		"0x0fac7d3063c9af55f46936778709db5e6a6873d8": "0x0fac7d3063c9af55f46936778709db5e6a6873d8",
		"0xb3a0e0f6615cf3a1ecdb6442b1461f2a3a596dd6": "0xb3a0e0f6615cf3a1ecdb6442b1461f2a3a596dd6",
		"0x7912de41323375c54358b50b1dc2291046fa5f36": "0x7912de41323375c54358b50b1dc2291046fa5f36",
		"0xc2320bfd048f615885a08d14c6183dd1e7501067": "0xc2320bfd048f615885a08d14c6183dd1e7501067",
		"0xc0da92f2231f4a2f2b17f0001d35c2038be004ea": "0xc0da92f2231f4a2f2b17f0001d35c2038be004ea",
		"0xe73923a20b60e8f64b57ca00098e082a936a47f5": "0xe73923a20b60e8f64b57ca00098e082a936a47f5",
		"0xf3039972a4c2bf89fa9372400acce9b23a9fad2b": "0xf3039972a4c2bf89fa9372400acce9b23a9fad2b",
		"0xe88d8cdadd61af6a026647b37120a736208c53cd": "0xe88d8cdadd61af6a026647b37120a736208c53cd",
		"0xa144222f68468b4687f62e87a695b8ea65269f6e": "0xa144222f68468b4687f62e87a695b8ea65269f6e",
		"0xed9c104bb5c98afc8ff7054924cb5f90038d4546": "0xed9c104bb5c98afc8ff7054924cb5f90038d4546",
		"0xa626ca42d34607b49d5dbcddeaee16848822d0fc": "0xa626ca42d34607b49d5dbcddeaee16848822d0fc",
		"0x06b44ea040a980c16b45e5e732b4616eb373f266": "0x06b44ea040a980c16b45e5e732b4616eb373f266",
		"0xd1487ed119964cb5a71b38df55006d262f240f02": "0xd1487ed119964cb5a71b38df55006d262f240f02",
		"0xf2435b8e583598c9f6e148df5fb14f08013c95c8": "0xf2435b8e583598c9f6e148df5fb14f08013c95c8",
		"0x4ca95e2020f8bf1b41e2becd594cc5e4cf48459e": "0x4ca95e2020f8bf1b41e2becd594cc5e4cf48459e",
		"0xbc2e250d73dbb0cf9ecc007ee9eaf57f635f03fd": "0xbc2e250d73dbb0cf9ecc007ee9eaf57f635f03fd",
		"0x655be2207a1c411921993aa462ab2c6e1e53b47d": "0x655be2207a1c411921993aa462ab2c6e1e53b47d",
		"0xd703cdd219b13120402ef53adb31cd0316fc9837": "0xd703cdd219b13120402ef53adb31cd0316fc9837",
		"0xd37b6569a3334f4dc819b36cbbb019635bd9dffe": "0xd37b6569a3334f4dc819b36cbbb019635bd9dffe",
		"0x9b1d14b1bb78cebbf9efa878570b53534f543acf": "0x9b1d14b1bb78cebbf9efa878570b53534f543acf",
		"0xdb33a4d047d8bce4a930a699a68cdf89c50038b9": "0xdb33a4d047d8bce4a930a699a68cdf89c50038b9",
		"0xecb37fde0d2c1ba0d1c5fd3ce299286da941ca5a": "0xecb37fde0d2c1ba0d1c5fd3ce299286da941ca5a",
		"0x9554562f3f3ccaa6b7beb069db13c240b7e860aa": "0x9554562f3f3ccaa6b7beb069db13c240b7e860aa",
		"0x707e6b33abe2e0a166477b29e4f2a706a2c6f467": "0x707e6b33abe2e0a166477b29e4f2a706a2c6f467",
		"0x96591def961a43075c836c613d4bd127cd4d849e": "0x96591def961a43075c836c613d4bd127cd4d849e",
		"0x17ca6c3d07fcd7de48bef0bc71aad9c454257437": "0x17ca6c3d07fcd7de48bef0bc71aad9c454257437",
		"0x5677f8c0627bf61c0a680bb0939b6fad650d3118": "0x5677f8c0627bf61c0a680bb0939b6fad650d3118",
		"0x2132801129b3f1b3f0a5062c289afb823110d5d7": "0x2132801129b3f1b3f0a5062c289afb823110d5d7",
		"0x6f0d1f0a58e8d96c2bac60870af56f9c16162c51": "0x6f0d1f0a58e8d96c2bac60870af56f9c16162c51",
		"0x6ef0c7e1ce7ba4f3ba42abcfc0946e156044c4c6": "0x6ef0c7e1ce7ba4f3ba42abcfc0946e156044c4c6",
		"0x8e8462dfe73cdc4fd37b5ebfcfb5e3c2ac5d5f16": "0x8e8462dfe73cdc4fd37b5ebfcfb5e3c2ac5d5f16",
		"0x2a6676b43b9ec0e22ddec44759d1267288a1decf": "0x2a6676b43b9ec0e22ddec44759d1267288a1decf",
		"0x93a120bab5b4965478b0c8fe543f154bbd6c8399": "0x93a120bab5b4965478b0c8fe543f154bbd6c8399",
		"0xd5fdee60153e08cf9885f962f20143d611892e34": "0xd5fdee60153e08cf9885f962f20143d611892e34",
		"0x6aed258b27458e8b18b38ac01077ab8ea3ed3154": "0x6aed258b27458e8b18b38ac01077ab8ea3ed3154",
		"0x2f2e65615ad785c4ce408bac67ee17db0511c346": "0x2f2e65615ad785c4ce408bac67ee17db0511c346",
		"0x8e4904b24c1febfd0f83153c392bfdd10a752170": "0x8e4904b24c1febfd0f83153c392bfdd10a752170",
		"0x65e8cd1bf1b972a9ba039423f5a5543f602ff522": "0x65e8cd1bf1b972a9ba039423f5a5543f602ff522",
		"0x9342ae273da6b76ef29099f6f7b6aa1da5b0c7e8": "0x9342ae273da6b76ef29099f6f7b6aa1da5b0c7e8",
		"0xce68fdc0ba07578a3cf0e0752f5f751a9852377f": "0xce68fdc0ba07578a3cf0e0752f5f751a9852377f",
		"0x7b0ef48c52cd47f4f64ba58bcea08bdb8e164a23": "0x7b0ef48c52cd47f4f64ba58bcea08bdb8e164a23",
		"0x80f181a60f3da839efcef1b23917e7a25cd05908": "0x80f181a60f3da839efcef1b23917e7a25cd05908",
		"0x5689887d58d50f4f5629858d812d0e6b20cfe17c": "0x5689887d58d50f4f5629858d812d0e6b20cfe17c",
		"0xc1344e6847b46959bb1bc0ef777ff91099b9541b": "0xc1344e6847b46959bb1bc0ef777ff91099b9541b",
		"0xacd8bf1fc56ff5847230ae75f37c622b649b9e85": "0xacd8bf1fc56ff5847230ae75f37c622b649b9e85",
		"0xe54ef502c935aa5e37984014557255a5e1ea67ba": "0xe54ef502c935aa5e37984014557255a5e1ea67ba",
		"0xb9e34086cf97935592c35bd27ed0a082cf27e9cf": "0xb9e34086cf97935592c35bd27ed0a082cf27e9cf",
		"0xe44866203b61e2cf4843fca215ec8c7c60bae1b4": "0xe44866203b61e2cf4843fca215ec8c7c60bae1b4",
		"0xf8f04f3983b75540533677377b300aaa82fd73d1": "0xf8f04f3983b75540533677377b300aaa82fd73d1",
		"0xd4ff8685df5e0d7083bbebe0cf0c2b23b49a6ef5": "0xd4ff8685df5e0d7083bbebe0cf0c2b23b49a6ef5",
		"0x01cbf4631104c759e53f17b1cb0d7a871c247b72": "0x01cbf4631104c759e53f17b1cb0d7a871c247b72",
		"0xeeb2d61ff573ca0ea4eb13f93f248639bb057608": "0xeeb2d61ff573ca0ea4eb13f93f248639bb057608",
		"0x2d6f6e563c3e2bca868b1bcc2f0d8ad963db3279": "0x2d6f6e563c3e2bca868b1bcc2f0d8ad963db3279",
		"0x669874d40ebcf25cf8cb5619b3fee7fa4e615db1": "0x669874d40ebcf25cf8cb5619b3fee7fa4e615db1",
		"0x1021cc3032b7b61b2e321dc8240cc5ee38f5c5a5": "0x1021cc3032b7b61b2e321dc8240cc5ee38f5c5a5",
		"0xfe8d7d2b12783ae699e318ed9e76bce7c3612b7e": "0xfe8d7d2b12783ae699e318ed9e76bce7c3612b7e",
		"0x9d421b9966bfd8292bbe137dc04b09cc6c37db89": "0x9d421b9966bfd8292bbe137dc04b09cc6c37db89",
		"0x9d1c19468514133b4852328854122ac462b30dc5": "0x9d1c19468514133b4852328854122ac462b30dc5",
		"0x7089d371d5507e27321e675c4f56daec2c471b3e": "0x7089d371d5507e27321e675c4f56daec2c471b3e",
		"0x7bce84f0cd4c202bb16bb8f9ecc2e19c81218c65": "0x7bce84f0cd4c202bb16bb8f9ecc2e19c81218c65",
		"0xe3a7e932404d511b093ef7b0c4ce092bbb160799": "0xe3a7e932404d511b093ef7b0c4ce092bbb160799",
		"0x7aa3a81afe888f137ebfaa1c16a97be984742520": "0x7aa3a81afe888f137ebfaa1c16a97be984742520",
		"0xb3b018826608a0211c7963cac4473cd5c298a20f": "0xb3b018826608a0211c7963cac4473cd5c298a20f",
		"0x45b99dba03421664c2f4caef15bbef630b667b0d": "0x45b99dba03421664c2f4caef15bbef630b667b0d",
		"0xb487ab4a9ec3362cc079b1797db56f7adbd39604": "0xb487ab4a9ec3362cc079b1797db56f7adbd39604",
		"0x115682e2d741aad0997c786cd6f147d5da617610": "0x115682e2d741aad0997c786cd6f147d5da617610",
		"0x5f8f46074131e2fdc55778410feb5046958aa09b": "0x5f8f46074131e2fdc55778410feb5046958aa09b",
		"0x96d6fc35d3248bdc45bed813346dfb52798e3eea": "0x96d6fc35d3248bdc45bed813346dfb52798e3eea",
	}

	fmt.Println("Number of account to be migrated: ", len(oldAddrList))

	for _, account := range oldAddrList {

		fmt.Println("Migrating account: ", account.Hex())

		email, err := oldWalletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, account, "email")
		if err != nil {
			fmt.Println("cannot get email from " + account.Hex())
			continue
		}
		fmt.Println("Account: ", account.Hex(), " email: ", email)
		isMerchantStr, err := oldWalletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, account, "isMerchant")
		if err != nil {
			fmt.Println("cannot get isMerchant from " + account.Hex())
			continue
		}
		fmt.Println("Account: ", account.Hex(), " isMerchantStr: ", isMerchantStr)

		isMetaMaskStr, err := oldWalletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, account, "isMetaMask")
		if err != nil {
			fmt.Println("cannot get isMetaMask from " + account.Hex())
			continue
		}
		fmt.Println("Account: ", account.Hex(), " isMetaMaskStr: ", isMetaMaskStr)

		kycLevelStr, err := oldWalletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, account, "kycLevel")
		if err != nil {
			fmt.Println("cannot get kycLevel from " + account.Hex())
			continue
		}
		fmt.Println("Account: ", account.Hex(), " kycLevel: ", kycLevelStr)

		newAddr, ok := accountMapping[strings.ToLower(account.Hex())]
		if !ok {
			fmt.Println("Account not found in map: ", account.Hex())
			continue
		}

		var isMerchant bool
		var isMetamask bool
		if isMerchantStr == "true" {
			isMerchant = true
		}
		if email != "" {
			isMetamask = true
		}
		transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
		if err != nil {
			fmt.Println("GetNewTransactorFromPrivateKey failed on account: ", account.Hex())
			continue
		}
		transOpt.GasLimit = 400000
		tx, err := newWalletAddressMap.AddWalletInfo(transOpt, common.HexToAddress(newAddr), email, isMerchant, isMetamask)
		if err != nil {
			fmt.Println("AddWalletInfo failed on account: ", account.Hex())
			continue
		}

		receipt, err := ethClient.QueryEthReceipt(tx)
		if err != nil {
			fmt.Println("query receipt error: ", err)
			continue
		}

		receiptJson, _ := json.Marshal(receipt)
		fmt.Println("Account: ", account.Hex())
		fmt.Println(string(receiptJson))

		if kycLevelStr != "0" {
			fmt.Println("Going to set KYC level for account: ", account.Hex())
			transOpt1, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
			if err != nil {
				fmt.Println("GetNewTransactorFromPrivateKey failed for kyclevel account: ", account.Hex())
				continue
			}
			tx1, err := newWalletAddressMap.SetWalletInfo(transOpt1, common.HexToAddress(newAddr), "kycLevel", kycLevelStr)
			receipt1, err := ethClient.QueryEthReceipt(tx1)
			if err != nil {
				fmt.Println("query receipt error: ", err)
				continue
			}

			receiptJson1, _ := json.Marshal(receipt1)
			fmt.Println("Account: ", account.Hex())
			fmt.Println(string(receiptJson1))
		}
	}
}

func TestWalletAddressMap_DeleteEntry(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")

	walletAddressMap, _ := contract.NewWalletAddressMap(addr, ethClient.Client)
	transOpt, err := ethClient.GetNewTransactorFromPrivateKey(testOwnerPrivateKey, ethClient.ChainID)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := walletAddressMap.RemoveWalletInfo(transOpt, common.HexToAddress("0xe5ae39a9f014230d488a2c4b072bc94d9c163bbb"))
	if err != nil {
		t.Fatal(err)
	}
	queryEthReceipt(t, &ethClient, tx)

}

func TestWalletAddressMap_GetEmail(t *testing.T) {
	ethClient := initEthClient(t)
	addr := getAddressBySmartContractName("OwnedUpgradeabilityProxy<WalletAddressMap>")

	walletAddressMap, _ := contract.NewWalletAddressMap(addr, ethClient.Client)

	val, err := walletAddressMap.GetWalletInfoValue(&bind.CallOpts{}, common.HexToAddress("0x35CD8522CFDA982e90AB6859a7FfE8C0b350B3D1"), "email")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Email: ", val)

}
