package bc_indexer

import (
	"strings"
	"time"
)

func DbAddTxIndex(context *blockChainProcessorContext, ext *ExtractedTransaction, isFrom bool) error {
	var err error = nil
	for i := 0; i < context.Config.GetRetryCount(); i++ {
		err = AddTxIndexToDB(context.Db, ext, int(context.EthClient.ChainID.Int64()), isFrom)
		if err == nil {
			break
		} else if err != nil && strings.Contains(err.Error(), "Database Network Error") {
			time.Sleep(time.Duration(context.Config.RetryInterval) * time.Second)
			continue
		} else if err != nil {
			return err
		}
	}
	return err
}
