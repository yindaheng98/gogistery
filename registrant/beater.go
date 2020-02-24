package registrant

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type beater struct {
	heart            *heart //此协议服务于哪个heart
	retryNController RetryNController
	stopChan         chan bool
	stoppedChan      chan bool
}

func newBeater(heart *heart, requestController RetryNController) *beater {
	stopChan := make(chan bool, 1)
	stoppedChan := make(chan bool, 1)
	close(stopChan)
	close(stoppedChan)
	return &beater{heart, requestController, stopChan, stoppedChan}
}
func (p *beater) Start() {
	p.Stop() //启动前必须先停止
	p.stopChan = make(chan bool, 1)
	p.stoppedChan = make(chan bool, 1)
}
func (p *beater) Stop() {
	defer func() { recover() }()
	p.stopChan <- true //发送停止信息
	<-p.stoppedChan    //等待已停止信息
	close(p.stopChan)  //关闭停止信息通道
}
func (p *beater) Beat(response protocol.Response, lastTimeout time.Duration, lastRetryN uint64, beat func(protocol.TobeSendRequest, time.Duration, uint64)) {
	request := p.heart.register(response)
	if response.IsReject() { //如果注册中心拒绝了连接请求
		defer func() { recover() }()
		p.stoppedChan <- true
		close(p.stoppedChan)
		close(p.stopChan)
		return //就直接断连退出
	}
	waitTime, sendTimeout, retryN := p.retryNController.GetWaitTimeoutRetryN(response, lastTimeout, lastRetryN)
	select {
	case <-time.After(waitTime): //等待一段时间再发，这里的等待时间应该小一点以免后续操作中发送时间不够
		request.Request.Disconnect = false
		beat(request, sendTimeout, retryN)
	case <-p.stopChan: //突然要求停机
		request.Request.Disconnect = true //那就发送断连信号
		beat(request, sendTimeout, retryN)
		p.stoppedChan <- true //发送已断连标志位
		defer func() { recover() }()
		close(p.stoppedChan)
	case <-p.stoppedChan: //如果已断连
		return //就直接退出
	}
}
