package Single

import "sync"

//调用Start()之后，process函数将被一次次循环调用直到调用了Stop()
//
//可以保证同时只会有一个process函数被调用
type Processor struct {
	thread    *Thread
	started   bool
	startedMu *sync.Mutex
	Callback  *callbacks
}

func NewProcessor() *Processor {
	p := &Processor{NewThread(), false, new(sync.Mutex), newCallbacks()}
	p.thread.Callback = p.Callback
	return p
}

func (p *Processor) Start(process func()) {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	p.thread.Callback = p.Callback
	if !p.started {
		p.started = true
		go p.thread.Run(func() {
			for p.started {
				process()
			}
		})
	}
}

func (p *Processor) Stop() {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	p.started = false
}
