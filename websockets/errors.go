package websockets

import (
	"fmt"
	"github.com/jericho-yu/aid/operation"
	"reflect"

	"github.com/jericho-yu/aid/myError"
)

type (
	WebsocketConnOption                                 struct{ myError.MyError }
	SyncMessageTimeout                                  struct{ myError.MyError }
	WebsocketOffline                                    struct{ myError.MyError }
	AsyncMessageCallbackEmpty                           struct{ myError.MyError }
	AsyncMessageTimeout                                 struct{ myError.MyError }
	WebsocketClientExist                                struct{ myError.MyError }
	WebsocketClientNotExist                             struct{ myError.MyError }
	WebsocketServerConnConditionFuncEmpty               struct{ myError.MyError }
	WebsocketServerConnTagEmpty                         struct{ myError.MyError }
	WebsocketServerConnTagExist                         struct{ myError.MyError }
	WebsocketServerOnReceiveMessageSuccessCallbackEmpty struct{ myError.MyError }
)

func (*WebsocketConnOption) New(msg string) myError.IMyError {
	return &WebsocketConnOption{myError.MyError{Msg: fmt.Sprintf("websocket链接参数错误：%s", msg)}}
}

func (*WebsocketConnOption) Wrap(err error) myError.IMyError {
	return &WebsocketConnOption{myError.MyError{Msg: fmt.Errorf("websocket链接参数错误：%w", err).Error()}}
}

func (my *WebsocketConnOption) Error() string { return my.Msg }

func (*WebsocketConnOption) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketConnOptionErr)
}

func (*SyncMessageTimeout) New(msg string) myError.IMyError {
	return &SyncMessageTimeout{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("同步消息超时：%s", msg), "同步消息超时")}}
}

func (*SyncMessageTimeout) Wrap(err error) myError.IMyError {
	return &SyncMessageTimeout{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("同步消息超时：%w", err).Error(), "同步消息超时")}}
}

func (my *SyncMessageTimeout) Error() string { return my.Msg }

func (*SyncMessageTimeout) Is(target error) bool {
	return reflect.DeepEqual(target, &SyncMessageTimeoutErr)
}

func (*WebsocketOffline) New(msg string) myError.IMyError {
	return &WebsocketOffline{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("链接不在线：%s", msg), "链接不在线")}}
}

func (*WebsocketOffline) Wrap(err error) myError.IMyError {
	return &WebsocketOffline{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("链接不在线：%w", err).Error(), "链接不在线")}}
}

func (my *WebsocketOffline) Error() string { return my.Msg }

func (*WebsocketOffline) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketOfflineErr)
}

func (*AsyncMessageCallbackEmpty) New(msg string) myError.IMyError {
	return &AsyncMessageCallbackEmpty{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("异步消息回调不能为空：%s", msg), "异步消息回调不能为空")}}
}

func (*AsyncMessageCallbackEmpty) Wrap(err error) myError.IMyError {
	return &AsyncMessageCallbackEmpty{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("异步消息回调不能为空：%w", err).Error(), "异步消息回调不能为空")}}
}

func (my *AsyncMessageCallbackEmpty) Error() string { return my.Msg }

func (*AsyncMessageCallbackEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &AsyncMessageCallbackEmptyErr)
}

func (*AsyncMessageTimeout) New(msg string) myError.IMyError {
	return &AsyncMessageTimeout{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("异步消息回调超时必须大于0：%s", msg), "异步消息回调超时必须大于0")}}
}

func (*AsyncMessageTimeout) Wrap(err error) myError.IMyError {
	return &AsyncMessageTimeout{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("异步消息回调超时必须大于0：%w", err).Error(), "异步消息回调超时必须大于0")}}
}

func (my *AsyncMessageTimeout) Error() string { return my.Msg }

func (*AsyncMessageTimeout) Is(target error) bool {
	return reflect.DeepEqual(target, &AsyncMessageTimeoutErr)
}

