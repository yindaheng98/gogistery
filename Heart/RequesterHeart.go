package Heart

type RequesterHeart struct {
	proto RequesterHeartProtocol
}

func NewRequesterHeart(proto RequesterHeartProtocol) *RequesterHeart {
	return &RequesterHeart{proto}
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequestBeat TobeSendRequest) error {
	request := initRequestBeat
	run := true
	for run {
		response, err := h.proto.Request(request)
		if err != nil {
			return err
		}
		run = false
		h.proto.Beat(response, func(heartbeat TobeSendRequest) {
			request = heartbeat
			run = true
		})
	}
	return nil
}
