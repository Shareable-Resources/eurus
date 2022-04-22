package user

import (
	"bytes"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"

	"github.com/ethereum/go-ethereum/common"
)

func ProcessRequestMerchantRefund(server *UserServer, req *RequestMerchantRefundRequest) *response.ResponseBase {

	userId, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid login token user Id", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.UserNotFound, err.Error())
	}
	user, err := DbGetUserById(userId, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to get user from DB: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if req.PurchaseTransHash != "" {
		tx, err := GetPurchaseTransaction(server, common.HexToHash(req.PurchaseTransHash))
		if err != nil {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Unable to get transaction: "+err.Error())
		}

		args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "EurusERC20", "purchase")
		if err != nil {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Transaction hash is not a purchase transaction")
		}

		receipt, err := server.EthClient.QueryEthReceipt(tx)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to get receipt. Trans: ", tx.Hash().String(), " nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Unable to get receipt: "+err.Error())
		}

		if receipt.Status == 0 {
			return response.CreateErrorResponse(req, foundation.EthereumError, "Purchase is unsuccessful")
		}

		dAppStockAbi := ethereum.DefaultABIDecoder.GetABI("DAppStockBase")
		purchasedEvent := dAppStockAbi.Events["purchasedEvent"]

		if err != nil {
			log.GetLogger(log.Name.Root).Error("Unable to insert merchant refund request: ", err, " nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}

		for _, logObj := range receipt.Logs {
			if bytes.Equal(logObj.Topics[0].Bytes(), purchasedEvent.ID.Bytes()) {
				buyer := ethereum.HashToAddress(&logObj.Topics[2])
				if !bytes.Equal(buyer.Bytes(), common.HexToAddress(user.WalletAddress).Bytes()) {
					return response.CreateErrorResponse(req, foundation.InvalidArgument, "Transaction is not make by user")
				}
				break
			}
		}

		merchantAddress, ok := args["dappAddress"].(common.Address)
		if !ok {
			log.GetLogger(log.Name.Root).Errorln("Unable to cast EurusERC20 purchase transaction dappAddress. Trans: ", tx.Hash().String(), " nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.InternalServerError, "Invalid purchase transaction")
		}
		merchantId, err := DbGetMerchantIdByWalletAddress(server.DefaultDatabase, merchantAddress)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("DbGetMerchantIdByWalletAddress error: ", err, " transaction hash: ", req.PurchaseTransHash, " nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}

		if merchantId != req.MerchantId {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Transaction is not for that merchant")
		}
	}

	requestId, err := DbInsertMerchantRefundRequest(server.DefaultDatabase, req, user)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to insert merchant refund request: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	resObj := new(RequestMerchantRefundResponse)
	resObj.RequestId = *requestId

	return response.CreateSuccessResponse(req, resObj)
}

func ProcessQueryMerchantRefundStatus(server *UserServer, req *QueryMerchantRefundStatusRequest) *response.ResponseBase {
	userId, err := UnmarshalUserIdFromLoginToken(&req.RequestBase)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, err.Error())
	}

	refundList, err := DbGetMerchantRefundStatus(server.SlaveDatabase, userId)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, refundList)
}
