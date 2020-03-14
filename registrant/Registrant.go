package registrant

import (
	"context"
	"errors"
	"github.com/yindaheng98/gogistry/heart/requester"
	"github.com/yindaheng98/gogistry/protocol"
	"sync"
	"sync/atomic"
	"time"
)

//Registrant stands for the "registrant" in gogistry.
//A "registrant" will send requests to registry in a loop to register itself.
type Registrant struct {

	//Info contains the information of this Registrant.
	Info          protocol.RegistrantInfo
	hearts        []*requester.Heart
	connections   []protocol.RegistryInfo
	connectionsMu *sync.RWMutex

	//Every registrant has a watch dog.
	//The watch dog will scan the connection list every once in a while,
	//and if the connection list is empty for too long,
	//watch dog will forcelly exit the Registrant.Run.
	//The scan time of the watch dog can be change through WatchdogTimeDelta.
	WatchdogTimeDelta time.Duration

	candidates RegistryCandidateList //候选服务器选择协议

	//A registrant will maintain a candidate registry list.
	//
	//Information of candidate registries exist in every response from registry.
	//
	//If registry did not send back a response during a specific period of time after send a request,
	//the registrant will no longer send requests to the registry,
	//but will begin to send requests to one of the registries in candidate registry list.
	//
	//CandidateBlacklist contains those registry that the registrant should not connect to.
	CandidateBlacklist chan []protocol.RegistryInfo //不可以进行连接的候选服务器

	//Events contains 5 emitters to record running events.
	Events *events
}

//New returns the pointer to a Registrant.
func New(Info protocol.RegistrantInfo, regitryN uint64, CandidateList RegistryCandidateList,
	retryNController WaitTimeoutRetryNController, RequestProto protocol.RequestProtocol) *Registrant {
	registrant := &Registrant{
		Info:          Info,
		hearts:        make([]*requester.Heart, regitryN),
		connections:   make([]protocol.RegistryInfo, regitryN),
		connectionsMu: new(sync.RWMutex),

		WatchdogTimeDelta: 1e9,

		candidates:         CandidateList,
		CandidateBlacklist: make(chan []protocol.RegistryInfo, 1),
		Events:             newEvents(),
	}
	for i := range registrant.hearts {
		heart := requester.NewHeart(newBeater(registrant, retryNController, uint64(i)), RequestProto)
		heart.Handlers.NewConnectionHandler = func(response protocol.Response) {
			registrant.Events.NewConnection.Emit(response.RegistryInfo)
		}
		heart.Handlers.UpdateConnectionHandler = func(response protocol.Response) {
			registrant.Events.UpdateConnection.Emit(response.RegistryInfo)
		}
		heart.Handlers.RetryHandler = func(request protocol.TobeSendRequest, err error) {
			registrant.Events.Retry.Emit(request, err)
		}
		heart.Handlers.DisconnectionHandler = func(response protocol.Response, err error) {
			registrant.Events.Disconnection.Emit(response.RegistryInfo, err)
		}
		registrant.hearts[i] = heart
	}
	registrant.CandidateBlacklist <- []protocol.RegistryInfo{}
	return registrant
}

func (r *Registrant) register(ctx context.Context, response protocol.Response, i uint64) bool {
	okChan := make(chan bool, 1)
	go func() {
		r.candidates.StoreCandidates(ctx, response.RegistryInfo.GetCandidates()) //先读取候选
		okChan <- true
	}()
	select {
	case <-okChan:
	case <-ctx.Done():
		return false
	}
	r.connectionsMu.Lock()
	defer r.connectionsMu.Unlock()
	if response.IsReject() { //如果拒绝连接
		r.connections[i] = nil //就删除
		return false           //然后断开连接
	} else {
		r.connections[i] = response.RegistryInfo //否则修改记录
		return true
	}
}

