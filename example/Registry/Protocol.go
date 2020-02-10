package Registry

import (
	"gogistery/Protocol"
	"gogistery/Registry"
	"math"
	"time"
)

type TimeoutProtocol struct {
	minT time.Duration //最小Timeout
	maxT time.Duration //最大Timeout
	cT   float64       //从最小到最大的增长系数
	tMap map[string]time.Duration
}

func NewTimeoutProtocol() *TimeoutProtocol {
	return &TimeoutProtocol{10e9, 100e9, 10, make(map[string]time.Duration)}
}

func (p *TimeoutProtocol) TimeoutForNew(request Registry.Request) time.Duration {
	p.tMap[request.GetRegistrantID()] = p.minT
	return p.minT
}

func (p *TimeoutProtocol) TimeoutForUpdate(request Registry.Request) time.Duration {
	t := p.tMap[request.GetRegistrantID()]
	t += time.Duration(math.Floor(float64(p.maxT-t) / p.cT))
	p.tMap[request.GetRegistrantID()] = t
	return t
}

type beatChanPair struct {
	requestChan  <-chan Protocol.ReceivedRequest
	responseChan chan<- Protocol.TobeSendResponse
}

type ResponseBeatProtocol struct {
	beatChanPairs chan beatChanPair
}

func NewResponseBeatProtocol() *ResponseBeatProtocol {
	return &ResponseBeatProtocol{make(chan beatChanPair)}
}

func (r *ResponseBeatProtocol) Response(requestChan chan<- Protocol.ReceivedRequest, responseChan <-chan Protocol.TobeSendResponse) {
	pair := <-r.beatChanPairs
	request := <-pair.requestChan
	requestChan <- request
	response := <-responseChan
	pair.responseChan <- response
}
