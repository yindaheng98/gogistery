package RequesterHeart

import (
	"gogistery/Protocol"
	"gogistery/util/emitters"
	"time"
)

type RequesterHeart struct {
	proto     RequesterHeartProtocol
	requester *requester
	Events    *events
}

func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestProtocol) *RequesterHeart {
	events := &events{
		NewConnection:    emitters.NewRegistryInfoEmitter(),
		UpdateConnection: emitters.NewRegistryInfoEmitter(),
		Retry:            emitters.NewTobeSendRequestErrorEmitter(),
		Disconnection:    emitters.NewRegistryInfoEmitter(),
	}
	heart := &RequesterHeart{heartProto, nil, events}
	heart.requester = newRequester(beatProto, heart)
	return heart
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequest Protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	request, timeout, retryN := initRequest, initTimeout, initRetryN
	run := true
	for run {
		response, err := h.requester.Send(request, timeout, retryN)
		if err != nil {
			return err
		}
		run = false
		h.proto.Beat(response, func(requestB Protocol.TobeSendRequest, timeoutB time.Duration, retryNB uint64) {
			request, timeout, retryN = requestB, timeoutB, retryNB
			run = true
		})
	}
	return nil
}
