package base

import "time"

//此接口仅负责记录在发送和接收器之间传递的消息
type ReceiverInfo interface {
	GetAddr() string           //获取下一次的请求地址
	GetTimeout() time.Duration //获取下一次请求的间隔时间
	GetRetryN() uint32         //获取最大重试次数
	IsReject() bool            //是否要求断开连接
}
