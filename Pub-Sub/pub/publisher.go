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
	//首先要启动本地的rabbitMQ服务器，然后使用amqp与服务器创建connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "failed to connect RabbitMQ")
	defer conn.Close()
	//由connection产生channel，负责之后的通信交互
	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to open a channel")
	defer ch.Close()
	//由channel负责声明exchange，这是RabbitMQ中负责消息交换的重要组件
	err = ch.ExchangeDeclare(
		"up1",    // name
		"fanout", // type，订阅发布者模式的交换器种类为扇出型
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	utils.FailOnError(err, "failed to declare an exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := utils.BodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx, "up1", "", false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain",
			Body: []byte(body)})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
