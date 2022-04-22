package bc_indexer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var bigZero *big.Int = big.NewInt(0)

func FilterERC20TransferTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (bool, *ExtractedTransaction, error) {
	return filterEurusERC20TransferTransactionImpl(context, tx, block, "ERC20")
}

func FilterEurusERC20TransferTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (bool, *ExtractedTransaction, error) {
	if context.IsMainnet {
		return false, nil, errors.New("Must be a sidechain transaction")
	}

	return filterEurusERC20TransferTransactionImpl(context, tx, block, "EurusERC20")
}

func filterEurusERC20TransferTransactionImpl(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block, contractName string) (bool, *ExtractedTransaction, error) {

	if tx.To() == nil {
		return false, nil, nil
	}

	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), contractName, "transfer")
	if err != nil {
		//Method not found
		return false, nil, nil
	}
	var ok bool
	ext := new(ExtractedTransaction)
	ext.OriginalTransaction = tx
	ext.IsMainnetTrans = context.IsMainnet
	ext.AssetName, ok = context.AssetAddressMap.AssetList[*tx.To()]
	ext.Block = block
	ext.TransactionType = asset.Transfer

	if !ok {
		return false, nil, nil
	}
	ext.Amount, ok = args["amount"].(*big.Int)
	if !ok {
		return false, nil, errors.New("Unable to get amount field from ERC20 transfer transaction")
	}
	ext.TxHash = tx.Hash().Hex()
	receipant, ok := args["recipient"].(common.Address)
	if !ok {
		log.GetLogger(context.LoggerName).Errorln("Unable to get convert argument to recipient address. Trans: ", tx.Hash().Hex())
		return false, nil, errors.New("Unable to get recipient field from ERC20 transfer transaction")
	}

	ext.SetTo(ethereum.ToLowerAddressString(receipant.String()))

	if context.IsMainnet && bytes.Equal(receipant.Bytes(), context.MainnetPlatformWalletAddress.Bytes()) {
		ext.ToUser = context.MainnetPlatformWalletUser
	} else {
		receiver, err := DbGetUserByWalletAddress(context, receipant.Hex())
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query receiver user for trans hash: ", ext.TxHash, " error: ", err)
		}
		if receiver != nil && receiver.Id > 0 {
			ext.ToUser = receiver
		}
	}

	senderAddr, err := ext.GetSender()
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Unable to get sender address for trans hash: ", ext.TxHash, " error: ", err)
	} else {
		senderUser, err := DbGetUserByWalletAddress(context, senderAddr)
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query sender user for trans hash: ", ext.TxHash, " error: ", err)
			return false, nil, err
		}
		if senderUser.Id > 0 {
			ext.FromUser = senderUser
		}
	}

	ext, _, err = extractReceipt(context, ext)
	if err != nil {
		return false, nil, err
	}

	ext.CreatedDate = time.Unix(int64(block.Time()), 0)

	return true, ext, nil
}

