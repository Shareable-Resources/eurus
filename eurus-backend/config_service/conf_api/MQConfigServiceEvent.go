package conf_api

import "eurus-backend/foundation/network"

type MQConfigServiceEvent struct {
	network.MQEvent
	Action PublishAction `json:"action"`
}

type PublishAction int

const (
	PublishActionUnknown PublishAction = iota
	PublishActionInsert
	PublishActionUpdate
	PublishActionDelete
)

func NewMQConfigServiceEvent(action PublishAction, data interface{}) *MQConfigServiceEvent {
	event := new(MQConfigServiceEvent)
	event.EventType = "ConfigServiceEvent"
	event.Action = action
	event.Data = data
	return event
}

func (action PublishAction) String() string {
	switch action {
	case PublishActionUnknown:
		return "unknown"
	case PublishActionInsert:
		return "insert"
	case PublishActionUpdate:
		return "update"
	case PublishActionDelete:
		return "delete"
	}
	return ""
}

func StringToPublishAction(str string) PublishAction {
	switch str {
	case "unknown":
		return PublishActionUnknown
	case "insert":
		return PublishActionInsert
	case "update":
		return PublishActionUpdate
	case "delete":
		return PublishActionDelete
	}
	return PublishActionUnknown
}
