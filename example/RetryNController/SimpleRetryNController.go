package RetryNController

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//SimpleRetryNController is a simple implementation of WaitTimeoutRetryNController.
type SimpleRetryNController struct{}

//sendTimeLimit=lastSendTime * 2
//
//
//nextRetryN=(lastRetryN + 3) * 2
//
//waitTime=response.GetTimeout() - nextRetryN * sendTimeLimit (if response.GetTimeout() > nextRetryN * sendTimeLimit)
//
//waitTime=0 (if response.GetTimeout() <= nextRetryN * sendTimeLimit)
func (c SimpleRetryNController) GetWaitTimeoutRetryN(response protocol.Response, lastSendTime time.Duration, lastRetryN uint64) (waitTime time.Duration, sendTimeLimit time.Duration, nextRetryN uint64) {
	totalTime := response.GetTimeout()
	expectRetryN := (lastRetryN + 3) * 2                          //预计的重试次数
	expectTimeout := lastSendTime * 2                             //预计的单次发送时长
	expectSendTime := time.Duration(expectRetryN) * expectTimeout //期望的发送总耗时
	if totalTime <= expectSendTime {
		return 0, expectTimeout, uint64(totalTime / expectTimeout)
	}
	return totalTime - expectSendTime, expectTimeout, expectRetryN
}
