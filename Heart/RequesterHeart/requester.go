package RequesterHeart

import (
	"errors"
	"gogistery/Protocol"
	"time"
)

type requester struct {
	proto Protocol.RequestProtocol
	heart *RequesterHeart
}

func newRequester(proto Protocol.RequestProtocol, heart *RequesterHeart) *requester {
	return &requester{proto, heart}
}

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *requester) Send(option Protocol.TobeSendRequest, timeout time.Duration, retryN uint64) (Protocol.Response, error) {
	for ; retryN > 0; retryN-- {
		response, err := r.SendOnce(option, timeout)
		if err == nil {
			return response, nil
		}
		r.heart.Event.Retry.Emit(option, err)
	}
	return Protocol.Response{}, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *requester) SendOnce(request Protocol.TobeSendRequest, timeout time.Duration) (Protocol.Response, error) {
	responseChan := make(chan Protocol.ReceivedResponse, 1)
	requestChan := make(chan Protocol.TobeSendRequest, 1)
	go r.proto.Request(requestChan, responseChan) //异步执行发送操作
	requestChan <- request
	select {
	case response := <-responseChan:
		return response.Response, response.Error
	case <-time.After(timeout):
		return Protocol.Response{}, errors.New("send timeout")
	}
}
