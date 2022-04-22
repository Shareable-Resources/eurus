package user

import (
	"eurus-backend/foundation/log"
	"eurus-backend/sign_service/sign_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func IsWalletAddressExist(address string, server *UserServer) (bool, error) {
	addr, err := GetWalletAddressMapAddress(server)
	if err != nil {
		return false, err
	}
	walletAddressInstance, err := contract.NewWalletAddressMap(*addr, server.EthClient.Client)
	if err != nil {
		return false, err
	}
	isExist, err := walletAddressInstance.IsWalletAddressExist(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}
	return true, nil
}

func GetWalletAddressMapAddress(server *UserServer) (*common.Address, error) {
	internalSC, err := contract.NewInternalSmartContractConfig(common.HexToAddress(server.ServerConfig.InternalSCConfigAddress), server.EthClient.Client)
	if err != nil {
		return nil, err
	}

	addr, err := internalSC.GetWalletAddressMap(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &addr, err
}

func AddUserAddressToWalletSC(server *UserServer, user *User, userIsMerchant bool, userIsMetaMask bool) error {

	userId := user.Id

	return retryFunction(server.Config, func() (bool, error) {

		walletAddressMap, err := contract.NewWalletAddressMap(common.HexToAddress(server.ServerConfig.WalletAddressAddress), server.EthClient.Client)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("NewWalletAddressMap error. userId: ", userId)
			return true, wrapErrorWithUserId(err, "NewWalletAddressMap error", userId)
		}

		isWalletExists, err := walletAddressMap.IsWalletAddressExist(&bind.CallOpts{}, common.HexToAddress(server.ServerConfig.WalletAddressAddress))
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("IsWalletAddressExist error: ", err, " userId: ", userId)
			return false, wrapErrorWithUserId(err, "IsWalletAddressExist", userId)
		}
		if isWalletExists {
			log.GetLogger(log.Name.Root).Infoln("Wallet address already exists at WalletAddressMap. userId: ", userId)
			return false, nil
		}

		transOpt, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromSignServer. userId: ", userId)
			return false, wrapErrorWithUserId(err, "GetNewTransactorFromSignServer", userId)
		}

		transOpt.GasLimit = uint64(server.Config.SideChainGasLimit)
		log.GetLogger(log.Name.Root).Debugln("AddUserAddressToWalletSC gas limit: ", transOpt.GasLimit)
		tx, err := walletAddressMap.AddWalletInfo(transOpt, common.HexToAddress(user.WalletAddress), user.Email, userIsMerchant, userIsMetaMask)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddWalletInfo error: ", err, " userId: ", userId)
			return false, wrapErrorWithUserId(err, "AddWalletInfo", userId)
		}
		log.GetLogger(log.Name.Root).Infoln("AddWalletInfo txHash: ", tx.Hash().Hex(), " userId: ", userId)

		receipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("AddWalletInfo query receipt error: ", err, " txHash: ", tx.Hash().Hex(), " userId: ", userId)
			return false, wrapErrorWithUserId(err, "AddWalletInfo query receipt error. txHash: "+tx.Hash().Hex(), userId)
		}

		if receipt.Status == 0 {
			log.GetLogger(log.Name.Root).Errorln("Query AddWalletInfo receipt status is 0. txHash: ", tx.Hash().Hex(), " userId: ", userId)
			return false, wrapErrorWithUserId(err, "Query AddWalletInfo receipt status is 0. txHash: "+tx.Hash().Hex(), userId)
		}
		return false, nil
	}, userId)
}

func SetUserWalletKycLevel(server *UserServer, userWalletAddress common.Address, walletAddressMapAddr common.Address, kycLevel string) (*types.Transaction, error) {
	transOpt, err := server.EthClient.GetNewTransactorFromSignServer(server.AuthClient, server.Config.SignServerUrl, sign_api.WalletKeyUserWalletOwner)
	if err != nil {
		return nil, err
	}

	userWalletSC, err := contract.NewWalletAddressMap(walletAddressMapAddr, server.EthClient.Client)
	if err != nil {
		return nil, err
	}
	transOpt.GasLimit = uint64(server.Config.SideChainGasLimit)
	tx, err := userWalletSC.SetWalletInfo(transOpt, userWalletAddress, "kycLevel", kycLevel)
	if err != nil {
		return nil, err
	}

	fmt.Println("set key level tx: ", tx.Hash())
	log.GetLogger(log.Name.Root).Infoln("set kyc level tx: ", tx.Hash())

	receipt, err := server.EthClient.QueryEthReceiptWithSetting(tx, 1, 20)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet owner")
		return nil, err
	}
	if receipt.Status != 1 {
		log.GetLogger(log.Name.Root).Error("transaction failed:")
		return nil, errors.New("set kycLevel failed for hx: " + receipt.TxHash.Hex())
	}

	return tx, nil

}
