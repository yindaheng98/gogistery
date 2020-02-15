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

func New(Info Protocol.RegistrantInfo, regitryN uint64, candProto CandidateRegistryProtocol, sendProto Protocol.RequestProtocol) *Registrant {
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
	for i := uint64(0); i < regitryN; i++ {
		registrant.hearts[i] = newHeart(registrant, sendProto)
	}
	return registrant
}

func (r *Registrant) Run() {
	r.stopChan = make(chan bool, 1)
	r.stoppedChan = make(chan bool, 1)
	for _, heart := range r.hearts {
		go func() {
			r.heartRoutine(heart)
			if atomic.LoadInt64(&r.runningN) <= 0 {
				r.stoppedChan <- true
				close(r.stoppedChan)
			}
		}()
	}
}

func (r *Registrant) Stop() {
	r.stopChan <- true
	close(r.stopChan)
	<-r.stoppedChan
}

func (r *Registrant) heartRoutine(h *heart) {
	atomic.AddInt64(&r.runningN, 1)
	defer atomic.AddInt64(&r.runningN, -1)
	for {
		errChan := make(chan error, 1)
		go func() { //新开一个线程运行注册程序
			initRequestSendOption, initTimeout, initRetryN := r.candProto.NewInitRequestSendOption()
			err := h.Run(Protocol.TobeSendRequest{
				Request: Protocol.Request{RegistrantInfo: r.Info, Disconnect: false},
				Option:  initRequestSendOption}, initTimeout, initRetryN)
			errChan <- err //掉线时发送掉线信息
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
