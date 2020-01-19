package Emitter

import (
	"gogistery/base"
	"gogistery/util/Emitter"
)

type LinkInfoEmitter struct {
	Emitter.Emitter
}

func NewLinkInfoEmitter() *LinkInfoEmitter {
	return &LinkInfoEmitter{*Emitter.New()}
}

func (e *LinkInfoEmitter) AddHandler(handler func(info base.LinkInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(base.LinkInfo))
	})
}

func (e *LinkInfoEmitter) Emit(err base.LinkInfo) {
	e.Emitter.Emit(err)
}
