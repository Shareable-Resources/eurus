package sign

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/secret"
	"eurus-backend/sign_service/sign_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func signTransaction(server *SignServer, req *sign_api.SignTransactionRequest) *response.ResponseBase {

	privateKey, err := crypto.HexToECDSA(server.Config.UserWalletOwnerPrivateKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get private key : " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	txBin := req.Data
	var tx *types.Transaction = new(types.Transaction)

	err = tx.UnmarshalBinary(txBin)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to UnmarshalBinary transaction :  " + err.Error())
		return response.CreateErrorResponse(req, foundation.InvalidArgument, err.Error())
	}

	serviceId, err := server.AuthClient.(*auth.AuthClient).GetServiceIdFromServerLoginToken(req.GetLoginToken())
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, err.Error())
	}

	isValid, serverErr := server.AuthClient.VerifySignature(txBin, req.Sign, serviceId)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Error("Unable to verify signature: ", serverErr.Message, " code: ", serverErr.ReturnCode)
		return response.CreateErrorResponse(req, serverErr.GetReturnCode(), serverErr.Message)
	}
	if !isValid {
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "Invalid signature")
	}
	keyPair, err := ethereum.GetEthKeyPair(server.Config.UserWalletOwnerPrivateKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get wallet address : " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	//Get the pending nonce here
	nonce, err := server.EthClient.Client.PendingNonceAt(context.Background(), keyPair.Address)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get pending nonce : " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	} else {
		log.GetLogger(log.Name.Root).Debug("User owner wallet pending nonce: ", nonce, " req: ", req.Nonce)
	}

	newTxData := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: tx.GasPrice(),
		Gas:      tx.Gas(),
		To:       tx.To(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	newTx := types.NewTx(newTxData)

	signer := types.LatestSignerForChainID(big.NewInt(int64(server.Config.EthClientChainID)))

	txHash := signer.Hash(newTx)

	signature, err := crypto.Sign(txHash.Bytes(), privateKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to sign transaction :  " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	res := new(sign_api.SignTransactionResponse)
	signedTx, err := newTx.WithSignature(signer, signature)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Create signed transaction error: ", err, " req: ", req.Nonce)
	}
	signedTxBin, err := signedTx.MarshalBinary()
	if err != nil {
		log.GetLogger(log.Name.Root).Error("MarshalBinary transaction error: ", err, " req: ", req.Nonce)
	}
	res.SignedTrans = signedTxBin
	return response.CreateSuccessResponse(req, res)

}

func getWalletAddress(server *SignServer, req *sign_api.QueryAddressRequest) *response.ResponseBase {

	var addr string
	if req.KeyType == sign_api.WalletKeyUserWalletOwner {
		keyPair, err := ethereum.GetEthKeyPair(server.Config.UserWalletOwnerPrivateKey)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to get wallet address : " + err.Error())
			return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		}
		addr = keyPair.Address.Hex()
	} else if req.KeyType == sign_api.WalletKeyInvoker {
		addr = server.Config.invokerAddress
	} else {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Unsupported wallet key type")
	}

	res := new(sign_api.QueryAddressResponse)
	res.Address = addr
	return response.CreateSuccessResponse(req, res)
}

func getPendingNonce(server *SignServer, req *sign_api.GetPendingNonceRequest) *response.ResponseBase {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	res := new(sign_api.GetNonceResponse)

	if req.KeyType == sign_api.WalletKeyUserWalletOwner {
		res.PendingNonce = 0
		//Decrepcated

		// if server.UserWalletOwnerNonce == 0 {
		// keyPair, err := ethereum.GetEthKeyPair(server.Config.UserWalletOwnerPrivateKey)
		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Error("Unable to get key pair : " + err.Error())
		// 	return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		// }
		// nonce, err := server.EthClient.Client.PendingNonceAt(context.Background(), keyPair.Address)

		// if err != nil {
		// 	log.GetLogger(log.Name.Root).Error("Unable to get pending nonce : " + err.Error())
		// 	return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		// }
		// res.PendingNonce = nonce
		// server.UserWalletOwnerNonce = nonce + 1
		// } else {
		// 	res.PendingNonce = server.UserWalletOwnerNonce
		// 	server.UserWalletOwnerNonce++
		// }
	} else if req.KeyType == sign_api.WalletKeyInvoker {
		if server.InvokerNonce == 0 {
			keyPair, err := ethereum.GetEthKeyPair(server.Config.InvokerPrivateKey)
			if err != nil {
				log.GetLogger(log.Name.Root).Error("Unable to get key pair : " + err.Error())
				return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
			}
			nonce, err := server.EthClient.Client.PendingNonceAt(context.Background(), keyPair.Address)

			if err != nil {
				log.GetLogger(log.Name.Root).Error("Unable to get pending nonce : " + err.Error())
				return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
			}
			res.PendingNonce = nonce
			server.InvokerNonce = nonce + 1
		} else {
			res.PendingNonce = server.InvokerNonce
			server.InvokerNonce++
		}
	} else {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Unsupported wallet key type")
	}

	return response.CreateSuccessResponse(req, res)
}

