package main

import (
	"encoding/json"
	"eurus-backend/auth_service/auth"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/network"
	"eurus-backend/foundation/server"
	"fmt"
	"net/http"
	"strings"

	"time"
)

type TestServer struct {
	server.ServerBase
}

func (me *TestServer) InitTerminal(args []string) (bool, error) {
	command := args[0]
	command = strings.ToLower(command)

	if command == "requestlogintoken" {
		loginToken, err := me.AuthClient.GenerateLoginToken("Hi")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(loginToken.GetToken())
		}

		return true, nil
	}

	return false, nil
}

func main() {
	testServer := new(TestServer)
	testServer.InitLog("TestServer.log")

	logger := log.GetLogger(log.Name.Root)
	logger.Debug("Hi")
	logger.Error("This is an error log")

	serverConfig := new(server.ServerConfigBase)
	arguments := &server.CommandLineArguments{}
	arguments.ParseCommandLineArgument()
	err := testServer.LoadConfig(arguments.ConfigFilePath, serverConfig)
	if err != nil {
		logger.Error(err.Error())
	}

	testServer.ServerConfig = serverConfig

	err = testServer.InitDBFromConfig(testServer.ServerConfig)
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = testServer.InitEthereumClientFromConfig(testServer.ServerConfig)
	if err != nil {
		logger.Error(err.Error())
	}
	err = testServer.InitHttpServer(testServer.ServerConfig)
	if err != nil {
		logger.Error(err.Error())
	}
	count := 0
	testServer.ServerBase.HttpServer.SetIndexResponse(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			processTime := time.Duration(0) * time.Second
			select {
			case <-ctx.Done():
				// this is the case that exceed the HTTP receive timeout time
				return
			case <-time.After(processTime):
				testServer.HttpServer.Logger.Error("test logger")
				logger.Error("count: ", count)
				count++
				w.Write([]byte("welcome"))
			}
		})

	go testServer.ServerBase.HttpServer.Listen()
	var client *auth.AuthClient = auth.NewAuthClient()
	client.SetLoginHandler(func(authClient network.IAuth) {
		addInfo := authClient.GetAdditionalInfo()

		err := json.Unmarshal([]byte(addInfo), testServer.ServerConfig)
		if err == nil {
			fmt.Println("Additional info: ", addInfo)
		}
	})
	testServer.InitAuth(client, testServer.ServerConfig, nil)

	go testServer.ServerBase.InitTerminal(nil)

	for {
		if testServer.AuthClient.IsLoggedIn() {
			loginToken, err := testServer.AuthClient.GenerateLoginToken("Some caller used identifier for the token")
			if err == nil {
				fmt.Println(loginToken)
				isValid, verifyLoginToken, err := testServer.AuthClient.VerifyLoginToken(loginToken.GetToken(), network.VerifyModeService|network.VerifyModeUser)
				fmt.Println("Verify result: ", isValid)
				if err == nil {

					fmt.Println(verifyLoginToken)
				} else {
					fmt.Println(err.Error())
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
