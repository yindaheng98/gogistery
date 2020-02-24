package registrant

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/util/emitters"
)

type events struct {
	NewConnection    *emitters.RegistryInfoEmitter
	UpdateConnection *emitters.RegistryInfoEmitter
	Disconnection    *emitters.TobeSendRequestErrorEmitter
	Error            *Emitter.ErrorEmitter
	Retry            *emitters.TobeSendRequestErrorEmitter
}

func newEvents() *events {
	return &events{
		emitters.NewRegistryInfoEmitter(),
		emitters.NewRegistryInfoEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
		Emitter.NewErrorEmitter(),
		emitters.NewTobeSendRequestErrorEmitter()}
}

func (e *events) EnableAll() {
	e.NewConnection.Enable()
	e.UpdateConnection.Enable()
	e.Disconnection.Enable()
	e.Error.Enable()
	e.Retry.Enable()
}

func (e *events) DisableAll() {
	e.NewConnection.Disable()
	e.UpdateConnection.Disable()
	e.Disconnection.Disable()
	e.Error.Enable()
	e.Retry.Disable()
}