func (*WebsocketClientExist) New(msg string) myError.IMyError {
	return &WebsocketClientExist{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket客户端已存在：%s", msg), "异步消息回调超时必须大于0")}}
}

func (*WebsocketClientExist) Wrap(err error) myError.IMyError {
	return &WebsocketClientExist{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket客户端已存在：%w", err).Error(), "websocket客户端已存在")}}
}

func (my *WebsocketClientExist) Error() string { return my.Msg }

func (*WebsocketClientExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketClientExistErr)
}

func (*WebsocketClientNotExist) New(msg string) myError.IMyError {
	return &WebsocketClientNotExist{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket客户端不存在：%s", msg), "websocket客户端不存在")}}
}

func (*WebsocketClientNotExist) Wrap(err error) myError.IMyError {
	return &WebsocketClientNotExist{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket客户端不存在：%w", err).Error(), "websocket客户端不存在")}}
}

func (my *WebsocketClientNotExist) Error() string { return my.Msg }

func (*WebsocketClientNotExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketClientExistErr)
}

func (*WebsocketServerConnConditionFuncEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerConnConditionFuncEmpty{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket服务端链接条件函数不能为空：%s", msg), "websocket服务端链接条件函数不能为空")}}
}

func (*WebsocketServerConnConditionFuncEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnConditionFuncEmpty{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket服务端链接条件函数不能为空：%w", err).Error(), "websocket服务端链接条件函数不能为空")}}
}

func (my *WebsocketServerConnConditionFuncEmpty) Error() string { return my.Msg }

func (*WebsocketServerConnConditionFuncEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnConditionFuncEmptyErr)
}

func (*WebsocketServerConnTagEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerConnTagEmpty{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket服务端连接标识不能为空：%s", msg), "websocket服务端连接标识不能为空")}}
}

func (*WebsocketServerConnTagEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnTagEmpty{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket服务端连接标识不能为空：%w", err).Error(), "websocket服务端连接标识不能为空")}}
}

func (my *WebsocketServerConnTagEmpty) Error() string { return my.Msg }

func (*WebsocketServerConnTagEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnTagEmptyErr)
}

func (*WebsocketServerConnTagExist) New(msg string) myError.IMyError {
	return &WebsocketServerConnTagExist{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket服务端连接标识已存在：%s", msg), "websocket服务端连接标识已存在")}}
}

func (*WebsocketServerConnTagExist) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnTagExist{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket服务端连接标识已存在：%w", err).Error(), "websocket服务端连接标识已存在")}}
}

func (my *WebsocketServerConnTagExist) Error() string { return my.Msg }

func (*WebsocketServerConnTagExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnTagExistErr)
}

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerOnReceiveMessageSuccessCallbackEmpty{myError.MyError{Msg: operation.Ternary[string](msg != "", fmt.Sprintf("websocket服务端接收消息成功回调不能为空：%s", msg), "websocket服务端接收消息成功回调不能为空")}}
}

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerOnReceiveMessageSuccessCallbackEmpty{myError.MyError{Msg: operation.Ternary[string](err != nil, fmt.Errorf("websocket服务端接收消息成功回调不能为空：%w", err).Error(), "websocket服务端接收消息成功回调不能为空")}}
}

func (my *WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Error() string { return my.Msg }

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerOnReceiveMessageSuccessCallbackEmptyErr)
}

var (
	WebsocketConnOptionErr                                 WebsocketConnOption
	SyncMessageTimeoutErr                                  SyncMessageTimeout
	WebsocketOfflineErr                                    WebsocketOffline
	AsyncMessageCallbackEmptyErr                           AsyncMessageCallbackEmpty
	AsyncMessageTimeoutErr                                 AsyncMessageTimeout
	WebsocketClientExistErr                                WebsocketClientExist
	WebsocketClientNotExistErr                             WebsocketClientNotExist
	WebsocketServerConnConditionFuncEmptyErr               WebsocketServerConnConditionFuncEmpty
	WebsocketServerConnTagEmptyErr                         WebsocketServerConnTagEmpty
	WebsocketServerConnTagExistErr                         WebsocketServerConnTagExist
	WebsocketServerOnReceiveMessageSuccessCallbackEmptyErr WebsocketServerOnReceiveMessageSuccessCallbackEmpty
)
