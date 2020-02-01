package proto

import (
	"gogistery/util/Emitter"
	"log"
)

//此类用于在服务端成功连接到一个客户端时进行的一些操作成功时执行的操作。下面列举几个可以用此类的情况
//
//比如：某个服务器成功收到了到一个客户端的连接请求，用户想定义一个处理此信息的操作
//
//又比如：某个服务器发现一个客户端下线了，用户想定义一个处理此情况的操作
type ClientEmitter struct {
	emitter *Emitter.Emitter
}

func (e *ClientEmitter) AddHandler(handler func(*ClientInfo)) {
	e.emitter.AddHandler(func(bytes []byte) {
		info, err := ParseClient(bytes)
		if err != nil {
			handler(info)
		} else {
			log.Println("Failed parsing a ClientInfo: " + string(bytes))
		}
	})
}

func (e *ClientEmitter) Emit(info ClientInfo) {
	e.emitter.Emit(info.String())
}

func NewClientEmitter() ClientEmitter {
	return ClientEmitter{[]func(ClientInfo){}}
}
