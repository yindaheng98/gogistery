package Heart

import (
	"gogistery/Protocol"
	"time"
)

type RequesterHeart struct {
	proto     RequesterHeartProtocol
	requester *requester
	Events    *requesterEvents
}

func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestProtocol) *RequesterHeart {
	requester := newRequester(beatProto)
	return &RequesterHeart{heartProto, requester, requester.Events}
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
