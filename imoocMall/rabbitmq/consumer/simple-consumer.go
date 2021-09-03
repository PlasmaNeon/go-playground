package main

import "example/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQSubscribe("" + "newSubscribe")
	mq.ConsumeSubscribe()
}
