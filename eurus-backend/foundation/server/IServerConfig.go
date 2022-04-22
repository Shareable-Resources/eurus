package server

import (
	"github.com/sirupsen/logrus"
)

type IServerConfig interface {
	ValidateField()
	SetHttpErrorLogger(logger *logrus.Logger)
	GetServerConfigBase() *ServerConfigBase
	GetParent() interface{}
	GetConfigFileOnlyFieldList() []string
}
