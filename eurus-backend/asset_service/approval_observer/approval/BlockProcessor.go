package approval

import (
	"encoding/hex"
	eurus_ethereum "eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ApprovalBlockProcessor struct {
	EthWebSocketClient *eurus_ethereum.EthClient
}
type ApprovalWalletTopic struct {
	SrcWallet        common.Address
	DestWallet       common.Address
	SubmitterAddress common.Address
	AssetName        string
	TransId          *big.Int
	Amount           *big.Int
	FeeAmount        *big.Int
	UserGasUsed      *big.Int
}

func (me *ApprovalObserver) TestExternal() (asset []string, address []common.Address, err error) {
	ethClient := me.EthClient
	externalSC, err := contract.NewExternalSmartContractConfig(common.HexToAddress(me.ServerConfig.ExternalSCConfigAddress), ethClient.Client)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Failed to call Smart Contract NewExternalSmartContractConfig : ", err)
	}
	assetList, addressList, err := externalSC.GetAssetAddress(&bind.CallOpts{})

	var allowAssetList []string
	var allowAddressList []common.Address
	for i, receAddress := range addressList {
		checkERC20, err := externalSC.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, assetList[i])
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("The assets can not check by ERC20 smart contract", err.Error())
		}
		addressHandle := hex.EncodeToString(checkERC20.Bytes())
		if addressHandle == "0000000000000000000000000000000000000000" {
			log.GetLogger(log.Name.Root).Errorln("This assets not accept")
			return allowAssetList, allowAddressList, err
		}
		log.GetLogger(log.Name.Root).Debugln("The asset type is verified. Addr: ", receAddress.String(), " currency: ", assetList[i])
		allowAssetList = append(allowAssetList, assetList[i])
		allowAddressList = append(allowAddressList, receAddress)
	}
	return allowAssetList, allowAddressList, nil
}
