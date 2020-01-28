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
	st.Callback.Started = func() {
		t.Log("I'm started.")
	}
	st.Callback.Stopped = func() {
		t.Log("I'm stopped.")
	}
	go st.Run(r)
	go st.Run(r)
	go st.Run(r)
	time.Sleep(1e9)
}

func TestSingleProcessor(t *testing.T) {
	i := int32(0)
	p := NewProcessor()
	p.Callback.Started = func() {
		t.Log("I'm started.")
	}
	p.Callback.Stopped = func() {
		t.Log("I'm stopped.")
	}
	pp := func() {
		atomic.AddInt32(&i, 1)
		t.Log(fmt.Sprintf("I'm No.%d", i))
		atomic.AddInt32(&i, -1)
	}
	go p.Start(pp)
	go p.Stop()
	go p.Start(pp)
	go p.Stop()
	go p.Start(pp)
	go p.Stop()
	go p.Stop()
	go p.Stop()
	go p.Stop()
	go p.Stop()
	time.Sleep(1e9)
}
