package registry

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
)

type beater struct {
	registry *Registry //服务于哪一个注册器
}

func (p *beater) Beat(ctx context.Context, request protocol.Request) protocol.TobeSendResponse {
	if p.registry.Info.GetServiceType() != request.RegistrantInfo.GetServiceType() { //类型检查不通过则拒绝连接
		return protocol.TobeSendResponse{
			Response: protocol.Response{
				RegistryInfo: p.registry.Info,
				Timeout:      3600 * 24 * 1e9,
				Reject:       true}, //拒绝连接
			Option: request.RegistrantInfo.GetResponseSendOption()}
	} else if timeout, ok := p.registry.register(request); !ok { //注册中心表示不能连接
		return protocol.TobeSendResponse{
			Response: protocol.Response{
				RegistryInfo: p.registry.Info,
				Timeout:      timeout,
				Reject:       true}, //拒绝连接
			Option: request.RegistrantInfo.GetResponseSendOption()}
	} else {
		return protocol.TobeSendResponse{
			Response: protocol.Response{
				RegistryInfo: p.registry.Info,
				Timeout:      timeout,
				Reject:       false}, //同意连接
			Option: request.RegistrantInfo.GetResponseSendOption()}
	}
}
