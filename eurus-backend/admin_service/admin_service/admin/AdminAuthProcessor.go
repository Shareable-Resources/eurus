package admin

import (
	"encoding/json"
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"fmt"
)

func getAdminUserIdFromLoginToken(loginToken auth_base.ILoginToken) (uint64, error) {
	userIdStr := loginToken.GetUserId()
	var tokenMap map[string]interface{} = make(map[string]interface{})

	err := json.Unmarshal([]byte(userIdStr), &tokenMap)
	if err != nil {
		return 0, err
	}

	userIdObj, ok := tokenMap["adminUserId"]
	if !ok {
		return 0, errors.New("Invalid login token")
	}

	if userId, ok1 := userIdObj.(float64); ok1 {
		return uint64(userId), nil
	}
	return 0, errors.New("Invalid admin user id format")
}

func generateLoginTokenUserInfo(userId uint64) string {
	return fmt.Sprintf("{\"adminUserId\":%d}", userId)
}

func checkPermission(featureId FeatureId, permissionId PermissionId, req *request.RequestBase, dbProcessor *AdminDBProcessor) *response.ResponseBase {
	adminId, err := getAdminUserIdFromLoginToken(req.LoginToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Error get admin id : ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, "Invalid login token")
	}
	isAllow, err := dbProcessor.DbAdminHasPermission(adminId, featureId, permissionId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbAdminHasPermission error: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	if !isAllow {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Permission denied")
	}

	return nil
}
