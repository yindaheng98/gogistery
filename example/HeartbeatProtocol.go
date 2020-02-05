package example

import (
	"errors"
	"fmt"
	"gogistery/Heartbeat"
	"math/rand"
	"sync/atomic"
	"time"
)

type Response struct {
	id string
}

type Request struct {
	id string
}

func (r Response) Print() string {
	return fmt.Sprintf("Response{id:%s}", r.id)
}
func (r Request) Print() string {
	return fmt.Sprintf("Request{id:%s}", r.id)
}

type RequestSendOption struct {
	id   string
	addr string
}
type ResponseSendOption struct {
	id string
}

func (o RequestSendOption) Print() string {
	return fmt.Sprintf("RequestSendOption{id:%s,addr:%s}", o.id, o.addr)
}
func (o ResponseSendOption) Print() string {
	return fmt.Sprintf("ResponseSendOption{id:%s}", o.id)
}

type RequestBeatProtocol struct {
	src       *rand.Source
	failRate  int32
	responseN uint32
}

func (t *RequestBeatProtocol) Request(requestChan <-chan Heartbeat.TobeSendRequest, responseChan chan<- Heartbeat.ReceivedResponse) {
	atomic.AddUint32(&t.responseN, 1)
	protoRequest := <-requestChan
	request, option := protoRequest.Request.(Request), protoRequest.Option.(RequestSendOption)
	s := "\n------RequestBeatProtocol------>"
	s += fmt.Sprintf("It was sending attempt %02d in protocol. %s is sending with %s. ",
		t.responseN, request.Print(), option.Print())
	timeout := time.Duration(rand.Int63n(1e3) * 1e6)
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
	responseChan <- Heartbeat.ReceivedResponse{Response: Response{fmt.Sprintf("%02d", t.responseN)}}
	fmt.Print(s + "This Sending was success.")
}

type ResponseBeatProtocol struct {
	src      *rand.Source
	failRate int32
	id       string
}

func (t ResponseBeatProtocol) Response(requestChan chan<- Heartbeat.ReceivedRequest, responseChan <-chan Heartbeat.TobeSendResponse) {
	time.Sleep(time.Duration(rand.Int31n(1e3) * 1e3))
	request := Request{t.id}
	s := "\n------ResponseBeatProtocol------>"
	s += fmt.Sprintf("A request %s arrived in protocol. ", request.Print())

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
		response, option := protoResponse.Response.(Response), protoResponse.Option.(ResponseSendOption)
		fmt.Print(s + fmt.Sprintf("And the Response is %s, with the option %s",
			response.Print(), option.Print()))
	}
}
