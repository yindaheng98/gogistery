package Heart

import (
	"gogistery/Protocol"
	"time"
)

type RequestSendOption interface {
	GetTimeout() time.Duration
	GetRetryN() uint64
	String() string
} //自定义请求发送设置

type RequesterHeartProtocol interface {
	//对接上层消息策略
	//
	//输入一个Beat数据响应和下一个Beat处理函数，处理响应并生成下一个Beat数据请求
	Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest))
}

type ResponserHeartProtocol interface {
	//对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Beat数据请求，处理请求并生成Beat数据响应
	Beat(request Protocol.Request) Protocol.TobeSendResponse
}
