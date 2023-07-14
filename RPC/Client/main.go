package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbit"
	"log"
	"os"
	"os/signal"
	"time"
)

func Send(ch *amqp.Channel, corID string, req interface{}) <-chan amqp.Delivery {
	//declare a queue AS ReplyTo queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "declare a queue failed")
	//receive response from ReplyTO queue
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")
	// create goroutine to handle request from producers
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//Marshal request into json form ,can use jsoniter to replace json
	data, _ := json.Marshal(req)
	//send request
	serr := ch.PublishWithContext(ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", CorrelationId: corID, Body: data})
	utils.FailOnError(serr, "failed to send request")
	return msgs
}
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "connect to amqp server failed")
	ch, err := conn.Channel()
	utils.FailOnError(err, "make channel failed")
	var req interface{}
	corid := "42"
	msgs := Send(ch, corid, req)
	//DONE: SMOOTHLY QUIT THE SERVER

	//Receive results
	for msg := range msgs {
		if msg.CorrelationId == corid {
			log.Println("the result received from server is ", msg)
			break
		}
	}
	//Smoothly quit
	quit := make(chan os.Signal)
	time.Sleep(5 * time.Second)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Quit rpc server successfully")
}
