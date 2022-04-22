package ws

import (
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type IWebSocketServerListener interface {
	OnIncomingMessage(server *WebSocketServer, conn *WebSocketConnection, message []byte)
	OnConnectionClosed(server *WebSocketServer, conn *WebSocketConnection)
	OnConnectionOpend(server *WebSocketServer, conn *WebSocketConnection)
}

type WebSocketServer struct {
	server.ServerBase
	Listener                IWebSocketServerListener
	upgrader                websocket.Upgrader
	pendingConnectionList   []*pendingConnection
	pendingConnectionMutex  sync.Mutex
	authorizedConnectionMap map[int64]*WebSocketConnection
	closeSignal             chan struct{}
	lastUsedSessionId       int64
	terminalCommandHandler  func(req *server.ControlRequestMessage) (bool, string, error)
}

type pendingConnection struct {
	Conn        *websocket.Conn
	SessionId   int64
	ConnectTime time.Time
}

func NewWebSocketServer(listener IWebSocketServerListener) *WebSocketServer {
	webSocketServer := new(WebSocketServer)
	webSocketServer.ServerConfig = new(server.ServerConfigBase)
	webSocketServer.pendingConnectionList = make([]*pendingConnection, 0)
	webSocketServer.authorizedConnectionMap = make(map[int64]*WebSocketConnection)
	webSocketServer.closeSignal = make(chan struct{})
	webSocketServer.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	webSocketServer.upgrader.Subprotocols = []string{"json"}
	webSocketServer.Listener = listener
	webSocketServer.ActualServer = webSocketServer

	return webSocketServer
}

func (me *WebSocketServer) InitHttpServer(httpConfig network.IHttpConfig) {
	panic("Call InitWebSocketServer function instead")
}

func (me *WebSocketServer) InitWebSocketServer(httpConfig network.IHttpConfig,
	webSocketPath string) error {

	err := me.ServerBase.InitHttpServer(httpConfig)
	if err != nil {
		return err
	}

	me.HttpServer.Router.HandleFunc(webSocketPath, me.WebSocketHandler)

	go me.housekeepPendingConnection()
	return err
}

func (me *WebSocketServer) WebSocketHandler(writer http.ResponseWriter, req *http.Request) {
	conn, err := me.upgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("Unable to upgrade websocket: ", err.Error())
		return
	}
	me.pendingConnectionMutex.Lock()
	me.lastUsedSessionId++

	pendingConn := &pendingConnection{Conn: conn, ConnectTime: time.Now(), SessionId: me.lastUsedSessionId}
	me.pendingConnectionMutex.Unlock()

	me.pendingConnectionList = append(me.pendingConnectionList, pendingConn)
	websocketConn := new(WebSocketConnection)
	websocketConn.Conn = conn
	websocketConn.SessionId = pendingConn.SessionId
	go me.MessageHandler(websocketConn)
}

func (me *WebSocketServer) MessageHandler(conn *WebSocketConnection) {
	defer me.Listener.OnConnectionClosed(me, conn)
	defer conn.Conn.Close()
	var buffer []byte = make([]byte, 4096)
	var doubleBuffer []byte = make([]byte, 0)
	for {
		if _, r, err := conn.Conn.NextReader(); err != nil {
			me.removeConnection(conn)
			return
		} else {
			doubleBuffer = doubleBuffer[0:0]
			var totalLen uint64
			for {
				len, err := r.Read(buffer)
				if err != nil && err == io.EOF {
					if len > 0 {
						totalLen += uint64(len)
						doubleBuffer = append(doubleBuffer[:], buffer[:len]...)
					}
					me.Listener.OnIncomingMessage(me, conn, doubleBuffer[:totalLen])

					break
				} else if err != nil {
					log.GetLogger(log.Name.Root).Error("Connection IP: ", conn.Conn.RemoteAddr().String(), " read error: ", err)
					break
				} else {
					totalLen += uint64(len)
					doubleBuffer = append(doubleBuffer[:], buffer[:len]...)
				}
			}
		}
	}
}

func (me *WebSocketServer) housekeepPendingConnection() {
	for {

		select {
		case <-me.closeSignal:
			break
		case <-time.After(time.Minute):
			me.processHousekeepPendingConnection()
		}
	}
}

func (me *WebSocketServer) processHousekeepPendingConnection() {
	me.pendingConnectionMutex.Lock()
	defer me.pendingConnectionMutex.Unlock()

	currentTime := time.Now()
	var closeIndex int = -1
	for i, pendingConn := range me.pendingConnectionList {
		duration := currentTime.Sub(pendingConn.ConnectTime)
		if duration.Minutes() >= 3.0 {
			pendingConn.Conn.Close()
			closeIndex = i
		}
	}
	if closeIndex >= 0 {
		me.pendingConnectionList = me.pendingConnectionList[:closeIndex+1]
	}

}

func (me *WebSocketServer) removeConnection(conn *WebSocketConnection) {
	if conn.IsAuthorized {
		delete(me.authorizedConnectionMap, conn.SessionId)

	} else {
		me.removePendingConnection(conn.SessionId)
	}
}

func (me *WebSocketServer) InitUDSControlServer(commandLineArgs *server.CommandLineArguments, handler func(req *server.ControlRequestMessage) (bool, string, error)) {

	me.ServerBase.InitUDSControlServer(commandLineArgs, me.TerminalFunction)
}

func (me *WebSocketServer) TerminalFunction(reg *server.ControlRequestMessage) (bool, string, error) {
	var isHandled bool = false
	var output string
	if me.terminalCommandHandler != nil {
		isHandled, output, _ = me.terminalCommandHandler(reg)
	}
	if !isHandled {
		switch reg.MethodName {
		case "displaywscount":
			output = fmt.Sprintf("Pending connection count: %d\r\nAuthroized connection count: %d\r\n",
				len(me.pendingConnectionList), len(me.authorizedConnectionMap))
			return true, output, nil
		case "displaywsdetail":
			output = fmt.Sprintln("Connection details: ")
			for sessionId, conn := range me.authorizedConnectionMap {
				output += fmt.Sprintln("Session ID\tRemote Address\tLocal address")
				if conn.Conn != nil {
					output += fmt.Sprintf("[%d]\t%s\t%s\r\n", sessionId, conn.Conn.RemoteAddr().String(), conn.Conn.LocalAddr().String())
				}
			}
			return true, output, nil
		case "help":
			output = me.displayCommandUsage()
		}
	}
	return false, output, nil
}

func (me *WebSocketServer) displayCommandUsage() string {
	var output string
	output = fmt.Sprintln("DisplayWSCount - Display total number of WS connection")
	output += fmt.Sprintln("DisplayWSDetail - Display connection details")
	return output
}

func (me *WebSocketServer) AuthorizeSession(conn *WebSocketConnection) bool {
	isFound := me.removePendingConnection(conn.SessionId)
	if !isFound {
		return isFound
	}
	conn.IsAuthorized = true
	me.authorizedConnectionMap[conn.SessionId] = conn
	return true
}

func (me *WebSocketServer) removePendingConnection(sessionId int64) bool {
	me.pendingConnectionMutex.Lock()
	defer me.pendingConnectionMutex.Unlock()
	for i, pendingConn := range me.pendingConnectionList {
		if pendingConn.SessionId == sessionId {
			me.pendingConnectionList = append(me.pendingConnectionList[:i], me.pendingConnectionList[i+1:]...)

			return true
		}
	}
	return false
}
