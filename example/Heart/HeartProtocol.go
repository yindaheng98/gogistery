package Heart

import (
	"fmt"
	"gogistery/Protocol"
	ExampleProtocol "gogistery/example/Protocol"
	"time"
)

type RequesterHeartProtocol struct {
	Info ExampleProtocol.RegistrantInfo
	n    int64
}

func NewRequesterHeartProtocol(info ExampleProtocol.RegistrantInfo, BeatN int64) *RequesterHeartProtocol {
	return &RequesterHeartProtocol{Info: info, n: BeatN}
}
func (r *RequesterHeartProtocol) Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest, time.Duration, uint64)) {
	s := "------RequesterHeartProtocol.Beat------>"
	defer func() { fmt.Print(s + "\n") }()
	s += fmt.Sprintf("No.%d beat was success with a response %s. ", r.n, response.String())
	if r.n--; r.n < 0 {
		s += "And it's the end of beating."
		return
	}
	request := Protocol.TobeSendRequest{
		Request: Protocol.Request{
			RegistrantInfo: r.Info,
			Disconnect:     false,
		},
		Option: response.RegistryInfo.GetRequestSendOption(),
	}
	s += fmt.Sprintf("And the next beat is %s. ", request.String())
	beat(request, response.Timeout, response.RetryN)
}

type ResponserHeartProtocol struct {
	Info    ExampleProtocol.RegistryInfo
	Timeout time.Duration
	RetryN  uint64
	n       uint64
}

func NewResponserHeartProtocol(info ExampleProtocol.RegistryInfo, Timeout time.Duration, RetryN uint64) *ResponserHeartProtocol {
	return &ResponserHeartProtocol{Info: info, Timeout: Timeout, RetryN: RetryN, n: 0}
}

func (r *ResponserHeartProtocol) Beat(request Protocol.Request) Protocol.TobeSendResponse {
	s := "------ResponserHeartProtocol.Beat------>"
	defer func() { fmt.Print(s + "\n") }()
	s += fmt.Sprintf("No.%d request %s arrived. ", r.n, request.String())
	r.n++
	response := Protocol.TobeSendResponse{
		Response: Protocol.Response{
			RegistryInfo: r.Info,
			Timeout:      r.Timeout,
			RetryN:       r.RetryN,
			Reject:       false,
		},
		Option: request.RegistrantInfo.GetResponseSendOption(),
	}
	s += fmt.Sprintf("And the response will be %s. ", response.String())
	return response
}
