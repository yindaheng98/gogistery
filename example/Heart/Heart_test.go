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
			Option:  RequestSendOption{Timeout: 1e9, RetryN: 10, ID: "0", Addr: "0.0.0.0"}},
		n: 0},
		NewRequestBeatProtocol())
	err := requester.RunBeating(
		Protocol.TobeSendRequest{
			Request: Request{ID: "1"},
			Option:  RequestSendOption{Timeout: 1e10, RetryN: 10, ID: "1", Addr: "1.1.1.1"}},
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
