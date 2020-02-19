package Registrant

import (
	"gogistery/Protocol"
	"sync/atomic"
)

type Registrant struct {
	Info   Protocol.RegistrantInfo //注册器的信息
	hearts []*heart                //heart列表，每一个heart负责向一个注册线程

	stopChan    chan bool //向各heart传递停止信息
	runningN    int64     //正在运行的线程数
	stoppedChan chan bool //线程全部停止则向这里传入已停止信息

	candProto CandidateRegistryProtocol //候选服务器选择协议
	Events    *events
}

func New(Info Protocol.RegistrantInfo, regitryN uint, candProto CandidateRegistryProtocol, sendProto Protocol.RequestProtocol) *Registrant {
	stopChan := make(chan bool, 1)
	stoppedChan := make(chan bool, 1)
	close(stopChan)
	close(stoppedChan)
	registrant := &Registrant{
		Info:   Info,
		hearts: make([]*heart, regitryN),

		stopChan:    stopChan,    //向各heart传递停止信息
		runningN:    0,           //正在运行的线程数
		stoppedChan: stoppedChan, //线程全部停止则向这里传入已停止信息

		candProto: candProto,
		Events:    newEvents(),
	}
	for i := uint(0); i < regitryN; i++ {
		registrant.hearts[i] = newHeart(registrant, sendProto)
	}
	return registrant
}

func (r *Registrant) Run() {
	r.stopChan = make(chan bool, 1)
	r.stoppedChan = make(chan bool, 1)
	connChan := make(chan []Protocol.RegistryInfo, 1)
	connChan <- make([]Protocol.RegistryInfo, len(r.hearts))
	for i, heart := range r.hearts {
		go func() {
			r.heartRoutine(heart, i, connChan)
			if atomic.LoadInt64(&r.runningN) <= 0 {
				r.stoppedChan <- true
				close(r.stoppedChan)
			}
		}()
	}
	<-r.stoppedChan
}

func (r *Registrant) Stop() {
	r.stopChan <- true
	close(r.stopChan)
	<-r.stoppedChan
}

func (r *Registrant) heartRoutine(h *heart, i int, connChan chan []Protocol.RegistryInfo) {
	atomic.AddInt64(&r.runningN, 1)
	defer atomic.AddInt64(&r.runningN, -1)
	for {
		errChan := make(chan error, 1)
		go func() { //新开一个线程运行注册程序
			var connections []Protocol.RegistryInfo
			connections = <-connChan           //从队列中取出已有连接列表
			var except []Protocol.RegistryInfo //去除空项
			for _, conn := range connections {
				if conn != nil {
					except = append(except, conn)
				}
			}
			initRegistryInfo, initTimeout, initRetryN := r.candProto.GetCandidate(except)
			//去除空项以此获取新连接
			connections[i] = initRegistryInfo //将新连接加入已有连接列表
			connChan <- connections           //已有连接列表放回队列
			err := h.Run(Protocol.TobeSendRequest{
				Request: Protocol.Request{RegistrantInfo: r.Info, Disconnect: false},
				Option:  initRegistryInfo.GetRequestSendOption()}, initTimeout, initRetryN)
			connections = <-connChan //从队列中取出已有连接列表
			connections[i] = nil     //将对应位置空
			connChan <- connections  //已有连接列表放回队列
			errChan <- err           //掉线时发送掉线信息
		}()
		select {
		case err := <-errChan: //若收到掉线信息
			if err != nil { //有错则报错
				r.Events.Error.Emit(err)
			}
		case <-r.stopChan: //若收到停止信息
			h.Stop() //则直接停止
			return   //并退出
		}
	}
}

func (r *Registrant) GetConnections() []Protocol.RegistryInfo {
	res := make([]Protocol.RegistryInfo, 0)
	for _, h := range r.hearts {
		if h.RegistryInfo != nil {
			res = append(res, h.RegistryInfo)
		}
	}
	return res
}
