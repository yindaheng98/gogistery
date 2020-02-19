package RequesterHeart

import (
	"github.com/yindaheng98/go-utility/Emitter"
	"gogistery/util/emitters"
)

type events struct {
	NewConnection    *emitters.ResponseEmitter
	UpdateConnection *emitters.ResponseEmitter
	Disconnection    *emitters.TobeSendRequestErrorEmitter
	Retry            *emitters.TobeSendRequestErrorEmitter
	Error            *Emitter.ErrorEmitter
}

func newEvents() *events {
	return &events{
		emitters.NewResponseEmitter(),
		emitters.NewResponseEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
		emitters.NewTobeSendRequestErrorEmitter(),
		Emitter.NewErrorEmitter()}
}
