package websocketPool

import (
	"sync"

	"github.com/gorilla/websocket"
)

type msgType struct{}

var MsgType msgType

func (msgType) Text() int {
	return websocket.TextMessage
}

func (msgType) Binary() int {
	return websocket.BinaryMessage
}

func (msgType) Close() int {
	return websocket.CloseMessage
}

func (msgType) Ping() int {
	return websocket.PingMessage
}

func (msgType) Pong() int {
	return websocket.PongMessage
}

var (
	ClientPoolApp     ClientPool
	ClientInstanceApp ClientInstance
	clientApp         Client
	clientPoolIns     *ClientPool
	clientPoolOnce    sync.Once
	HeartOpt          Heart
	MessageTimeoutOpt MessageTimeout
)
