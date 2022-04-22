package test

import (
	"log"
	"testing"

	"github.com/streadway/amqp"
)

//To test Publish and Subscribe
// 1. go test eurus-backend/user_service/kyc_service/kyc/test -v -run Test_3_1_Receive_Logs
// 2. go test eurus-backend/user_service/kyc_service/kyc/test -v -run Test_3_1_Receive_Logs > logs_from_rabbit.log
// 3. go test eurus-backend/user_service/kyc_service/kyc/test -v -run Test_3_2_Emit_Logs
// 4. rabbitmqctl list_bindings to see bingings
// Publisher and Exchange
// [Publisher] ---- [Exchange]
// go test eurus-backend/user_service/kyc_service/kyc/test -v -run Test_3_1_Receive_Logs
// Define a exchange name=logs
// Define a publisher which send message to exchanges logs
func Test_3_1_Receive_Logs(t *testing.T) {
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
		"",    // name, "" means auto generated name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// bind exchange [logs] with queue
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
