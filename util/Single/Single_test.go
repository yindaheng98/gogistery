package Single

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleThread(t *testing.T) {
	i := int32(0)
	st := NewThread()
	r := func() {
		atomic.AddInt32(&i, 1)
		t.Log(fmt.Sprintf("I'm No.%d", i))
		atomic.AddInt32(&i, -1)
		time.Sleep(1e8)
	}
	rr := func() {
		if st.IsRunning() {
			t.Log("I'm running")
		} else {
			t.Log("I'm not running")
		}
	}
	go st.Run(r)
	go rr()
	go rr()
	go rr()
	go st.Run(r)
	go rr()
	go rr()
	go rr()
	go st.Run(r)
	go rr()
	go rr()
	go rr()
	time.Sleep(1e9)
	go rr()
	go rr()
	go rr()
	time.Sleep(1e9)
}

func TestSingleProcessor(t *testing.T) {
	i := int32(0)
	p := NewProcessor(func() {
		atomic.AddInt32(&i, 1)
		t.Log(fmt.Sprintf("I'm No.%d", i))
		atomic.AddInt32(&i, -1)
	})
	pp := func() {
		if p.IsRunning() {
			t.Log("I'm running")
		} else {
			t.Log("I'm not running")
		}
	}
	go p.Start()
	go pp()
	go p.Stop()
	go pp()
	go p.Start()
	go pp()
	go p.Stop()
	go pp()
	go p.Start()
	go pp()
	go p.Stop()
	go pp()
	go p.Stop()
	go p.Stop()
	go p.Stop()
	go p.Stop()
	go pp()
	go pp()
	time.Sleep(1e9)
	go pp()
}
