package user

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"eurus-backend/auth_service/auth"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/mainnet_smart_contract/build/golang/mainnet_contract"
	"eurus-backend/marketing/reward"
	"eurus-backend/secret"
	"eurus-backend/sign_service/sign_api"
	"eurus-backend/smartcontract/build/golang/contract"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"

	go_crypto "crypto"

	go_ethereum_crypto "github.com/ethereum/go-ethereum/crypto"
)

func TransferEUNToUser(server *UserServer, userId uint64, walletAddress string) error {

	if server.Config.InitialFundExactAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	_, receipt, err := server.rewardProcessor.TransferEUN(walletAddress, server.Config.MarketRegWalletAddress, server.Config.InitialFundExactAmount)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to reward user: ", err)
		return err
	}
	if receipt.Status == 0 {
		receiptData, _ := json.Marshal(receipt)
		log.GetLogger(log.Name.Root).Errorln("Reward user receipt status is 0: ", string(receiptData))
		return errors.New("Reward user receipt status is 0:")
	}

	distributedToken := new(reward.DistributedToken)
	distributedToken.AssetName = "EUN"
	distributedToken.Amount = decimal.NewFromBigInt(server.Config.InitialFundExactAmount, 0)
	chainId := server.EthClient.ChainID.Uint64()
	distributedToken.Chain = &chainId
	distributedToken.DistributedType = reward.DistributedRegistration
	distributedToken.TriggerType = reward.TriggerNotApplicable
	distributedToken.UserId = userId
	distributedToken.TxHash = strings.ToLower(receipt.TxHash.Hex())
	distributedToken.FromAddress = server.Config.HdWalletAddress
	distributedToken.ToAddress = ethereum.ToLowerAddressString(walletAddress)
	gasPriceDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0)
	distributedToken.GasPrice = &gasPriceDec
	distributedToken.GasUsed = receipt.GasUsed

	bigGasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
	gasFeeDec := decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0).Mul(decimal.NewFromBigInt(bigGasUsed, 0))

	distributedToken.GasFee = &gasFeeDec
	distributedToken.InitDate()

	err = server.rewardProcessor.DbProcessor.DbInsertDistributedToken(distributedToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to insert distributed_token: ", err)
	}

	return nil
}

func ImportWallet(server *UserServer, req *ImportWalletRequest, remoteAddr string) (*response.ResponseBase, bool) {
	verified, err := importWalletVerifySign(req.WalletBaseRequest)
	if err != nil {
		res := response.CreateErrorResponse(req, foundation.InvalidSignature, err.Error())
		return res, false
	}
	req.WalletAddress = ethereum.ToLowerAddressString(req.WalletAddress)

	decompressedPubKeyStr, err := crypto.DecompressPubKey(req.PublicKey)
	if err != nil {
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res, false
	}

	decompressedPubKey, _ := hex.DecodeString(decompressedPubKeyStr)
	pubKey, err := go_ethereum_crypto.UnmarshalPubkey(decompressedPubKey)
	if err != nil {
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, err.Error())
		return res, false
	}

	userLoginAddress := go_ethereum_crypto.PubkeyToAddress(*pubKey)
	if req.WalletAddress != ethereum.ToLowerAddressString(userLoginAddress.Hex()) {
		res := response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
		return res, false
	}

	var res *response.ResponseBase = nil
	var isValid bool = false
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to verify the sign: "+err.Error(), "\nRequest Params: "+string(reqStr))
		res = response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	} else if verified {
		res = walletAddressHandler(server, req, remoteAddr)
		if res.ReturnCode == int64(foundation.Success) {
			isValid = true
		}
	} else {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Invalid Sign: ", string(reqStr))
		res = response.CreateErrorResponse(req, foundation.SignMatchError, "Invalid Sign!")
	}
	return res, isValid

}

// var WalletSCTxRetryNum = 5
// var durationStr string = "5s"

// func TxRetryerDuration() time.Duration {
// 	duration, _ := time.ParseDuration(durationStr)
// 	return duration
// }

// func WalletAddressInsertionRetryer(addUserAddressFn func(*UserServer, *User, bool, bool) (string, error), server *UserServer, user *User, currentRetryCount int, retryConfig foundation.IRetrySetting, userIsMerchant bool, userIsMetaMask bool) (string, error) {
// 	txHash, err := addUserAddressFn(server, user, false, true)
// 	if currentRetryCount != 0 && err != nil {
// 		log.GetLogger(log.Name.Root).Error("AddUserAddressToWalletSC failed: ", err, " retry count: ", currentRetryCount, " user id: ", user.Id, " wallet address: ", user.WalletAddress)
// 		currentRetryCount--
// 		time.Sleep(time.Second * retryConfig.GetRetryInterval())
// 		return WalletAddressInsertionRetryer(addUserAddressFn, server, user, currentRetryCount, retryConfig, userIsMerchant, userIsMetaMask)
// 	} else if err != nil {
// 		log.GetLogger(log.Name.Root).Error("AddUserAddressToWalletSC failed: ", err, " retry count: 0 user id: ", user.Id, " wallet address: ", user.WalletAddress)
// 		return "", err
// 	}
// 	return txHash, nil
// }

func walletAddressHandler(server *UserServer, req *ImportWalletRequest, remoteAddr string) *response.ResponseBase {
	isExisted, err := isWalletAddressExisted(req.WalletAddress, server)
	var token auth_base.ILoginToken
	var res *response.ResponseBase = nil
	var ires *ImportWalletResponse = nil
	var txHash string
	user := new(User)
	var dbRecordFound bool = false
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to check if WalletAddress is Existed:"+err.Error(), "\nRequest Params: "+string(reqStr))
		res = response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
	} else if isExisted {
		user, err = DbGetUserByWalletAddress(req.WalletAddress, server.SlaveDatabase)
		if err != nil {
			if serverErr, ok := err.(*foundation.ServerError); ok && serverErr.ReturnCode == foundation.UserNotFound {
				log.GetLogger(log.Name.Root).Debugln("User not found")
				dbRecordFound = false
				err = nil
				user = new(User)
			} else {
				log.GetLogger(log.Name.Root).Debugln("DbGetUserByWalletAddress error: ", err)
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to get user from db:"+err.Error(), "\nRequest Params: "+string(reqStr))
				res = response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
			}
		} else {
			dbRecordFound = true
			log.GetLogger(log.Name.Root).Debugln("DbGetUserByWalletAddress record found")
			token, err = getAuthLoginToken(server, user)
			if err != nil {
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to generate new token:"+err.Error(), "\nRequest Params: "+string(reqStr))
				res = response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
			}
		}
	}

	if err == nil && !dbRecordFound {
		user.WalletAddress = req.WalletAddress
		user.LoginAddress = req.WalletAddress

		if !isExisted {
			//Add user wallet address to WalletAddressMap
			err = AddUserAddressToWalletSC(server, user, false, true)
			if err != nil {
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to add new user to WalletAddress smart contract:"+err.Error(), "\nRequest Params: "+string(reqStr))
				res = response.CreateErrorResponse(req, foundation.NetworkError, err.Error())
			}
		}

		if err == nil {
			var err1 error
			user, err1 = DbAddNewUser(req.WalletAddress, server.DefaultDatabase, false, "", true)
			if err1 != nil {
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to add new user:"+err1.Error(), "\nRequest Params: "+string(reqStr))
				res = response.CreateErrorResponse(req, foundation.DatabaseError, err1.Error())
			} else {
				err = TransferEUNToUser(server, user.Id, req.WalletAddress)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Transfer EUN failed:" + err.Error())
				}
			}
		}
	}

	if res == nil {

		token, err = getAuthLoginToken(server, user)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("Unable to generate new token:"+err.Error(), "\nRequest Params: "+string(reqStr))
			res = response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		}

		_, err := DbUpdateLoginTime(req.WalletAddress, server.DefaultDatabase)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("Unable to update user login time:"+err.Error(), "\nRequest Params: "+string(reqStr))
		}

		ires = NewImportWalletResponse(token, txHash, user, req)
		res = response.CreateSuccessResponse(req, ires)
		go func() {
			importWalletData := NewLoginDataFromLoginBySignature(req.LoginLogDetail, user.Id, remoteAddr)
			server.elasticLogger.InsertLog(importWalletData)
		}()
	}

	return res
}

type TokenString struct {
	LoginAddress string `json:"loginAddress"`
	UserId       uint64 `json:"userId"`
}

func getAuthLoginToken(server *UserServer, user *User) (auth_base.ILoginToken, error) {
	if user.Id == 0 {
		return nil, foundation.NewError(foundation.UserNotFound)
	}
	tokenStr := TokenString{LoginAddress: user.LoginAddress, UserId: user.Id}
	tokenBytes, err := json.Marshal(tokenStr)
	if err != nil {
		return nil, err
	}
	token, err := server.AuthClient.GenerateLoginToken(string(tokenBytes))
	if err != nil {
		return nil, err
	}
	return token, err
}

func isWalletAddressExisted(address string, server *UserServer) (bool, error) {
	isExist, err := IsWalletAddressExist(address, server)
	return isExist, err
}

func verifyToken(authClient auth_base.IAuth, token string) (bool, auth_base.ILoginToken, error) {
	isValid, loginToken, err := authClient.VerifyLoginToken(token, auth_base.VerifyModeUser)
	if err == nil {
		return isValid, loginToken, nil
	}
	return isValid, loginToken, err
}

