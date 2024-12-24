package websockets

import "github.com/gorilla/websocket"

type (
	callbackFn func(groupName, name string, message []byte)

	standardSuccessFn func(groupName, name string, conn *websocket.Conn)

	standardFailFn func(groupName, name string, conn *websocket.Conn, err error)

	receiveMessageSuccessFn func(groupName, name string, prototypeMessage []byte)

	heartFn func(groupName, name string, client *Client)

	timeoutFn func(groupName, name string, prototypeMessage []byte)

	pingFn func(conn *websocket.Conn) error
)
