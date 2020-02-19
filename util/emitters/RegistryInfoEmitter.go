package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type RegistryInfoEmitter struct {
	*Emitter.Emitter
}

func NewRegistryInfoEmitter() *RegistryInfoEmitter {
	return &RegistryInfoEmitter{Emitter.NewEmitter()}
}

func (e *RegistryInfoEmitter) AddHandler(handler func(info Protocol.RegistryInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Protocol.RegistryInfo))
	})
}

func (e *RegistryInfoEmitter) Emit(info Protocol.RegistryInfo) {
	e.Emitter.Emit(info)
}
