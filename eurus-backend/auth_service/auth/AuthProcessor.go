package auth

import (
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/ws/ws_message"
	"strconv"
	"strings"
	"time"
)

func verifySignature(dataSource *AuthDataSource, serverId int64, message string, sign string, timestamp int64) bool {
	//TODO
	if _, ok := dataSource.ServiceInfo[serverId]; !ok {
		log.GetLogger(log.Name.Root).Info("req.Timestamp:", timestamp, " no this serviceId ")
		return false
	}

	publicBase64Key := dataSource.ServiceInfo[serverId].PubKey

	isVerify, err := crypto.VerifyRSASignFromBase64(publicBase64Key, message, sign)
	if err != nil {
		log.GetLogger(log.Name.Root).Info("req.Timestamp:", timestamp, " err.Error(): ", err.Error())
		return false
	}
	return isVerify

}

func RequestLoginToken(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage, serviceId int64) (*RequestLoginTokenResponse, *foundation.ServerError) {

	type jwtUserInfo struct {
		ClientInfo string
		ServiceId  int64
	}

	var loginReq *RequestLoginTokenRequest = req.Request.Data.(*RequestLoginTokenRequest)
	userInfo := &jwtUserInfo{ClientInfo: loginReq.UserId}

	userInfoData, err := json.Marshal(userInfo)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.RequestMalformat, "User data invalid")
	}
	var hasTimeout bool = true

	if serviceId == 0 {
		hasTimeout = false
	}

	var expiredIn int64
	expiredIn = 0
	if loginReq.ExpiredDuration != 0 {
		expiredIn = loginReq.ExpiredDuration
	}

	loginToken, expireTime, err := GenerateJwtToken("EurusAuth", string(userInfoData), hasTimeout, expiredIn)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.RequestMalformat, err.Error())
	}

	err = DbInsertUserSession(dataSource.DB, loginToken, serviceId, expireTime, loginReq.UserId, RefreshableToken)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbInsertUserSession failed: ", err.Error())
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}

	resData := new(RequestLoginTokenResponse)
	resData.Token = loginToken
	resData.ExpiredTime = expireTime
	resData.Type = int16(RefreshableToken)
	resData.LastModifiedDate = time.Now().Unix()

	if dataSource.Config.ResponseDelayInterval > 0 {
		time.Sleep(time.Duration(dataSource.Config.ResponseDelayInterval) * time.Millisecond)
	}
	return resData, nil
}

func VerifyLoginToken(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage, serviceId int64) (*VerifyLoginTokenResponse, *foundation.ServerError) {
	actualReq := req.Request.Data.(*VerifyLoginTokenRequest)

	token := strings.Trim(actualReq.Token, " ")
	verifyMode := actualReq.VerifyMode

	returnCode, _, _ := ValidateJwtoken(token)
	if returnCode != foundation.Success && returnCode != foundation.LoginTokenExired {
		return nil, foundation.NewError(returnCode)
	}

	serviceInfo, ok := dataSource.ServiceInfo[serviceId]
	if !ok {
		serverErr := foundation.NewErrorWithMessage(foundation.LoginTokenInvalid, "Service ID not found. Service ID: "+strconv.FormatInt(serviceId, 10))
		log.GetLogger(log.Name.Root).Errorln(serverErr.Message)
		return nil, serverErr
	}
	var serviceIdList []int64 = make([]int64, 0)

	groupId := serviceInfo.ServiceGroupId
	for serviceId, info := range dataSource.ServiceInfo {
		if info.ServiceGroupId == groupId {
			serviceIdList = append(serviceIdList, serviceId)
		}
	}
	userSession, err := DbQueryUserSession(dataSource.SlaveDB, token, serviceIdList)
	if userSession == nil && err == nil {
		log.GetLogger(log.Name.Root).Errorln("Invalid token or no permission to retrieve")
		serverErr := foundation.NewErrorWithMessage(foundation.LoginTokenInvalid, "Invalid token or no permission to retrieve")
		return nil, serverErr
	}
	if err != nil {
		serverErr := foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
		return nil, serverErr
	}

	if verifyMode == int(auth_base.VerifyModeService) {
		if userSession.ServiceId != 0 {
			serverErr := foundation.NewErrorWithMessage(foundation.LoginTokenInvalid, "Login token is not a server token")
			return nil, serverErr
		}
	} else if verifyMode == int(auth_base.VerifyModeUser) {
		if userSession.ServiceId == 0 {
			serverErr := foundation.NewErrorWithMessage(foundation.LoginTokenInvalid, "Login token is not a user token")
			return nil, serverErr
		}
	}

	if userSession.ExpiredTime != nil {
		// nil case means token never expire
		if userSession.ExpiredTime.Before(time.Now()) {
			serverErr := foundation.NewError(foundation.LoginTokenExired)
			return nil, serverErr
		}
	}

	response := new(VerifyLoginTokenResponse)
	response.UserId = userSession.UserId
	response.Type = userSession.Type
	if userSession.ExpiredTime != nil {
		response.ExpiredTime = userSession.ExpiredTime.Unix()
	}

	return response, nil
}

