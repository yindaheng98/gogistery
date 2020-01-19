package Emitter

import (
	"testing"
	"time"
)

type event struct {
	name string
}

func TestEmitter(t *testing.T) {
	emitter := New()
	emitter.AddHandler(func(e interface{}) {
		t.Log("Here is a handler, I'm handling: " + e.(event).name)
	})
	emitter.Enable()
	go emitter.Emit(event{"I'm a event"})
	emitter.AddHandler(func(e interface{}) {
		t.Log("Here is another handler, I'm handling: " + e.(event).name)
	})
	go emitter.Emit(event{"I'm another event"})
	go emitter.Enable()
	go emitter.Emit(event{"I'm event2"})
	emitter.AddHandler(func(e interface{}) {
		t.Log("Here is handler2, I'm handling: " + e.(event).name)
	})
	go emitter.Emit(event{"I'm event3"})
	go emitter.Disable()
	go emitter.Enable()
	go emitter.Emit(event{"I'm event4"})
	go emitter.Disable()
	emitter.AddHandler(func(e interface{}) {
		t.Log("Here is handler3, I'm handling: " + e.(event).name)
	})
	go emitter.Emit(event{"I'm event5"})
	go emitter.Disable()
	time.Sleep(1e9 * 3)
}
