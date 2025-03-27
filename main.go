package main

import (
	"log"

	"github.com/jericho-yu/aid/messageQueue/rabbit"
)

func main() {
	r := rabbit.RabbitApp.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
	defer func() { _ = r.Close() }()

	r.NewQueue("message")
	consumer := r.Consume("message", "", func(prototypeMessage []byte) error {
		message := rabbit.MessageApp.Parse(prototypeMessage)
		log.Printf("收到消息：%s -> %s", message.Content, prototypeMessage)
		return nil
	})
	msgs := consumer.Go()
	go func() {
		for msg := range msgs {
			// 消费消息
			log.Printf("收到消息：%s", msg.Body)
			msg.Ack(true)
		}
	}()

	// 模拟程序运行一段时间后停止消费者
	select {
	case <-msgs:
		log.Println("消息通道已关闭")
	default:
		log.Println("停止消费者...")
		consumer.Stop()
	}
}
