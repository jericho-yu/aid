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
func (my *HttpClientMultiple) Add(hc *HttpClient) *HttpClientMultiple {
	my.clients = append(my.clients, hc)
	return my
}

// SetClients 设置httpClient对象
func (my *HttpClientMultiple) SetClients(clients []*HttpClient) *HttpClientMultiple {
	my.clients = clients
	return my
}

// Send 批量发送
func (my *HttpClientMultiple) Send() *HttpClientMultiple {
	if len(my.clients) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(my.clients))

		for _, client := range my.clients {
			go func(client *HttpClient) {
				defer wg.Done()

				client.Send()
			}(client)
		}

		wg.Wait()
	}

	return my
}

// GetClients 获取链接池
func (my *HttpClientMultiple) GetClients() []*HttpClient {
	return my.clients
}
