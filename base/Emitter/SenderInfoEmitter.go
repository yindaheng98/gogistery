package Emitter

import (
	"gogistery/base"
	"gogistery/util/Emitter"
)

//事件格式为base.SenderInfo
type SenderInfoEmitter struct {
	Emitter.Emitter
}

func NewSenderInfoEmitter() *SenderInfoEmitter {
	return &SenderInfoEmitter{*Emitter.New()}
}

func (e *SenderInfoEmitter) AddHandler(handler func(info base.SenderInfo)) {
	e.Emitter.AddHandler(func(i interface{}) {
		handler(i.(base.SenderInfo))
	})
}

func (e *SenderInfoEmitter) Emit(err base.SenderInfo) {
	e.Emitter.Emit(err)
}
