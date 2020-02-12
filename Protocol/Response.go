package Protocol

import (
	"fmt"
	"time"
)

//自定义请求发送设置
type RequestSendOption interface {
	String() string
}

//记录服务器信息
type RegistryInfo interface {
	GetRegistryID() string
	GetRequestSendOption() RequestSendOption //此服务端接收何种请求
	String() string
}

//心跳数据响应基础类
type Response struct {
	RegistryInfo RegistryInfo
	Timeout      time.Duration //下一次连接的时间限制
	Reject       bool          //是否拒绝连接
}

func (r Response) String() string {
	return fmt.Sprintf("Registry.Response{RegistryInfo:%s,Timeout:%d,Reject:%t}",
		r.RegistryInfo.String(), r.Timeout, r.Reject)
}