func VerifySignature(deviceId string, timestamp int64, address string, signature string, isPersonalSign bool, publicKey string) (bool, error) {
	deviceIdStr := "deviceId=" + deviceId
	timestampStr := "timestamp=" + strconv.FormatInt(timestamp, 10)
	addrStr := "walletAddress=" + address
	baseStr := deviceIdStr + "&" + timestampStr + "&" + addrStr

	if isPersonalSign {
		prefix := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(baseStr))
		baseStr = prefix + baseStr
	}
	log.GetLogger(log.Name.Root).Debugln("Signature base string: ", baseStr)

	sign, err := hex.DecodeString(signature)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to verify the sign: "+err.Error(), "\ndeviceIdStr: "+deviceIdStr+" timestampStr: "+timestampStr+" addrStr: "+addrStr)
		return false, err
	}

	if len(sign) > 64 {
		sign = sign[:64]
	}
	verified, err := crypto.VerifyECDSASignatureByHexKey(publicKey, []byte(baseStr), sign)
	return verified, err
}

func RegisterVerifySign(req RegistrationRequest) (bool, error) {
	isValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.LoginAddress, req.Signature, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	return true, nil
}

func importWalletVerifySign(req WalletBaseRequest) (bool, error) {
	isValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.WalletAddress, req.Sign, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	return true, nil
}

func RefreshLoginToken(authClient auth_base.IAuth, req *RefreshTokenRequest) *response.ResponseBase {
	loginToken, err := authClient.RefreshLoginToken(req.LoginToken.GetToken())
	if err != nil {
		// jsonByte, err1 := json.Marshal(req)
		// var jsonStr string
		// if err1 == nil {
		// 	jsonStr = string(jsonByte)
		// }
		// log.GetLogger(log.Name.Root).Error("Unable to refresh token: ", err.Error(), " req: ", jsonStr)
		return response.CreateErrorResponse(req, err.GetReturnCode(), err.Error())
	}

	res := RefreshTokenResponse{}
	res.ExpireTime = loginToken.GetExpiredTime()
	res.CreatedDate = loginToken.GetCreatedDate()
	res.LastModifiedDate = loginToken.GetLastModifiedDate()
	res.Token = loginToken.GetToken()

	return response.CreateSuccessResponse(req, res)
}

func LoginBySignature(server *UserServer, req *UserLoginBySignatureRequest, remoteAddr string) *response.ResponseBase {
	isValid, err := importWalletVerifySign(req.WalletBaseRequest)
	req.WalletAddress = ethereum.ToLowerAddressString(req.WalletAddress) //this is login address

	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to verify the sign: "+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	if !isValid {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Invalid Sign: ", string(reqStr))
		return response.CreateErrorResponse(req, foundation.SignMatchError, "Invalid Sign!")
	}

	checkUser, err := DbGetUserByLoginAddress(req.WalletAddress, server.SlaveDatabase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	if checkUser.Id == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	} else if checkUser.Status != UserStatusNormal {
		var code foundation.ServerReturnCode

		switch checkUser.Status {
		case UserStatusDeleted:
			code = foundation.UserNotFound
		case UserStatusSuspended:
			code = foundation.UserStatusForbidden
		default:
			code = foundation.UserNotFound
		}
		return response.CreateErrorResponse(req, code, code.String())
	}

	token, err := getAuthLoginToken(server, checkUser)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to update user login time:"+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	expiryTime := time.Unix(token.GetExpiredTime(), 0)
	// expiryTimeInRFC3339 := expiryTime.Format(time.RFC3339) // converts utc time to RFC3339 format
	mnemonicPhase := ""
	if !checkUser.IsMetamaskAddr {
		// Public key will be empty string if the device id is invalid
		publicKey, err := DbGetUserDevicePublicKey(server.SlaveDatabase, checkUser.Id, int16(CustomerUser), req.DeviceId)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("Unable to get user public key:"+err.Error(), "\nRequest Params: "+string(reqStr))
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}

		// If public key cannot be found, just keep mnemonic phase send back to user blank
		if publicKey != "" {
			plainTextMnemonic, err := decryptMnemonic(checkUser.Mnemonic, server.Config.UserMnenomicAesKey)
			if err != nil {
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to get decrypt user payment owner mnemonic: "+err.Error(), " create a dummy one.\nRequest Params: "+string(reqStr))
				plainTextMnemonic = crypto.Generate128BitsEntropyMnemonicPhrase()
				encryptedUserPaymentMnemonic, err := encryptMnemonicWithRSAAesKey(plainTextMnemonic, server.Config.UserMnenomicAesKey)
				if err != nil {
					log.GetLogger(log.Name.Root).Error("Unable to encrypt dummy user payment owner mnenomic phase: ", err)
				} else {
					err = DbUpdateUserMnemonicPhase(server.DefaultDatabase, checkUser.Id, encryptedUserPaymentMnemonic)
					if err != nil {
						log.GetLogger(log.Name.Root).Error("DbUpdateUserMnemonicPhase failed: ", err)
					}
				}
			}

			encryptedForUser, err := encryptMnemonicWithPublicKey(plainTextMnemonic, publicKey)
			if err != nil {
				reqStr, _ := json.Marshal(req)
				log.GetLogger(log.Name.Root).Error("Unable to encrypt user mnemonic:"+err.Error(), "\nRequest Params: "+string(reqStr))
				return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
			}

			mnemonicPhase = encryptedForUser
		}
	}

	user, err := DbUpdateLoginTime(req.WalletAddress, server.DefaultDatabase)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to update user login time:"+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	if user.Id == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	}

	if user.Id > 0 {
		go func() {
			loginData := NewLoginDataFromLoginBySignature(req.LoginLogDetail, user.Id, remoteAddr)
			server.elasticLogger.InsertLog(loginData)

		}()
	}

	ures := UserLoginBySignatureResponse{Token: token.GetToken(), ExpiryTime: expiryTime.Unix(), LastLoginTime: user.LastLoginTime.Unix(), Status: int(UserStatusNormal), LastModifiedDate: token.GetLastModifiedDate(), Mnemonic: mnemonicPhase, WalletAddress: user.WalletAddress, MainnetWalletAddress: user.MainnetWalletAddress, OwnerWalletAddress: user.OwnerWalletAddress, IsMetaMaskUser: user.IsMetamaskAddr}
	return response.CreateSuccessResponse(req, ures)
}

func GetServerConfig(req *QueryServerConfigRequest, server *UserServer) *response.ResponseBase {
	serverConfig := new(ServerConfig)
	serverConfig.EurusRPCProtocol = server.Config.UserEthClientProtocol
	serverConfig.EurusRPCDomain = server.Config.UserEthClientIP
	serverConfig.EurusRPCPort = server.Config.UserEthClientPort
	serverConfig.EurusChainId = server.ServerConfig.EthClientChainID
	serverConfig.ExternalSmartContractConfigAddress = server.ServerConfig.ExternalSCConfigAddress
	serverConfig.EurusInternalConfigAddress = server.ServerConfig.EurusInternalConfigAddress
	serverConfig.MainnetRPCDomain = server.Config.UserMainnetEthClientIP
	serverConfig.MainnetRPCPort = server.Config.UserMainnetEthClientPort
	serverConfig.MainnetRPCProtocol = server.Config.UserMainnetEthClientProtocol
	serverConfig.MainnetChainId = server.Config.MainnetEthClientChainID

	return response.CreateSuccessResponse(req, serverConfig)
}

func GetUserDetails(server *UserServer, req *QueryUserDetailsRequest) *response.ResponseBase {
	var res *response.ResponseBase = nil

	isValid, token, err := verifyToken(server.AuthClient, req.Token)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to verify token: ", string(reqStr))
		res = response.CreateErrorResponse(req, foundation.InternalServerError, "Unable to verify token")
	} else if !isValid {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Invalid token", string(reqStr))
		res = response.CreateErrorResponse(req, foundation.LoginTokenInvalid, "Invalid Token!")
	} else {
		userIdStr := token.GetUserId()
		userIdObj := new(UserLoginId)
		json.Unmarshal([]byte(userIdStr), userIdObj)
		userId := userIdObj.UserId
		user, err := DbGetUserById(userId, server.SlaveDatabase)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("Unable to get user from database", string(reqStr))
			res = response.CreateErrorResponse(req, foundation.DatabaseError, "Unable to get user from database")
		} else {
			res = response.CreateSuccessResponse(req, user)
		}
	}
	return res
}

func GetWithdrawAdminFee(server *UserServer, req *QueryWithdrawAdminFee, symbol string) *response.ResponseBase {
	//TODO Check symbol if EUN
	isValid, _, err := verifyToken(server.AuthClient, req.Token)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to verify token: ", string(reqStr))
		return response.CreateErrorResponse(req, foundation.InternalServerError, "Unable to verify token")
	} else if !isValid {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Invalid token", string(reqStr))
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, "Invalid Token!")
	} else {
		withdrawAdminFee := new(WithdrawAdminFee)

		decimal, err := DbGetAdminFeeDecimal(server.SlaveDatabase, symbol)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("cannot get asset decimal from DB ", string(reqStr), " symbol: ", symbol)
			return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get asset decimal from DB")
		}

		fee, err := server.GetAdminFeeFromSC(symbol)
		if err != nil {
			reqStr, _ := json.Marshal(req)
			log.GetLogger(log.Name.Root).Error("cannot get admin fee from SC", string(reqStr))
			return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get admin fee from SC")
		}

		withdrawAdminFee.Currency = symbol
		withdrawAdminFee.Fee = *fee
		withdrawAdminFee.Decimal = decimal
		return response.CreateSuccessResponse(req, withdrawAdminFee)
	}
}

