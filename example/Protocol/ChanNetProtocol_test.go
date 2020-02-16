package Protocol

import (
	"fmt"
	"gogistery/Protocol"
	"math/rand"
	"testing"
	"time"
)

func RequestTest(t *testing.T, proto ChanNetRequestProtocol, RegistrantID string, addr string) {
	s := "(RequestTest)->"
	defer func() { t.Log(s) }()
	requestChan := make(chan Protocol.TobeSendRequest, 1)
	responseChan := make(chan Protocol.ReceivedResponse, 1)
	req := Protocol.TobeSendRequest{
		Request: Protocol.Request{
			RegistrantInfo: RegistrantInfo{ID: RegistrantID, Option: ResponseSendOption{Timestamp: time.Now()}},
			Disconnect:     false,
		},
		Option: RequestSendOption{RequestAddr: addr, Timestamp: time.Now()},
	}
	requestChan <- req
	go proto.Request(requestChan, responseChan)
	s += fmt.Sprintf("Request sent: %s. ", req.Request.String())
	select {
	case <-time.After(1e9):
		s += fmt.Sprintf("But timeout. ")
	case res := <-responseChan:
		response, err := res.Response, res.Error
		if err != nil {
			s += fmt.Sprintf("But an error occurred: %s", err.Error())
		} else {
			s += fmt.Sprintf("And the response is %s", response.String())
		}
	}
}

func ResponseTest(t *testing.T, proto ChanNetResponseProtocol) {
	s := "(ResponseTest)->"
	defer func() { t.Log(s) }()
	requestChan := make(chan Protocol.ReceivedRequest, 1)
	responseChan := make(chan Protocol.TobeSendResponse, 1)
	go proto.Response(requestChan, responseChan)
	req := <-requestChan
	s += fmt.Sprintf("A request arrived at '%s'. ", proto.GetAddr())
	request, err := req.Request, req.Error
	if err != nil {
		s += fmt.Sprintf("But an error occurred: %s. ", err.Error())
		return
	}
	s += fmt.Sprintf("It is %s. ", request.String())
	response := Protocol.TobeSendResponse{
		Response: Protocol.Response{
			RegistryInfo: RegistryInfo{
				ID: proto.GetAddr(),
				Option: RequestSendOption{
					RequestAddr: proto.GetAddr(),
					Timestamp:   time.Now(),
				},
				Candidates: []Protocol.RegistryInfo{},
			},
			Timeout: 0,
			RetryN:  0,
			Reject:  false,
		},
		Option: request.RegistrantInfo.GetResponseSendOption(),
	}
	s += fmt.Sprintf("And the response will be %s. ", response.String())
	responseChan <- response
}

const SERVERN = 10
const TESTN = 30

func TestChanNetRequestProtocol(t *testing.T) {
	servers := make([]ChanNetResponseProtocol, SERVERN)
	servern := make([]int, SERVERN)
	for i := 0; i < SERVERN; i++ {
		servers[i] = NewChanNetResponseProtocol()
		servern[i] = 0
	}
	for i := 0; i < TESTN; i++ {
		si := rand.Intn(SERVERN)
		servern[si] += 1
		proto := NewChanNetRequestProtocol()
		go func(i int) {
			RequestTest(t, proto, fmt.Sprintf("REGISTRANT_%02d", i), servers[si].GetAddr())
		}(i)
	}
	for i := 0; i < SERVERN; i++ {
		proto := servers[i]
		for j := 0; j < servern[i]; j++ {
			go func() {
				ResponseTest(t, proto)
			}()
		}
	}
	time.Sleep(2e9)
}
