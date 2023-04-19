package main

import (
	"context"
	"encoding/json"
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
		"dsm", false, false, false, false, nil)
	utils.FailOnError(err, "failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//msg := "Alice fell into rabbit hole"
	req := utils.Request{1, 100, "pay", "cash", "现金交付"}
	msg, _ := json.Marshal(&req)
	err = ch.PublishWithContext(ctx, "", q.Name, false, false,
		amqp.Publishing{ContentType: "*/*", Body: []byte(msg)})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)
}
