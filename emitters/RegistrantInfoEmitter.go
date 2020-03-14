package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

//RegistrantInfoEmitter use protocol.RegistrantInfo as event payload
type RegistrantInfoEmitter struct {
	Emitter.Emitter
}

//NewRegistrantInfoEmitter returns the pointer to a sync RegistrantInfoEmitter
func NewRegistrantInfoEmitter() *RegistrantInfoEmitter {
	return &RegistrantInfoEmitter{Emitter.NewSyncEmitter()}
}

//Implementation of Emitter.AddHandler
func (e *RegistrantInfoEmitter) AddHandler(handler func(info protocol.RegistrantInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(protocol.RegistrantInfo))
	})
}

//Implementation of Emitter.Emit
func (e *RegistrantInfoEmitter) Emit(info protocol.RegistrantInfo) {
	e.Emitter.Emit(info)
}
