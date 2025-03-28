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
		log.Printf("收到消息OK：%s", message.Content)
		return nil
	})
	consumer.Start()

	select {}
}
