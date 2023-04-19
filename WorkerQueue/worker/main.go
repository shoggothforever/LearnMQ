package main

import (
	"bytes"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbit"
	"log"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "failed to connect RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"dsm_durable", true, false, false, false, nil)
	utils.FailOnError(err, "failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(q.Name, "dsm", false, false, false, false, nil)

	utils.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			fmt.Println(t)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
