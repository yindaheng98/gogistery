package Emitter

import (
	"gogistery/base"
	"gogistery/util/Emitter"
)

//事件格式为base.ReceiverInfo
type ReceiverInfoEmitter struct {
	*Emitter.Emitter
}

func NewReceiverInfoEmitter() *ReceiverInfoEmitter {
	return &ReceiverInfoEmitter{Emitter.New()}
}

func (e *ReceiverInfoEmitter) AddHandler(handler func(info base.ReceiverInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(base.ReceiverInfo))
	})
}

func (e *ReceiverInfoEmitter) Emit(info base.ReceiverInfo) {
	e.Emitter.Emit(info)
}
