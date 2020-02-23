package RegistryRegistrant

import (
	"fmt"
	ExampleProtocol "gogistery/example/protocol"
	"gogistery/protocol"
	"gogistery/registrant"
	"gogistery/registry"
	"gogistery/util/CandidateList"
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
	registry := registry.New(info, 5, NewTimeoutController(1e9, 3e9, 2), proto)
	registry.Events.NewConnection.AddHandler(func(i protocol.RegistrantInfo) {
		t.Log(fmt.Sprintf("RegistryTest:%s--NewConnection--%s", info.GetRegistryID(), i.GetRegistrantID()))
	})
	registry.Events.NewConnection.Enable()
	registry.Events.Disconnection.AddHandler(func(i protocol.RegistrantInfo) {
		t.Log(fmt.Sprintf("RegistryTest:%s--Disconnection--%s", info.GetRegistryID(), i.GetRegistrantID()))
	})
	registry.Events.Disconnection.Enable()
	go func() {
		go registry.Run()
		time.Sleep(15e9)
		registry.Stop()
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
	registrant := registrant.New(info, 5,
		CandidateList.NewSimpleCandidateList(LastRegistryInfo, SERVERN, 1e9, 3),
		RetryNController{}, proto)
	registrant.Events.NewConnection.AddHandler(func(i protocol.RegistryInfo) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--NewConnection--%s", info.GetRegistrantID(), i.GetRegistryID()))
	})
	registrant.Events.NewConnection.Enable()
	registrant.Events.Disconnection.AddHandler(func(request protocol.TobeSendRequest, err error) {
		t.Log(fmt.Sprintf("RegistrantTest:%s--Disconnection--%s. error:%s",
			info.GetRegistrantID(), request.Option.String(), err))
	})
	registrant.Events.Disconnection.Enable()
	go func() {
		registrant.Run()
		time.Sleep(18e9)
		registrant.Stop()
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
