package Emitter

import "sync"

//线程安全的触发器类，多线程输入事件->单线程处理事件
type Emitter struct {
	status     bool           //触发器状态（是否启动）
	statusMu   *sync.RWMutex  //触发器状态读写锁
	handlers   []func([]byte) //事件处理器列表
	handlersMu *sync.RWMutex  //事件处理器列表读写锁
	events     chan []byte    //事件队列
}

//新建触发器
func New() *Emitter {
	return &Emitter{false,
		new(sync.RWMutex),
		[]func([]byte){},
		new(sync.RWMutex),
		make(chan []byte)}
}

//添加一个事件处理函数
func (e *Emitter) AddHandler(handler func([]byte)) {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()
	e.handlers = append(e.handlers, handler)
}

//触发事件
func (e *Emitter) Emit(info []byte) {
	defer func() {
		if recover() != nil {
			e.statusMu.Lock()
			e.status = false
			e.statusMu.Unlock()
		}
	}()
	e.events <- info
}

//启动事件循环
func (e *Emitter) Start() {
	e.statusMu.Lock()
	defer e.statusMu.Unlock()
	if !e.status { //处于停止状态才启动
		e.status = true
		go e.routine() //启动事件处理循环
	}
}

//停止事件循环
func (e *Emitter) Stop() {
	e.statusMu.Lock()
	defer e.statusMu.Unlock()
	e.status = false
	close(e.events)
	e.events = make(chan []byte)
}

//goroutine循环调用事件处理函数
func (e *Emitter) routine() {
	for {
		e.statusMu.RLock()
		if !e.status { //如果要停止循环
			e.statusMu.RUnlock()
			break //那就停止循环
		}
		e.eventLoop() //事件处理循环
		e.statusMu.RUnlock()
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
