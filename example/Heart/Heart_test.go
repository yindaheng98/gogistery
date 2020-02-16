package Heart

import (
	"fmt"
	"gogistery/Heart"
	"gogistery/Protocol"
	ExampleProtocol "gogistery/example/Protocol"
	"math/rand"
	"testing"
	"time"
)

func ChanNetRequesterHeartTest(t *testing.T, RegistrantID string, initAddr string) {
	s := fmt.Sprintf("--ChanNetRequesterHeartTest(RegistrantID:%s,initAddr:%s)-->", RegistrantID, initAddr)
	info := ExampleProtocol.RegistrantInfo{
		ID:     RegistrantID,
		Option: ExampleProtocol.ResponseSendOption{Timestamp: time.Now()},
	}
	requester := Heart.NewRequesterHeart(
		NewRequesterHeartProtocol(info, 10),
		ExampleProtocol.NewChanNetRequestProtocol())
	requester.Events.Retry.AddHandler(func(o Protocol.TobeSendRequest, err error) {
		t.Log(s + fmt.Sprintf("A request %s retryed because %s", o.String(), err.Error()))
	})
	requester.Events.Retry.Enable()
	go func() {
		t.Log(s + fmt.Sprintf("RequesterHeart started with info %s.", info.String()))
		err := requester.RunBeating(Protocol.TobeSendRequest{
			Request: Protocol.Request{RegistrantInfo: info, Disconnect: false},
			Option:  ExampleProtocol.RequestSendOption{RequestAddr: initAddr, Timestamp: time.Now()},
		}, 10e9, 10)
		if err != nil {
			t.Log(s+"RequesterHeart stopped with error: %s.", err.Error())
		} else {
			t.Log(s + "RequesterHeart stopped normally.")
		}
	}()
}

func ChanNetResponserHeartTest(t *testing.T, RegistryID string) string {
	s := fmt.Sprintf("--ChanNetRequesterHeartTest(RegistryID:%s)-->", RegistryID)
	proto := ExampleProtocol.NewChanNetResponseProtocol()
	info := ExampleProtocol.RegistryInfo{
		ID:         RegistryID,
		Option:     ExampleProtocol.RequestSendOption{RequestAddr: proto.GetAddr(), Timestamp: time.Now()},
		Candidates: []Protocol.RegistryInfo{}}
	responser := Heart.NewResponserHeart(NewResponserHeartProtocol(info, 2e9, 5), proto)
	responser.Event.Error.AddHandler(func(err error) {
		t.Log(s + fmt.Sprintf("An error occurred: %s", err.Error()))
	})
	responser.Event.Error.Enable()
	go func() {
		t.Log(s + fmt.Sprintf("ResponserHeart started with info %s.", info.String()))
		go responser.RunBeating()
		time.Sleep(5e9)
		responser.Stop()
		responser.Stop()
		go responser.RunBeating()
		time.Sleep(1e9)
		responser.Stop()
	}()
	return proto.GetAddr()
}

const REQUESTERN = 30
const RESPONSERN = 10

func TestChanNetHeart(t *testing.T) {
	responsers := make([]string, RESPONSERN)
	for i := 0; i < RESPONSERN; i++ {
		responsers[i] = ChanNetResponserHeartTest(t, fmt.Sprintf("RESPONSER_%02d", i))
	}
	for i := 0; i < REQUESTERN; i++ {
		ChanNetRequesterHeartTest(t, fmt.Sprintf("REQUESTER_%02d", i), responsers[rand.Intn(RESPONSERN)])
	}
	time.Sleep(8e9)
}
