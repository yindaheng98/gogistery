package TimeoutMap

import (
	"fmt"
	"sync"
	"time"
)

type timeoutValue struct {
	element     Element
	mu          *sync.RWMutex
	updatedTime time.Time
	timeouted   bool
}

func newValue(element Element) *timeoutValue {
	return &timeoutValue{element, new(sync.RWMutex), time.Now(), false}
}

func (v *timeoutValue) Update(el Element) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if el != nil {
		v.element = el
	}
	v.updatedTime = time.Now()
	v.timeouted = false
}

func (v *timeoutValue) IsTimeout(timeout time.Duration) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	now := time.Now()
	fmt.Printf("%d\n", now.Sub(v.updatedTime))
	return v.timeouted || timeout < now.Sub(v.updatedTime)
}

func (v *timeoutValue) MakeTimeout() {
	v.timeouted = true
}
