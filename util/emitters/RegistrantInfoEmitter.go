package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type RegistrantInfoEmitter struct {
	*Emitter.Emitter
}

func NewRegistrantInfoEmitter() *RegistrantInfoEmitter {
	return &RegistrantInfoEmitter{Emitter.NewEmitter()}
}

func (e *RegistrantInfoEmitter) AddHandler(handler func(info Protocol.RegistrantInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Protocol.RegistrantInfo))
	})
}

func (e *RegistrantInfoEmitter) Emit(info Protocol.RegistrantInfo) {
	e.Emitter.Emit(info)
}
