package Sender

import (
	"gogistery/base"
	"time"
)

type Protocol interface {
	Send(senderInfo base.SenderInfo, addr string, timeout time.Duration) (base.ReceiverInfo, error)
	//这个里面的timeout项用于指定超时时间，超过此时间即判定为发送失败
}
