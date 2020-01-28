package Receiver

import "gogistery/base/Emitter"

type events struct {
	Start        *Emitter.EmptyEmitter
	Stop         *Emitter.EmptyEmitter
	Connected    *Emitter.SenderInfoEmitter
	Disconnected *Emitter.SenderInfoEmitter
}

func newEvents() *events {
	return &events{Emitter.NewEmptyEmitter(),
		Emitter.NewEmptyEmitter(),
		Emitter.NewSenderInfoEmitter(),
		Emitter.NewSenderInfoEmitter()}
}

func (e *events) EnableAll() {
	e.Start.Enable()
	e.Stop.Enable()
	e.Connected.Enable()
	e.Disconnected.Enable()
}

func (e *events) DisableAll() {
	e.Start.Enable()
	e.Stop.Enable()
	e.Connected.Enable()
	e.Disconnected.Enable()
}
