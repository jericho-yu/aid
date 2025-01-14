package websockets

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jericho-yu/aid/dict"
)

type (
	ServerPool struct {
		connections             *dict.AnyDict[string, *Server]
		authIdToAddr            *dict.AnyDict[string, string]
		addrToAuth              *dict.AnyDict[string, string]
		onConnectionFail        serverConnectionFailFn
		onConnectionSuccess     serverConnectionSuccessFn
		onSendMessageSuccess    serverSendMessageSuccessFn
		onSendMessageFail       serverSendMessageFailFn
		onReceiveMessageFail    serverReceiveMessageFailFn
		onReceiveMessageSuccess serverReceiveMessageSuccessFn
	}
)

var (
	serverPoolOnce sync.Once
	serverPool     *ServerPool
)

// OnceServer 单例化：websocket服务端
func OnceServer(serverCallbackConfig ServerCallbackConfig) *ServerPool {
	serverPoolOnce.Do(func() {
		serverPool = &ServerPool{
			connections:             dict.MakeAnyDict[string, *Server](),
			authIdToAddr:            dict.MakeAnyDict[string, string](),
			addrToAuth:              dict.MakeAnyDict[string, string](),
			onConnectionFail:        serverCallbackConfig.OnConnectionFail,
			onConnectionSuccess:     serverCallbackConfig.OnConnectionSuccess,
			onSendMessageSuccess:    serverCallbackConfig.OnSendMessageSuccess,
			onSendMessageFail:       serverCallbackConfig.OnSendMessageFail,
			onReceiveMessageFail:    serverCallbackConfig.OnReceiveMessageFail,
			onReceiveMessageSuccess: serverCallbackConfig.OnReceiveMessageSuccess,
		}
	})

	return serverPool
}

// appendConn 增加连接
func (*ServerPool) appendConn(authId *string, conn *websocket.Conn) {
	server := NewServer(conn)
	serverPool.authIdToAddr.Set(*authId, conn.RemoteAddr().String())
	serverPool.addrToAuth.Set(conn.RemoteAddr().String(), *authId)
	serverPool.connections.Set(conn.RemoteAddr().String(), server)
}

// SendMessageByAddr 发送消息：通过地址
func (*ServerPool) SendMessageByAddr(addr *string, prototypeMessage []byte) {
	if server, ok := serverPool.connections.Get(*addr); ok {
		server.AsyncMessage(prototypeMessage, serverPool.onSendMessageSuccess, serverPool.onSendMessageFail)
	} else {
		if serverPool.onSendMessageFail != nil {
			serverPool.onSendMessageFail(fmt.Errorf("没有找到连接：%s", *addr))
		}
	}
}

// SendMessageByAuthId 发送消息：通过认证ID
func (*ServerPool) SendMessageByAuthId(authId *string, prototypeMessage []byte) {
	serverPool.authIdToAddr.GetKeysByValue(*authId).Each(func(idx int, item string) {
		if server, ok := serverPool.connections.Get(item); ok {
			server.AsyncMessage(prototypeMessage, serverPool.onSendMessageSuccess, serverPool.onSendMessageFail)
		} else {
			if serverPool.onSendMessageFail != nil {
				serverPool.onSendMessageFail(fmt.Errorf("没有找到连接：%s -> %s", *authId, item))
			}
		}
	})
}

// SetOnConnectionSuccess 设置回调：当连接成功
func (*ServerPool) SetOnConnectionSuccess(onConnectionSuccess serverConnectionSuccessFn) *ServerPool {
	serverPool.onConnectionSuccess = onConnectionSuccess
	return serverPool
}

// SetOnConnectionFail 设置回调：当连接失败
func (*ServerPool) SetOnConnectionFail(onConnectionFail serverConnectionFailFn) *ServerPool {
	serverPool.onConnectionFail = onConnectionFail
	return serverPool
}

// SetOnSendMessageSuccess 设置回调：当发送消息成功
func (*ServerPool) SetOnSendMessageSuccess(onSendMessageSuccess serverSendMessageSuccessFn) *ServerPool {
	serverPool.onSendMessageSuccess = onSendMessageSuccess
	return serverPool
}

// SetOnSendMessageFail 设置回调：当发送消息失败
func (*ServerPool) SetOnSendMessageFail(onSendMessageFail serverSendMessageFailFn) *ServerPool {
	serverPool.onSendMessageFail = onSendMessageFail
	return serverPool
}

// SetOnReceiveMessageSuccess 设置回调：当接收消息成功
func (*ServerPool) SetOnReceiveMessageSuccess(onReceiveMessageSuccess serverReceiveMessageSuccessFn) *ServerPool {
	serverPool.onReceiveMessageSuccess = onReceiveMessageSuccess
	return serverPool
}

// SetOnReceiveMessageFail 设置回调：当接收消息失败
func (*ServerPool) SetOnReceiveMessageFail(onReceiveMessageFail serverReceiveMessageFailFn) *ServerPool {
	serverPool.onReceiveMessageFail = onReceiveMessageFail
	return serverPool
}

// Handle 消息处理
func (my *ServerPool) Handle(
	writer http.ResponseWriter,
	req *http.Request,
	header http.Header,
	condition serverConnectionCheckFn,
) {
	var (
		err  error
		conn *websocket.Conn
	)

	// 升级协议
	conn, err = upgrader.Upgrade(writer, req, header)
	if err != nil {
		if serverPool.onConnectionFail != nil {
			serverPool.onConnectionFail(err)
		}
	}

	// 验证连接
	if condition != nil {
		if err = condition(header); err != nil && serverPool.onConnectionFail != nil {
			serverPool.onConnectionFail(err)
			return
		}
	}

	// 加入连接池
	authId := header.Get("Auth-Id")
	my.appendConn(&authId, conn)

	// 开启接收消息

}
