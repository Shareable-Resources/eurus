package main

import (
	"bufio"
	"encoding/json"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/uds"
	"eurus-backend/foundation/ws/ws_message"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type ServerControlClient struct {
	messageReceivedChan chan *uds.MessageFrame
	client              *uds.UDSMessageClient
	FinishChan          chan bool
	config              *ServerControlClientConfig
}

func NewServerControlClient() *ServerControlClient {
	client := new(ServerControlClient)
	client.messageReceivedChan = make(chan *uds.MessageFrame)
	client.FinishChan = make(chan bool)
	client.config = new(ServerControlClientConfig)
	client.config.IsEncrypted = true
	return client
}

func (me *ServerControlClient) LoadConfig(configPath string) error {

	pathByte := []byte(configPath)
	if configPath == "" || (len(pathByte) == 2 && pathByte[0] == 34 && pathByte[1] == 34) {
		currExecPath, err := os.Getwd()
		if err != nil {
			logger := log.GetLogger(log.Name.Root)
			logger.Error(err.Error())
			return err
		}
		configPath = path.Join(currExecPath, "config", "ServerControlClientConfig.json")
	}

	fmt.Println("Loading config file at path: ", configPath)
	configByte, loadErr := ioutil.ReadFile(configPath)
	if loadErr != nil {
		logger := log.GetLogger(log.Name.Root)
		logger.Error(loadErr.Error())
		return loadErr
	}

	err := json.Unmarshal(configByte, &me.config)
	if err != nil {
		logger := log.GetLogger(log.Name.Root)
		logger.Error("Log config error: ", err.Error())
		return err
	}

	return err
}

func (me *ServerControlClient) Connect(socketPath string) error {
	me.client = uds.NewUDSMessageClient(log.GetLogger(log.Name.Root), me.messageReceivedChan)

	err := me.client.Connect(socketPath)
	if err != nil {
		return err
	}

	go me.onMessageReceived()

	return nil
}

func (me *ServerControlClient) InitTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		commands := strings.Split(text, " ")
		var args []string = make([]string, 0)
		if len(commands) > 1 {
			args = commands[1:]
		}
		if commands[0] == "" {
			continue
		}

		commands[0] = strings.ToLower(commands[0])
		masterReq, err := ws_message.CreateMasterRequestMessage(args, commands[0])
		if err != nil {
			fmt.Println("Cannot create request: ", err)
			continue
		}
		fmt.Println("Sending request: ", masterReq.Message)
		err = me.client.SendMessage(masterReq)
		if err != nil {
			fmt.Println("Cannot send command to server: ", err)
		}
	}
}

func (me *ServerControlClient) onMessageReceived() {
	for {
		messageReceived := <-me.messageReceivedChan
		if messageReceived.Error != nil {
			fmt.Println("Error received: ", messageReceived.Error)
			continue
		} else if messageReceived.IsConnectionEstablishedEvent {
			continue
		} else if messageReceived.ResponseMessage != nil {
			rawJson, ok := messageReceived.ResponseMessage.Data.(*json.RawMessage)
			if ok {
				val := string(*rawJson)
				val = strings.ReplaceAll(val, `\"`, "\"")
				val = strings.ReplaceAll(val, "\\n", "\n")
				val = strings.ReplaceAll(val, "\\t", "\t")
				val = strings.ReplaceAll(val, "\\r", "\r")
				val = strings.ReplaceAll(val, "\\\\", "\\")
				switch messageReceived.ResponseMessage.MethodName {
				case "welcomeMessage":
					fmt.Println(val + " connected")
				default:
					fmt.Println(val)
				}
			} else {
				fmt.Println("Messge received: ", messageReceived.ResponseMessage.Message)
			}
		}
	}
}
