package Heartbeat

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type TestResponseProtocol struct {
	src      *rand.Source
	failRate int32
	id       string
}

type TestResponseOption struct {
	id string
}

func (t TestResponseProtocol) Response(requestChan chan<- ReceivedRequest, responseChan <-chan ProtocolResponseSendOption) {
	time.Sleep(time.Duration(rand.Int31n(1e3) * 1e3))
	request := TestRequest{t.id}
	s := fmt.Sprintf("\nA request TestRequest{id:%s} arrived in protocol. ", request.id)

	r := rand.New(*t.src).Int31n(100)
	if r < t.failRate {
		s += "This Receiving was failed. "
		requestChan <- ReceivedRequest{nil, errors.New(fmt.Sprintf(
			"Your fail rate is %d%%, but this random output is %02d, so failed.", t.failRate, r))}
	} else {
		requestChan <- ReceivedRequest{request, nil}
		s += "This Receiving was success. "
	}

	protoResponse, ok := <-responseChan
	if !ok {
		fmt.Print(s + "But the Response was timeouted.")
	} else {
		response, option := protoResponse.response, protoResponse.option
		fmt.Print(s + fmt.Sprintf("And the Response is TestResponse{id:%s}, with the option TestResponseOption{id:%s}",
			response.(TestResponse).id,
			option.(TestResponseOption).id))
	}
}

func testRes(i uint64, logger func(string)) {
	responser := NewResponser(TestResponseProtocol{&src, 30, fmt.Sprintf("%d", i)})
	responseChan := make(chan ProtocolResponseSendOption, 1)
	request, err := responser.Recv(responseChan, 1e7) /*********将该值调低可模拟超时情况**********/
	d := time.Duration(rand.Int31n(1e3) * 1e3)
	if err != nil {
		logger(err.Error())
		time.Sleep(d)
		responseChan <- ProtocolResponseSendOption{TestResponse{fmt.Sprintf("error%02d", i)},
			TestResponseOption{fmt.Sprintf("error%02d", i)}}
	} else {
		logger(fmt.Sprintf("A request TestRequest{id:%s} arrived. Response will be sent back in %d",
			request.(TestRequest).id, d))
		time.Sleep(d)
		responseChan <- ProtocolResponseSendOption{TestResponse{fmt.Sprintf("%02d", i)},
			TestResponseOption{fmt.Sprintf("%02d", i)}}
	}
	close(responseChan)
}

//单次Heartbeat
func TestResponser(t *testing.T) {
	for i := uint64(0); i < 30; i++ {
		testRes(i, func(s string) {
			t.Log(s)
		})
	}
}
