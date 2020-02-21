package requester

import (
	"errors"
	"gogistery/Protocol"
	"time"
)

type requester struct {
	proto        Protocol.RequestProtocol
	RetryHandler func(Protocol.TobeSendRequest, error)
}

func newRequester(proto Protocol.RequestProtocol) *requester {
	return &requester{proto, func(Protocol.TobeSendRequest, error) {}}
}

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *requester) Send(option Protocol.TobeSendRequest, timeout time.Duration, retryN uint64) (Protocol.Response, error) {
	for ; retryN > 0; retryN-- {
		response, err := r.SendOnce(option, timeout)
		if err == nil {
			return response, nil
		}
		r.RetryHandler(option, err)
	}
	return Protocol.Response{}, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *requester) SendOnce(request Protocol.TobeSendRequest, timeout time.Duration) (Protocol.Response, error) {
	responseChan := make(chan Protocol.ReceivedResponse, 1)
	defer func() {
		defer func() { recover() }()
		close(responseChan)
	}()
	requestChan := make(chan Protocol.TobeSendRequest, 1)
	go r.proto.Request(requestChan, responseChan) //异步执行发送操作
	requestChan <- request
	close(requestChan)
	select {
	case response := <-responseChan:
		return response.Response, response.Error
	case <-time.After(timeout):
		return Protocol.Response{}, errors.New("send timeout")
	}
}
