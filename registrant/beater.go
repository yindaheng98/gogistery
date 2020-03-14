package registrant

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type beater struct {
	registrant       *Registrant //此协议服务于哪个heart
	retryNController WaitTimeoutRetryNController
	i                uint64 //此协议编号
}

func newBeater(registrant *Registrant, retryNController WaitTimeoutRetryNController, i uint64) *beater {
	return &beater{registrant, retryNController, i}
}
func (p *beater) Beat(ctx context.Context, response protocol.Response, lastTimeout time.Duration, lastRetryN uint64, beat func(protocol.TobeSendRequest, time.Duration, uint64)) {
	if response.RegistryInfo.GetServiceType() != p.registrant.Info.GetServiceType() { //类型检查不通过
		return //则断开连接
	}
	if ok := p.registrant.register(ctx, response, p.i); !ok { //如果heart拒绝了连接请求
		return //就直接断连退出
	}
	waitTime, sendTimeout, retryN := p.retryNController.GetWaitTimeoutRetryN(response, lastTimeout, lastRetryN)
	time.Sleep(waitTime) //等待一段时间再发
	beat(protocol.TobeSendRequest{
		Request: protocol.Request{
			RegistrantInfo: p.registrant.Info,
			Disconnect:     false,
		},
		Option: response.RegistryInfo.GetRequestSendOption(),
	}, sendTimeout, retryN)
}
