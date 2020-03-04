package responser

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
)

type Heart struct {
	beater       HeartBeater
	responser    *responser
	ErrorHandler func(error)
}

func NewHeart(beater HeartBeater, ResponseProto protocol.ResponseProtocol) *Heart {
	return &Heart{beater,
		newResponser(ResponseProto),
		func(error) {}}
}

//开始接收心跳，直到主动停止
func (h *Heart) RunBeating(ctx context.Context) {
	for {
		var request protocol.Request
		var err error
		var responseFunc func(protocol.TobeSendResponse)
		okChan := make(chan bool, 1)
		go func() {
			request, err, responseFunc = h.responser.Recv(ctx)
			okChan <- true
		}()
		select {
		case <-okChan:
			if err != nil {
				h.ErrorHandler(err)
			} else {
				responseFunc(h.beater.Beat(request))
			}
		case <-ctx.Done():
			return
		}
	}
}
