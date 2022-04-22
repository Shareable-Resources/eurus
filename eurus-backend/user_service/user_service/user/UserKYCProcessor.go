package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/auth_base"
	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
)

func GetKYCCountryListFromKYCServer(server *UserServer) ([]*kyc_model.KYCCountryCode, error) {
	var err error
	req := NewKYCCRequest()
	req.RequestPath = KYCServerPath + RootPath + GetKYCCountryListPath
	kycRes := kyc_model.NewResponseGetKYCCountryList()
	reqRes := api.NewRequestResponse(req, kycRes)
	urlObj := server.Config.KYCServerUrl + req.RequestPath
	_, err = server.SendApiRequest(urlObj, reqRes)
	if err != nil {
		return nil, err
	}
	var data = kycRes.Data
	return data, err
}

func CreateKYCStatusFromKYCServer(server *UserServer, req *kyc_model.RequestCreateKYCStatus) (*response.ResponseBase, error) {
	var err error
	req.RequestPath = KYCServerPath + RootPath + CreateKYCStatusPath
	kycRes := new(response.ResponseBase)
	reqRes := api.NewRequestResponse(req, kycRes)
	urlObj := server.Config.KYCServerUrl + req.RequestPath
	_, err = server.SendApiRequest(urlObj, reqRes)
	//Direct Error Message from KYC Server
	if err != nil {
		return nil, err
	}
	if kycRes.ReturnCode != int64(foundation.Success) {
		return kycRes, errors.New(kycRes.GetMessage())
	}
	/*

	 */
	return kycRes, nil
}

func SubmitKYCApproval(server *UserServer, req *kyc_model.RequestSubmitKYCApproval) (*response.ResponseBase, error) {
	var err error
	req.RequestPath = KYCServerPath + RootPath + SubmitKYCApprovalPath
	kycRes := new(response.ResponseBase)
	reqRes := api.NewRequestResponse(req, kycRes)
	urlObj := server.Config.KYCServerUrl + req.RequestPath
	_, err = server.SendApiRequest(urlObj, reqRes)
	//Direct Error Message from KYC Server
	if err != nil {
		return kycRes, err
	}
	if kycRes.ReturnCode != int64(foundation.Success) {
		return kycRes, errors.New(kycRes.GetMessage())
	}
	return kycRes, nil
}

func GetKYCStatusOfUser(server *UserServer, req *kyc_model.RequestGetKYCStatusOfUser) (*response.ResponseBase, error) {
	var err error
	req.RequestPath += "/" + strconv.FormatUint(req.UserId, 10)
	kycRes := new(response.ResponseBase)
	reqRes := api.NewRequestResponse(req, kycRes)
	urlObj := server.Config.KYCServerUrl + req.RequestPath
	_, err = server.SendApiRequest(urlObj, reqRes)
	//Direct Error Message from KYC Server
	if err != nil {
		return kycRes, err
	}

	if kycRes.ReturnCode != int64(foundation.Success) {
		return kycRes, errors.New(kycRes.GetMessage())
	}
	/*
		//Convert model to jsonString
		jsonString, _ := json.Marshal(res.GetData())
		fmt.Println(string(jsonString))
		// convert json to struct
		parsedModel := []kyc_model.UserKYCStatus{}
		json.Unmarshal(jsonString, &parsedModel)*/
	return kycRes, nil
}

func ProcessKYCDocument(server *UserServer, multiPartReader *multipart.Reader, loginToken auth_base.ILoginToken) (request.IRequest, *response.ResponseBase) {
	emptyReqObj := new(request.RequestBase)
	part, err := multiPartReader.NextPart()
	if err != nil {

		return emptyReqObj, response.CreateErrorResponse(emptyReqObj, foundation.BadRequest, "Cannot read form data part: "+err.Error())
	}

	if part.FormName() != "requestJson" {
		return emptyReqObj, response.CreateErrorResponse(emptyReqObj, foundation.BadRequest, "Form name requestJson expected")
	}

	formData, err := io.ReadAll(part)
	if err != nil {
		return emptyReqObj, response.CreateErrorResponse(emptyReqObj, foundation.BadRequest, "Read form data requestJson error: "+err.Error())
	}

	reqObj := new(kyc_model.RequestSubmitKYCDocument)
	reqObj.LoginToken = loginToken
	err = json.Unmarshal(formData, reqObj)
	if err != nil {
		return reqObj, response.CreateErrorResponse(emptyReqObj, foundation.RequestMalformat, err.Error())
	}

	submitImageReq := new(kyc_model.KYCSubmitImage)
	submitImageReq.KYCStatusId = reqObj.UserKYCStatusId
	userId, err := GetUserIdFromLoginToken(reqObj.LoginToken)
	if err != nil {
		return reqObj, response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, err.Error())
	}

	submitImageReq.UserId = userId
	submitImageReq.ImageType = reqObj.ImageType
	submitImageReq.FileExtension = reqObj.FileExtension

	reqJson, _ := json.Marshal(submitImageReq)

	buffer := new(bytes.Buffer)
	partWriter := multipart.NewWriter(buffer)
	partWriter.SetBoundary("20210729")
	defer partWriter.Close()

	fieldWriter, _ := partWriter.CreateFormField("kycSubmitImage")
	_, err = fieldWriter.Write(reqJson)
	if err != nil {
		return reqObj, response.CreateErrorResponse(reqObj, foundation.LoginTokenInvalid, err.Error())
	}
	fileName := fmt.Sprintf("kyc_%d_%d_%d.%s", submitImageReq.KYCStatusId, submitImageReq.UserId, submitImageReq.ImageType, submitImageReq.FileExtension)
	fileWriter, _ := partWriter.CreateFormFile("image", fileName)
	part2, err := multiPartReader.NextPart()
	if err != nil {
		return reqObj, response.CreateErrorResponse(reqObj, foundation.InternalServerError, "Cannot read image form data part: "+err.Error())
	}

	n, err := io.Copy(fileWriter, part2)
	if err != nil {
		return reqObj, response.CreateErrorResponse(reqObj, foundation.InternalServerError, "Cannot copy image data: "+err.Error())

	}
	if n == 0 {
		return reqObj, response.CreateErrorResponse(reqObj, foundation.InvalidArgument, "Empty image")
	}

	partWriter.Close()
	var headers map[string]interface{} = make(map[string]interface{})
	headers["Boundary"] = partWriter.Boundary()
	err = server.kycMq.PublishWithContentType("", "multipart/form-data", buffer.Bytes(), headers)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to publish data to MQ. Error: ", err)
		return reqObj, response.CreateErrorResponse(reqObj, foundation.InternalServerError, err.Error())

	}
	return reqObj, response.CreateSuccessResponse(reqObj, nil)
}
