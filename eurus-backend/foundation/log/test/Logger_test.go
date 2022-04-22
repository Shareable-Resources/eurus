package test

import (
	log "eurus-backend/foundation/log"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	logger, _ := log.NewLogger("TestLogger", "/Users/duncanto/Documents/test_logger.log", logrus.DebugLevel)
	logger.Debug("Hello ", "World")

	log.GetDefaultLogger().Debug("Hi")
	fmt.Print("abc")
}
