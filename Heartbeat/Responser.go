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

type ResponserOption interface {
	GetResponseOption() ResponseOption
	GetTimeout() time.Duration
}

//传入一个响应channel，此channel将在函数返回接收到的Request后等待要响应的Response
func (r *Responser) Recv(responseChan chan Response, option ResponserOption) (Request, error) {
	requestProtoChan := make(chan RequestChanElement, 1)
	responseProtoChan := make(chan Response, 1)                                      //构造Protocol中要用的channel
	go r.proto.Recv(responseProtoChan, option.GetResponseOption(), requestProtoChan) //调用Protocol的接收协议
	requestEl, ok := <-requestProtoChan                                              //等待接收数据到达
	if !ok {                                                                         //如果通道已关闭
		return nil, errors.New("request channel closed unexpectedly") //则返回错误
	}
	close(requestProtoChan)                                  //接收数据正确到达后关闭Protocol中的接收channel
	go waitResponse(responseChan, responseProtoChan, option) //然后新开一个线程等待响应数据
	return requestEl.request, requestEl.error                //并返回收到的Request
}

//等待等待响应数据的线程
//
//等待高层responseChan传入的Response，解析后传入到底层协议responseProtoChan中
func waitResponse(responseChan chan Response, responseProtoChan chan Response, option ResponserOption) {
	defer func() { recover() }()
	defer func() {
		defer func() { recover() }()
		close(responseChan)
	}()
	defer func() {
		defer func() { recover() }()
		close(responseProtoChan)
	}()
	select {
	case response, ok := <-responseChan: //等待高层channel返回响应数据
		if ok { //如果正常返回
			responseProtoChan <- response //则传入到底层协议
		}
	case <-time.After(option.GetTimeout()): //或超时
	}
}
