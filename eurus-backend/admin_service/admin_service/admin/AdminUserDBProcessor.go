package admin

import (
	"eurus-backend/foundation/database"
)

type AdminUserDBProcessor struct {
	config *AdminServerConfig
	db     *database.Database
}

func NewAdminUserDBProcessor(config *AdminServerConfig, db *database.Database) *AdminUserDBProcessor {
	processor := new(AdminUserDBProcessor)
	processor.config = config
	processor.db = db
	return processor
}
