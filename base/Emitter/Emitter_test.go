package Emitter

import (
	"errors"
	"fmt"
	"gogistery/base"
	"testing"
	"time"
)

func TestErrorEmitter(t *testing.T) {
	emitter := NewErrorEmitter()
	emitter.AddHandler(func(e error) {
		t.Log("Here is a handler, I'm handling: " + e.Error())
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
	time.Sleep(1e9 * 3)
}

type TestInfo struct {
	addr    string
	timeout time.Duration
}

func (e *TestInfo) Send(addr string, timeout time.Duration) (base.ReceiverInfo, error) {
	fmt.Printf("I'm sending messages to %s with duration %s", addr, timeout)
	return &TestInfo{e.addr, timeout}, nil
}

func (e *TestInfo) GetAddr() string {
	return e.addr
}

func (e *TestInfo) GetTimeout() time.Duration {
	return e.timeout
}

func TestSenderInfoEmitter(t *testing.T) {
	emitter := NewSenderInfoEmitter()
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log("Handler1->Here is a SenderInfo handler, a SenderInfo has just arrived, ")
		receiverInfo, err := senderInfo.Send("addr_in_handler_1", 100)
		t.Log(fmt.Sprintf("Handler1->SenderInfo have just been sended to addr_in_handler_1"))
		if err != nil {
			t.Log("Handler1->An error occurred: " + err.Error())
		}
		t.Log(fmt.Sprintf("Handler1->Receiver replied a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_1", 1})
	emitter.AddHandler(func(senderInfo base.SenderInfo) {
		t.Log("Handler2->Here is a SenderInfo handler, a SenderInfo has just arrived, ")
		receiverInfo, err := senderInfo.Send("addr_in_handler_1", 100)
		t.Log(fmt.Sprintf("Handler2->SenderInfo have just been sended to addr_in_handler_2"))
		if err != nil {
			t.Log("Handler2->An error occurred: " + err.Error())
		}
		t.Log(fmt.Sprintf("Handler2->Receiver replied a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(&TestInfo{"addr_in_event_2", 2})
	go emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3})
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
	go emitter.Emit(&TestInfo{"addr_in_event_1", 1})
	emitter.AddHandler(func(receiverInfo base.ReceiverInfo) {
		t.Log(fmt.Sprintf("Handler2->a ReceiverInfo has just arrived, it has a addr %s and a timeout %s",
			receiverInfo.GetAddr(),
			receiverInfo.GetTimeout()))
	})
	go emitter.Emit(&TestInfo{"addr_in_event_2", 2})
	go emitter.Start()
	go emitter.Emit(&TestInfo{"addr_in_event_3", 3})
	go emitter.Stop()
	time.Sleep(1e9 * 3)
}
