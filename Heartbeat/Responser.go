package Heartbeat

import (
	"errors"
	"time"
)

type Responser struct {
	proto ResponseProtocol
}

func NewResponser(proto ResponseProtocol) *Responser {
	return &Responser{proto}
}

type ResponserOption struct {
	responseOption ResponseOption
	timeout        time.Duration
}

//传入一个响应channel，此channel将在函数返回接收到的Request后等待要响应的Response
func (r *Responser) Recv(responseChan <-chan Response, option ResponserOption) (Request, error) {
	requestProtoChan := make(chan RequestChanElement, 1)
	defer func() {
		defer func() { recover() }()
		close(requestProtoChan) //退出时关闭通道
	}()
	responseProtoChan := make(chan Response, 1)
	go r.proto.Response(requestProtoChan, option.responseOption, responseProtoChan) //异步执行Protocol的接收协议
	requestEl, ok := <-requestProtoChan                                             //等待接收数据到达
	if !ok {                                                                        //如果通道已关闭
		return nil, errors.New("request channel closed unexpectedly") //则返回错误
	}

	go func() { //然后新开一个线程等待响应数据
		defer func() {
			defer func() { recover() }()
			close(responseProtoChan) //退出时关闭通道
		}()
		select {
		case response, ok := <-responseChan: //等待高层channel返回响应数据
			if ok { //如果正常返回
				responseProtoChan <- response //则传入到底层协议
			}
		case <-time.After(option.timeout): //或超时
		}
	}()
	return requestEl.request, requestEl.error //返回收到的Request
}
