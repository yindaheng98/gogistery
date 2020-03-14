package CandidateList

import (
	"context"
	"github.com/yindaheng98/go-utility/SortedSet"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//PINGer defines a customized PING method for the measurement of network latency in PingerCandidateList.
type PINGer interface {
	PING(ctx context.Context, info protocol.RegistryInfo) bool
}

type pingTimer struct {
	pinger         PINGer
	maxPingTimeout time.Duration
}

//返回一次ping操作的时间
func (p *pingTimer) PINGTime(ctx context.Context, el element) (time.Duration, bool) {
	pingChan := make(chan bool, 1)
	t := time.Now() //记录ping操作开始的时间
	pingCtx, cancel := context.WithTimeout(ctx, p.maxPingTimeout)
	defer cancel()
	go func() { pingChan <- p.pinger.PING(pingCtx, el.RegistryInfo) }()
	var ok bool
	select {
	case ok = <-pingChan:
		ok = true
	case <-pingCtx.Done(): //超时则失败
		ok = false
	}
	return time.Now().Sub(t), ok //返回ping操作时长和是否成功
}

//PingerCandidateList is a implementation of RegistryCandidateList.
//This implementation sort registries by a their network latency.
type PingerCandidateList struct {
	SimpleCandidateList
	timer pingTimer

	size        uint64                  //每次进行多少个候选ping操作
	pingingList chan map[string]element //正在被ping的注册中心列表
}

//NewPingerCandidateList returns the pointer to a PingerCandidateList.
func NewPingerCandidateList(size uint64, pinger PINGer, maxPingTimeout time.Duration,
	initRegistry protocol.RegistryInfo, initTimeout time.Duration, initRetryN uint64) *PingerCandidateList {
	list := &PingerCandidateList{
		SimpleCandidateList: SimpleCandidateList{
			initTimeout: initTimeout,
			initRetryN:  initRetryN,
			set:         make(chan *SortedSet.SortedSet, 1),
			waitGroup:   make(chan bool),
		},
		timer: pingTimer{
			pinger:         pinger,
			maxPingTimeout: maxPingTimeout,
		},

		size:        size,
		pingingList: make(chan map[string]element, 1),
	}

	list.set <- SortedSet.New(size)
	list.pingingList <- make(map[string]element, size)
	list.ping(context.Background(), element{initRegistry})
	return list
}

//进行一次ping操作并按照结果更新集合
func (list *PingerCandidateList) ping(ctx context.Context, el element) {
	pingingList := <-list.pingingList
	if _, exists := pingingList[el.GetName()]; exists { //如果已经有其他进程在ping了
		list.pingingList <- pingingList
		return //就直接退出
	}
	pingingList[el.GetName()] = el //开始ping操作前先入pinging列表
	list.pingingList <- pingingList

	set := <-list.set //取出set
	var timeout time.Duration
	okChan := make(chan bool, 1)
	go func() {
		var ok bool
		timeout, ok = list.timer.PINGTime(ctx, el) //进行ping操作
		okChan <- ok
	}()
	ok := false
	select {
	case ok = <-okChan:
	case <-ctx.Done():
	}
	if ok { //如果成功
		set.Update(el, float64(timeout))    //则更新
		close(list.waitGroup)               //唤醒所有等待进程
		list.waitGroup = make(chan bool, 1) //然后再阻塞之
	} else {
		set.Remove(el)
	}
	list.set <- set //完成后放回

	pingingList = <-list.pingingList
	delete(pingingList, el.GetName()) //操作结束后再删pinging列表
	list.pingingList <- pingingList
}

//StoreCandidates adds or updates registries in the list candidates and change the registries' weight
//according to their network latency.
func (list *PingerCandidateList) StoreCandidates(ctx context.Context, candidates []protocol.RegistryInfo) {
	set := <-list.set
	defer func() { list.set <- set }()
	set.DeltaUpdateAll(-1)
	for _, info := range candidates { //遍历候选元素
		el := element{info}
		if w, ok := set.GetWeight(el); ok { //如果已存在
			set.Update(el, w) //则只更新元素内容
		} else { //否则进行ping和权重更新操作
			go list.ping(ctx, el)
		}
	}
}
