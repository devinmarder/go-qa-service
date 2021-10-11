package event

import (
	"log"

	"github.com/streadway/amqp"
)

//queueName is the queue that the consumer listens to.
const queueName = "test_queue"

//RunEventlistenter consumes and logs events on the queue.
func RunEventlistener() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "faild to declare queue")

	err = ch.QueueBind(
		queueName,    // queue name
		"",           // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "faild to bind queue")

	msgs, err := ch.Consume(
		queueName,      // queue
		"testConsumer", // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("event recieved: %s", d.Body)
	}
}
