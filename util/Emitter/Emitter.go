package Emitter

import (
	"sync"
	"sync/atomic"
)

//线程安全的触发器类，多线程输入事件->单线程处理事件
type Emitter struct {
	started    uint32              //触发器状态（是否启动）
	handlers   []func(interface{}) //事件处理器列表
	handlersMu *sync.RWMutex       //事件处理器列表读写锁
	events     chan interface{}    //事件队列
}

//新建触发器
func New() *Emitter {
	return &Emitter{0,
		[]func(interface{}){},
		new(sync.RWMutex),
		make(chan interface{})}
}

//添加一个事件处理函数
func (e *Emitter) AddHandler(handler func(interface{})) {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()
	e.handlers = append(e.handlers, handler)
}

//触发事件
func (e *Emitter) Emit(info interface{}) {
	defer func() {
		if recover() != nil {
			e.Stop()
		}
	}()
	e.events <- info
}

//启动事件循环
func (e *Emitter) Start() {
	if atomic.CompareAndSwapUint32(&e.started, 0, 1) { //处于停止状态才启动
		go e.routine() //启动事件处理循环
	}
}

//停止事件循环
func (e *Emitter) Stop() {
	if atomic.CompareAndSwapUint32(&e.started, 1, 0) { //处于启动状态才进行停止操作
		close(e.events)
		e.events = make(chan interface{})
	}
}

//goroutine循环调用事件处理函数
func (e *Emitter) routine() {
	for {
		if atomic.CompareAndSwapUint32(&e.started, 0, 0) { //如果要停止循环
			break //那就停止循环
		}
		e.eventLoop() //事件处理循环
	}
}

//事件处理循环：出队列处理事件
func (e *Emitter) eventLoop() {
	info, ok := <-e.events
	if !ok {
		return
	}
	e.handlersMu.RLock()
	defer e.handlersMu.RUnlock()
	for _, handler := range e.handlers {
		handler(info)
	}
}
