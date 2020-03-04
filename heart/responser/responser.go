package responser

import (
	"context"
	"errors"
	"github.com/yindaheng98/gogistry/protocol"
)

type responser struct {
	proto protocol.ResponseProtocol
}

func newResponser(proto protocol.ResponseProtocol) *responser {
	return &responser{proto}
}

//此channel将返回接收到的Request和一个处理Response的函数
func (r *responser) Recv(ctx context.Context) (protocol.Request, error, func(protocol.TobeSendResponse)) {
	requestProtoChan := make(chan protocol.ReceivedRequest, 1)
	defer close(requestProtoChan) //退出时关闭通道
	responseProtoChan := make(chan protocol.TobeSendResponse, 1)
	go r.proto.Response(ctx, requestProtoChan, responseProtoChan) //异步执行Protocol的接收协议
	select {
	case request := <-requestProtoChan: //等待接收数据到达
		return request.Request, request.Error, func(response protocol.TobeSendResponse) { //response处理函数
			responseProtoChan <- response //传入到底层协议
			close(responseProtoChan)      //退出时关闭通道
		}
	case <-ctx.Done():
		err := errors.New("exited by context")
		return protocol.Request{}, err, func(protocol.TobeSendResponse) {} //则返回错误
	}
}
