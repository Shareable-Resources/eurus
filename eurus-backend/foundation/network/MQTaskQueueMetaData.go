package network

type MQTaskQueueMetaData struct {
	QueueName    string
	IsExclusive  bool
	IsAutoDelete bool
	IsAutoAck    bool
}
