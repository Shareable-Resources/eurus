package test

import (
	"bytes"
	"encoding/json"

	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"testing"
)

var testData TestData = TestData{ResponseData: "123"}

func TestInitBaseRequest(t *testing.T) {
	requestBase := request.InitRequestBase()

	requestBase.Nonce = "requestNonce1"
	strResult := requestBase.GetNonce()
	if strResult != "requestNonce1" {
		t.Errorf("GetNonce() failed, got %s", strResult)
	} else {
		t.Logf("GetNonce() succeeded, got %s", strResult)
	}
	requestBase.SetNonce("requestNonce2")
	strResult = requestBase.GetNonce()
	if strResult != "requestNonce2" {
		t.Errorf("SetNonce() failed, got %s", strResult)
	} else {
		t.Logf("SetNonce() succeeded, got %s", strResult)
	}
}

func TestInitBaseResponse(t *testing.T) {
	responseBase := response.InitResponseBase()
	responseBase.MethodName = "responseMethodName"
	strResult := responseBase.GetMethodName()
	if strResult != "responseMethodName" {
		t.Errorf("GetMethodName() failed, got %s", strResult)
	} else {
		t.Logf("GetMethodName() succeeded, got %s", strResult)
	}
	responseBase.ReturnCode = 123
	intResult := responseBase.GetReturnCode()
	if intResult != 123 {
		t.Errorf("GetReturnCode() failed, got %v", intResult)
	} else {
		t.Logf("GetReturnCode() succeeded, got %v", intResult)
	}
	responseBase.Nonce = "responseNonce"
	strResult = responseBase.GetNonce()
	if strResult != "responseNonce" {
		t.Errorf("GetNonce() failed, got %s", strResult)
	} else {
		t.Logf("GetNonce() succeeded, got %s", strResult)
	}

	responseBase.Data = testData
	byteResult, err := json.Marshal(responseBase.GetData())
	expectedResult, _ := json.Marshal(testData)
	if err != nil {
		t.Errorf("%x", err)
	} else if bytes.Compare(byteResult, expectedResult) != 0 {
		t.Errorf("GetData() failed, got %s", string(byteResult))
	} else {
		t.Logf("GetData() succeeded, got %s", string(byteResult))
	}
}

func TestInitRequestResponseByBase(t *testing.T) {
	responseBase := response.InitResponseBase()
	responseBase.MethodName = "responseMethodName"
	responseBase.ReturnCode = 123
	responseBase.Nonce = "responseNonce"
	responseBase.Data = testData

	requestBase := request.InitRequestBase()
	requestBase.MethodName = "requestMethodName"
	requestBase.Nonce = "requestNonce1"

	reqRes := api.NewRequestResponse(requestBase, responseBase)
	strResult, _ := reqRes.RequestToJson()
	expectedResult, _ := json.Marshal(requestBase)
	if string(expectedResult) != strResult {
		t.Errorf("RequestToJson() failed, got %s", strResult)
	} else {
		t.Logf("RequestToJson() succeeded, got %s", strResult)

	}
	strResult, _ = reqRes.ResponseToJson()
	expectedResult, _ = json.Marshal(responseBase)
	if string(expectedResult) != strResult {
		t.Errorf("ResponseToJson() failed, got %s", strResult)
	} else {
		t.Logf("ResponseToJson() succeeded, got %s", strResult)

	}
}
