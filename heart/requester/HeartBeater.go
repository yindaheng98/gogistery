package requester

import (
	"gogistery/protocol"
	"time"
)

type HeartBeater interface {
	//对接上层消息策略
	//
	//输入一个Beat数据响应和下一个Beat处理函数，处理响应并生成下一个Beat数据请求
	Beat(response protocol.Response, beat func(protocol.TobeSendRequest, time.Duration, uint64))
}
