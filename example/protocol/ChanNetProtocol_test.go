package protocol

import (
	"context"
	"fmt"
	"github.com/yindaheng98/gogistry/protocol"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func RequestTest(t *testing.T, timeoutCtx context.Context, proto ChanNetRequestProtocol, RegistrantID string, addr string) {
	s := "(RequestTest)->"
	defer func() { t.Log(s) }()
	requestChan := make(chan protocol.TobeSendRequest, 1)
	responseChan := make(chan protocol.ReceivedResponse, 1)
	req := protocol.TobeSendRequest{
		Request: protocol.Request{
			RegistrantInfo: RegistrantInfo{
				ID:     RegistrantID,
				Type:   "REGISTRANT_TYPE_0",
				Option: ResponseSendOption{Timestamp: time.Now()}},
			Disconnect: false,
		},
		Option: RequestSendOption{RequestAddr: addr, Timestamp: time.Now()},
	}
	requestChan <- req
	go proto.Request(timeoutCtx, requestChan, responseChan)
	s += fmt.Sprintf("Request sent: %s. ", req.Request.String())
	select {
	case <-timeoutCtx.Done():
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

func ResponseTest(t *testing.T, ctx context.Context, proto ChanNetResponseProtocol) {
	s := "(ResponseTest)->"
	defer func() { t.Log(s) }()
	requestChan := make(chan protocol.ReceivedRequest, 1)
	responseChan := make(chan protocol.TobeSendResponse, 1)
	go proto.Response(ctx, requestChan, responseChan)
	select {
	case <-ctx.Done():
		s += fmt.Sprintf("Exited by context. ")
	case req := <-requestChan:
		s += fmt.Sprintf("A request arrived at '%s'. ", proto.GetAddr())
		request, err := req.Request, req.Error
		if err != nil {
			s += fmt.Sprintf("But an error occurred: %s. ", err.Error())
			return
		}
		s += fmt.Sprintf("It is %s. ", request.String())
		response := protocol.TobeSendResponse{
			Response: protocol.Response{
				RegistryInfo: RegistryInfo{
					ID:   proto.GetAddr(),
					Type: "REGISTRY_TYPE_0",
					Option: RequestSendOption{
						RequestAddr: proto.GetAddr(),
						Timestamp:   time.Now(),
					},
					Candidates: []protocol.RegistryInfo{},
				},
				Timeout: 0,
				Reject:  false,
			},
			Option: request.RegistrantInfo.GetResponseSendOption(),
		}
		s += fmt.Sprintf("And the response will be %s. ", response.String())
		responseChan <- response
	}
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
	wg := new(sync.WaitGroup)
	wg.Add(TESTN * 2)
	ctx, cancel := context.WithTimeout(context.Background(), 1e9)
	for i := 0; i < TESTN; i++ {
		si := rand.Intn(SERVERN)
		servern[si] += 1
		proto := NewChanNetRequestProtocol()
		go func(i int) {
			defer wg.Done()
			timeoutCtx, cancel := context.WithTimeout(ctx, 1e9)
			defer cancel()
			RequestTest(t, timeoutCtx, proto, fmt.Sprintf("REGISTRANT_%02d", i), servers[si].GetAddr())
		}(i)
	}
	for i := 0; i < SERVERN; i++ {
		proto := servers[i]
		for j := 0; j < servern[i]; j++ {
			go func() {
				defer wg.Done()
				timeoutCtx, cancel := context.WithTimeout(ctx, 1e9)
				defer cancel()
				ResponseTest(t, timeoutCtx, proto)
			}()
		}
	}
	time.Sleep(1e8)
	//cancel()
	wg.Wait()
	cancel()
}
