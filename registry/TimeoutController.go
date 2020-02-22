package registry

import (
	"gogistery/protocol"
	"time"
)

type TimeoutController interface {
	TimeoutForNew(request protocol.Request) time.Duration
	TimeoutForUpdate(request protocol.Request) time.Duration
}
