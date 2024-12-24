package websockets

type ClientCallbackConfig struct {
	OnConnSuccessCallback           standardSuccessFn
	OnConnFailCallback              standardFailFn
	OnCloseSuccessCallback          standardSuccessFn
	OnCloseFailCallback             standardFailFn
	OnReceiveMessageSuccessCallback receiveMessageSuccessFn
	OnReceiveMessageFailCallback    standardFailFn
	OnSendMessageFailCallback       standardFailFn
}
