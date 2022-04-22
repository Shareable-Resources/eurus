package deposit

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/ethereum"
)

type DepositProcessorContext struct {
	db                     *database.Database
	slaveDb                *database.ReadOnlyDatabase
	retrySetting           foundation.IRetrySetting
	LoggerName             string
	MainnetRescanCounter   *ethereum.ScanBlockCounter
	SideChainRescanCounter *ethereum.ScanBlockCounter
}

func NewDepositProcessorContext(db *database.Database, slaveDb *database.ReadOnlyDatabase, retrySetting foundation.IRetrySetting, loggerName string,
	mainnetRescanCounter, sideChainRescanCounter *ethereum.ScanBlockCounter) *DepositProcessorContext {
	processorContext := new(DepositProcessorContext)
	processorContext.db = db
	processorContext.slaveDb = slaveDb
	processorContext.retrySetting = retrySetting
	processorContext.LoggerName = loggerName
	processorContext.MainnetRescanCounter = mainnetRescanCounter
	processorContext.SideChainRescanCounter = sideChainRescanCounter
	return processorContext
}
