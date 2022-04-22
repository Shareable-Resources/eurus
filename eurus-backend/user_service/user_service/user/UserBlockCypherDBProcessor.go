package user

import (
	"errors"
	"eurus-backend/block_cypher_service/block_cypher"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"

	"gorm.io/gorm"
)

func DbConnectBlockCypher(config *UserServerConfig) (*database.Database, error) {

	db := &database.Database{
		ReadOnlyDatabase: database.ReadOnlyDatabase{
			IP:         config.BlockCypherDBServerIP,
			Port:       config.BlockCypherDBPort,
			UserName:   config.BlockCypherDBUserName,
			Password:   config.BlockCypherDBPassword,
			SchemaName: config.BlockCypherDBSchemaName,
			DBName:     config.BlockCypherDBDatabaseName,
			Logger:     log.GetLogger(log.Name.Root),
		},
	}

	db.SetAESKey(config.DBAESKey)

	_, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DbGetAndUpdateBlockCypherToken(coin string, chain string, db *database.Database) (*block_cypher.BlockCypherToken, error) {
	conn, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	blockCypherToken := new(block_cypher.BlockCypherToken)
	err = conn.Transaction(func(dbTx *gorm.DB) error {
		tx := dbTx.Session(&gorm.Session{}).Where("coin = ? AND chain = ? AND is_enabled = ?", coin, chain, true).
			Order("score DESC, used_count ASC").FirstOrInit(&blockCypherToken)
		if tx.Error != nil {
			return tx.Error
		}
		if blockCypherToken.Token == "" {
			return errors.New("No available token")
		}
		tx = dbTx.Session(&gorm.Session{}).Model(blockCypherToken).
			Where("email = ? AND coin = ? AND chain = ?", blockCypherToken.Email, blockCypherToken.Coin, blockCypherToken.Chain).
			UpdateColumn("used_count", gorm.Expr("used_count + 1"))
		return tx.Error
	})

	if err != nil {
		return nil, err
	}
	return blockCypherToken, nil
}
