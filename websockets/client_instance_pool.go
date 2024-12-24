package websockets

import (
	"sync"

	"github.com/jericho-yu/aid/dict"
)

type ClientInstancePool struct {
	pool *dict.AnyDict[string, *ClientInstance]
}

var (
	clientInstancePoolOnce sync.Once
	clientInstancePool     *ClientInstancePool
)

func OnceClientInstancePool() *ClientInstancePool {
	clientInstancePoolOnce.Do(func() { clientInstancePool = &ClientInstancePool{pool: dict.MakeAnyDict[string, *ClientInstance]()} })

	return clientInstancePool
}

// Append 增加客户端
func (*ClientInstancePool) Append(clientInstance *ClientInstance) error {
	if clientInstance.pool.Has(clientInstance.name) {
		return WebsocketClientExistErr
	}

	clientInstancePool.pool.Set(clientInstance.name, clientInstance)

	return nil
}

// Remove 删除客户端
func (*ClientInstancePool) Remove(name string) error {
	if !clientInstance.pool.Has(name) {
		return WebsocketClientNotExistErr
	}

	clientInstance.pool.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (*ClientInstancePool) Get(name string) (*ClientInstance, error) {
	if !clientInstance.pool.Has(name) {
		return nil, WebsocketClientNotExistErr
	}

	clientInstance, _ := clientInstancePool.pool.Get(name)

	return clientInstance, nil
}

// Has 检查客户端是否存在
func (*ClientInstancePool) Has(name string) bool {
	return clientInstancePool.pool.Has(name)
}

// Close 关闭客户端
func (*ClientInstancePool) Close(name string) error {
	if clientInstance, err := clientInstance.Get(name); err != nil {
		return err
	} else {
		return clientInstance.Close().Error()
	}
}
