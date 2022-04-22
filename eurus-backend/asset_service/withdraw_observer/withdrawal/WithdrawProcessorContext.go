package withdrawal

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
)

// Cannot use NewWithdrawProcessorContext to construct it, because db cannot be known at first
type WithdrawProcessorContext struct {
	db           *database.Database
	slaveDb      *database.ReadOnlyDatabase
	retrySetting foundation.IRetrySetting
	loggerName   string
}
