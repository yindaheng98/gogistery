package example

import (
	"fmt"
	"gogistery/Heartbeat"
	"gogistery/Protocol"
	"time"
)

type RequesterBeat struct {
	Request
	n uint64
}
type ResponserBeat struct {
	Response
	n uint64
}

func (r ResponserBeat) String() string {
	return fmt.Sprintf("ResponserBeat{%s,n:%d}", r.Response.String(), r.n)
}
func (r RequesterBeat) String() string {
	return fmt.Sprintf("RequesterBeat{%s,n:%d}", r.Request.String(), r.n)
}

type RequesterBeatSendOption struct {
	RequestSendOption
	timeout time.Duration
	retryN  int64
}
type ResponserBeatSendOption struct {
	ResponseSendOption
}

func (r RequesterBeatSendOption) String() string {
	return fmt.Sprintf("ResponserBeat{%s,timeout:%d,retryN:%d}", r.RequestSendOption.String(), r.timeout, r.retryN)
}
func (r ResponserBeatSendOption) String() string {
	return fmt.Sprintf("RequesterBeat{%s}", r.ResponseSendOption.String())
}

type RequesterHeartProtocol struct {
	requester *Heartbeat.Requester
}

func (r RequesterHeartProtocol) Request(request Protocol.TobeSendRequest) (Protocol.Response, error) {
	s := "\n------RequesterHeartProtocol.Request------>"
	req, opt := request.Request.(RequesterBeat), request.Option.(RequesterBeatSendOption)
	s += fmt.Sprintf("A %s is sending with %s. ", req.String(), opt.String())
	response, err := r.requester.Send(
		Protocol.TobeSendRequest{Request: req.Request, Option: opt.RequestSendOption},
		opt.timeout, opt.retryN)
	if err != nil {
		fmt.Print(s + fmt.Sprintf("Send failed, and the error is: %s", err.Error()))
		return nil, err
	}
	fmt.Print(s + fmt.Sprintf("Send success, and the response is: %s", response.(Response).String()))
	return ResponserBeat{response.(Response), req.n}, nil
}
func (r RequesterHeartProtocol) Beat(request Protocol.TobeSendRequest, response Protocol.Response, beat func(Protocol.TobeSendRequest)) {
	s := "\n------RequesterHeartProtocol.Beat------>"
	req, opt := request.Request.(RequesterBeat), request.Option.(RequesterBeatSendOption)
	s += fmt.Sprintf("A beat %s,%s->%s was success. ", req.String(), opt.String(), response.String())
	req.n++
	if req.n < 10 {
		nextReq := Protocol.TobeSendRequest{Request: req, Option: opt}
		s += fmt.Sprintf("And the next beat is %s,%s",
			nextReq.Request.(RequesterBeat).String(),
			nextReq.Option.(RequesterBeatSendOption).String())
		beat(nextReq)
	}
}

/*
type ResponserHeartProtocol struct {
	responser *Heartbeat.Responser
}

func (r ResponserHeartProtocol) Response() (Protocol.Request, error, func(Protocol.TobeSendResponse)) {
	beat, err, resF := r.responser.Recv()
	request := beat.(RequesterBeat)
	return beat.(RequesterBeat)
}

func (r ResponserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {

}
*/
