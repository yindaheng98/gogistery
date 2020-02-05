package example

import (
	"fmt"
	"gogistery/Heart"
	"gogistery/Heartbeat"
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

func (r ResponserBeat) Print() string {
	return fmt.Sprintf("ResponserBeat{%s,n:%d}", r.Response.Print(), r.n)
}
func (r RequesterBeat) Print() string {
	return fmt.Sprintf("RequesterBeat{%s,n:%d}", r.Request.Print(), r.n)
}

type RequesterBeatSendOption struct {
	RequestSendOption
	timeout time.Duration
	retryN  int64
}
type ResponserBeatSendOption struct {
	ResponseSendOption
}

func (r RequesterBeatSendOption) Print() string {
	return fmt.Sprintf("ResponserBeat{%s,timeout:%d,retryN:%d}", r.RequestSendOption.Print(), r.timeout, r.retryN)
}
func (r ResponserBeatSendOption) Print() string {
	return fmt.Sprintf("RequesterBeat{%s}", r.ResponseSendOption.Print())
}

type RequesterHeartProtocol struct {
	requester *Heartbeat.Requester
}

func (r RequesterHeartProtocol) Request(beat Heart.TobeSendRequesterBeat) (Heart.ResponserBeat, error) {
	s := "\n------RequesterHeartProtocol.Request------>"
	b := beat.RequesterBeat.(RequesterBeat)
	o := beat.SendOption.(RequesterBeatSendOption)
	s += fmt.Sprintf("A %s is sending with %s. ", b.Print(), o.Print())
	response, err := r.requester.Send(Heartbeat.TobeSendRequest{Request: b.Request, Option: o.RequestSendOption}, o.timeout, o.retryN)
	if err != nil {
		fmt.Print(s + fmt.Sprintf("Send failed, and the error is: %s", err.Error()))
		return nil, err
	}
	fmt.Print(s + fmt.Sprintf("Send success, and the response is: %s", response.(Response).Print()))
	return ResponserBeat{response.(Response), b.n}, nil
}
func (r RequesterHeartProtocol) Beat(request Heart.TobeSendRequesterBeat, response Heart.ResponserBeat, beat func(Heart.TobeSendRequesterBeat)) {
	s := "\n------RequesterHeartProtocol.Beat------>"
	s += fmt.Sprintf("A beat %s,%s->%s was success. ",
		request.RequesterBeat.(RequesterBeat).Print(),
		request.SendOption.(RequesterBeatSendOption).Print(),
		response.(ResponserBeat).Print())
	req := request.RequesterBeat.(RequesterBeat)
	req.n++
	if req.n < 10 {
		nextReq := Heart.TobeSendRequesterBeat{RequesterBeat: req, SendOption: request.SendOption}
		s += fmt.Sprintf("And the next beat is %s,%s",
			nextReq.RequesterBeat.(RequesterBeat).Print(),
			nextReq.SendOption.(RequesterBeatSendOption).Print())
		beat(nextReq)
	}
}

/*
type ResponserHeartProtocol struct {
}

func (r ResponserHeartProtocol) Response() (Heart.RequesterBeat, error, func(Heart.TobeSendResponserBeat)) {

}

func (r ResponserHeartProtocol) Beat(request Heart.RequesterBeat) Heart.TobeSendResponserBeat {

}
*/
