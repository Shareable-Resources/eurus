package test

import (
	"log"
	"os"
	"testing"

	"github.com/streadway/amqp"
)

// Publisher and Exchange
// [Publisher] ---- [Exchange]
// go test eurus-backend/user_service/kyc_service/kyc/test -v -run Test_3_2_Emit_Logs
// Define a exchange name=logs
// Define a publisher which send message to exchanges logs
func Test_3_2_Emit_Logs(t *testing.T) {

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

	body := lastArg(os.Args)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key, because fanout doesn' t require routing keyt
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
