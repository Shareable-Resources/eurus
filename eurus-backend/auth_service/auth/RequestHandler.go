package auth

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/ws"
	"eurus-backend/foundation/ws/ws_message"
	"fmt"
)

func AuthenticateHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {
	authServer := server.ActualServer.(*AuthServer)
	request := req.Request
	innerRequest := request.Data.(*AuthenticateRequest)

	if req.Sign == "" {
		ws.ResponseErrorCode(conn, request, foundation.InvalidSignature)
		return nil
	} else {

		isValid := verifySignature(authServer.DataSource, innerRequest.ServiceId, req.Message, req.Sign, req.Timestamp)
		if !isValid {
			ws.ResponseErrorCode(conn, request, foundation.SignMatchError)
			return nil
		}
	}

	if conn.IsAuthorized {
		response := &AuthenticateResponse{
			SessionId: conn.SessionId,
		}

		ws.Response(conn, request, response)
		return nil
	}

	loginReq := &RequestLoginTokenRequest{}
	userId := fmt.Sprintf(`{"serviceId": %d, "sessionId":%d}`, innerRequest.ServiceId, conn.SessionId)
	loginReq.UserId = userId
	masterLoginReq, err := ws_message.CreateMasterRequestMessage(loginReq, req.Request.Method)
	if err != nil {
		ws.ResponseError(conn, request, foundation.InternalServerError, err.Error())
		return nil
	}

	res, serverError := RequestLoginToken(authServer.DataSource, masterLoginReq, 0)
	if serverError != nil {
		ws.ResponseError(conn, request, serverError.ReturnCode, serverError.Message)
		return nil
	}

	isSuccess := server.AuthorizeSession(conn)
	if !isSuccess {
		ws.ResponseErrorCode(conn, request, foundation.MethodNotFound)
		return nil
	}

	conn.ServiceId = innerRequest.ServiceId

	addInfoStr := fmt.Sprintf(`{"configServerIp":"%s", "configServerPort":%d}`, server.ServerBase.ServerConfig.ConfigServerIP, server.ServerConfig.ConfigServerPort)

	response := &AuthenticateResponse{
		SessionId:      conn.SessionId,
		Token:          res.Token,
		AdditionalInfo: addInfoStr,
	}

	ws.Response(conn, request, response)
	return nil
}

func RequestLoginTokenHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {
	authServer := server.ActualServer.(*AuthServer)

	data, err := RequestLoginToken(authServer.DataSource, req, conn.ServiceId)
	if err != nil {
		ws.ResponseError(conn, req.Request, err.ReturnCode, err.Error())
		return nil
	}

	ws.Response(conn, req.Request, data)

	return nil
}

func VerifyLoginTokenHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {

	authServer := server.ActualServer.(*AuthServer)

	data, err := VerifyLoginToken(authServer.DataSource, req, conn.ServiceId)
	if err != nil {
		ws.ResponseError(conn, req.Request, err.ReturnCode, err.Error())
		return nil
	}

	ws.Response(conn, req.Request, data)

	return nil
}

//RefreshLoginTokenHandler

func RefreshLoginTokenHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {

	authServer := server.ActualServer.(*AuthServer)

	data, err := RefreshLoginToken(authServer.DataSource, req, conn.ServiceId)
	if err != nil {
		ws.ResponseError(conn, req.Request, err.ReturnCode, err.Error())
		return nil
	}

	ws.Response(conn, req.Request, data)

	return nil
}

func RevokeLoginTokenHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {
	authServer := server.ActualServer.(*AuthServer)
	data, err := RevokeLoginToken(authServer.DataSource, req, conn.ServiceId)
	if err != nil {
		ws.ResponseError(conn, req.Request, err.ReturnCode, err.Error())
		return nil
	}
	ws.Response(conn, req.Request, data)
	return nil
}

func RequestNonRefreshableLoginTokenHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {
	//request := req.Request
	authServer := server.ActualServer.(*AuthServer)
	loginReq := &NonRefreshableLoginTokenRequest{}
	masterLoginReq, err := ws_message.CreateMasterRequestMessage(loginReq, req.Request.Method)
	masterLoginReq.Request.Data = req.Request.Data
	if err != nil {
		ws.ResponseError(conn, req.Request, foundation.InternalServerError, err.Error())
		return nil
	}
	data, err2 := RequestNonRefreshableLoginToken(authServer.DataSource, masterLoginReq, conn.ServiceId)
	if err2 != nil {
		ws.ResponseError(conn, req.Request, err2.ReturnCode, err.Error())
		return nil
	}

	ws.Response(conn, req.Request, data)
	return nil
}

func VerifySignHandler(server *ws.WebSocketServer, conn *ws.WebSocketConnection, req *ws_message.MasterRequestMessage) error {
	authServer := server.ActualServer.(*AuthServer)
	res, err := VerifyInputSign(authServer.DataSource, req)
	if err != nil {
		ws.ResponseError(conn, req.Request, err.ReturnCode, err.Error())
		return nil
	}
	ws.Response(conn, req.Request, res)
	return nil
}

// func CreateSuccessResponse(data interface{}, req *ws_message.MasterRequestMessage) *ws.ResponseMessage {
// 	res := new(ws.ResponseMessage)
// 	res.Nonce = req.Request.GetNonce()
// 	res.MethodName = req.Request.MethodName
// 	var code foundation.ServerReturnCode = foundation.Success
// 	res.ReturnCode = int(code)
// 	res.Message = code.String()
// 	res.Timestamp = time.Now().UnixNano() / 1000000.0
// 	res.Data = data
// 	return res
// }
