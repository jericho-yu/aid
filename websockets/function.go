package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	clientCallbackFn              func(groupName, name string, message []byte)
	clientStandardSuccessFn       func(groupName, name string, conn *websocket.Conn)
	clientStandardFailFn          func(groupName, name string, conn *websocket.Conn, err error)
	clientReceiveMessageSuccessFn func(groupName, name string, prototypeMessage []byte)
	heartFn                       func(groupName, name string, client *Client)
	pingFn                        func(conn *websocket.Conn) error
	serverConnConditionFn         func() (string, error)
	serverReceiveMessageSuccessFn func(prototypeMessage []byte, ws *websocket.Conn)
	serverReceiveMessageFailFn    func(conn *websocket.Conn, err error)
	serverReceivePingFn           func(ws *websocket.Conn)
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
