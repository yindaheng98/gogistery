package RequesterHeart

import (
	"gogistery/Protocol"
	"time"
)

type RequesterHeart struct {
	proto     RequesterHeartProtocol
	requester *requester
	Event     *EventList
}

func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestProtocol) *RequesterHeart {
	heart := &RequesterHeart{heartProto, nil, NewEventList()}
	heart.requester = newRequester(beatProto, heart)
	return heart
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequest Protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	request, timeout, retryN := initRequest, initTimeout, initRetryN
	var err error = nil
	defer func() { h.Event.Disconnection.Emit(request, err) }()
	established := false
	run := true
	for run {
		response, err := h.requester.Send(request, timeout, retryN)
		if err != nil {
			h.Event.Error.Emit(err)
			return err
		}
		if !established {
			h.Event.NewConnection.Emit(response.RegistryInfo)
			established = true
		} else {
			h.Event.UpdateConnection.Emit(response.RegistryInfo)
		}
		run = false
		h.proto.Beat(response, func(requestB Protocol.TobeSendRequest, timeoutB time.Duration, retryNB uint64) {
			request, timeout, retryN = requestB, timeoutB, retryNB
			run = true
		})
	}
	return nil
}
