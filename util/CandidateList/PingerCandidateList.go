package CandidateList

import (
	"github.com/yindaheng98/go-utility/SortedSet"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type PINGer interface {
	PING(info protocol.RegistryInfo) bool
}

type pingTimer struct {
	pinger         PINGer
	maxPingTimeout time.Duration
}

//返回一次ping操作的时间
func (p *pingTimer) PINGTime(el element) (time.Duration, bool) {
	pingChan := make(chan bool, 1)
	t := time.Now() //记录ping操作开始的时间
	go func() { pingChan <- p.pinger.PING(el.RegistryInfo) }()
	var ok bool
	select {
	case ok = <-pingChan:
	case <-time.After(p.maxPingTimeout):
		ok = false //超时则失败
	}
	return time.Now().Sub(t), ok //返回ping操作时长和是否成功
}

//使用PING操作更新候选列表顺序的CandidateList
type PingerCandidateList struct {
	SimpleCandidateList
	timer pingTimer

	size        uint64                  //每次进行多少个候选ping操作
	pingingList chan map[string]element //正在被ping的注册中心列表
}

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
	list.ping(element{initRegistry})
	return list
}

//进行一次ping操作并按照结果更新集合
func (list *PingerCandidateList) ping(el element) {
	pingingList := <-list.pingingList
	if _, exists := pingingList[el.GetName()]; exists { //如果已经有其他进程在ping了
		list.pingingList <- pingingList
		return //就直接退出
	}
	pingingList[el.GetName()] = el //开始ping操作前先入pinging列表
	list.pingingList <- pingingList

	set := <-list.set                      //取出set
	timeout, ok := list.timer.PINGTime(el) //进行ping操作
	if ok {                                //如果成功
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

func (list *PingerCandidateList) StoreCandidates(candidates []protocol.RegistryInfo) {
	set := <-list.set
	defer func() { list.set <- set }()
	set.DeltaUpdateAll(-1)
	for _, info := range candidates { //遍历候选元素
		el := element{info}
		if w, ok := set.GetWeight(el); ok { //如果已存在
			set.Update(el, w) //则只更新元素内容
		} else { //否则进行ping和权重更新操作
			go list.ping(el)
		}
	}
}
