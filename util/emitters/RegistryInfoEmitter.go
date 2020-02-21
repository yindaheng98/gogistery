package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/protocol"
)

type RegistryInfoEmitter struct {
	*Emitter.Emitter
}

func NewRegistryInfoEmitter() *RegistryInfoEmitter {
	return &RegistryInfoEmitter{Emitter.NewEmitter()}
}

func (e *RegistryInfoEmitter) AddHandler(handler func(info protocol.RegistryInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(protocol.RegistryInfo))
	})
}

func (e *RegistryInfoEmitter) Emit(info protocol.RegistryInfo) {
	e.Emitter.Emit(info)
}
