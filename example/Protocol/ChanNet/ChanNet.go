package ChanNet

import (
	"errors"
	"fmt"
	"gogistery/Protocol"
	"math/rand"
	"time"
)

type ChanNet struct {
	MaxTimeout time.Duration
	FailRate   int
	AddrFormat string
	servers    map[string]*chanPairServer
	src        rand.Source
}

func New(MaxTimeout time.Duration, FailRate int, AddrFormat string, randSeed int64) *ChanNet {
	return &ChanNet{MaxTimeout: MaxTimeout, FailRate: FailRate, AddrFormat: AddrFormat,
		servers: make(map[string]*chanPairServer), src: rand.NewSource(randSeed)}
}

func (n *ChanNet) NewServer() string {
	addr := fmt.Sprintf(n.AddrFormat, len(n.servers))
	n.servers[addr] = &chanPairServer{processChan: make(chan chanPair)}
	return addr
}

func (n *ChanNet) Request(addr string, request Protocol.Request) (Protocol.Response, error) {
	s := "(ChanNet)->"
	s += fmt.Sprintf("A request %s is transmitting to server in address '%s'. ", request.String(), addr)
	defer func() { fmt.Print(s + "\n") }()
	server, exists := n.servers[addr]
	if !exists {
		s += fmt.Sprintf("But the address '%s' is not exists. ", addr)
		return Protocol.Response{}, errors.New("404 not found")
	}
	failN := rand.New(n.src).Intn(100)
	if failN <= n.FailRate {
		s += fmt.Sprintf("This transmition will fail. ")
		return Protocol.Response{}, errors.New(fmt.Sprintf("send failed (failRate:%d,failN:%d)", n.FailRate, failN))
	}
	timeout := rand.New(n.src).Int63n(int64(n.MaxTimeout))
	s += fmt.Sprintf("This transmition will arrived in %f second. ", float64(timeout)/1e9)
	time.Sleep(time.Duration(timeout))
	return server.Request(request), nil
}

func (n *ChanNet) Response(addr string) (Protocol.Request, error, chan<- Protocol.Response) {
	s := "(ChanNet)->"
	server, exists := n.servers[addr]
	if !exists {
		return Protocol.Request{}, errors.New(fmt.Sprintf("addr '%s' not exists", addr)), nil
	}
	request, responseChan := server.Response()
	s += fmt.Sprintf("A request %s was arrived at server in address '%s'", request.String(), addr)
	defer func() { fmt.Print(s + "\n") }()
	failN := rand.New(n.src).Intn(100)
	if failN <= n.FailRate {
		s += fmt.Sprintf("This transmition will fail. ")
		return Protocol.Request{}, errors.New(fmt.Sprintf("recv failed (failRate:%d,failN:%d)", n.FailRate, failN)), nil
	}
	s += fmt.Sprintf("This transmition will success. ")
	return request, nil, responseChan
}
