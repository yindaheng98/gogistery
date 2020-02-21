package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/protocol"
)

type RegistrantInfoEmitter struct {
	*Emitter.Emitter
}

func NewRegistrantInfoEmitter() *RegistrantInfoEmitter {
	return &RegistrantInfoEmitter{Emitter.NewEmitter()}
}

func (e *RegistrantInfoEmitter) AddHandler(handler func(info protocol.RegistrantInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(protocol.RegistrantInfo))
	})
}

func (e *RegistrantInfoEmitter) Emit(info protocol.RegistrantInfo) {
	e.Emitter.Emit(info)
}
