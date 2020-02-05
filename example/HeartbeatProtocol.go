package example

import (
	"errors"
	"fmt"
	"gogistery/Heartbeat"
	"math/rand"
	"sync/atomic"
	"time"
)

type MyResponse struct {
	id string
}

type MyRequest struct {
	id string
}

type MyRequestOption struct {
	id   string
	addr string
}

type MyResponseOption struct {
	id string
}

type MyRequestBeatProtocol struct {
	src       *rand.Source
	failRate  int32
	responseN uint32
}

func (t *MyRequestBeatProtocol) Request(requestChan <-chan Heartbeat.TobeSendRequest, responseChan chan<- Heartbeat.ReceivedResponse) {
	atomic.AddUint32(&t.responseN, 1)
	protoRequest := <-requestChan
	request, option := protoRequest.Request, protoRequest.Option
	s := "\n------MyRequestBeatProtocol------>"
	s += fmt.Sprintf("It was sending attempt %02d in protocol. MyRequest{id:%s} is sending to %s. ",
		t.responseN, request.(MyRequest).id, option.(MyRequestOption).addr)
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
		responseChan <- Heartbeat.ReceivedResponse{Error: errors.New(fmt.Sprintf(
			"Your fail rate is %d%%, but this random output is %02d, so failed.", t.failRate, r))}
		return
	}
	time.Sleep(timeout)
	responseChan <- Heartbeat.ReceivedResponse{Response: MyResponse{fmt.Sprintf("%02d", t.responseN)}}
	fmt.Print(s + "This Sending was success.")
}

type MyResponseBeatProtocol struct {
	src      *rand.Source
	failRate int32
	id       string
}

func (t MyResponseBeatProtocol) Response(requestChan chan<- Heartbeat.ReceivedRequest, responseChan <-chan Heartbeat.TobeSendResponse) {
	time.Sleep(time.Duration(rand.Int31n(1e3) * 1e3))
	request := MyRequest{t.id}
	s := "\n------MyResponseBeatProtocol------>"
	s += fmt.Sprintf("A request MyRequest{id:%s} arrived in protocol. ", request.id)

	r := rand.New(*t.src).Int31n(100)
	if r < t.failRate {
		s += "This Receiving was failed. "
		requestChan <- Heartbeat.ReceivedRequest{Error: errors.New(fmt.Sprintf(
			"Your fail rate is %d%%, but this random output is %02d, so failed.", t.failRate, r))}
	} else {
		requestChan <- Heartbeat.ReceivedRequest{Request: request}
		s += "This Receiving was success. "
	}

	protoResponse, ok := <-responseChan
	if !ok {
		fmt.Print(s + "But the Response was timeouted.")
	} else {
		response, option := protoResponse.Response, protoResponse.Option
		fmt.Print(s + fmt.Sprintf("And the Response is MyResponse{id:%s}, with the option MyResponseOption{id:%s}",
			response.(MyResponse).id,
			option.(MyResponseOption).id))
	}
}
