package Heart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type TobeSendRequestErrorEmitter struct {
	*Emitter.ErrorInfoEmitter
}

func newTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewErrorInfoEmitter()}
}

func (e *TobeSendRequestErrorEmitter) AddHandler(handler func(o Protocol.TobeSendRequest, err error)) {
	e.ErrorInfoEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(Protocol.TobeSendRequest), err)
	})
}

func (e *TobeSendRequestErrorEmitter) Emit(o Protocol.TobeSendRequest, err error) {
	e.ErrorInfoEmitter.Emit(o, err)
}
