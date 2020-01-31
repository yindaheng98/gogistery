package base

//此接口仅负责记录在发送和接收器之间传递的消息
type SenderInfo interface {
	GetID() string      //获取ID
	IsDisconnect() bool //是否请求中断连接
}
