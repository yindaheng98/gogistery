package CandidateList

import (
	"github.com/yindaheng98/go-utility/SortedSet"
	"gogistery/protocol"
	"sync"
	"time"
)

type SimpleCandidateList struct {
	defaultTimeout time.Duration
	defaultRetryN  uint64
	set            *SortedSet.SortedSet
	setMu          *sync.RWMutex
	waitChan       chan bool
}

func NewSimpleCandidateList(initRegistry protocol.RegistryInfo, CandidateN uint64, defaultTimeout time.Duration, defaultRetryN uint64) *SimpleCandidateList {
	proto := &SimpleCandidateList{
		defaultTimeout: defaultTimeout,
		defaultRetryN:  defaultRetryN,
		set:            SortedSet.New(CandidateN),
		setMu:          new(sync.RWMutex),
		waitChan:       make(chan bool, 1),
	}
	proto.set.Update(element{initRegistry}, 0)
	return proto
}

func (p *SimpleCandidateList) StoreCandidates(response protocol.Response) {
	p.setMu.Lock()
	defer p.setMu.Unlock()
	p.set.DeltaUpdateAll(-1)
	candidates := response.RegistryInfo.GetCandidates()
	for _, candidate := range candidates {
		el := element{candidate}
		if _, exists := p.set.GetWeight(el); exists {
			p.set.DeltaUpdate(el, 1)
		} else {
			p.set.Update(el, 0)
		}
	}
	close(p.waitChan)               //唤醒所有等待进程
	p.waitChan = make(chan bool, 1) //然后再阻塞之
}

func (p *SimpleCandidateList) GetCandidate(excepts []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64) {
	for {
		p.setMu.Lock()
		for _, except := range excepts {
			p.set.Remove(element{except})
		}

		if els := p.set.Sorted(1); len(els) > 0 {
			el := els[0].(element)
			p.set.Remove(el)
			p.setMu.Unlock()
			return el.RegistryInfo, p.defaultTimeout, p.defaultRetryN
		}
		p.setMu.Unlock()
		<-p.waitChan
	}
}
