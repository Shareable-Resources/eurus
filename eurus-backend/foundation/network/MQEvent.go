package network

type IMQEvent interface {
	GetEventType() string
	SetData(data interface{})
	GetData() interface{}
}

type MQEvent struct { //implements IMQEvent
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

func (me *MQEvent) GetEventType() string {
	return me.EventType
}

func (me *MQEvent) SetData(data interface{}) {
	me.Data = data
}

func (me *MQEvent) GetData() interface{} {
	return me.Data
}
