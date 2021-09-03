package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const MQURL = "amqp://imooc:9@localhost:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	// queue name
	QueueName string
	Exchange  string
	// binding key
	Key   string
	Mqurl string
}

// NewRabbitMQ create general instance
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	mq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	// create connection
	mq.conn, err = amqp.Dial(mq.Mqurl)
	mq.failOnErr(err, "connection failed")
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "get channel error")
	return mq
}

// Destroy close RabbitMQ
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s, %s", message, err))
	}
}

/* Simple Mode */

// NewRabbitMQSimple create RabbitMQ instance in simple mode
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// PublishSimple is for producer publishing messages
func (r *RabbitMQ) PublishSimple(message string) {
	// apply for a queue, if not exist
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 持久化
		false, // 自动删除
		false,
		false, // 是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send message to queue
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // mandatory 如果为true则会根据exchange类型和routekey寻找符合条件的队列，如果无法找到，则会把消息退回给发送者
		false, // immediate 如果为true，当exchange 发送消息到队列后发现队列上没有绑定消费者，则退回
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) ConsumeSimple() {
	// apply for a queue, if not exist
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 持久化
		false, // 自动删除
		false,
		false, // 是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    // 区分多个consumer
		true,  // 是否自动应答
		false, // exclusive true:只有自己可见
		false, // noLocal true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 实现要处理的函数
			log.Printf("received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] waiting for messages...")
	<-forever
}

/* Work Mode, a message can only be consumed by one consumer
e.g: When producing rate is much higher than consuming rate.
*/

/* Subscribe Mode */

func NewRabbitMQSubscribe(exchange string) *RabbitMQ {
	return NewRabbitMQ("", exchange, "")
}

func (r *RabbitMQ) PublishSubscribe(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // exchange 类型
		true,
		false,
		false, // internal 表示这个exchange不可以被client用来推送消息，仅用来进行exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnErr(err, "failed to send messages")
}

func (r *RabbitMQ) ConsumeSubscribe() {
	// create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // exchange 类型
		true,
		false,
		false, // internal 表示这个exchange不可以被client用来推送消息，仅用来进行exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	// create queue
	q, err := r.channel.QueueDeclare(
		"", // random queue name
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	// bind queue and exchange
	err = r.channel.QueueBind(
		q.Name,
		"", // in subscribe method, key is ""
		r.Exchange,
		false,
		nil,
	)
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	<-forever
}
