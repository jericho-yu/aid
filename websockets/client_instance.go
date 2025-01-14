package websockets

import (
	"sync"

	"github.com/jericho-yu/aid/dict"
)

type ClientInstance struct {
	name        string
	connections *dict.AnyDict[string, *Client]
}

var (
	clientInstanceOnce sync.Once
	clientInstance     *ClientInstance
)

func OnceClientInstance(name string) *ClientInstance {
	clientInstanceOnce.Do(func() { clientInstance = &ClientInstance{name: name, connections: dict.MakeAnyDict[string, *Client]()} })

	return clientInstance
}

// Append 增加客户端
func (*ClientInstance) Append(client *Client) error {
	if clientInstance.connections.Has(client.name) {
		return WebsocketClientExistErr
	}

	clientInstance.connections.Set(client.name, client)

	return nil
}

// Remove 删除客户端
func (*ClientInstance) Remove(name string) error {
	if !clientInstance.connections.Has(name) {
		return WebsocketClientNotExistErr
	}

	clientInstance.connections.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (*ClientInstance) Get(name string) (*Client, error) {
	if !clientInstance.connections.Has(name) {
		return nil, WebsocketClientNotExistErr
	}

	client, _ := clientInstance.connections.Get(name)

	return client, nil
}

// Has 检查客户端是否存在
func (*ClientInstance) Has(name string) bool {
	return clientInstance.connections.Has(name)
}

// Close 关闭客户端
func (*ClientInstance) Close(name string) error {
	if client, err := clientInstance.Get(name); err != nil {
		return err
	} else {
		err = client.Close().Error()
		clientInstance.connections.RemoveByKey(client.name)
		return err
	}
}

// Clean 清空客户端
func (*ClientInstance) Clean() []error {
	var errorList []error
	for _, client := range clientInstance.connections.All() {
		if err := client.Close().Error(); err != nil {
			errorList = append(errorList, err)
		} else {
			clientInstance.connections.RemoveByKey(client.name)
		}
	}

	return errorList
}
