package Heart

import (
	"fmt"
	"gogistery/Protocol"
	"time"
)

type RequestSendOption struct {
	Timeout time.Duration
	RetryN  int64
} //自定义请求发送设置

func (r RequestSendOption) String() string {
	return fmt.Sprintf("Heart.RequestSendOption{timeout:%d,retryN:%d}", r.Timeout, r.RetryN)
}

type RequesterHeartProtocol interface {
	//对接上层消息策略
	//
	//输入一个Beat数据响应和下一个Beat处理函数，处理响应并生成下一个Beat数据请求
	Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest, RequestSendOption))
}

type ResponserHeartProtocol interface {
	//对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Beat数据请求，处理请求并生成Beat数据响应
	Beat(request Protocol.Request) Protocol.TobeSendResponse
}
