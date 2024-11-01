package websocketPool

import (
	"context"
	"errors"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type (
	// Client websocket 客户端链接
	Client struct {
		url                url.URL
		InstanceName, Name string
		Conn               *websocket.Conn
		mu                 sync.Mutex    // 同步锁
		closeChan          chan struct{} // 关闭信号
		syncChan           chan []byte   // 同步消息
		onReceiveMsg       func(instanceName, clientName string, prototypeMsg []byte) ([]byte, error)
		heart              *Heart
		timeout            *MessageTimeout
	}

	// PendingRequest 待处理请求
	PendingRequest struct {
		Uuid uuid.UUID
		Chan chan []byte
		Done chan struct{}
		Err  error
	}
)

// NewClient 实例化：websocket 客户端链接
func NewClient(instanceName, name, host, path string, receiveMessageFunc func(instanceName, clientName string, prototypeMsg []byte) ([]byte, error)) (*Client, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   path,
	}
	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		InstanceName: instanceName,
		Name:         name,
		url:          u,
		Conn:         client,
		onReceiveMsg: receiveMessageFunc,
	}, nil
}

// SendMsg 发送消息：通过链接
func (my *Client) SendMsg(msgType int, msg []byte) ([]byte, error) {
	var (
		err error
		res = make([]byte, 0)
	)
	my.syncChan = make(chan []byte)

	if my.timeout == nil || my.timeout.interval == 0 {
		clientPoolIns.Error = errors.New("同步消息，需要设置超时时间")
		return nil, errors.New("同步消息，需要设置超时时间")
	}

	my.mu.Lock()
	defer my.mu.Unlock()

	err = my.Conn.WriteMessage(msgType, msg)
	if err != nil {
		if clientPoolIns.onSendMsgErr != nil {
			clientPoolIns.onSendMsgErr(my.InstanceName, my.Name, err)
		}
		clientPoolIns.Error = err
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), my.timeout.interval)
	defer cancel()
	select {
	case <-ctx.Done():
		clientPoolIns.Error = errors.New("请求超时")
		return nil, errors.New("请求超时")
	case res = <-my.syncChan:
		if my.onReceiveMsg != nil {
			return my.onReceiveMsg(my.InstanceName, my.Name, res)
		}
		return res, nil
	}
}

// Close 关闭链接
func (my *Client) Close() error {
	var err error

	if err = my.Conn.Close(); err != nil {
		if clientPoolIns.onCloseErr != nil {
			clientPoolIns.onCloseErr(my.InstanceName, my.Name, err)
		}
		my.closeChan <- struct{}{}
		return err
	}

	return nil
}
