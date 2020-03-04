package heart

import (
	"context"
	"fmt"
	ExampleProtocol "github.com/yindaheng98/gogistry/example/protocol"
	"github.com/yindaheng98/gogistry/heart/requester"
	"github.com/yindaheng98/gogistry/heart/responser"
	"github.com/yindaheng98/gogistry/protocol"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func ChanNetRequesterHeartTest(t *testing.T, ctx context.Context, RegistrantID string, initAddr string, wg *sync.WaitGroup) {
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
	heart.Handlers.DisconnectionHandler = func(response protocol.Response, err error) {
		t.Log(s + fmt.Sprintf("Disonnection-->%s,%s", err, response.String()))
	}
	heart.Handlers.RetryHandler = func(o protocol.TobeSendRequest, err error) {
		t.Log(s + fmt.Sprintf("A request %s retryed because %s", o.String(), err.Error()))
	}
	go func() {
		t.Log(s + fmt.Sprintf("RequesterHeart started with info %s.", info.String()))
		err := heart.RunBeating(ctx, protocol.TobeSendRequest{
			Request: protocol.Request{RegistrantInfo: info, Disconnect: false},
			Option:  ExampleProtocol.RequestSendOption{RequestAddr: initAddr, Timestamp: time.Now()},
		}, 1e9, 3)
		if err != nil {
			t.Log(s+"RequesterHeart stopped with error: %s.", err.Error())
		} else {
			t.Log(s + "RequesterHeart stopped normally.")
		}
		wg.Done()
	}()
}

const TEST_TIMEOUT = 1e9

func ChanNetResponserHeartTest(t *testing.T, ctx context.Context, RegistryID string, wg *sync.WaitGroup) string {
	s := fmt.Sprintf("--ChanNetRequesterHeartTest(RegistryID:%s)-->", RegistryID)
	proto := ExampleProtocol.NewChanNetResponseProtocol()
	info := ExampleProtocol.RegistryInfo{
		ID:         RegistryID,
		Option:     ExampleProtocol.RequestSendOption{RequestAddr: proto.GetAddr(), Timestamp: time.Now()},
		Candidates: []protocol.RegistryInfo{}}
	heart := responser.NewHeart(NewResponserHeartBeater(info, TEST_TIMEOUT), proto)
	heart.ErrorHandler = func(err error) {
		t.Log(s + fmt.Sprintf("An error occurred: %s", err.Error()))
	}
	go func() {
		t.Log(s + fmt.Sprintf("ResponserHeart started with info %s.", info.String()))
		heart.RunBeating(ctx)
		wg.Done()
	}()
	return proto.GetAddr()
}

const REQUESTERN = 30
const RESPONSERN = 10

func TestChanNetHeart(t *testing.T) {
	responsers := make([]string, RESPONSERN)
	responsersWG := new(sync.WaitGroup)
	responsersWG.Add(RESPONSERN)
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 1e9)
	for i := 0; i < RESPONSERN; i++ {
		responsers[i] = ChanNetResponserHeartTest(t, timeoutCtx, fmt.Sprintf("RESPONSER_%02d", i), responsersWG)
	}
	requestersWG := new(sync.WaitGroup)
	requestersWG.Add(REQUESTERN)
	for i := 0; i < REQUESTERN; i++ {
		ChanNetRequesterHeartTest(t, timeoutCtx, fmt.Sprintf("REQUESTER_%02d", i), responsers[rand.Intn(RESPONSERN)], requestersWG)
	}
	time.Sleep(1e8)
	cancel()
	responsersWG.Wait()
	requestersWG.Wait()
	cancel()
}