func Faucet(server *UserServer, req *request.RequestBase, symbol string) *response.ResponseBase {

	mainnetClient := &ethereum.EthClient{
		Protocol: server.Config.MainnetEthClientProtocol,
		IP:       server.Config.MainnetEthClientIP,
		Port:     server.Config.MainnetEthClientPort,
		ChainID:  big.NewInt(int64(4)),
	}
	mainnetClient.Connect()

	userId, err := GetUserIdFromLoginToken(req.LoginToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	userFaucetsInstance, canFaucet, isPending, err := DbCheckIfUserAllowFaucet(server.SlaveDatabase, symbol, userId)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot check user faucet validity", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if isPending {
		return response.CreateSuccessResponse(req, FaucetResponse{Status: 2, TxHash: userFaucetsInstance.TransHash})
	}

	if !canFaucet {
		return response.CreateSuccessResponse(req, FaucetResponse{Status: 3, TxHash: userFaucetsInstance.TransHash})
	}

	userFromDB, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user from db", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	//TODO: 👇🏻 below code should be put in SCProcessor
	eurusInternalConfig, err := mainnet_contract.NewEurusInternalConfig(common.HexToAddress(server.Config.EurusInternalConfigAddress), mainnetClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot create externalSCConfig instance ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	addr, err := eurusInternalConfig.GetErc20SmartContractAddrByAssetName(&bind.CallOpts{}, symbol)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get erc 20 ", symbol, err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	erc20, err := contract.NewTestERC20(addr, mainnetClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot create erc20 instance ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	faucetConfigs, err := getFaucetConfigFromConfigServer(server)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to GetFaucetConfig: " + err.Error())
		return response.CreateErrorResponse(req, foundation.MethodNotFound, err.Error())
	}
	for _, e := range faucetConfigs {
		if e.Key == symbol {

			var tx *types.Transaction
			if symbol == "ETH" {
				_, tx, err = mainnetClient.TransferETH(secret.FaucetHdWalletPrivateKey, userFromDB.MainnetWalletAddress, big.NewInt(int64(e.Amount)))
				if err != nil {
					log.GetLogger(log.Name.Root).Error("Unable to transfer eth: " + err.Error())
					return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
				}
			} else {

				transOpt, err := mainnetClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, mainnetClient.ChainID)
				if err != nil {
					log.GetLogger(log.Name.Root).Error("cannot get transOpt ", err.Error())
					return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
				}
				if server.Config.FaucetGasLimit > 0 {
					transOpt.GasLimit = uint64(server.Config.FaucetGasLimit)
				}

				tx, err = erc20.Mint(transOpt, common.HexToAddress(userFromDB.MainnetWalletAddress), big.NewInt(int64(e.Amount)))
				if err != nil {
					log.GetLogger(log.Name.Root).Error("Unable to mint: " + err.Error())
					return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
				}
			}

			go func() {
				err := DbAddUserFaucetRecord(server.DefaultDatabase, symbol, userId, tx.Hash().Hex(), 1)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Unable to insert faucet record to DB. Tx: ", tx.Hash().Hex(), " Error: ", err)
					return
				}
				receipt, err := mainnetClient.QueryEthReceiptWithSetting(tx, 5, 10)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Unable to get receipt: " + err.Error())
					err := DbAddUserFaucetRecord(server.DefaultDatabase, symbol, userId, tx.Hash().Hex(), -1)
					if err != nil {
						log.GetLogger(log.Name.Root).Errorln("Unable to insert faucet record to DB for query receipt error. Tx: ", tx.Hash().Hex(), " Error: ", err)
					}
					return
				}

				if receipt.Status != 1 {
					err := DbAddUserFaucetRecord(server.DefaultDatabase, symbol, userId, tx.Hash().Hex(), -1)
					if err != nil {
						log.GetLogger(log.Name.Root).Errorln("Unable to insert faucet record to DB for receipt status = -1. Tx: ", tx.Hash().Hex(), " Error: ", err)
					}
					log.GetLogger(log.Name.Root).Errorln("faucet failed: ", receipt.TxHash.Hex())
					return
				}
				err = DbAddUserFaucetRecord(server.DefaultDatabase, symbol, userId, tx.Hash().Hex(), 0)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("add faucet record fail: ", receipt.TxHash.Hex(), " Error: ", err)
				}
				//TODO: TEST purpose
				//afterAmount,_:=erc20.BalanceOf(&bind.CallOpts{},common.HexToAddress("0x0b601ce20491145abcb994f61bc04eb75c18e851"))
				//fmt.Println("after amount: ",afterAmount)
			}()

			faucetResponse := new(FaucetResponse)
			faucetResponse.TxHash = tx.Hash().Hex()
			faucetResponse.Status = 1
			return response.CreateSuccessResponse(req, faucetResponse)
		}
	}
	log.GetLogger(log.Name.Root).Errorln("cannot find currency symbol: " + symbol)
	return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot found "+symbol)

}

func GetFaucetConfig(server *UserServer, req *request.RequestBase) *response.ResponseBase {
	faucetConfig, err := getFaucetConfigFromConfigServer(server)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get faucet config: " + err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get faucet config")
	}

	return response.CreateSuccessResponse(req, faucetConfig)
}

func Register(server *UserServer, req *RegistrationRequest) *response.ResponseBase {
	req.Email = strings.ToLower(req.Email)

	// Some special characters may not be allowed, currently block any plus sign inside email address
	if strings.ContainsAny(req.Email, "+") {
		log.GetLogger(log.Name.Root).Errorln("Invalid character(s) in email addess: ", req.Email)
		return response.CreateErrorResponse(req, foundation.BadRequest, "Email address contains invalid character(s)")
	}

	isExist, userObj, err := DbCheckEmailExist(server.DefaultDatabase, req.Email)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot check email existence : ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot check email existence")
	}

	// Check if the email is not yet verified
	if isExist && userObj.Status != UserStatusVerifiedNotSetPaymentAddress && userObj.Status != UserStatusNotVerify {
		log.GetLogger(log.Name.Root).Errorln("email already register : ", req.Email)
		return response.CreateErrorResponse(req, foundation.BadRequest, "email already register")
	}

	verified, err := RegisterVerifySign(*req)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify signature for email :  ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot verify signature")
	}
	if !verified {
		log.GetLogger(log.Name.Root).Errorln("invalid signature for email : ", req.Email)
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "invalid signature")
	}

	req.LoginAddress = ethereum.ToLowerAddressString(req.LoginAddress)
	if !isExist {

		userObj, err = DbAddNewUser(req.LoginAddress, server.DefaultDatabase, true, req.Email, false)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("cannot add user to db :  ", err.Error(), " email: ", req.Email)
			return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot add user to db")
		}
		log.GetLogger(log.Name.Root).Infoln("New User Add to DB : ", userObj.Id, " email: ", req.Email)
	} else {
		userObj, err = DbUpdateNewCentralizedUser(req.LoginAddress, server.DefaultDatabase, req.Email)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("cannot update user to db :  ", err.Error(), " email: ", req.Email)
			return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot add user to db")
		}
		log.GetLogger(log.Name.Root).Infoln("Update User Add to DB : ", userObj.Id, " email: ", req.Email)
	}

	verification, serverErr := DbAddNewVerification(userObj.Id, server.DefaultDatabase, server.Config.VerificationDuration, VerificationRegistration)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot add verification to db for user :  ", serverErr.Error(), " user Id: ", userObj.Id)
		return response.CreateErrorResponse(req, serverErr.GetReturnCode(), serverErr.Error())
	}
	emailTemplate, err := GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get email template : ", err.Error(), " user Id: ", userObj.Id)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	go SendEmail(server.Config, userObj.Email, "Eurus Verification Code", verification.Code, emailTemplate)
	userRes := new(RegistrationResponse)
	userRes.UserId = userObj.Id

	if secret.Tag == "dev" || secret.Tag == "default" {
		userRes.Code = verification.Code
	} else {
		verification.Code = ""
	}

	return response.CreateSuccessResponse(req, userRes)
}

