package Heart

import (
	"gogistery/Protocol"
)

type RequesterHeart struct {
	proto RequesterHeartProtocol
}

func NewRequesterHeart(proto RequesterHeartProtocol) *RequesterHeart {
	return &RequesterHeart{proto}
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequestBeat Protocol.TobeSendRequest) error {
	request := initRequestBeat
	for request.Request != nil {
		if nextRequest, err := requesterBeat(request, h.proto); err != nil {
			return err
		} else {
			request = nextRequest
		}
	}
	return nil
}

//输入协议进行一次beat，返回下一次beat所需的数据
func requesterBeat(request Protocol.TobeSendRequest, proto RequesterHeartProtocol) (Protocol.TobeSendRequest, error) {
	response, err := proto.Request(request)
	if err != nil {
		return Protocol.TobeSendRequest{}, err
	}
	nextRequest := Protocol.TobeSendRequest{}
	proto.Beat(request, response, func(heartbeat Protocol.TobeSendRequest) {
		nextRequest = heartbeat
	})
	return nextRequest, nil
}
