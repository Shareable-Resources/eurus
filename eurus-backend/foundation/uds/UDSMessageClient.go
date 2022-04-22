package uds

import (
	"encoding/binary"
	"encoding/json"
	"eurus-backend/foundation/ws/ws_message"
	"net"

	"github.com/sirupsen/logrus"
)

type UDSMessageClient struct {
	Conn   net.Conn
	Logger *logrus.Logger
	// messageReceivedHandler func(*UDSMessageClient, *ws.ResponseMessage, error)
	messageReceivedChan chan (*MessageFrame)
	CloseChan           chan (bool)
	manualCloseConn     bool
}

func NewUDSMessageClient(logger *logrus.Logger, messageReceivedChan chan (*MessageFrame)) *UDSMessageClient {
	client := new(UDSMessageClient)
	client.Logger = logger
	client.messageReceivedChan = messageReceivedChan
	client.CloseChan = make(chan bool)
	return client
}

func (me *UDSMessageClient) Connect(socketPath string) error {
	var err error
	me.Conn, err = net.Dial("unix", socketPath)
	if err != nil {
		return err
	}

	// go func() {
	// 	for {
	// 		messageFrame := <-me.messageReceivedChan
	// 		if messageFrame.Error != nil {
	// 			me.messageReceivedHandler(me, nil, messageFrame.Error)
	// 			if messageFrame.Error == io.EOF {
	// 				me.Logger.Debugln("UDS client connection closed")
	// 				break
	// 			}
	// 		} else {
	// 			me.messageReceivedHandler(me, messageFrame.ResponseMessage, messageFrame.Error)
	// 		}
	// 	}
	// }()

	go func() {
		err = readConnection(me.Conn, me.messageReceivedChan, false, me.Logger)
		if err != nil {
			if !me.manualCloseConn {
				me.Logger.Errorln("Connection error: ", err)
			}
			me.CloseChan <- true
		}
	}()

	return nil
}

func (me *UDSMessageClient) SendMessage(reqMessage *ws_message.MasterRequestMessage) error {
	data, err := json.Marshal(reqMessage)
	if err != nil {
		return err
	}

	var header []byte = make([]byte, 8)
	header[0] = 0xB
	header[1] = 0xa
	header[2] = 0xb
	header[3] = 0xd

	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	_, err = me.Conn.Write(header)
	if err != nil {
		return err
	}
	_, err = me.Conn.Write(data)
	return err
}

func (me *UDSMessageClient) Close() {
	me.manualCloseConn = true

	me.Conn.Close()
}
