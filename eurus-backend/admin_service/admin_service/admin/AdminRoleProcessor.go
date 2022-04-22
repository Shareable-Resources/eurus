package admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

func processQueryRoleList(dbProcessor *AdminRoleDBProcessor, req *admin_common.QueryRoleListRequest) *response.ResponseBase {

	res := checkPermission(FeatureRoleManagement, PermissionQuery, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}

	roleList, err := dbProcessor.DbQueryAdminRoleList(req.RoleName, req.State)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbQueryAdminRoleList error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, roleList)
}

func processQueryRoleDetail(dbProcessor *AdminRoleDBProcessor, req *admin_common.QueryRoleDetailRequest) *response.ResponseBase {
	res := checkPermission(FeatureRoleManagement, PermissionQuery, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}

	role, featureList, err := dbProcessor.DbQueryRoleDetail(req.RoleId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbQueryRoleDetail failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	resObj := &admin_common.QueryRoleDetailResponse{
		Role:           admin_common.Role(*role),
		PermissionList: featureList,
	}

	return response.CreateSuccessResponse(req, resObj)
}

func processCreateRole(dbProcessor *AdminRoleDBProcessor, req *admin_common.CreateRoleRequest) *response.ResponseBase {
	res := checkPermission(FeatureRoleManagement, PermissionNew, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}

	//Since check permission already check the login token is valid, so no need to check again
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)
	if req.RoleName == "" {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Empty role name")
	}
	if req.PermissionList == nil {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Feature permission list not found")
	}

	roleId, err := dbProcessor.DbCreateRole(req.RoleName, req.Description, adminId, req.PermissionList)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbCreateRole error: ", err)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	return response.CreateSuccessResponse(req, roleId)
}

func processUpdateRole(dbProcessor *AdminRoleDBProcessor, req *admin_common.UpdateRoleRequest) *response.ResponseBase {
	res := checkPermission(FeatureRoleManagement, PermissionUpdate, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}
	//Since check permission already check the login token is valid, so no need to check again
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)
	if req.UpdateField == admin_common.RoleFieldNone {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "No field update")
	}

	if req.UpdateField&admin_common.RoleFieldName > 0 {
		if req.RoleName == "" {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Role name is empty")
		}
	}

	var state admin_common.RoleState
	if req.IsEnabled {
		state = admin_common.RoleEnabled
	} else {
		state = admin_common.RoleDisabled
	}

	err := dbProcessor.DbUpdateRole(req.RoleId, adminId, req.UpdateField, req.RoleName, req.Description, state, req.FeaturePermissionList)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	return response.CreateSuccessResponse(req, nil)

}

func processDeleteRole(dbProcessor *AdminRoleDBProcessor, req *admin_common.DeleteRoleRequest) *response.ResponseBase {
	res := checkPermission(FeatureRoleManagement, PermissionDelete, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}
	//Since check permission already check the login token is valid, so no need to check again
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)

	err := dbProcessor.DbDeleteRole(req.RoleId, adminId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to update role: ", err)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, nil)
}
