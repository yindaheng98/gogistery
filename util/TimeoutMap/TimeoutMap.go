package TimeoutMap

import (
	"sync"
	"time"
)

type TimeoutMap struct {
	elements map[string]*timeoutValue //存储数据
	mu       *sync.RWMutex            //读写锁
	timeout  time.Duration            //指定超时时间

	delBufferN uint64   //删除缓存的数量
	deli       uint64   //删到第几个了
	delBuffer  []string //要删的id列表
}

//输入超时时间和删除缓存的数量新建发送器列表
func New(timeout time.Duration, delBufferN uint64) *TimeoutMap {
	return &TimeoutMap{make(map[string]*timeoutValue), new(sync.RWMutex), timeout,
		delBufferN, 0, make([]string, delBufferN)}
}

func (m *TimeoutMap) GetTimeout() time.Duration {
	return m.timeout
}

//通过id进行更新，仅更新时间
func (m *TimeoutMap) UpdateID(id string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.elements[id] //查询此id是否存在
	if exists {                     //如果存在
		value.Update(nil) //则更新
	}
}

//通过一个Element进行更新，更新存储的信息
func (m *TimeoutMap) UpdateInfo(el Element) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id := el.GetID()                //先获取发送器信息中的id
	value, exists := m.elements[id] //查询此id是否存在
	if exists {                     //如果存在
		value.Update(el) //则更新
	} else {
		m.elements[id] = newValue(el) //否则新建
	}
}

func (m *TimeoutMap) delete(id string) {
	value, ok := m.elements[id] //先查找
	if ok {                     //如果有
		value.MakeTimeout()                         //则使其超时
		m.delBuffer[m.deli] = value.element.GetID() //然后放入删除队列
		m.deli = (m.deli + 1) % m.delBufferN        //更新删除数量
		if m.deli <= 0 {                            //删除缓存累积到指定数量
			go m.delRoutine(m.delBuffer)               //就启动删除线程
			m.delBuffer = make([]string, m.delBufferN) //然后刷新删除缓存
		}
	}
}

//删除线程
func (m *TimeoutMap) delRoutine(ids []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		if element, ok := m.elements[id]; ok && element.IsTimeout(m.timeout) {
			delete(m.elements, id)
		}
	}
}

func (m *TimeoutMap) getElement(id string) (Element, bool) {
	var el Element = nil
	value, ok := m.elements[id] //先查找
	if ok {                     //如果找得到
		if value.IsTimeout(m.timeout) { //那就看是否超时
			m.delete(id) //超时则删
			ok = false
		} else {
			el = value.element //不超时则返回结果
		}
	}
	return el, ok
}

//获取某个id对应的信息
func (m *TimeoutMap) GetElement(id string) (Element, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getElement(id)
}

//获取所有的信息
func (m *TimeoutMap) GetAll() []Element {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var els []Element
	for id := range m.elements {
		el, ok := m.getElement(id)
		if ok {
			els = append(els, el)
		}
	}
	return els
}

func (m *TimeoutMap) Delete(id string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.delete(id)
}