func EmailVerification(server *UserServer, req *VerificationRequest) *response.ResponseBase {
	if req.Email == "" || req.Code == "" || req.DeviceId == "" || req.PublicKey == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Missing argument")
	}
	req.Email = strings.ToLower(req.Email)
	user, err := DbGetUserByEmail(server.SlaveDatabase, req.Email)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get user :  ", err.Error(), " email: ", req.Email)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user")
	}

	if user.Status == UserStatusNormal {
		log.GetLogger(log.Name.Root).Errorln("user already verified :  ", req.Email)
		return response.CreateErrorResponse(req, foundation.BadRequest, "user already verified")
	}

	if user.Status != UserStatusNotVerify && user.Status != UserStatusVerifiedNotSetPaymentAddress {
		log.GetLogger(log.Name.Root).Errorln("user id : ", user.Id, " status is not pending :  ", user.Status)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "user account status problem")
	}

	verified, err := DbVerifyCode(server.DefaultDatabase, user.Id, req.Code, int(VerificationRegistration))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify user : ", req.Email, err.Error(), " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if !verified {
		log.GetLogger(log.Name.Root).Errorln("invalid verification code for user : ", req.Email, " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid verification code")
	}

	// Here customer type must be user
	err = DbAddOrUpdateUserDevice(server.DefaultDatabase, user.Id, int16(CustomerUser), req.DeviceId, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot store user device information : ", req.Email, " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	mnemonic := crypto.Generate128BitsEntropyMnemonicPhrase()
	//address := crypto.GenerateWalletAddressFromMnemonic(mnemonic)

	// Encrypt mnemonic by RSA AES key, store it so that server side has the ability to recover it
	encrypted, err := encryptMnemonicWithRSAAesKey(mnemonic, server.Config.UserMnenomicAesKey)
	log.GetLogger(log.Name.Root).Debugln("Encrypted: ", encrypted, " user Id: ", user.Id)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot encrypt : ", err.Error(), " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot encrypt")
	}

	// Mnemonic encrypted by user-provided public key, just send back to user and no need to store
	encryptedForUser, err := encryptMnemonicWithPublicKey(mnemonic, req.PublicKey)
	log.GetLogger(log.Name.Root).Debugln("Encrypted for user: ", encryptedForUser, " user Id: ", user.Id)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot encrypt for user : ", err.Error(), " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot encrypt")
	}

	mainnetAddressReq := sign_api.NewGetCentralizedUserMainnetAddressRequest()
	mainnetAddressReq.UserId = user.Id

	mainnetAddressRes := new(sign_api.GetCentralizedUserMainnetAddressFullResponse)
	reqRes := api.NewRequestResponse(mainnetAddressReq, mainnetAddressRes)
	url, err := url.Parse(server.Config.SignServerUrl)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Sign server URL error: ", err)
		return response.CreateErrorResponse(req, foundation.NetworkError, "Internal config error")
	}
	_, err = api.SendApiRequest(*url, reqRes, server.AuthClient)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Cannot query sign server to get mainnet address: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.NetworkError, "Internal network error")
	}
	if reqRes.Res.GetReturnCode() != int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("Sign server response error code for GetCentralizedUserMainnetAddressFullResponse: ",
			reqRes.Res.GetMessage(), " code: ", reqRes.Res.GetReturnCode())

		return response.CreateErrorResponse(req, foundation.InternalServerError, reqRes.Res.GetMessage())
	}

	mainnetAddress := mainnetAddressRes.Data

	err = DbUpdateUserRegisterSuccessful(server.DefaultDatabase, user.Id, encrypted, strings.ToLower(mainnetAddress))

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot update user successful verify : ", err.Error(), " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot update user successful verify")
	}

	verificationResponse := new(VerificationResponse)
	verificationResponse.UserId = user.Id
	verificationResponse.Email = user.Email
	verificationResponse.Mnemonic = encryptedForUser

	token, err := server.AuthClient.GenerateLoginToken(fmt.Sprintf("{\"loginAddress\":\"%v\",\"userId\":%v}", user.LoginAddress, user.Id))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get token : ", err.Error(), " user Id: ", user.Id)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get token")
	}

	verificationResponse.Token = token.GetToken()
	verificationResponse.ExpiredTime = token.GetExpiredTime()

	return response.CreateSuccessResponse(req, verificationResponse)
}

func SetupPaymentWallet(server *UserServer, req *SetupPaymentWalletRequest, remoteAddr string) *response.ResponseBase {

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	dBUser, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user from db:  ", err.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if dBUser.Status == UserStatusNotVerify {
		log.GetLogger(log.Name.Root).Error("user not verify : ", dBUser.Id)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "user not verify")
	}

	if dBUser.Status != UserStatusVerifiedNotSetPaymentAddress {
		log.GetLogger(log.Name.Root).Error("user not eligible for payment wallet setup : ", dBUser.Id)
		return response.CreateErrorResponse(req, foundation.BadRequest, "user not eligible for payment wallet setup")
	}

	ownerWalletAddr := ethereum.ToLowerAddressString(req.UserWalletOwnerAddress)

	err = DbUpdateUserOwnerWalletAddress(server.DefaultDatabase, userID, ownerWalletAddr)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot update user owner wallet address :  ", err.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	var userWalletProxyAddr common.Address
	var deployUserWalletErr error
	for i := 0; i < server.Config.GetRetryCount(); i++ {
		userWalletProxyAddr, deployUserWalletErr = DeployUserWallet(server, userID)
		if deployUserWalletErr == nil {
			break
		} else {
			log.GetLogger(log.Name.Root).Error("deploy user wallet error: ", deployUserWalletErr.Error(), " retry: ", i, " user Id: ", userID)
		}
	}
	if deployUserWalletErr != nil || bytes.Equal(common.Address{}.Bytes(), userWalletProxyAddr.Bytes()) {
		log.GetLogger(log.Name.Root).Error("cannot deploy user wallet : ", deployUserWalletErr, " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.InternalServerError, deployUserWalletErr.Error())
	}

	dBUser.WalletAddress = userWalletProxyAddr.Hex()
	err = AddUserAddressToWalletSC(server, dBUser, false, false)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot add user address to wallet address map user id: ", userID)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot add user address to wallet sc")
	}

	var setUserWalletOwnerError error

	for i := 0; i < server.Config.GetRetryCount(); i++ {
		setUserWalletOwnerError = SetUserWalletOwner(server, req.UserWalletOwnerAddress, userWalletProxyAddr, userID)
		if setUserWalletOwnerError != nil {
			log.GetLogger(log.Name.Root).Error("cannot set user wallet owner : ", setUserWalletOwnerError.Error(), " user Id: ", userID)
		} else {
			break
		}
	}

	if setUserWalletOwnerError != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet owner : ", err.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, err.Error())
	}

	//setUserWalletInternalSmartContractConfig
	var setUserWalletInternalSCConfigError error
	var setUserWalletInternalSCConfigTx *types.Transaction
	for i := 0; i < server.Config.GetRetryCount(); i++ {
		setUserWalletInternalSCConfigTx, setUserWalletInternalSCConfigError = SetUserWalletInternalSmartContractConfig(server, userWalletProxyAddr, common.HexToAddress(server.Config.InternalSCConfigAddress))
		if setUserWalletInternalSCConfigError != nil {
			log.GetLogger(log.Name.Root).Error("cannot set user wallet internal SC Config : ", setUserWalletInternalSCConfigError.Error(), " user Id: ", userID)
		} else {
			break
		}
	}
	if setUserWalletInternalSCConfigError != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet internal SC Config : ", setUserWalletInternalSCConfigError.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, setUserWalletInternalSCConfigError.Error())
	}

	log.GetLogger(log.Name.Root).Infoln("set user wallet internal sc config transaction hash : ", setUserWalletInternalSCConfigTx.Hash(), " user Id: ", userID)

	err = DbUpdateUserWalletAddress(server.DefaultDatabase, userID, strings.ToLower(userWalletProxyAddr.Hex()))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot update wallet address to db  : ", err.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	token, err := server.AuthClient.GenerateLoginToken(fmt.Sprintf("{\"loginAddress\":\"%v\",\"userId\":%v}", dBUser.LoginAddress, dBUser.Id))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot generate user login token", " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	setupUserWalletResponse := new(SetupUserWalletResponse)
	setupUserWalletResponse.Token = token.GetToken()
	setupUserWalletResponse.WalletAddress = dBUser.WalletAddress
	setupUserWalletResponse.MainnetWalletAddress = dBUser.MainnetWalletAddress
	setupUserWalletResponse.IsMetamaskAddr = dBUser.IsMetamaskAddr

	err = TransferEUNToUser(server, dBUser.Id, dBUser.WalletAddress)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to transfer EUN to newly registered user. Wallet address: ", dBUser.WalletAddress, " Error: ", err, " user Id: ", userID)
	}

	go func() {
		loginData := NewLoginDataFromLoginBySignature(req.LoginLogDetail, userID, remoteAddr)
		_ = server.elasticLogger.InsertLog(loginData)
	}()

	return response.CreateSuccessResponse(req, setupUserWalletResponse)

}

func ResendVerificationEmail(server *UserServer, req *ResendVerificationEmailRequest) *response.ResponseBase {
	isExist, err := DbCheckUnVerifiedUserExistById(server.SlaveDatabase, req.UserId)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot user id existence : ", req.UserId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot check email existence")
	}

	if !isExist {
		log.GetLogger(log.Name.Root).Error("cannot found unverified user : ", req.UserId)
		return response.CreateErrorResponse(req, foundation.BadRequest, "unverified user not exist")
	}

	verification, serverErr := DbUpdateVerification(req.UserId, server.DefaultDatabase, 300, VerificationRegistration)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Error("cannot update verification : ", req.UserId)
		return response.CreateErrorResponse(req, serverErr.GetReturnCode(), serverErr.Error())
	}

	user, err := DbGetUserById(req.UserId, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user by user id : ", req.UserId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user")
	}

	emailTemplate, err := GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get email template : ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	var verCode string = verification.Code
	go SendEmail(server.Config, user.Email, "Eurus Verification Code", verCode, emailTemplate)

	if secret.Tag != "dev" && secret.Tag != "default" {
		verification.Code = ""
	}
	return response.CreateSuccessResponse(req, verification)

}

func RequestLoginRequestToken(server *UserServer, req *request.RequestBase) *response.ResponseBase {
	loginRequestTokenMap, err := DbInsertLoginRequestToken(server.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot request login request token")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot request login request token")
	}

	return response.CreateSuccessResponse(req, RequestLoginRequestTokenResponse{LoginRequestToken: loginRequestTokenMap.LoginRequestToken, ExpiredTime: *loginRequestTokenMap.ExpiredTime})

}

