package registrant

import (
	"errors"
	"gogistery/protocol"
	"sync"
	"sync/atomic"
	"time"
)

type Registrant struct {
	Info              protocol.RegistrantInfo //注册器的信息
	hearts            []*heart                //heart列表，每一个heart负责向一个注册线程
	stopChan          chan bool               //传递停止信息
	WatchdogTimeDelta time.Duration

	candidates RegistryCandidateList //候选服务器选择协议
	Events     *events
}

func New(Info protocol.RegistrantInfo, regitryN uint, CandidateList RegistryCandidateList,
	retryNController RetryNController, RequestProto protocol.RequestProtocol) *Registrant {
	registrant := &Registrant{
		Info:              Info,
		hearts:            make([]*heart, regitryN),
		stopChan:          make(chan bool, 1),
		WatchdogTimeDelta: 1e9,

		candidates: CandidateList,
		Events:     newEvents(),
	}
	close(registrant.stopChan)
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
	r.stopChan = make(chan bool, 1)
	wg := new(sync.WaitGroup)
	connChan := make(chan []protocol.RegistryInfo, 1)
	connChan <- make([]protocol.RegistryInfo, len(r.hearts))
	beatingN := int64(0)
	for i, h := range r.hearts {
		wg.Add(1)
		go func(i int, h *heart) {
			r.heartRoutine(h, i, connChan, &beatingN)
			wg.Done()
		}(i, h)
	}
	go r.watchDog(&beatingN)
	wg.Wait()
}

func (r *Registrant) watchDog(beatingN *int64) {
	lastBite := false //记录上一次watch是否达标
	for {
		select {
		case <-r.stopChan: //如果要停机
			return //就直接退出
		case <-time.After(r.WatchdogTimeDelta): //每隔一段时间watch一次
			if atomic.LoadInt64(beatingN) <= 0 { //如果本次watch没达标
				if lastBite { //并且上一次watch也没达标
					r.Events.Error.Emit(errors.New("all the sending goroutine fall asleep"))
					r.Stop() //就直接停掉
				}
				lastBite = true //将达标标记置为true
			} else { //此次watch达标
				lastBite = false //达标标记置为false
			}
		}
	}
}

func (r *Registrant) heartRoutine(h *heart, i int, connChan chan []protocol.RegistryInfo, beatingN *int64) {
	for {
		conn := <-connChan                 //从队列中取出已有连接列表
		var except []protocol.RegistryInfo //去除空项
		for _, c := range conn {
			if c != nil {
				except = append(except, c)
			}
		}

		var initRegistryInfo protocol.RegistryInfo
		var initTimeout time.Duration
		var initRetryN uint64
		done := make(chan bool, 1)
		go func() {
			initRegistryInfo, initTimeout, initRetryN = r.candidates.GetCandidate(except) //获取新连接
			done <- true
		}()
		select {
		case <-done: //完成则进行下一步
			conn[i] = initRegistryInfo //将新连接加入已有连接列表
			connChan <- conn           //已有连接列表放回队列
		case <-r.stopChan: //突然要停止
			connChan <- conn //已有连接列表放回队列
			return           //直接退出
		}

		atomic.AddInt64(beatingN, 1)
		err := h.RunBeating(protocol.TobeSendRequest{
			Request: protocol.Request{RegistrantInfo: r.Info, Disconnect: false},
			Option:  initRegistryInfo.GetRequestSendOption()}, initTimeout, initRetryN)
		if err != nil {
			r.Events.Error.Emit(err)
		}
		atomic.AddInt64(beatingN, -1)

		conn = <-connChan //从队列中取出已有连接列表
		conn[i] = nil     //将对应位置空
		connChan <- conn  //已有连接列表放回队列
	}
}

func (r *Registrant) Stop() {
	func() {
		defer func() { recover() }()
		close(r.stopChan)
	}()
	for _, h := range r.hearts {
		h.Stop()
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
