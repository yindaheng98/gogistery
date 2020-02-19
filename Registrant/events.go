package Registrant

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/util/emitters"
)

type events struct {
	NewConnection     *emitters.RegistryInfoEmitter
	UpdateConnection  *emitters.RegistryInfoEmitter
	ConnectionTimeout *emitters.RegistryInfoEmitter
	Disconnection     *emitters.RegistryInfoEmitter
	Retry             *emitters.TobeSendRequestErrorEmitter
	Error             *Emitter.ErrorEmitter
}

func newEvents() *events {
	return &events{
		emitters.NewRegistryInfoEmitter(),
		emitters.NewRegistryInfoEmitter(),
		emitters.NewRegistryInfoEmitter(),
		emitters.NewRegistryInfoEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
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
