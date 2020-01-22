package SingleThread

import "sync"

type SingleThread struct {
	started   bool
	startedMu *sync.Mutex
}

func New() *SingleThread {
	return &SingleThread{false, new(sync.Mutex)}
}

//向此函数中输入的routine同一时刻只有一个会运行
//
//在有一个routine没有运行完成时向此函数中输入的routine会被丢弃
func (s *SingleThread) Run(routine func()) {
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
	routine() //然后启动协程
}

func (s *SingleThread) IsRunning() bool {
	return s.started
}
