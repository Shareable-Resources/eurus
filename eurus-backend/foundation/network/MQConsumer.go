package network

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type consumerMode int

const (
	consumerTopic consumerMode = iota
	consumerTaskQueue
)

type MQConsumer struct {
	Logger     *logrus.Logger
	mq         *amqp.Channel
	closeEvent chan *amqp.Error

	mode              consumerMode
	mqUrl             string
	subscribeTopic    string
	exchangeMetaData  *MQExchangeMetaData
	taskQueueMetaData *MQTaskQueueMetaData

	queueName   string
	isExclusive bool
	isAutoAck   bool

	eventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)
}

func (me *MQConsumer) SubscribeTopic(mqUrl string, subscribeTopic string, exchangeMetaData *MQExchangeMetaData,
	taskQueueMetaData *MQTaskQueueMetaData, eventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)) error {

	me.mode = consumerTopic
	me.mqUrl = mqUrl
	me.subscribeTopic = subscribeTopic
	me.exchangeMetaData = exchangeMetaData
	me.taskQueueMetaData = taskQueueMetaData
	me.eventHandler = eventHandler

	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		return err
	}

	me.mq, err = conn.Channel()
	if err != nil {
		return err
	}

	err = me.mq.ExchangeDeclare(
		exchangeMetaData.ExchangeName,
		"topic",
		true,
		exchangeMetaData.IsAutoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := me.mq.QueueDeclare(taskQueueMetaData.QueueName, true, taskQueueMetaData.IsAutoAck, taskQueueMetaData.IsExclusive, false, nil)
	if err != nil {
		return err
	}

	err = me.mq.QueueBind(
		queue.Name,
		subscribeTopic,
		exchangeMetaData.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	me.closeEvent = make(chan *amqp.Error)
	me.closeEvent = me.mq.NotifyClose(me.closeEvent)

	return me.consume(queue.Name, taskQueueMetaData.IsExclusive, taskQueueMetaData.IsAutoAck, eventHandler)
}

func (me *MQConsumer) SubscribeTaskQueue(mqUrl string, taskQueueMetaData *MQTaskQueueMetaData, eventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)) error {

	me.mode = consumerTaskQueue
	me.mqUrl = mqUrl
	me.eventHandler = eventHandler
	me.taskQueueMetaData = taskQueueMetaData

	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		return err
	}

	me.mq, err = conn.Channel()
	if err != nil {
		return err
	}

	queue, err := me.mq.QueueDeclare(taskQueueMetaData.QueueName, true, taskQueueMetaData.IsAutoDelete, taskQueueMetaData.IsExclusive, false, nil)
	if err != nil {
		return err
	}

	me.closeEvent = make(chan *amqp.Error)
	me.closeEvent = me.mq.NotifyClose(me.closeEvent)

	return me.consume(queue.Name, taskQueueMetaData.IsExclusive, taskQueueMetaData.IsAutoAck, eventHandler)
}

func (me *MQConsumer) consume(queueName string, isExclusive bool, isAutoAck bool, eventHandler func(message *amqp.Delivery, topic string, contentType string, content []byte)) error {
	me.queueName = queueName
	me.isExclusive = isExclusive
	me.isAutoAck = isAutoAck
	me.eventHandler = eventHandler

	msgs, err := me.mq.Consume(
		queueName,
		"",
		isAutoAck,
		isExclusive,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case message := <-msgs:
				eventHandler(&message, message.RoutingKey, message.ContentType, message.Body)
			case errEvent := <-me.closeEvent:
				if errEvent != nil {
					if me.Logger != nil {
						me.Logger.Errorln(queueName, " disconnected. Code: ", errEvent.Code, " Reason: ",
							errEvent.Reason, " IsServer: ", errEvent.Server, " Can recover: ", errEvent.Recover, ". Trying to reconnect")
					}
					go me.reconnect()
				}
				break

			}
		}
	}()
	return nil
}

func (me *MQConsumer) reconnect() {
	for {
		if me.mode == consumerTopic {
			err := me.SubscribeTopic(me.mqUrl, me.subscribeTopic, me.exchangeMetaData, me.taskQueueMetaData, me.eventHandler)
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
		} else {
			err := me.SubscribeTaskQueue(me.mqUrl, me.taskQueueMetaData, me.eventHandler)
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
		}
		break
	}
	if me.Logger != nil {
		me.Logger.Errorln("MQ reconnected: ", me.queueName)
	}
}
