package heart

import (
	"fmt"
	ExampleProtocol "gogistery/example/protocol"
	"gogistery/protocol"
	"time"
)

type RequesterHeartBeater struct {
	Info ExampleProtocol.RegistrantInfo
	n    int64
}

func NewRequesterHeartBeater(info ExampleProtocol.RegistrantInfo, BeatN int64) *RequesterHeartBeater {
	return &RequesterHeartBeater{Info: info, n: BeatN}
}
func (r *RequesterHeartBeater) Beat(response protocol.Response, timeout time.Duration, retryN uint64,
	beat func(protocol.TobeSendRequest, time.Duration, uint64)) {
	s := "------RequesterHeartProtocol.Beat------>"
	defer func() { fmt.Print(s + "\n") }()
	s += fmt.Sprintf("No.%d beat was success with a response %s. ", r.n, response.String())
	s += fmt.Sprintf("It takes %d times, at last time it takes %d", retryN, timeout)
	if r.n--; r.n < 0 {
		s += "And it's the end of beating."
		return
	}
	request := protocol.TobeSendRequest{
		Request: protocol.Request{
			RegistrantInfo: r.Info,
			Disconnect:     false,
		},
		Option: response.RegistryInfo.GetRequestSendOption(),
	}
	s += fmt.Sprintf("And the next beat is %s. ", request.String())
	beat(request, response.Timeout, (retryN+3)*2)
}

type ResponserHeartBeater struct {
	Info    ExampleProtocol.RegistryInfo
	Timeout time.Duration
	n       uint64
}

func NewResponserHeartBeater(info ExampleProtocol.RegistryInfo, Timeout time.Duration) *ResponserHeartBeater {
	return &ResponserHeartBeater{Info: info, Timeout: Timeout, n: 0}
}

func (r *ResponserHeartBeater) Beat(request protocol.Request) protocol.TobeSendResponse {
	s := "------ResponserHeartProtocol.Beat------>"
	defer func() { fmt.Print(s + "\n") }()
	s += fmt.Sprintf("No.%d request %s arrived. ", r.n, request.String())
	r.n++
	response := protocol.TobeSendResponse{
		Response: protocol.Response{
			RegistryInfo: r.Info,
			Timeout:      r.Timeout,
			Reject:       false,
		},
		Option: request.RegistrantInfo.GetResponseSendOption(),
	}
	s += fmt.Sprintf("And the response will be %s. ", response.String())
	return response
}
