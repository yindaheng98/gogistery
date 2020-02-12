package Heart

import (
	"fmt"
	"gogistery/Protocol"
	"time"
)

type RequesterHeartProtocol struct {
	lastReq     Protocol.TobeSendRequest
	lastTimeout time.Duration
	lastRetryN  uint64
	n           uint64
}

func (r *RequesterHeartProtocol) Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest, time.Duration, uint64)) {
	s := "\n------RequesterHeartProtocol.Beat------>"
	request := r.lastReq
	s += fmt.Sprintf("No.%d beat was success with a response %s. ", r.n, response.String())
	r.n++
	if r.n < 10 {
		fmt.Print(s + fmt.Sprintf("And the next beat is %s. ", request.String()))
		beat(request, r.lastTimeout, r.lastRetryN)
	} else {
		fmt.Print(s + "And it's the end of beating.")
	}
}

type ResponserHeartProtocol struct {
	lastRes Protocol.TobeSendResponse
	n       uint64
}

func (r *ResponserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	s := "\n------ResponserHeartProtocol.Beat------>"
	s += fmt.Sprintf("No.%d request %s arrived. ", r.n, request.String())
	r.n++
	res := r.lastRes
	fmt.Print(s + fmt.Sprintf("And the response will be %s with option %s. ",
		res.Response.String(), res.Option.String()))
	return res
}
