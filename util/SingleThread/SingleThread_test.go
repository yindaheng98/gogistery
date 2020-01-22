package SingleThread

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleThread(t *testing.T) {
	i := int32(0)
	st := New(func() {
		atomic.AddInt32(&i, 1)
		t.Log(fmt.Sprintf("I'm No.%d", i))
		atomic.AddInt32(&i, -1)
	})
	go st.Run()
	go st.Run()
	go st.Run()
	go st.Run()
	go st.Run()
	time.Sleep(1e9)
}
