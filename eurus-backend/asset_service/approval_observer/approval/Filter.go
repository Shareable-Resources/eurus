package approval

// func FilterTransfer(ext *ExtractedTransaction, txData []byte) *ExtractedTransaction {
// 	var err error
// 	dataMap, err, state := ExtractTransferData(txData)
// 	if state == ExtractTransferState()("Failed") {
// 		log.GetLogger(log.Name.Root).Error("Failed!!! Unable to input extract transfer data from tx data: ", err.Error(), "Transaction hash is : ", strings.ToLower(ext.TxHash))
// 		ext = nil
// 	} else if state == ExtractTransferState()("Success") {
// 		_, err = ext.GetSender()
// 		if err != nil {
// 			log.GetLogger(log.Name.Root).Errorln("Unable to input extract transfer data from tx data: ", err.Error(), "Transaction hash is : ", strings.ToLower(ext.TxHash))
// 		}
// 		//ext.extractTransactionFrom()
// 		ext.To = dataMap["recipient"].(common.Address).Hex()
// 	}
// 	return ext
// }
