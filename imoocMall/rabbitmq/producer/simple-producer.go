package main

import (
	"example/rabbitmq"
	"fmt"
)

func main() {
	mq := rabbitmq.NewRabbitMQSubscribe("" + "newSubscribe")
	//mq.PublishSimple("Hello world!")
	for i := 0; i < 20; i++ {
		mq.PublishSubscribe(fmt.Sprintf("%d", i))
		fmt.Println(i)
	}

	fmt.Println("sent")

}
