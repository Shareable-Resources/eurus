package ws

import (
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/ws/ws_message"
	"time"

	"github.com/gorilla/websocket"
)

type RequestHandlerInfo struct {
	Handler func(*WebSocketServer, *WebSocketConnection, *ws_message.MasterRequestMessage) error
}

type WebSocketMessageDispatcher struct { //implements server.IWebSocketServerListener
	requestHandlerMap map[string]*RequestHandlerInfo
	authFuncName      string
}

func NewWebSocketMessageDispatcher(authFuncName string) *WebSocketMessageDispatcher {
	dispatcher := new(WebSocketMessageDispatcher)
	dispatcher.requestHandlerMap = make(map[string]*RequestHandlerInfo)
	dispatcher.authFuncName = authFuncName
	return dispatcher
}

func (me *WebSocketMessageDispatcher) AddMessageHandler(messageName string, handlerInfo *RequestHandlerInfo) {
	me.requestHandlerMap[messageName] = handlerInfo
}

func (me *WebSocketMessageDispatcher) OnIncomingMessage(server *WebSocketServer, conn *WebSocketConnection, message []byte) {
	me.parseRequest(server, string(message), conn)
}

func (me *WebSocketMessageDispatcher) OnConnectionClosed(server *WebSocketServer, conn *WebSocketConnection) {

}
func (me *WebSocketMessageDispatcher) OnConnectionOpend(server *WebSocketServer, conn *WebSocketConnection) {

}

func (me *WebSocketMessageDispatcher) parseRequest(server *WebSocketServer, message string, conn *WebSocketConnection) {
	var masterRequest = ws_message.MasterRequestMessage{}
	err := json.Unmarshal([]byte(message), &masterRequest)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to parse request: ", err, " message: ", message)
		ResponseError(conn, nil, foundation.RequestMalformat, err.Error())
		return
	}

	var request = ws_message.RequestMessage{}
	err = json.Unmarshal([]byte(masterRequest.Message), &request)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to unmarshal request: ", err, " message: ", message)
		ResponseError(conn, &request, foundation.RequestMalformat, err.Error())
		return
	}

	handlerInfo, ok := me.requestHandlerMap[request.MethodName]
	if !ok {
		ResponseErrorCode(conn, &request, foundation.MethodNotFound)
		return
	}

	masterRequest.Request = &request
	if !conn.IsAuthorized && request.MethodName != me.authFuncName {
		ResponseErrorCode(conn, &request, foundation.UnauthorizedAccess)
		return
	}

	err = handlerInfo.Handler(server, conn, &masterRequest)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Handler error: ", request.MethodName, " error: ", err.Error(), " nonce: ", request.Nonce)
		ResponseError(conn, &request, foundation.InternalServerError, "Internal server error")
	}
}

func Response(conn *WebSocketConnection, req *ws_message.RequestMessage, resData interface{}) {
	var code foundation.ServerReturnCode = foundation.Success

	ResponseWithCode(conn, req, code, code.String(), resData)
}

func ResponseWithCode(conn *WebSocketConnection, req *ws_message.RequestMessage, code foundation.ServerReturnCode, message string, resData interface{}) {
	response := ws_message.ResponseMessage{
		ResponseBase: response.ResponseBase{
			ReturnCode: int64(code),
			Message:    message,
			Nonce:      req.Nonce,
			Data:       resData,
		},
		MethodName: req.MethodName,
		Timestamp:  int64(time.Now().UnixNano() / 1000000.0),
	}
	sendResponse(conn, &response)
}

func ResponseSuccess(conn *WebSocketConnection, req *ws_message.RequestMessage) {
	ResponseErrorCode(conn, req, foundation.Success)
}

func ResponseErrorCode(conn *WebSocketConnection, req *ws_message.RequestMessage, returnCode foundation.ServerReturnCode) {
	ResponseError(conn, req, returnCode, returnCode.String())
}

func ResponseError(conn *WebSocketConnection, req *ws_message.RequestMessage, returnCode foundation.ServerReturnCode, message string) {
	var messageType, nonce string

	if req != nil {
		messageType = req.MethodName
		nonce = req.Nonce
	}

	response := ws_message.ResponseMessage{
		ResponseBase: response.ResponseBase{
			Nonce:      nonce,
			ReturnCode: int64(returnCode),
			Message:    message,
		},
		MethodName: messageType,
		Timestamp:  0,
	}

	timestamp := time.Now().UnixNano() / 1000000.0
	response.Timestamp = timestamp

	sendResponse(conn, &response)
}

func sendResponse(conn *WebSocketConnection, response *ws_message.ResponseMessage) {

	writer, err := conn.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Cannot get writer from websocket connection. Error: ", err.Error(), " Session ID: ", conn.SessionId,
			" service id: ", conn.ServiceId, " nonce: ", response.Nonce, " message type: ", response.MethodName)
		return
	}
	defer writer.Close()
	data, _ := json.Marshal(response)
	index := 0
	for {

		length, err := writer.Write(data)
		if err != nil {
			log.GetLogger(log.Name.Root).Error("Cannot write to websocket. Error: ", err.Error(), " Session ID: ", conn.SessionId,
				" service id: ", conn.ServiceId, " nonce: ", response.Nonce, " message type: ", response.MethodName)
			return
		}

		if length < len(data) {
			data = data[index:length]
		} else {
			break
		}
		index += length
	}
}
