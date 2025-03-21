package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	//1.建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//2.配置channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//3.配置队列  ，要和发布者统一起来
	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable true表示持久，false 表示不持久化
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//4.消费消息
	msgs, err := ch.Consume(q.Name, "", //消费者标签，用于区分不同的消费者
		true,  //自动确认(消费消息后，会从队列里直接删除) false的话 需要手动确认
		false, //独占，独享。这个队列只能一个消费者使用
		false, //如果为true, 表示生产者和消费者不能是同一个connect
		false, //是否阻塞。true表示阻塞。   阻塞：表示创建交换机的请求发送后，阻塞信息。 非阻塞： 不会阻塞等待RMQ Server的返回消息。（不推荐使用）
		nil,   //其它参数 直接写nil
	)
	failOnError(err, "Failed to register a consumer")

	//5.输出信息
	var forever chan struct{}
	go func() {
		for msg := range msgs { //-chan amqp091.Delivery (msgs也是一个管道)
			log.Printf("Received a message: %s", msg.Body)
			//msg.Ack(true) //如果上面配置了false这里就需要手动确认
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
