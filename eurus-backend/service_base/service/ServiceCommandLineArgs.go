package service

import (
	"eurus-backend/foundation/server"
	"flag"
)

type ServiceCommandLineArgs struct {
	server.CommandLineArguments
	UDSPath            string
	PasswordServerPath string
}

func (me *ServiceCommandLineArgs) ParseCommandLineArgument() {

	flag.StringVar(&me.UDSPath, "uds", "", "Open Unix Domain Socket path to input config decryption password")
	flag.StringVar(&me.PasswordServerPath, "pwServer", "", "Password server path. Get password from password server.")
	me.CommandLineArguments.ParseCommandLineArgument()

}
