package Heart

import (
	"gogistery/Heart"
	"gogistery/Protocol"
	ExampleHeartbeat "gogistery/example/Heartbeat"
	"testing"
	"time"
)

func TestRequesterHeart(t *testing.T) {
	requester := Heart.NewRequesterHeart(&RequesterHeartProtocol{
		lastReq: Protocol.TobeSendRequest{
			Request: ExampleHeartbeat.Request{ID: "0"},
			Option:  ExampleHeartbeat.RequestSendOption{ID: "0", Addr: "0.0.0.0"}},
		lastOpt: Heart.RequestSendOption{Timeout: 1e9, RetryN: 10},
		n:       0},
		ExampleHeartbeat.NewRequestBeatProtocol())
	err := requester.RunBeating(
		Protocol.TobeSendRequest{
			Request: ExampleHeartbeat.Request{ID: "1"},
			Option:  ExampleHeartbeat.RequestSendOption{ID: "1", Addr: "1.1.1.1"}},
		Heart.RequestSendOption{Timeout: 1e9, RetryN: 10})
	t.Log(err)
}

func TestResponserHeart(t *testing.T) {
	responser := Heart.NewResponserHeart(&ResponserHeartProtocol{
		lastRes: Protocol.TobeSendResponse{
			Response: ExampleHeartbeat.Response{ID: "0"},
			Option:   ExampleHeartbeat.ResponseSendOption{ID: "0"}},
	}, ExampleHeartbeat.NewResponseBeatProtocol("0"))
	go responser.RunBeating()
	time.Sleep(1e9)
	responser.Stop()
}
