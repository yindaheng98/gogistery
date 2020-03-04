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

func (s *chanPairServer) Request(ctx context.Context, request protocol.Request) (r protocol.Response, e error) {
	requestChan := make(chan protocol.Request, 1)
	responseChan := make(chan protocol.Response, 1)
	e = errors.New("exited by context")
	select {
	case s.processChan <- chanPair{ctx: ctx, requestChan: requestChan, responseChan: responseChan}:
	case <-ctx.Done():
		return
	}
	requestChan <- request
	select {
	case response := <-responseChan:
		return response, nil
	case <-ctx.Done():
		return
	}
}

func (s *chanPairServer) Response(ctx context.Context) (r protocol.Request, e error, c chan<- protocol.Response) {
	var pair chanPair
	e = errors.New("exited by context")
	select {
	case pair = <-s.processChan:
	case <-ctx.Done():
		return
	}
	requestChan, responseChan := pair.requestChan, pair.responseChan
	select {
	case request := <-requestChan:
		return request, nil, responseChan
	case <-pair.ctx.Done():
		e = errors.New("exited by requester's context")
		return
	case <-ctx.Done():
		return
	}
}
