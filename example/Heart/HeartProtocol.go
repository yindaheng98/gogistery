package Heart

import (
	"fmt"
	"gogistery/Heart"
	"gogistery/Heartbeat"
	"gogistery/Protocol"
	ExampleHeartbeat "gogistery/example/Heartbeat"
	"time"
)

type RequestSendOption struct {
	id      string
	timeout time.Duration
	retryN  int64
	n       uint64
}
type ResponseSendOption struct {
	id string
	n  uint64
}

func (o RequestSendOption) Print() string {
	return fmt.Sprintf("RequestSendOption{id:%s,timeout:%d,retryN:%d,n:%d}", o.id, o.timeout, o.retryN, o.n)
}
func (o ResponseSendOption) Print() string {
	return fmt.Sprintf("ResponseSendOption{id:%s,n:%d}", o.id, o.n)
}

type RequesterHeartProtocol struct {
	requester *Heartbeat.Requester
	lastReq   Heart.TobeSendRequest
}

func (r RequesterHeartProtocol) Request(request Heart.TobeSendRequest) (Protocol.Response, error) {
	s := "\n------RequesterHeartProtocol.Request------>"
	r.lastReq = request
	req, opt := request.Request, request.Option.(RequestSendOption)
	s += fmt.Sprintf("A %s is sending with %s. ", req.String(), opt.Print())
	response, err := r.requester.Send(req, opt.timeout, opt.retryN)
	if err != nil {
		fmt.Print(s + fmt.Sprintf("Send failed, and the error is: %s", err.Error()))
		return nil, err
	}
	fmt.Print(s + fmt.Sprintf("Send success, and the response is: %s", response.String()))
	return response, nil
}

func (r RequesterHeartProtocol) Beat(response Protocol.Response, beat func(Heart.TobeSendRequest)) {
	s := "\n------RequesterHeartProtocol.Beat------>"
	req, opt := r.lastReq.Request, r.lastReq.Option.(RequestSendOption)
	s += fmt.Sprintf("A beat was success with a response %s. ", response.String())
	opt.n++
	if opt.n < 10 {
		fmt.Print(s + fmt.Sprintf("And the next beat is %s,%s. ", req.String(), opt.Print()))
		beat(Heart.TobeSendRequest{Request: req, Option: opt})
	} else {
		fmt.Print(s + "And it's the end of beating.")
	}
}

type ResponserHeartProtocol struct {
	responser *Heartbeat.Responser
	lastRes   Heart.TobeSendResponse
}

func (r ResponserHeartProtocol) Response() (Protocol.Request, error, func(Heart.TobeSendResponse)) {
	s := "\n------ResponserHeartProtocol.Response------>"
	req, err, resF := r.responser.Recv()
	if err == nil {
		fmt.Print(s + fmt.Sprintf("A request %s arrived. ", req.String()))
	} else {
		fmt.Print(s + fmt.Sprintf("A error %s occurred. ", err.Error()))
	}
	return req, err, func(response Heart.TobeSendResponse) {
		r.lastRes = response
		fmt.Printf("\n------>A response %s will be sent back with an option %s",
			response.Response.String(), response.Option.(ResponseSendOption).Print())
		resF(response.Response)
	}
}

func (r ResponserHeartProtocol) Beat(request Protocol.Request) Heart.TobeSendResponse {
	s := "\n------ResponserHeartProtocol.Beat------>"
	s += fmt.Sprintf("A request %s arrived. ", request.String())
	if r.lastRes.Option == nil {
		r.lastRes = Heart.TobeSendResponse{
			Response: Protocol.TobeSendResponse{
				Response: ExampleHeartbeat.Response{ID: "0"},
				Option:   ExampleHeartbeat.ResponseSendOption{ID: "0"}},
			Option: ResponseSendOption{"0", 0}}
	}
	o := r.lastRes.Option.(ResponseSendOption)
	o.n++
	res := Heart.TobeSendResponse{Response: r.lastRes.Response, Option: o}
	fmt.Print(s + fmt.Sprintf("And the response will be %s with an option %s", res.Response.String(), o.Print()))
	return res
}
