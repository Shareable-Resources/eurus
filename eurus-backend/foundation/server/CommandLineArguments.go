package server

import "flag"

type CommandLineArguments struct {
	ConfigFilePath string
	IsNoUDSControl bool
	UDSPath        string
}

func (me *CommandLineArguments) ParseCommandLineArgument() {
	flag.StringVar(&me.ConfigFilePath, "config", "", "Config file path")
	flag.BoolVar(&me.IsNoUDSControl, "noUds", false, "Disable UDS control server")
	flag.StringVar(&me.UDSPath, "udsPath", "", "Specified UDS control server path")
	flag.Parse()
}
