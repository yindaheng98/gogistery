package TimeoutValue

import (
	"fmt"
	"testing"
	"time"
)

type TestElement struct {
	id string
}

func (e *TestElement) NewAddedHandler() {
	fmt.Printf("Element %s was added.\n", e.id)
}

func (e *TestElement) TimeoutHandler() {
	fmt.Printf("Element %s is timeout.\n", e.id)
}

func TestTimeoutValue(t *testing.T) {
	v := New(&TestElement{"001"}, 1e8, func() {
		t.Log("Value 001 is timeout")
	})
	go v.Start()
	go v.Update(nil)
	go v.Start()
	go v.Update(nil)
	go v.Stop()
	go v.Update(nil)
	go v.Stop()
	go v.Update(nil)
	//go v.Stop()
	go v.Update(nil)
	//go v.Stop()
	go v.Update(nil)
	time.Sleep(1e8)
	//v.Stop()
}
