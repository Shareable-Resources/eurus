package merchant_admin

import (
	"encoding/json"
	"eurus-backend/admin_service/merchant_common"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
)

type MerchantAdminProcessor struct {
	smartContractProcessor *MerchantAdminSCProcessor
	dbProcessor            *MerchantAdminDBProcessor
	config                 *MerchantAdminServerConfig
}

func NewMerchantAdminProcessor(dbProcessor *MerchantAdminDBProcessor, scProcessor *MerchantAdminSCProcessor, config *MerchantAdminServerConfig) *MerchantAdminProcessor {
	processor := new(MerchantAdminProcessor)
	processor.smartContractProcessor = scProcessor
	processor.dbProcessor = dbProcessor
	processor.config = config
	return processor
}

func (me *MerchantAdminProcessor) GetRefundRequestList(req *GetMerchantRefundRequest) *response.ResponseBase {
	if req.PageNum == 0 {
		req.PageNum = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var requestStatus *merchant_common.RefundRequestStatus
	if req.Status != -1 {
		requestStatus = (*merchant_common.RefundRequestStatus)(&req.Status)
	} else {
		requestStatus = nil
	}

	loginToken, err := me.UnmarshalMerchantLoginToken(req.GetLoginToken().GetUserId())
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unauthroized access: ", err.Error(), " ", req.GetLoginToken().GetUserId())
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, err.Error())
	}

	count, list, err := me.dbProcessor.QueryRefundRequestList(loginToken.MerchantId, requestStatus, req.PageNum, req.PageSize)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to query refund request list: ", err.Error(), " merchant id: ", loginToken.MerchantId)
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	res := new(GetMerchantRefundResponse)
	res.RecordCount = count
	res.List = list

	return response.CreateSuccessResponse(req, res)
}

func (me *MerchantAdminProcessor) processRefundRequest(req *RefundRequest) *response.ResponseBase {
	loginToken, err := me.UnmarshalMerchantLoginToken(req.LoginToken.GetUserId())
	if err != nil {
		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, err.Error())
	}

	if req.Answer == merchant_common.RefundAccepted {
		if req.RefundTransHash == "" {
			return response.CreateErrorResponse(req, foundation.InvalidArgument, "Missing refund trans hash")
		}
	}

	err = me.dbProcessor.UpdateRefundRequest(loginToken.MerchantId, loginToken.OperatorId, req)
	if err != nil {
		return response.CreateErrorResponse(req, foundation.DatabaseError, "Unable to update refund request: "+err.Error())
	}

	return response.CreateSuccessResponse(req, nil)
}

func (me *MerchantAdminProcessor) DummyLogin(authClient auth_base.IAuth, req *request.RequestBase) *response.ResponseBase {
	merchantLoginToken := new(MerchantLoginToken)
	merchantLoginToken.MerchantId = 1
	merchantLoginToken.OperatorId = 1

	tokenJson, _ := json.Marshal(merchantLoginToken)

	loginToken, err := authClient.GenerateLoginToken(string(tokenJson))
	if err != nil {
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}
	var responseMap map[string]interface{} = make(map[string]interface{}, 0)
	responseMap["loginToken"] = loginToken.GetToken()

	return response.CreateSuccessResponse(req, responseMap)
}

func (me *MerchantAdminProcessor) UnmarshalMerchantLoginToken(token string) (*MerchantLoginToken, error) {
	var loginToken *MerchantLoginToken = new(MerchantLoginToken)
	err := json.Unmarshal([]byte(token), &loginToken)
	if err != nil {
		return nil, err
	}
	return loginToken, nil
}

func (me *MerchantAdminProcessor) processLogin(authClient auth_base.IAuth, req *MerchantLoginRequest) *response.ResponseBase {
	merchantAdmin, err := me.dbProcessor.getMerchantByUserName(req.MerchantId, req.UserName)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to getMerchantByUserName. Error: ", err.Error(), " nonce: ", req.Nonce)
		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
	}

	if merchantAdmin.MerchantId == 0 {
		return response.CreateErrorResponse(req, foundation.UserNotFound, "Merchant admin not found")
	}

	if merchantAdmin.PasswordHash != req.PasswordHash {
		return response.CreateErrorResponse(req, foundation.UserNotFound, "Merchant admin not found")
	}

	merchantLoginToken := new(MerchantLoginToken)
	merchantLoginToken.MerchantId = req.MerchantId
	merchantLoginToken.OperatorId = merchantAdmin.OperatorId
	merchantLoginToken.AccountState = merchantAdmin.Status

	tokenJson, _ := json.Marshal(merchantLoginToken)

	loginToken, err := authClient.GenerateLoginToken(string(tokenJson))
	if err != nil {
		return response.CreateErrorResponse(req, foundation.InternalServerError, err.Error())
	}

	res := new(MerchantLoginResponse)
	res.LoginToken = loginToken.GetToken()
	res.AccountState = merchantAdmin.Status

	return response.CreateSuccessResponse(req, res)
}

func (me *MerchantAdminProcessor) ProcessRefreshLoginToken(authClient auth_base.IAuth, req *MerchantRefreshLoginTokenRequest) *response.ResponseBase {
	loginToken, err := authClient.RefreshLoginToken(req.LoginToken.GetToken())
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to refresh login token: ", err)
		return response.CreateErrorResponse(req, err.GetReturnCode(), err.Error())
	}

	res := new(MerchantRefreshLoginTokenResponse)
	res.ExpireTime = loginToken.GetExpiredTime()
	res.CreatedDate = loginToken.GetCreatedDate()
	res.LastModifiedDate = loginToken.GetLastModifiedDate()
	res.Token = loginToken.GetToken()

	return response.CreateSuccessResponse(req, res)
}

// func (me *MerchantAdminProcessor) processChangePassword(authClient network.IAuth, req *MerchantChangePasswordRequest) *response.ResponseBase {

// 	merchantLoginToken := new(MerchantLoginToken)
// 	err := json.Unmarshal([]byte(req.LoginToken.GetUserId()), &merchantLoginToken)
// 	if err != nil {
// 		log.GetLogger(log.Name.Root).Errorln("Unable to unmarshal login token: ", err, " nonce: ", req.Nonce)
// 		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Invalid login token")
// 	}

// 	merchantAdmin, err := me.dbProcessor.getMerchantByUserName(req.MerchantId, req.UserName)
// 	if err != nil {
// 		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
// 	}

// 	if merchantAdmin.OperatorId == 0 {
// 		return response.CreateErrorResponse(req, foundation.UserNotFound, "User not found")
// 	}

// 	if merchantAdmin.OperatorId != merchantLoginToken.OperatorId {
// 		return response.CreateErrorResponse(req, foundation.UnauthorizedAccess, "Unauthorize access")
// 	}

// 	err = me.dbProcessor.updateMerchantPassword(req.MerchantId, req.UserName, req.NewPasswordHash)
// 	if err != nil {
// 		log.GetLogger(log.Name.Root).Errorln("Unable to update merchant password: ", err)
// 		return response.CreateErrorResponse(req, foundation.DatabaseError, err.Error())
// 	}
// 	return response.CreateSuccessResponse(req, nil)
// }
