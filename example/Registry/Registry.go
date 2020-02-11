package Registry

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"gogistery/Protocol"
	"gogistery/Registry"
	"time"
)

type RegistrantInfo struct {
	ID   string
	Addr string
}

func (info RegistrantInfo) GetRegistrantID() string {
	return info.ID
}

type Request struct {
	RegistrantInfo
	toDisconnect bool
	RequestTime  time.Time
	LastRetryN   uint64
}

func (r Request) ToDisconnect() bool {
	return r.toDisconnect
}
func (r Request) String() string {
	return fmt.Sprintf("Request{RegistrantInfo{id:%s,addr:%s},LastRetryN:%d}", r.ID, r.Addr, r.LastRetryN)
}

type Info struct {
	ID   string
	Addr string
}

func (info Info) String() string {
	return fmt.Sprintf("RegistryInfo{ID:%s,Addr:%s}", info.ID, info.Addr)
}

func RegistryHandler(ctx iris.Context, beatProto *ResponseBeatProtocol) {
	request := Request{
		RegistrantInfo: RegistrantInfo{
			ID:   ctx.URLParam("RegistrantID"),
			Addr: ctx.URLParam("RegistrantAddr"),
		},
		LastRetryN:  uint64(ctx.URLParamInt64Default("LastRetryN", 0)),
		RequestTime: time.Unix(ctx.URLParamInt64Default("RequestTime", 0), 0),
	}
	requestChan := make(chan Protocol.ReceivedRequest, 1)
	responseChan := make(chan Protocol.TobeSendResponse, 1)
	beatProto.beatChanPairs <- beatChanPair{requestChan, responseChan}
	requestChan <- Protocol.ReceivedRequest{Request: request}
	resp := <-responseChan
	response, option := resp.Response.(Registry.Response), resp.Option
	fmt.Print(option.String())
	ctx.WriteString(response.String())
	ctx.EndRequest()
}
