package auth

import (
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation/database"
)

type AuthDataSource struct {
	Config      *AuthServerConfig
	DB          *database.Database
	SlaveDB     *database.ReadOnlyDatabase
	ServiceInfo map[int64]conf_api.AuthService
}
