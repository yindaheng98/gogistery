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

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *Requester) Send(option ProtocolRequestSendOption, timeout time.Duration, retryN int64) (Response, error) {
	for i := retryN; i > 0; i-- {
		response, err := r.SendOnce(option, timeout)
		if err == nil {
			return response, nil
		}
		r.Events.Retry.Emit(option, err)
	}
	return nil, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *Requester) SendOnce(option ProtocolRequestSendOption, timeout time.Duration) (Response, error) {
	responseChan := make(chan ReceivedResponse, 1)
	defer func() {
		defer func() { recover() }()
		close(responseChan) //退出时关闭通道
	}()
	requestChan := make(chan ProtocolRequestSendOption, 1)
	defer func() {
		defer func() { recover() }()
		close(requestChan) //退出时关闭通道
	}()
	go r.proto.Request(requestChan, responseChan) //异步执行发送操作
	requestChan <- option
	select {
	case response := <-responseChan:
		return response.response, response.error
	case <-time.After(timeout):
		return nil, errors.New("send timeout")
	}
}
