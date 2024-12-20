package websockets

import "github.com/gorilla/websocket"

type (
	callbackFn func(groupName, name string, message []byte)

	standardSuccessFn func(groupName, name string, conn *websocket.Conn)

	standardFailFn func(groupName, name string, conn *websocket.Conn, err error)
)