func RequestLoginTokenFromLoginRequestToken(server *UserServer, req *RequestLoginTokenRequest) *response.ResponseBase {
	isValid, err := DbCheckIfLoginRequestTokenValid(server.SlaveDatabase, req.LoginRequestToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot verify token ")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot verify token")
	}

	if !isValid {
		log.GetLogger(log.Name.Root).Error("invalid token ")
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, "invalid token")
	}

	token, _ := server.AuthClient.GenerateLoginToken(req.LoginToken.GetUserId())
	err = DbUpdateLoginRequestToken(server.DefaultDatabase, token.GetToken(), req.LoginRequestToken, req.LoginToken.GetUserId())
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get token ")
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get token")
	}
	return response.CreateSuccessResponse(req, RequestLoginTokenResponse{LoginToken: token.GetToken()})

}

func ChangePaymentPassword(server *UserServer, req *RequestChangePasswordRequest) *response.ResponseBase {

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	// verify new signature
	newSignIsValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.OwnerWalletAddress, req.Sign, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify new signature: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy new signature")
	}

	if !newSignIsValid {
		log.GetLogger(log.Name.Root).Errorln("invalid new signature", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid new signature")
	}

	//verify old signature
	oldSignIsValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.OldOwnerWalletAddress, req.OldSign, req.IsPersonalSign, req.OldPublicKey)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify old signature: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy old signature")
	}

	if !oldSignIsValid {
		log.GetLogger(log.Name.Root).Errorln("invalid old signature", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid old signature")
	}

	//check new address
	decompressedPubKey, err := crypto.DecompressPubKey(req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot decompress new pub key. Error: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot decompress new public key")
	}
	newAddr, err := crypto.PubKeyStringToAddress(decompressedPubKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("invalid new public key. Error: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "invalid new public key")
	}

	if ethereum.ToLowerAddressString(req.OwnerWalletAddress) != ethereum.ToLowerAddressString(newAddr.Hex()) {
		log.GetLogger(log.Name.Root).Errorln("new public key address not match", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Incorrect new login address")
	}
	//check old address
	decompressedOldPubKey, err := crypto.DecompressPubKey(req.OldPublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot decompress old pub key. Error: ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot decompress old public key")
	}
	oldAddr, err := crypto.PubKeyStringToAddress(decompressedOldPubKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("invalid old public key. Error: ", err.Error())
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "invalid old public key")
	}

	if ethereum.ToLowerAddressString(req.OldOwnerWalletAddress) != ethereum.ToLowerAddressString(oldAddr.Hex()) {
		log.GetLogger(log.Name.Root).Errorln("old public key address not match. Error: ", err.Error())
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Incorrect old login address")
	}

	dBUser, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get user :", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user from DB")
	}

	if ethereum.ToLowerAddressString(dBUser.OwnerWalletAddress) != ethereum.ToLowerAddressString(req.OldOwnerWalletAddress) {
		log.GetLogger(log.Name.Root).Errorln("Old owner wallet address not does match with user wallet owner wallet address. userId: ", dBUser.Id)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Old owner wallet address not does match with user wallet owner wallet address")
	}

	err = SetUserWalletOwner(server, req.OwnerWalletAddress, common.HexToAddress(dBUser.WalletAddress), userID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln(err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	err = DbChangeUserOwnerWalletAddress(server.DefaultDatabase, userID, req.OwnerWalletAddress)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot update user owner wallet address :", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot update user owner wallet address")
	}

	return response.CreateSuccessResponse(req, nil)

}

func ChangeLoginPassword(server *UserServer, req *RequestChangeLoginPasswordRequest) *response.ResponseBase {

	// verify new signature
	isValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.LoginAddress, req.Sign, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Infoln("cannot verify signature: ", err.Error())
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy signature")
	}
	if !isValid {
		log.GetLogger(log.Name.Root).Infoln("invalid signature")
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid signature")
	}

	//verify old signature
	oldSignIsValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.OldLoginAddress, req.OldSign, req.IsPersonalSign, req.OldPublicKey)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify old signature: ", err.Error())
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy old signature")
	}

	if !oldSignIsValid {
		log.GetLogger(log.Name.Root).Errorln("invalid old signature")
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid old signature")
	}

	//check new address

	decompressedPubKey, err := crypto.DecompressPubKey(req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot decompress new pub key")
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot decompress new public key")
	}

	addr, err := crypto.PubKeyStringToAddress(decompressedPubKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Infoln("invalid public key")
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid public key")
	}

	if ethereum.ToLowerAddressString(req.LoginAddress) != ethereum.ToLowerAddressString(addr.Hex()) {
		log.GetLogger(log.Name.Root).Infoln("public key address not match")
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}
	//check old address
	oldDecompressedPublicKey, err := crypto.DecompressPubKey(req.OldPublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot decompress old pub key")
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot decompress old public key")
	}
	oldAddr, err := crypto.PubKeyStringToAddress(oldDecompressedPublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("invalid old public key")
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "invalid old public key")
	}

	if ethereum.ToLowerAddressString(req.OldLoginAddress) != ethereum.ToLowerAddressString(oldAddr.Hex()) {
		log.GetLogger(log.Name.Root).Errorln("old public key address not match")
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Incorrect old login address")
	}

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}
	dBUser, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user by id: ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user by id")
	}

	if ethereum.ToLowerAddressString(dBUser.LoginAddress) != ethereum.ToLowerAddressString(req.OldLoginAddress) {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "old login address not match")
	}

	err = DbUpdateUserLoginAddress(server.DefaultDatabase, userID, req.LoginAddress)
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot update login address")
	}
	return response.CreateSuccessResponse(req, nil)
}

func RequestPaymentLoginToken(server *UserServer, req *request.RequestBase) *response.ResponseBase {
	token, err := server.AuthClient.RequestNonRefreshableLoginToken(req.LoginToken.GetUserId(), 300, int16(auth.NonRefreshableToken))
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get payment login token")
	}

	return response.CreateSuccessResponse(req, RequestLoginTokenResponse{LoginToken: token.GetToken()})
}

func NewLoginDataFromLoginBySignature(reqObj LoginLogDetail, userID uint64, remoteAddr string) *elasticLoginData {
	userIP := remoteAddr[:strings.IndexByte(remoteAddr, ':')]
	loginData := newElasticLoginData()
	loginData.UserId = userID
	loginData.Ip = userIP
	loginData.AppVersion = reqObj.AppVersion
	loginData.Os = reqObj.Os
	loginData.RegistrationSource = reqObj.RegistrationSource
	return loginData
}

func RequestGetUserStorage(server *UserServer, req *RequestUserStorage) *response.ResponseBase {
	var err error
	req.UserId, err = UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	getData, err := DbGetUserStorageByID(req, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get User storage")
	}

	return response.CreateSuccessResponse(req, GetUserStorageResponse{UserId: req.UserId, Sequence: int(getData.Sequence), Platform: req.Platform, Storage: getData.Storage})
}

func RequestSetUserStorage(server *UserServer, req *RequestUserStorage) *response.ResponseBase {
	var err error
	req.UserId, err = UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	updateData, err := DbUpdateUserStorage(req, server.DefaultDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get User storage")
	}
	return response.CreateSuccessResponse(req, UserStorageSequenceResponse{Sequence: int(updateData.Sequence)})
}

func UnmarshalUserIdFromLoginToken(req *request.RequestBase) (uint64, error) {
	type user struct {
		UserId       uint64 `json:"userId"`
		LoginAddress string `json:"loginAddress"`
	}

	userInstance := new(user)
	userId := req.LoginToken.GetUserId()
	err := json.Unmarshal([]byte(userId), userInstance)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot unmarshal : ", err.Error())
		return 0, err
	}
	return userInstance.UserId, nil
}

func ForgetLoginPassword(server *UserServer, req *ForgetLoginPasswordRequest) *response.ResponseBase {
	emailTemplate, err := GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get email template : ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	user, err := DbGetUserByEmail(server.SlaveDatabase, req.Email)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Cannot get user by email: ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get User by email")
	}

	checkRes := checkUserStatus(req, user)
	if checkRes != nil {
		return checkRes
	}

	var verification *Verification
	verification, serverErr := DbAddNewVerification(user.Id, server.DefaultDatabase, 300, VerificationForgetLoginPassword)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Error("cannot add/update verification to db : ", serverErr.Error())
		return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Error())
	}

	emailTemplate, err = GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get email template : ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	verCode := verification.Code
	go SendEmail(server.Config, user.Email, "Eurus Verification Code", verCode, emailTemplate)

	if secret.Tag != "dev" && secret.Tag != "default" {
		verification.Code = ""
	}

	return response.CreateSuccessResponse(req, verification)
}

func VerifyForgetLoginPasswordCode(server *UserServer, req *VerifyForgetLoginPasswordRequest) *response.ResponseBase {
	user, err := DbGetUserByEmail(server.SlaveDatabase, req.Email)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user by email : ", req.Email, err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "email not found")
	}

	verified, err := DbVerifyCode(server.DefaultDatabase, user.Id, req.Code, int(VerificationForgetLoginPassword))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot verify user : ", req.Email, err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if !verified {
		log.GetLogger(log.Name.Root).Error("invalid verification code for user : ", req.Email)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid verification code")
	}

	token, getAuthTokenErr := server.AuthClient.RequestNonRefreshableLoginToken("{\"userId\":"+strconv.Itoa(int(user.Id))+",\"loginAddress\":\"\"}", 300, int16(auth.ResetLoginPasswordToken))
	if getAuthTokenErr != nil {
		log.GetLogger(log.Name.Root).Error("cannot get non-refreshable login token : ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get non-refreshable login token")
	}
	verifyForgetLoginPasswordResponse := new(VerifyForgetPasswordResponse)
	verifyForgetLoginPasswordResponse.Token = token.GetToken()
	return response.CreateSuccessResponse(req, verifyForgetLoginPasswordResponse)

}

