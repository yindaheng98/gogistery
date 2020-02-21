package Registrant

import (
	"gogistery/heart/requester"
	"gogistery/protocol"
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
func (h *heart) register(response protocol.Response) protocol.TobeSendRequest {
	request := h.registrant.register(response)
	if response.IsReject() {
		h.RegistryInfo = nil
	} else {
		h.RegistryInfo = response.RegistryInfo
	}
	return protocol.TobeSendRequest{
		Request: request,
		Option:  response.RegistryInfo.GetRequestSendOption(),
	}
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
