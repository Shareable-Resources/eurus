package response

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"
)

type ResponseBase struct { // implements response.IResponse
	ReturnCode int64       `json:"returnCode"`
	Message    string      `json:"message"`
	Nonce      string      `json:"nonce"`
	Data       interface{} `json:"data"`
}

func (me *ResponseBase) IsInterfaceNil() bool {
	if nil == me {
		return true
	}
	return false
}

func (me *ResponseBase) GetReturnCode() int64 {
	return me.ReturnCode
}

func (me *ResponseBase) GetMessage() string {
	return me.Message
}

func (me *ResponseBase) GetNonce() string {
	return me.Nonce
}

func (me *ResponseBase) GetData() interface{} {
	return me.Data
}

func InitResponseBase() *ResponseBase {
	res := &ResponseBase{}
	return res
}

func CreateErrorResponse(req request.IRequest, code foundation.ServerReturnCode, message string) *ResponseBase {
	res := new(ResponseBase)
	res.Nonce = req.GetNonce()
	res.ReturnCode = int64(code)
	res.Message = message
	return res
}

func CreateSuccessResponse(req request.IRequest, data interface{}) *ResponseBase {
	res := new(ResponseBase)
	res.Nonce = req.GetNonce()
	var code foundation.ServerReturnCode = foundation.Success
	res.ReturnCode = int64(code)
	res.Message = code.String()
	res.Data = data
	return res
}

func CreateSuccessResponseWithMessage(req request.IRequest, data interface{}, message string) *ResponseBase {
	res := new(ResponseBase)
	res.Nonce = req.GetNonce()
	var code foundation.ServerReturnCode = foundation.Success
	res.ReturnCode = int64(code)
	res.Message = message
	res.Data = data
	return res
}