func calibrateNonce(server *SignServer, req *sign_api.CalibrateNonceRequest) *response.ResponseBase {

	nonce, err := calibrate(server, req.KeyType)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to calibrate nonce : " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	res := &sign_api.CalibrateNonceResponse{
		Nonce:   nonce,
		KeyType: req.KeyType,
	}
	return response.CreateSuccessResponse(req, res)
}

func calibrate(server *SignServer, keyType sign_api.WalletKeyType) (uint64, error) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	var keyPair *ethereum.ECDSAKeyPair
	var err error
	if keyType == sign_api.WalletKeyUserWalletOwner {
		keyPair, err = ethereum.GetEthKeyPair(server.Config.UserWalletOwnerPrivateKey)
	} else if keyType == sign_api.WalletKeyInvoker {
		keyPair, err = ethereum.GetEthKeyPair(server.Config.InvokerPrivateKey)
	} else {
		return 0, errors.New("Unsupported wallet key type")
	}

	if err != nil {
		return 0, err
	}
	nonce, err := server.EthClient.Client.NonceAt(context.Background(), keyPair.Address, nil)
	if err != nil {
		return 0, err
	}
	if keyType == sign_api.WalletKeyUserWalletOwner {
		// server.UserWalletOwnerNonce = nonce
	} else {
		server.InvokerNonce = nonce
	}

	return nonce, nil
}

func signUserWalletTransaction(server *SignServer, req *sign_api.SignUserWalletTransactionRequest) *response.ResponseBase {
	paramStr := strings.TrimPrefix(req.InputFunction, "0x")
	inputParam, err := hex.DecodeString(paramStr)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Unable to decode inputFunction")
	}

	userWalletAbi := ethereum.DefaultABIDecoder.GetABI("UserWallet")
	if len(inputParam) < 4 {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid inputFunction")
	}
	var isFound bool
	var functionName string
	for funcName, method := range userWalletAbi.Methods {

		if bytes.Equal(inputParam[:4], method.ID) {
			isFound = true
			functionName = funcName
			break
		}
	}
	if !isFound {
		return response.CreateErrorResponse(req, foundation.MethodNotFound, "Contract function not found")
	}

	_, err, _ = ethereum.DefaultABIDecoder.DecodeABIInputArgument(inputParam, "UserWallet", functionName)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid inputFunction argument")
	}

	if req.To == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "To address is empty")
	}

	userWallet, err := contract.NewUserWallet(common.HexToAddress(req.To), server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid wallet address at request nonce: ", req.GetNonce(), " Error: ", err.Error())
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid wallet address at request nonce: "+req.GetNonce()+" Error: "+err.Error())
	}

	isWriter, err := userWallet.IsWriter(&bind.CallOpts{}, common.HexToAddress(server.Config.invokerAddress))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Validate user wallet failed at nonce: ", req.GetNonce(), " Error: ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "Validate user wallet failed at nonce: "+req.GetNonce()+" Error: "+err.Error())

	}

	if !isWriter {
		log.GetLogger(log.Name.Root).Errorln("No permission to call user wallet at nonce: ", req.GetNonce())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "No permission to call user wallet at nonce: "+req.GetNonce())

	}

	transOpt, err := server.EthClient.GetNewTransactorFromPrivateKey(server.Config.InvokerPrivateKey, server.EthClient.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey failed at request nonce: ", req.GetNonce())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	val := big.NewInt(0)
	if req.Value != "" {
		val.SetString(req.Value, 10)
	}

	gasPrice := big.NewInt(0)
	gasPrice = gasPrice.SetUint64(req.GasPrice)

	tx := types.NewTransaction(transOpt.Nonce.Uint64(), common.HexToAddress(req.To), val, server.Config.SideChainGasLimit, gasPrice, inputParam)

	signedTx, err := transOpt.Signer(common.HexToAddress(server.Config.invokerAddress), tx)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Sign transaction error on nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "Sign error: "+err.Error())
	}

	rlpEncoded, err := signedTx.MarshalBinary()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Encode signed transaction error on nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "Encode signed transaction error: "+err.Error())
	}

	rlpTxStr := hex.EncodeToString(rlpEncoded)
	rlpTxStr = "0x" + rlpTxStr
	res := &sign_api.SignedUserWalletTransactionResponse{SignedTx: rlpTxStr}
	return response.CreateSuccessResponse(req, res)
}

func getCentralizedUserMainnetAddress(config *SignServerConfig, req *sign_api.GetCentralizedUserMainnetAddressRequest) *response.ResponseBase {
	_, addr, err := secret.GenerateCentralizedUserMainnetPrivateKey(config.CentralizedUserWalletMnemonicPhase, req.UserId)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, err.Error())
	}

	return response.CreateSuccessResponse(req, addr)
}
