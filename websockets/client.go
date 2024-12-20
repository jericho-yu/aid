package websockets

import (
	"github.com/jericho-yu/aid/dict"
	"github.com/jericho-yu/aid/reflection"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type (
	Client struct {
		err                               error
		requestHeader                     http.Header
		uri                               url.URL
		groupName                         string
		name                              string
		conn                              *websocket.Conn
		status                            WebsocketConnStatus
		receiveMessageChan                chan []byte
		receiveMessageCallback            callbackFn
		returnCallback                    callbackFn
		asyncReceiveCallbackDict          *dict.AnyDict[string, callbackFn]
		onReceiveMessageSuccess           standardSuccessFn
		onCloseFail, onReceiveMessageFail standardFailFn
		syncMessageTimeout                time.Duration
	}
)

// NewClient 实例化：链接
func NewClient(groupName, name string, uri url.URL, options ...any) (client *Client, err error) {
	if reflection.New(uri).IsZero {
		return nil, NewWebsocketConnOptionErr("链接地址为空")
	}
	if name == "" {
		return nil, NewWebsocketConnOptionErr("链接名称为空")
	}

	client = &Client{
		uri:                      uri,
		groupName:                groupName,
		name:                     name,
		conn:                     &websocket.Conn{},
		status:                   Offline,
		receiveMessageChan:       make(chan []byte, 1),
		receiveMessageCallback:   nil,
		returnCallback:           nil,
		asyncReceiveCallbackDict: dict.MakeAnyDict[string, callbackFn](),
		syncMessageTimeout:       5 * time.Second,
	}

	if len(options) > 0 {
		for i := 0; i < len(options); i++ {
			if v, ok := options[i].(http.Header); ok {
				client.requestHeader = v
			}
		}
	}

	return
}

// GetStatus 获取链接状态
func (my *Client) GetStatus() WebsocketConnStatus { return my.status }

// GetName 获取链接名称
func (my *Client) GetName() string { return my.name }

// GetUri 获取链接地址
func (my *Client) GetUri() url.URL { return my.uri }

// GetConn 获取链接本体
func (my *Client) GetConn() *websocket.Conn { return my.conn }

// GetRequestHeader 获取请求头
func (my *Client) GetRequestHeader() http.Header { return my.requestHeader }

// SetRequestHeader 设置请求头
func (my *Client) SetRequestHeader(header http.Header) *Client {
	my.requestHeader = header
	return my
}

// SetOnCloseFail 设置关闭链接失败回调
func (my *Client) SetOnCloseFail(fn standardFailFn) *Client {
	my.onCloseFail = fn
	return my
}

// SetOnReceiveMessageSuccess 设置接收消息成功回调
func (my *Client) SetOnReceiveMessageSuccess(fn standardSuccessFn) *Client {
	my.onReceiveMessageSuccess = fn
	return my
}

// SetOnReceiveMessageFail 设置接收消息失败回调
func (my *Client) SetOnReceiveMessageFail(fn standardFailFn) *Client {
	my.onReceiveMessageFail = fn
	return my
}

// Conn 启动链接，并打开监听
func (my *Client) Conn(fn ...standardSuccessFn) *Client {
	var receiveMessage []byte

	my.conn, _, my.err = websocket.DefaultDialer.Dial(my.uri.String(), my.requestHeader)
	if my.err != nil {
		return my
	}

	// 执行链接成功回调
	if len(fn) > 0 {
		fn[0](my.groupName, my.name, my.conn)
	}

	// 开启监听
	go func(client *Client) {
		for {
			_, receiveMessage, client.err = client.conn.ReadMessage()
			if client.err != nil {
				if client.onReceiveMessageFail != nil {
					client.onReceiveMessageFail(client.groupName, client.name, client.conn, client.err)
				}
				return
			}

			// 解析消息
			message := ParseMessage(receiveMessage)
			client.receiveMessageCallback(client.groupName, client.name, message.GetMessage())
			if message.GetAsync() { // 异步消息
				if callback, ok := client.asyncReceiveCallbackDict.Get(message.GetMessageId()); ok {
					callback(client.groupName, client.name, message.GetMessage())
				}
			} else { // 同步消息
				client.receiveMessageChan <- message.GetMessage()
			}
		}
	}(my)

	my.status = Online

	return my
}

// AsyncMessage 发送消息：异步
func (my *Client) AsyncMessage(message []byte, fn callbackFn) *Client {
	msg := NewMessage(true, message)

	my.asyncReceiveCallbackDict.Set(msg.GetMessageId(), fn) // 配置异步回调

	my.err = my.conn.WriteMessage(websocket.TextMessage, msg.GetMessage()) // 发送消息

	return my
}

// SyncMessage 发送消息：同步
func (my *Client) SyncMessage(message []byte, options ...any) ([]byte, error) {
	var (
		timeout = my.syncMessageTimeout
		msg     = NewMessage(false, message)
	)

	if my.conn == nil || my.status == Offline {
		return nil, NewWebsocketOfflineErr("链接断开不在线")
	}

	my.err = my.conn.WriteMessage(websocket.TextMessage, msg.GetMessage()) // 发送消息
	if my.err != nil {
		return nil, my.err
	}

	for i := 0; i < len(options); i++ {
		if v, ok := options[i].(time.Duration); ok && v > 0 {
			timeout = v
		}
	}

	timeoutTimer := time.After(timeout)

	select {
	case receiveMessage := <-my.receiveMessageChan:
		return receiveMessage, nil
	case <-timeoutTimer:
		return nil, NewSyncMessageTimeoutErr("消息超时")
	}
}

// Close 关闭链接
func (my *Client) Close() *Client {
	if my.conn != nil && my.status == Online {
		my.err = my.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		if my.err != nil {
			if my.onCloseFail != nil {
				my.onCloseFail(my.groupName, my.name, my.conn, my.err)
			}
			my.status = Online
		} else {
			my.err = my.conn.Close()
			if my.err != nil {
				if my.onCloseFail != nil {
					my.onCloseFail(my.groupName, my.name, my.conn, my.err)
				}
				my.status = Online
			} else {
				my.status = Offline
				close(my.receiveMessageChan) // 关闭同步消息通道
			}
		}
	} else {
		my.conn = nil
		my.status = Offline
		close(my.receiveMessageChan)
	}

	return my
}

// Error 获取错误
func (my *Client) Error() error { return my.err }
