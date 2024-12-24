package websockets

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/jericho-yu/aid/operation"
	"strings"
)

type Message struct {
	async            bool
	messageId        string
	message          []byte
	prototypeMessage []byte
}

// NewMessage 新建消息
func NewMessage(async bool, message []byte) Message {
	u := uuid.Must(uuid.NewV6()).String()
	b := bytes.Buffer{}
	b.Write([]byte(u))
	b.WriteByte(':')
	b.Write(message)
	return Message{
		async:            async,
		messageId:        operation.Ternary[string](async, u, ""),
		message:          operation.Ternary[[]byte](async, b.Bytes(), message),
		prototypeMessage: message,
	}
}

// ParseMessage 解析消息
func ParseMessage(prototypeMessage []byte) Message {
	var (
		messages = strings.Split(string(prototypeMessage), ":")
		wm       = Message{}
	)

	if len(messages) == 2 {
		wm.messageId = messages[0]
		wm.message = []byte(messages[1])
		wm.prototypeMessage = prototypeMessage
		wm.async = true
	} else {
		wm.message = prototypeMessage
		wm.prototypeMessage = prototypeMessage
	}

	return wm
}

// GetAsync 获取同步类型
func (my *Message) GetAsync() bool { return my.async }

// GetMessageId 获取消息编号
func (my *Message) GetMessageId() string { return my.messageId }

// GetMessage 获取消息
func (my *Message) GetMessage() []byte {
	return operation.Ternary[[]byte](my.async, my.message, my.prototypeMessage)
}

// GetPrototypeMessage 获取原始消息
func (my *Message) GetPrototypeMessage() []byte { return my.prototypeMessage }