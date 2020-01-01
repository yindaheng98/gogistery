package proto

//此类用于让用户定义一些在客户端或服务端对另一个服务端进行的一些操作成功时执行的操作。
//
//比如：某个客户端或服务器成功连接到了另一个服务器，收到了另一个服务器传回的服务器信息，用户想定义一个处理此信息的操作
//
//又比如：某个客户端或服务器发现自己连接的服务器下线了，用户想定义一个处理此情况的操作
type ServerEmitter struct {
	handlers []func(ServerInfo)
}

func (emitter *ServerEmitter) AddHandler(handler func(ServerInfo)) {
	emitter.handlers = append(emitter.handlers, handler)
}

func (emitter *ServerEmitter) Emit(info ServerInfo) {
	for _, handler := range emitter.handlers {
		handler(info)
	}
}

func NewServerEmitter() ServerEmitter {
	return ServerEmitter{[]func(ServerInfo){}}
}