func FilterEurusERC20MerchantTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (*ExtractedTransaction, error) {
	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "EurusERC20", "depositToDApp")
	var ext *ExtractedTransaction
	var ok bool

	if err != nil {
		args, err, _ = ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "EurusERC20", "purchase")
		if err != nil {
			return nil, nil
		} else {
			//Merchant purchase transaction
			ext = new(ExtractedTransaction)
			ext.Amount, ok = args["amount"].(*big.Int)
			if !ok {
				log.GetLogger(context.LoggerName).Error("Unable to get amount from tx input argument. Tx hash: ", tx.Hash().Hex())
				return nil, errors.New("Amount not found")
			}
			ext.Quantity, ok = args["quantity"].(*big.Int)
			if !ok {
				log.GetLogger(context.LoggerName).Error("Unable to get quantity from tx input argument. Tx hash: ", tx.Hash().Hex())
				return nil, errors.New("Quantity not found")
			}
			ext.ProductId, ok = args["productId"].(*big.Int)
			if !ok {
				log.GetLogger(context.LoggerName).Error("Unable to get quantity from tx input argument. Tx hash: ", tx.Hash().Hex())
				return nil, errors.New("Quantity not found")
			}
			var addr common.Address
			addr, ok = args["dappAddress"].(common.Address)
			if !ok {
				log.GetLogger(context.LoggerName).Error("Unable to get dappAddress from tx input argument. Tx hash: ", tx.Hash().Hex())
				return nil, errors.New("Address not found")
			}
			ext.to = ethereum.ToLowerAddressString(addr.Hex())

			extraData, _ := args["extraData"].([]byte)
			if extraData != nil {
				ext.Remarks = hex.EncodeToString(extraData)
			}
			ext.TransactionType = asset.Purchase
		}
	} else {
		//DApp deposit transaction
		ext = new(ExtractedTransaction)
		ext.Amount, ok = args["amount"].(*big.Int)
		if !ok {
			log.GetLogger(context.LoggerName).Error("Unable to get amount from tx input argument. Tx hash: ", tx.Hash().Hex())
			return nil, errors.New("Amount not found")
		}
		var addr common.Address
		addr, ok = args["dappAddress"].(common.Address)
		if !ok {
			log.GetLogger(context.LoggerName).Error("Unable to get dappAddress from tx input argument. Tx hash: ", tx.Hash().Hex())
			return nil, errors.New("Address not found")
		}
		ext.to = ethereum.ToLowerAddressString(addr.Hex())
		ext.TransactionType = asset.MerchantDeposit
		extraData, _ := args["extraData"].([]byte)
		if extraData != nil {
			ext.Remarks = hex.EncodeToString(extraData)
		}
		ext.TransactionType = asset.MerchantDeposit
	}
	ext.Block = block
	ext.AssetName, ok = context.AssetAddressMap.AssetList[*tx.To()]
	if !ok {
		log.GetLogger(context.LoggerName).Error("Unknown EurusERC20 address. tx.To(): ", tx.To().Hex(), " tx hash: ", tx.Hash().Hex())
		return nil, errors.New("Asset name not found")
	}
	ext.OriginalTransaction = tx
	ext.IsMainnetTrans = false
	ext.CreatedDate = time.Unix(int64(block.Time()), 0)
	ext.TxHash = tx.Hash().Hex()
	sender, err := ext.GetSender()
	if err != nil {
		log.GetLogger(context.LoggerName).Error("Unable to get sender from transaction. Trans hash: ", tx.Hash().Hex(), " error: ", err)
		return nil, err
	}

	user, err := DbGetUserByWalletAddress(context, ethereum.ToLowerAddressString(sender))
	if err != nil {
		log.GetLogger(context.LoggerName).Error("Unable to get user from transaction. Trans hash: ", tx.Hash().Hex(), " error: ", err)
		return nil, err
	}

	ext.FromUser = user
	ext, _, err = extractReceipt(context, ext)
	return ext, err
}

func FilterCentralizedUserTransferTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (*ExtractedTransaction, error) {
	if tx.To() == nil {
		return nil, nil
	}
	toAddr := strings.ToLower(tx.To().Hex())
	//Check if the to address is a centralized user wallet address (a smart contract address)
	centralizedUser, err := DbGetUserByWalletAddress(context, toAddr)
	if err != nil {
		return nil, err
	}
	if centralizedUser == nil || centralizedUser.Id == 0 || centralizedUser.IsMetamaskAddr {
		return nil, nil
	}

	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "UserWallet", "requestTransferV1")

	if err != nil {
		if strings.Contains(err.Error(), "Method not found") {
			args, err, _ = ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "UserWallet", "directRequestTransfer")
			if err != nil {
				if strings.Contains(err.Error(), "Method not found") {
					return nil, nil
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err == nil {
		//Transaction is a transferRequest smart contract call
		receipt, err := context.EthClient.QueryEthReceipt(tx)
		if err != nil {
			return nil, err
		}

		ext := new(ExtractedTransaction)
		ext.IsMainnetTrans = context.IsMainnet
		ext.CreatedDate = time.Unix(int64(block.Time()), 0)
		ext.FromUser = centralizedUser
		ext.TransGasUsed = receipt.GasUsed
		ext.Block = block
		ext.TransactionType = asset.Transfer

		ext.SetSender(toAddr)
		if receipt.Status != 0 {
			ext.Status = int16(asset.TransferStatusPending)

			transferRequestEvent, ok := ethereum.DefaultABIDecoder.GetABI("UserWallet").Events["TransferRequestEvent"]
			if !ok {
				return nil, errors.New("ABI UserWallet TransferRequestEvent not found")
			}

			var found bool
			for _, logObj := range receipt.Logs {
				if logObj.Topics[0] == transferRequestEvent.ID {
					if len(logObj.Topics) >= 2 {
						ext.RequestTransId = logObj.Topics[1].Big()
						ext.UserGasUsed = receipt.GasUsed
						found = true
						break
					} else {
						return nil, errors.New("RequestTransId not found at TransferRequestEvent. Trans hash: " + ext.TxHash)
					}
				}
			}
			if !found {
				return nil, errors.New("Neither TransferRequestEvent nor TransferRequestFailed event found. Trans hash: " + ext.TxHash)
			}

		} else {
			ext.Status = int16(asset.TransferStatusPending * asset.TransferStatusError)
			receiptByte, _ := json.Marshal(receipt)
			if receiptByte != nil {
				remarks := asset.NewRemarksJsonFromReceipt(receiptByte)
				ext.Remarks = remarks.String()
			}
		}

		ext.TxHash = tx.Hash().Hex()
		ext.OriginalTransaction = tx
		amount, ok := args["amount"].(*big.Int)
		if !ok {
			return nil, errors.New("Unable to get amount from input argument")
		}
		ext.Amount = amount
		assetName, ok := args["assetName"].(string)
		if !ok {
			return nil, errors.New("Unable to get assetName from input argument")
		}
		ext.AssetName = assetName
		toUserAddr, ok := args["dest"].(common.Address)
		if !ok {
			return nil, errors.New("Unable to get dest from input argument")
		}
		ext.SetTo(toUserAddr.Hex())
		receiverUser, err := DbGetUserByWalletAddress(context, ext.GetTo())
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query receiver user for trans hash: ", ext.TxHash, " Error: ", err)
		}
		if receiverUser != nil && receiverUser.Id > 0 {
			ext.ToUser = receiverUser
		}
		return ext, nil
	}

	if strings.Contains(err.Error(), "Method not found") {
		return nil, nil
	}
	return nil, err

}

func FilterCentralizedUserConfirmTransferTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (*ExtractedTransaction, error) {
	if tx.To() == nil {
		return nil, nil
	}
	toAddr := strings.ToLower(tx.To().Hex())

	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "UserWallet", "confirmTransaction")
	if err != nil {
		//Method not found
		return nil, nil

	}

	//Check if the to address is a centralized user wallet address (a smart contract address)
	centralizedUser, err := DbGetUserByWalletAddress(context, tx.To().Hex())
	if err != nil {
		return nil, err
	}

	if centralizedUser == nil || centralizedUser.Id == 0 || centralizedUser.IsMetamaskAddr {
		return nil, nil
	}

	transId, ok := args["transactionId"].(*big.Int)
	if !ok {
		log.GetLogger(context.LoggerName).Errorln("Unable to get transactionId argument from transaction data: Trans hash: " + tx.Hash().Hex())
		return nil, errors.New("Unable to get transactionId argument from transaction data: Trans hash: " + tx.Hash().Hex())
	}

	receipt, err := context.EthClient.QueryEthReceipt(tx)
	if err != nil {
		log.GetLogger(context.LoggerName).Errorln("Unable to query receipt for transfer confirmation. Trans hash: " + tx.Hash().Hex())
		return nil, errors.Wrap(err, "Unable to query receipt for transfer confirmation. Trans hash: "+tx.Hash().Hex())
	}

	if receipt.Status == 0 {
		return nil, nil
	}

	userWalletAbi := ethereum.DefaultABIDecoder.GetABI("UserWallet")
	transferEvent, ok := userWalletAbi.Events["TransferEvent"]
	if !ok {
		log.GetLogger(context.LoggerName).Errorln("Cannot find transfer event from ABI. Trans hash: " + tx.Hash().Hex())
		return nil, errors.Wrap(err, "Cannot find transfer event from ABI. Trans hash: "+tx.Hash().Hex())
	}

	var isTransferEventFound bool
	for _, logField := range receipt.Logs {
		if len(logField.Topics) > 0 {
			if logField.Topics[0] == transferEvent.ID {
				isTransferEventFound = true
				break
			}
		}
	}

	if !isTransferEventFound {
		return nil, nil
	}

	ext := new(ExtractedTransaction)
	ext.RequestTransId = transId
	ext.FromUser = centralizedUser
	ext.Block = block
	ext.SetSender(toAddr)
	if receipt.Status > 0 {
		ext.Status = int16(asset.TransferStatusConfirmed)
	} else {
		ext.Status = int16(asset.TransferStatusConfirmed * asset.TransferStatusError)
	}
	ext.ConfirmTransHash = tx.Hash().Hex()
	ext.TransactionType = asset.Transfer
	return ext, nil
}

