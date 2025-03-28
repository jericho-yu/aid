package main

import (
	"log"

	"github.com/jericho-yu/aid/httpClient"
	"github.com/jericho-yu/aid/messageQueue/rabbit"
)

func a() {
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
func main() {
	hc := httpClient.
		App.
		NewPost("http://127.0.0.1:9000/team/all").
		SetHeaderContentType(httpClient.ContentTypeJson).
		SetHeaderAccept(httpClient.AcceptJson).
		Send()
	if hc.Err != nil {
		log.Fatalf("错误：%v", hc.Err)
	}

	log.Printf("响应体：%s", hc.GetResponseRawBody())
}
