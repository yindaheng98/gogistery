package example

import (
	"fmt"
	"gogistery/Heart"
	"gogistery/Heartbeat"
	"gogistery/Protocol"
	"testing"
	"time"
)

func TestRequesterHeart(t *testing.T) {
	i := 1
	requester := Heart.NewRequesterHeart(
		RequesterHeartProtocol{
			Heartbeat.NewRequester(&RequestBeatProtocol{&src, 30, 0})})
	err := requester.RunBeating(
		Protocol.TobeSendRequest{
			Request: RequesterBeat{
				Request: Request{id: fmt.Sprintf("%d", i)},
				n:       0},
			Option: RequesterBeatSendOption{
				RequestSendOption: RequestSendOption{
					id:   fmt.Sprintf("%d", i),
					addr: fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)},
				timeout: time.Duration(5e8), /*********将该值调低可模拟超时情况**********/
				retryN:  10}})
	t.Log(err)
}

func TestResponserHeart(t *testing.T) {

}
