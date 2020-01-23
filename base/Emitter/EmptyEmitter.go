package Emitter

import "gogistery/util/Emitter"

//没有格式的空事件
type EmptyEmitter struct {
	*Emitter.Emitter
}

func NewEmptyEmitter() *EmptyEmitter {
	return &EmptyEmitter{Emitter.New()}
}

func (e *EmptyEmitter) AddHandler(handler func()) {
	e.Emitter.AddHandler(func(interface{}) {
		handler()
	})
}

func (e *EmptyEmitter) Emit() {
	e.Emitter.Emit(nil)
}
