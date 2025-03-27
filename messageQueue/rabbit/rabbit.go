package rabbit

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type (
	Rabbit struct {
		err         error
		username    string
		password    string
		host        string
		port        string
		virtualHost string
		conn        *amqp.Connection
		ch          *amqp.Channel
		queues      map[string]amqp.Queue
		mu          sync.RWMutex
	}
)

var RabbitApp Rabbit

// New 创建一个 Rabbit 实例
func (*Rabbit) New(
	username,
	password,
	host,
	port,
	virtualHost string,
) *Rabbit {
	ins := &Rabbit{
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		virtualHost: virtualHost,
		queues:      make(map[string]amqp.Queue),
	}

	// 连接到 RabbitMQ
	ins.conn, ins.err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", ins.username, ins.password, ins.host, ins.port, ins.virtualHost))
	if ins.err != nil {
		ins.err = ConnRabbitErr.Wrap(ins.err)
	}

	ins.NewChannel() // 创建频道

	return ins
}

// 获取链接
func (my *Rabbit) getConn() *amqp.Connection { return my.conn }

// GetConn 获取链接
func (my *Rabbit) GetConn() *amqp.Connection {
	my.mu.RLock()
	defer my.mu.RLock()

	return my.getConn()
}

// Error 获取错误
func (my *Rabbit) Error() error { return my.err }

// Close 关闭链接
func (my *Rabbit) Close() error {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.conn != nil {
		my.closeChannel()
		return my.conn.Close()
	}

	return nil
}

// NewChannel 创建频道
func (my *Rabbit) NewChannel() *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.ch, my.err = my.getConn().Channel()

	return my
}

// closeChannel 关闭频道
func (my *Rabbit) closeChannel() { my.err = my.ch.Close() }

// CloseChannel 关闭频道
func (my *Rabbit) CloseChannel() *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.closeChannel()

	return my
}

// NewQueue 创建队列
func (my *Rabbit) NewQueue(queueName string) *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.ch == nil {
		my.NewChannel()
	}

	if my.err != nil {
		return my
	}

	var queue amqp.Queue

	// 声明一个队列
	queue, my.err = my.ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 附加属性
	)
	if my.err != nil {
		fmt.Println("QueueDeclare err:", my.err)
	}

	my.queues[queue.Name] = queue

	return my
}

// 推送消息
func (my *Rabbit) Publish(queueName string, body string) *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.err != nil {
		return my
	}

	queue, exist := my.queues[queueName]
	if !exist {
		my.err = QueueNotExistErr.New(queueName)
		return my
	}

	// 发送消息
	my.err = my.ch.Publish(
		"",         // 默认交换机
		queue.Name, // 路由键，使用队列名称
		false,      // 是否立即发送
		false,      // 是否持久化
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if my.err != nil {
		my.err = PublishMessageErr.Wrap(my.err)
	}

	return my
}
