package protocol

import "fmt"

//自定义响应发送设置
type ResponseSendOption interface {
	String() string
}

//记录服务器端收到的注册器信息
type RegistrantInfo interface {
	GetRegistrantID() string
	GetResponseSendOption() ResponseSendOption //此服务端接收何种请求
	String() string
}

//心跳数据请求基础类
type Request struct { //服务器端收到的请求
	RegistrantInfo RegistrantInfo
	Disconnect     bool
}

func (r Request) IsDisconnect() bool {
	return r.Disconnect
}

func (r Request) String() string {
	return fmt.Sprintf("Request{RegistrantInfo:%s,Disconnect:%t}", r.RegistrantInfo.String(), r.Disconnect)
}
