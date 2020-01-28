package Receiver

import "gogistery/base"

type Protocol interface {
	//info用于指定返回的数据，返回接收到的数据
	//
	//此接口的实现中等待连接到达并接收到完整SenderInfo或出错时才会返回
	Receive(info base.ReceiverInfo) (base.SenderInfo, error)
}
