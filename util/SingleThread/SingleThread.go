package SingleThread

import "sync"

type SingleThread struct {
	routine   func()
	started   bool
	startedMu *sync.Mutex
}

func New(routine func()) *SingleThread {
	return &SingleThread{routine, false, new(sync.Mutex)}
}

func (s *SingleThread) Run() {
	s.startedMu.Lock()
	if s.started { //如果已经启动过了
		return //就直接返回
	}
	s.started = true //否则就进入已启动状态
	s.startedMu.Unlock()
	defer func() {
		s.startedMu.Lock()
		s.started = false //在程序退出时重新回到未启动状态
		s.startedMu.Unlock()
	}()
	s.routine() //然后启动协程
}
