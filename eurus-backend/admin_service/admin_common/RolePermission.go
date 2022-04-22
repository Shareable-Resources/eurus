package admin_common

import "eurus-backend/foundation/database"

type RolePermission struct {
	database.DbModel
	FeaturePermissionPair
	RoleId uint64 `json:"roleId"`
}
