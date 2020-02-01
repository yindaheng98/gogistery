package Sender

import (
	"gogistery/base"
)

type ProtoChanElement struct {
	info  base.ReceiverInfo
	error error
}

type Protocol interface {
	//将senderInfo发往addr所指的地址，并将返回信息放入protoChan中
	//
	//如果超时，protoChan将关闭
	Send(senderInfo base.SenderInfo, addr string, protoChan chan ProtoChanElement)

	//发送“停止连接”的消息
	SendDisconnect(senderInfo base.SenderInfo, addr string)
}
