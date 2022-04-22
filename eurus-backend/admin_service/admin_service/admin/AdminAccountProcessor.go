package admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
)

func processQueryAccountList(dbProcessor *AdminAccountDBProcessor, req *admin_common.QueryAccountListRequest) *response.ResponseBase {
	res := checkPermission(FeatureAccountManagement, PermissionQuery, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}

	accountList, err := dbProcessor.DbQueryAccountList(req.UserName, req.RoleName, req.State)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbQueryAccountList error: ", err)
		res = response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	} else {
		res = response.CreateSuccessResponse(req, accountList)
	}

	return res
}

func processCreateAccount(dbProcessor *AdminAccountDBProcessor, req *admin_common.CreateAccountRequest) *response.ResponseBase {
	res := checkPermission(FeatureAccountManagement, PermissionNew, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)

	newId, err := dbProcessor.DbCreateAccount(req, adminId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbCreateAccount failed: ", err)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	return response.CreateSuccessResponse(req, newId)
}

func processUpdateAccount(dbProcessor *AdminAccountDBProcessor, req *admin_common.UpdateAccountRequest) *response.ResponseBase {
	res := checkPermission(FeatureAccountManagement, PermissionUpdate, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)

	admin, serverErr := dbProcessor.DbQueryAdminById(req.AdminId)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("Query admin error: ", serverErr.Error())
		return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Message)
	}
	salt := reverse(admin.Username)

	if req.UpdateField&admin_common.AccountFieldPassword > 0 {
		if req.Password == "" {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Missing password")
		}
	}

	err := dbProcessor.DbUpdateAccount(req, salt, adminId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Update admin failed: ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, nil)

}

func processDeleteAccount(dbProcessor *AdminAccountDBProcessor, req *admin_common.DeleteAccountRequest) *response.ResponseBase {
	res := checkPermission(FeatureAccountManagement, PermissionDelete, &req.RequestBase, &dbProcessor.AdminDBProcessor)
	if res != nil {
		return res
	}
	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)

	_, serverErr := dbProcessor.DbQueryAdminById(req.AdminId)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("Query admin error: ", serverErr.Error())
		return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Message)
	}

	err := dbProcessor.DbDeleteAccount(req.AdminId, adminId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Update admin failed: ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, nil)

}

func processChangePassword(dbProcessor *AdminAccountDBProcessor, req *admin_common.ChangePasswordRequest) *response.ResponseBase {

	adminId, _ := getAdminUserIdFromLoginToken(req.LoginToken)
	req.AdminId = adminId

	admin, serverErr := dbProcessor.DbQueryAdminById(req.AdminId)
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("Query admin error: ", serverErr.Error())
		return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Message)
	}
	salt := reverse(admin.Username)

	updateReq := new(admin_common.UpdateAccountRequest)
	updateReq.AdminId = adminId
	updateReq.Password = req.Password
	updateReq.UpdateField = admin_common.AccountFieldPassword

	err := dbProcessor.DbUpdateAccount(updateReq, salt, adminId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Update admin failed: ", err.Error())
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, nil)
}
