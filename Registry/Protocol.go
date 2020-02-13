package Registry

import (
	"gogistery/Protocol"
	"time"
)

type responserHeartProtocol struct {
	registry *Registry //服务于哪一个注册器
}

func (p *responserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	if timeout, retryN, ok := p.registry.register(request); ok {
		return Protocol.TobeSendResponse{
			Response: Protocol.Response{
				RegistryInfo: p.registry.Info,
				Timeout:      timeout,
				RetryN:       retryN,
				Reject:       false}, //同意连接
			Option: request.RegistrantInfo.GetResponseSendOption()}
	} else {
		return Protocol.TobeSendResponse{
			Response: Protocol.Response{
				RegistryInfo: p.registry.Info,
				Timeout:      timeout,
				RetryN:       retryN,
				Reject:       true}, //拒绝连接
			Option: request.RegistrantInfo.GetResponseSendOption()}
	}
}

type RegistrantControlProtocol interface {
	TimeoutRetryNForNew(request Protocol.Request) (time.Duration, uint64)
	TimeoutRetryNForUpdate(request Protocol.Request) (time.Duration, uint64)
}
