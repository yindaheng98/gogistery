package Registry

import (
	"fmt"
	"gogistery/Protocol"
	"time"
)

type Info interface {
	String() string
} //记录服务器信息

type Response struct { //服务器端的响应
	Info
	Timeout time.Duration //下一次连接的时间限制
	Reject  bool          //是否拒绝连接
}

func (r Response) String() string {
	return fmt.Sprintf("Registry.Response{%s,Timeout:%d,Reject:%t}", r.Info.String(), r.Timeout, r.Reject)
}

type ResponseSendOption struct{}

func (o ResponseSendOption) String() string {
	return "Registry.ResponseSendOption{}"
}

//记录服务器端收到的注册器信息
type RegistrantInfo interface {
	GetRegistrantID() string
}

type Request interface { //服务器端收到的请求
	RegistrantInfo
	String() string
}

type responserHeartProtocol struct {
	registry *Registry //服务于哪一个注册器
}

func (p *responserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	if timeout, ok := p.registry.register(request.(Request)); ok {
		return Protocol.TobeSendResponse{
			Response: Response{p.registry.info, timeout, false}, //同意连接
			Option:   ResponseSendOption{}}
	} else {
		return Protocol.TobeSendResponse{
			Response: Response{p.registry.info, timeout, true}, //拒绝连接
			Option:   ResponseSendOption{}}
	}
}

type RegistrantTimeoutProtocol interface {
	TimeoutForNew(request Request) time.Duration
	TimeoutForUpdate(request Request) time.Duration
}
