package Heart

import (
	"errors"
	"gogistery/Protocol"
	"math"
	"time"
)

type requesterEvents struct {
	Retry *TobeSendRequestErrorEmitter
}

type Requester struct {
	proto  Protocol.RequestBeatProtocol
	Events *requesterEvents
}

func NewRequester(proto Protocol.RequestBeatProtocol) *Requester {
	return &Requester{proto, &requesterEvents{newTobeSendRequestErrorEmitter()}}
}

//多次重试发送并等待回复，直到成功或达到重试次数上限
func (r *Requester) Send(option Protocol.TobeSendRequest, timeout time.Duration, retryN uint64) (Protocol.Response, error) {
	timeout = time.Duration(math.Abs(float64(timeout)))
	outTime := time.Now().Add(timeout)                                           //何时过期
	timeoutOnce := time.Duration(math.Floor(float64(timeout) / float64(retryN))) //单次发送的超时时间
	deltaTimeoutOnce := time.Duration(0)                                         //一次发送超时到下一次发送开始的时间差
	lastTime := time.Now()                                                       //上一次发送结束的时间
	for i := retryN; i > 0; i-- {
		response, err := r.SendOnce(option, timeoutOnce-deltaTimeoutOnce)
		if err == nil {
			return response, nil
		}
		r.Events.Retry.Emit(option, err)
		t := time.Now()                                  //当前发送结束的时间
		deltaTimeoutOnce = t.Sub(lastTime) - timeoutOnce //比预计多花了多少时间
		timeoutOnce = time.Duration(math.Floor(float64(outTime.Sub(t)) / float64(i)))
		lastTime = t
	}
	return nil, errors.New("connection failed")
}

//发送并等待回复，直到成功或超时
func (r *Requester) SendOnce(option Protocol.TobeSendRequest, timeout time.Duration) (Protocol.Response, error) {
	responseChan := make(chan Protocol.ReceivedResponse, 1)
	defer func() {
		defer func() { recover() }()
		close(responseChan) //退出时关闭通道
	}()
	requestChan := make(chan Protocol.TobeSendRequest, 1)
	defer func() {
		defer func() { recover() }()
		close(requestChan) //退出时关闭通道
	}()
	go r.proto.Request(requestChan, responseChan) //异步执行发送操作
	requestChan <- option
	select {
	case response := <-responseChan:
		return response.Response, response.Error
	case <-time.After(timeout):
		return nil, errors.New("send timeout")
	}
}
