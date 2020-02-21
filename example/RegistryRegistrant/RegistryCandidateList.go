package RegistryRegistrant

import (
	"github.com/yindaheng98/go-utility/SortedSet"
	"gogistery/protocol"
	"sync"
	"time"
)

type CandidateElement struct {
	RegistryInfo protocol.RegistryInfo
}

func (e CandidateElement) GetName() string {
	return e.RegistryInfo.GetRegistryID()
}

type RegistryCandidateList struct {
	defaultTimeout time.Duration
	defaultRetryN  uint64
	set            *SortedSet.SortedSet
	setMu          *sync.RWMutex
	waitChan       chan bool
}

func NewRegistryCandidateList(initRegistry protocol.RegistryInfo, CandidateN uint64, defaultTimeout time.Duration, defaultRetryN uint64) *RegistryCandidateList {
	proto := &RegistryCandidateList{
		defaultTimeout: defaultTimeout,
		defaultRetryN:  defaultRetryN,
		set:            SortedSet.New(CandidateN),
		setMu:          new(sync.RWMutex),
		waitChan:       make(chan bool, 1),
	}
	proto.set.Update(CandidateElement{initRegistry}, 0)
	return proto
}

func (p *RegistryCandidateList) StoreCandidates(response protocol.Response) {
	p.setMu.Lock()
	defer p.setMu.Unlock()
	p.set.DeltaUpdateAll(-1)
	candidates := response.RegistryInfo.GetCandidates()
	for _, candidate := range candidates {
		el := CandidateElement{candidate}
		if _, exists := p.set.GetWeight(el); exists {
			p.set.DeltaUpdate(el, 1)
		} else {
			p.set.Update(el, 0)
		}
	}
	close(p.waitChan)               //唤醒所有等待进程
	p.waitChan = make(chan bool, 1) //然后再阻塞之
}

func (p *RegistryCandidateList) GetCandidate(excepts []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64) {
	for {
		p.setMu.Lock()
		for _, except := range excepts {
			p.set.Remove(CandidateElement{except})
		}

		if els := p.set.Sorted(1); len(els) > 0 {
			el := els[0].(CandidateElement)
			p.set.Remove(el)
			p.setMu.Unlock()
			return el.RegistryInfo, p.defaultTimeout, p.defaultRetryN
		}
		p.setMu.Unlock()
		<-p.waitChan
	}
}
