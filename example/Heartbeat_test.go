package example

import (
	"fmt"
	"gogistery/Heartbeat"
	"math/rand"
	"testing"
	"time"
)

var src = rand.NewSource(10)

func testReq(i uint64, logger func(string)) {
	s := "------TestMyRequester------>"
	requester := Heartbeat.NewRequester(&MyRequestProtocol{&src, 30, 0})
	requester.Events.Retry.AddHandler(func(o Heartbeat.ProtocolRequestSendOption, err error) {
		logger(s + fmt.Sprintf("An retry was occured. error: %s", err.Error()))
	})
	requester.Events.Retry.Enable()
	response, err := requester.Send(Heartbeat.ProtocolRequestSendOption{
		Request: MyRequest{fmt.Sprintf("%02d", i)},
		Option: MyRequestOption{
			fmt.Sprintf("%02d", i),
			fmt.Sprintf("%02d.%02d.%02d.%02d", i, i, i, i)}},
		time.Duration(1e6), /*********将该值调低可模拟超时情况**********/
		10)
	if err != nil {
		logger(s + fmt.Sprintf("No.%02d test failed. err is %s", i, err.Error()))
		return
	}
	logger(s + fmt.Sprintf("No.%02d sending test succeed. response is MyResponse{id:%s}", i, response.(MyResponse).id))
}

//单次Heartbeat
func TestMyRequester(t *testing.T) {
	for i := uint64(0); i < 30; i++ {
		testReq(i, func(s string) {
			t.Log(s)
		})
	}
	time.Sleep(1000)
}

func testRes(i uint64, logger func(string)) {
	s := "------TestMyResponser------>"
	responser := Heartbeat.NewResponser(MyResponseProtocol{&src, 30, fmt.Sprintf("%d", i)})
	request, err, responseFunc := responser.Recv()
	d := time.Duration(rand.Int31n(1e3) * 1e3)
	if err != nil {
		logger(s + err.Error())
		time.Sleep(d)
		responseFunc(Heartbeat.ProtocolResponseSendOption{Response: MyResponse{fmt.Sprintf("error%02d", i)},
			Option: MyResponseOption{fmt.Sprintf("error%02d", i)}})
	} else {
		logger(s + fmt.Sprintf("A request MyRequest{id:%s} arrived. Response will be sent back in %d",
			request.(MyRequest).id, d))
		time.Sleep(d)
		responseFunc(Heartbeat.ProtocolResponseSendOption{Response: MyResponse{fmt.Sprintf("%02d", i)},
			Option: MyResponseOption{fmt.Sprintf("%02d", i)}})
	}
}

//单次Heartbeat
func TestMyResponser(t *testing.T) {
	for i := uint64(0); i < 30; i++ {
		testRes(i, func(s string) {
			t.Log(s)
		})
	}
}
