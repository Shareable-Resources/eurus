package admin

import (
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/two_fa/ga"
	"time"
)

func processAdminLogin(server *AdminServer, req *AdminLoginRequest) *response.ResponseBase {
	adminUser, err := server.dbProcessor.DbVerifyAdminPassword(req.UserName, req.Password)
	if err != nil || adminUser == nil {
		log.GetLogger(log.Name.Root).Errorln("DbVerifyAdminPassword failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.UserNotFound, "User not found or incorrect password")
	}

	res := new(AdminLoginResponse)

	if adminUser.Secret == "" || adminUser.Status == AdminWaitForBindGA {
		res.NextAction = BindGA

		var secret, qrCode string
		if adminUser.Secret == "" {
			secret, qrCode, err = ga.EnableTwoFA(adminUser.Username, "Eurus", "admin")
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("EnableTwoFA failed: ", err, " nonce: ", req.Nonce)
				return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
			}

			err = server.dbProcessor.DbUpdateAdminSecret(adminUser.Id, secret)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("DbUpdateAdminSecret failed: ", err, " nonce: ", req.Nonce)
				return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
			}
		} else {
			qrCode = ga.GetQRCode(adminUser.Username, "admin", "Eurus", adminUser.Secret)
		}
		res.GASecretQRCode = qrCode
	} else {
		res.NextAction = LoginGA
	}
	userInfo := generateLoginTokenUserInfo(adminUser.Id)
	loginToken, serverErr := server.AuthClient.RequestNonRefreshableLoginToken(userInfo, 300, int16(auth.NonRefreshableToken))
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("RequestNonRefreshableLoginToken error: ", serverErr.Message, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, serverErr.ReturnCode, serverErr.Message)
	}
	res.AccessToken = loginToken.GetToken()
	return response.CreateSuccessResponse(req, res)
}

func processVerifyGA(server *AdminServer, req *VerifyGARequest) *response.ResponseBase {

	userId, err := getAdminUserIdFromLoginToken(req.LoginToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GetAdminUserIdFromLoginToken failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.LoginTokenInvalid, err.Error())
	}

	adminUser, err := server.dbProcessor.DbGetAdminUserById(userId)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbGetAdminUserById failed: ", err, " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}
	isSuccess := ga.VerifyTwoFACode(adminUser.Secret, req.GACode)
	if !isSuccess {
		return response.CreateErrorResponse(req, foundation.InvalidArgument, "Invalid GA code")
	}

	_, serverErr := server.AuthClient.RevokeLoginToken(req.LoginToken.GetToken())
	if serverErr != nil {
		log.GetLogger(log.Name.Root).Errorln("RevokeLoginToken failed: ", serverErr.Message, " Nonce: ", req.Nonce)
	}
	if adminUser.Status == AdminWaitForBindGA {
		err = server.dbProcessor.DbUpdateAdminUserStatus(adminUser.Id, AdminNormal)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to DbUpdateAdminUserStatus: ", err, " Nonce: ", req.Nonce)
			return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
		}
	}

	userInfo := generateLoginTokenUserInfo(userId)
	newLoginToken, err := server.AuthClient.GenerateLoginToken(userInfo)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GenerateLoginToken failed: ", serverErr, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	loginTime := time.Now()
	err = server.dbProcessor.DbUpdateAdminUserLoginInfo(adminUser.Id, req.LoginIp, loginTime)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to DbUpdateAdminUserStatus: ", err, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	res := new(VerifyGAResponse)
	res.AccessToken = newLoginToken.GetToken()
	featureList, err := server.dbProcessor.DbQueryAdminEffectivePermission(adminUser.Id)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbQueryAdminEffectivePermission error: ", err, " Nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if server.Config.ElasticSearchPath != "" {
		go func() {
			loginLog := new(ElasticAdminLoginData)
			loginLog.AdminId = adminUser.Id
			loginLog.LoginIp = req.LoginIp
			loginLog.LoginTime = loginTime
			loginLog.UserName = adminUser.Username
			err := server.elasticSearch.InsertLog("/admin_user/login_data", loginLog)
			if err != nil {
				log.GetLogger(log.Name.Root).Warnln("Elastic search insert log failed: ", err)
			}
		}()
	}
	res.FeaturePermissionList = featureList

	return response.CreateSuccessResponse(req, res)
}

func processQueryFeatureList(dbProcessor *AdminDBProcessor, req *request.RequestBase) *response.ResponseBase {
	featureList, err := dbProcessor.DbQueryAllFeaturePermission()
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	return response.CreateSuccessResponse(req, featureList)

}
