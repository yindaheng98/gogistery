package requester

import (
	"gogistery/protocol"
	"time"
)

type HeartBeater interface {
	//上一次发送重试了lastRetryN次，最后一次请求花费了lastTimeout发送完毕
	//
	//处理响应并生成下一个Beat数据请求，输入到处理函数beat中
	Beat(response protocol.Response, lastTimeout time.Duration, lastRetryN uint64,
		beat func(protocol.TobeSendRequest, time.Duration, uint64))
}
