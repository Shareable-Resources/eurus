package userObserver

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
)

type UserObserverContext struct {
	db           *database.Database
	slaveDb      *database.ReadOnlyDatabase
	retrySetting foundation.IRetrySetting
	LoggerName   string
}

func NewUserObserverContext(db *database.Database, slaveDb *database.ReadOnlyDatabase, retrySetting foundation.IRetrySetting, loggerName string) *UserObserverContext {
	processorContext := new(UserObserverContext)
	processorContext.db = db
	processorContext.slaveDb = slaveDb
	processorContext.retrySetting = retrySetting
	processorContext.LoggerName = loggerName
	return processorContext
}
