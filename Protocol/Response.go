package Protocol

import (
	"fmt"
	"time"
)

//记录服务器信息
type RegistryInfo interface {
	GetRegistryID() string
	String() string
}

//心跳数据响应基础类
type Response struct {
	RegistryInfo
	Timeout time.Duration //下一次连接的时间限制
	Reject  bool          //是否拒绝连接
}

func (r Response) String() string {
	return fmt.Sprintf("Registry.Response{RegistryInfo:%s,Timeout:%d,Reject:%t}",
		r.RegistryInfo.String(), r.Timeout, r.Reject)
}

//自定义响应发送设置
type ResponseSendOption interface {
	String() string
}
