package ChanNet

import (
	"context"
	"errors"
	"github.com/yindaheng98/gogistry/protocol"
)

type chanPair struct {
	ctx          context.Context
	requestChan  <-chan protocol.Request
	responseChan chan<- protocol.Response
}

type chanPairServer struct {
	processChan chan chanPair
}

func (s *chanPairServer) Request(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	requestChan := make(chan protocol.Request)
	responseChan := make(chan protocol.Response)
	s.processChan <- chanPair{ctx: ctx, requestChan: requestChan, responseChan: responseChan}
	requestChan <- request
	select {
	case response := <-responseChan:
		return response, nil
	case <-ctx.Done():
		return protocol.Response{}, errors.New("exited by context")
	}
}

func (s *chanPairServer) Response(ctx context.Context) (protocol.Request, error, chan<- protocol.Response) {
	pair := <-s.processChan
	requestChan, responseChan := pair.requestChan, pair.responseChan
	select {
	case request := <-requestChan:
		return request, nil, responseChan
	case <-pair.ctx.Done():
		return protocol.Request{}, errors.New("exited by requester's context"), responseChan
	case <-ctx.Done():
		return protocol.Request{}, errors.New("exited by context"), responseChan

	}
}
