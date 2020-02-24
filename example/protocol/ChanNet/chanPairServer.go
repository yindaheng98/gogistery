package ChanNet

import "github.com/yindaheng98/gogistry/protocol"

type chanPair struct {
	requestChan  <-chan protocol.Request
	responseChan chan<- protocol.Response
}

type chanPairServer struct {
	processChan chan chanPair
}

func (s *chanPairServer) Request(request protocol.Request) protocol.Response {
	requestChan := make(chan protocol.Request)
	responseChan := make(chan protocol.Response)
	s.processChan <- chanPair{requestChan: requestChan, responseChan: responseChan}
	requestChan <- request
	return <-responseChan
}

func (s *chanPairServer) Response() (protocol.Request, chan<- protocol.Response) {
	pair := <-s.processChan
	requestChan, responseChan := pair.requestChan, pair.responseChan
	return <-requestChan, responseChan
}
