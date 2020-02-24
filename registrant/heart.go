package registrant

import (
	"github.com/yindaheng98/gogistry/heart/requester"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type heart struct {
	*requester.Heart
	beater       *beater
	registrant   *Registrant
	RegistryInfo protocol.RegistryInfo //此heart当前连接着哪个注册中心

	stoppedChan chan bool
}

func newHeart(registrant *Registrant, retryNController RetryNController, RequestProto protocol.RequestProtocol) *heart {
	stoppedChan := make(chan bool, 1)
	close(stoppedChan)
	heart := &heart{
		Heart:        nil,
		beater:       nil,
		registrant:   registrant,
		RegistryInfo: nil,
		stoppedChan:  stoppedChan,
	}
	heart.beater = newBeater(heart, retryNController)
	heart.Heart = requester.NewHeart(heart.beater, RequestProto)
	return heart
}

//For the struct beater
func (h *heart) register(response protocol.Response) (protocol.TobeSendRequest, bool) {
	request, ok := h.registrant.register(response)
	if (!ok) || response.IsReject() { //如果不可响应或是拒绝连接
		h.RegistryInfo = nil //就清空已连接记录
	} else {
		h.RegistryInfo = response.RegistryInfo //否则写入已连接记录
	}
	return protocol.TobeSendRequest{
		Request: request,
		Option:  response.RegistryInfo.GetRequestSendOption(),
	}, ok
}

func (h *heart) RunBeating(initRequest protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	h.stoppedChan = make(chan bool, 1)
	h.beater.Start()
	err := h.Heart.RunBeating(initRequest, initTimeout, initRetryN)
	h.RegistryInfo = nil
	h.stoppedChan <- true
	close(h.stoppedChan)
	return err
}

func (h *heart) Stop() {
	h.beater.Stop()
	<-h.stoppedChan
}