func ResetLoginPassword(server *UserServer, req *ResetLoginPasswordReqeust) *response.ResponseBase {
	if req.LoginToken.GetTokenType() != int16(auth.ResetLoginPasswordToken) {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "token type unmatch")
	}

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}
	// verify signature
	isValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.LoginAddress, req.Sign, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Infoln("cannot verify signature: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy signature")
	}

	if !isValid {
		log.GetLogger(log.Name.Root).Infoln("invalid signature", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid signature")
	}

	decompressedPubKey, err := crypto.DecompressPubKey(req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to decompress public key: ", err, " userId: ", userID)
	}
	// check public key
	addr, err := crypto.PubKeyStringToAddress(decompressedPubKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("invalid public key", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "invalid public key")
	}

	if ethereum.ToLowerAddressString(req.LoginAddress) != ethereum.ToLowerAddressString(addr.Hex()) {
		log.GetLogger(log.Name.Root).Errorln("public key address not match", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Incorrect login address")
	}

	err = DbUpdateUserLoginAddress(server.DefaultDatabase, userID, req.LoginAddress)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("DbUpdateUserLoginAddress error: ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot update login address")
	}
	server.AuthClient.RevokeLoginToken(req.LoginToken.GetToken())
	return response.CreateSuccessResponse(req, nil)
}

func ForgetPaymentPassword(server *UserServer, req *request.RequestBase) *response.ResponseBase {
	userID, err := UnmarshalUserIdFromLoginToken(req)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}

	user, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error(err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user by userid")
	}
	checkRes := checkUserStatus(req, user)
	if checkRes != nil {
		return checkRes
	}
	var verification *Verification

	verification, serverErr := DbAddNewVerification(user.Id, server.DefaultDatabase, 300, VerificationForgetPaymentPassword)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Error("cannot add verification to db : ", serverErr.Error())
		return response.CreateErrorResponse(req, serverErr.GetReturnCode(), serverErr.Error())
	}

	emailTemplate, err := GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get email template : ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	verCode := verification.Code
	go SendEmail(server.Config, user.Email, "Eurus Verification Code", verCode, emailTemplate)

	if secret.Tag != "dev" && secret.Tag != "default" {
		verification.Code = ""
	}
	return response.CreateSuccessResponse(req, verification)
}

// This method will verify if the user already apply forget password request from user/forgetPaymentPassword
// The user must user/forgetPaymentPassword first
// It fetchs record from verifications table where type = 2 and userId=token.userId
// If there is record, this function will select user public key from user_devices
// It will then generate new payment mnemonic phase, by decrypt original mnemonic phase in DB
// Using public key to encrypt it with AES key, and updated the new mnemonic phase to users table
// Use public key to encrypt the new mnemonic phase in RA format, then respond
func VerifyForgetPaymentPasswordCode(server *UserServer, req *VerifyForgetPaymentPasswordRequest) *response.ResponseBase {

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}

	verified, err := DbVerifyCode(server.DefaultDatabase, userID, req.Code, int(VerificationForgetPaymentPassword))
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot verify user : ", userID, err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if !verified {
		log.GetLogger(log.Name.Root).Error("invalid verification code for user : ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid verification code")
	}

	dbUser, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	checkRes := checkUserStatus(req, dbUser)
	if checkRes != nil {
		return checkRes
	}

	token, getAuthTokenErr := server.AuthClient.RequestNonRefreshableLoginToken("{\"userId\":"+strconv.FormatUint(userID, 10)+",\"loginAddress\":\"\"}", 300, int16(auth.ResetPaymentPasswordToken))
	if getAuthTokenErr != nil {
		log.GetLogger(log.Name.Root).Error("cannot get non-refreshable login token : ", err.Error())
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot get non-refreshable login token")
	}

	if req.DeviceId == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Missing device id")
	}

	// Issue EU-947:
	// 1. Fetch public key from table user_devices
	publicKey, err := DbGetUserDevicePublicKey(server.SlaveDatabase, userID, int16(CustomerUser), req.DeviceId)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to get user public key:"+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	// If public key cannot be found, just return error
	if publicKey == "" {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Public key does not exist for this device id", "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Public key does not exist for this device id")
	}

	// 2. Generate new payment mnemonic phase
	plainTextMnemonic := crypto.Generate128BitsEntropyMnemonicPhrase()

	// 3. Encrypt new generated mnemonic by RSA AES key
	encryptedMnemonicAES, err := encryptMnemonicWithRSAAesKey(plainTextMnemonic, server.Config.UserMnenomicAesKey)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot encrypt : ", err.Error(), " user Id: ", userID)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "cannot encrypt")
	}
	// 4. Update the encrypted mnemonic to User table
	err = DbUpdateUserMnemonicPhase(server.DefaultDatabase, userID, encryptedMnemonicAES)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot update verification : ", err.Error())
		return response.CreateErrorResponse(req, foundation.RequestTooFrequenct, err.Error())
	}
	//5. Encrypt new generated mnemonic phase by public key in "RA" format
	encryptedMnemonicAESByRAAndPublicKey, err := encryptMnemonicWithPublicKey(plainTextMnemonic, publicKey)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to encrypt user mnemonic:"+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	//6. Send the response to client side
	verifyForgetPaymentPasswordResponse := new(VerifyForgetPasswordResponse)
	verifyForgetPaymentPasswordResponse.Token = token.GetToken()
	verifyForgetPaymentPasswordResponse.Mnemonic = encryptedMnemonicAESByRAAndPublicKey
	return response.CreateSuccessResponse(req, verifyForgetPaymentPasswordResponse)

}

func ResetPaymentPassword(server *UserServer, req *ResetPaymentPasswordReqeust) *response.ResponseBase {
	if req.LoginToken.GetTokenType() != int16(auth.ResetPaymentPasswordToken) {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "token type unmatch")
	}
	// verify signature
	isValid, err := VerifySignature(req.DeviceId, req.Timestamp, req.OwnerWalletAddress, req.Sign, req.IsPersonalSign, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Infoln("cannot verify signature: ", err.Error())
		return response.CreateErrorResponse(req, foundation.RequestMalformat, "cannot verifiy signature")
	}

	if !isValid {
		log.GetLogger(log.Name.Root).Infoln("invalid signature")
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid signature")
	}

	userID, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "Incorrect wallet address")
	}

	dbUser, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	checkRes := checkUserStatus(req, dbUser)
	if checkRes != nil {
		return checkRes
	}

	// check public key
	decompressedPubKey, err := crypto.DecompressPubKey(req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to decompress public key: ", err, " userId: ", userID)
	}
	addr, err := crypto.PubKeyStringToAddress(decompressedPubKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("invalid public key", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.RequestParamsValidationError, "invalid public key")
	}

	if ethereum.ToLowerAddressString(req.OwnerWalletAddress) != ethereum.ToLowerAddressString(addr.Hex()) {
		log.GetLogger(log.Name.Root).Errorln("public key address not match", " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Incorrect login address")
	}

	user, err := DbGetUserById(userID, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot get user :  ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user")
	}

	var setUserWalletOwnerError error
	for i := 0; i < server.Config.GetRetryCount(); i++ {
		setUserWalletOwnerError = SetUserWalletOwner(server, req.OwnerWalletAddress, common.HexToAddress(user.WalletAddress), userID)
		if setUserWalletOwnerError != nil {
			log.GetLogger(log.Name.Root).Error("cannot set user wallet owner : ", setUserWalletOwnerError.Error(), " userId: ", userID)
		} else {
			break
		}
	}

	if setUserWalletOwnerError != nil {
		log.GetLogger(log.Name.Root).Error("cannot set user wallet owner : ", setUserWalletOwnerError.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, setUserWalletOwnerError.Error())
	}

	err = DbChangeUserOwnerWalletAddress(server.DefaultDatabase, userID, req.OwnerWalletAddress)

	if err != nil {
		log.GetLogger(log.Name.Root).Error("cannot update user onwer wallet address :  ", err.Error(), " userId: ", userID)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	server.AuthClient.RevokeLoginToken(req.LoginToken.GetToken())
	return response.CreateSuccessResponse(req, nil)

}

