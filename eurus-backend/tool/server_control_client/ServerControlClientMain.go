package main

import (
	"eurus-backend/foundation/log"
	"flag"
)

func main() {

	controlClient := NewServerControlClient()
	var socketPath string
	flag.StringVar(&socketPath, "uds", "", "Unix domain socket path")
	flag.Parse()

	if socketPath == "" {
		panic("uds is mandatory argument")
	}

	var err error
	err = controlClient.Connect(socketPath)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to server via ", socketPath, " error: ", err)
		return
	}

	controlClient.InitTerminal()

}
