package heart

import (
	"fmt"
	ExampleProtocol "gogistery/example/protocol"
	"gogistery/heart/requester"
	"gogistery/heart/responser"
	"gogistery/protocol"
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
	heart := requester.NewHeart(
		NewRequesterHeartBeater(info, 10),
		ExampleProtocol.NewChanNetRequestProtocol())
	heart.Handlers.NewConnectionHandler = func(info protocol.Response) {
		t.Log(s + fmt.Sprintf("New Connection-->%s", info.String()))
	}
	heart.Handlers.UpdateConnectionHandler = func(info protocol.Response) {
		t.Log(s + fmt.Sprintf("Update Connection-->%s", info.String()))
	}
	heart.Handlers.DisconnectionHandler = func(request protocol.TobeSendRequest, err error) {
		t.Log(s + fmt.Sprintf("Disonnection-->%s,%s", err, request.String()))
	}
	heart.Handlers.RetryHandler = func(o protocol.TobeSendRequest, err error) {
		t.Log(s + fmt.Sprintf("A request %s retryed because %s", o.String(), err.Error()))
	}
	go func() {
		t.Log(s + fmt.Sprintf("RequesterHeart started with info %s.", info.String()))
		err := heart.RunBeating(protocol.TobeSendRequest{
			Request: protocol.Request{RegistrantInfo: info, Disconnect: false},
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
		Candidates: []protocol.RegistryInfo{}}
	heart := responser.NewHeart(NewResponserHeartBeater(info, 2e9, 5), proto)
	heart.ErrorHandler = func(err error) {
		t.Log(s + fmt.Sprintf("An error occurred: %s", err.Error()))
	}
	go func() {
		t.Log(s + fmt.Sprintf("ResponserHeart started with info %s.", info.String()))
		go heart.RunBeating()
		time.Sleep(5e9)
		heart.Stop()
		heart.Stop()
		go heart.RunBeating()
		time.Sleep(1e9)
		heart.Stop()
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
