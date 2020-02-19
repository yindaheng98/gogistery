package Registrant

import (
	"gogistery/Heart"
	"gogistery/Protocol"
	"time"
)

type heart struct {
	registrant   *Registrant           //此heart服务于哪个注册器
	RegistryInfo Protocol.RegistryInfo //此heart当前连接着哪个注册中心

	heartProto *requesterHeartProtocol
	requester  *Heart.RequesterHeart

	stoppedChan chan bool
}

func newHeart(registrant *Registrant, sendProto Protocol.RequestProtocol) *heart {
	stoppedChan := make(chan bool, 1)
	close(stoppedChan)

	h := &heart{
		registrant:   registrant,
		RegistryInfo: nil,

		heartProto: nil,
		requester:  nil,

		stoppedChan: stoppedChan}

	heartProto := newRequesterHeartProtocol(h)
	h.heartProto = heartProto
	requester := Heart.NewRequesterHeart(heartProto, sendProto)
	requester.Events.Retry = registrant.Events.Retry
	h.requester = requester
	return h
}

func (h *heart) beatResponse(response Protocol.Response) Protocol.TobeSendRequest {
	h.registrant.candProto.StoreCandidates(response)
	if response.IsReject() {
		h.RegistryInfo = nil
	} else {
		h.RegistryInfo = response.RegistryInfo
	}
	return Protocol.TobeSendRequest{
		Request: Protocol.Request{RegistrantInfo: h.registrant.Info},
		Option:  response.RegistryInfo.GetRequestSendOption(),
	}
}

func (h *heart) Run(initRequest Protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	h.stoppedChan = make(chan bool, 1)
	h.heartProto.start()
	err := h.requester.RunBeating(initRequest, initTimeout, initRetryN)
	h.RegistryInfo = nil
	h.stoppedChan <- true
	close(h.stoppedChan)
	return err
}

func (h *heart) Stop() {
	h.heartProto.stop()
	<-h.stoppedChan
}

type CandidateRegistryProtocol interface {
	//存入一组候选注册中心
	StoreCandidates(response Protocol.Response)

	//选出一个用于初始化的注册中心信息，并且不能是except中列出的这几个
	GetCandidate(except []Protocol.RegistryInfo) (Protocol.RegistryInfo, time.Duration, uint64)
}
