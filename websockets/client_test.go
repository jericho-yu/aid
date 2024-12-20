package websockets

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
	"time"
)

func onLine() (*Client, error) {
	client, err := NewClient("", "client-test1", url.URL{Scheme: "ws", Host: "127.0.0.1:12345"})
	if err != nil {
		return nil, err

	}
	if client == nil {
		return nil, errors.New("创建链接失败")
	}

	return client, nil
}

func offLine(client *Client) error {
	if err := client.Close().Error(); err != nil {
		return fmt.Errorf("关闭链接失败：%v", err)
	}
	return nil
}

// Test1Conn 测试：创建和关闭链接
func Test1Conn(t *testing.T) {
	t.Run("websocket客户端测试", func(t *testing.T) {
		client, err := onLine()
		if err != nil {
			t.Error(err)
		}

		client.Conn(func(groupName, name string, conn *websocket.Conn) {
			log.Printf("链接成功 -> [%s:%s]", groupName, name)
		})

		if err = offLine(client); err != nil {
			t.Error(err)
		}
	})
}

// Test2Sync 测试：同步消息
func Test2Sync(t *testing.T) {
	client, err := onLine()
	if err != nil {
		t.Error(err)
	}

	_, err = client.Conn().SyncMessage([]byte("hello"), time.Second) // 1秒超时
	if err != nil {
		if !errors.Is(err, NewSyncMessageTimeoutErr("")) {
			t.Errorf("发送消息失败：%v", err)
		}
	}

	if err = offLine(client); err != nil {
		t.Error(err)
	}
}
