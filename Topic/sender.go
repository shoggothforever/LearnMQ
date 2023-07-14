package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	mes "go-rabbit"
	"go-rabbit/common"
	"os"
	"time"
)

func main() {
	rbc := common.Rbc
	//发送消息，获取路由键
	key := mes.RouterKey(os.Args)
	rq, err := common.NewRabbit(rbc.BaseName, rbc.BaseExchangeName, common.Topic, key, rbc.Mqurl)
	mes.FailOnError(err, "获取mq失败")
	defer rq.ReleaseMQ()
	body := os.Args[1]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 1; i <= 10; i++ {
		rq.Channel.PublishWithContext(ctx, rbc.BaseExchangeName, rq.RoutingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Expiration:  "1000",
			Body:        []byte(body),
		})
	}
	fmt.Println(body, "消息发送成功")
}
