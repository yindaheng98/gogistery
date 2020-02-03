package Heartbeat

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

type TestResponse struct {
	id string
}

type TestRequest struct {
	id string
}

type TestRequestOption struct {
	id   string
	addr string
}

type TestRequestProtocol struct {
	src       *rand.Source
	failRate  int32
	responseN uint32
}

func (t *TestRequestProtocol) Request(requestChan <-chan ProtocolRequestSendOption, responseChan chan<- ReceivedResponse) {
	atomic.AddUint32(&t.responseN, 1)
	protoRequest := <-requestChan
	request, option := protoRequest.request, protoRequest.option
	s := fmt.Sprintf("\nIt was sending attempt %02d in protocol. TestRequest{id:%s} is sending to %s. ",
		t.responseN, request.(TestRequest).id, option.(TestRequestOption).addr)
	timeout := time.Duration(rand.Int63n(1e3) * 1e3)
	s += fmt.Sprintf("Response will arrived in %d. ", timeout)
	defer func() {
		if recover() != nil {
			fmt.Print(s + "This Sending was timeout.")
		}
	}()
	r := rand.New(*t.src).Int31n(100)
	if r < t.failRate {
		fmt.Print(s + "This Sending was failed.")
		responseChan <- ReceivedResponse{nil, errors.New(fmt.Sprintf(
			"Your fail rate is %d%%, but this random output is %02d, so failed.", t.failRate, r))}
		return
	}
	time.Sleep(timeout)
	responseChan <- ReceivedResponse{TestResponse{fmt.Sprintf("%02d", t.responseN)}, nil}
	fmt.Print(s + "This Sending was success.")
}

var src = rand.NewSource(10)

func testReq(i uint64, logger func(string)) {
	requester := NewRequester(&TestRequestProtocol{&src, 30, 0})
	requester.Events.Retry.AddHandler(func(o ProtocolRequestSendOption, err error) {
		logger(fmt.Sprintf("An retry was occured. error: %s", err.Error()))
	})
	requester.Events.Retry.Enable()
	response, err := requester.Send(ProtocolRequestSendOption{
		TestRequest{fmt.Sprintf("%02d", i)},
		TestRequestOption{
			fmt.Sprintf("%02d", i),
			fmt.Sprintf("%02d.%02d.%02d.%02d", i, i, i, i)}},
		time.Duration(5e7), /*********将该值调低可模拟超时情况**********/
		10)
	if err != nil {
		logger(fmt.Sprintf("No.%02d test failed. err is %s", i, err.Error()))
		return
	}
	logger(fmt.Sprintf("No.%02d sending test succeed. response is TestResponse{id:%s}", i, response.(TestResponse).id))
}

//单次Heartbeat
func TestRequester(t *testing.T) {
	for i := uint64(0); i < 30; i++ {
		testReq(i, func(s string) {
			t.Log(s)
		})
	}
	time.Sleep(1000)
}
