package registrant

import (
	"gogistery/protocol"
	"sync/atomic"
)

type Registrant struct {
	Info     protocol.RegistrantInfo //注册器的信息
	hearts   []*heart                //heart列表，每一个heart负责向一个注册线程
	runningN int64                   //正在运行的线程数
	toStop   bool

	candidates RegistryCandidateList //候选服务器选择协议
	Events     *events
}

func New(Info protocol.RegistrantInfo, regitryN uint, CandidateList RegistryCandidateList,
	retryNController RetryNController, RequestProto protocol.RequestProtocol) *Registrant {
	registrant := &Registrant{
		Info:     Info,
		hearts:   make([]*heart, regitryN),
		runningN: 0, //正在运行的线程数
		toStop:   false,

		candidates: CandidateList,
		Events:     newEvents(),
	}
	for i := uint(0); i < regitryN; i++ {
		registrant.hearts[i] = newHeart(registrant, retryNController, RequestProto)
	}
	return registrant
}

//For the struct heart
func (r *Registrant) register(response protocol.Response) protocol.Request {
	r.candidates.StoreCandidates(response)
	return protocol.Request{
		RegistrantInfo: r.Info,
		Disconnect:     response.IsReject(),
	}
}

func (r *Registrant) Run() {
	r.toStop = false
	stoppedChan := make(chan bool, 1)
	connChan := make(chan []protocol.RegistryInfo, 1)
	connChan <- make([]protocol.RegistryInfo, len(r.hearts))
	for i, h := range r.hearts {
		go func(i int, h *heart) {
			atomic.AddInt64(&r.runningN, 1)
			r.heartRoutine(h, i, connChan)
			atomic.AddInt64(&r.runningN, -1)
			if atomic.LoadInt64(&r.runningN) <= 0 {
				stoppedChan <- true
				close(stoppedChan)
			}
		}(i, h)
	}
	<-stoppedChan
}

func (r *Registrant) Stop() {
	r.toStop = true
	for _, h := range r.hearts {
		h.Stop()
	}
}

func (r *Registrant) heartRoutine(h *heart, i int, connChan chan []protocol.RegistryInfo) {
	for !r.toStop {
		conn := <-connChan                 //从队列中取出已有连接列表
		var except []protocol.RegistryInfo //去除空项
		for _, c := range conn {
			if c != nil {
				except = append(except, c)
			}
		}
		initRegistryInfo, initTimeout, initRetryN := r.candidates.GetCandidate(except) //以此获取新连接
		conn[i] = initRegistryInfo                                                     //将新连接加入已有连接列表
		connChan <- conn                                                               //已有连接列表放回队列
		err := h.RunBeating(protocol.TobeSendRequest{
			Request: protocol.Request{RegistrantInfo: r.Info, Disconnect: false},
			Option:  initRegistryInfo.GetRequestSendOption()}, initTimeout, initRetryN)
		if err != nil {
			r.Events.Error.Emit(err)
		}
		conn = <-connChan //从队列中取出已有连接列表
		conn[i] = nil     //将对应位置空
		connChan <- conn  //已有连接列表放回队列
	}
}

func (r *Registrant) GetConnections() []protocol.RegistryInfo {
	res := make([]protocol.RegistryInfo, 0)
	for _, h := range r.hearts {
		if h.RegistryInfo != nil {
			res = append(res, h.RegistryInfo)
		}
	}
	return res
}
