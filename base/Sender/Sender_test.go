package Sender

import (
	"errors"
	"fmt"
	"gogistery/base"
	"gogistery/base/Errors"
	"math/rand"
	"testing"
	"time"
)

type TestReceiverInfo struct {
	addr    string
	timeout time.Duration
	retryN  uint32
}

var src = rand.NewSource(10)

func NewTestReceiverInfo() *TestReceiverInfo {
	return &TestReceiverInfo{
		addr: fmt.Sprintf("%d.%d.%d.%d:%d",
			rand.New(src).Int31n(255),
			rand.New(src).Int31n(255),
			rand.New(src).Int31n(255),
			rand.New(src).Int31n(255),
			rand.New(src).Int31n(25565)),
		timeout: time.Duration(rand.New(src).Int63n(100) * 1e6),
		retryN:  rand.New(src).Uint32() % 10}
}

func (info *TestReceiverInfo) GetAddr() string {
	return info.addr
}

func (info *TestReceiverInfo) GetTimeout() time.Duration {
	return info.timeout
}

func (info *TestReceiverInfo) GetRetryN() uint32 {
	return info.retryN
}

type TestSenderInfo struct {
	src      *rand.Source
	failRate int32
}

func (info *TestSenderInfo) Send(addr string, timeout time.Duration) (base.ReceiverInfo, error) {
	fmt.Printf("TestSenderInfo is sending to %s with timeout %s.", addr, timeout)
	r := rand.New(*info.src).Int31n(100)
	if r < info.failRate {
		fmt.Printf("This Send will failed.\n")
		return nil, errors.New(fmt.Sprintf(
			"Your fail rate is %d%%, and this random output is %d, so failed.", info.failRate, r))
	}
	fmt.Printf("This Send will success.\n")
	return NewTestReceiverInfo(), nil
}

func TestSender(t *testing.T) {
	testSenderInfo := TestSenderInfo{&src, 30}
	sender := New(&testSenderInfo, "initAddr:0", 0, 10)
	sender.Events.Start.AddHandler(func() {
		t.Log("A start event occurred.")
	})
	sender.Events.Stop.AddHandler(func() {
		t.Log("A stop event occurred.")
	})
	sender.Events.Connect.AddHandler(func(info base.ReceiverInfo) {
		t.Log(
			fmt.Sprintf("A connect event occurred: base.ReceiverInfo{addr:%s,timeout:%s,retryN:%d}",
				info.GetAddr(),
				info.GetTimeout(),
				info.GetRetryN()))
	})
	sender.Events.Disconnect.AddHandler(func(e Errors.LinkError) {
		t.Log(
			fmt.Sprintf("A disconnect event occurred, its error message is %s, and its receiver addr is \" %s \"",
				e.Error(),
				e.LinkInfo().ReceiverInfo().GetAddr()))
	})
	sender.Events.Retry.AddHandler(func(e Errors.LinkError) {
		t.Log(
			fmt.Sprintf("A retry event occurred, its error message is %s, and its receiver addr is \" %s \"",
				e.Error(),
				e.LinkInfo().ReceiverInfo().GetAddr()))
	})
	sender.Events.Error.AddHandler(func(e Errors.LinkError) {
		t.Log(
			fmt.Sprintf("A error occurred, its error message is %s, and its receiver addr is \" %s \"",
				e.Error(),
				e.LinkInfo().ReceiverInfo().GetAddr()))
	})
	go sender.Events.EnableAll()
	go sender.Connect()
	time.Sleep(1e9 * 3)
	go sender.Events.DisableAll()
	time.Sleep(1e6)
	go sender.Events.EnableAll()
	sender.Disconnect()
	go sender.Events.DisableAll()
	time.Sleep(1e9 * 1)
}