func FindEmailWalletAddress(server *UserServer, req *EmailWalletAddressRequest) *response.ResponseBase {

	if req.WalletAddress == "" && req.Email == "" {
		return response.CreateErrorResponse(req, foundation.BadRequest, "Invalid request body")
	}

	if req.WalletAddress != "" {
		req.WalletAddress = ethereum.ToLowerAddressString(req.WalletAddress)
		user, err := DbGetUserByWalletAddress(req.WalletAddress, server.SlaveDatabase)

		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get email by wallet address : ", err.Error())
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}
		emailWalletAddressResponse := new(EmailWalletAddressResponse)
		emailWalletAddressResponse.Email = user.Email
		emailWalletAddressResponse.WalletAddress = req.WalletAddress
		if user.IsMetamaskAddr {
			emailWalletAddressResponse.UserType = int(DecentralizedUser)
		} else {
			emailWalletAddressResponse.UserType = int(CentralizedUser)
		}

		return response.CreateSuccessResponse(req, emailWalletAddressResponse)
	} else {
		user, err := DbGetUserByEmail(server.SlaveDatabase, req.Email)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("cannot get wallet address by email : ", err.Error())
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}

		checkRes := checkUserStatus(req, user)
		if checkRes != nil {
			return checkRes
		}
		emailWalletAddressResponse := new(EmailWalletAddressResponse)
		emailWalletAddressResponse.Email = req.Email
		emailWalletAddressResponse.WalletAddress = user.WalletAddress

		if user.IsMetamaskAddr {
			emailWalletAddressResponse.UserType = int(DecentralizedUser)
		} else {
			emailWalletAddressResponse.UserType = int(CentralizedUser)
		}

		return response.CreateSuccessResponse(req, emailWalletAddressResponse)

	}

}

func RegisterDevice(server *UserServer, req *RegisterDeviceRequest) *response.ResponseBase {
	userId, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	user, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get user :", err.Error(), " userId: ", userId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user from DB")
	}

	checkRes := checkUserStatus(req, user)
	if checkRes != nil {
		return checkRes
	}

	// This feature is only for centralized user
	if user.IsMetamaskAddr {
		log.GetLogger(log.Name.Root).Errorln("Decentralized user attempt to register device: userId: ", userId)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Not applicable for decentralized user")
	}

	// Add a new verification code
	// User should provide the new device id, public key and this code to show he or she really want to register a new device
	verification, serverErr := DbAddNewVerification(userId, server.DefaultDatabase, server.Config.VerificationDuration, VerificationRegisterDevice)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot add verification to db for user :  ", serverErr.Error(), " user Id: ", userId)
		return response.CreateErrorResponse(req, serverErr.GetReturnCode(), serverErr.Error())
	}

	emailTemplate, err := GetEmailTemplate(&server.ServerBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get email template : ", err.Error(), " user Id: ", userId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	go SendEmail(server.Config, user.Email, "Eurus Verification Code", verification.Code, emailTemplate)

	userRes := new(RegisterDeviceResponse)

	// Response in dev env can include the verification code directly
	if secret.Tag == "dev" || secret.Tag == "default" {
		userRes.Code = verification.Code
	}

	return response.CreateSuccessResponse(req, userRes)
}

func VerifyDevice(server *UserServer, req *VerifyDeviceRequest) *response.ResponseBase {
	if req.Code == "" || req.DeviceId == "" || req.PublicKey == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Missing argument")
	}

	userId, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid login token format")
	}

	user, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get user :", err.Error(), " userId: ", userId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "cannot get user from DB")
	}
	checkRes := checkUserStatus(req, user)
	if checkRes != nil {
		return checkRes
	}

	// Again registering new device is only for centralized user
	if user.IsMetamaskAddr {
		log.GetLogger(log.Name.Root).Errorln("Decentralized user attempt to verify for new device: userId: ", userId)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Not applicable for decentralized user")
	}

	verified, err := DbVerifyCode(server.DefaultDatabase, userId, req.Code, int(VerificationRegisterDevice))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot verify user : ", err.Error(), " user Id: ", userId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if !verified {
		log.GetLogger(log.Name.Root).Errorln("invalid verification code for user: user Id: ", userId)
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "invalid verification code")
	}

	err = DbAddOrUpdateUserDevice(server.DefaultDatabase, userId, int16(CustomerUser), req.DeviceId, req.PublicKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot store user device information: user Id: ", userId)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	plainTextMnemonic, err := decryptMnemonic(user.Mnemonic, server.Config.UserMnenomicAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to decrypt payment mnemonic, generate a dummy one for user: ", userId, " error: ", err.Error())
		plainTextMnemonic = crypto.Generate128BitsEntropyMnemonicPhrase()
	}

	encryptedForUser, err := encryptMnemonicWithPublicKey(plainTextMnemonic, req.PublicKey)
	if err != nil {
		reqStr, _ := json.Marshal(req)
		log.GetLogger(log.Name.Root).Error("Unable to encrypt user mnemonic:"+err.Error(), "\nRequest Params: "+string(reqStr))
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	userRes := new(VerifyDeviceResponse)
	userRes.Mnemonic = encryptedForUser
	return response.CreateSuccessResponse(req, userRes)
}

func SignTransaction(server *UserServer, req *SignTransactionRequest) *response.ResponseBase {
	userId, err := GetUserIdFromLoginToken(req.LoginToken)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
	}

	userObj, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	}

	if userObj.Id == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	}

	if userObj.Status != UserStatusNormal {
		return response.CreateErrorResponse(req, foundation.UserStatusForbidden, foundation.UserStatusForbidden.String())
	}

	signReq := sign_api.NewSignUserWalletTransactionRequest()
	signReq.GasPrice = req.GasPrice
	signReq.InputFunction = req.InputFunction
	signReq.Value = req.Value
	signReq.Nonce = req.GetNonce()
	signReq.To = userObj.WalletAddress

	signRes := &sign_api.SignedUserWalletTransactionFullResponse{}

	reqRes := api.NewRequestResponse(signReq, signRes)
	_, err = server.SendApiRequest(server.Config.SignServerUrl, reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to sign user wallet transaction. Nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	if reqRes.Res.GetReturnCode() != int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("Unable to sign user wallet transaction. Return code: ", reqRes.Res.GetReturnCode(), " Nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.ServerReturnCode(reqRes.Res.GetReturnCode()), signRes.Message)
	}
	res := new(SignTransactionResponse)
	res.SignedTx = signRes.Data.SignedTx

	return response.CreateSuccessResponse(req, res)
}

func processTopUpPaymentWallet(server *UserServer, req *TopUpPaymentWalletRequest) *response.ResponseBase {

	userId, err := GetUserIdFromLoginToken(req.LoginToken)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, foundation.LoginTokenInvalid.String())
	}

	userObj, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	}

	if userObj.Id == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	}

	if userObj.Status != UserStatusNormal {
		return response.CreateErrorResponse(req, foundation.UserStatusForbidden, foundation.UserStatusForbidden.String())
	}

	userWallet, err := contract.NewUserWallet(common.HexToAddress(userObj.WalletAddress), server.EthClient.Client)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create user wallet smart contract instance: ", err, " userId: ", userId)
		return response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
	}

	if req.Signature == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Signature is a required field")
	}
	signByte, err := hex.DecodeString(req.Signature)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid signature format: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "Invalid signature format")
	}

	addrType, err := abi.NewType("address", "address", nil)
	arg := abi.Argument{
		Name: "walletOwnerAddress",
		Type: addrType,
	}

	argList := abi.Arguments{arg}
	packed, err := argList.Pack(common.HexToAddress(userObj.OwnerWalletAddress))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to Pack argument when creating signature. Error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "Unable to generate signature")
	}

	hashData := go_ethereum_crypto.Keccak256(packed)
	if len(signByte) < 64 {
		log.GetLogger(log.Name.Root).Errorln("Invalid signature length nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid signature length")
	}
	var isAlteredRecoverId bool
	if signByte[64] >= 27 {
		isAlteredRecoverId = true
		signByte[64] = signByte[64] - 27
	}
	//Check signature
	pubKey, err := go_ethereum_crypto.Ecrecover(hashData, signByte)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to recover public key from signature. Error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "Unable to recover public key from signature")
	}

	if isAlteredRecoverId {
		signByte[64] = signByte[64] + 27
	}
	addr, err := crypto.PubKeyStringToAddress(hex.EncodeToString(pubKey))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid public key: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "Invalid public key")
	}
	if ethereum.ToLowerAddressString(addr.Hex()) != userObj.OwnerWalletAddress {
		log.GetLogger(log.Name.Root).Errorln("Signature is not sign by wallet owner. Address is: ", addr.Hex(), " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InvalidSignature, "Signature is not sign by wallet owner")
	}
	targetGas := big.NewInt(0)
	targetGas.SetUint64(req.TargetGasAmount)

	maxTopUpAmount, err := GetMaxTopUpGasAmount(server)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetMaxTopUpGasAmount failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	if maxTopUpAmount.Cmp(targetGas) < 0 {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Target gas amount exceed limit")
	}

	//Estimate gas needed , simulate send by invoker (sign server invoker private key)
	if len(server.Config.InvokerAddressList) == 0 {
		log.GetLogger(log.Name.Root).Errorln("No invoker address in config. Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, "No invoker address")
	}
	transOpt, err := server.EthClient.GetNewTransactorFromPrivateKey(server.ServerConfig.HdWalletPrivateKey, ethereum.EurusChainConfig.ChainID)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetNewTransactorFromPrivateKey: ", err, " userId: ", userId)
		return response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
	}
	//Just need to create a function call abipacked bytes
	transOpt.Signer = func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		//No need to sign at user server, sign tx will be carried out at sign server
		return tx, nil
	}
	transOpt.Nonce = nil //Dummy txn nonce
	transOpt.NoSend = true
	transOpt.From = common.HexToAddress(server.Config.InvokerAddressList[0]) //Simulate send by sign server invoker address

	estimatedTx, err := userWallet.TopUpPaymentWallet(transOpt, common.HexToAddress(userObj.OwnerWalletAddress), targetGas,
		signByte)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Estimate gas failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InsufficientAccountBalance, foundation.InsufficientAccountBalance.String())
	}

	userWalletAmount, err := server.EthClient.GetBalance(common.HexToAddress(userObj.WalletAddress))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query user wallet balance: ", err, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.EthereumError, "Unable to query user wallet balance: "+err.Error())
	}

	userOwnerWalletAmount, err := server.EthClient.GetBalance(common.HexToAddress(userObj.OwnerWalletAddress))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query owner user wallet balance: ", err, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.EthereumError, "Unable to query owner user wallet balance: "+err.Error())
	}

	gasPrice, err := server.EthClient.GetGasPrice()
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query gas price: ", err, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.EthereumError, "Unable to query gas price: "+err.Error())
	}

	targetEun := new(big.Int).Mul(targetGas, gasPrice)
	eunTopUp := new(big.Int).Sub(targetEun, userOwnerWalletAmount)

	if eunTopUp.Sign() <= 0 {
		log.GetLogger(log.Name.Root).Errorln("Payment wallet has the target gas amount already. Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.TargetGasAmountReach, foundation.TargetGasAmountReach.String())
	}

	estimatedGas := new(big.Int).SetUint64(estimatedTx.Gas())
	gasUsed := new(big.Int).Div(
		new(big.Int).Mul(
			new(big.Int).Mul(estimatedGas, big.NewInt(100)),
			big.NewInt(2)),
		big.NewInt(100))

	eunUsed := new(big.Int).Add(new(big.Int).Mul(gasUsed, gasPrice), eunTopUp)
	if userWalletAmount.Cmp(eunUsed) < 0 {
		log.GetLogger(log.Name.Root).Errorln("Account balance: ", userWalletAmount, "does not enough to pay gas fee: ", eunUsed, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InsufficientAccountBalance, foundation.InsufficientAccountBalance.String())
	}
	//Actually pack the transaction
	transOpt.GasLimit = 1000000
	tx, err := userWallet.TopUpPaymentWallet(transOpt, common.HexToAddress(userObj.OwnerWalletAddress), targetGas,
		signByte)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Create TopUpPaymentWallet transaction failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	signReq := sign_api.NewSignUserWalletTransactionRequest()
	signReq.GasPrice = 2400000000
	signReq.InputFunction = hex.EncodeToString(tx.Data())
	signReq.Value = "0"
	signReq.Nonce = req.GetNonce() //this is the API request nonce, not the ETH transaction nonce
	signReq.To = userObj.WalletAddress

	signRes := &sign_api.SignedUserWalletTransactionFullResponse{}

	reqRes := api.NewRequestResponse(signReq, signRes)
	_, err = server.SendApiRequest(server.Config.SignServerUrl, reqRes)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to sign user wallet transaction. Nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	if reqRes.Res.GetReturnCode() != int64(foundation.Success) {
		log.GetLogger(log.Name.Root).Errorln("Unable to sign user wallet transaction. Return code: ", reqRes.Res.GetReturnCode(), " Nonce: ", req.GetNonce(), " Error: ", err)
		return response.CreateErrorResponse(req, foundation.ServerReturnCode(reqRes.Res.GetReturnCode()), signRes.Message)
	}
	res := new(TopUpPaymentWalletResponse)
	res.Tx = signRes.Data.SignedTx
	res.EstimatedGasUsed = eunUsed.Uint64()

	return response.CreateSuccessResponse(req, res)
}

