package Heartbeat

import (
	"errors"
)

type Responser struct {
	proto ResponseBeatProtocol
}

func NewResponser(proto ResponseBeatProtocol) *Responser {
	return &Responser{proto}
}

//此channel将返回接收到的Request和一个处理Response的函数
func (r *Responser) Recv() (Request, error, func(TobeSendResponse)) {
	requestProtoChan := make(chan ReceivedRequest, 1)
	defer func() {
		defer func() { recover() }()
		close(requestProtoChan) //退出时关闭通道
	}()
	responseProtoChan := make(chan TobeSendResponse, 1)
	go r.proto.Response(requestProtoChan, responseProtoChan) //异步执行Protocol的接收协议
	request, ok := <-requestProtoChan                        //等待接收数据到达

	//response处理函数
	responseFunc := func(response TobeSendResponse) {
		defer func() { recover() }()
		responseProtoChan <- response //传入到底层协议
		close(responseProtoChan)      //退出时关闭通道
	}

	if !ok { //如果通道已关闭
		return nil, errors.New("request channel closed unexpectedly"), responseFunc //则返回错误
	}
	return request.Request, request.Error, responseFunc //返回收到的Request
}
