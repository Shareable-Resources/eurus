package uds

import (
	"errors"
	"eurus-backend/foundation/ws/ws_message"
	"io"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

type UDSMessageServer struct {
	SocketPath     string
	Logger         *logrus.Logger
	MessageHandler chan (*MessageFrame)
	server         net.Listener
}

type MessageFrame struct {
	Conn                         net.Conn
	IsConnectionEstablishedEvent bool
	Frame                        []byte
	MasterRequestMessage         *ws_message.MasterRequestMessage
	ResponseMessage              *ws_message.ResponseMessage
	Error                        error
}

func NewUDSMessageServer(socketPath string, logger *logrus.Logger, messageHandler chan (*MessageFrame)) (*UDSMessageServer, error) {
	server := new(UDSMessageServer)
	if socketPath == "/" || socketPath == "" {
		logger.Errorln("Invalid socket path: ", socketPath)
		return nil, errors.New("Invalid socket path")
	}
	server.SocketPath = socketPath
	server.Logger = logger
	server.MessageHandler = messageHandler
	return server, nil
}

func (me *UDSMessageServer) Listen() error {
	var err error
	if err = os.RemoveAll(me.SocketPath); err != nil {
		me.Logger.Errorln("Unable to remove socket: ", err)
		return err
	}

	me.server, err = net.Listen("unix", me.SocketPath)
	if err != nil {
		me.Logger.Errorln("Unable to open UDS: ", err)
		return err
	}

	var conn net.Conn

	for {
		conn, err = me.server.Accept()
		if err != nil {
			me.Logger.Warn("Unable to accept new connection: ", err)
			continue
		}

		me.Logger.Infoln("Connection accepted: remote IP: ", conn.RemoteAddr().String())

		go func() {
			me.MessageHandler <- &MessageFrame{
				Conn:                         conn,
				IsConnectionEstablishedEvent: true,
			}
			err = readConnection(conn, me.MessageHandler, true, me.Logger)
			if err != nil && err != io.EOF {
				me.Logger.Errorln("Connection error: ", err, " conn remote IP: ", conn.RemoteAddr().String())
			}
		}()
	}

	return nil
}

func (me *UDSMessageServer) Close() {
	me.server.Close()
}

func (me *UDSMessageServer) SendJsonMessage(conn net.Conn, data interface{}) error {
	return SendUDSJsonMessage(conn, data, me.Logger)
}
