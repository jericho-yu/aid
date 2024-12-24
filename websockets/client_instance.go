package websockets

import (
	"sync"

	"github.com/jericho-yu/aid/dict"
)

type ClientInstance struct {
	name string
	pool *dict.AnyDict[string, *Client]
}

var (
	clientInstanceOnce sync.Once
	clientInstance     *ClientInstance
)

func OnceClientInstance(name string) *ClientInstance {
	clientInstanceOnce.Do(func() { clientInstance = &ClientInstance{name: name, pool: dict.MakeAnyDict[string, *Client]()} })

	return clientInstance
}

// Append 增加客户端
func (*ClientInstance) Append(client *Client) error {
	if clientInstance.pool.Has(client.name) {
		return WebsocketClientExistErr
	}

	clientInstance.pool.Set(client.name, client)

	return nil
}

// Remove 删除客户端
func (*ClientInstance) Remove(name string) error {
	if !clientInstance.pool.Has(name) {
		return WebsocketClientNotExistErr
	}

	clientInstance.pool.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (*ClientInstance) Get(name string) (*Client, error) {
	if !clientInstance.pool.Has(name) {
		return nil, WebsocketClientNotExistErr
	}

	client, _ := clientInstance.pool.Get(name)

	return client, nil
}

// Has 检查客户端是否存在
func (*ClientInstance) Has(name string) bool {
	return clientInstance.pool.Has(name)
}

// Close 关闭客户端
func (*ClientInstance) Close(name string) error {
	if client, err := clientInstance.Get(name); err != nil {
		return err
	} else {
		return client.Close().Error()
	}
}
