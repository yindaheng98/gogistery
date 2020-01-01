package proto

//此类用于在服务端成功连接到一个客户端时进行的一些操作成功时执行的操作。下面列举几个可以用此类的情况
//
//比如：某个服务器成功收到了到一个客户端的连接请求，用户想定义一个处理此信息的操作
//
//又比如：某个服务器发现一个客户端下线了，用户想定义一个处理此情况的操作
type EmitterClient struct {
	handlers []func(ClientInfo)
}

func (emitter *EmitterClient) AddHandler(handler func(ClientInfo)) {
	emitter.handlers = append(emitter.handlers, handler)
}

func (emitter *EmitterClient) Emit(info ClientInfo) {
	for _, handler := range emitter.handlers {
		handler(info)
	}
}

func NewEmitterClient() EmitterClient {
	return EmitterClient{[]func(ClientInfo){}}
}
