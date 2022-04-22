package test

// func TestSetMainnetWalletAddress(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	// configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	txHash, err := conf.SetMainnetWalletAddressFromSC(configServer, "0xf204a4Ef082f5c04bB89F7D5E6568B796096735a")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	t.Logf("tx hash: %s", txHash)
// }

// func TestGetMainnetWalletAddress(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	address, err := conf.GetMainnetWalletAddressFromSC(configServer)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	t.Logf("mainnet address: %s", address)
// }

// func TestSetInnetWalletAddress(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	txHash, err := conf.SetInnetWalletAddressFromSC(configServer, "0x0")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	t.Logf("tx hash: %s", txHash)
// }

// func TestGetInnetWalletAddress(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	address, err := conf.GetInnetWalletAddressFromSC(configServer)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	fmt.Printf("innet address: %s\r\n", address)
// }

// func TestAddCurrencyInfo(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	txHash, err := conf.AddCurrencyInfoFromSC(configServer, "0xf204a4Ef082f5c04bB89F7D5E6568B796096735a", "testing")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	t.Logf("tx hash: %s", txHash)
// }

// func TestDelCurrencyInfo(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	txHash, err := conf.DelCurrencyInfoFromSC(configServer, "testing")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	t.Logf("tx hash: %s", txHash)
// }
