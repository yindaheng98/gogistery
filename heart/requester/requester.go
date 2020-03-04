package requester

import (
	"context"
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
func (r *requester) Send(ctx context.Context, option protocol.TobeSendRequest, timeout time.Duration, retryN uint64) (
	res protocol.Response, e error, lastTimeout time.Duration, totalRetryN uint64) {
	for totalRetryN = 0; totalRetryN <= retryN; totalRetryN++ {
		res, e, lastTimeout = r.SendOnce(ctx, option, timeout) //试一次
		if e == nil {                                          //不出错
			return //就返回
		}
		r.RetryHandler(option, e) //出错就报错
	}
	return protocol.Response{}, errors.New("connection failed"), lastTimeout, totalRetryN
}

//发送并等待回复，直到成功或超时
func (r *requester) SendOnce(ctx context.Context, request protocol.TobeSendRequest, timeout time.Duration) (
	res protocol.Response, e error, totalTimeout time.Duration) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)     //先构造一个超时ctx
	defer cancel()                                              //退出时cancel
	startTime := time.Now()                                     //记下启动时间
	defer func() { totalTimeout = time.Now().Sub(startTime) }() //退出时记录启动到退出的时间
	responseChan := make(chan protocol.ReceivedResponse, 1)
	defer close(responseChan)
	requestChan := make(chan protocol.TobeSendRequest, 1)
	go r.proto.Request(timeoutCtx, requestChan, responseChan) //异步执行发送操作
	requestChan <- request                                    //放入请求
	close(requestChan)                                        //然后关闭
	select {
	case response := <-responseChan:
		return response.Response, response.Error, timeout
	case <-timeoutCtx.Done():
		return protocol.Response{}, errors.New("send timeout"), timeout
	}
}
