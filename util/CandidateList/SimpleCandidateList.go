package CandidateList

import (
	"github.com/yindaheng98/go-utility/SortedSet"
	"gogistery/protocol"
	"time"
)

type SimpleCandidateList struct {
	initTimeout time.Duration
	initRetryN  uint64
	set         chan *SortedSet.SortedSet
	waitGroup   chan bool
}

func NewSimpleCandidateList(size uint64, initRegistry protocol.RegistryInfo, initTimeout time.Duration, initRetryN uint64) *SimpleCandidateList {
	list := &SimpleCandidateList{
		initTimeout: initTimeout,
		initRetryN:  initRetryN,
		set:         make(chan *SortedSet.SortedSet, 1),
		waitGroup:   make(chan bool),
	}
	set := SortedSet.New(size)
	set.Update(element{initRegistry}, 0)
	list.set <- set
	return list
}

func (list *SimpleCandidateList) StoreCandidates(response protocol.Response) {
	set := <-list.set
	defer func() {
		list.set <- set                  //完成后放回队列
		close(list.waitGroup)            //唤醒所有等待的获取队列进程
		list.waitGroup = make(chan bool) //然后再阻塞之
	}()
	set.DeltaUpdateAll(-1) //先让所有元素优先级下降1
	candidates := response.RegistryInfo.GetCandidates()
	for _, candidate := range candidates {
		el := element{candidate}
		if _, exists := set.GetWeight(el); exists {
			set.DeltaUpdate(el, 1) //然后让响应中给出的元素优先级上升1
		} else {
			set.Update(el, 0) //或者加入新元素
		}
	}
}

func (list *SimpleCandidateList) GetCandidate(excepts []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64) {
	for {
		set := <-list.set
		for _, except := range excepts {
			set.Remove(element{except}) //先移除所有除外元素
		}
		if els := set.Sorted(1); len(els) > 0 { //然后看是否还有剩余元素
			el := els[0].(element) //有就作为返回结果
			set.Remove(el)         //然后删除返回结果
			list.set <- set        //并将集合放回队列
			return el.RegistryInfo, list.initTimeout, list.initRetryN
		}
		list.set <- set  //并将集合放回队列
		<-list.waitGroup //等待更新
	}
}
