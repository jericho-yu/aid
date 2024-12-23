package websockets

import (
	"errors"
)

var (
	WebsocketConnOptionErr = errors.New("websocket链接错误：参数错误")
	SyncMessageTimeoutErr  = errors.New("同步消息超时")
	WebsocketOfflineErr    = errors.New("链接不在线")
	PingErr                = errors.New("ping失败")
)
