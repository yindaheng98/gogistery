package Heartbeat

import (
	"github.com/yindaheng98/go-utility/Emitter"
)

type TobeSendRequestErrorEmitter struct {
	*Emitter.ErrorEmitter
}

func newTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewErrorEmitter()}
}

func (e *TobeSendRequestErrorEmitter) AddHandler(handler func(o TobeSendRequest, err error)) {
	e.ErrorEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(TobeSendRequest), err)
	})
}

func (e *TobeSendRequestErrorEmitter) Emit(o TobeSendRequest, err error) {
	e.ErrorEmitter.Emit(o, err)
}
