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
	addr    string
	timeout time.Duration
	retryN  uint32
}

func (e *TestInfo) Send(addr string, timeout time.Duration) (base.ReceiverInfo, error) {
	fmt.Printf("I'm sending messages to %s", addr)
	return &TestInfo{e.addr, timeout, 100}, nil
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

func NewTestInfo(addr string, i uint32) *TestInfo {
	return &TestInfo{addr: addr, timeout: time.Duration(i), retryN: i}
}

func NewTestLinkInfo(i uint32) base.LinkInfo {
	return base.NewLinkInfo(NewTestInfo("sender"+string(i), i), NewTestInfo("receiver"+string(i), i))
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
	go emitter.Disable()
	go emitter.Enable()
	go emitter.Emit(Errors.NewLinkError(errors.New("error3"), NewTestLinkInfo(3)))
	emitter.AddHandler(func(e Errors.LinkError) {
		t.Log(fmt.Sprintf("Here is handler2, I'm handling: %s.", e.Error()))
	})
	go emitter.Emit(Errors.NewLinkError(errors.New("error4"), NewTestLinkInfo(4)))
	go emitter.Disable()
	time.Sleep(1e9 * 3)
}

func TestSenderInfoEmitter(t *testing.T) {
	emitter := NewSenderInfoEmitter()
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log("Handler1->Here is a SenderInfo handler, a SenderInfo has just arrived, ")
		receiverInfo, err := senderInfo.Send("addr_in_handler_1", 1)
		t.Log(fmt.Sprintf("Handler1->SenderInfo have just been sended to addr_in_handler_1"))
		if err != nil {
			t.Log("Handler1->An error occurred: " + err.Error())
		}
		t.Log(fmt.Sprintf("Handler1->Receiver replied a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	emitter.Enable()
	go emitter.Emit(&TestInfo{"addr_in_event_1", 1, 1})
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log("Handler2->Here is a SenderInfo handler, a SenderInfo has just arrived, ")
		receiverInfo, err := senderInfo.Send("addr_in_handler_2", 2)
		t.Log(fmt.Sprintf("Handler2->SenderInfo have just been sended to addr_in_handler_2"))
		if err != nil {
			t.Log("Handler2->An error occurred: " + err.Error())
		}
		t.Log(fmt.Sprintf("Handler2->Receiver replied a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(&TestInfo{"addr_in_event_2", 2, 2})
	go emitter.Enable()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3, 3})
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
	go emitter.Emit(&TestInfo{"addr_in_event_1", 1, 1})
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler2->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(&TestInfo{"addr_in_event_2", 2, 2})
	go emitter.Enable()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3, 3})
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
