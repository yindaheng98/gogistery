package Heartbeat

import (
	"github.com/yindaheng98/go-utility/Emitter"
)

type ProtocolRequestSendOptionErrorEmitter struct {
	*Emitter.ErrorEmitter
}

func newProtocolRequestSendOptionErrorEmitter() *ProtocolRequestSendOptionErrorEmitter {
	return &ProtocolRequestSendOptionErrorEmitter{Emitter.NewErrorEmitter()}
}

func (e *ProtocolRequestSendOptionErrorEmitter) AddHandler(handler func(o ProtocolRequestSendOption, err error)) {
	e.ErrorEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(ProtocolRequestSendOption), err)
	})
}

func (e *ProtocolRequestSendOptionErrorEmitter) Emit(o ProtocolRequestSendOption, err error) {
	e.ErrorEmitter.Emit(o, err)
}
