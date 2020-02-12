package Heart

import (
	"gogistery/Heart"
	"gogistery/Protocol"
	"testing"
	"time"
)

func TestRequesterHeart(t *testing.T) {
	requester := Heart.NewRequesterHeart(&RequesterHeartProtocol{
		lastReq: Protocol.TobeSendRequest{
			Request: Request{ID: "0"},
			Option:  RequestSendOption{ID: "0", Addr: "0.0.0.0"}},
		lastTimeout: 10e9,
		lastRetryN:  10,
		n:           0},
		NewRequestBeatProtocol())
	err := requester.RunBeating(
		Protocol.TobeSendRequest{
			Request: Request{ID: "1"},
			Option:  RequestSendOption{ID: "1", Addr: "1.1.1.1"}},
		10e9, 10,
	)
	t.Log(err)
}

func TestResponserHeart(t *testing.T) {
	responser := Heart.NewResponserHeart(&ResponserHeartProtocol{
		lastRes: Protocol.TobeSendResponse{
			Response: Response{ID: "0"},
			Option:   ResponseSendOption{ID: "0"}},
	}, NewResponseBeatProtocol("0"))
	go responser.RunBeating()
	time.Sleep(1e9)
	responser.Stop()
	responser.Stop()
	go responser.RunBeating()
	time.Sleep(1e9)
	responser.Stop()
}
