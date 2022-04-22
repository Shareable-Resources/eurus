package test

import (
	"bytes"
	"eurus-backend/foundation/database"
	eurus_log "eurus-backend/foundation/log"
	"eurus-backend/foundation/server"
	"eurus-backend/secret"
	"eurus-backend/user_service/kyc_service/kyc"
	kyc_model "eurus-backend/user_service/kyc_service/kyc/model"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func lastArg(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = args[len(args)-1]
	}
	return s
}

func TestPublisherSendToExchange(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

// Consumer and Exchange
// [Exchange] ---- [Queue] ---- [Consumer]
// go test eurus-backend/user_service/kyc_service/kyc/test -v -run TestConsumerReceiveFromExchangeQueue
// Define a exchange name=logs
// Define a queue which bind to the exchange
// Define a consumer which consume the queue messages
func TestConsumerReceiveFromExchangeQueue(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// Publisher
// go test eurus-backend/user_service/kyc_service/kyc/test -v -run TestPublisherSendToQueue
// [Publisher] ---- [Queue]
// Define a routing key, which is the queue
// Define a publisher which send message to queue
func TestPublisherSendToQueue(t *testing.T) {
	fmt.Println("TestPublisherSendToQueue")
	body := bodyFrom(os.Args)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //This will save msg to disk
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

}

// Define Queue and consumer
// Same worker can be run in different console
// go test eurus-backend/user_service/kyc_service/kyc/test -v -run TestConsumerReceiveFromQueue
func TestConsumerReceiveFromQueue(t *testing.T) {
	fmt.Println("TestConsumerReceiveFromQueue")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Deliever message to less heavry worker
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Queue 1: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			// Acknowlegement will send to rabbitMQ to notify them that a queue is processed
			// or setting auto-ack to true whe init from ch.Consume
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

// go test eurus-backend/user_service/kyc_service/kyc/test -v -run TestConsumerReceiveFromQueueDurable
func TestConsumerReceiveFromQueueDurable(t *testing.T) {
	fmt.Println("TestConsumerReceiveFromQueueDurable")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q2, err := ch.QueueDeclare(
		"hello2", // name
		true,     // durable - Queues and exchanges needs to be configured as durable in order to survive a broker restart. A durable queue only means that the queue definition will survive a server restart, not the messages in it.
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q2.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Queue 2: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			//d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func TestInsertKYCImage(t *testing.T) {
	config := server.ServerConfigBase{
		DBAESKey: "rLoBACaeYMQea0w/8bouX/2KyjPXUhrLG87lU7pJIegOSjIKC7OTJrZYDowMOmiB1XfWMMFGGVZEPFLTPm59atRoGa5oJGH72ZaSkTn04anhdGgY+VFUslAZ8RG/IRh2WuS8tMd8GPxX7ObV59i9DcRXJL1nR155QDEoR1+Cw3hiciAwnBDc+ZkCQuA6X9NcOS/2+ko+GZL8TKQUU77PLq5QdoZOfQCW89qMaJcW6UJOF91NqPjj7i8odAXC3GPIRU+ctA==",
	}
	secret.DecryptSensitiveConfig(&config)

	db := database.Database{
		IP:         "18.167.146.127",
		Port:       5432,
		UserName:   "admin",
		Password:   "Xy1JtrlrWLlMHNHlRYpINC/f75Tk5zJRwc0i2huveik=",
		DBName:     "postgres",
		SchemaName: "public",
		Logger:     eurus_log.GetDefaultLogger(),
	}

	db.SetAESKey(config.DBAESKey)

	_, err := kyc.DbInsertKYCImage(&db, 1, kyc_model.KYCImgePassport)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPointer(t *testing.T) {
	nPtr := new(int)
	fmt.Println(*nPtr)
}
