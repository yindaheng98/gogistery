package gogistery

import (
	"context"
	"fmt"
	"github.com/yindaheng98/gogistry/example/CandidateList"
	"github.com/yindaheng98/gogistry/example/RetryNController"
	"github.com/yindaheng98/gogistry/example/TimeoutController"
	ExampleProtocol "github.com/yindaheng98/gogistry/example/protocol"
	"github.com/yindaheng98/gogistry/protocol"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var RegistryInfos = make(map[string]ExampleProtocol.RegistryInfo)
var LastRegistryInfo protocol.RegistryInfo

func RegistryTest(t *testing.T, ctx context.Context, wg *sync.WaitGroup) {
	proto := ExampleProtocol.NewChanNetResponseProtocol()
	info := ExampleProtocol.RegistryInfo{
		ID: "REGISTRY_" + proto.GetAddr(),
		Option: ExampleProtocol.RequestSendOption{
			RequestAddr: proto.GetAddr(),
			Timestamp:   time.Now(),
		},
		Candidates: nil,
	}
	for _, RegistryInfo := range RegistryInfos {
		info.Candidates = append(info.Candidates, RegistryInfo)
	}
	RegistryInfos[proto.GetAddr()] = info
	LastRegistryInfo = info
	r := NewRegistry(info, 5,
		TimeoutController.NewLogTimeoutController(1e9, 3e9, 2),
		proto)
	r.Events.NewConnection.AddHandler(func(i protocol.RegistrantInfo) {
		t.Log(fmt.Sprintf("RegistryTest(connections: %d) %s--NewConnection--%s", len(r.GetConnections()), info.GetRegistryID(), i.GetRegistrantID()))
	})
	r.Events.NewConnection.Enable()
	r.Events.Disconnection.AddHandler(func(i protocol.RegistrantInfo) {
		t.Log(fmt.Sprintf("RegistryTest:%s--Disconnection--%s", info.GetRegistryID(), i.GetRegistrantID()))
	})
	r.Events.Disconnection.Enable()
	go func() {
		defer wg.Done()
		r.Run(ctx)
		fmt.Printf("%s stopped itself.\n", info.ID)
	}()
}

const SERVERN = 5
const CLIENTN = 30

type TestPINGer struct {
	failRate uint8
	src      rand.Source
	maxT     time.Duration
}

func NewTestPINGer(failRate uint8, maxT time.Duration) *TestPINGer {
	return &TestPINGer{failRate, rand.NewSource(10), maxT}
}

func (p *TestPINGer) PING(ctx context.Context, info protocol.RegistryInfo) bool {
	s := fmt.Sprintf("TestPINGer.PING(%s)-->", info.String())
	r := rand.New(p.src).Int31n(100)
	timeout := time.Duration(rand.New(p.src).Uint64()) % p.maxT
	s += fmt.Sprintf("This PING will return in %d. ", timeout)
	if uint8(r) < p.failRate {
		s += fmt.Sprintf("But it was failed(failRate:%d,r:%d).", p.failRate, r)
	} else {
		s += "And it succeed"
	}
	fmt.Println(s)
	select {
	case <-ctx.Done():
		return false
	case <-time.After(timeout):
	}
	time.Sleep(timeout)
	return uint8(r) >= p.failRate
}

func RegistrantTest(t *testing.T, ctx context.Context, i int, wg *sync.WaitGroup) {
	proto := ExampleProtocol.NewChanNetRequestProtocol()
	info := ExampleProtocol.RegistrantInfo{
		ID: fmt.Sprintf("REGISTRANT_%02d", i),
		//Type:   "XXX", //模拟类型不一样时的连接拒绝过程
		Option: ExampleProtocol.ResponseSendOption{},
	}
	r := NewRegistrant(info, 5,
		//CandidateList.NewSimpleCandidateList(SERVERN, LastRegistryInfo, 2e9, 10),
		CandidateList.NewPingerCandidateList(SERVERN, NewTestPINGer(30, 1e9), 1e9, LastRegistryInfo, 2e9, 10),
		RetryNController.NewLinearRetryNController(), proto)
	r.Events.NewConnection.AddHandler(func(i protocol.RegistryInfo) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--NewConnection--%s", info.GetRegistrantID(), i.GetRegistryID()))
	})
	r.Events.NewConnection.Enable()
	r.Events.Disconnection.AddHandler(func(i protocol.RegistryInfo, err error) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--Disconnection--%s. error:%s",
			info.GetRegistrantID(), i.GetRegistryID(), err))
	})
	r.Events.Disconnection.Enable()
	r.Events.Error.AddHandler(func(err error) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--Error--%s", info.GetRegistrantID(), err))
	})
	r.Events.Error.Enable()
	go func() {
		defer wg.Done()
		r.Run(ctx)
		fmt.Printf("%s stopped itself.\n", info.ID)
	}()
}

func TestRegistryRegistrant(t *testing.T) {
	wgRegistry := new(sync.WaitGroup)
	wgRegistry.Add(SERVERN)
	ctxRegistry, cancelRegistry := context.WithTimeout(context.Background(), 10e9)
	for i := 0; i < SERVERN; i++ {
		RegistryTest(t, ctxRegistry, wgRegistry)
	}
	time.Sleep(1e9)
	wgRegistrant := new(sync.WaitGroup)
	wgRegistrant.Add(CLIENTN)
	ctxRegistrant, cancelRegistrant := context.WithTimeout(context.Background(), 10e9)
	for i := 0; i < CLIENTN; i++ {
		RegistrantTest(t, ctxRegistrant, i, wgRegistrant)
	}
	time.Sleep(20e9)
	cancelRegistry()
	cancelRegistrant()
	wgRegistry.Wait()
	wgRegistrant.Wait()
}
