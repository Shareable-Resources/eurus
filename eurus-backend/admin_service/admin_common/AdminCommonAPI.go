package admin_common

import (
	"eurus-backend/foundation/api/request"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type PaginateRequest struct {
	PageNum  int `json:"-"`
	PageSize int `json:"-"`
}

func (me *PaginateRequest) ParseQueryString(queryUrl *url.URL) (*PaginateRequest, error) {
	pageNumStr := queryUrl.Query().Get("page")
	pageSizeStr := queryUrl.Query().Get("page_size")
	var err error
	if pageNumStr == "" {
		me.PageNum = 1
	} else {
		me.PageNum, err = strconv.Atoi(pageNumStr)
		if err != nil {
			return me, errors.Wrap(err, "Invalid page num")
		}
	}

	if pageSizeStr == "" {
		me.PageSize = 10
	} else {
		me.PageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return me, errors.Wrap(err, "Invalid page size")
		}
	}

	return me, nil
}

type QueryRoleListRequest struct {
	request.RequestBase
	RoleName string    `json:"roleName"`
	State    RoleState `json:"state"`
}
type QueryRoleListResponse []*RoleEx

type QueryRoleDetailRequest struct {
	request.RequestBase
	RoleId uint64 `json:"roleId,string"`
}

type QueryRoleDetailResponse struct {
	Role
	PermissionList []*FeaturePermissionPair `json:"featurePermissionList"`
}

type CreateRoleRequest struct {
	request.RequestBase
	RoleName       string                   `json:"roleName"`
	Description    string                   `json:"description"`
	PermissionList []*FeaturePermissionPair `json:"featurePermissionList"`
}

type CreateRoleResponse uint64

type UpdateRoleRequest struct {
	request.RequestBase
	RoleId                uint64                   `json:"roleId"`
	UpdateField           UpdateRoleField          `json:"updateField"`           //bitwise
	RoleName              string                   `json:"roleName"`              //optional
	Description           string                   `json:"description"`           //optional
	IsEnabled             bool                     `json:"isEnabled"`             //optional
	FeaturePermissionList []*FeaturePermissionPair `json:"featurePermissionList"` //optional
}

type UpdateRoleField uint32

const (
	RoleFieldNone UpdateRoleField = 0
	RoleFieldName UpdateRoleField = 1 << (iota - 1)
	RoleFieldDescription
	RoleFieldState
	RoleFieldPermission
)

type DeleteRoleRequest struct {
	request.RequestBase
	RoleId uint64 `json:"roleId"`
}

type QueryAccountListRequest struct {
	request.RequestBase
	UserName string       `json:"userName"`
	RoleName string       `json:"roleName"`
	State    AccountState `json:"state"`
}

type CreateAccountRequest struct {
	request.RequestBase
	UserName   string `json:"userName" validate:"required,min=4,excludesall= "`
	Password   string `json:"password" validate:"required,min=4,printascii"`
	RoleIdList []int  `json:"roleIdList"`
}

type CreateAccountResponse uint64

type UpdateAccountField uint32

const (
	AccountFieldNone     UpdateAccountField = 0
	AccountFieldPassword UpdateAccountField = 1 << (iota - 1)
	AccountFieldRole
	AccountFieldState
)

type UpdateAccountRequest struct {
	request.RequestBase
	AdminId     uint64
	UpdateField UpdateAccountField `json:"updateField" validate:"required"`
	Password    string             `json:"password" validate:"omitempty,min=4,printascii"`
	RoleIdList  []uint64           `json:"roleIdList"`
	State       int                `json:"state" validate:"omitempty,numeric"`
}

type DeleteAccountRequest struct {
	request.RequestBase
	AdminId uint64
}

type ChangePasswordRequest struct {
	request.RequestBase
	Password string `json:"password" validate:"required,min=4,printascii"`
	AdminId  uint64
}