//Run will start the request sending and wake up the watch dog.
//The information of registries that the registrant should connect was given by candidate registry list,
//so make sure when registry is needed, the candidate registry list is not empty.
func (r *Registrant) Run(ctx context.Context) {
	connectionsChan := make(chan []protocol.RegistryInfo, 1)
	connections := make([]protocol.RegistryInfo, len(r.connections))
	connectionsChan <- connections
	GetExcept := func() []protocol.RegistryInfo { //锁定并读取除外项
		connections = <-connectionsChan
		var except []protocol.RegistryInfo
		for _, c := range connections { //去除空项
			if c != nil {
				except = append(except, c)
			}
		}
		blacklist := <-r.CandidateBlacklist //取不可连接列表
		except = append(except, blacklist...)
		r.CandidateBlacklist <- blacklist
		return except
	}
	PutExcept := func(info protocol.RegistryInfo, i uint64) { //放入并释放除外项
		connections[i] = info
		connectionsChan <- connections
	}
	ResetExcept := func(i uint64) { //锁定清除释放一条龙服务
		connections := <-connectionsChan
		connections[i] = nil
		connectionsChan <- connections
	}

	ctx, cancel := context.WithCancel(ctx)
	wg := new(sync.WaitGroup)
	beatingN := int64(0)
	for i := range r.connections {
		r.connections[i] = nil
		wg.Add(1)
		go func(i uint64) {
			r.heartRoutine(ctx, r.hearts[i], &beatingN,
				GetExcept,
				func(info protocol.RegistryInfo) { PutExcept(info, i) },
				func() { ResetExcept(i) })
			wg.Done()
		}(uint64(i))
	}
	go r.watchDog(ctx, &beatingN, cancel)
	wg.Wait()
	cancel()
}

func (r *Registrant) heartRoutine(ctx context.Context, h *requester.Heart, beatingN *int64,
	GetExcept func() []protocol.RegistryInfo, PutExcept func(protocol.RegistryInfo), ResetExcept func()) {
	for {
		var initRegistryInfo protocol.RegistryInfo
		var initTimeout time.Duration
		var initRetryN uint64
		done := make(chan bool, 1)
		go func() {
			initRegistryInfo, initTimeout, initRetryN = r.candidates.GetCandidate(ctx, GetExcept()) //获取新连接
			done <- true
		}()
		select {
		case <-done: //完成则进行下一步
			PutExcept(initRegistryInfo)
		case <-ctx.Done(): //突然要停止
			PutExcept(nil) //删除除外表
			return         //然后退出
		}

		atomic.AddInt64(beatingN, 1)
		err := h.RunBeating(ctx, protocol.TobeSendRequest{
			Request: protocol.Request{RegistrantInfo: r.Info, Disconnect: false},
			Option:  initRegistryInfo.GetRequestSendOption()}, initTimeout, initRetryN)
		ResetExcept() //清空除外表
		if err != nil {
			r.Events.Error.Emit(err)
		}
		atomic.AddInt64(beatingN, -1)
	}
}

func (r *Registrant) watchDog(ctx context.Context, beatingN *int64, cancel func()) {
	lastBite := false //记录上一次watch是否达标
	for {
		select {
		case <-ctx.Done(): //如果要停机
			cancel()
			return //就直接退出
		case <-time.After(r.WatchdogTimeDelta): //每隔一段时间watch一次
			if atomic.LoadInt64(beatingN) <= 0 { //如果本次watch没达标
				if lastBite { //并且上一次watch也没达标
					r.Events.Error.Emit(errors.New("all the sending goroutine fall asleep"))
					cancel() //就直接停掉
				}
				lastBite = true //将达标标记置为true
			} else { //此次watch达标
				lastBite = false //达标标记置为false
			}
		}
	}
}

//Registrant maintains a connection list,
//which can be accessed through GetConnections.
//This connection list contains the information (`protocol.RegistryInfo`) of the registries connecting by the `Registrant`.
func (r *Registrant) GetConnections() []protocol.RegistryInfo {
	connections := make([]protocol.RegistryInfo, 0)
	r.connectionsMu.RLock()
	defer r.connectionsMu.RUnlock()
	for _, c := range r.connections {
		if c != nil {
			connections = append(connections, c)
		}
	}
	return connections
}

//AddCandidates can add the information of a registries to candidate registry list.
func (r *Registrant) AddCandidates(ctx context.Context, candidates []protocol.RegistryInfo) {
	r.candidates.StoreCandidates(ctx, candidates)
}
