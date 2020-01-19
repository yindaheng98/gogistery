package Emitter

import (
	"gogistery/base/Error"
	"gogistery/util/Emitter"
)

type ConnectionErrorEmitter struct {
	Emitter.Emitter
}

func NewConnectionErrorEmitter() *ConnectionErrorEmitter {
	return &ConnectionErrorEmitter{*Emitter.New()}
}

func (e *ConnectionErrorEmitter) AddHandler(handler func(Error.ConnectionError)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Error.ConnectionError))
	})
}

func (e *ConnectionErrorEmitter) Emit(err Error.ConnectionError) {
	e.Emitter.Emit(err)
}
