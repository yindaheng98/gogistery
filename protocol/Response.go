package protocol

import (
	"fmt"
	"time"
)

//RequestSendOption is the option information for request sending (encoding, encryption, etc).
type RequestSendOption interface {
	String() string
}

//RegistryInfo contains the information for registry.
//It will be send back from registry to registrant within the request.
//It should be implement by user.
type RegistryInfo interface {

	//Returns the unique ID of the registry
	GetRegistryID() string

	//Returns the type of the service
	GetServiceType() string

	//Returns the option when the registrant send the request
	GetRequestSendOption() RequestSendOption

	//Returns the information of candidate registries
	GetCandidates() []RegistryInfo

	String() string
}

//Response is the response that registry send to registrant.
//It contains the information for registry "RegistryInfo", a connection flag "Reject" and a timeout value "Timeout".
type Response struct {
	RegistryInfo RegistryInfo
	Timeout      time.Duration //下一次连接的时间限制
	Reject       bool          //是否拒绝连接
}

//Get the value of connection flag "Reject".
func (r Response) IsReject() bool {
	return r.Reject
}

//Get the value of timeout value "Timeout".
func (r Response) GetTimeout() time.Duration {
	return r.Timeout
}
func (r Response) String() string {
	return fmt.Sprintf("Registry.Response{RegistryInfo:%s,Timeout:%d,Reject:%t}",
		r.RegistryInfo.String(), r.Timeout, r.Reject)
}
