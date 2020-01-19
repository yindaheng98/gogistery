package Emitter

import (
	"gogistery/base/Errors"
	"gogistery/util/Emitter"
)

type LinkErrorEmitter struct {
	Emitter.Emitter
}

func NewLinkErrorEmitter() *LinkErrorEmitter {
	return &LinkErrorEmitter{*Emitter.New()}
}

func (e *LinkErrorEmitter) AddHandler(handler func(Errors.LinkError)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(Errors.LinkError))
	})
}

func (e *LinkErrorEmitter) Emit(err Errors.LinkError) {
	e.Emitter.Emit(err)
}
