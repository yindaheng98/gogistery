package registrant

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type RetryNController interface {
	GetWaitTimeoutRetryN(response protocol.Response, timeout time.Duration, retryN uint64) (time.Duration, time.Duration, uint64)
}
