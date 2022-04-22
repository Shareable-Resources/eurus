package admin

import (
	"eurus-backend/foundation/database"
	"time"
)

type AdminUserStatus int

const (
	AdminDisabled AdminUserStatus = iota
	AdminNormal
	AdminWaitForBindGA
	AdminDeleted
)

//DB model
type AdminUser struct {
	database.DbModel
	Id            uint64          `json:"id"`
	Username      string          `json:"userName"`
	Password      string          `json:"-"`
	Secret        string          `json:"-"`
	Status        AdminUserStatus `json:"status"`
	ModifiedBy    uint64          `json:"modifiedBy"`
	LastLoginTime *time.Time      `json:"lastLoginTime"`
	LoginIp       string          `json:"loginIp"`
}

type SimplifiedAdminRole struct {
	AdminId  uint64 `json:"-"`
	RoleId   uint64 `json:"roleId"`
	RoleName string `json:"roleName"`
}
type AdminUserEx struct {
	AdminUser
	ModifiedByUser string                 `json:"modifiedByUser"`
	RoleList       []*SimplifiedAdminRole `gorm:"-" json:"roleList"`
}