func GetUserIdFromLoginToken(loginToken auth_base.ILoginToken) (uint64, error) {
	str := loginToken.GetUserId()
	var tokenInfo map[string]interface{}
	err := json.Unmarshal([]byte(str), &tokenInfo)
	if err != nil {
		return 0, err
	}

	userId, ok := tokenInfo["userId"].(float64)
	if !ok {
		return 0, errors.New("User Id not found")
	}

	return uint64(userId), nil
}

func encryptMnemonicWithPublicKey(plainTextMnemonic string, base64PublicKey string) (string, error) {
	encrypted, err := crypto.EncryptRAFormat([]byte(plainTextMnemonic), base64PublicKey, go_crypto.SHA1)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func encryptMnemonicWithRSAAesKey(plainTextMnemonic string, mnemonicRSAAesKey string) (string, error) {
	decryptedRSAAesKey, err := secret.DecryptConfigValue(mnemonicRSAAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decrypt RSA Aes Key:" + err.Error())
		return "", err
	}
	mnemonicAesKeyBytes, err := base64.StdEncoding.DecodeString(decryptedRSAAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decode string:" + err.Error())
		return "", err
	}
	mnemonicAesKey := string(mnemonicAesKeyBytes)
	aesKey, err := base64.StdEncoding.DecodeString(mnemonicAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot get AES key for encryption : ", err.Error())
		return "", err
	}

	encrypted, err := crypto.EncryptAES([]byte(plainTextMnemonic), aesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("cannot encrypt : ", err.Error())
		return "", err
	}

	return encrypted, nil
}

func decryptMnemonic(base64EncryptedMnemonic string, mnemonicRSAAesKey string) (string, error) {
	decryptedRSAAesKey, err := secret.DecryptConfigValue(mnemonicRSAAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decrypt RSA Aes Key:" + err.Error())
		return "", err
	}
	mnemonicAesKeyBytes, err := base64.StdEncoding.DecodeString(decryptedRSAAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decode string:" + err.Error())
		return "", err
	}
	mnemonicAesKey := string(mnemonicAesKeyBytes)
	// aesKey is in raw format
	aesKey, err := base64.StdEncoding.DecodeString(mnemonicAesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get AES key for decryption:" + err.Error())
		return "", err
	}

	// Get the encrypted byte[] first
	rawEncryptedMenmonic, err := base64.StdEncoding.DecodeString(base64EncryptedMnemonic)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decode encrypted user mnemonic:" + err.Error())
		return "", err
	}

	// DecryptAES returns base64 encoded message, so need to decode again
	decrypted, err := crypto.DecryptAES(rawEncryptedMenmonic, aesKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decrypt user mnemonic:" + err.Error())
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(decrypted)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to decode user mnemonic:" + err.Error())
		return "", err
	}

	return string(decoded), nil
}

func GetEmailTemplate(server *server.ServerBase) (string, error) {

	var err error
	var res *conf_api.QueryConfigResponse = new(conf_api.QueryConfigResponse)
	for i := 0; i < server.ServerConfig.RetryCount; i++ {
		req := conf_api.NewQueryConfigRequest()
		req.Id = server.ServerConfig.ServiceId
		req.Key = "emailTemplate"

		reqRes := api.NewRequestResponse(req, res)
		reqRes, err = server.SendConfigApiRequest(reqRes)
		if err != nil {
			if i < server.ServerConfig.RetryCount {
				time.Sleep(time.Second * server.ServerConfig.GetRetryInterval())
			}
			continue
		}

		if res.ReturnCode != int64(foundation.Success) {
			err = errors.New("Server error code: " + strconv.FormatInt(res.ReturnCode, 10))
			if i < server.ServerConfig.RetryCount {
				time.Sleep(time.Second * server.ServerConfig.GetRetryInterval())
			}
		} else {
			break
		}
	}

	if err == nil {
		if len(res.Data.ConfigData) > 0 {
			configMap := res.Data.ConfigData[0]
			return configMap.Value, nil
		} else {
			err = nil
		}
	}

	return "", err
}

func getFaucetConfigFromConfigServer(server *UserServer) ([]*conf_api.FaucetConfig, error) {

	req := conf_api.NewQueryFaucetConfig()
	res := &conf_api.QueryFaucetConfigFullResponse{}
	reqRes := api.NewRequestResponse(req, res)
	_, err := server.SendConfigApiRequest(reqRes)

	if err != nil {
		return nil, err
	}
	if res.ReturnCode != int64(foundation.Success) {
		return nil, errors.New(res.GetMessage())
	}

	return res.Data, nil
}

func processGetWalletAddressByLoginAddress(db *database.ReadOnlyDatabase, req *GetWalletAddressRequest) *response.ResponseBase {
	addr := ethereum.ToLowerAddressString(req.LoginAddress)
	user, err := DbGetUserByLoginAddress(addr, db)

	if err != nil {
		if serverErr, ok := err.(*foundation.ServerError); ok {
			if serverErr.ReturnCode == foundation.UserNotFound {
				return response.CreateSuccessResponse(req, addr)
			} else {
				return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Error())
			}
		} else {
			return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
		}
	}
	res := checkUserStatus(req, user)
	if res != nil {
		return res
	}
	return response.CreateSuccessResponse(req, user.WalletAddress)
}

func processGetAssetAddressList(config *UserServerConfig, chainName string, req *request.RequestBase) *response.ResponseBase {
	switch chainName {
	case "bsc":
		if config.BSCAssetAddressList != "" {
			res := new(AssetAddressListResponse)
			res.Data = config.BSCAssetAddressList
			return response.CreateSuccessResponse(req, res)
		}
	}
	return response.CreateErrorResponse(req, foundation.InvalidArgument, "Chain not supported")
}

func checkUserStatus(req request.IRequest, user *User) *response.ResponseBase {

	if user.Id == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, foundation.UserNotFound.String())
	} else if user.Status != UserStatusNormal {
		var code foundation.ServerReturnCode

		switch user.Status {
		case UserStatusDeleted:
			code = foundation.UserNotFound
		case UserStatusSuspended:
			code = foundation.UserStatusForbidden
		default:
			code = foundation.UserNotFound
		}
		return response.CreateErrorResponse(req, code, code.String())
	}
	return nil
}
