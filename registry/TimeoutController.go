package registry

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type TimeoutController interface {
	TimeoutForNew(request protocol.Request) time.Duration
	TimeoutForUpdate(request protocol.Request) time.Duration
}
