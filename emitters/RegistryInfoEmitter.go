package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

//RegistryInfoEmitter use protocol.RegistryInfo as event payload
type RegistryInfoEmitter struct {
	Emitter.Emitter
}

//NewRegistryInfoEmitter returns the pointer to a sync RegistryInfoEmitter
func NewRegistryInfoEmitter() *RegistryInfoEmitter {
	return &RegistryInfoEmitter{Emitter.NewSyncEmitter()}
}

//Implementation of Emitter.AddHandler
func (e *RegistryInfoEmitter) AddHandler(handler func(info protocol.RegistryInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(protocol.RegistryInfo))
	})
}

//Implementation of Emitter.Emit
func (e *RegistryInfoEmitter) Emit(info protocol.RegistryInfo) {
	e.Emitter.Emit(info)
}
