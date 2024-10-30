package websocketPool

import (
	"errors"

	"github.com/jericho-yu/aid/dict"
)

// ClientInstance websocket 客户端链接实例
type ClientInstance struct {
	Name    string
	Clients *dict.AnyDict[string, *Client]
}

// New 实例化：websocket 客户端实例
func (ClientInstance) New(instanceName string) *ClientInstance {
	return &ClientInstance{Name: instanceName, Clients: dict.MakeAnyDict[string, *Client]()}
}

// GetClient 获取websocket客户端链接
func (r *ClientInstance) GetClient(clientName string) (*Client, bool) {
	websocketClient, exist := r.Clients.Get(clientName)
	if !exist {
		return nil, exist
	}
	return websocketClient, true
}

// SetClient 创建新链接
func (r *ClientInstance) SetClient(
	clientName, host, path string,
	receiveMessageFn func(instanceName, clientName string, propertyMessage []byte) ([]byte, error),
	heart *Heart,
	timeout *MessageTimeout,
) (*Client, error) {
	var (
		err          error
		exist        bool
		client       *Client
		prototypeMsg []byte
	)

	if heart == nil {
		return nil, errors.New("心跳设置不能为空")
	}

	client, exist = r.Clients.Get(clientName)
	if exist {
		if err = client.Conn.Close(); err != nil {
			return nil, err
		}
		r.Clients.RemoveByKey(clientName)
	}

	if client, err = NewClient(r.Name, clientName, host, path, receiveMessageFn); err != nil {
		return nil, err
	}
	r.Clients.Set(clientName, client)

	if clientPoolIns.onConnect != nil {
		clientPoolIns.onConnect(r.Name, clientName)
	}

	client.heart = heart
	if timeout != nil {
		client.timeout = timeout
	}

	// 开启协程：接收消息
	go func() {
		for {
			select {
			case <-client.closeChan:
				// 关闭链接
				client.heart.ticker.Stop()
				r.Clients.RemoveByKey(clientName)
				return
			case <-client.heart.ticker.C:
				// 执行心跳
				if client.heart.fn != nil {
					client.heart.fn(client)
				}
			default:
				_, prototypeMsg, err = client.Conn.ReadMessage()
				if err != nil {
					if clientPoolIns.onReceiveMsgErr != nil {
						clientPoolIns.onReceiveMsgErr(r.Name, clientName, prototypeMsg, err)
					}
					return
				}

				client.syncChan <- prototypeMsg
			}
		}
	}()

	return client, nil
}

// SendMsgByName 发送消息：通过名称
func (r *ClientInstance) SendMsgByName(clientName string, msgType int, msg []byte) ([]byte, error) {
	var (
		exist  bool
		client *Client
	)

	client, exist = r.Clients.Get(clientName)
	if !exist {
		if clientPoolIns.onSendMsgErr != nil {
			clientPoolIns.onSendMsgErr(r.Name, clientName, errors.New("没有找到客户端链接"))
		}
	}

	return client.SendMsg(msgType, msg)
}

// Close 关闭客户端实例
func (r *ClientInstance) Close() {
	for _, conn := range r.Clients.All() {
		_ = conn.Close()
	}

	r.Clients.Clean()
	clientPoolIns.clientInstances.RemoveByKey(r.Name)
}
