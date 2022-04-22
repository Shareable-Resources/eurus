package admin

type FeatureId uint64

const (
	FeatureRoleManagement    FeatureId = 20
	FeatureAccountManagement FeatureId = 21
)

type PermissionId uint64

const (
	PermissionQuery PermissionId = iota + 1
	PermissionNew
	PermissionUpdate
	PermissionDelete
	PermissionApproval
)
