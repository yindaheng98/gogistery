package requester

import (
	"github.com/yindaheng98/gogistry/protocol"
)

type handlers struct {
	NewConnectionHandler    func(protocol.Response)
	UpdateConnectionHandler func(protocol.Response)
	DisconnectionHandler    func(protocol.Response, error)
	RetryHandler            func(protocol.TobeSendRequest, error)
}

func newEvents() *handlers {
	return &handlers{
		func(protocol.Response) {},
		func(protocol.Response) {},
		func(protocol.Response, error) {},
		func(protocol.TobeSendRequest, error) {},
	}
}
