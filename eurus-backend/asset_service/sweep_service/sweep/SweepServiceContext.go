package sweep

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"

	"github.com/sirupsen/logrus"
)

type SweepServiceContext struct {
	db           *database.Database
	retrySetting foundation.IRetrySetting
	logger       *logrus.Logger
}

func NewSweepServiceContext(db *database.Database, retrySetting foundation.IRetrySetting, logger *logrus.Logger) *SweepServiceContext {
	processorContext := new(SweepServiceContext)
	processorContext.db = db
	processorContext.retrySetting = retrySetting
	processorContext.logger = logger
	return processorContext
}
