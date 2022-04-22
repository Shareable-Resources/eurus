package auth

///Fill up the req or create an error response when error occurs
// func RequestToModelHandler(httpReq *http.Request, req request.IRequest) response.IResponse {
// 	err := api.HttpRequestToModel(httpReq, req, false)

// 	if err != nil {
// 		errStr := fmt.Sprintf("%v", err)
// 		return response.CreateErrorResponse(req, foundation.InternalServerError, errStr)
// 	}
// 	return nil
// }
