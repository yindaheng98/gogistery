package ChanNet

import (
	"fmt"
	"github.com/yindaheng98/gogistry/protocol"
	"math/rand"
	"testing"
	"time"
)

type TestResponseSendOption struct {
	ID      string
	Encrypt string
}

func (o TestResponseSendOption) String() string {
	return fmt.Sprintf("TestResponseSendOption{ID:%s,Encrypt:%s}", o.ID, o.Encrypt)
}

type TestRegistrantInfo struct {
	ID     string
	Type   string
	Option TestResponseSendOption
}

func (info TestRegistrantInfo) GetRegistrantID() string {
	return info.ID
}
func (info TestRegistrantInfo) GetServiceType() string {
	return info.Type
}
func (info TestRegistrantInfo) GetResponseSendOption() protocol.ResponseSendOption {
	return info.Option
}
func (info TestRegistrantInfo) String() string {
	return fmt.Sprintf("TestRegistrantInfo{ID:%s,Type:%s,Option:%s}", info.ID, info.Type, info.Option.String())
}

type TestRequestSendOption struct {
	ID      string
	Encrypt string
}

func (o TestRequestSendOption) String() string {
	return fmt.Sprintf("TestRequestSendOption{ID:%s,Encrypt:%s}", o.ID, o.Encrypt)
}

type TestRegistryInfo struct {
	ID         string
	Type       string
	Option     TestRequestSendOption
	Candidates []protocol.RegistryInfo
}

func (info TestRegistryInfo) GetRegistryID() string {
	return info.ID
}
func (info TestRegistryInfo) GetServiceType() string {
	return info.Type
}
func (info TestRegistryInfo) GetRequestSendOption() protocol.RequestSendOption {
	return info.Option
}
func (info TestRegistryInfo) GetCandidates() []protocol.RegistryInfo {
	return info.Candidates
}
func (info TestRegistryInfo) String() string {
	Candidates := ""
	for _, RegistryInfo := range info.Candidates {
		Candidates += RegistryInfo.String() + ","
	}
	return fmt.Sprintf("TestRegistryInfo{ID:%s,Type:%s,Option:%s,Candidates:[%s]}",
		info.ID, info.Type, info.Option.String(), Candidates)
}

func RequestTest(t *testing.T, addr string, chanNet *ChanNet, i int) {
	request := protocol.Request{
		RegistrantInfo: TestRegistrantInfo{
			ID:   fmt.Sprintf("Registrant_%s", addr),
			Type: "REGISTRANT_TYPE_0",
			Option: TestResponseSendOption{
				ID:      fmt.Sprintf("RES_OPT_%02d", i),
				Encrypt: "AES_RESPONSE",
			},
		},
		Disconnect: false,
	}
	response, err := chanNet.Request(addr, request)
	s := "(RequestTest)"
	s += fmt.Sprintf("Request sent to '%s': %s, ", addr, request.String())
	if err != nil {
		t.Log(s + fmt.Sprintf("But an error occurred: %s.", err.Error()))
	} else {
		t.Log(s + fmt.Sprintf("And the response is %s.", response.String()))
	}
}

func ResponseTest(t *testing.T, addr string, chanNet *ChanNet) {
	request, err, responseChan := chanNet.Response(addr)
	s := "(ResponseTest)"
	if err != nil {
		t.Log(s + fmt.Sprintf("An error occurred: %s", err.Error()))
		return
	}
	s += fmt.Sprintf("Request arrived at '%s': %s, ", addr, request.String())
	response := protocol.Response{
		RegistryInfo: TestRegistryInfo{
			ID:   fmt.Sprintf("Registry_%s", addr),
			Type: "REGISTRY_TYPE_0",
			Option: TestRequestSendOption{
				ID:      "REQ_OPT_TO_" + request.RegistrantInfo.GetResponseSendOption().(TestResponseSendOption).ID,
				Encrypt: "AES_REQUEST",
			},
			Candidates: nil,
		},
		Timeout: 0,
		Reject:  false,
	}
	t.Log(s + fmt.Sprintf("And the response will be %s", response.String()))
	responseChan <- response
}

const TESTADDRN = 10
const TESTREQN = 30

func TestChanNet(t *testing.T) {
	chanNet := New(1e9, 30, "%02d.service.chanNet", 100)
	addrs := make([]string, TESTADDRN)
	reqN := make([]int64, TESTADDRN)
	for i := 0; i < TESTADDRN; i++ { //新建指定数量的chanPair服务器
		addrs[i] = chanNet.NewServer()
		reqN[i] = 0
	}
	for i := 0; i < TESTREQN; i++ { //发送指定数量个请求
		reqi := rand.Intn(TESTADDRN) //随机选择向谁发
		reqN[reqi] += 1              //记录请求发送次数
		go func(testi int) {
			RequestTest(t, addrs[reqi], chanNet, testi)
		}(i)
	}
	for i := 0; i < TESTADDRN; i++ { //每个服务器都执行响应操作
		for j := reqN[i]; j > 0; j-- {
			go func(addri int) {
				ResponseTest(t, addrs[addri], chanNet)
			}(i)
		}
	}
	time.Sleep(1e9)
}
