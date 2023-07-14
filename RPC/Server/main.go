package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbit"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "connect to amqp server failed")
	ch, err := conn.Channel()
	utils.FailOnError(err, "make channel failed")
	//declare a queue to receive request
	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "declare a queue failed")
	//Using Qos function to make every worker fairly handle requests
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	reqs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")
	// create goroutine to handle request from producers
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		for req := range reqs {
			//1:HANDEL MSG BODY
			data := req.Body
			log.Println(data)
			//produce result
			var result string

			//3:Send results to Producer through msg.Replyto queue
			serr := ch.PublishWithContext(ctx, req.ReplyTo, "", false, false, amqp.Publishing{ContentType: "text/plain", CorrelationId: req.CorrelationId, Body: []byte(result)})
			utils.FailOnError(serr, "failed to send response")
			//4:manually ack
			req.Ack(false)
		}
	}()
	//DONE: SMOOTHLY QUIT THE SERVER
	quit := make(chan os.Signal)
	time.Sleep(5 * time.Second)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Quit rpc server successfully")
}
