package Receiver

import (
	"gogistery/base"
	"gogistery/base/Emitter"
)

type element struct {
	sender base.SenderInfo

	connected    *Emitter.SenderInfoEmitter
	disconnected *Emitter.SenderInfoEmitter
}

func (e *element) NewAddedHandler() {
	e.connected.Emit(e.sender)
}

func (e *element) TimeoutHandler() {
	e.disconnected.Emit(e.sender)
}

func (e *element) GetID() string {
	return e.sender.GetID()
}
