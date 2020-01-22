package base

type SenderInfo interface {
	GetID() string      //获取ID
	IsDisconnect() bool //是否请求中断连接
}
