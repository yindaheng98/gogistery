package Heart

import "gogistery/Protocol"

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
func (h *RequesterHeart) RunBeating(initRequest Protocol.TobeSendRequest) error {
	request, option := initRequest, initRequest.Option.(RequestSendOption)
	run := true
	for run {
		response, err := h.requester.Send(request, option.GetTimeout(), option.GetRetryN())
		if err != nil {
			return err
		}
		run = false
		h.proto.Beat(response, func(req Protocol.TobeSendRequest) {
			request, option = req, req.Option.(RequestSendOption)
			run = true
		})
	}
	return nil
}
