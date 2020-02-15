package Registrant

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Heart"
	"gogistery/Protocol"
)

//事件格式为base.RegistrantInfo
type RegistryInfoEmitter struct {
	*Emitter.Emitter
}

func (e *RegistryInfoEmitter) AddHandler(handler func(info Protocol.RegistryInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Protocol.RegistryInfo))
	})
}

func (e *RegistryInfoEmitter) Emit(info Protocol.RegistryInfo) {
	e.Emitter.Emit(info)
}

type events struct {
	NewConnection     *RegistryInfoEmitter
	UpdateConnection  *RegistryInfoEmitter
	ConnectionTimeout *RegistryInfoEmitter
	Disconnection     *RegistryInfoEmitter
	Retry             *Heart.TobeSendRequestErrorEmitter
	Error             *Emitter.ErrorEmitter
}

func newEvents() *events {
	return &events{
		&RegistryInfoEmitter{Emitter.NewEmitter()},
		&RegistryInfoEmitter{Emitter.NewEmitter()},
		&RegistryInfoEmitter{Emitter.NewEmitter()},
		&RegistryInfoEmitter{Emitter.NewEmitter()},
		&Heart.TobeSendRequestErrorEmitter{ErrorInfoEmitter: Emitter.NewErrorInfoEmitter()},
		Emitter.NewErrorEmitter()}
}

func (e *events) EnableAll() {
	e.NewConnection.Enable()
	e.UpdateConnection.Enable()
	e.ConnectionTimeout.Enable()
	e.Disconnection.Enable()
	e.Retry.Enable()
	e.Error.Enable()
}

func (e *events) DisableAll() {
	e.NewConnection.Disable()
	e.UpdateConnection.Disable()
	e.ConnectionTimeout.Disable()
	e.Disconnection.Disable()
	e.Retry.Enable()
	e.Error.Enable()
}
