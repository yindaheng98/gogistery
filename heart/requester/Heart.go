package requester

import (
	"gogistery/Protocol"
	"time"
)

type Heart struct {
	beater    HeartBeater
	requester *requester
	Handlers  *handlers
}

func NewHeart(beater HeartBeater, RequestProto Protocol.RequestProtocol) *Heart {
	heart := &Heart{beater, nil, newEvents()}
	heart.requester = newRequester(RequestProto)
	heart.requester.RetryHandler = func(request Protocol.TobeSendRequest, err error) {
		heart.Handlers.RetryHandler(request, err)
	}
	return heart
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *Heart) RunBeating(initRequest Protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	request, timeout, retryN := initRequest, initTimeout, initRetryN
	var err error = nil
	established := false
	defer func() {
		if established {
			h.Handlers.DisconnectionHandler(request, err)
		}
	}()
	run := true
	for run {
		response, err := h.requester.Send(request, timeout, retryN)
		if err != nil {
			return err
		}
		if established { //如果已经达成过连接就触发更新事件
			h.Handlers.UpdateConnectionHandler(response)
		}
		run = false
		h.beater.Beat(response, func(requestB Protocol.TobeSendRequest, timeoutB time.Duration, retryNB uint64) {
			request, timeout, retryN = requestB, timeoutB, retryNB
			run = true
		})
		if run { //只有上级协议判定可以继续进行接下来的连接才能视为连接达成
			if !established { //此时可以触发新建连接事件
				h.Handlers.NewConnectionHandler(response)
			}
			established = true //并且设置连接达成标记
		}
	}
	return nil
}
