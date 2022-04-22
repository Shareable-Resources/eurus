package admin_common

import "eurus-backend/foundation/database"

//DB model base class
type Role struct {
	database.DbModel
	Id          uint64    `json:"roleId"`
	RoleName    string    `json:"roleName"`
	ModifiedBy  uint64    `json:"modifiedBy"`
	Description string    `json:"description"`
	State       RoleState `json:"state"`
}

type RoleEx struct {
	Role
	Username string `json:"modifiedByUser"`
}
