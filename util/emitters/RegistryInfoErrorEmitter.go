package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

type RegistryInfoErrorEmitter struct {
	*Emitter.ErrorInfoEmitter
}

func NewRegistryInfoErrorEmitter() *RegistryInfoErrorEmitter {
	return &RegistryInfoErrorEmitter{Emitter.NewAsyncErrorInfoEmitter()}
}

func (e *RegistryInfoErrorEmitter) AddHandler(handler func(info protocol.RegistryInfo, err error)) {
	e.ErrorInfoEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(protocol.RegistryInfo), err)
	})
}
func (e *RegistryInfoErrorEmitter) Emit(info protocol.RegistryInfo, err error) {
	e.ErrorInfoEmitter.Emit(info, err)
}
