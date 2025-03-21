package config

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func InitRabbitMQ(msg string) {
	//1.建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "连接RabbitMQ失败")
	defer conn.Close()
	//2.创建channel
	channel, err := conn.Channel()
	failOnError(err, "创建channel失败")
	//3.简单模式中需要创建队列
	q, err := channel.QueueDeclare(
		"hello", // name
		true,    // durable true表示持久，false 表示不持久化
		false,   // delete when unused 是否自动删除队列，如果设置为true ，则当最后一个消费者取消订阅时会自动删除队列
		false,   // exclusive 是否独享队列
		false,   // no-wait 等待服务器确认
		nil,     // arguments 额外参数
	)
	//4.发送消息
	body := "this is a rabbitMQ Message :" + msg
	channel.PublishWithContext(
		context.Background(),
		"",
		q.Name, //队列名称
		false,  // false表示如果交换机无法找到一个符合条件的队列，就把消息丢弃掉
		false,  //是否立即被消费者接收 false 不需要立即被消费者接收
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
