package asset

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type AssetLog struct {
}

func ConstructRescanLoggerName(loggerName string, serviceId int64, eventId decimal.Decimal) string {
	return loggerName + "_" + strconv.FormatInt(serviceId, 10) + "_" + eventId.String()
}

func ConstructLogFilePath(loggerName string, configFilePath string) string {
	var fileNamePrefix string
	var logFilePath string

	index := strings.LastIndex(configFilePath, ".")
	if index >= 0 {
		fileNamePrefix = configFilePath[:index]
		logFilePath = fileNamePrefix + "_" + loggerName + configFilePath[index:]
	} else {
		fileNamePrefix = configFilePath
		logFilePath = fileNamePrefix + "_" + loggerName
	}
	return logFilePath
}
