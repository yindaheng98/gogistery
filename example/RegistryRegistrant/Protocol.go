package RegistryRegistrant

import (
	"github.com/yindaheng98/go-utility/SortedSet"
	"gogistery/Protocol"
	"math"
	"sync"
	"time"
)

type RegistrantControlProtocol struct {
	minT   time.Duration //最小Timeout
	maxT   time.Duration //最大Timeout
	cT     float64       //从最小到最大的增长系数
	tMap   map[string]time.Duration
	retryN uint64
}

func NewRegistrantControlProtocol(minT time.Duration, maxT time.Duration, cT float64, retryN uint64) *RegistrantControlProtocol {
	return &RegistrantControlProtocol{minT, maxT, cT,
		make(map[string]time.Duration), retryN}
}

func (p RegistrantControlProtocol) TimeoutRetryNForNew(request Protocol.Request) (time.Duration, uint64) {
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = p.minT
	return p.minT, p.retryN

}
func (p RegistrantControlProtocol) TimeoutRetryNForUpdate(request Protocol.Request) (time.Duration, uint64) {
	t := p.tMap[request.RegistrantInfo.GetRegistrantID()]
	t += time.Duration(math.Floor(float64(p.maxT-t) / p.cT))
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = t
	return t, p.retryN

}

type CandidateElement struct {
	RegistryInfo Protocol.RegistryInfo
}

func (e CandidateElement) GetName() string {
	return e.RegistryInfo.GetRegistryID()
}

type CandidateRegistryProtocol struct {
	defaultTimeout time.Duration
	defaultRetryN  uint64
	set            *SortedSet.SortedSet
	setMu          *sync.RWMutex
	waitChan       chan bool
}

func NewCandidateRegistryProtocol(initRegistry Protocol.RegistryInfo, CandidateN uint64, defaultTimeout time.Duration, defaultRetryN uint64) *CandidateRegistryProtocol {
	proto := CandidateRegistryProtocol{
		defaultTimeout: defaultTimeout,
		defaultRetryN:  defaultRetryN,
		set:            SortedSet.New(CandidateN),
		setMu:          new(sync.RWMutex),
		waitChan:       make(chan bool, 1),
	}
	proto.set.Update(CandidateElement{initRegistry}, 0)
	return &proto
}

func (p *CandidateRegistryProtocol) StoreCandidates(response Protocol.Response) {
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

func (p *CandidateRegistryProtocol) GetCandidate(excepts []Protocol.RegistryInfo) (Protocol.RegistryInfo, time.Duration, uint64) {
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
