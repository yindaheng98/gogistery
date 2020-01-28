package Single

import "sync"

type Thread struct {
	started   bool
	startedMu *sync.Mutex
	Callback  *callbacks
}

func NewThread() *Thread {
	return &Thread{false, new(sync.Mutex), newCallbacks()}
}

//向此函数中输入的routine同一时刻只有一个会运行
//
//在有一个routine没有运行完成时向此函数中输入的routine会被丢弃
func (s *Thread) Run(routine func()) {
	s.startedMu.Lock()
	if s.started { //如果已经启动过了
		s.startedMu.Unlock()
		return //就直接返回
	}
	s.started = true //否则就进入已启动状态
	s.Callback.Started()
	s.startedMu.Unlock()
	defer func() {
		s.startedMu.Lock()
		s.started = false //在程序退出时重新回到未启动状态
		s.Callback.Stopped()
		s.startedMu.Unlock()
	}()
	routine() //然后启动协程
}
