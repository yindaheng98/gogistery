package Heartbeat

import (
	"github.com/yindaheng98/go-utility/Emitter"
)

type RequestOptionErrorEmitter struct {
	*Emitter.ErrorEmitter
}

func NewRequestOptionErrorEmitter() *RequestOptionErrorEmitter {
	return &RequestOptionErrorEmitter{Emitter.NewErrorEmitter()}
}

func (e *RequestOptionErrorEmitter) AddHandler(handler func(o RequestOption, err error)) {
	e.ErrorEmitter.AddHandler(func(i interface{}, err error) {
		handler(i.(RequestOption), err)
	})
}

func (e *RequestOptionErrorEmitter) Emit(o RequestOption, err error) {
	e.ErrorEmitter.Emit(o, err)
}
