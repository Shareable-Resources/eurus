package network

import (
	"encoding/json"
	"eurus-backend/foundation/log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type MQMode int

const (
	MQModeTopic MQMode = iota
	MQModeTaskQueue
)

type MQPublisher struct {
	MQUrl               string
	ExchangeName        string
	QueueName           string
	Logger              *logrus.Logger
	mq                  *amqp.Channel
	closeEvent          chan *amqp.Error
	isAutoDeleteMessage bool
}

/// Create a new publisher
/// mqUrl - MQ end point URL
/// exchangeName - exchange name, empty string to create publisher without topic
func NewMQPublisher(mqUrl string, publisherMode MQMode, exchangeOrQueueName string, logger *logrus.Logger) *MQPublisher {
	publisher := new(MQPublisher)
	publisher.MQUrl = mqUrl
	publisher.Logger = logger

	switch publisherMode {
	case MQModeTopic:
		publisher.ExchangeName = exchangeOrQueueName
	case MQModeTaskQueue:
		publisher.QueueName = exchangeOrQueueName
	default:
		return nil
	}

	return publisher
}

func (me *MQPublisher) InitPublisher(isAutoDeleteMessage bool) error {

	me.isAutoDeleteMessage = isAutoDeleteMessage

	conn, err := amqp.Dial(me.MQUrl)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to connect to MQ: ", err)
		return err
	}

	me.mq, err = conn.Channel()
	if err != nil {
		return err
	}

	if me.ExchangeName != "" {
		err = me.mq.ExchangeDeclare(
			me.ExchangeName,
			"topic",
			true,
			isAutoDeleteMessage,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	me.closeEvent = make(chan *amqp.Error)

	me.mq.NotifyClose(me.closeEvent)
	go me.monitorConnection()

	return nil
}

func (me *MQPublisher) PublishPlainText(topic string, content string, headers map[string]interface{}) error {
	return me.PublishWithContentType(topic, "text/plain", []byte(content), headers)
}

func (me *MQPublisher) PublishJson(topic string, object interface{}, headers map[string]interface{}) error {
	content, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return me.PublishWithContentType(topic, "application/json", []byte(content), headers)
}

func (me *MQPublisher) PublishWithContentType(topic string, contentType string, content []byte, headers map[string]interface{}) error {
	routingKey := topic
	if me.QueueName != "" {
		routingKey = me.QueueName
	}

	err := me.mq.Publish(me.ExchangeName, routingKey, false, false, amqp.Publishing{
		Headers:     headers,
		ContentType: contentType,
		Body:        content,
	})

	return err
}

func (me *MQPublisher) monitorConnection() {

	for {
		connEvent := <-me.closeEvent
		if me.closeEvent == nil {
			continue
		}
		if me.Logger != nil {
			me.Logger.Errorln("MQ connection error: ", me.QueueName, " Code: ", connEvent.Code, " Reason: ", connEvent.Reason, " Can recover: ", connEvent.Recover, " Is server: ", connEvent.Server)
		}

		for {
			err := me.InitPublisher(me.isAutoDeleteMessage)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			if me.Logger != nil {
				me.Logger.Infoln("MQ connection resumed: ", me.QueueName)
			}
			return
		}
	}

}
