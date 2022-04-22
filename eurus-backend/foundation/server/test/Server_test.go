package test

import (
	"eurus-backend/authen/auth"
	"eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"net/http"
	"time"

	"testing"
)

type TestServer struct {
	server.ServerBase
}

func TestNewServer(t *testing.T) {
	testServer := new(TestServer)
	testServer.InitLog("TestServer.log")

	logger := log.GetLogger(log.Name.Root)
	logger.Debug("Hi")
	logger.Error("This is an error log")

	serverConfig := new(server.ServerConfigBase)
	err := testServer.LoadConfig("TestServerConfig.json", serverConfig)
	if err != nil {
		t.Error(err.Error())
	}
	testServer.ServerConfig = serverConfig

	_, err = testServer.InitDBFromConfig(testServer.ServerConfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = testServer.InitEthereumClientFromConfig(testServer.ServerConfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.InitHttpServer(testServer.ServerConfig)
	if err != nil {
		t.Errorf(err.Error())
	}
	testServer.ServerBase.HttpServer.SetIndexResponse(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			processTime := time.Duration(0) * time.Second
			select {
			case <-ctx.Done():
				// this is the case that exceed the HTTP receive timeout time
				logger.Error("exceed the HTTP receive timeout time")
				return
			case <-time.After(processTime):
				w.Write([]byte("welcome"))
			}
		})

	authClient := auth.NewAuthClient()
	testServer.InitAuth(authClient, testServer.ServerConfig)

	// fmt.Println("sleep for 8 seconds, you may use teminal to see if the server port is opened by command line 'netstat -an'")
	// time.Sleep(8 * time.Second)
	// fmt.Println("end sleep")
	err = testServer.ServerBase.HttpServer.Listen()
	if err != nil {
		t.Errorf(err.Error())
	}

}
