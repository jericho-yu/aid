package rabbit

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// 链接
	rabbit := RabbitApp.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
	defer rabbit.CloseChannel().Close()

	pool := PoolApp.Once().Set("default", rabbit)

	rabbit = pool.Get("default")
	if rabbit == nil {
		t.Fatalf("没有找到链接：%s", "default")
	}

	rabbit.NewChannel()
	if rabbit.Error() != nil {
		t.Fatalf("创建channel失败：%v", rabbit.Error())
	}

	rabbit.NewQueue("message")
	if rabbit.Error() != nil {
		t.Fatalf("创建队列失败：%v", rabbit.Error())
	}

	rabbit.Publish("message", "hello world"+time.Now().Format(time.DateTime))
}
