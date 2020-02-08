package Heart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type responserEvent struct {
	Error *Emitter.ErrorEmitter
}

type ResponserHeart struct {
	proto         ResponserHeartProtocol
	responser     *Responser
	Event         *responserEvent
	interruptChan chan bool
}

func NewResponserHeart(heartProto ResponserHeartProtocol, beatProto Protocol.ResponseBeatProtocol) *ResponserHeart {
	return &ResponserHeart{heartProto,
		NewResponser(beatProto),
		&responserEvent{Emitter.NewErrorEmitter()},
		make(chan bool, 1)}
}

//开始接收心跳，直到主动停止
func (h *ResponserHeart) RunBeating() {
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
	h.interruptChan <- true
}
