package deposit

import (
	"eurus-backend/foundation/log"

	"github.com/ethereum/go-ethereum/core/types"
)

// const retryGetCurrencyNameNum int = 5

func FilterERC20TransactionOnly(tx *AssetTransferTransaction, assetInfo AssetAddressMap) *AssetTransferTransaction {
	if tx == nil {
		log.GetLogger(log.Name.Root).Errorln("AssetTransferTransaction is null ")
	} else if tx.OriginalTransaction == nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid asset trans")
	}

	if tx.TransferLog == nil {
		log.GetLogger(log.Name.Root).Infoln("Transfer event log is null, ignore this transaction")
		return nil
	}

	var ok bool
	var supportedTokenLogList []types.Log = make([]types.Log, 0)

	for _, transferLog := range tx.TransferLog {
		//Check if the ERC20 is a supported token type
		tx.AssetName, ok = assetInfo[transferLog.Address]
		if ok {
			supportedTokenLogList = append(supportedTokenLogList, transferLog)
		}
	}

	if len(supportedTokenLogList) > 0 {
		//Only need supported ERC20 token log objects
		tx.TransferLog = supportedTokenLogList
		return tx
	}
	return nil
}

// func FilterTransferMethodOnly(ext *AssetTransferTransaction, txData []byte) *AssetTransferTransaction {
// 	var err error
// 	dataMap, err, state := ethereum.DefaultABIDecoder.DecodeABIInputArgument(txData, "EurusERC20", "transfer")
// 	if state == ethereum.ExtractFailed {
// 		log.GetLogger(log.Name.Root).Error("Unable to input extract transfer data from tx data: ", err.Error())
// 		ext = nil
// 	} else if state == ethereum.ExtractSuccess {
// 		_, err = ext.GetSender()
// 		if err != nil {
// 			log.GetLogger(log.Name.Root).Error("Unable to input extract transfer data from tx data: ", err.Error())
// 			return nil
// 		}

// 		var obj interface{}
// 		var ok bool
// 		if obj, ok = dataMap["_to"]; !ok {
// 			log.GetLogger(log.Name.Root).Error("_to argument not found from")
// 			return nil
// 		}
// 		addr, ok := obj.(common.Address)
// 		if !ok {
// 			log.GetLogger(log.Name.Root).Error("Invalid argument type: ", err, " trans hash: ", ext.Hash().Hex())
// 			return nil
// 		}
// 		ext.Receiptant = addr.Hex()
// 	}
// 	return ext
// }
