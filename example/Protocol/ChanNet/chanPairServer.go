package ChanNet

import "gogistery/Protocol"

type chanPair struct {
	requestChan  <-chan Protocol.Request
	responseChan chan<- Protocol.Response
}

type chanPairServer struct {
	processChan chan chanPair
}

func (s *chanPairServer) Request(request Protocol.Request) Protocol.Response {
	requestChan := make(chan Protocol.Request)
	responseChan := make(chan Protocol.Response)
	s.processChan <- chanPair{requestChan: requestChan, responseChan: responseChan}
	requestChan <- request
	return <-responseChan
}

func (s *chanPairServer) Response() (Protocol.Request, chan<- Protocol.Response) {
	pair := <-s.processChan
	requestChan, responseChan := pair.requestChan, pair.responseChan
	return <-requestChan, responseChan
}
