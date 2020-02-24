package protocol

import (
	ChanNet2 "github.com/yindaheng98/gogistry/example/protocol/ChanNet"
	"github.com/yindaheng98/gogistry/protocol"
)

var ChanNet = ChanNet2.New(1e9, 30, "%02d.service.chanNet", 100)

type ChanNetRequestProtocol struct {
}

func NewChanNetRequestProtocol() ChanNetRequestProtocol {
	return ChanNetRequestProtocol{}
}

func (proto ChanNetRequestProtocol) Request(requestChan <-chan protocol.TobeSendRequest, responseChan chan<- protocol.ReceivedResponse) {
	defer func() { recover() }()
	r := <-requestChan
	request, option := r.Request, r.Option.(RequestSendOption)
	response, err := ChanNet.Request(option.RequestAddr, request)
	responseChan <- protocol.ReceivedResponse{Response: response, Error: err}
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
func (proto ChanNetResponseProtocol) Response(requestChan chan<- protocol.ReceivedRequest, responseChan <-chan protocol.TobeSendResponse) {
	defer func() { recover() }()
	request, err, respChan := ChanNet.Response(proto.addr)
	requestChan <- protocol.ReceivedRequest{Request: request, Error: err}
	response, ok := <-responseChan
	if ok {
		respChan <- response.Response
	}
}
