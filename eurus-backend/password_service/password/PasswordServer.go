package password

import (
	go_crypto "crypto"
	"encoding/json"
	"eurus-backend/env"
	"eurus-backend/foundation"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/foundation/uds"
	"eurus-backend/foundation/ws/ws_message"
	"eurus-backend/password_service/password_api"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

type PasswordServer struct {
	server.ServerBase
	Config         *PasswordServerConfig
	udsServer      *uds.UDSMessageServer
	messageHandler chan (*uds.MessageFrame)
}

func NewPasswordServer() *PasswordServer {
	server := new(PasswordServer)
	server.Config = NewPasswordServerConfig()
	server.ServerConfig = &server.Config.ServerConfigBase
	server.messageHandler = make(chan *uds.MessageFrame)
	return server
}

func (me *PasswordServer) InitAll() {
	var err error
	me.udsServer, err = uds.NewUDSMessageServer(me.Config.UDSPath, log.GetLogger(log.Name.Root), me.messageHandler)
	if err != nil {
		panic(err)
	}

	var pw []byte
	for {
		fmt.Print("Input password (screen is hidden): ")
		pw, err = terminal.ReadPassword(0)
		if err != nil {
			if env.Tag != "dev" {
				panic("Running debug mode for input password is not supported. Error: " + err.Error())
			} else {
				pw = []byte{'a', 'b', 'c', 'd', '1', '2', '3', '4'}
			}
		}

		if len(pw) > 0 {
			//Encrypt the password by password client public key
			encrypted, err := crypto.EncryptRAFormat(pw, passwordClientPublicKey, go_crypto.SHA256)
			if err != nil {
				fmt.Println("Unable to encrypt password: ", err.Error())
				continue
			}
			me.Config.Password = encrypted
			break
		}
	}

	go func() {
		err := me.udsServer.Listen()
		if err != nil {
			panic(err)
		}
	}()

	go me.receiveMessage()
}

func (me *PasswordServer) receiveMessage() {
	for {
		messageFrame := <-me.messageHandler
		if messageFrame.Error != nil {
			log.GetLogger(log.Name.Root).Errorln("Socket error received: ", messageFrame.Error, " conn: ", messageFrame.Conn.RemoteAddr().String())
			continue
		} else if messageFrame.IsConnectionEstablishedEvent {
			continue
		} else if messageFrame.MasterRequestMessage != nil {
			req := new(ws_message.RequestMessage)
			err := json.Unmarshal([]byte(messageFrame.MasterRequestMessage.Message), &req)
			if err != nil {

				res := new(ws_message.ResponseMessage)
				res.MethodName = messageFrame.MasterRequestMessage.Request.MethodName
				res.ReturnCode = int64(foundation.RequestMalformat)
				res.Message = err.Error()
				res.Timestamp = time.Now().Unix()

				me.sendResponse(res, messageFrame.Conn)
				continue
			}

			err = VerifySignature(me.Config, messageFrame.MasterRequestMessage, req)
			if err != nil {
				me.sendErrorResponse(req, messageFrame.Conn, int64(foundation.InvalidSignature), err.Error())
				continue
			}

			if req.MethodName == password_api.PasswordRequestMethodName {
				res := processPasswordRequest(me.Config, req)
				me.sendResponse(res, messageFrame.Conn)
			} else {
				me.sendErrorResponse(req, messageFrame.Conn, int64(foundation.MethodNotFound), foundation.MethodNotFound.String())
			}
		}
	}
}

func (me *PasswordServer) sendErrorResponse(req *ws_message.RequestMessage, conn net.Conn, returnCode int64, message string) {

	res := createErrorResponse(req, returnCode, message)
	me.sendResponse(res, conn)
}

func (me *PasswordServer) sendResponse(res *ws_message.ResponseMessage, conn net.Conn) {

	data, err := json.Marshal(res.ResponseBase)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("ResponseBase marshal error ", err)
		return
	}

	sign, err := crypto.GenerateRSASignFromBase64(passwordServerPrivateKey, string(data))
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("GenerateRSASignFromBase64 error ", err)
		return
	}
	res.Sign = sign
	err = me.udsServer.SendJsonMessage(conn, res)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Send response error: ", err)
	}
}
