package RequesterHeart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/util/emitters"
)

type EventList struct {
	NewConnection    *emitters.RegistryInfoEmitter
	UpdateConnection *emitters.RegistryInfoEmitter
	Disconnection    *emitters.TobeSendRequestErrorEmitter
	Retry            *emitters.TobeSendRequestErrorEmitter
	Error            *Emitter.ErrorEmitter
}

func NewEventList() *EventList {
	return &EventList{
		emitters.NewRegistryInfoEmitter(),
		emitters.NewRegistryInfoEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
		Emitter.NewErrorEmitter()}
}

func (e *EventList) EnableAll() {
	e.NewConnection.Enable()
	e.UpdateConnection.Enable()
	e.Disconnection.Enable()
	e.Retry.Enable()
	e.Error.Enable()
}

func (e *EventList) DisableAll() {
	e.NewConnection.Disable()
	e.UpdateConnection.Disable()
	e.Disconnection.Disable()
	e.Retry.Enable()
	e.Error.Enable()
}
