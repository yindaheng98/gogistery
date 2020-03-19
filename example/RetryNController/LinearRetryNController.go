package RetryNController

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//LinearRetryNController is a simple implementation of WaitTimeoutRetryNController.
type LinearRetryNController struct {

	//nextRetryN=lastRetryN * K_RetryN + B_RetryN
	K_RetryN, B_RetryN uint64

	//nextSendTimeLimit=lastSendTime * K_SendTime + B_SendTime
	K_SendTime, B_SendTime time.Duration
}

func NewLinearRetryNController() *LinearRetryNController {
	return &LinearRetryNController{
		K_RetryN: 2, B_RetryN: 1,
		K_SendTime: 10, B_SendTime: 1e9,
	}
}

//waitTime=response.GetTimeout() - nextRetryN * sendTimeLimit (if response.GetTimeout() > nextRetryN * sendTimeLimit)
//
//waitTime=0 (if response.GetTimeout() <= nextRetryN * sendTimeLimit)
func (c LinearRetryNController) GetWaitTimeoutRetryN(response protocol.Response, lastSendTime time.Duration, lastRetryN uint64) (waitTime time.Duration, sendTimeLimit time.Duration, nextRetryN uint64) {
	totalTime := response.GetTimeout()
	expectRetryN := lastRetryN*c.K_RetryN + c.B_RetryN            //预计的重试次数
	expectTimeout := lastSendTime*c.K_SendTime + c.B_SendTime     //预计的单次发送时长
	expectSendTime := time.Duration(expectRetryN) * expectTimeout //期望的发送总耗时
	if totalTime <= expectSendTime {
		return 0, expectTimeout, uint64(totalTime / expectTimeout)
	}
	return totalTime - expectSendTime, expectTimeout, expectRetryN
}
