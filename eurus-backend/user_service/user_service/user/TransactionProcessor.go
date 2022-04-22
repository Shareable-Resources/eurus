package user

import (
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/env"
	"eurus-backend/marketing/reward"

	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"

	"github.com/shopspring/decimal"
)

func GetWalletAddressByUserId(req *QueryRecentTransactionDetailsRequest, database *database.ReadOnlyDatabase) (string, error) {
	walletAddress, err := DBGetWalletAddressByUserId(req.UserId, database)
	if err != nil {
		return "", err
	}
	return walletAddress, nil
}

func GetRecentTransaction(server *UserServer, req *QueryRecentTransactionDetailsRequest) *response.ResponseBase {
	tokenString := new(TokenString)
	err := json.Unmarshal([]byte(req.LoginToken), tokenString)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unmarshal error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	req.UserId = tokenString.UserId
	if req.ChainId == 0 {
		req.ChainId = env.DefaultEurusChainId
	}
	transList, err := DBGetRecentTransferTransaction(req, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DBGetRecentTransferTransaction error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	recentTx := new(RecentTransaction)
	for _, trans := range transList {
		recentTxDetails := new(TransferDetail)
		recentTxDetails.TxHash = trans.TxHash
		recentTxDetails.TransDate = trans.TransactionDate.Unix()
		recentTxDetails.Status = trans.Status
		recentTxDetails.TransType = asset.Transfer
		recentTxDetails.ChainLocation = Innet
		recentTxDetails.Amount = trans.Amount.String()
		recentTxDetails.GasPrice = trans.GasPrice.BigInt().Uint64()
		if trans.UserGasUsed != nil {
			recentTxDetails.GasUsed = *trans.UserGasUsed
		} else if trans.TransGasUsed != nil {
			recentTxDetails.GasUsed = *trans.TransGasUsed
		}
		if trans.Chain == server.Config.EthClientChainID {
			recentTxDetails.ChainLocation = Innet
		} else {
			recentTxDetails.ChainLocation = Mainnet
		}
		recentTxDetails.ToAddress = trans.ToAddress
		recentTxDetails.FromAddress = trans.FromAddress
		recentTxDetails.Remarks = trans.Remarks
		recentTxDetails.IsSend = trans.IsSend

		recentTx.TransList = append(recentTx.TransList, recentTxDetails)
	}

	withdrawTxs, err := DBGetRecentWithdrawTransaction(req, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DBGetRecentWithdrawTransaction error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	RecentWithdrawTransactionsHandler(recentTx, withdrawTxs)

	depositTxs, err := DBGetRecentDepositTransaction(req, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DBGetRecentDepositTransaction error: ", err, " nonce: ", req.Nonce)

		return response.CreateErrorResponse(req, foundation.DatabaseError, "Getting RecentDepositTransaction failed: "+err.Error())
	}
	RecentDepositTransactionsHandler(recentTx, depositTxs)

	rewardTxs, err := DBGetRecentDistributedTokenTransaction(req.UserId, req.CurrencySymbol, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DBGetRecentDistributedTokenTransaction error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Getting DistributedTokenTransaction failed: "+err.Error())
	}
	RecentDistributedTokenTransactionsHandler(recentTx, rewardTxs)

	purchaseTxs, err := DBGetRecentPurchaseTransaction(req, server.SlaveDatabase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DBGetRecentPurchaseTransaction error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	RecentPurchaseTransactionHandler(recentTx, purchaseTxs)

	if req.CurrencySymbol == asset.EurusTokenName {
		topUpTxs, err := DBGetRecentTopUpTransaction(req.UserId, server.SlaveDatabase)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("DBGetRecentTopUpTransaction error: ", err, " Nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.DatabaseError, "Getting Top up transaction error: "+err.Error())
		}
		RecentTopUpTransactionsHandler(recentTx, topUpTxs)
	}

	if req.CurrencySymbol != asset.EurusTokenName {
		decimals, err := server.GetDecimalPlaceFromSC(req.CurrencySymbol)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("GetDecimalPlaceFromSC error: ", err, " nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.EthereumError, err.Error())
		} else {
			recentTx.DecimalPlace = decimals
			return response.CreateSuccessResponse(req, recentTx)
		}
	} else {
		recentTx.DecimalPlace = asset.EurusTokenDecimalPoint
	}

	return response.CreateSuccessResponse(req, recentTx)
}

func RecentDepositTransactionsHandler(recentTx *RecentTransaction, depositTxs []*asset.DepositTransaction) {
	for i := 0; i < len(depositTxs); i++ {
		RecentDepositTransactionHandler(recentTx, depositTxs[i])
	}
}

func RecentDepositTransactionHandler(recentTx *RecentTransaction, depositTx *asset.DepositTransaction) {
	status := depositTx.Status
	negative := false

	if status < 0 {
		status = -status
		negative = true
	}
	for ; status >= 0; status -= 10 {
		depositDetail := new(DepositDetail)
		depositDetail.TransType = asset.Deposit
		depositDetail.Amount = depositTx.Amount.String()
		depositDetail.FromAddress = depositTx.MainnetFromAddress
		depositDetail.ToAddress = depositTx.MainnetToAddress
		depositDetail.DestAddress = depositTx.InnetToAddress
		depositDetail.Remarks = depositTx.Remarks
		depositDetail.Status = status
		depositDetail.DepositTxHash = depositTx.MainnetTransHash
		depositDetail.TransDate = depositTx.CreatedDate.Unix()
		if !depositTx.MainnetGasUsed.Equals(decimal.Zero) {
			depositDetail.GasPrice = depositTx.MainnetGasFee.Div(depositTx.MainnetGasUsed).BigInt().Uint64()
		}
		depositDetail.GasUsed = depositTx.MainnetGasUsed.BigInt().Uint64()
		if negative {
			depositDetail.Status = -status
			negative = false
		}
		switch status {
		case asset.DepositReceiptCollected:
			depositDetail.ChainLocation = Mainnet
			depositDetail.TransDate = depositTx.MainnetTransDate.Unix()
		case asset.DepositAssetCollected:
			depositDetail.ChainLocation = Mainnet
			if depositTx.MainnetCollectTransHash != "" {
				depositDetail.TxHash = depositTx.MainnetCollectTransHash
				depositDetail.TransDate = depositTx.MainnetCollectTransDate.Unix()
			}
		case asset.DepositMintRequesting:
			depositDetail.ChainLocation = Innet

			if depositTx.MintTransId != nil {
				depositDetail.MintTransId = *depositTx.MintTransId
			} else {
				depositDetail.MintTransId = 0
			}
			depositDetail.TransDate = depositTx.MainnetTransDate.Unix()
		case asset.DepositCompleted:
			depositDetail.TxHash = depositTx.MintTransHash
			if depositTx.MintDate != nil {
				depositDetail.TransDate = depositTx.MintDate.Unix()
			}
			depositDetail.ChainLocation = Innet
		}
		recentTx.TransList = append(recentTx.TransList, depositDetail)
	}
}

func RecentWithdrawTransactionsHandler(recentTx *RecentTransaction, withdrawTxs []*asset.WithdrawTransaction) {
	for i := 0; i < len(withdrawTxs); i++ {
		RecentWithdrawTransactionHandler(recentTx, withdrawTxs[i])
	}
}

func RecentWithdrawTransactionHandler(recentTx *RecentTransaction, withdrawTx *asset.WithdrawTransaction) {
	status := withdrawTx.Status
	negative := false
	if status < 0 {
		status = -status
		negative = true
	}

	for ; status >= 0; status -= 10 {
		withdrawDetail := new(WithdrawDetail)
		withdrawDetail.TransType = asset.Withdraw
		withdrawDetail.Amount = withdrawTx.Amount.String()
		withdrawDetail.FromAddress = withdrawTx.InnetFromAddress
		withdrawDetail.ToAddress = withdrawTx.ApprovalWalletAddress
		withdrawDetail.Remarks = withdrawTx.Remarks
		withdrawDetail.Status = status
		withdrawDetail.TargetAddress = withdrawTx.MainnetToAddress
		withdrawDetail.WithdrawTxHash = withdrawTx.RequestTransHash
		withdrawDetail.AdminFee, _ = withdrawTx.AdminFee.Float64()
		withdrawDetail.TransDate = withdrawTx.CreatedDate.Unix()
		withdrawDetail.GasPrice = withdrawTx.GasPrice.BigInt().Uint64()
		withdrawDetail.GasUsed = withdrawTx.UserGasUsed.BigInt().Uint64()

		if negative {
			withdrawDetail.Status = -status
			negative = false
		}
		switch status {
		case asset.StatusPendingApproval:
			withdrawDetail.ChainLocation = Innet
			if withdrawTx.RequestTransId != nil {
				withdrawDetail.RequestTransId = *withdrawTx.RequestTransId
			}
			withdrawDetail.TxHash = withdrawTx.RequestTransHash
			withdrawDetail.TransDate = withdrawTx.CreatedDate.Unix()
			break
		case asset.StatusApproved:
		case asset.StatusRejected:
			withdrawDetail.TransDate = withdrawTx.ReviewDate.Unix()
			withdrawDetail.TxHash = withdrawTx.ReviewTransHash
			withdrawDetail.ChainLocation = Innet
			break
		case asset.StatusBurnConfirming:
			withdrawDetail.BurnTransId = withdrawTx.BurnTransId
			withdrawDetail.ChainLocation = Innet
			withdrawDetail.TxHash = withdrawTx.ReviewTransHash
			withdrawDetail.TransDate = withdrawTx.ReviewDate.Unix()
		case asset.StatusBurned:
			withdrawDetail.TxHash = withdrawTx.BurnTransHash
			withdrawDetail.TransDate = withdrawTx.BurnDate.Unix()
			break
		case asset.StatusCompleted:
			withdrawDetail.TxHash = withdrawTx.MainnetTransHash
			withdrawDetail.TransDate = withdrawTx.MainnetTransDate.Unix()
			withdrawDetail.ChainLocation = Mainnet
			break
		}
		recentTx.TransList = append(recentTx.TransList, withdrawDetail)
	}

}

func RecentPurchaseTransactionHandler(recentTx *RecentTransaction, purchaseTxs []*asset.PurchaseTransaction) {
	for _, trans := range purchaseTxs {
		purchaseDetail := new(PurchaseDetail)
		purchaseDetail.TransType = trans.PurchaseType
		purchaseDetail.FromAddress = trans.FromAddress
		purchaseDetail.ToAddress = trans.ToAddress
		purchaseDetail.Amount = trans.Amount.String()
		purchaseDetail.ChainLocation = Innet
		purchaseDetail.GasFee = trans.GasFee.BigInt().Uint64()
		purchaseDetail.GasPrice = trans.GasPrice.BigInt().Uint64()
		purchaseDetail.GasUsed = trans.UserGasUsed
		purchaseDetail.TransDate = trans.CreatedDate.Unix()
		purchaseDetail.ProductId = trans.ProductId
		if trans.Quantity != nil {
			quantity := trans.Quantity.BigInt().Uint64()
			purchaseDetail.Quantity = &quantity
		}
		purchaseDetail.Status = trans.Status

		recentTx.TransList = append(recentTx.TransList, purchaseDetail)
	}
}

func RecentDistributedTokenTransactionsHandler(recentTx *RecentTransaction, distributedTokenList []*reward.DistributedToken) {
	for _, distributedToken := range distributedTokenList {
		rewardDetail := new(DistributedTokenDetail)
		rewardDetail.TxHash = distributedToken.TxHash
		rewardDetail.FromAddress = distributedToken.FromAddress
		rewardDetail.ToAddress = distributedToken.ToAddress
		rewardDetail.Amount = distributedToken.Amount.String()
		if *distributedToken.Chain == uint64(env.DefaultEurusChainId) {
			rewardDetail.ChainLocation = Innet
		} else {
			rewardDetail.ChainLocation = Mainnet
		}

		rewardDetail.TransDate = distributedToken.CreatedDate.Unix()
		rewardDetail.TransType = asset.DistributedToken
		rewardDetail.DistributedType = distributedToken.DistributedType
		rewardDetail.TriggerType = distributedToken.TriggerType

		recentTx.TransList = append(recentTx.TransList, rewardDetail)
	}
}

func RecentTopUpTransactionsHandler(recentTx *RecentTransaction, topUpTxList []*asset.TopUpTransaction) {
	for i := 0; i < len(topUpTxList); i++ {
		RecentTopUpTransactionHandler(recentTx, topUpTxList[i])
	}
}

func RecentTopUpTransactionHandler(recentTx *RecentTransaction, topUpTx *asset.TopUpTransaction) {

	detail := new(TopUpDetail)
	detail.TransType = asset.TopUp
	detail.ChainLocation = Innet
	detail.FromAddress = topUpTx.FromAddress
	detail.ToAddress = topUpTx.ToAddress
	detail.TxHash = topUpTx.TxHash
	detail.IsDirectTopUp = topUpTx.IsDirectTopUp
	detail.GasPrice = topUpTx.GasPrice.BigInt().Uint64()
	detail.Status = topUpTx.Status
	detail.Remarks = topUpTx.Remarks
	detail.TransDate = topUpTx.TransactionDate.Unix()
	detail.GasUsed = topUpTx.UserGasUsed
	detail.TransGasUsed = topUpTx.TransGasUsed
	detail.TransferGas = topUpTx.TransferGas.BigInt().Uint64()
	detail.TargetGas = topUpTx.TargetGas.BigInt().Uint64()

	recentTx.TransList = append(recentTx.TransList, detail)
}
