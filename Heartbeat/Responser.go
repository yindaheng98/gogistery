package Heartbeat

import (
	"errors"
)

type Responser struct {
	proto ResponseProtocol
}

func NewResponser(proto ResponseProtocol) *Responser {
	return &Responser{proto}
}

//传入一个响应channel和超时时间，此channel将在函数返回接收到的Request后等待要响应的Response
func (r *Responser) Recv(responseChan <-chan ProtocolResponseSendOption) (Request, error) {
	requestProtoChan := make(chan ReceivedRequest, 1)
	defer func() {
		defer func() { recover() }()
		close(requestProtoChan) //退出时关闭通道
	}()
	responseProtoChan := make(chan ProtocolResponseSendOption, 1)
	go r.proto.Response(requestProtoChan, responseProtoChan) //异步执行Protocol的接收协议
	request, ok := <-requestProtoChan                        //等待接收数据到达

	go func() { //数据到达后立即新开一个线程等待上层响应数据
		defer func() {
			defer func() { recover() }()
			close(responseProtoChan) //退出时关闭通道
		}()
		response := <-responseChan    //等待高层channel返回响应数据
		responseProtoChan <- response //传入到底层协议
	}()

	if !ok { //如果通道已关闭
		return nil, errors.New("request channel closed unexpectedly") //则返回错误
	}
	return request.request, request.error //返回收到的Request
}
