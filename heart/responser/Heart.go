package responser

import (
	"gogistery/Protocol"
)

type Heart struct {
	beater          HeartBeater
	responser       *responser
	ErrorHandler    func(error)
	interruptChan   chan bool
	interruptedChan chan bool
}

func NewHeart(beater HeartBeater, ResponseProto Protocol.ResponseProtocol) *Heart {
	interruptChan := make(chan bool, 1)
	interruptedChan := make(chan bool, 1)
	close(interruptChan)
	close(interruptedChan)
	return &Heart{beater,
		newResponser(ResponseProto),
		func(error) {},
		interruptChan, interruptedChan}
}

//开始接收心跳，直到主动停止
func (h *Heart) RunBeating() {
	h.interruptChan = make(chan bool, 1)
	h.interruptedChan = make(chan bool, 1)
	defer func() {
		h.interruptedChan <- true
		close(h.interruptChan)
		close(h.interruptedChan)
	}()
	for {
		var request Protocol.Request
		var err error
		var responseFunc func(Protocol.TobeSendResponse)
		responseChan := make(chan bool, 1)
		go func() {
			request, err, responseFunc = h.responser.Recv()
			if err != nil {
				h.ErrorHandler(err)
			} else {
				responseFunc(h.beater.Beat(request))
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

func (h *Heart) Stop() {
	defer func() { recover() }()
	h.interruptChan <- true
	<-h.interruptedChan
}
