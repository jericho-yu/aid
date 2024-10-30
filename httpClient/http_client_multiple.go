package httpClient

import (
	"sync"
)

type HttpClientMultiple struct {
	clients []*HttpClient
}

// New New 实例化：批量请求对象
func (HttpClientMultiple) New() *HttpClientMultiple {
	return &HttpClientMultiple{}
}

// Add 添加httpClient对象
func (receiver *HttpClientMultiple) Add(hc *HttpClient) *HttpClientMultiple {
	receiver.clients = append(receiver.clients, hc)
	return receiver
}

// SetClients 设置httpClient对象
func (receiver *HttpClientMultiple) SetClients(clients []*HttpClient) *HttpClientMultiple {
	receiver.clients = clients
	return receiver
}

// Send 批量发送
func (receiver *HttpClientMultiple) Send() *HttpClientMultiple {
	if len(receiver.clients) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(receiver.clients))

		for _, client := range receiver.clients {
			go func(client *HttpClient) {
				defer wg.Done()

				client.Send()
			}(client)
		}

		wg.Wait()
	}

	return receiver
}

// GetClients 获取链接池
func (receiver *HttpClientMultiple) GetClients() []*HttpClient {
	return receiver.clients
}
