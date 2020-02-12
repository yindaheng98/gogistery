package Protocol

//记录服务器端收到的注册器信息
type RegistrantInfo interface {
	GetRegistrantID() string
}

//心跳数据请求基础类
type Request interface { //服务器端收到的请求
	RegistrantInfo
	GetResponseSendOption() ResponseSendOption
	IsDisconnect() bool
	String() string
}

//自定义请求发送设置
type RequestSendOption interface {
	String() string
}
