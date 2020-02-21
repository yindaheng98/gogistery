package responser

import (
	"gogistery/Protocol"
)

type HeartBeater interface {
	//对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Beat数据请求，处理请求并生成Beat数据响应
	Beat(request Protocol.Request) Protocol.TobeSendResponse
}
