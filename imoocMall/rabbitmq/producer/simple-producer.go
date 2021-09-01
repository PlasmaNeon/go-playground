package main

import (
	"example/rabbitmq"
	"fmt"
)

func main() {
	mq := rabbitmq.NewRabbitMQSimple("imoocSimple")
	//mq.PublishSimple("Hello world!")
	for i := 0; i < 20; i++ {
		mq.PublishSimple(fmt.Sprintf("%d", i))
	}

	fmt.Println("sent")
}
