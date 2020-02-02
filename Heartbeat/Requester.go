package Heartbeat

import (
	"errors"
	"time"
)

type requesterEvents struct {
	Retry *RequestOptionErrorEmitter
}

type Requester struct {
	proto  RequestProtocol
	Events *requesterEvents
}

func NewRequester(proto RequestProtocol) *Requester {
	return &Requester{proto, &requesterEvents{NewRequestOptionErrorEmitter()}}
}

type RequesterOption struct {
	requestOption RequestOption //发送设置
	timeout       time.Duration //超时时间
	retryN        int64         //重试次数
}

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *Requester) Send(request Request, option RequesterOption) (Response, error) {
	for i := option.retryN; i > 0; i-- {
		response, err := r.SendOnce(request, option)
		if err == nil {
			return response, nil
		}
		r.Events.Retry.Emit(option, err)
	}
	return nil, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *Requester) SendOnce(request Request, option RequesterOption) (Response, error) {
	responseChan := make(chan ResponseChanElement, 1)
	go func() { //异步执行发送函数
		response, err := r.proto.Send(request, option.requestOption) //异步执行发送操作
		responseChan <- ResponseChanElement{response, err}
	}()
	go func() { //异步执行超时检测函数
		defer func() {
			_ = recover()
		}()
		time.Sleep(option.timeout)                                           //等待一段时间
		responseChan <- ResponseChanElement{nil, errors.New("send timeout")} //发送超时信息
	}()
	res := <-responseChan //接收信息
	close(responseChan)   //然后关闭通道
	return res.response, res.error
}
