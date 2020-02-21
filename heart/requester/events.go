package requester

import (
	"gogistery/Protocol"
)

type handlers struct {
	NewConnectionHandler    func(Protocol.Response)
	UpdateConnectionHandler func(Protocol.Response)
	DisconnectionHandler    func(Protocol.TobeSendRequest, error)
	RetryHandler            func(Protocol.TobeSendRequest, error)
}

func newEvents() *handlers {
	return &handlers{func(Protocol.Response) {},
		func(Protocol.Response) {},
		func(Protocol.TobeSendRequest, error) {},
		func(Protocol.TobeSendRequest, error) {}}
}
