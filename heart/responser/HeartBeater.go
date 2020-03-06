package responser

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
)

type HeartBeater interface {
	//对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Beat数据请求，处理请求并生成Beat数据响应
	Beat(ctx context.Context, request protocol.Request) protocol.TobeSendResponse
}
