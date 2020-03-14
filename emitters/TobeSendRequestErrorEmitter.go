package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

//TobeSendRequestErrorEmitter use protocol.TobeSendRequest as event payload
type TobeSendRequestErrorEmitter struct {
	*Emitter.ErrorInfoEmitter
}

//NewSyncTobeSendRequestErrorEmitter returns the pointer to a sync TobeSendRequestErrorEmitter
func NewSyncTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewSyncErrorInfoEmitter()}
}

//NewAsyncTobeSendRequestErrorEmitter returns the pointer to a async TobeSendRequestErrorEmitter
func NewAsyncTobeSendRequestErrorEmitter() *TobeSendRequestErrorEmitter {
	return &TobeSendRequestErrorEmitter{Emitter.NewAsyncErrorInfoEmitter()}
}

//Implementation of Emitter.AddHandler
func (e *TobeSendRequestErrorEmitter) AddHandler(handler func(protocol.TobeSendRequest, error)) {
	e.ErrorInfoEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(protocol.TobeSendRequest), err)
	})
}

//Implementation of Emitter.Emit
func (e *TobeSendRequestErrorEmitter) Emit(o protocol.TobeSendRequest, err error) {
	e.ErrorInfoEmitter.Emit(o, err)
}
