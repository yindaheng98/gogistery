package Heart

type responserEvent struct {
	Error *ErrorEmitter
}

type ResponserHeart struct {
	proto         ResponserHeartbeatProtocol
	Event         *responserEvent
	interruptChan chan bool
}

//开始接收心跳，直到主动停止
func (h *ResponserHeart) RunBeating() {
	for {
		var request RequesterHeartbeat
		var err error
		var responseFunc func(ResponserHeartbeat)
		responseChan := make(chan bool, 1)
		go func() {
			request, err, responseFunc = h.proto.Response()
			if err != nil {
				h.Event.Error.Emit(err)
			} else {
				responseFunc(h.proto.Beat(request))
			}
			responseChan <- true
		}()

		select {
		case <-responseChan:
		case <-h.interruptChan:
			return
		}
	}
}

func (h *ResponserHeart) Stop() {
	h.interruptChan <- true
}
