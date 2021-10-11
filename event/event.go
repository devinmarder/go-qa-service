package event

import (
	"log"

	"github.com/streadway/amqp"
)

//ExchangeName is the name used for the exchange that the event is published to.
const ExchangeName = "qa.events"

//RunEventProducer listens to mesages sent to the channel and publishes them to the Exchange.
func RunEventProducer(msgs chan string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed declare exchange")

	//Listen for events on channel.
	for msg := range msgs {
		err := ch.Publish(
			ExchangeName, // exchange
			"",           // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "json",
				Body:        []byte(msg),
			},
		)
		failOnError(err, "Failed to send message")
	}
}

//failOnError logs the error and Terminates the process.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
