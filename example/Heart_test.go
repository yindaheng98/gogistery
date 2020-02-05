package example

import (
	"gogistery/Heart"
	"gogistery/Heartbeat"
	"gogistery/Protocol"
	ExampleHeartbeat "gogistery/example/Heartbeat"
	"testing"
	"time"
)

func TestRequesterHeart(t *testing.T) {
	requester := Heart.NewRequesterHeart(
		RequesterHeartProtocol{requester: Heartbeat.NewRequester(ExampleHeartbeat.NewRequestBeatProtocol())})
	err := requester.RunBeating(
		Heart.TobeSendRequest{
			Request: Protocol.TobeSendRequest{
				Request: ExampleHeartbeat.Request{ID: "0"},
				Option:  ExampleHeartbeat.RequestSendOption{ID: "0", Addr: "0.0.0.0"}},
			Option: RequestSendOption{"0", time.Duration(1e9), 10, 0}})
	t.Log(err)
}

func TestResponserHeart(t *testing.T) {
	responser := Heart.NewResponserHeart(
		ResponserHeartProtocol{responser: Heartbeat.NewResponser(ExampleHeartbeat.NewResponseBeatProtocol("0"))})
	go responser.RunBeating()
	time.Sleep(2e9)
	responser.Stop()
}
