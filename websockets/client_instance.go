package websockets

import (
	"github.com/jericho-yu/aid/dict"
)

type ClientInstance struct {
	name        string
	connections *dict.AnyDict[string, *Client]
}

// NewClientInstance 实例化：websocket客户端实例
func NewClientInstance(name string) *ClientInstance {
	return &ClientInstance{name: name, connections: dict.MakeAnyDict[string, *Client]()}
}

// Append 增加客户端
func (my *ClientInstance) Append(client *Client) error {
	if my.connections.Has(client.name) {
		return WebsocketClientExistErr
	}

	my.connections.Set(client.name, client)

	return nil
}

// Remove 删除客户端
func (my *ClientInstance) Remove(name string) error {
	if !my.connections.Has(name) {
		return WebsocketClientNotExistErr
	}

	my.connections.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (my *ClientInstance) Get(name string) (*Client, error) {
	if !my.connections.Has(name) {
		return nil, WebsocketClientNotExistErr
	}

	client, _ := my.connections.Get(name)

	return client, nil
}

// Has 检查客户端是否存在
func (my *ClientInstance) Has(name string) bool {
	return my.connections.Has(name)
}

// Close 关闭客户端
func (my *ClientInstance) Close(name string) error {
	if client, err := my.Get(name); err != nil {
		return err
	} else {
		err = client.Close().Error()
		my.connections.RemoveByKey(client.name)
		return err
	}
}

// Clean 清空客户端
func (my *ClientInstance) Clean() []error {
	var errorList []error
	for _, client := range my.connections.All() {
		if err := client.Close().Error(); err != nil {
			errorList = append(errorList, err)
		} else {
			my.connections.RemoveByKey(client.name)
		}
	}

	return errorList
}
