package websocketPool

import "time"

type (
	// Heart 链接心跳
	Heart struct {
		ticker *time.Ticker
		fn     func(*Client)
	}

	// MessageTimeout 通信超时
	MessageTimeout struct {
		interval time.Duration
	}
)

// New 实例化：链接心跳
func (Heart) New() *Heart {
	return &Heart{}
}

// SetInterval 设置定时器
func (r *Heart) SetInterval(interval time.Duration) *Heart {
	if r.ticker != nil {
		r.ticker.Reset(interval)
	} else {
		r.ticker = time.NewTicker(interval)
	}
	return r
}

// SetFn 设置回调：定时器执行内容
func (r *Heart) SetFn(fn func(client *Client)) *Heart {
	r.fn = fn
	return r
}

// Stop 停止定时器
func (r *Heart) Stop() *Heart {
	r.ticker.Stop()
	return r
}

// Default 默认心跳：10秒
func (Heart) Default() *Heart {
	return HeartOpt.New().SetInterval(time.Second * 10).SetFn(func(client *Client) {
		_, _ = client.SendMsg(MsgType.Ping(), []byte("ping"))
	})
}

// New 实例化：链接超时
func (MessageTimeout) New() *MessageTimeout {
	return &MessageTimeout{}
}

// SetInterval 设置定时器时间
func (r *MessageTimeout) SetInterval(interval time.Duration) *MessageTimeout {
	r.interval = interval
	return r
}

// Default 默认消息超时：5秒
func (MessageTimeout) Default() *MessageTimeout {
	return MessageTimeoutOpt.New().SetInterval(time.Second * 5)
}
