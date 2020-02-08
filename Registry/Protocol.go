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
	Reject bool //是否拒绝连接
}

func (r Response) String() string {
	return fmt.Sprintf("Registry.Response{%s,Reject:%t}", r.Info.String(), r.Reject)
}

type ResponseSendOption struct{}

func (o ResponseSendOption) String() string {
	return "Registry.ResponseSendOption{}"
}

//记录服务器端收到的注册器信息
type RegistrantInfo interface {
	GetID() string
}

type Request interface { //服务器端收到的请求
	RegistrantInfo
	String() string
}

type responserHeartProtocol struct {
	registry     *Registry       //服务于哪一个注册器
	timeoutProto TimeoutProtocol //如何选择timeout
}

func (p *responserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	registrantInfo := request.(RegistrantInfo)
	if p.registry.register(registrantInfo, p.timeoutProto.DecideNextTimeout(request.(Request))) {
		return Protocol.TobeSendResponse{
			Response: Response{p.registry.info, true}, //同意连接
			Option:   ResponseSendOption{}}
	} else {
		return Protocol.TobeSendResponse{
			Response: Response{p.registry.info, false}, //拒绝连接
			Option:   ResponseSendOption{}}
	}
}

type TimeoutProtocol interface {
	DecideNextTimeout(info Request) time.Duration
}
