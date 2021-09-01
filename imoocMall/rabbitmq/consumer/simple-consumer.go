package main

import "example/rabbitmq"

func main() {
	mq := rabbitmq.NewRabbitMQSimple("imoocSimple")
	mq.ConsumeSimple()
}
