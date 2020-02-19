package RequesterHeart

import (
	"gogistery/Protocol"
	"time"
)

type RequesterHeart struct {
	proto     RequesterHeartProtocol
	requester *requester
	Event     *events
}

func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestProtocol) *RequesterHeart {
	heart := &RequesterHeart{heartProto, nil, newEvents()}
	heart.requester = newRequester(beatProto, heart)
	return heart
}

//开始心跳，直到最后由协议主动停止心跳或出错才返回
func (h *RequesterHeart) RunBeating(initRequest Protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error {
	request, timeout, retryN := initRequest, initTimeout, initRetryN
	var err error = nil
	established := false
	defer func() {
		if established {
			h.Event.Disconnection.Emit(request, err)
		}
	}()
	run := true
	for run {
		response, err := h.requester.Send(request, timeout, retryN)
		if err != nil {
			h.Event.Error.Emit(err)
			return err
		}
		if established { //如果已经达成过连接就触发更新事件
			h.Event.UpdateConnection.Emit(response)
		}
		run = false
		h.proto.Beat(response, func(requestB Protocol.TobeSendRequest, timeoutB time.Duration, retryNB uint64) {
			request, timeout, retryN = requestB, timeoutB, retryNB
			run = true
		})
		if run { //只有上级协议判定可以继续进行接下来的连接才能视为连接达成
			if !established { //此时可以触发新建连接事件
				h.Event.NewConnection.Emit(response)
			}
			established = true //并且设置连接达成标记
		}
	}
	return nil
}
