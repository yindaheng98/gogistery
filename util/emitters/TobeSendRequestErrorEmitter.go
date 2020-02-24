package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

type TobeSendRequestErrorEmitter struct {
	*Emitter.ErrorInfoEmitter
}

func NewTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewErrorInfoEmitter()}
}

func (e *TobeSendRequestErrorEmitter) AddHandler(handler func(protocol.TobeSendRequest, error)) {
	e.ErrorInfoEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(protocol.TobeSendRequest), err)
	})
}

func (e *TobeSendRequestErrorEmitter) Emit(o protocol.TobeSendRequest, err error) {
	e.ErrorInfoEmitter.Emit(o, err)
}
