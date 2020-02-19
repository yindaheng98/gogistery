package emitters

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/Protocol"
)

type ResponseEmitter struct {
	*Emitter.Emitter
}

func NewResponseEmitter() *ResponseEmitter {
	return &ResponseEmitter{Emitter.NewEmitter()}
}

func (e *ResponseEmitter) AddHandler(handler func(Protocol.Response)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Protocol.Response))
	})
}

func (e *ResponseEmitter) Emit(response Protocol.Response) {
	e.Emitter.Emit(response)
}
