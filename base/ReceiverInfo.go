package base

import "time"

type ReceiverInfo interface {
	GetAddr() string           //获取下一次的请求地址
	GetTimeout() time.Duration //获取下一次请求的间隔时间
	GetRetryN() uint32         //获取最大重试次数
	IsDisconnect() bool        //是否要求断开连接
}
