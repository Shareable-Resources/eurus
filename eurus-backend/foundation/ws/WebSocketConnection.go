package ws

import "github.com/gorilla/websocket"

type WebSocketConnection struct {
	SessionId    int64
	ServiceId    int64
	Conn         *websocket.Conn
	IsAuthorized bool
}
