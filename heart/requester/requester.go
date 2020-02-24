package requester

import (
	"errors"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type requester struct {
	proto        protocol.RequestProtocol
	RetryHandler func(protocol.TobeSendRequest, error)
}

func newRequester(proto protocol.RequestProtocol) *requester {
	return &requester{proto, func(protocol.TobeSendRequest, error) {}}
}

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *requester) Send(option protocol.TobeSendRequest, timeout *time.Duration, retryN *uint64) (protocol.Response, error) {
	lastTimeout := time.Duration(0)
	lastRetryN := uint64(0)
	defer func() {
		*timeout = lastTimeout
		*retryN = lastRetryN
	}()
	for ; lastRetryN <= *retryN; lastRetryN++ {
		lastTimeout = *timeout
		response, err := r.SendOnce(option, &lastTimeout)
		if err == nil {
			return response, nil
		}
		r.RetryHandler(option, err)
	}
	return protocol.Response{}, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *requester) SendOnce(request protocol.TobeSendRequest, timeout *time.Duration) (protocol.Response, error) {
	startTime := time.Now()
	defer func() { *timeout = time.Now().Sub(startTime) }()
	responseChan := make(chan protocol.ReceivedResponse, 1)
	defer close(responseChan)
	requestChan := make(chan protocol.TobeSendRequest, 1)
	go r.proto.Request(requestChan, responseChan) //异步执行发送操作
	requestChan <- request
	close(requestChan)
	select {
	case response := <-responseChan:
		return response.Response, response.Error
	case <-time.After(*timeout):
		return protocol.Response{}, errors.New("send timeout")
	}
}
