package gogistery

import (
	"fmt"
	ExampleProtocol "github.com/yindaheng98/gogistry/example/protocol"
	"github.com/yindaheng98/gogistry/protocol"
	"github.com/yindaheng98/gogistry/util/CandidateList"
	"github.com/yindaheng98/gogistry/util/RetryNController"
	"github.com/yindaheng98/gogistry/util/TimeoutController"
	"testing"
	"time"
)

var RegistryInfos = make(map[string]ExampleProtocol.RegistryInfo)
var LastRegistryInfo protocol.RegistryInfo

func RegistryTest(t *testing.T) {
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
		t.Log(fmt.Sprintf("RegistryTest:%s--NewConnection--%s", info.GetRegistryID(), i.GetRegistrantID()))
	})
	r.Events.NewConnection.Enable()
	r.Events.Disconnection.AddHandler(func(i protocol.RegistrantInfo) {
		t.Log(fmt.Sprintf("RegistryTest:%s--Disconnection--%s", info.GetRegistryID(), i.GetRegistrantID()))
	})
	r.Events.Disconnection.Enable()
	go func() {
		r.Run()
		fmt.Printf("%s stopped itself.\n", info.ID)
	}()
	go func() {
		time.Sleep(15e9)
		r.Stop()
		t.Log(fmt.Sprintf("%s stopped manually.", info.ID))
	}()
}

const SERVERN = 5
const CLIENTN = 30

func RegistrantTest(t *testing.T, i int) {
	proto := ExampleProtocol.NewChanNetRequestProtocol()
	info := ExampleProtocol.RegistrantInfo{
		ID:     fmt.Sprintf("REGISTRANT_%02d", i),
		Option: ExampleProtocol.ResponseSendOption{},
	}
	r := NewRegistrant(info, 5,
		CandidateList.NewSimpleCandidateList(SERVERN, LastRegistryInfo, 2e9, 10),
		RetryNController.SimpleRetryNController{}, proto)
	r.Events.NewConnection.AddHandler(func(i protocol.RegistryInfo) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--NewConnection--%s", info.GetRegistrantID(), i.GetRegistryID()))
	})
	r.Events.NewConnection.Enable()
	r.Events.Disconnection.AddHandler(func(request protocol.TobeSendRequest, err error) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--Disconnection--%s. error:%s",
			info.GetRegistrantID(), request.Option.String(), err))
	})
	r.Events.Disconnection.Enable()
	r.Events.Error.AddHandler(func(err error) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--Error--%s", info.GetRegistrantID(), err))
	})
	r.Events.Error.Enable()
	go func() {
		r.Run()
		fmt.Printf("%s stopped itself.\n", info.ID)
	}()
	go func() {
		time.Sleep(10e9)
		r.Stop()
		t.Log(fmt.Sprintf("%s stopped manually.", info.ID))
	}()
}

func TestRegistryRegistrant(t *testing.T) {
	for i := 0; i < SERVERN; i++ {
		RegistryTest(t)
	}
	time.Sleep(1e9)
	for i := 0; i < CLIENTN; i++ {
		RegistrantTest(t, i)
	}
	time.Sleep(20e9)
}
