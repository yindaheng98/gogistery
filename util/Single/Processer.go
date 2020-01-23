package Single

import "sync"

//调用Start()之后，process函数将被一次次循环调用直到调用了Stop()
//
//可以保证同时只会有一个process函数被调用
type Processor struct {
	thread    *Thread
	process   func()
	started   bool
	startedMu *sync.Mutex
}

func NewProcessor(process func()) *Processor {
	return &Processor{NewThread(), process, false, new(sync.Mutex)}
}

func (p *Processor) Start() {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	if !p.started {
		p.started = true
		go p.thread.Run(p.routine)
	}
}

func (p *Processor) Stop() {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	p.started = false
}

func (p *Processor) routine() {
	for p.started {
		p.process()
	}
}

func (p *Processor) IsRunning() bool {
	return p.thread.IsRunning()
}
