package requester

import (
	"context"
	"errors"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//Heart is the request controller.
//It can send request ("beat") to registry in a loop,
//until the registry reject it or occured some error.
type Heart struct {
	beater    HeartBeater
	requester *requester

	//Handlers consists of a series of functions.
	//The function will be called when:
	//
	//* `NewConnectionHandler`: the 2nd message to a registry sent successfully
	//
	//* `UpdateConnectionHandler`: the 3rd and later message to a registry sent successfully
	//
	//* `DisconnectionHandler`: `Heart.RunBeating(...)` exited
	//
	//* `RetryHandler`: message sent failed and before retransmission
	Handlers *handlers
}

//NewHeart returns the pointer to a Heart
func NewHeart(beater HeartBeater, RequestProto protocol.RequestProtocol) *Heart {
	heart := &Heart{beater,
		newRequester(RequestProto),
		newEvents()}
	heart.requester.RetryHandler = func(request protocol.TobeSendRequest, err error) {
		heart.Handlers.RetryHandler(request, err)
	}
	return heart
}

//The method RunBeating will send "beat" to registry in a loop,
//until the registry reject it or occured some error.
//It will do following things:
//
//1. Call the method `protocol.RequestProtocol.Request` you have implemented to send `initRequest` to registry, and get a response
//
//2. Call the method `HeartBeater.Beat` to get next request
//
//3. Call the method `protocol.RequestProtocol.Request` you have implemented to send the request to registry, and get the other response
//
//4. Loop step 2 and step 3 until some error occured or function `beat` was not called in `HeartBeater.Beat`
//
//When called, this function will block your goroutine until returned.
func (h *Heart) RunBeating(ctx context.Context,
	initRequest protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	request, Timeout, RetryN := initRequest, initTimeout, initRetryN
	established := false
	var lastResponse protocol.Response
	var err error = nil
	defer func() {
		if established {
			h.Handlers.DisconnectionHandler(lastResponse, err)
		}
	}()
	run := true
	for run {
		okChan := make(chan error, 1)
		go func() {
			defer close(okChan)
			response, err, timeout, retryN := h.requester.Send(ctx, request, Timeout, RetryN)
			if err != nil {
				okChan <- err
				return
			}
			lastResponse = response
			if established { //如果已经达成过连接就触发更新事件
				h.Handlers.UpdateConnectionHandler(response)
			}
			run = false
			h.beater.Beat(ctx, response, timeout, retryN,
				func(requestB protocol.TobeSendRequest, timeoutB time.Duration, retryNB uint64) {
					request, Timeout, RetryN = requestB, timeoutB, retryNB
					run = true
				})
			if run { //只有上级协议判定可以继续进行接下来的连接才能视为连接达成
				if !established { //此时可以触发新建连接事件
					h.Handlers.NewConnectionHandler(response)
				}
				established = true //并且设置连接达成标记
			}
		}()
		select {
		case err := <-okChan:
			if err != nil {
				return err
			}
		case <-ctx.Done(): //被要求退出
			request.Request.Disconnect = true //就发送断连信号
			_, _, _, _ = h.requester.Send(context.Background(), request, Timeout, RetryN)
			return errors.New("exited by context") //发完再退出
		}
	}
	return nil
}
