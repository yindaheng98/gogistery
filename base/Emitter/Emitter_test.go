package Emitter

import (
	"errors"
	"fmt"
	"gogistery/base"
	"gogistery/base/Error"
	"testing"
	"time"
)

func TestErrorEmitter(t *testing.T) {
	emitter := NewConnectionErrorEmitter()
	emitter.AddHandler(func(e Error.ConnectionError) {
		t.Log(fmt.Sprintf("Here is a handler, I'm handling: %s. Its code is %d", e.Error(), e.Code))
	})
	emitter.Start()
	go emitter.Emit(Error.NewConnectionError(errors.New("error1"), Error.CONNECTION_ERROR_InitFailed))
	emitter.AddHandler(func(e Error.ConnectionError) {
		t.Log(fmt.Sprintf("Here is another handler, I'm handling: %s. Its code is %d", e.Error(), e.Code))
	})
	go emitter.Emit(Error.NewConnectionError(errors.New("error2"), Error.CONNECTION_ERROR_ConnectionInterrupt))
	go emitter.Stop()
	go emitter.Start()
	go emitter.Emit(Error.NewConnectionError(errors.New("error3"), Error.CONNECTION_ERROR_RetryFailed))
	emitter.AddHandler(func(e Error.ConnectionError) {
		t.Log(fmt.Sprintf("Here is handler2, I'm handling: %s. Its code is %d", e.Error(), e.Code))
	})
	go emitter.Emit(Error.NewConnectionError(errors.New("error4"), Error.CONNECTION_ERROR_InitFailed))
	go emitter.Stop()
	time.Sleep(1e9 * 3)
}

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
	emitter.Start()
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
	go emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3, 3})
	go emitter.Stop()
	time.Sleep(1e9 * 3)
}

func TestReceiverInfoEmitter(t *testing.T) {
	emitter := NewReceiverInfoEmitter()
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler1->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_1", 1, 1})
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler2->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(&TestInfo{"addr_in_event_2", 2, 2})
	go emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3, 3})
	go emitter.Stop()
	time.Sleep(1e9 * 3)
}

func TestEmptyEmitter(t *testing.T) {
	emitter := NewEmptyEmitter()
	emitter.AddHandler(func() {
		t.Log("I'm handler 1")
	})
	emitter.Start()
	go emitter.Emit()
	emitter.AddHandler(func() {
		t.Log("I'm handler 2")
	})
	go emitter.Emit()
	go emitter.Stop()
	emitter.AddHandler(func() {
		t.Log("I'm handler 3")
	})
	emitter.AddHandler(func() {
		t.Log("I'm handler 4")
	})
	go emitter.Start()
	go emitter.Emit()
	go emitter.Emit()
}
