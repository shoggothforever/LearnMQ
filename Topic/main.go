package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	mes "go-rabbit"
	"go-rabbit/common"
	"log"
	"os"
)

func dead() {
	rbc := common.Rbc
	rq, err := common.NewRabbit(rbc.DeadName, rbc.DeadExchangeName, common.Topic, rbc.DeadRouting, rbc.Mqurl)
	if err != nil {
		mes.FailOnError(err, "创建死信连接失败")
	}
	defer rq.ReleaseMQ()
	//common.NewRabbit(rbc.DeadName, rbc.DeadExchangeName, rbc.DeadExchangeType, "", rbc.Mqurl)
	//common.NewRabbit(rbc.DeadName, rbc.DeadExchangeName, rbc.DeadExchangeType, "", rbc.Mqurl)
	//common.NewRabbit(rbc.DeadName, rbc.DeadExchangeName, rbc.DeadExchangeType, "", rbc.Mqurl)
	//声明死信交换机
	err = rq.Channel.ExchangeDeclare(rbc.DeadExchangeName, common.Topic, true, true, false, false, nil)
	mes.FailOnError(err, "声明死信交换机失败")
	//声明死信队列
	deadq, err := rq.CreateDeadQueue(rbc.DeadQueueName, rbc.DeadExchangeName, rbc.DeadRouting, true)
	mes.FailOnError(err, "创建死信队列失败")
	//fmt.Println(deadq)
	//消费死信队列中的信息
	msgs, err := rq.Channel.Consume(
		deadq.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	mes.FailOnError(err, "Failed to register a consumer")
	var forever chan struct{}
	go func() {
		//业务逻辑写在这里，可以将任务丢给线程池处理
		for d := range msgs {
			//time.Sleep(5 * time.Second)
			log.Printf("有消息过期了，消息来源:%s,消息内容: %s,", d.Exchange, d.Body)
		}
	}()
	<-forever
}

// 创建消息接收者
func getBindKey(s []string) []string {
	if len(s) <= 1 {
		return []string{"anonymous.infos"}
	} else {
		return s[1:]
	}
}

// go run main.go dream.*
func main() {
	go dead()
	rbc := common.Rbc
	rq, err := common.NewRabbit(rbc.BaseName, rbc.BaseExchangeName, common.Topic, "", rbc.Mqurl)
	mes.FailOnError(err, "获取mq失败")
	defer rq.ReleaseMQ()
	//声明原始队列
	q, err := rq.Channel.QueueDeclare("", true, true, true, false, amqp.Table{
		common.Publishdeadexchange: rbc.DeadExchangeName, // 死信交换机
		common.Publishdeadrouter:   rbc.DeadRouting,      // 路由键
		"x-message-ttl":            common.Xmessagettl,
		"x-max-length":             common.Xmaxlength,
	})
	mes.FailOnError(err, "创建队列失败")
	//将原始队列与死信交换机绑定
	//err = rq.Channel.QueueBind(q.Name, "", rbc.DeadExchangeName, false, nil)
	//mes.FailOnError(err, "绑定死信队列失败")
	bindkeys := getBindKey(os.Args)
	for _, s := range bindkeys {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, rbc.BaseExchangeName, s)
		err = rq.Channel.QueueBind(q.Name, s, rq.ExchangeName, false, nil)
		mes.FailOnError(err, "绑定键路由失败")
	}
	msgs, err := rq.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	mes.FailOnError(err, "Failed to register a consumer")
	var forever chan struct{}
	go func() {
		//业务逻辑写在这里，可以将任务丢给线程池处理
		for d := range msgs {
			//d.Nack(false, false)
			log.Printf("正常接受到了消息: %s", d.Body)
		}
	}()
	log.Printf("Waiting for logs. To exit press CTRL+C")
	<-forever
}
