package conf

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/request"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/auth_base"

	"fmt"
	"net/http"
)

///Fill up the req or create an error response when error occurs
func RequestToModelHandler(httpReq *http.Request, req request.IRequest) response.IResponse {

	err := api.HttpRequestToModel(httpReq, req, false)

	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		return response.CreateErrorResponse(req, foundation.InternalServerError, errStr)
	}

	obj := httpReq.Context().Value("loginToken")
	if obj != nil {
		// fmt.Println("httpReq.Context().Value().(network.ILoginToken)", httpReq.Context().Value("loginToken").(network.ILoginToken))
		loginToken := obj.(auth_base.ILoginToken)
		req.SetLoginToken(loginToken)
	}
	return nil
}
