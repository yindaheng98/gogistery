package Emitter

import (
	"errors"
	"fmt"
	"gogistery/base"
	"gogistery/base/Errors"
	"testing"
	"time"
)

type TestInfo struct {
	id      string
	addr    string
	timeout time.Duration
	retryN  uint32
}

func (e *TestInfo) IsDisconnect() bool {
	return false
}

func (e *TestInfo) GetID() string {
	return e.id
}

func (e *TestInfo) GetAddr() string {
	return e.addr
}

func (e *TestInfo) GetTimeout() time.Duration {
	return e.timeout
}

func (e *TestInfo) GetRetryN() uint32 {
	return e.retryN
}

func NewTestInfo(id, addr string, i uint32) *TestInfo {
	return &TestInfo{id: id, addr: addr, timeout: time.Duration(i), retryN: i}
}

func NewTestLinkInfo(i uint32) base.LinkInfo {
	si := string(i)
	return base.NewLinkInfo(NewTestInfo("sender"+si, si+".send", i), NewTestInfo("receiver"+si, si+".recv", i))
}

func TestLinkErrorEmitter(t *testing.T) {
	emitter := NewLinkErrorEmitter()
	emitter.AddHandler(func(e Errors.LinkError) {
		t.Log(fmt.Sprintf("Here is a handler, I'm handling: %s.", e.Error()))
	})
	emitter.Enable()
	go emitter.Emit(Errors.NewLinkError(errors.New("error1"), NewTestLinkInfo(1)))
	emitter.AddHandler(func(e Errors.LinkError) {
		t.Log(fmt.Sprintf("Here is another handler, I'm handling: %s.", e.Error()))
	})
	go emitter.Emit(Errors.NewLinkError(errors.New("error2"), NewTestLinkInfo(2)))
	//go emitter.Disable()
	go emitter.Enable()
	go emitter.Emit(Errors.NewLinkError(errors.New("error3"), NewTestLinkInfo(3)))
	emitter.AddHandler(func(e Errors.LinkError) {
		t.Log(fmt.Sprintf("Here is handler2, I'm handling: %s.", e.Error()))
	})
	go emitter.Emit(Errors.NewLinkError(errors.New("error4"), NewTestLinkInfo(4)))
	time.Sleep(1e9 * 3)
	go emitter.Disable()
}

func TestSenderInfoEmitter(t *testing.T) {
	emitter := NewSenderInfoEmitter()
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log(
			fmt.Sprintf("Handler1->A SenderInfo{id:%s} has just arrived.",
				senderInfo.GetID()))
	})
	emitter.Enable()
	go emitter.Emit(NewTestInfo("event0", "0.event", 1))
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log(
			fmt.Sprintf("Handler1->A SenderInfo{id:%s} has just arrived.",
				senderInfo.GetID()))
	})
	go emitter.Emit(NewTestInfo("event1", "1.event", 1))
	go emitter.Enable()
	go emitter.Emit(NewTestInfo("event2", "2.event", 2))
	go emitter.Disable()
	time.Sleep(1e9 * 3)
}

func TestReceiverInfoEmitter(t *testing.T) {
	emitter := NewReceiverInfoEmitter()
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler1->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	emitter.Enable()
	go emitter.Emit(NewTestInfo("event0", "0.event", 1))
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler2->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(NewTestInfo("event1", "1.event", 1))
	go emitter.Enable()
	go emitter.Emit(NewTestInfo("event2", "2.event", 2))
	go emitter.Disable()
	time.Sleep(1e9 * 3)
}

func TestEmptyEmitter(t *testing.T) {
	emitter := NewEmptyEmitter()
	emitter.AddHandler(func() {
		t.Log("I'm handler 1")
	})
	emitter.Enable()
	go emitter.Emit()
	emitter.AddHandler(func() {
		t.Log("I'm handler 2")
	})
	go emitter.Emit()
	go emitter.Disable()
	emitter.AddHandler(func() {
		t.Log("I'm handler 3")
	})
	emitter.AddHandler(func() {
		t.Log("I'm handler 4")
	})
	go emitter.Enable()
	go emitter.Emit()
	go emitter.Emit()
	go emitter.Disable()
	time.Sleep(1e9 * 3)
}
