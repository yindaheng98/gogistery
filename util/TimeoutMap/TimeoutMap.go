package TimeoutMap

import (
	"gogistery/util/TimeoutMap/TimeoutValue"
	"sync"
	"time"
)

type TimeoutMap struct {
	elements map[string]*TimeoutValue.TimeoutValue //存储数据
	mu       *sync.RWMutex                         //读写锁
}

//输入超时时间和删除缓存的数量新建发送器列表
func New() *TimeoutMap {
	return &TimeoutMap{make(map[string]*TimeoutValue.TimeoutValue), new(sync.RWMutex)}
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
func (m *TimeoutMap) UpdateInfo(el Element, timeout time.Duration) {
	id := el.GetID() //先获取发送器信息中的id
	m.mu.RLock()
	value, exists := m.elements[id] //查询此id是否存在
	m.mu.RUnlock()
	if exists { //如果存在
		value.Update(el) //则更新
	} else {
		m.mu.Lock()
		defer m.mu.Unlock()
		value := TimeoutValue.New(el, timeout, func() {
			m.delete(id) //超时则删除
		}) //否则新建
		m.elements[id] = value
		value.Start()
	}
}

func (m *TimeoutMap) delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.elements, id)
}

func (m *TimeoutMap) getElement(id string) (Element, bool) {
	var el Element = nil
	value, ok := m.elements[id] //先查找
	if ok {                     //如果找得到
		el = value.GetElement().(Element) //则返回结果
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
	value, ok := m.elements[id] //先查找
	if ok {                     //如果找得到
		value.Stop() //则使其停止
		m.mu.RUnlock()
		m.delete(id) //并删除
	} else {
		m.mu.RUnlock()
	}
}
