package Protocol

import (
	"gogistery/Protocol"
	ChanNet2 "gogistery/example/Protocol/ChanNet"
)

var ChanNet = ChanNet2.New(1e9, 30, "%02d.service.chanNet", 100)

type ChanNetRequestProtocol struct {
}

func NewChanNetRequestProtocol() ChanNetRequestProtocol {
	return ChanNetRequestProtocol{}
}

func (proto ChanNetRequestProtocol) Request(requestChan <-chan Protocol.TobeSendRequest, responseChan chan<- Protocol.ReceivedResponse) {
	r := <-requestChan
	request, option := r.Request, r.Option.(RequestSendOption)
	response, err := ChanNet.Request(option.RequestAddr, request)
	responseChan <- Protocol.ReceivedResponse{Response: response, Error: err}
}

type ChanNetResponseProtocol struct {
	addr string
}

func NewChanNetResponseProtocol() ChanNetResponseProtocol {
	return ChanNetResponseProtocol{ChanNet.NewServer()}
}
func (proto ChanNetResponseProtocol) GetAddr() string {
	return proto.addr
}
func (proto ChanNetResponseProtocol) Response(requestChan chan<- Protocol.ReceivedRequest, responseChan <-chan Protocol.TobeSendResponse) {
	request, err, respChan := ChanNet.Response(proto.addr)
	requestChan <- Protocol.ReceivedRequest{Request: request, Error: err}
	response, ok := <-responseChan
	if ok {
		respChan <- response.Response
	}
}
