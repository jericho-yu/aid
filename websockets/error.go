package websockets

import "errors"

var (
	WebsocketConnOptionErr                                 = errors.New("websocket链接错误：参数错误")
	SyncMessageTimeoutErr                                  = errors.New("同步消息超时")
	WebsocketOfflineErr                                    = errors.New("链接不在线")
	AsyncMessageCallbackEmptyErr                           = errors.New("异步消息回调不能为空")
	AsyncMessageTimeoutErr                                 = errors.New("异步消息回调超时必须大于0")
	WebsocketClientExistErr                                = errors.New("websocket客户端已存在")
	WebsocketClientNotExistErr                             = errors.New("websocket客户端不存在")
	WebsocketServerConnConditionFuncEmptyErr               = errors.New("websocket服务端链接条件函数不能为空")
	WebsocketServerConnTagEmptyErr                         = errors.New("websocket服务端连接标识不能为空")
	WebsocketServerConnTagExistErr                         = errors.New("websocket服务端连接标识已存在")
	WebsocketServerOnReceiveMessageSuccessCallbackEmptyErr = errors.New("websocket服务端接收消息成功回调不能为空")
	// WebsocketServerOnReceiveMessageSuccessCallbackEmptyEr  = errors.New("websocket服务端接收消息成功回调不能为空")
)
