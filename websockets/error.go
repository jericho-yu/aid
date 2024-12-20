package websockets

import (
	"errors"
	"github.com/jericho-yu/aid/exception"
	"strings"
)

type (
	WebsocketConnOptionErr error

	SyncMessageTimeoutErr error

	WebsocketOfflineErr error

	ErrWebsocketNewOption exception.Exception
)

func NewErrWebsocketNewOption(values ...any) *ErrWebsocketNewOption {
	e := &ErrWebsocketNewOption{}
	for i := 0; i < len(values); i++ {
		if v, ok := values[i].(error); ok {
			e.Err = v
		}
		if v, ok := values[i].(string); ok {
			e.Msg = v
		}
	}

	return e
}

func (my *ErrWebsocketNewOption) Error() string {
	b := strings.Builder{}

	if my.Msg != "" {
		b.WriteString(my.Msg)
		b.WriteByte('：')
	}
	b.WriteString(my.Err.Error())

	return b.String()
}

// 实现 Is 方法，用于 errors.Is 的类型比较
func (e *MyError) Is(target error) bool {
	_, ok := target.(*MyError)
	return ok
}

// NewWebsocketConnOptionErr websocket链接错误
func NewWebsocketConnOptionErr(msg string) error {
	return errors.New(msg)
}

// NewSyncMessageTimeoutErr 同步消息超时错误
func NewSyncMessageTimeoutErr(msg string) error {
	return errors.New(msg)
}

// NewWebsocketOfflineErr 链接不在线
func NewWebsocketOfflineErr(msg string) error {
	return errors.New(msg)
}