func RefreshLoginToken(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage, serviceId int64) (*RefreshLoginTokenResponse, *foundation.ServerError) {

	actualReq := req.Request.Data.(*RefreshLoginTokenRequest)

	token := actualReq.Token

	userSession, err := DbRefreshUserSession(dataSource.DB, token, serviceId)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbRefreshUserSession failed: ", err.Error())
		serverErr := foundation.NewErrorWithMessage(foundation.RevokeTokenError, err.Error())
		return nil, serverErr
	}

	response := new(RefreshLoginTokenResponse)
	if userSession.ExpiredTime != nil {
		response.ExpiredTime = userSession.ExpiredTime.Unix()
	}
	if userSession.CreatedDate.Unix() != 0 {
		response.CreatedDate = userSession.CreatedDate.Unix()
	}

	if userSession.LastModifiedDate.Unix() != 0 {
		response.LastModifiedDate = userSession.LastModifiedDate.Unix()
	}
	response.UserId = userSession.UserId

	if dataSource.Config.ResponseDelayInterval > 0 {
		time.Sleep(time.Duration(dataSource.Config.ResponseDelayInterval) * time.Millisecond)
	}

	return response, nil
}

func RevokeLoginToken(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage, serviceId int64) (*RevokeLoginTokenResponse, *foundation.ServerError) {
	actualReq := req.Request.Data.(*RevokeLoginTokenRequest)

	token := actualReq.Token

	userSession, err := DbRevokeUserSession(dataSource.DB, token, serviceId)

	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbRevokeUserSession failed: ", err.Error())
		serverErr := foundation.NewErrorWithMessage(foundation.RevokeTokenError, err.Error())
		return nil, serverErr
	}

	response := new(RevokeLoginTokenResponse)
	if userSession.ExpiredTime != nil {
		response.ExpiredTime = userSession.ExpiredTime.Unix()
	}
	response.UserId = userSession.UserId

	if dataSource.Config.ResponseDelayInterval > 0 {
		time.Sleep(time.Duration(dataSource.Config.ResponseDelayInterval) * time.Millisecond)
	}
	return response, nil
}

func RequestNonRefreshableLoginToken(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage, serviceId int64) (*NonRefreshableLoginTokenResponse, *foundation.ServerError) {

	type jwtUserInfo struct {
		ClientInfo     string
		ServiceId      int64
		IsPaymentToken bool
	}

	var loginReq *NonRefreshableLoginTokenRequest = req.Request.Data.(*NonRefreshableLoginTokenRequest)
	userInfo := &jwtUserInfo{ClientInfo: loginReq.UserId, IsPaymentToken: true}

	userInfoData, err := json.Marshal(userInfo)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.RequestMalformat, "User data invalid")
	}

	if loginReq.Duration == 0 {
		loginReq.Duration = 300
	}

	loginToken, expireTime, err := GenerateJwtToken("EurusAuth", string(userInfoData), true, loginReq.Duration)
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.RequestMalformat, err.Error())
	}

	err = DbInsertUserSession(dataSource.DB, loginToken, serviceId, expireTime, loginReq.UserId, loginReq.TokenType)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("DbInsertUserSession failed: ", err.Error())
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}

	resData := new(NonRefreshableLoginTokenResponse)
	resData.Token = loginToken
	resData.Type = loginReq.TokenType
	resData.ExpiredTime = expireTime
	resData.SetUserId(loginReq.UserId)
	resData.CreatedDate = time.Now().Unix()
	resData.LastModifiedDate = time.Now().Unix()

	if dataSource.Config.ResponseDelayInterval > 0 {
		time.Sleep(time.Duration(dataSource.Config.ResponseDelayInterval) * time.Millisecond)
	}

	return resData, nil
}

func VerifyInputSign(dataSource *AuthDataSource, req *ws_message.MasterRequestMessage) (*VerifySignResponse, *foundation.ServerError) {
	verifyReq := req.Request.Data.(*VerifySignRequest)
	authService, ok := dataSource.ServiceInfo[verifyReq.ServiceId]
	if !ok {
		return nil, foundation.NewErrorWithMessage(foundation.ServerIdUnmatch, "Service id not found")
	}

	isSuccess, _ := crypto.VerifyRSASignFromBase64(authService.PubKey, string(verifyReq.Data), verifyReq.Sign)
	res := new(VerifySignResponse)
	res.IsSuccess = isSuccess
	return res, nil
}
