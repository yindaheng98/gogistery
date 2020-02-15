package Registrant

import (
	"gogistery/Protocol"
	"time"
)

type requesterHeartProtocol struct {
	heart            *heart //此协议服务于哪个heart
	stopChan         chan bool
	disconnectedChan chan bool
}

func newRequesterHeartProtocol(heart *heart) *requesterHeartProtocol {
	stopChan := make(chan bool, 1)
	disconnectedChan := make(chan bool, 1)
	close(stopChan)
	close(disconnectedChan)
	return &requesterHeartProtocol{heart: heart, stopChan: stopChan, disconnectedChan: disconnectedChan}
}
func (p *requesterHeartProtocol) start() {
	p.stop() //启动前必须先停止
	p.stopChan = make(chan bool, 1)
	p.disconnectedChan = make(chan bool, 1)
}
func (p *requesterHeartProtocol) stop() {
	defer func() { recover() }()
	p.stopChan <- true   //发送停止信息
	<-p.disconnectedChan //等待已停止信息
	close(p.stopChan)    //关闭停止信息通道
}
func (p *requesterHeartProtocol) Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest, time.Duration, uint64)) {
	request := p.heart.beatResponse(response)
	if response.IsReject() { //如果注册中心拒绝了连接请求
		p.disconnectedChan <- true
		close(p.disconnectedChan)
		defer func() { recover() }()
		close(p.stopChan)
		return //就直接断连退出
	}
	select {
	case <-time.After(response.GetTimeout() / 2): //等待一段时间再发，这里的等待时间应该小一点以免后续操作中发送时间不够
		request.Request.Disconnect = false
		beat(request, response.GetTimeout()/2, response.GetRetryN())
	case <-p.stopChan: //突然要求停机
		request.Request.Disconnect = true //那就发送断连信号
		beat(request, response.GetTimeout()/2, response.GetRetryN())
		p.disconnectedChan <- true //已断连标志位
		defer func() { recover() }()
		close(p.disconnectedChan)
	case <-p.disconnectedChan: //如果已断连
		return //就直接退出
	}
}
