package Receiver

import "gogistery/base"

type Protocol interface {
	//用一个chan向处理线程传递接收到的数据，另一个chan用于指定要返回的数据
	Receive(senderInfoChan chan base.SenderInfo, recvInfoChan chan base.ReceiverInfo)
	Reject()
}
