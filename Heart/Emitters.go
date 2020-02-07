package Heart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type ErrorEmitter struct {
	*Emitter.Emitter
}

func newErrorEmitter() *ErrorEmitter {
	return &ErrorEmitter{Emitter.NewEmitter()}
}

func (e *ErrorEmitter) AddHandler(handler func(err error)) {
	e.Emitter.AddHandler(func(err interface{}) {
		handler(err.(error))
	})
}

func (e *ErrorEmitter) Emit(err error) {
	e.Emitter.Emit(err)
}

type TobeSendRequestErrorEmitter struct {
	*Emitter.ErrorEmitter
}

func newTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewErrorEmitter()}
}

func (e *TobeSendRequestErrorEmitter) AddHandler(handler func(o Protocol.TobeSendRequest, err error)) {
	e.ErrorEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(Protocol.TobeSendRequest), err)
	})
}

func (e *TobeSendRequestErrorEmitter) Emit(o Protocol.TobeSendRequest, err error) {
	e.ErrorEmitter.Emit(o, err)
}
