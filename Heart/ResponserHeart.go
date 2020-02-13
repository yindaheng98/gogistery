package Heart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type responserEvents struct {
	Error *Emitter.ErrorEmitter
}

type ResponserHeart struct {
	proto           ResponserHeartProtocol
	responser       *responser
	Event           *responserEvents
	interruptChan   chan bool
	interruptedChan chan bool
}

func NewResponserHeart(heartProto ResponserHeartProtocol, beatProto Protocol.ResponseProtocol) *ResponserHeart {
	interruptChan := make(chan bool, 1)
	interruptedChan := make(chan bool, 1)
	close(interruptChan)
	close(interruptedChan)
	return &ResponserHeart{heartProto,
		newResponser(beatProto),
		&responserEvents{Emitter.NewErrorEmitter()},
		interruptChan, interruptedChan}
}

//开始接收心跳，直到主动停止
func (h *ResponserHeart) RunBeating() {
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
	defer func() { recover() }()
	h.interruptChan <- true
	<-h.interruptedChan
}
