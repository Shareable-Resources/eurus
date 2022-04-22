package block_cypher_bg

import (
	"eurus-backend/block_cypher_service/block_cypher"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/server"
)

type BlockCypherBgDBProcessor struct {
	Db     *database.Database
	Config *server.ServerConfigBase
}

func NewDbProcessor(db *database.Database, config *server.ServerConfigBase) *BlockCypherBgDBProcessor {
	processor := new(BlockCypherBgDBProcessor)
	processor.Config = processor.Config
	processor.Db = db

	return processor
}

func (me *BlockCypherBgDBProcessor) DbQueryTokenList(coin string, chain string) ([]*block_cypher.BlockCypherToken, error) {
	dbConn, err := me.Db.GetConn()
	if err != nil {
		return nil, err
	}
	var tokenList []*block_cypher.BlockCypherToken = make([]*block_cypher.BlockCypherToken, 0)
	tx := dbConn.Model(block_cypher.BlockCypherToken{}).Where("coin = ? AND chain = ? AND is_enabled = ?", coin, chain, true).Order("email").Find(&tokenList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tokenList, nil
}

func (me *BlockCypherBgDBProcessor) UpdateToken(tokenInfo *block_cypher.BlockCypherToken) error {
	dbConn, err := me.Db.GetConn()
	if err != nil {
		return err
	}

	tx := dbConn.Where("email = ? ", tokenInfo.Email).
		Select("HitsApiPerHour", "HitsApiPerDay", "HitsConfidencePerHour", "LimitApiPerHour", "LimitApiPerDay", "LimitConfidencePerHour", "LastModifiedDate", "Score", "UsedCount").
		Updates(tokenInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
