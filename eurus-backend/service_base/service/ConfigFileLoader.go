package service

import (
	"encoding/base64"
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/uds"
	"eurus-backend/foundation/ws/ws_message"
	"eurus-backend/password_service/password_api"
	"eurus-backend/secret"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func LoadConfigFile(configFilePath string, isEncrypted bool, isDeleteAfterUsed bool, passordServerPath string, openUdsPath string) ([]byte, error) {
	configByte, loadErr := ioutil.ReadFile(configFilePath)
	if loadErr != nil {
		logger := log.GetLogger(log.Name.Root)
		logger.Error(loadErr.Error())
		fmt.Println("Read config file error: ", loadErr)
		return nil, errors.Wrap(loadErr, "Read config file error")
	}
	var password string
	if isEncrypted {
		if passordServerPath == "" && openUdsPath == "" {
			var pw []byte
			var err error
			for {
				fmt.Print("Input password to decrypt config file (screen is hidden): ")
				pw, err = terminal.ReadPassword(0)
				if err != nil {
					panic(err)
				}
				if len(pw) > 0 {
					password = string(pw)
					break
				}
			}
		} else if passordServerPath != "" {
			fmt.Println("Query password from password server")
			var err error
			password, err = queryPasswordFromPasswordServer(passordServerPath)
			if err != nil {
				fmt.Println("queryPasswordFromPasswordServer error: ", err)
				return nil, errors.Wrap(err, "queryPasswordFromPasswordServer error")
			}
			fmt.Println("Password received via Password server")
		} else {
			fmt.Println("Waiting for password from UDS port")
			var err error
			password, err = getPasswordFromUDS(openUdsPath)
			if err != nil {
				return nil, errors.Wrap(err, "getPasswordFromUDS error")
			}
			fmt.Println("Password received via UDS port")
		}

		fmt.Println("\nDecrypting config file")
		decrypted, err := crypto.DecryptByPassword(password, string(configByte))
		if err != nil {
			fmt.Println("DecryptByPassword failed: ", err)
			return nil, errors.Wrap(err, "Unable to decrypt config file")
		}
		configByte, err = base64.StdEncoding.DecodeString(decrypted)
		if err != nil {
			fmt.Println("Base64 decode decrypted config file failed: ", err)
			return nil, errors.Wrap(err, "Base64 decode decrypted config file failed")
		}
	}

	if isDeleteAfterUsed {
		os.Remove(configFilePath)
	}
	return configByte, nil
}

func queryPasswordFromPasswordServer(serverPath string) (string, error) {

	for {
		password, err := queryPasswordFromPasswordServerImpl(serverPath)
		if err == nil {
			return password, nil
		}
		log.GetLogger(log.Name.Root).Infoln("Wait for 1 second to retry")
		time.Sleep(time.Second)
		continue
	}
}

func queryPasswordFromPasswordServerImpl(serverPath string) (string, error) {

	var messageEventHandler chan (*uds.MessageFrame) = make(chan *uds.MessageFrame)
	client := uds.NewUDSMessageClient(log.GetLogger(log.Name.Root), messageEventHandler)

	err := client.Connect(serverPath)
	if err != nil {
		fmt.Println("Connect password server failed: ", err)
		return "", err
	}

	defer func() {
		client.Close()
	looping:
		for {
			select {
			case <-messageEventHandler:
				continue
			case <-client.CloseChan:
				break looping
			}
		}
		close(messageEventHandler)
	}()

	passReq := new(password_api.PasswordRequest)
	masterReq, err := ws_message.CreateMasterRequestMessage(passReq, passReq.MethodName())
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - CreateMasterRequestMessage failed: ", err)
		return "", err
	}

	masterReq.Sign, err = secret.GeneratePasswordClientSignature(masterReq.Message)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Sign request failed: ", err)
		return "", err
	}

	err = client.SendMessage(masterReq)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - SendMessage failed: ", err)
		return "", err
	}
	log.GetLogger(log.Name.Root).Infoln("Send request to password server")

	var messageFrame *uds.MessageFrame

	select {
	case messageFrame = <-messageEventHandler:
		break
	case <-time.After(time.Second * 15):
		break
	}

	if messageFrame == nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Receive message timeout:  wait for 1 second")
		return "", errors.New("Receive timeout")
	}
	if messageFrame.Error != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Password server response error: ", messageFrame.Error)
		return "", messageFrame.Error
	}

	if messageFrame.ResponseMessage == nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Unexpected message receive from password server")
		return "", errors.New("Unexpected message receive from password server")
	}

	passRes, ok := messageFrame.ResponseMessage.Data.(*password_api.PasswordResponse)
	if !ok {
		fmt.Println("queryPasswordFromPasswordServerImpl - Password server does not response expected data")
		return "", errors.New("Password server does not response expected data")
	}

	if messageFrame.ResponseMessage.Sign == "" {
		fmt.Println("queryPasswordFromPasswordServerImpl - Response does not contains signature")
		return "", errors.New("Response does not contains signature")
	}

	data, err := json.Marshal(messageFrame.ResponseMessage.ResponseBase)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Marshal ResponseBase failed: ", err)
		return "", err
	}

	isValid, err := secret.VerifyPasswordServerSignature(string(data), messageFrame.ResponseMessage.Sign)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Verify sign failed: ", err)
		return "", err
	}
	if !isValid {
		fmt.Println("queryPasswordFromPasswordServerImpl - VerifyPasswordServerSignature Invalid signature")
		return "", errors.New("Invalid signature")
	}

	password, err := secret.DecryptPasswordServerData(passRes.EncryptedPassword)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Unable to decrypted response password. Error: ", err.Error())
		return "", err
	}
	val, _ := base64.StdEncoding.DecodeString(password)
	password = string(val)
	return password, nil

}

