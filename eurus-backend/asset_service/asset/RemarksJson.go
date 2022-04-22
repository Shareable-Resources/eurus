package asset

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type RemarksJson map[string]interface{}

func NewRemarksJsonFromError(err error) RemarksJson {
	var causeError error = errors.Cause(err)
	var remarksJson RemarksJson = make(RemarksJson)

	var message string
	if causeError != nil {
		message = causeError.Error()
	} else {
		message = err.Error()
	}
	if index := strings.Index(causeError.Error(), "Inner error JSON:"); index >= 0 {
		message = causeError.Error()[index+len("Inner error JSON:"):]
		index = strings.Index(message, "}:")
		if index >= 0 {
			message = message[:index+1]
			var messageMap RemarksJson
			err := json.Unmarshal([]byte(message), &messageMap)
			if err != nil {
				remarksJson["remarksType"] = "errorMessage"
				remarksJson["message"] = message
			} else {
				remarksJson["remarksType"] = "errorJson"
				remarksJson["message"] = messageMap
			}
		}
	} else {
		remarksJson["remarksType"] = "errorMessage"
		remarksJson["message"] = message
		remarksJson["inner"] = causeError
	}
	return remarksJson
}

func NewRemarksJsonFromReceipt(receipt []byte) RemarksJson {
	var remarksJson RemarksJson = make(RemarksJson)
	remarksJson["remarksType"] = "receipt"

	var receiptMap RemarksJson
	err := json.Unmarshal(receipt, &receiptMap)
	if err == nil {
		remarksJson["receipt"] = receiptMap
	}
	return remarksJson
}

func NewRemarksJsonFromString(message string) RemarksJson {
	var remarksJson map[string]interface{} = make(map[string]interface{})
	remarksJson["remarksType"] = "errorMessage"
	remarksJson["message"] = message
	return remarksJson
}

func (me *RemarksJson) String() string {
	data, _ := json.Marshal(me)
	return string(data)
}
