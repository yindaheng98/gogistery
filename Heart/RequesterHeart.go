package Heart

import "gogistery/Protocol"

type RequesterHeart struct {
	proto     RequesterHeartProtocol
	requester *Requester
}

func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestBeatProtocol) *RequesterHeart {
	return &RequesterHeart{heartProto, NewRequester(beatProto)}
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequestBeat Protocol.TobeSendRequest, initRequestOption RequestSendOption) error {
	request := initRequestBeat
	option := initRequestOption
	run := true
	for run {
		response, err := h.requester.Send(request, option.Timeout, option.RetryN)
		if err != nil {
			return err
		}
		run = false
		h.proto.Beat(response, func(heartbeat Protocol.TobeSendRequest, sendOption RequestSendOption) {
			request = heartbeat
			option = sendOption
			run = true
		})
	}
	return nil
}
