package requester

import (
	"gogistery/protocol"
)

type handlers struct {
	NewConnectionHandler    func(protocol.Response)
	UpdateConnectionHandler func(protocol.Response)
	DisconnectionHandler    func(protocol.TobeSendRequest, error)
	RetryHandler            func(protocol.TobeSendRequest, error)
}

func newEvents() *handlers {
	return &handlers{func(protocol.Response) {},
		func(protocol.Response) {},
		func(protocol.TobeSendRequest, error) {},
		func(protocol.TobeSendRequest, error) {}}
}
