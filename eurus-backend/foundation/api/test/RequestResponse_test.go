package test

import (
	"encoding/json"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"

	"testing"
)

func (me *TestResponse) GetData() interface{} {
	data := &TestData{ResponseData: "on9"}
	return data
}

func (me *TestResponse) GetMethodName() string {
	me.MethodName = "test123"
	return me.MethodName
}

func (me *TestResponse) GetNonce() string {
	me.Nonce = "test456"
	return me.Nonce
}

func (me *TestResponse) GetReturnCode() int {
	me.ReturnCode = 123
	return me.ReturnCode
}

func (me *TestRequest) GetMethodName() string {
	return "testMethod"
}

func (me *TestRequest) SetNonce(nonce string) {
	me.Nonce = nonce
}

func (me *TestRequest) GetNonce() string {
	return me.Nonce
}

func TestInitRequestResponse(t *testing.T) {
	var testIresponse response.IResponse
	var testResponse TestResponse
	testIresponse = &testResponse
	var testIrequest request.IRequest
	var testRequest TestRequest
	testRequest.MethodName = "789"
	testIrequest = &testRequest

	requestResponse := api.NewRequestResponse(testIrequest, testIresponse)
	result1 := requestResponse.Res.GetMethodName()
	if result1 == "test123" {
		t.Logf("GetMethodName() succeeded, got %s", result1)
	} else {
		t.Errorf("GetMethodName() failed, got %s", result1)
	}
	result1 = requestResponse.Res.GetNonce()
	if result1 == "test456" {
		t.Logf("GetNonce() succeeded, got %s", result1)
	} else {
		t.Errorf("GetNonce() failed, got %s", result1)
	}
	result2 := requestResponse.Res.GetReturnCode()
	if result2 == 123 {
		t.Logf("GetReturnCode() succeeded, got %v", result2)
	} else {
		t.Errorf("GetReturnCode() failed, got %v", result2)
	}

	result1 = requestResponse.Req.GetMethodName()
	if result1 == "testMethod" {
		t.Logf("GetMethodName() succeeded, got %s", result1)
	} else {
		t.Errorf("GetMethodName() failed, got %s", result1)
	}

	requestResponse.Req.SetNonce("test456")
	result1 = requestResponse.Req.GetNonce()
	if result1 == "test456" {
		t.Logf("GetNonce() succeeded, got %s", result1)
	} else {
		t.Errorf("GetNonce() failed, got %s", result1)
	}

	result3 := requestResponse.Res.GetData()
	if result3.(*TestData).ResponseData == "on9" {
		t.Logf("GetData() succeeded, got %s", result3)
	} else {
		t.Errorf("GetData() failed, got %s", result3)

	}
	result1, err := requestResponse.RequestToJson()
	if err != nil {
		t.Errorf("%x", err)
	} else if result1 == "{\"nonce\":\"test456\",\"methodName\":\"789\"}" {
		t.Logf("RequestToJson() succeeded, got %s", result1)
	} else {
		t.Errorf("RequestToJson() failed, got %s", result1)
	}

	jsonStr := "{\"methodName\":\"123\",\"returnCode\":456,\"nonce\":\"789\",\"data\":{\"responseData\":\"Hello\"}}"
	testRequest1 := TestRequest{}
	testResponse1 := TestResponse{}
	testResponse1.Data = TestData{}
	requestResponse1 := api.NewRequestResponse(&testRequest1, &testResponse1)

	err = requestResponse1.ResponseJsonToModel(jsonStr)
	if err != nil {
		t.Errorf("%x", err)
	} else {
		byteStr, _ := json.Marshal(requestResponse1.Res)
		t.Logf("ResponseJsonToModel() succeeded, got %s", string(byteStr))

	}

}
