package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbit"
	"log"
	"os"
	"reflect"
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
	msgs, err := ch.Consume(q.Name, "dsm", true, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var rec utils.Request
			json.Unmarshal(d.Body, &rec)
			value := reflect.ValueOf(&rec).Elem()
			log.Printf("Received a message:%v  %v ,%v", rec, value, value.Type())

			os.WriteFile("consumer"+q.Name+".txt", d.Body, 0666)
			value.Field(2).Set(reflect.ValueOf("refund"))
			log.Printf("Received a message:%v  %v ,%v", rec, value, value.Type())
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