func getPasswordFromUDS(socketPath string) (string, error) {
	var messageEventHandler chan (*uds.MessageFrame) = make(chan *uds.MessageFrame)
	var passwordChan chan (string) = make(chan string)
	var serverClosedChan chan (error) = make(chan error)
	server, err := uds.NewUDSMessageServer(socketPath, log.GetLogger(log.Name.Root), messageEventHandler)
	if err != nil {
		fmt.Println("queryPasswordFromPasswordServerImpl - Unable to init UDS server: ", err)
		return "", err
	}

	go onMessageReceived(server, messageEventHandler, passwordChan)

	go func() {
		err = server.Listen()
		if err != nil && err != io.EOF {
			fmt.Println("Unable to bind UDS server: ", err)
		}
		serverClosedChan <- err
	}()

	select {
	case pw := <-passwordChan:
		return pw, nil
	case sererErr := <-serverClosedChan:
		return "", sererErr
	}
}

func onMessageReceived(server *uds.UDSMessageServer, messageEventHandler chan (*uds.MessageFrame), passwordChan chan (string)) {

	for {
		messageEvent := <-messageEventHandler
		if messageEvent.Error != nil {
			return
		}

		if messageEvent.MasterRequestMessage != nil {
			req := new(ws_message.RequestMessage)
			err := json.Unmarshal([]byte(messageEvent.MasterRequestMessage.Message), &req)
			if err != nil {
				fmt.Println("Receive invalid message: ", messageEvent.MasterRequestMessage.Message)
				continue
			}

			if passwordReq, ok := req.Data.(*PasswordRequest); ok {
				var password string
				res := new(ws_message.ResponseMessage)

				if passwordReq.IsEncrypted {
					val, err := secret.DecryptConfigValue(passwordReq.Password)
					if err != nil {
						res.ReturnCode = int64(foundation.InvalidEncryption)
						res.Message = foundation.InvalidEncryption.String()
					} else {
						data, _ := base64.StdEncoding.DecodeString(val)
						password = string(data)
						res.ReturnCode = int64(foundation.Success)
						res.Message = foundation.Success.String()
						passwordChan <- password
					}

				} else {
					password = passwordReq.Password
					res.ReturnCode = int64(foundation.Success)
					res.Message = foundation.Success.String()
					passwordChan <- password
				}

				_ = server.SendJsonMessage(messageEvent.Conn, res)
				messageEvent.Conn.Close()
				server.Close()
			}
		} else if messageEvent.Error != nil {
			log.GetLogger(log.Name.Root).Error("Error receive: ", messageEvent.Error)
		}
	}

}
