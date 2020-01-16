package Emitter

import "gogistery/util/Emitter"

type ErrorEmitter struct {
	Emitter.Emitter
}

func NewErrorEmitter() *ErrorEmitter {
	return &ErrorEmitter{*Emitter.New()}
}

func (e *ErrorEmitter) AddHandler(handler func(error)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(error))
	})
}

func (e *ErrorEmitter) Emit(err error) {
	e.Emitter.Emit(err)
}
