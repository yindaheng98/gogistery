package Heart

type RequesterHeart struct {
	proto RequesterHeartProtocol
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequest RequesterHeartbeat) error {
	request := initRequest
	for request != nil {
		if nextRequest, err := requesterBeat(request, h.proto); err != nil {
			return err
		} else {
			request = nextRequest
		}
	}
	return nil
}

//输入协议进行一次beat，返回下一次beat所需的数据
func requesterBeat(request RequesterHeartbeat, proto RequesterHeartProtocol) (RequesterHeartbeat, error) {
	response, err := proto.Request(request)
	if err != nil {
		return nil, err
	}
	var nextRequest RequesterHeartbeat = nil
	proto.Beat(request, response, func(heartbeat RequesterHeartbeat) {
		nextRequest = heartbeat
	})
	return nextRequest, nil
}
