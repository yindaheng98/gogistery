package Registry

import (
	"github.com/yindaheng98/go-utility/Emitter"
)

//事件格式为base.RegistrantInfo
type RegistrantInfoEmitter struct {
	*Emitter.Emitter
}

func (e *RegistrantInfoEmitter) AddHandler(handler func(info RegistrantInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(RegistrantInfo))
	})
}

func (e *RegistrantInfoEmitter) Emit(info RegistrantInfo) {
	e.Emitter.Emit(info)
}

type events struct {
	NewConnection    *RegistrantInfoEmitter
	UpdateConnection *RegistrantInfoEmitter
	Disconnection    *RegistrantInfoEmitter
	Error            *Emitter.ErrorEmitter
}

func newEvents() *events {
	return &events{
		&RegistrantInfoEmitter{Emitter.NewEmitter()},
		&RegistrantInfoEmitter{Emitter.NewEmitter()},
		&RegistrantInfoEmitter{Emitter.NewEmitter()},
		Emitter.NewErrorEmitter()}
}

func (e *events) EnableAll() {
	e.NewConnection.Enable()
	e.UpdateConnection.Enable()
	e.Disconnection.Enable()
	e.Error.Enable()
}

func (e *events) DisableAll() {
	e.NewConnection.Disable()
	e.UpdateConnection.Disable()
	e.Disconnection.Disable()
	e.Error.Enable()
}