/// This function is used to filter decentralized user transfer EUN only
func FilterEUNTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (*ExtractedTransaction, error) {

	if len(tx.Data()) == 0 && tx.Value().Cmp(bigZero) > 0 {
		ext := new(ExtractedTransaction)
		ext.IsMainnetTrans = false
		ext.OriginalTransaction = tx
		ext.AssetName = asset.EurusTokenName
		ext.TxHash = strings.ToLower(tx.Hash().Hex())
		ext.to = tx.To().String()
		ext.Amount = tx.Value()
		ext.Status = int16(asset.TransferStatusConfirmed)
		ext.CreatedDate = time.Unix(int64(block.Time()), 0)
		ext.Block = block
		ext.TransactionType = asset.Transfer

		senderAddr, _ := ext.GetSender()
		senderUser, err := DbGetUserByWalletAddress(context, senderAddr)
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query sender user for trans hash: ", ext.TxHash, " Error: ", err)
			return nil, err
		}
		if senderUser != nil && senderUser.Id > 0 {
			ext.FromUser = senderUser
		}

		receiverUser, err := DbGetUserByWalletAddress(context, ext.GetTo())
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to query receiver user for trans hash: ", ext.TxHash, " Error: ", err)
			return nil, err
		}
		if receiverUser != nil && receiverUser.Id > 0 {
			ext.ToUser = receiverUser
		}

		if ext.ToUser == nil && ext.FromUser == nil {
			return nil, nil
		}

		receipt, err := context.EthClient.QueryEthReceiptWithSetting(tx, 3, 40)
		if err == nil {
			if receipt.Status == 0 {
				ext.Status = int16(asset.TransferStatusError * asset.TransferStatusConfirmed)
			}
			ext.TransGasUsed = receipt.GasUsed
			ext.UserGasUsed = receipt.GasUsed
		}

		return ext, nil
	}
	return nil, nil
}

func FilterMainNetETHTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) *ExtractedTransaction {
	if tx.Value().Cmp(bigZero) > 0 {
		if tx.To() == nil {
			return nil
		}

		if bytes.Equal(tx.To().Bytes(), context.EurusUserDepositAddress.Bytes()) {
			//This is a deposit transaction, skip it for deposit observer processing
			return nil
		}

		ext := new(ExtractedTransaction)
		ext.OriginalTransaction = tx
		ext.IsMainnetTrans = true
		ext.AssetName = "ETH"
		ext.Block = block
		ext.Amount = tx.Value()
		if tx.To() == nil {
			return nil
		}
		ext.SetTo(tx.To().String())

		ext.TransactionType = asset.Transfer
		ext.TxHash = strings.ToLower(tx.Hash().Hex())
		ext.Status = int16(asset.TransferStatusConfirmed)
		ext.CreatedDate = time.Unix(int64(block.Time()), 0)
		ext, _, err := extractReceipt(context, ext)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to query receipt for tx hash: ", tx.Hash().Hex())
		}
		return ext
	}
	return nil
}

func FilterSweepDecentralizedTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) *ExtractedTransaction {
	inputArgs, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "EurusUserDeposit", "sweep")
	if err != nil {
		return nil
	}
	ext := new(ExtractedTransaction)
	ext.OriginalTransaction = tx
	ext.IsMainnetTrans = true
	ext.AssetName = inputArgs["assetName"].(string)
	ext.Block = block
	ext.Amount = inputArgs["amount"].(*big.Int)

	ext.TransactionType = asset.Transfer
	ext.TxHash = strings.ToLower(tx.Hash().Hex())
	ext.Status = int16(asset.TransferStatusConfirmed)
	ext.CreatedDate = time.Unix(int64(block.Time()), 0)

	systemUser := new(user.User)
	systemUser.Id = 0
	systemUser.MainnetWalletAddress = ethereum.ToLowerAddressString(context.MainnetPlatformWalletAddress.Hex())

	ext.ToUser = systemUser
	ext.SetSender(ethereum.ToLowerAddressString(context.EurusUserDepositAddress.Hex()))
	ext, _, err = extractReceipt(context, ext)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query receipt for tx hash: ", tx.Hash().Hex())
	}
	return ext

}

func FilterMainNetSender(context *blockChainProcessorContext, ext *ExtractedTransaction) (*ExtractedTransaction, error) {
	sender, _ := ext.GetSender()

	if bytes.Equal(common.HexToAddress(sender).Bytes(), context.SweepServiceInvokerAddress.Bytes()) {
		ext.FromUser = context.SweepServiceInvokerUser
	} else {
		//Only check if the sender is a decentralized user
		senderIsUser, err := DbGetUserByMainnetWalletAddress(context, sender) //Check receiver
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Failed to find user from db. Error : ", err)
			return ext, nil
		}
		if senderIsUser.Id > 0 && senderIsUser.MainnetWalletAddress != "" {
			ext.FromUser = senderIsUser
		}
	}
	return ext, nil
}

func FilterMainNetReceiver(context *blockChainProcessorContext, ext *ExtractedTransaction) (*ExtractedTransaction, error) {
	toAddr := strings.ToLower(ext.GetTo())

	if bytes.Equal(common.HexToAddress(ext.GetTo()).Bytes(), context.MainnetPlatformWalletAddress.Bytes()) {
		ext.ToUser = context.MainnetPlatformWalletUser
	} else {

		receiverUser, err := DbGetUserByMainnetWalletAddress(context, toAddr)
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Failed to find user from db. Error : ", err)
			return ext, nil
		}

		if receiverUser.Id > 0 && receiverUser.MainnetWalletAddress != "" {
			ext.ToUser = receiverUser
		}
	}
	return ext, nil
}

func FilterTransferMethodOnly(ext *ExtractedTransaction, txData []byte) *ExtractedTransaction {
	var err error
	dataMap, err, state := ethereum.DefaultABIDecoder.DecodeABIInputArgument(txData, "EurusERC20", "transfer")
	if dataMap == nil {
		return nil
	}
	if ext.AssetName != asset.EurusTokenName {
		getAmount := dataMap["amount"]
		ext.Amount = getAmount.(*big.Int)
	}
	if state == ethereum.ExtractFailed {
		log.GetLogger(log.Name.Root).Error("Tx Hash: "+ext.OriginalTransaction.Hash().Hex()+" Unable to input extract transfer data from tx data: ", err.Error())
		ext = nil
	} else if state == ethereum.ExtractSuccess {
		_, err = ext.GetSender()
		if err != nil {
			log.GetLogger(log.Name.Root).Fatal("Tx Hash: "+ext.OriginalTransaction.Hash().Hex()+"Unable to input extract transfer data from tx data: ", err.Error())
		}
		ext.SetTo(dataMap["recipient"].(common.Address).Hex())
	}
	return ext
}

