package httpClient

import (
	"sync"
)

type Multiple struct {
	clients []*Client
}

// NewHttpClientMultiple 实例化：批量请求对象
func NewHttpClientMultiple() *Multiple { return &Multiple{} }

// Add 添加httpClient对象
func (my *Multiple) Add(hc *Client) *Multiple {
	my.clients = append(my.clients, hc)
	return my
}

// SetClients 设置httpClient对象
func (my *Multiple) SetClients(clients []*Client) *Multiple {
	my.clients = clients
	return my
}

// Send 批量发送
func (my *Multiple) Send() *Multiple {
	if len(my.clients) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(my.clients))

		for _, client := range my.clients {
			go func(client *Client) {
				defer wg.Done()

				client.Send()
			}(client)
		}

		wg.Wait()
	}

	return my
}

// GetClients 获取链接池
func (my *Multiple) GetClients() []*Client {
	return my.clients
}
