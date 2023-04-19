package Topic

import (
	amqp "github.com/rabbitmq/amqp091-go"
	mes "go-rabbit"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:15672/")
	defer conn.Close()
	mes.FailOnError(err, "1")
	ch, err := conn.Channel()
	defer ch.Close()
	mes.FailOnError(err, "2")
	exchanger1 := ch.ExchangeDeclare(
		"",
		"topic",
		true,
		false,
		false,
		false,
		nil)
	queue1 := ch.QueueDeclare(
		"")

}
