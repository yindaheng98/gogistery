package base

import "time"

type SenderInfo interface {
	Send(addr string, timeout time.Duration) (ReceiverInfo, error)
	//这个里面的timeout项用于指定超时时间，超过此时间即判定为发送失败
}
