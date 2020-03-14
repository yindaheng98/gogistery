package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"github.com/yindaheng98/gogistry/protocol"
)

//RegistryInfoErrorEmitter use protocol.RegistryInfo and error as event payload
type RegistryInfoErrorEmitter struct {
	*Emitter.ErrorInfoEmitter
}

//NewRegistryInfoErrorEmitter returns the pointer to a async RegistryInfoErrorEmitter
func NewRegistryInfoErrorEmitter() *RegistryInfoErrorEmitter {
	return &RegistryInfoErrorEmitter{Emitter.NewAsyncErrorInfoEmitter()}
}

//Implementation of Emitter.AddHandler
func (e *RegistryInfoErrorEmitter) AddHandler(handler func(info protocol.RegistryInfo, err error)) {
	e.ErrorInfoEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(protocol.RegistryInfo), err)
	})
}

//Implementation of Emitter.Emit
func (e *RegistryInfoErrorEmitter) Emit(info protocol.RegistryInfo, err error) {
	e.ErrorInfoEmitter.Emit(info, err)
}
