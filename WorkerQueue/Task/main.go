package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbit"
	"log"
	"os"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := utils.BodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx, "", q.Name, false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain",
			Body: []byte(body)})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
