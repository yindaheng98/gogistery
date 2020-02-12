package Registry

import (
	"gogistery/Protocol"
)

type responserHeartProtocol struct {
	registry *Registry //服务于哪一个注册器
}

func (p *responserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	if timeout, ok := p.registry.register(request); ok {
		return Protocol.TobeSendResponse{
			Response: Protocol.Response{RegistryInfo: p.registry.info, Timeout: timeout, Reject: false}, //同意连接
			Option:   request.RegistrantInfo.GetResponseSendOption()}
	} else {
		return Protocol.TobeSendResponse{
			Response: Protocol.Response{RegistryInfo: p.registry.info, Timeout: timeout, Reject: true}, //拒绝连接
			Option:   request.RegistrantInfo.GetResponseSendOption()}
	}
}
