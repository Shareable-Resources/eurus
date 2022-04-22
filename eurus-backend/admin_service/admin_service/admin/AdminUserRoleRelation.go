package admin

import "eurus-backend/foundation/database"

type AdminUserRoleRelation struct {
	database.DbModel
	AdminId   uint64 `json:"adminId"`
	RoleId    uint64 `json:"roleId"`
	CreatedBy uint64 `json:"createdBy"`
}
