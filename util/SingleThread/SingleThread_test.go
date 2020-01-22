package SingleThread

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleThread(t *testing.T) {
	i := int32(0)
	st := New()
	r := func() {
		atomic.AddInt32(&i, 1)
		t.Log(fmt.Sprintf("I'm No.%d", i))
		atomic.AddInt32(&i, -1)
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
	go st.Run(r)
	go rr()
	go st.Run(r)
	go rr()
	go st.Run(r)
	go rr()
	go st.Run(r)
	go rr()
	time.Sleep(1e9)
}
