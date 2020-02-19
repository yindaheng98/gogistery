package RequesterHeart

import (
	"gogistery/util/emitters"
)

type events struct {
	NewConnection    *emitters.RegistryInfoEmitter
	UpdateConnection *emitters.RegistryInfoEmitter
	Retry            *emitters.TobeSendRequestErrorEmitter
	Disconnection    *emitters.RegistryInfoEmitter
}
