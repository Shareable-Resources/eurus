package kyc_const

import "eurus-backend/foundation/network"

var TaskQueueMetaData = network.MQTaskQueueMetaData{
	QueueName:    "kyc_task_queue",
	IsExclusive:  false,
	IsAutoAck:    false,
	IsAutoDelete: true,
}

var DBQueryConfig = struct {
	RowLimit int
}{
	RowLimit: 50,
}
