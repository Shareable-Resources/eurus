package userObserver

import (
	"eurus-backend/foundation/server"
	"time"
)

type UserObserverConfig struct { //implements IRetrySetting
	server.ServerConfigBase
	SideChainGasLimit uint64 `json:"sideChainGasLimit"`
}

func NewUserObserverConfig() *UserObserverConfig {
	config := new(UserObserverConfig)
	config.ActualConfig = config
	return config
}

func (me *UserObserverConfig) ValidateEthClientField() error {
	var err error
	err = me.ServerConfigBase.ValidateEthClientField()
	if err != nil {
		return err
	}

	return err
}

func (me *UserObserverConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *UserObserverConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *UserObserverConfig) GetRetryCount() int {
	return me.RetryCount
}

func (me *UserObserverConfig) GetRetryInterval() time.Duration {
	return time.Duration(me.RetryInterval)
}