func FilterEurusUserDepositTransaction(context *blockChainProcessorContext, ext *ExtractedTransaction) *ExtractedTransaction {
	if context.EurusUserDepositAddress == nil {
		return ext
	}
	if bytes.Equal(common.HexToAddress(ext.GetTo()).Bytes(), context.EurusUserDepositAddress.Bytes()) {
		//Omit ERC20 transfer to EurusUserDeposit, leave this to deposit observer to do processing
		return nil
	}

	return ext
}

func FilterTopUpTransaction(context *blockChainProcessorContext, tx *types.Transaction, block *types.Block) (*ExtractedTransaction, error) {

	var toAddr *common.Address
	args, err, _ := ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "UserWallet", "topUpPaymentWallet")
	if err != nil {
		if err.Error() == ethereum.ABIMethodNotFoundError.Error() {
			args, err, _ = ethereum.DefaultABIDecoder.DecodeABIInputArgument(tx.Data(), "UserWallet", "directTopUpPaymentWallet")
			if err != nil {
				if err.Error() == ethereum.ABIMethodNotFoundError.Error() {
					return nil, nil
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		toAddr = new(common.Address)
		*toAddr = args["paymentWalletAddr"].(common.Address)
	}

	//Check if the to address is a centralized user wallet address (a smart contract address)
	centralizedUser, err := DbGetUserByWalletAddress(context, tx.To().Hex())
	if err != nil {
		return nil, err
	}

	if centralizedUser == nil || centralizedUser.Id == 0 || centralizedUser.IsMetamaskAddr {
		return nil, nil
	}

	ext := NewTopUpExtractedTransaction()
	ext.TransactionType = asset.TopUp
	ext.OriginalTransaction = tx

	gas, ok := args["targetGasWei"].(*big.Int)
	if !ok {
		log.GetLogger(context.LoggerName).Errorln("Unable to get targetGasWei from transaction hash: ", tx.Hash().Hex())
		return nil, errors.New("Unable to get targetGasWei from transaction")
	}
	ext.TargetGas = gas

	ext.Block = block
	if toAddr != nil {
		// topUpPaymentWallet function, dest address comes from input argument
		ext.to = ethereum.ToLowerAddressString(toAddr.Hex())

		ext.IsDirectTopUp = false
	} else {
		// directTopUpPaymentWallet function, dest address = sender address
		ext.to, err = ext.GetSender()
		if err != nil {
			return nil, err
		}
		ext.IsDirectTopUp = true
	}

	ext.TxHash = tx.Hash().Hex()
	ext.FromUser = centralizedUser
	ext.IsMainnetTrans = false

	ext.CreatedDate = time.Unix(int64(block.Time()), 0)
	ext.AssetName = asset.EurusTokenName
	var receipt *ethereum.BesuReceipt
	_, receipt, err = extractReceipt(context, &ext.ExtractedTransaction)
	if err != nil {
		return nil, err
	}

	abiObj := ethereum.DefaultABIDecoder.GetABI("UserWalletProxy")
	gasFeeTransferEvent, ok := abiObj.Events["GasFeeTransferred"]
	if !ok {
		panic("UserWalletProxy GasFeeTransferred event not found at ABI")
	}
	topUpPaymentFailedEvent, ok := abiObj.Events["TopUpPaymentWalletFailed"]
	if !ok {
		panic("UserWalletProxy TopUpPaymentWalletFailed event not found at ABI")
	}

	userWaleltAbiObj := ethereum.DefaultABIDecoder.GetABI("UserWallet")
	topUpPaymentEvent, ok := userWaleltAbiObj.Events["TopUpPaymentWalletEvent"]
	if !ok {
		panic("UserWallet TopUpPaymentEvent event not found at ABI")
	}

	if ext.Status == int16(asset.TopUpSuccess) {

		for _, logObj := range receipt.Logs {
			if logObj.Topics[0] == gasFeeTransferEvent.ID {
				if len(logObj.Topics) < 2 {
					log.GetLogger(context.LoggerName).Errorln("GasFeeTransferred event item count invalid. Trans hash: ", ext.TxHash)
					break
				}
				gasUsedHash := logObj.Topics[2]

				gasUsed := big.NewInt(0)
				gasUsed = gasUsed.SetBytes(gasUsedHash.Bytes())
				ext.UserGasUsed = gasUsed.Uint64()

			} else if logObj.Topics[0] == topUpPaymentFailedEvent.ID {
				if len(logObj.Topics) < 2 {
					log.GetLogger(context.LoggerName).Errorln("topUpPaymentFailedEvent event item count invalid. Trans hash: ", ext.TxHash)
					break
				}
				gasUsedHash := logObj.Topics[2]

				gasUsed := big.NewInt(0)
				gasUsed = gasUsed.SetBytes(gasUsedHash.Bytes())
				ext.UserGasUsed = gasUsed.Uint64()
				ext.Status = int16(asset.TopUpError)
				ext.Amount = big.NewInt(0)
				ext.Remarks = parseRevertReason(logObj.Data)

			} else if logObj.Topics[0] == topUpPaymentEvent.ID {
				gasTransferred := logObj.Topics[3]
				gas := big.NewInt(0)
				gas.SetBytes(gasTransferred.Bytes())
				ext.Amount = gas
			}
		}

	} else {
		//Direct topup failed
		reason, err := hex.DecodeString(receipt.RevertReason[2:])
		if err != nil {
			log.GetLogger(context.LoggerName).Errorln("Unable to decode revert reason hex. Error: ", err, " Trans Hash: ", tx.Hash().Hex())
		} else {
			if len(reason) > 4 {
				ext.Remarks = parseRevertReason(reason[4:])
			}
		}
		if ext.Remarks == "" {
			ext.Remarks = receipt.RevertReason
		}
		ext.Amount = big.NewInt(0)
	}

	return &ext.ExtractedTransaction, nil
}

func parseRevertReason(reason []byte) string {
	if len(reason) > 32 {
		offset := big.NewInt(0)
		offset = offset.SetBytes(reason[:32])

		strLen := big.NewInt(0)
		strLen = strLen.SetBytes(reason[32 : 32+offset.Uint64()])
		if strLen.Cmp(bigZero) != 0 {
			return string(reason[32+offset.Uint64() : 32+offset.Uint64()+strLen.Uint64()])
		}
	}
	return ""
}

func extractReceipt(context *blockChainProcessorContext, ext *ExtractedTransaction) (*ExtractedTransaction, *ethereum.BesuReceipt, error) {
	receipt, err := context.EthClient.QueryEthReceiptWithSetting(ext.OriginalTransaction, 3, 40)
	if err != nil {
		return ext, nil, err
	}

	if ext.TransactionType == asset.Transfer {
		if receipt.Status == 0 {
			ext.Status = int16(asset.TransferStatusConfirmed * asset.TransferStatusError)
		} else {
			ext.Status = int16(asset.TransferStatusConfirmed)
		}
	} else if ext.TransactionType == asset.MerchantDeposit || ext.TransactionType == asset.Purchase {
		if receipt.Status == 0 {
			ext.Status = int16(asset.PurchaseStatusError)
		} else {
			ext.Status = int16(asset.PurchaseStatusConfirmed)
		}
	} else if ext.TransactionType == asset.TopUp {
		if receipt.Status == 0 {
			ext.Status = int16(asset.TopUpError)
		} else {
			ext.Status = int16(asset.TopUpSuccess)
		}
	}

	ext.TransGasUsed = receipt.GasUsed
	ext.UserGasUsed = receipt.GasUsed
	ext.EffectiveGasPrice = receipt.EffectiveGasPrice

	return ext, receipt, nil
}
