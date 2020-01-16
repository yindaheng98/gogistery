package Emitter

import (
	"errors"
	"testing"
)

func TestErrorEmitter(t *testing.T) {
	emitter := NewErrorEmitter()
	emitter.AddHandler(func(e error) {
		t.Log("Here is another handler, I'm handling: " + e.Error())
	})
	emitter.Start()
	go emitter.Emit(errors.New("error1"))
	emitter.AddHandler(func(e error) {
		t.Log("Here is another handler, I'm handling: " + e.Error())
	})
	go emitter.Emit(errors.New("error2"))
	go emitter.Start()
	go emitter.Emit(errors.New("error3"))
	emitter.AddHandler(func(e error) {
		t.Log("Here is handler2, I'm handling: " + e.Error())
	})
	go emitter.Emit(errors.New("error4"))
	go emitter.Stop()
}
