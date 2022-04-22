package block_cypher_bg

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/log"
	"time"

	"github.com/blockcypher/gobcy"
)

type BlockCypherBgProcessor struct {
	dbProcessor *BlockCypherBgDBProcessor
	config      *BlockCypherBackgroundConfig
}

func NewBlockCypherBgProcessor(dbProcessor *BlockCypherBgDBProcessor, config *BlockCypherBackgroundConfig) *BlockCypherBgProcessor {
	processor := &BlockCypherBgProcessor{dbProcessor: dbProcessor, config: config}
	return processor
}

func (me *BlockCypherBgProcessor) processUpdateToken() {

	for _, coin := range me.config.CoinList {
		tokenList, err := me.dbProcessor.DbQueryTokenList(coin, me.config.Chain)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to query block cypher token list: ", err)
			return
		}

		for _, apiTokenModel := range tokenList {
			btcApi := gobcy.API{apiTokenModel.Token, coin, me.config.Chain}

			err := retryFunction(me.config, func() (bool, error) {

				usage, err := btcApi.CheckUsage()
				if err != nil {
					return false, err
				}
				apiTokenModel.HitsApiPerDay = usage.Hits.PerDay
				apiTokenModel.HitsApiPerHour = usage.Hits.PerHour
				apiTokenModel.HitsConfidencePerHour = usage.Hits.ConfPerHour
				apiTokenModel.LimitApiPerDay = usage.Limits.PerDay
				apiTokenModel.LimitApiPerHour = usage.Limits.PerHour
				apiTokenModel.LimitConfidencePerHour = usage.Limits.ConfPerHour
				apiTokenModel.UsedCount = 0
				apiTokenModel.LastModifiedDate = time.Now()

				return false, nil
			})
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to query usage: ", err, " email: ", apiTokenModel.Email, " coin: ", coin)
			}

			var score int = int((apiTokenModel.LimitApiPerHour-apiTokenModel.HitsApiPerHour)*100000 + (apiTokenModel.LimitApiPerDay - apiTokenModel.HitsApiPerDay))
			apiTokenModel.Score = score

			err = me.dbProcessor.UpdateToken(apiTokenModel)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("UpdateToken to db failed: ", err, " email: ", apiTokenModel.Email, " coin: ", coin)
			}
		}
	}
}

func retryFunction(retryConfig foundation.IRetrySetting, invoker func() (bool, error)) error {
	var err error
	var isFatal bool
	for i := 0; i < retryConfig.GetRetryCount(); i++ {
		isFatal, err = invoker()
		if isFatal {
			break
		}
		if err != nil {
			time.Sleep(retryConfig.GetRetryInterval() * time.Second)
			continue
		} else {
			return nil
		}
	}

	return err
}
