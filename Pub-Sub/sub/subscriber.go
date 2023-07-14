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

	//不指定队列名字，由rabbitMQ服务器自动创建
	q, err := ch.QueueDeclare(
		"", false, false, true, false, nil)
	utils.FailOnError(err, "failed to declare a queue")

	err = ch.QueueBind(q.Name, "", "up1", false, nil)
	utils.FailOnError(err, "Failed to Bind Queue:"+q.Name)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

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
		}
	}()

	log.Printf(" [*] Waiting for radio  messages. To exit press CTRL+C")
	<-forever
}
