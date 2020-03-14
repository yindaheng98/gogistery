package registrant

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//WaitTimeoutRetryNController controls the wait time, send time limit and retry limit of the request sending.
type WaitTimeoutRetryNController interface {

	//This method receive a response,
	//the time consuming from send last request to receive the response,
	//the retry time before received the response, and returns next wait time,
	//the next send time limit, and retry limit.
	GetWaitTimeoutRetryN(response protocol.Response, lastSendTime time.Duration, lastRetryN uint64) (waitTime time.Duration, sendTimeLimit time.Duration, nextRetryN uint64)
}
