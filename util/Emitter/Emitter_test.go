package Emitter

import (
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	emitter := New()
	emitter.AddHandler(func(bytes []byte) {
		t.Log("Here is a handler, I'm handling: " + string(bytes))
	})
	emitter.Start()
	go emitter.Emit([]byte("I'm a event"))
	emitter.AddHandler(func(bytes []byte) {
		t.Log("Here is another handler, I'm handling: " + string(bytes))
	})
	go emitter.Emit([]byte("I'm another event"))
	go emitter.Start()
	go emitter.Emit([]byte("I'm event2"))
	emitter.AddHandler(func(bytes []byte) {
		t.Log("Here is the other handler, I'm handling: " + string(bytes))
	})
	go emitter.Emit([]byte("I'm event3"))
	//go emitter.Stop()
	go emitter.Start()
	go emitter.Emit([]byte("I'm event4"))
	emitter.AddHandler(func(bytes []byte) {
		t.Log("Here is the other handler, I'm handling: " + string(bytes))
	})
	go emitter.Emit([]byte("I'm event5"))
	time.Sleep(1e9 * 3)
}
