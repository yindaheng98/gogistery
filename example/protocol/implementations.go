package protocol

import (
	"fmt"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type ResponseSendOption struct {
	Timestamp time.Time
}

func (o ResponseSendOption) String() string {
	return fmt.Sprintf("ResponseSendOption{Timestamp:%s}", o.Timestamp)
}

type RegistrantInfo struct {
	ID     string
	Option ResponseSendOption
}

func (info RegistrantInfo) GetRegistrantID() string {
	return info.ID
}
func (info RegistrantInfo) GetResponseSendOption() protocol.ResponseSendOption {
	return info.Option
}
func (info RegistrantInfo) String() string {
	return fmt.Sprintf("RegistrantInfo{ID:%s,Option:%s}", info.ID, info.Option.String())
}

type RequestSendOption struct {
	RequestAddr string
	Timestamp   time.Time
}

func (o RequestSendOption) String() string {
	return fmt.Sprintf("RequestSendOption{RequestAddr:%s,Timestamp:%s}", o.RequestAddr, o.Timestamp)
}

type RegistryInfo struct {
	ID         string
	Option     RequestSendOption
	Candidates []protocol.RegistryInfo
}

func (info RegistryInfo) GetRegistryID() string {
	return info.ID
}
func (info RegistryInfo) GetRequestSendOption() protocol.RequestSendOption {
	return info.Option
}
func (info RegistryInfo) GetCandidates() []protocol.RegistryInfo {
	return info.Candidates
}
func (info RegistryInfo) String() string {
	Candidates := ""
	for _, RegistryInfo := range info.Candidates {
		Candidates += RegistryInfo.String() + ","
	}
	return fmt.Sprintf("RegistryInfo{ID:%s,Option:%s,Candidates:[%s]}",
		info.ID, info.Option.String(), Candidates)
}
