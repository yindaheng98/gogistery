package protocol

import (
	"context"
	ChanNet2 "github.com/yindaheng98/gogistry/example/protocol/ChanNet"
	"github.com/yindaheng98/gogistry/protocol"
)

var ChanNet = ChanNet2.New(1e9, 10, "%02d.service.chanNet", 100)

type ChanNetRequestProtocol struct {
}

func NewChanNetRequestProtocol() ChanNetRequestProtocol {
	return ChanNetRequestProtocol{}
}

func (proto ChanNetRequestProtocol) Request(ctx context.Context, requestChan <-chan protocol.TobeSendRequest, responseChan chan<- protocol.ReceivedResponse) {
	select {
	case r := <-requestChan:
		request, option := r.Request, r.Option.(RequestSendOption)
		response, err := ChanNet.Request(ctx, option.RequestAddr, request)
		defer func() { recover() }()
		responseChan <- protocol.ReceivedResponse{Response: response, Error: err}
	case <-ctx.Done():
	}
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
func (proto ChanNetResponseProtocol) Response(ctx context.Context, requestChan chan<- protocol.ReceivedRequest, responseChan <-chan protocol.TobeSendResponse) {
	request, err, respChan := ChanNet.Response(ctx, proto.addr)
	defer func() { recover() }()
	select {
	case requestChan <- protocol.ReceivedRequest{Request: request, Error: err}:
		select {
		case response, ok := <-responseChan:
			if ok {
				respChan <- response.Response
			}
		case <-ctx.Done():
		}
	case <-ctx.Done():
	}
}
