package RetryNController

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type SimpleRetryNController struct{}

func (c SimpleRetryNController) GetWaitTimeoutRetryN(response protocol.Response, lastTimeout time.Duration, lastRetryN uint64) (time.Duration, time.Duration, uint64) {
	totalTime := response.GetTimeout()
	expectRetryN := (lastRetryN + 3) * 2                          //预计的重试次数
	expectTimeout := lastTimeout * 2                              //预计的单次发送时长
	expectSendTime := time.Duration(expectRetryN) * expectTimeout //期望的发送总耗时
	if totalTime <= expectSendTime {
		return 0, expectTimeout, uint64(totalTime / expectTimeout)
	}
	return totalTime - expectSendTime, expectTimeout, expectRetryN
}
