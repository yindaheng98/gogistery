package TimeoutValue

import (
	"gogistery/util/Single"
	"sync"
	"time"
)

type TimeoutValue struct {
	timeout        time.Duration //固定超时时间
	runner         *Single.Processor
	timeoutHandler func() //超时时的回调函数
	//以上是在TimeoutValue创建即固定的变量

	element    Element     //值内容
	isUpdated  bool        //更新标记，每次更新都刷新一次
	updateTime time.Time   //上次刷新更新标记的时间
	updateMu   *sync.Mutex //更新锁，同一时刻只能有一个线程在更新
	//以上是在一次update操作中会被修改的变量

	toStop bool        //停止标记——是否要停止
	stopMu *sync.Mutex //停止锁，同意时刻只能有一个线程在修改停止位标记
	//以上是在一次停止操作中会被修改的变量
}

func New(element Element, timeout time.Duration, timeoutHandler func()) *TimeoutValue {
	return &TimeoutValue{timeout, Single.NewProcessor(), timeoutHandler,
		element, true, time.Now(), new(sync.Mutex),
		false, new(sync.Mutex)}
}

func (v *TimeoutValue) GetElement() Element {
	return v.element
}

//启动检查线程
func (v *TimeoutValue) Start() {
	v.stopMu.Lock()
	defer v.stopMu.Unlock()
	v.toStop = false
	v.runner.Callback.Started = func() {
		v.element.NewAddedHandler()
	}
	v.runner.Start(func() {
		v.updateMu.Lock()
		v.isUpdated = false //重置更新标记
		v.updateMu.Unlock()
		time.Sleep(v.timeout - time.Now().Sub(v.updateTime)) //等待一定时间
		v.stopMu.Lock()
		defer v.stopMu.Unlock()
		if (!v.isUpdated) && (!v.toStop) { //如果一定时间后更新标记还没有更新
			v.runner.Stop()    //就停止检查线程
			v.timeoutHandler() //并调用超时回调
			v.element.TimeoutHandler()
		}
	})
}

func (v *TimeoutValue) Update(el Element) {
	v.updateMu.Lock()
	defer v.updateMu.Unlock()
	v.updateTime = time.Now() //先更新刷新时间
	v.isUpdated = true        //再更新刷新标记
	if el != nil {
		v.element = el //再更新值内容
	}
	v.Start() //确保检查进程的持续运行
}

func (v *TimeoutValue) Stop() {
	v.stopMu.Lock()
	defer v.stopMu.Unlock()
	v.toStop = true
	v.runner.Stop()
}
