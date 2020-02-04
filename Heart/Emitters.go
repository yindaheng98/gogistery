package Heart

import (
	"github.com/yindaheng98/go-utility/Emitter"
)

type ErrorEmitter struct {
	*Emitter.Emitter
}

func newErrorEmitter() *ErrorEmitter {
	return &ErrorEmitter{Emitter.NewEmitter()}
}

func (e *ErrorEmitter) AddHandler(handler func(err error)) {
	e.Emitter.AddHandler(func(err interface{}) {
		handler(err.(error))
	})
}

func (e *ErrorEmitter) Emit(err error) {
	e.Emitter.Emit(err)
}
