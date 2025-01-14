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
	if clientInstance.connections.Has(clientInstance.name) {
		return WebsocketClientExistErr
	}

	clientInstancePool.pool.Set(clientInstance.name, clientInstance)

	return nil
}

// Remove 删除客户端
func (*ClientInstancePool) Remove(name string) error {
	if !clientInstance.connections.Has(name) {
		return WebsocketClientNotExistErr
	}

	clientInstance.connections.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (*ClientInstancePool) Get(name string) (*ClientInstance, error) {
	if !clientInstance.connections.Has(name) {
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
		err = clientInstance.Close().Error()
		clientInstancePool.pool.RemoveByKey(clientInstance.name)
		return err
	}
}

// Clean 清空客户端实例
func (*ClientInstancePool) Clean() []error {
	var errorList []error
	for _, clientInstance := range clientInstancePool.pool.All() {
		errTmp := clientInstance.Clean()
		if len(errTmp) > 0 {
			errorList = append(errorList, errTmp...)
		} else {
			clientInstance.connections.RemoveByKey(clientInstance.name)
		}
	}

	return errorList
}
