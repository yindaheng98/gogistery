package Heart

type RequesterHeart struct {
	proto RequesterHeartProtocol
}

func NewRequesterHeart(proto RequesterHeartProtocol) *RequesterHeart {
	return &RequesterHeart{proto}
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequestBeat TobeSendRequesterBeat) error {
	request := initRequestBeat
	for request.RequesterBeat != nil {
		if nextRequest, err := requesterBeat(request, h.proto); err != nil {
			return err
		} else {
			request = nextRequest
		}
	}
	return nil
}

//输入协议进行一次beat，返回下一次beat所需的数据
func requesterBeat(request TobeSendRequesterBeat, proto RequesterHeartProtocol) (TobeSendRequesterBeat, error) {
	response, err := proto.Request(request)
	if err != nil {
		return TobeSendRequesterBeat{nil, nil}, err
	}
	nextRequest := TobeSendRequesterBeat{nil, nil}
	proto.Beat(request, response, func(heartbeat TobeSendRequesterBeat) {
		nextRequest = heartbeat
	})
	return nextRequest, nil
}
