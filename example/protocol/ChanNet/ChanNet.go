package ChanNet

import (
	"context"
	"errors"
	"fmt"
	"github.com/yindaheng98/gogistry/protocol"
	"math/rand"
	"time"
)

//A network simulator based on go channel
type ChanNet struct {

	//ChanNet can simulate the random latency in the network,
	//and this value defined the limitation of the latency
	MaxTimeout time.Duration

	//ChanNet can simulate the random sending failure in the network,
	//and this value defined the probability of the sending failure.
	//Probability of the sending failure=FailRate/100
	FailRate int

	//ChanNet can generate address for server, and AddrFormat is its format.
	//For example: "%02d.service.chanNet"
	AddrFormat string

	servers map[string]*chanPairServer
	src     rand.Source
}

//New Returns the pointer to a ChanNet
func New(MaxTimeout time.Duration, FailRate int, AddrFormat string, randSeed int64) *ChanNet {
	return &ChanNet{MaxTimeout: MaxTimeout, FailRate: FailRate, AddrFormat: AddrFormat,
		servers: make(map[string]*chanPairServer), src: rand.NewSource(randSeed)}
}

//NewServer creates a server and return its address
func (n *ChanNet) NewServer() string {
	addr := fmt.Sprintf(n.AddrFormat, len(n.servers))
	n.servers[addr] = &chanPairServer{processChan: make(chan chanPair)}
	return addr
}

//Send a request to a server and returns the response or error
func (n *ChanNet) Request(ctx context.Context, addr string, request protocol.Request) (protocol.Response, error) {
	s := "(ChanNet.Request)->"
	s += fmt.Sprintf("A request %s is transmitting to server in address '%s'. ", request.String(), addr)
	defer func() { fmt.Print(s + "\n") }()
	server, exists := n.servers[addr]
	if !exists {
		s += fmt.Sprintf("But the address '%s' is not exists. ", addr)
		return protocol.Response{}, errors.New("404 not found")
	}
	failN := rand.New(n.src).Intn(100)
	if failN <= n.FailRate {
		s += fmt.Sprintf("This transmition will fail. ")
		return protocol.Response{}, errors.New(fmt.Sprintf("send failed (failRate:%d,failN:%d)", n.FailRate, failN))
	}
	timeout := rand.New(n.src).Int63n(int64(n.MaxTimeout))
	s += fmt.Sprintf("This transmition will arrived in %f second. ", float64(timeout)/1e9)
	time.Sleep(time.Duration(timeout))
	return server.Request(ctx, request)
}

//Receive a request and send back response via a chan
func (n *ChanNet) Response(ctx context.Context, addr string) (protocol.Request, error, chan<- protocol.Response) {
	s := "(ChanNet.Response)->"
	server, exists := n.servers[addr]
	if !exists {
		return protocol.Request{}, errors.New(fmt.Sprintf("addr '%s' not exists", addr)), nil
	}
	request, err, responseChan := server.Response(ctx)
	if err != nil {
		return request, err, responseChan
	}
	s += fmt.Sprintf("A request from %s was arrived at server in address '%s'. ", request.RegistrantInfo.GetRegistrantID(), addr)
	failN := rand.New(n.src).Intn(100)
	if failN <= n.FailRate {
		defer func() { fmt.Print(s + "\n") }()
		s += fmt.Sprintf("This transmition failed. ")
		return protocol.Request{}, errors.New(fmt.Sprintf("recv failed (failRate:%d,failN:%d)", n.FailRate, failN)), nil
	}
	s += fmt.Sprintf("This transmition succeeded. ")
	responseTranChan := make(chan protocol.Response)
	go func() {
		defer func() { fmt.Print(s + "\n") }()
		select {
		case <-ctx.Done():
			s += "But canceled by context. "
		case response := <-responseTranChan:
			s += fmt.Sprintf("And the response is %s. ", response.String())
			responseChan <- response
		}
	}()
	return request, nil, responseTranChan
}
